// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

// iter is a package for operating on arbitrary collections of elements.
package iter

import "reflect"

// Iterable represents any user-defined struct which can be iterated over.
//
// Slices, maps and channels can also be iterated,
// but in order to get an iterator call iter.OverSlice, iter.OverMap or iter.OverChannel.
type Iterable[T any] interface {
	Iter() Iterator[T]
}

// Iterator represents a state of iteration over elements of type T.
//
// Calling iterator.Get() or iterator.Next() the second time
// after iterator.Next() returned false is unspecified behavior and will usually panic.
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
	Get() T

	// Err returns any error encountered by the iterator.
	//
	// Must not be called before Next() returns false.
	//
	// Some iterator may never error out, and their Err() method may always return nil.
	// See the documentation of relevant iterators.
	Err() error
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

// Pair is a utility type containing two possibly heterogenous elements.
//
// Use e.g. by OverMap() or Pairwise().
type Pair[T any, U any] struct {
	First  T
	Second U
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

// OverMap returns an iterator over keys of a map.
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

// OverMap returns an iterator over keys of a map.
//
// Keys are generated in an arbitrary order.
//
// The Err() method always returns nil.
func OverMapValues[K comparable, V any](m map[K]V) Iterator[V] {
	return &mapValueIterator[V]{i: reflect.ValueOf(m).MapRange()}
}

type emptyIterator[T any] struct{}

func (i emptyIterator[T]) Next() bool { return false }
func (i emptyIterator[T]) Get() T     { panic("can't get from an empty iterator") }
func (i emptyIterator[T]) Err() error { return nil }

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
