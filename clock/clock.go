// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

// clock is a package for providing time.
package clock

import "time"

// Interface defines a clock, which can provide time.
type Interface interface {
	Now() time.Time
}

// System is a clock which uses time.Now() to provide time.
var System Interface = systemClock{}

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

// Specific is a clock providing specific times, in sequence.
//
// &Specific{Time: ...} is ready to use.
type Specific struct {
	// Times is the slice of times to generate.
	Times []time.Time

	// WrapAround controls the behavior once the Specific clock runs out of times.
	// If false (default), Now() panics. Otherwise, Now() wraps around and starts
	// back at the beginning of the provided slice.
	WrapAround bool

	i int
}

func (s *Specific) Now() time.Time {
	if s.i >= len(s.Times) {
		panic("clock.Specific overflow")
	}
	t := s.Times[s.i]

	s.i++
	if s.WrapAround && s.i == len(s.Times) {
		s.i = 0
	}

	return t
}
