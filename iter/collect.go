// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package iter

import "strings"

// IntoSlice collects all elements from an iterator into a single slice.
//
// If the provided iterator implements [VolatileIterator], uses GetCopy() instead of Get().
func IntoSlice[T any](i Iterator[T]) []T {
	it := ToNonVolatile(i)
	s := make([]T, 0)
	for it.Next() {
		s = append(s, it.Get())
	}
	return s
}

// IntoMap collects all elements from an iterator into a map.
//
// If the provided iterator implements [VolatileIterator], uses GetCopy() instead of Get().
func IntoMap[K comparable, V any](i Iterator[Pair[K, V]]) map[K]V {
	it := ToNonVolatile(i)
	m := make(map[K]V)
	for it.Next() {
		elem := it.Get()
		m[elem.First] = elem.Second
	}
	return m
}

// IntoChannel spawns a new goroutine which sends all elements
// from an iterator over a returned channel.
//
// After the iterator is exhausted the returned channel is closed.
//
// If the provided iterator implements [VolatileIterator], uses GetCopy() instead of Get().
func IntoChannel[T any](i Iterator[T]) <-chan T {
	it := ToNonVolatile(i)
	ch := make(chan T)
	go func() {
		defer close(ch)
		for it.Next() {
			ch <- it.Get()
		}
	}()
	return ch
}

// IntoString collects all codepoints and returns a UTF-8 string containing those codepoints.
//
// If the provided iterator implements [VolatileIterator], uses GetCopy() instead of Get().
func IntoString(i Iterator[rune]) string {
	it := ToNonVolatile(i)
	b := strings.Builder{}
	for it.Next() {
		b.WriteRune(it.Get())
	}
	return b.String()
}

// SendOver sends all elements from an iterator over a provided channel.
// Blocks until the iterator is exhausted.
// The provided channel is *not* closed.
//
// If the provided iterator implements [VolatileIterator], uses GetCopy() instead of Get().
func SendOver[T any](i Iterator[T], out chan<- T) {
	it := ToNonVolatile(i)
	for it.Next() {
		out <- it.Get()
	}
}
