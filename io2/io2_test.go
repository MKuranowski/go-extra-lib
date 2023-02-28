// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package io2_test

import (
	"bytes"
	"io"
	"testing"
	"testing/iotest"

	"github.com/MKuranowski/go-extra-lib/io2"
	"github.com/MKuranowski/go-extra-lib/testing2/check"
)

func TestRepeated(t *testing.T) {
	p, err := io.ReadAll(io2.Repeated("foo", 3))
	check.NoErr(t, err)
	check.EqMsg(t, string(p), "foofoofoo", "Repeated(\"foo\", 3)")

	p, err = io.ReadAll(io2.Repeated("foo", 0))
	check.NoErr(t, err)
	check.EqMsg(t, string(p), "", "Repeated(\"foo\", 0)")

	p, err = io.ReadAll(io2.Repeated("", 3))
	check.NoErr(t, err)
	check.EqMsg(t, string(p), "", "Repeated(\"\", 3)")

	err = iotest.TestReader(
		io2.Repeated("FooBarBaz", 20),
		bytes.Repeat([]byte("FooBarBaz"), 20),
	)
	if err != nil {
		t.Error(err)
	}
}
