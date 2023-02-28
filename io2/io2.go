// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

// io2 is an extension of the io module, and mostly contains
// occasionally useful io.Reader implementations.
package io2

import "io"

type repeated struct {
	s             string
	times, offset int
}

func (r *repeated) Read(p []byte) (n int, err error) {
	for {
		if r.times <= 0 || r.s == "" {
			err = io.EOF
			return
		}

		copied := copy(p, r.s[r.offset:])
		n += copied
		r.offset += copied
		p = p[copied:]

		if r.offset == len(r.s) {
			r.offset, r.times = 0, r.times-1
		}

		if len(p) == 0 {
			return
		}
	}
}

// Repeated yields the provided string n number of times.
func Repeated(s string, n int) io.Reader { return &repeated{s: s, times: n} }
