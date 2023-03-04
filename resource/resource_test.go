// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package resource_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MKuranowski/go-extra-lib/clock"
	"github.com/MKuranowski/go-extra-lib/resource"
	"github.com/MKuranowski/go-extra-lib/testing2/assert"
)

const fixtureContent = "Hello, world!\n"

func assertResourceFetched(t *testing.T, r resource.Interface, conditional bool, sequence int) {
	msgPrefix := fmt.Sprintf("FetchIfChanged-%d: ", sequence)

	f, _, err := r.Fetch(conditional)
	assert.NoErrMsg(t, err, msgPrefix+"error")

	if f == nil {
		t.Fatal(msgPrefix + "content: got nil, expected non-nil")
	}

	defer f.Close()

	content, err := io.ReadAll(f)
	assert.NoErrMsg(t, err, msgPrefix+"ReadAll: error")

	assert.EqMsg(t, string(content), fixtureContent, msgPrefix+"content")
}

func assertResourceNotFetched(t *testing.T, r resource.Interface, conditional bool, sequence int) {
	f, _, err := r.Fetch(conditional)
	if err != nil {
		t.Fatalf("FetchIfChanged-%d: error: got: %v, expected: nil", sequence, err)
	}
	if f != nil {
		f.Close()
		t.Fatalf("FetchIfChanged-%d: content: got: %T, expected: nil", sequence, f)
	}
}

func testResource(t *testing.T, r resource.Interface, sleepBeforeChecking time.Duration, refresh func()) {
	refresh()
	assertResourceFetched(t, r, resource.Conditional, 1)

	time.Sleep(sleepBeforeChecking)
	assertResourceNotFetched(t, r, resource.Conditional, 2)

	time.Sleep(sleepBeforeChecking)
	refresh()
	assertResourceFetched(t, r, resource.Conditional, 3)

	time.Sleep(sleepBeforeChecking)
	assertResourceNotFetched(t, r, resource.Conditional, 4)

	time.Sleep(sleepBeforeChecking)
	assertResourceFetched(t, r, resource.Unconditional, 5)
}

func TestLocal(t *testing.T) {
	path := filepath.Join(t.TempDir(), "local_test.txt")
	testResource(
		t,
		resource.Local(path),
		10*time.Millisecond, // a few ms are required for the file system to have a different modification time
		func() {
			err := os.WriteFile(path, []byte(fixtureContent), 0o600)
			if err != nil {
				t.Fatalf("os.WriteFile(fixture): %s", err)
			}
		},
	)
}

func TestOnFS(t *testing.T) {
	tempDir := t.TempDir()
	const name = "on_fs_test.txt"

	fs := os.DirFS(tempDir)
	path := filepath.Join(tempDir, name)

	testResource(
		t,
		resource.OnFS(fs, name),
		10*time.Millisecond, // a few ms are required for the file system to have a different modification time
		func() {
			err := os.WriteFile(path, []byte(fixtureContent), 0o600)
			if err != nil {
				t.Fatalf("os.WriteFile(fixture): %s", err)
			}
		},
	)
}

func TestHTTPLastModified(t *testing.T) {
	// Use a different clock to overcome Last-Modified's resolution of 1 second -
	// otherwise the test takes ages to run.
	c := &clock.EvenlySpaced{
		T:     time.Date(2022, 11, 11, 10, 0, 0, 0, time.UTC),
		Delta: 30 * time.Second,
	}

	refreshTime := c.Now()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Last-Modified", refreshTime.Format(resource.HTTPTimestampFormat))

		ifModifiedSinceString := r.Header.Get("If-Modified-Since")

		// Try to parse the If-Modified-Since and check it
		if ifModifiedSinceString != "" {
			ifModifiedSince, err := time.Parse(resource.HTTPTimestampFormat, ifModifiedSinceString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if !refreshTime.After(ifModifiedSince) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		w.Write([]byte(fixtureContent))
	}))
	defer ts.Close()

	res := resource.HTTPGet(ts.URL)
	res.Clock = c

	testResource(t, res, 0, func() { refreshTime = c.Now().UTC() })
}

func TestHTTPEtag(t *testing.T) {
	refreshTime := time.Now().UTC()
	etagCounter := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		etag := fmt.Sprintf("\"%d\"", etagCounter)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Last-Modified", refreshTime.Format(resource.HTTPTimestampFormat))
		w.Header().Set("ETag", etag)

		// NOTE: this is not proper If-None-Match support; but it's enough to get the test passing
		ifNoneMatch := r.Header.Get("If-None-Match")
		if ifNoneMatch == etag {
			w.WriteHeader(http.StatusNotModified)
		} else {
			w.Write([]byte(fixtureContent))
		}
	}))
	defer ts.Close()

	testResource(
		t,
		resource.HTTPGet(ts.URL),
		0, // ETags are not time-dependent
		func() { refreshTime, etagCounter = time.Now().UTC(), etagCounter+1 },
	)
}

type fixtureResource struct {
	Clock clock.Interface

	fetchTime, lastModified time.Time
}

func (f *fixtureResource) Fetch(conditional bool) (content io.ReadCloser, hasChanged bool, err error) {
	hasChanged = f.lastModified.After(f.fetchTime)
	if !conditional || hasChanged {
		f.fetchTime = f.Clock.Now()
		content = io.NopCloser(strings.NewReader(fixtureContent))
	}
	return
}

func (f *fixtureResource) FetchTime() time.Time    { return f.fetchTime }
func (f *fixtureResource) LastModified() time.Time { return f.lastModified }
func (f *fixtureResource) Refresh()                { f.lastModified = f.Clock.Now() }

func TestTimeLimited(t *testing.T) {
	c := &clock.Specific{Times: []time.Time{
		time.Date(2022, 11, 11, 10, 0, 0, 0, time.UTC),  // 1st call to refresh
		time.Date(2022, 11, 11, 10, 0, 0, 0, time.UTC),  // 1st call to fetch (initial)
		time.Date(2022, 11, 11, 10, 0, 0, 0, time.UTC),  // 1st fetchTime set
		time.Date(2022, 11, 11, 10, 0, 15, 0, time.UTC), // 2nd call to refresh
		time.Date(2022, 11, 11, 10, 0, 30, 0, time.UTC), // 2nd call to fetch (time limited; changed)
		time.Date(2022, 11, 11, 10, 1, 30, 0, time.UTC), // 3rd call to fetch (not limited; changed)
		time.Date(2022, 11, 11, 10, 1, 30, 0, time.UTC), // 2nd fetchTime set
		time.Date(2022, 11, 11, 10, 2, 0, 0, time.UTC),  // 4th call to fetch (time limited; not changed)
		time.Date(2022, 11, 11, 10, 5, 0, 0, time.UTC),  // 5th call to fetch (not limited; not changed)
		time.Date(2022, 11, 11, 10, 6, 0, 0, time.UTC),  // 6th call to fetch (unconditional)
		time.Date(2022, 11, 11, 10, 6, 0, 0, time.UTC),  // 3rd fetchTime set
	}}

	r := fixtureResource{Clock: c}
	tl := &resource.TimeLimited{
		R:                  &r,
		MinimalTimeBetween: time.Minute,
		Clock:              c,
	}

	r.Refresh()
	assertResourceFetched(t, tl, resource.Conditional, 1)

	r.Refresh()
	assertResourceNotFetched(t, tl, resource.Conditional, 2)
	assertResourceFetched(t, tl, resource.Conditional, 3)
	assertResourceNotFetched(t, tl, resource.Conditional, 4)
	assertResourceNotFetched(t, tl, resource.Conditional, 5)
	assertResourceFetched(t, tl, resource.Unconditional, 6)
}
