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

	"github.com/MKuranowski/go-extra-lib/clock"
)

const (
	// Unconditional means that the resource must be fetched regardless if it has changed.
	Unconditional = false

	// Conditional means that the resource must be fetched only if it has changed.
	Conditional = true
)

// Interface represents an external resource, which may change as the program is running.
type Interface interface {
	// Fetch returns the content of the resource if the resource has changed,
	// or the fetch is Unconditional.
	//
	// If an error occurs, (nil, false, err) is returned.
	//
	// For Conditional fetches, the hasChanged return value can be ignored -
	// if the resource changed, returns its content; nil if the resource has not changed.
	//
	// For Unconditional fetches, the content is always non-nil, and the hasChanged
	// flag is set appropriately.
	Fetch(conditional bool) (content io.ReadCloser, hasChanged bool, err error)

	// FetchTime returns the time when the resource was successfully fetched.
	FetchTime() time.Time

	// LastModified returns the time when the external resource was changed as of FetchTime().
	LastModified() time.Time
}

// File is a resource which supports the [fs.File] interface.
//
// Files are considered as changed if their modification time (fs.File.Stat().ModTime())
// has advanced forward.
//
// &File{Open: ...} is ready to use. See also helper [OnFS] and [Local] functions.
type File struct {
	// Open opens the resource file and must behave like [fs.FS]'s Open.
	Open func() (fs.File, error)

	// Clock is the interface used to provide the fetchTime.
	// In nil, [clock.SystemClock] will be used.
	Clock clock.Interface

	fetchTime, lastModified time.Time
}

var _ Interface = &File{}

// Fetch opens the file, stats it and returns it if either unconditionally is set to true,
// or the modification time has advanced.
//
// Returned content, if non-nil, will be exactly what Open() has returned.
func (r *File) Fetch(conditional bool) (content io.ReadCloser, hasChanged bool, err error) {
	// Ensure we have a clock
	if r.Clock == nil {
		r.Clock = clock.System
	}

	// Try to open the file
	f, err := r.Open()
	if err != nil {
		err = fmt.Errorf("resource: Open: %w", err)
		return
	}

	// Ensure file is closed, unless it is returned
	defer func() {
		if content == nil {
			f.Close()
		}
	}()

	// Try to stat the file
	stat, err := f.Stat()
	if err != nil {
		err = fmt.Errorf("resource: Stat: %w", err)
		return
	}

	// Return the file if modification time has advanced,
	// or this is supposed to be an unconditional fetch.
	modTime := stat.ModTime()
	hasChanged = modTime.After(r.lastModified)
	if !conditional || hasChanged {
		r.fetchTime, r.lastModified = r.Clock.Now(), modTime
		content = f
	}

	return
}

// FetchTime returns the latest time when Fetch() returned a non-nil content.
func (r *File) FetchTime() time.Time { return r.fetchTime }

// LastModified returns the modification time of the resource when it was last opened.
func (r *File) LastModified() time.Time { return r.lastModified }

// OnFS creates a [*File] resource which calls fileSystem.Open(name).
func OnFS(fileSystem fs.FS, name string) *File {
	return &File{Open: func() (fs.File, error) { return fileSystem.Open(name) }}
}

// Local creates a [*File] resource which calls [os.Open](name).
func Local(name string) *File {
	return &File{Open: func() (fs.File, error) { return os.Open(name) }}
}

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

	// Clock is the interface used to provide the fetchTime.
	// In nil, clock.SystemClock will be used.
	Clock clock.Interface

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

// Returned by [HTTP.Fetch] if the response is missing the Last-Modified header.
var ErrHTTPNoLastModified = errors.New("server did not return the Last-Modified header")

var _ Interface = &HTTP{} // check that HTTP implements the interface

// Fetch tries to fetch the resource.
//
// If the fetch is Unconditional, makes use of the conditional request headers.
// In this case, the server must reply with 304 Input Not Modified in order
// for this function to detect that the content has not changed.
//
// If the fetch is Conditional, hasChanged will be set if the ETag has changed
// (if the server returned one), or it the Last-Modified time has advanced.
//
// On any errors, including 4xx, 5xx and 3xx status codes, (nil, false, err) is returned.
func (r *HTTP) Fetch(conditional bool) (body io.ReadCloser, hasChanged bool, err error) {
	// Ensure a http.Client is present
	if r.Client == nil {
		r.Client = http.DefaultClient
	}

	// Ensure a clock is present
	if r.Clock == nil {
		r.Clock = clock.System
	}

	// Set the If-None-Match and If-Modified-Since headers
	// If both are set, the If-None-Match takes precedence, according to HTTP.
	// (see https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/If-None-Match)
	//
	// On the first request r.etag will be empty and lasModified will be zero,
	// thus none of the headers will be set and the server will respond with the content.
	if conditional && r.etag != "" {
		r.Request.Header.Set("If-None-Match", r.etag)
	} else {
		r.Request.Header.Del("If-None-Match")
	}

	if conditional && !r.lastModified.IsZero() {
		r.Request.Header.Set("If-Modified-Since", r.lastModified.Format(HTTPTimestampFormat))
	} else {
		r.Request.Header.Del("If-Modified-Since")
	}

	// Run the request
	requestTime := r.Clock.Now()
	resp, err := r.Client.Do(r.Request)
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

	// 304 Input Not Modified - report that nothing has changed;
	// but only for conditional requests.
	if conditional && resp.StatusCode == http.StatusNotModified {
		return
	}

	// Only return the content if the response was successful
	if resp.StatusCode >= 300 {
		err = &HTTPError{r.Request, resp}
		return
	}

	// Try to parse Last-Modified
	lastModifiedString := resp.Header.Get("Last-Modified")
	if lastModifiedString == "" {
		err = ErrHTTPNoLastModified
		return
	}
	lastModified, err := time.Parse(HTTPTimestampFormat, lastModifiedString)
	if err != nil {
		err = fmt.Errorf("invalid Last-Modified: %w", err)
		return
	}

	// Get the etag
	etag := resp.Header.Get("ETag")

	// Try to detect if the resource has changed for unconditional requests.
	// Conditional requests at this point must have been modified,
	// because of an early return if the server returned 304 Input Not Modified.
	if !conditional && etag != "" {
		hasChanged = etag == r.etag
	} else if !conditional {
		hasChanged = lastModified.After(r.lastModified)
	} else {
		hasChanged = true
	}

	r.fetchTime, r.etag, r.lastModified = requestTime, etag, lastModified
	body = resp.Body
	return
}

// FetchTime returns the time of the latest fetch which returned a non-nil body.
func (r *HTTP) FetchTime() time.Time { return r.fetchTime }

// LastModified returns the value of Last-Modified as of FetchTime().
func (r *HTTP) LastModified() time.Time { return r.lastModified }

// ETag returns the value of the ETag header as of FetchTime().
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
// ensures R.Fetch() is called when at least MinimalTimeBetween has passed.
type TimeLimited struct {
	R                  Interface
	MinimalTimeBetween time.Duration

	// Clock is the interface used to provide time - to decide when to fetch
	// forward the fetches to the underlying resource.
	// In nil, [clock.SystemClock] will be used.
	Clock clock.Interface

	nextCheck time.Time
}

var _ Interface = &TimeLimited{} // check that TimeLimited implements the interface

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

// Fetch forwards the call to R.Fetch(conditional) only if it is time to check the resource
// (see [TimeLimited.ShouldCheck]) or the fetch is Unconditional.
func (t *TimeLimited) Fetch(conditional bool) (content io.ReadCloser, hasChanged bool, err error) {
	// Ensure a clock is present
	if t.Clock == nil {
		t.Clock = clock.System
	}

	now := t.Clock.Now()
	if !conditional || now.After(t.nextCheck) {
		t.nextCheck = now.Add(t.MinimalTimeBetween)
		return t.R.Fetch(conditional)
	}
	return
}

// FetchTime returns the last time the resource was actually fetched;
// alias for t.R.FetchTime().
func (t *TimeLimited) FetchTime() time.Time { return t.R.FetchTime() }

// FetchTime returns the last time the resource was last modified (as of FetchTime());
// alias for t.R.LastModified().
func (t *TimeLimited) LastModified() time.Time { return t.R.LastModified() }
