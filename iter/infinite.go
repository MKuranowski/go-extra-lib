// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package iter

import "fmt"

type cycleIterator[T any] struct {
	items          []T
	i, loop, loops int
}

func (i *cycleIterator[T]) Next() bool {
	i.i++
	if i.i == len(i.items) {
		i.loop++
		i.i = 0
	}
	return i.loop < i.loops
}

func (i *cycleIterator[T]) Get() T     { return i.items[i.i] }
func (i *cycleIterator[T]) Err() error { return nil }

// Cycle cycles trough the provided elements `n` times.
//
//	Cycle(3, "a", "b", "c") → ["a" "b" "c" "a" "b" "c" "a" "b" "c"]
//	Cycle(0, "a", "b", "c") → []
//	Cycle(2) → []
//
// Panics if n is negative.
//
// See also [Repeat] which cycles indefinitely and [CycleIter] which accepts an Iterator.
//
// The Err() method always returns nil.
func Cycle[T any](n int, items ...T) Iterator[T] {
	if n < 0 {
		panic(fmt.Sprintf("can't cycle %d (negative) times", n))
	}

	if len(items) == 0 || n == 0 {
		return Empty[T]()
	}
	return &cycleIterator[T]{items: items, i: -1, loops: n}
}

// CycleIter collects elements from the iterator into a slice,
// and then cycles trough elements of the provided iterator `n` times.
//
// Equivalent to `Cycle(n, IntoSlice(i)...)`.
//
//	CycleIter("abc", 3) → "abcabcabc"
//	CycleIter("abc", 0) → ""
//	CycleIter("", 2) → ""
//
// See also [Cycle] and [RepeatIter] which cycles indefinitely.
//
// The Err() method always returns nil.
func CycleIter[T any](i Iterator[T], n int) Iterator[T] {
	return Cycle(n, IntoSlice(i)...)
}

type infiniteRangeIterator[T Numeric] struct {
	current, delta T
	started        bool
}

func (i *infiniteRangeIterator[T]) Next() bool {
	if !i.started {
		i.started = true
	} else {
		i.current += i.delta
	}
	return true
}

func (i *infiniteRangeIterator[T]) Get() T     { return i.current }
func (i *infiniteRangeIterator[T]) Err() error { return nil }

// InfiniteRange generates every number starting from 0, adding 1 on every iteration,
// equivalent to InfiniteRangeWithStep(0, 1).
//
// The iterator performs addition on every step, which may lead to
// accumulating floating-point inaccuracy for non-integer types.
//
//	InfiniteRange[int]() → [0 1 2 ... math.MaxInt math.MinInt ... -1 0 1 ...]
//	InfiniteRange[uint]() → [0 1 2 ... math.MaxUint 0 ...]
//
// See also [InfiniteRangeFrom] and [InfiniteRangeWithStep] and the [Range] family of functions.
//
// The Err() method always returns nil.
func InfiniteRange[T Numeric]() Iterator[T] {
	return InfiniteRangeWithStep(T(0), T(1))
}

// InfiniteRangeFrom generates every number starting from the provided number,
// adding 1 on every iteration. Equivalent to InfiniteRangeWithStep(start, 1).
//
// The iterator performs addition on every step, which may lead to
// accumulating floating-point inaccuracy for non-integer types.
//
//	InfiniteRangeFrom[int](10) → [10 11 12 ... math.MaxInt math.MinInt ... -1 0 1 ...]
//	InfiniteRangeFrom[uint](10) → [10 11 12 ... math.MaxUint 0 1 ...]
//
// See also [InfiniteRange] and [InfiniteRangeWithStep] and the [Range] family of functions.
//
// The Err() method always returns nil.
func InfiniteRangeFrom[T Numeric](start T) Iterator[T] {
	return InfiniteRangeWithStep(start, T(1))
}

// InfiniteRangeWithStep generates every number starting from the `start`
// adding `delta` on every step.
//
// The iterator performs addition on every step, which may lead to
// accumulating floating-point inaccuracy for non-integer types.
//
//	InfiniteRangeWithStep[int](10, 2) → [10 12 14 ... math.MaxInt-1 math.MinInt ... 8 10 12 ...]
//	InfiniteRangeWithStep[int](10, -2) → [10 8 6 ... math.MinInt math.MaxInt-1 ... 12 10 8 ...]
//	InfiniteRangeWithStep[uint](10, 2) → [10 12 14 ... math.MaxUint-1 0 ... 8 10 12 ...]
//
// See also [InfiniteRange] and [InfiniteRangeFrom] and the [Range] family of functions.
//
// The Err() method always returns nil.
func InfiniteRangeWithStep[T Numeric](start, delta T) Iterator[T] {
	return &infiniteRangeIterator[T]{current: start, delta: delta}
}

type rangeIterator[T NumericComparable] struct {
	current, stop, delta T
	started              bool
}

func (i *rangeIterator[T]) Next() bool {
	if !i.started {
		i.started = true
	} else {
		i.current += i.delta
	}
	return i.current < i.stop
}

func (i *rangeIterator[T]) Get() T     { return i.current }
func (i *rangeIterator[T]) Err() error { return nil }

// InfiniteRange generates every number starting from 0, adding 1 on every iteration,
// as long as the current value is smaller than `stop`. Equivalent to RangeWithStep(0, stop, 1).
//
// The iterator performs addition on every step, which may lead to
// accumulating floating-point inaccuracy for non-integer types.
//
//	Range(5) → [0 1 2 3 4]
//	Range(0) → []
//
// See also [RangeFrom] and [RangeWithStep] and the family of [InfiniteRange] functions.
//
// The Err() method always returns nil.
func Range[T NumericComparable](stop T) Iterator[T] {
	return RangeWithStep(T(0), stop, T(1))
}

// RangeFrom generates every number starting from `start`, adding 1 on every iteration,
// as long as the current value is smaller than `stop`. Equivalent to RangeWithStep(start, stop, 1).
//
// The iterator performs addition on every step, which may lead to
// accumulating floating-point inaccuracy for non-integer types.
//
//	Range(5, 10) → [5 6 7 8 9]
//	Range(5, 5) → []
//
// See also [Range] and [RangeWithStep] and the family of [InfiniteRange] functions.
//
// The Err() method always returns nil.
func RangeFrom[T NumericComparable](start, stop T) Iterator[T] {
	return RangeWithStep(start, stop, T(1))
}

// RangeWithStep generates every number starting from `start`, adding `delta` on every iteration,
// as long as the current value is smaller than `stop`.
//
// The iterator performs addition on every step, which may lead to
// accumulating floating-point inaccuracy for non-integer types.
//
//	RangeWithStep(5, 11, 2) → [5 7 9]
//	RangeWithStep(5, 12, 2) → [5 7 9 11]
//	RangeWithStep(5, 5, 2) → [5 7 9 11]
//
// See also [Range] and [RangeFrom] and the family of [InfiniteRange] functions.
//
// The Err() method always returns nil.
func RangeWithStep[T NumericComparable](start, stop, delta T) Iterator[T] {
	return &rangeIterator[T]{current: start, stop: stop, delta: delta}
}

type repeatedlyApplyIterator[T any] struct {
	v       T
	f       func(T) T
	started bool
}

func (i *repeatedlyApplyIterator[T]) Next() bool {
	if !i.started {
		i.started = true
	} else {
		i.v = i.f(i.v)
	}
	return true
}

func (i *repeatedlyApplyIterator[T]) Get() T     { return i.v }
func (i *repeatedlyApplyIterator[T]) Err() error { return nil }

// RepeatedlyApply generates an infinite sequence of continuously
// applying the provided function to its output, starting with `v`.
//
// RepeatedlyApply(x => x + 5, 0) → [0 5 10 15 20 25 ...]
func RepeatedlyApply[T any](f func(T) T, v T) Iterator[T] {
	return &repeatedlyApplyIterator[T]{v: v, f: f}
}

type repeatIterator[T any] struct {
	items []T
	i     int
}

func (i *repeatIterator[T]) Next() bool {
	i.i = (i.i + 1) % len(i.items)
	return true
}

func (i *repeatIterator[T]) Get() T     { return i.items[i.i] }
func (i *repeatIterator[T]) Err() error { return nil }

// Repeat cycles through given elements indefinitely.
//
//	Repeat(1 2 3) → [1 2 3 1 2 3 1 2 3 ...]
//	Repeat(1) → [1 1 1 ...]
//
// Panics if there are no elements.
//
// See also [RepeatIter], which accepts an iterator and [Cycle], which stops after given
// number of cycles.
//
// The Err() method always returns nil.
func Repeat[T any](items ...T) Iterator[T] {
	if len(items) == 0 {
		panic("can't repeat zero elements")
	}
	return &repeatIterator[T]{items: items, i: -1}
}

// RepeatIter collects items from i into a slice, and cycles through them indefinitely.
//
//	RepeatIter([1 2 3]) → [1 2 3 1 2 3 1 2 3 ...]
//	RepeatIter([1]) → [1 1 1 ...]
//
// Equivalent to `Repeat(IntoSlice(i)...)`
//
// Panics if there are no elements.
//
// See also [Repeat] and [CycleIter], which stops after given number of cycles.
//
// The Err() method always returns nil.
func RepeatIter[T any](i Iterator[T]) Iterator[T] {
	return Repeat(IntoSlice(i)...)
}
