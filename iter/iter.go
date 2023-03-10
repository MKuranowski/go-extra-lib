// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

// iter is a package for operating on arbitrary collections of elements.
package iter

import (
	"errors"
	"io"
	"reflect"
	"unicode/utf8"
)

// Iterable represents any user-defined struct which can be iterated over.
//
// Slices, maps, channels and strings can also be iterated,
// but in order to get an iterator use one of the OverXxx methods.
type Iterable[T any] interface {
	Iter() Iterator[T]
}

// Iterator represents a state of iteration over elements of type T.
//
// Iterators should start in a state before-the-first-element,
// as iteration is performed in the following pattern:
//
//	for iterator.Next() {
//		elem := iterator.Get()
//	}
//	err := iterator.Err()
type Iterator[T any] interface {
	// Next tries to advance the iterator to the next element,
	// or the first element if it's the first call.
	//
	// Returns true if there's an element available.
	//
	// Must not be called for the second time after false is returned.
	Next() bool

	// Get retrieves the current element of the iterator.
	//
	// Must not be called after Next() returned false.
	//
	// While it is not forbidden to call Get() multiple times without advancing
	// the iterator, is is strongly advised not to do so. Calls to Get()
	// can be arbitrarily complex; see e.g. [Map].
	//
	// All functionalities in the iter module ensure only a single call to Get() is made,
	// as long as the caller also makes a single call to Get().
	//
	// If an Iterator also implements [VolatileIterator], subsequent calls to Get()
	// may return the same element (usually a slice), just mutated.
	Get() T

	// Err returns any error encountered by the iterator.
	//
	// Must not be called before Next() returns false.
	//
	// Some iterator may never error out, and their Err() method may always return nil.
	// See the documentation of relevant iterators.
	Err() error
}

// VolatileIterator is an extension of the Iterator protocol,
// used by iterators, whose return value of Get() is mutated between iterations.
//
// VolatileIterator is usually returned by functions which transform Iterator[T] to Iterator[[]T].
// This allows those iterator to skip allocating a new slice with each call to Get().
//
// Use [ToNonVolatile] if newly-allocated elements are required with each iterator advancement.
type VolatileIterator[T any] interface {
	Iterator[T]

	// GetCopy() return a (usually shallow) copy of the element,
	// which would be returned by Get().
	GetCopy() T
}

type sliceIterator[T any] struct {
	s []T
	i int
}

func (i *sliceIterator[T]) Next() bool {
	if i.i >= len(i.s) {
		return false
	}
	i.i++
	return i.i < len(i.s)
}

func (i *sliceIterator[T]) Get() T {
	return i.s[i.i]
}

func (i *sliceIterator[T]) Err() error {
	return nil
}

// OverSlice returns an iterator over slice elements.
//
// The Err() method always returns nil.
func OverSlice[T any](s []T) Iterator[T] {
	return &sliceIterator[T]{s, -1}
}

// Over returns an iterator over the provided elements.
//
// Equivalent to OverSlice.
//
// The Err() method always returns nil.
func Over[T any](s ...T) Iterator[T] {
	return &sliceIterator[T]{s: s, i: -1}
}

type channelIterator[T any] struct {
	ch <-chan T
	e  T
}

func (i *channelIterator[T]) Next() bool {
	var ok bool
	i.e, ok = <-i.ch
	return ok
}

func (i *channelIterator[T]) Get() T {
	return i.e
}

func (i *channelIterator[T]) Err() error {
	return nil
}

// OverChannel returns an iterator over channel elements.
//
// Note that it may be possible to leak goroutines if the iterator is not exhausted.
// See consumer documentation whether they short-circuit and leave iterators non-exhausted.
// The leak may happen if the producer goroutine(s) assume that all elements
// will be received before exiting, e.g.:
//
//	go func(){
//		for i := 0; i < 100; i++ {
//			ch <- i
//		}
//		close(ch)
//	}()
//
// The Next() method blocks until an element is available.
//
// The Err() method always returns nil.
func OverChannel[T any](ch <-chan T) Iterator[T] {
	return &channelIterator[T]{ch: ch}
}

type mapIterator[K comparable, V any] struct {
	i *reflect.MapIter
}

func (i *mapIterator[K, V]) Next() bool {
	return i.i.Next()
}

func (i *mapIterator[K, V]) Get() Pair[K, V] {
	return Pair[K, V]{
		First:  i.i.Key().Interface().(K),
		Second: i.i.Value().Interface().(V),
	}
}

func (i *mapIterator[K, V]) Err() error {
	return nil
}

// OverMap returns an iterator over key-values pair of a map.
//
// Elements are generated in an arbitrary order.
//
// The Err() method always returns nil.
func OverMap[K comparable, V any](m map[K]V) Iterator[Pair[K, V]] {
	return &mapIterator[K, V]{i: reflect.ValueOf(m).MapRange()}
}

type mapKeyIterator[K comparable] struct {
	i *reflect.MapIter
}

func (i *mapKeyIterator[K]) Next() bool { return i.i.Next() }
func (i *mapKeyIterator[K]) Get() K     { return i.i.Key().Interface().(K) }
func (i *mapKeyIterator[K]) Err() error { return nil }

// OverMapKeys returns an iterator over keys of a map.
//
// Keys are generated in an arbitrary order.
//
// The Err() method always returns nil.
func OverMapKeys[K comparable, V any](m map[K]V) Iterator[K] {
	return &mapKeyIterator[K]{i: reflect.ValueOf(m).MapRange()}
}

type mapValueIterator[V any] struct {
	i *reflect.MapIter
}

func (i *mapValueIterator[V]) Next() bool { return i.i.Next() }
func (i *mapValueIterator[V]) Get() V     { return i.i.Value().Interface().(V) }
func (i *mapValueIterator[V]) Err() error { return nil }

// OverMapValues returns an iterator over values of a map.
//
// Keys are generated in an arbitrary order.
//
// The Err() method always returns nil.
func OverMapValues[K comparable, V any](m map[K]V) Iterator[V] {
	return &mapValueIterator[V]{i: reflect.ValueOf(m).MapRange()}
}

type stringIterator struct {
	rest string
	c    rune
}

func (i *stringIterator) Next() bool {
	if len(i.rest) == 0 {
		return false
	}

	var size int
	i.c, size = utf8.DecodeRuneInString(i.rest)
	i.rest = i.rest[size:]
	return true
}

func (i *stringIterator) Get() rune  { return i.c }
func (i *stringIterator) Err() error { return nil }

// OverString returns an iterator over UTF-8 codepoints in the string.
//
// If the string contains invalid UTF-8 sequences, the replacement character (U+FFFD)
// is returned, and iteration advances over a single byte.
//
// The Err() method always returns nil.
func OverString(s string) Iterator[rune] {
	return &stringIterator{rest: s}
}

type emptyIterator[T any] struct{}

func (emptyIterator[T]) Next() bool { return false }
func (emptyIterator[T]) Get() T     { panic("can't get from an empty iterator") }
func (emptyIterator[T]) Err() error { return nil }

// Empty returns an iterator which never generates any elements,
// and which never returns an error.
func Empty[T any]() Iterator[T] { return emptyIterator[T]{} }

type errorIterator[T any] struct {
	err error
}

func (i errorIterator[T]) Next() bool { return false }
func (i errorIterator[T]) Get() T     { panic("can't get from an error iterator") }
func (i errorIterator[T]) Err() error { return i.err }

// Error returns an iterator, which never generates any elements,
// but whose Err method returns a provided error.
func Error[T any](err error) Iterator[T] { return errorIterator[T]{err} }

type nonVolatileIterator[T any] struct {
	i VolatileIterator[T]
}

func (i nonVolatileIterator[T]) Next() bool { return i.i.Next() }
func (i nonVolatileIterator[T]) Get() T     { return i.i.GetCopy() }
func (i nonVolatileIterator[T]) Err() error { return nil }

// ToNonVolatile ensures that the returned iterator will return newly-allocated
// elements on each call to Get().
//
// See the description of [VolatileIterator].
//
// For VolatileIterator inputs, returns an iterator whose Get() method calls through to GetCopy().
// For other inputs, simply returns the input iterator.
//
// All of the IntoXxx functions automatically call ToNonVolatile.
func ToNonVolatile[T any](i Iterator[T]) Iterator[T] {
	if v, ok := i.(VolatileIterator[T]); ok {
		return nonVolatileIterator[T]{v}
	}
	return i
}

type ioIterator[T any] struct {
	r   IOReader[T]
	v   T
	err error
}

func (i *ioIterator[T]) Next() bool {
	i.v, i.err = i.r.Read()
	if errors.Is(i.err, io.EOF) {
		i.err = nil
		return false
	} else if i.err != nil {
		return false
	}
	return true
}

func (i *ioIterator[T]) Get() T     { return i.v }
func (i *ioIterator[T]) Err() error { return i.err }

// OverIOReader wraps an IOReader into an Iterator.
//
// See [IOReader] for a more detailed explanation;
// but an example implementation of an IOReader is [csv.Reader].
func OverIOReader[T any](r IOReader[T]) Iterator[T] {
	return &ioIterator[T]{r: r}
}
