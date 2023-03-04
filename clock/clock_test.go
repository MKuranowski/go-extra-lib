// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package clock_test

import (
	"testing"
	"time"

	"github.com/MKuranowski/go-extra-lib/clock"
)

func checkSameTime(t *testing.T, got, expected time.Time, msg string) {
	if !got.Equal(expected) {
		t.Errorf("%s: got: %v, expected: %v", msg, got, expected)
	}

}

func TestSpecific(t *testing.T) {
	times := []time.Time{
		time.Date(2005, 5, 3, 15, 30, 0, 0, time.UTC),
		time.Date(2005, 5, 3, 15, 31, 0, 0, time.UTC),
		time.Date(2005, 5, 3, 15, 32, 0, 0, time.UTC),
	}
	c := &clock.Specific{Times: times}

	checkSameTime(t, c.Now(), times[0], "1")
	checkSameTime(t, c.Now(), times[1], "2")
	checkSameTime(t, c.Now(), times[2], "3")
}

func TestSpecificWrapAround(t *testing.T) {
	times := []time.Time{
		time.Date(2005, 5, 3, 15, 30, 0, 0, time.UTC),
		time.Date(2005, 5, 3, 15, 31, 0, 0, time.UTC),
		time.Date(2005, 5, 3, 15, 32, 0, 0, time.UTC),
	}
	c := &clock.Specific{Times: times, WrapAround: true}

	checkSameTime(t, c.Now(), times[0], "1")
	checkSameTime(t, c.Now(), times[1], "2")
	checkSameTime(t, c.Now(), times[2], "3")
	checkSameTime(t, c.Now(), times[0], "4")
	checkSameTime(t, c.Now(), times[1], "5")
	checkSameTime(t, c.Now(), times[2], "6")
	checkSameTime(t, c.Now(), times[0], "7")
	checkSameTime(t, c.Now(), times[1], "8")
	checkSameTime(t, c.Now(), times[2], "9")
}

func TestEvenlySpaced(t *testing.T) {
	c := &clock.EvenlySpaced{
		T:     time.Date(2005, 5, 3, 15, 30, 0, 0, time.UTC),
		Delta: time.Minute,
	}

	checkSameTime(t, c.Now(), time.Date(2005, 5, 3, 15, 30, 0, 0, time.UTC), "1")
	checkSameTime(t, c.Now(), time.Date(2005, 5, 3, 15, 31, 0, 0, time.UTC), "2")
	checkSameTime(t, c.Now(), time.Date(2005, 5, 3, 15, 32, 0, 0, time.UTC), "3")
}
