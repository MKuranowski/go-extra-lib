// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

// resource is a package for working with external data
// which may be updated as the program is running.
package resource

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Interface represents an external resource, which may change as the program is running.
type Interface interface {
	// FetchIfChanged first checks whether the underlying resource has changed.
	//
	// If it did, opens the resource and returns it, alongside a nil error.
	// If it did not, returns (nil, nil).
	// If any error has occurred, returns nil and that error.
	FetchIfChanged() (io.ReadCloser, error)

	// FetchTime returns the time when the resource was actually opened.
	FetchTime() time.Time

	// LastModified returns the time when the external resource was changed as of FetchTime().
	LastModified() time.Time
}

// OnFS is a resource that lives on an [fs.FS].
// Files are consider as changed if their modification time advances.
//
// &OnFS{FS: fs, Name: name} is ready to be used.
type OnFS struct {
	// FS is the filesystem on which the resource lives.
	FS fs.FS

	// Name is the argument passed to FS.Open when checking
	// if a resource has changed and for returning its content.
	Name string

	fetchTime    time.Time
	lastModified time.Time
}

var _ Interface = &OnFS{} // check that OnFS implements the interface

// FetchIfChanged returns the opened file on the filesystem,
// if that file's modification time has advanced forward.
//
// If the modification time stayed the same (or moved backwards), returns (nil, nil).
//
// FetchTime and LastModified are updated accordingly.
func (r *OnFS) FetchIfChanged() (content io.ReadCloser, err error) {
	f, err := r.FS.Open(r.Name)
	if err != nil {
		err = fmt.Errorf("resource: Open: %w", err)
		return
	}

	// Ensure the file is closed, unless it is returned - even on panic
	defer func() {
		if content == nil {
			f.Close()
		}
	}()

	// Check the modified time of the file
	stat, err := f.Stat()
	if err != nil {
		err = fmt.Errorf("resource: Stat: %w", err)
		return
	}

	// Compare the modification time - if it has advanced, return the content
	currentModificationTime := stat.ModTime()
	if currentModificationTime.After(r.lastModified) {
		r.lastModified = currentModificationTime
		r.fetchTime = time.Now()
		content = f // assign file to content so that defer won't close it
		return
	}

	return
}

// FetchTime returns the time when the resource was actually opened.
func (r *OnFS) FetchTime() time.Time { return r.fetchTime }

// LastModified returns the modification time of the resource when it was last opened.
func (r *OnFS) LastModified() time.Time { return r.lastModified }

// Local is a resource that lives on the local file system.
// Files are consider as changed if their modification time advances.
//
// Uses [os.Open] internally.
//
// &Local{Path: path} is ready to use.
type Local struct {
	Path string

	fetchTime    time.Time
	lastModified time.Time
}

var _ Interface = &Local{} // check that Local implements the interface

// FetchIfChanged returns the opened [*os.File], if
// the file's modification time has advanced.
//
// If the modification time stayed the same (or moved backwards), returns (nil, nil).
//
// FetchTime and LastModified are updated accordingly.
func (r *Local) FetchIfChanged() (content io.ReadCloser, err error) {
	f, err := os.Open(r.Path)
	if err != nil {
		err = fmt.Errorf("resource: Open: %w", err)
		return
	}

	// Ensure the file is closed, unless it is returned - even on panic
	defer func() {
		if content == nil {
			f.Close()
		}
	}()

	// Check the modified time of the file
	stat, err := f.Stat()
	if err != nil {
		err = fmt.Errorf("resource: Stat: %w", err)
		return
	}

	// Compare the modification time - if it has advanced, return the content
	currentModificationTime := stat.ModTime()
	if currentModificationTime.After(r.lastModified) {
		r.lastModified = currentModificationTime
		r.fetchTime = time.Now()
		content = f // assign file to content so that defer won't close it
		return
	}

	return
}

// FetchTime returns the time when the resource was actually opened.
func (r *Local) FetchTime() time.Time { return r.fetchTime }

// LastModified returns the modification time of the file when it was last opened.
func (r *Local) LastModified() time.Time { return r.lastModified }

// HTTP is a resource which lives on a remote HTTP or a HTTPS server.
//
// If the server uses [ETags], the resource is considered as changed if their [ETag] changes.
// Otherwise, the resource is considered as changed if [Last-Modified] has advanced.
//
// The server must include the [Last-Modified] header and must support the
// [If-Modified-Since] (and [If-None-Match] if using [ETags]).
//
// `&HTTP{Request: request}` is ready to use.
//
// [ETag]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/ETag
// [ETags]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/ETag
// [Last-Modified]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Last-Modified
// [If-Modified-Since]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/If-Modified-Since
// [If-None-Match]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/If-None-Match
type HTTP struct {
	// Request is sent to the server on every call to the FetchIfChanged.
	// The request object is reused between calls - and if the request is supposed to have
	// some body, the GetBody field must be used; not Body.
	//
	// FetchIfChanged will modify 2 headers in the request: If-Modified-Since and If-None-Match.
	Request *http.Request

	// Client is used the object used to actually perform the request.
	// If nil, [http.DefaultClient] will be used.
	Client *http.Client

	fetchTime    time.Time
	lastModified time.Time
	etag         string
}

// HTTPError is an error returned when a HTTP server returns an unsuccessful response;
// usually one with [status codes] 3xx, 4xx or 5xx.
//
// [status codes]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
type HTTPError struct {
	Request  *http.Request
	Response *http.Response
}

func (h HTTPError) Error() string {
	return fmt.Sprintf("%s: %s", h.Request.Host, h.Response.Status)
}

// Time format expected by e.g. the Last-Modified or If-Modified-Since headers
const HTTPTimestampFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

// Returned by [HTTP.FetchIfChanged] if the response is missing the Last-Modified header.
var ErrHTTPNoLastModified = errors.New("server did not return the Last-Modified header")

var _ Interface = &HTTP{} // check that HTTP implements the interface

// FetchIfChanged tries to make a conditional request to the server,
// using the If-None-Match and If-Modified-Since headers.
//
// If the server returned 304 Not Modified - returns (nil, nil)
// Otherwise, if the server returned 1xx or 2xx - returns (response.Body, nil).
//
// On any errors, including 4xx, 5xx and other 3xx status codes;
// nil is returned alongside the raised error.
func (r *HTTP) FetchIfChanged() (body io.ReadCloser, err error) {
	// Decide which http.Client to use
	c := r.Client
	if c == nil {
		c = http.DefaultClient
	}

	// Set the If-None-Match and If-Modified-Since headers
	// If both are set, the If-None-Match takes precedence, according to HTTP.
	// (see https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/If-None-Match)
	//
	// On the first request r.etag will be empty and lasModified will be zero,
	// thus none of the headers will be set and the server will respond with the content.
	if r.etag != "" {
		r.Request.Header.Set("If-None-Match", r.etag)
	} else {
		r.Request.Header.Del("If-None-Match")
	}

	if !r.lastModified.IsZero() {
		r.Request.Header.Set("If-Modified-Since", r.lastModified.Format(HTTPTimestampFormat))
	} else {
		r.Request.Header.Del("If-Modified-Since")
	}

	// Run the request
	requestTime := time.Now()
	resp, err := c.Do(r.Request)
	if err != nil {
		err = fmt.Errorf("resource: Do request: %w", err)
		return
	}

	// Ensure the body is closed, unless it is returned - even on panic
	defer func() {
		if body == nil {
			resp.Body.Close()
		}
	}()

	// 304 Input Not Modified - report that nothing has changed
	if resp.StatusCode == http.StatusNotModified {
		return
	}

	// Only return the content if the response was successful (2xx)
	if resp.StatusCode >= 300 {
		err = &HTTPError{r.Request, resp}
		return
	}

	// Try to parse Last-Modified
	currentLastModifiedString := resp.Header.Get("Last-Modified")
	if currentLastModifiedString == "" {
		err = ErrHTTPNoLastModified
		return
	}
	r.lastModified, err = time.Parse(HTTPTimestampFormat, currentLastModifiedString)
	if err != nil {
		err = fmt.Errorf("invalid Last-Modified: %w", err)
		return
	}

	r.etag = resp.Header.Get("ETag")
	r.fetchTime = requestTime
	body = resp.Body
	return
}

// FetchTime returns the last time when the resource was successfully fetched.
func (r *HTTP) FetchTime() time.Time { return r.fetchTime }

// LastModified returns the value of Last-Modified as of last successful fetch.
func (r *HTTP) LastModified() time.Time { return r.lastModified }

// ETag returns the value of the ETag header as of last successful fetch.
func (r *HTTP) ETag() string { return r.etag }

// HTTPGet creates a simple HTTP resource performing GET requests to the specified URL
// using [http.DefaultClient].
//
// The Request and Client may be further customized.
//
// Panics if [http.NewRequest] fails.
func HTTPGet(url string) *HTTP {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Errorf("http.NewRequest failed for %s: %w", url, err))
	}
	return &HTTP{Request: req}
}

// HTTPGetUrl creates a simple HTTP resource performing GET requests to the specified [*url.URL]
// using [http.DefaultClient].
//
// The Request and Client may be further customized.
//
// Panics if [http.NewRequest] fails.
func HTTPGetURL(url *url.URL) *HTTP {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		panic(fmt.Errorf("http.NewRequest failed for https://example.com: %w", err))
	}
	req.URL, req.Host = url, url.Host
	return &HTTP{Request: req}
}

// HTTPPost creates a simple HTTP resource performing POST requests to the specified URL
// using [http.DefaultClient]. Request.GetBody must be set for the body to be properly
// replied with every call to [HTTP.FetchIfChanged]; it is automatically set
// for [*bytes.Buffer], [*bytes.Reader] and [*strings.Reader] (by [http.NewRequestWithContext]).
//
// The Request and Client may be further customized.
//
// Panics if [http.NewRequest] fails.
func HTTPPost(url string, contentType string, body io.Reader) *HTTP {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		panic(fmt.Errorf("http.NewRequest failed for %s: %w", url, err))
	}
	req.Header.Set("Content-Type", contentType)
	return &HTTP{Request: req}
}

// HTTPPost creates a simple HTTP resource performing POST requests to the specified [*url.URL]
// using [http.DefaultClient]. Request.GetBody must be set for the body to be properly
// replied with every call to [HTTP.FetchIfChanged]; it is automatically set
// for [*bytes.Buffer], [*bytes.Reader] and [*strings.Reader] (by [http.NewRequestWithContext]).
//
// The Request and Client may be further customized.
//
// Panics if [http.NewRequest] fails.
func HTTPPostURL(url *url.URL, contentType string, body io.Reader) *HTTP {
	req, err := http.NewRequest("POST", "https://example.com", body)
	if err != nil {
		panic(fmt.Errorf("http.NewRequest failed for https://example.com: %w", err))
	}
	req.URL, req.Host = url, url.Host
	req.Header.Set("Content-Type", contentType)
	return &HTTP{Request: req}
}

// HTTPPostForm creates a simple HTTP resource performing POST requests to the specified url
// with the specified form using [http.DefaultClient].
//
// Content-Type is set to "application/x-www-form-urlencoded" and
// Request.GetBody is set to return the same data every time.
//
// The Request and Client may be further customized.
//
// Panics if [http.NewRequest] fails.
func HTTPPostForm(url string, data url.Values) *HTTP {
	return HTTPPost(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// HTTPPostForm creates a simple HTTP resource performing POST requests to the specified [*url.URL]
// with the specified form using [http.DefaultClient].
//
// Content-Type is set to "application/x-www-form-urlencoded" and
// Request.GetBody is set to return the same data every time.
//
// The Request and Client may be further customized.
//
// Panics if [http.NewRequest] fails.
func HTTPPostFormURL(url *url.URL, data url.Values) *HTTP {
	return HTTPPostURL(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// TimeLimited is a simple rate-limiting mechanizm for Resources -
// ensures R.FetchIfChanged() is called when at least MinimalTimeBetween has passed.
type TimeLimited struct {
	R                  Interface
	MinimalTimeBetween time.Duration

	nextCheck time.Time
}

var _ Interface = &TimeLimited{} // check that TimeLimited implements the interface

// ShouldCheck returns true if sufficient time has passed from last call to
// [Resource.FetchIfChanged] to facilitate another call.
func (t *TimeLimited) ShouldCheck() bool {
	return time.Now().After(t.nextCheck)
}

// NextCheck returns the time when the resource should be checked,
// or a zero-value time.Time if the resource was never checked.
func (t *TimeLimited) NextCheck() time.Time { return t.nextCheck }

// LastCheck returns the time when the resource was last checked,
// or a zero-value time.Time if the resource was never checked.
func (t *TimeLimited) LastCheck() time.Time {
	if t.nextCheck.IsZero() {
		return time.Time{}
	}
	return t.nextCheck.Add(-t.MinimalTimeBetween)
}

// ForceFetchIfChanged bypasses the timer checks and always calls R.FetchIfChanged.
// Also updates the NextCheck() and LastCheck() fields.
func (t *TimeLimited) ForceFetchIfChanged() (io.ReadCloser, error) {
	t.nextCheck = time.Now().Add(t.MinimalTimeBetween)
	return t.R.FetchIfChanged()
}

// FetchIfChanged checks if it is time to check the resource -
// and returns the result of calling R.FetchIfChanged(); or (nil, nil) otherwise.
//
// Shorthand for:
//
//	if t.ShouldCheck() {
//		return t.ForceFetchIfChanged()
//	}
//	return (nil, nil)
func (t *TimeLimited) FetchIfChanged() (io.ReadCloser, error) {
	if t.ShouldCheck() {
		return t.ForceFetchIfChanged()
	}
	return nil, nil
}

// FetchTime returns the last time the resource was actually fetched;
// alias for t.R.FetchTime().
func (t *TimeLimited) FetchTime() time.Time { return t.R.FetchTime() }

// FetchTime returns the last time the resource was last modified (as of FetchTime());
// alias for t.R.LastModified().
func (t *TimeLimited) LastModified() time.Time { return t.R.LastModified() }
