// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

// slices2 is an extension of golang.org/x/exp/slices (https://pkg.go.dev/golang.org/x/exp/slices),
// adding a few more common slice operations, most from https://github.com/golang/go/wiki/SliceTricks.
package slices2

// Batches partitions slice S into ceil(s / batchSize) parts,
// each containing at most batchSize elements.
//
// Example:
//
//	Batches([]int{1, 2, 3, 4, 5}, 2)  // → [[1 2] [3 4] [5]]
//
// Based on https://github.com/golang/go/wiki/SliceTricks#batching-with-minimal-allocation
func Batches[S ~[]E, E any](s S, batchSize int) []S {
	batches := make([]S, 0, (len(s)+batchSize-1)/batchSize)
	for len(s) > batchSize {
		s, batches = s[batchSize:], append(batches, s[0:batchSize:batchSize])
	}
	batches = append(batches, s)
	return batches
}

// DeleteAndSetToZero performs the same operation as slices.Delete (https://pkg.go.dev/golang.org/x/exp/slices#Delete),
// except that deleted elements are set to the zero-value of type E.
//
// Based on https://github.com/golang/go/wiki/SliceTricks#cut
func DeleteAndSetToZero[S ~[]E, E any](s S, i int, j int) S {
	_ = s[i:j] // bounds-check
	var zeroE E
	copy(s[i:], s[j:])
	for k, n := len(s)-j+i, len(s); k < n; k++ {
		s[k] = zeroE
	}
	s = s[:len(s)-j+i]
	return s
}

// Expand inserts n zero-value elements starting at index i
//
// Based on https://github.com/golang/go/wiki/SliceTricks#expand
func Expand[S ~[]E, E any](s S, i, n int) S {
	return append(s[:i], append(make([]E, n), s[i:]...)...)
}

// Extend inserts n zero-value elements at the end of a slice
//
// Based on https://github.com/golang/go/wiki/SliceTricks#expand
func Extend[S ~[]E, E any](s S, n int) S {
	return append(s, make(S, n)...)
}

// Filter modifies in-place a slice by removing elements for which
// keep(x) returns false.
//
// Use FilterAndSetToZero if elements contain pointers to other elements
// to avoid memory leaks.
//
// Based on https://github.com/golang/go/wiki/SliceTricks#filter-in-place
func Filter[S ~[]E, E any](s S, keep func(E) bool) S {
	n := 0
	for _, x := range s {
		if keep(x) {
			s[n] = x
			n++
		}
	}
	return s[:n]
}

// FilterAndSetToZero modifies in-place a slice by removing elements
// for which keep(x) returns false. Removed elements are set to zero.
//
// Based on https://github.com/golang/go/wiki/SliceTricks#filter-in-place
func FilterAndSetToZero[S ~[]E, E any](s S, keep func(E) bool) S {
	n := 0
	for _, x := range s {
		if keep(x) {
			s[n] = x
			n++
		}
	}

	var zeroE E
	for i, end := n, len(s); i < end; i++ {
		s[i] = zeroE
	}

	return s[:n]
}

// Reverse reverses the order of a slice, in-place.
//
// Based on https://github.com/golang/go/wiki/SliceTricks#reversing
func Reverse[E any](s []E) {
	for left, right := 0, len(s)-1; left < right; left, right = left+1, right-1 {
		s[left], s[right] = s[right], s[left]
	}
}

// SlidingWindow returns all slices of s of length windowSize.
//
// If s is smaller than windowSize, returns a single window.
//
// Based on https://github.com/golang/go/wiki/SliceTricks#sliding-window
func SlidingWindow[S ~[]E, E any](s S, windowSize int) []S {
	if len(s) < windowSize {
		return []S{s}
	}

	r := make([]S, 0, len(s)-windowSize+1)
	for i, j, end := 0, windowSize, len(s); j <= end; i, j = i+1, j+1 {
		r = append(r, s[i:j])
	}
	return r
}
