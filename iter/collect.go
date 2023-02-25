// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package iter

import "strings"

// IntoSlice collects all elements from an iterator into a single slice.
func IntoSlice[T any](i Iterator[T]) []T {
	s := make([]T, 0)
	for i.Next() {
		s = append(s, i.Get())
	}
	return s
}

// IntoMap collects all elements from an iterator into a map.
func IntoMap[K comparable, V any](i Iterator[Pair[K, V]]) map[K]V {
	m := make(map[K]V)
	for i.Next() {
		elem := i.Get()
		m[elem.First] = elem.Second
	}
	return m
}

// IntoChannel spawns a new goroutine which sends all elements
// from an iterator over a returned channel.
//
// After the iterator is exhausted the returned channel is closed.
func IntoChannel[T any](i Iterator[T]) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for i.Next() {
			ch <- i.Get()
		}
	}()
	return ch
}

// IntoString collects all codepoints and returns a UTF-8 string containing those codepoints.
func IntoString(i Iterator[rune]) string {
	b := strings.Builder{}
	for i.Next() {
		b.WriteRune(i.Get())
	}
	return b.String()
}

// SendOver sends all elements from an iterator over a provided channel.
// Blocks until the iterator is exhausted.
// The provided channel is *not* closed.
func SendOver[T any](i Iterator[T], out chan<- T) {
	for i.Next() {
		out <- i.Get()
	}
}
