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

	"github.com/MKuranowski/go-extra-lib/resource"
	"github.com/MKuranowski/go-extra-lib/testing2/assert"
)

const fixtureContent = "Hello, world!\n"

func assertResourceFetched(t *testing.T, r resource.Interface, sequence int) {
	msgPrefix := fmt.Sprintf("FetchIfChanged-%d: ", sequence)

	f, err := r.FetchIfChanged()
	assert.NoErrMsg(t, err, msgPrefix+"error")

	if f == nil {
		t.Fatal(msgPrefix + "content: got nil, expected non-nil")
	}

	defer f.Close()

	content, err := io.ReadAll(f)
	assert.NoErrMsg(t, err, msgPrefix+"ReadAll: error")

	assert.EqMsg(t, string(content), fixtureContent, msgPrefix+"content")
}

func assertResourceNotFetched(t *testing.T, r resource.Interface, sequence int) {
	f, err := r.FetchIfChanged()
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
	assertResourceFetched(t, r, 1)

	time.Sleep(sleepBeforeChecking)
	assertResourceNotFetched(t, r, 2)

	time.Sleep(sleepBeforeChecking)
	refresh()
	assertResourceFetched(t, r, 3)

	time.Sleep(sleepBeforeChecking)
	assertResourceNotFetched(t, r, 4)
}

func TestLocal(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "local_test.txt")
	testResource(
		t,
		&resource.Local{Path: path},
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
	t.Parallel()

	tempDir := t.TempDir()
	const name = "on_fs_test.txt"

	fs := os.DirFS(tempDir)
	path := filepath.Join(tempDir, name)

	testResource(
		t,
		&resource.OnFS{FS: fs, Name: name},
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
	t.Parallel()

	refreshTime := time.Now().UTC().Truncate(time.Second)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Last-Modified", refreshTime.Format(resource.HTTPTimestampFormat))

		ifModifiedSinceString := r.Header.Get("If-Modified-Since")

		// Try to parse the If-Modified-Since and check
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

	testResource(
		t,
		resource.HTTPGet(ts.URL),
		1005*time.Millisecond, // Last-Modified has a resolution of 1s
		func() { refreshTime = time.Now().UTC().Truncate(time.Second) },
	)
}

func TestHTTPEtag(t *testing.T) {
	t.Parallel()

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
	fetchTime, lastModified time.Time
}

func (f *fixtureResource) FetchIfChanged() (io.ReadCloser, error) {
	if f.lastModified.After(f.fetchTime) {
		f.fetchTime = time.Now()
		return io.NopCloser(strings.NewReader(fixtureContent)), nil
	}
	return nil, nil
}

func (f *fixtureResource) FetchTime() time.Time    { return f.fetchTime }
func (f *fixtureResource) LastModified() time.Time { return f.lastModified }
func (f *fixtureResource) Refresh()                { f.lastModified = time.Now() }

func TestTimeLimited(t *testing.T) {
	r := fixtureResource{}
	r.Refresh()

	tl := &resource.TimeLimited{R: &r, MinimalTimeBetween: time.Millisecond}
	assertResourceFetched(t, tl, 1)

	r.Refresh()
	assertResourceNotFetched(t, tl, 2)

	time.Sleep(time.Millisecond)
	assertResourceFetched(t, tl, 3)
}
