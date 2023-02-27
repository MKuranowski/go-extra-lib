// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package iter

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type accumulateIteratorState uint8

const (
	accumulateIteratorStateBase accumulateIteratorState = iota
	accumulateIteratorStateInitial
	accumulateIteratorStateFinished
)

type accumulateIterator[T, R any] struct {
	i     Iterator[T]
	f     func(accumulator R, element T) R
	state accumulateIteratorState

	acc R
}

func (i *accumulateIterator[T, R]) Next() bool {
	switch i.state {
	case accumulateIteratorStateBase:
		if i.i.Next() {
			i.acc = i.f(i.acc, i.i.Get())
			return true
		} else {
			i.state = accumulateIteratorStateFinished
			return false
		}

	case accumulateIteratorStateInitial:
		i.state = accumulateIteratorStateBase
		return true

	case accumulateIteratorStateFinished:
		return false

	default:
		panic(fmt.Sprintf("invalid accumulateIteratorState state: %d", i.state))
	}
}

func (i *accumulateIterator[T, R]) Get() R { return i.acc }

func (i *accumulateIterator[T, R]) Err() error { return i.i.Err() }

// Accumulate returns an iterator over accumulated ("partial")
// results of applying a binary function.
//
// Assuming that the function provided is summation (`(a, b) => a + b`),
// these are the results of Accumulate:
//
//	Accumulate([1 2 3 4 5], sum) → [1 3 6 10 15]
//	Accumulate([1], sum) → [1]
//	Accumulate([], sum) → []
//
// See function AccumulateWithInitial, which is accepts an initial accumulator value.
//
// See function Reduce, which only returns the last element.
func Accumulate[T any](i Iterator[T], f func(accumulator T, element T) T) Iterator[T] {
	if !i.Next() {
		// empty iterator
		return &emptyIterator[T]{}
	}
	return &accumulateIterator[T, T]{i: i, f: f, acc: i.Get(), state: accumulateIteratorStateInitial}
}

// AccumulateWithInitial returns an iterator over accumulated ("partial")
// results of applying a binary function, starting with the provided initial value.
//
// Assuming that the function provided is summation (`(a, b) => a + b`),
// these are the results of AccumulateWithInitial:
//
//	AccumulateWithInitial([1 2 3 4 5], sum, 0) → [0 1 3 6 10 15]
//	AccumulateWithInitial([1 2 3 4 5], sum, 5) → [5 6 8 11 15 20]
//	AccumulateWithInitial([1], sum, 5) → [5 6]
//	AccumulateWithInitial([], sum, 5) → [5]
//
// See function Accumulate, which is assumes the first element of an iterator
// is the initial accumulator value.
//
// See function ReduceWithInitial, which only returns the last element.
func AccumulateWithInitial[T, R any](i Iterator[T], f func(accumulator R, element T) R, initial R) Iterator[R] {
	return &accumulateIterator[T, R]{i: i, f: f, acc: initial, state: accumulateIteratorStateInitial}
}

// AggregateBy collects elements from an iterable, and groups them by the `key` function.
//
// Similar to [GroupBy], except that this function does work like SQL's GROUP BY construct
// and therefore does not care whether the elements are sorted by the key.
//
//	names := ["Alice" "Andrew" "Bob" "Casey" "Adam" "Amelia" "Chloe" "Craig" "Brian"]
//	AggregateBy(names, name => name[0])
//	→ map["A":["Alice" "Andrew" "Adam" "Amelia"] "B":["Bob" "Brian"] "C":["Casey" "Chloe" "Craig"]]
func AggregateBy[K comparable, V any](i Iterator[V], key func(V) K) map[K][]V {
	r := make(map[K][]V)
	for i.Next() {
		v := i.Get()
		k := key(v)
		r[k] = append(r[k], v)
	}
	return r
}

// Any returns true if any element for the iterator is true.
//
//	Any([false true false]) → true
//	Any([false false]) → false
//	Any([true true]) → true
//	Any([]) → false
//
// See functions All and None; or AnyFunc which accepts objects of arbitrary type.
//
// This function short-circuits and may not exhaust the provided iterator.
func Any(i Iterator[bool]) bool {
	for i.Next() {
		if i.Get() {
			return true
		}
	}
	return false
}

// AnyFunc returns true if there's at least one element for which `f(elem)“ is true.
//
//	type Person struct { Name string; Age int }
//	AnyFunc([Person{"Alice", 30} Person{"Bob", 16}], p => p.Age >= 18) → true
//	AnyFunc([Person{"Bob", 16} Person{"Charlie", 17}], p => p.Age >= 18) → false
//	AnyFunc([Person{"Alice", 30} Person{"Deborah", 47]), p => p.Age >= 18 → true
//	AnyFunc([], p => p.Age >= 18) → false
//
// See functions AllFunc and NoneFunc; or Any which accepts iterators over booleans.
//
// This function short-circuits and may not exhaust the provided iterator.
func AnyFunc[T any](i Iterator[T], f func(T) bool) bool {
	for i.Next() {
		if f(i.Get()) {
			return true
		}
	}
	return false
}

// All returns true if all elements from the iterator are true.
//
//	All([false true false]) → false
//	All([false false]) → false
//	All([true true]) → true
//	All([]) → true
//
// See functions Any and None; or AllFunc which accepts objects of arbitrary type.
//
// This function short-circuits and may not exhaust the provided iterator.
func All(i Iterator[bool]) bool {
	for i.Next() {
		if !i.Get() {
			return false
		}
	}
	return true
}

// All returns true if all for all elements `f(elem)` returns true.
//
//	type Person struct { Name string; Age int }
//	AllFunc([Person{"Alice", 30} Person{"Bob", 16}], p => p.Age >= 18) → false
//	AllFunc([Person{"Bob", 16} Person{"Charlie", 17}], p => p.Age >= 18) → false
//	AllFunc([Person{"Alice", 30} Person{"Deborah", 47]), p => p.Age >= 18 → true
//	AllFunc([], p => p.Age >= 18) → true
//
// See functions AnyFunc and NoneFunc; or Any which accepts iterators over booleans.
//
// This function short-circuits and may not exhaust the provided iterator.
func AllFunc[T any](i Iterator[T], f func(T) bool) bool {
	for i.Next() {
		if !f(i.Get()) {
			return false
		}
	}
	return true
}

// Count exhausts the iterator and returns the number of elements encountered.
//
//	Count([1 2 3]) → 3
//	Count([]) → 0
func Count[T any](i Iterator[T]) int {
	n := 0
	for i.Next() {
		n++
	}
	return n
}

type dropWhileIterator[T any] struct {
	i       Iterator[T]
	pred    func(T) bool
	e       T
	skipped bool
}

func (i *dropWhileIterator[T]) Next() bool {
	for i.i.Next() {
		i.e = i.i.Get()
		if i.skipped {
			return true
		} else if !i.pred(i.e) {
			i.skipped = true
			return true
		}
	}

	return false
}

func (i *dropWhileIterator[T]) Get() T     { return i.e }
func (i *dropWhileIterator[T]) Err() error { return i.i.Err() }

// DropWhile drops the first elements for which `pred(elem)` is true.
// Afterwards, all elements are returned (regardless for the result of pred)
//
//	DropWhile([1 2 3 2 1], x => x < 3) → [3 2 1]
//	DropWhile([1 2 3 2 1], x => x < 5) → []
//	DropWhile([1 2 3 2 1], x => x < 0) → [1 2 3 2 1]
func DropWhile[T any](i Iterator[T], pred func(T) bool) Iterator[T] {
	return &dropWhileIterator[T]{i: i, pred: pred}
}

type enumerateIterator[T any] struct {
	i Iterator[T]
	n int
}

func (i *enumerateIterator[T]) Next() bool {
	has := i.i.Next()
	if has {
		i.n++
	}
	return has
}

func (i *enumerateIterator[T]) Get() Pair[int, T] {
	return Pair[int, T]{i.n, i.i.Get()}
}

func (i *enumerateIterator[T]) Err() error { return i.i.Err() }

// Enumerate generates pairs of elements from i and their corresponding indices
// (offset by start).
//
// Enumerate(["a" "b" "c"], 0) → [Pair{0 "a"} Pair{1 "b"} Pair{2 "c"}]
// Enumerate(["a" "b" "c"], 42) → [Pair{42 "a"} Pair{43 "b"} Pair{44 "c"}]
func Enumerate[T any](i Iterator[T], start int) Iterator[Pair[int, T]] {
	return &enumerateIterator[T]{i, start - 1}
}

// Exhaust exhausts the iterator ignoring any elements in it.
func Exhaust[T any](i Iterator[T]) {
	for i.Next() {
	}
}

type filterIterator[T any] struct {
	i    Iterator[T]
	keep func(T) bool

	e T
}

func (i *filterIterator[T]) Next() bool {
	for i.i.Next() {
		i.e = i.i.Get()
		if i.keep(i.e) {
			return true
		}
	}
	return false
}

func (i *filterIterator[T]) Get() T {
	return i.e
}

func (i *filterIterator[T]) Err() error {
	return i.i.Err()
}

// Filter returns an iterator over elements for which `keep(elem)` returns true.
//
// Filter([1 2 3 4 5 6], isOdd) → [1 3 5]
// Filter([2 4 6], isOdd) → []
func Filter[T any](i Iterator[T], keep func(T) bool) Iterator[T] {
	return &filterIterator[T]{i: i, keep: keep}
}

// ForEach calls the provided function on every element of an iterator, exhausting it.
//
//	ForEach([1 2 3], fmt.Print)
//	 // Prints "123"
func ForEach[T any](i Iterator[T], f func(T)) {
	for i.Next() {
		f(i.Get())
	}
}

// ForEachWithError calls the provided function on every element of an iterator,
// stopping once f returns an error or the iterator is exhausted.
//
// Errors from the iterator are not checked.
//
//	func Foo(i int) error {
//		if i < 0 {
//			return errors.New("i can't be negative")
//		}
//		fmt.Print(i)
//		return nil
//	}
//	i := ForEachWithError([1 -1 2], Foo)
//	// Prints "1"
//	i.Err() → "i can't be negative"
//
// This function short-circuits and may not exhaust the provided iterator.
func ForEachWithError[T any](i Iterator[T], f func(T) error) error {
	for i.Next() {
		err := f(i.Get())
		if err != nil {
			return err
		}
	}
	return nil
}

type limitIterator[T any] struct {
	i    Iterator[T]
	left int
}

func (i *limitIterator[T]) Next() bool {
	if i.i.Next() && i.left > 0 {
		i.left--
		return true
	}
	return false
}

func (i *limitIterator[T]) Get() T     { return i.i.Get() }
func (i *limitIterator[T]) Err() error { return i.i.Err() }

// Limit generates up to n first elements from the provided iterator.
//
// Panics if n is negative.
//
//	Limit([1 2 3 4 5], 3) → [1 2 3]
//	Limit([1 2 3], 5) → [1 2 3]
//	Limit([1 2 3], 0) → []
//
// See functions Skip and Slice.
//
// This function short-circuits and may not exhaust the provided iterator.
func Limit[T any](i Iterator[T], n int) Iterator[T] {
	if n < 0 {
		panic("Limit count can't be negative")
	}
	return &limitIterator[T]{i: i, left: n}
}

type functionMapIterator[T, U any] struct {
	i Iterator[T]
	f func(T) U
}

func (i *functionMapIterator[T, U]) Next() bool {
	return i.i.Next()
}

func (i *functionMapIterator[T, U]) Get() U {
	return i.f(i.i.Get())
}

func (i *functionMapIterator[T, U]) Err() error { return i.i.Err() }

// Map generates the results of applying a function to every element of an iterable.
//
// Every call to Get() results in a call to `f` - which might pose a problem
// if `f` may have side effects. All functionality from the iter module
// guarantees a single call to Get() if the caller also performs only a single call to Get()
// per iteration.
//
//	Map([1 2 3], x => x + 5) → [6 7 8]
func Map[T, U any](i Iterator[T], f func(T) U) Iterator[U] {
	return &functionMapIterator[T, U]{i, f}
}

type functionMapWithErrorIterator[T, U any] struct {
	i Iterator[T]
	f func(T) (U, error)

	e   U
	err error
}

func (i *functionMapWithErrorIterator[T, U]) Next() bool {
	if i.err != nil {
		return false
	} else if i.i.Next() {
		i.e, i.err = i.f(i.i.Get())
		return i.err == nil
	} else {
		i.err = i.i.Err()
		return false
	}
}

func (i *functionMapWithErrorIterator[T, U]) Get() U     { return i.e }
func (i *functionMapWithErrorIterator[T, U]) Err() error { return i.err }

// MapWithError generates the results of applying a function to every element of an iterable,
// stopping once the function returns an error.
//
//	func Foo(i int) (int, error) {
//		if i < 0 {
//			return 0, errors.New("i can't be negative")
//		}
//		return i + 5, nil
//	}
//	MapWithError([1 2 -1 -2 3 4], Foo) → [6 7]
//	// iterator's Err() returns "i can't be negative"
func MapWithError[T, U any](i Iterator[T], f func(T) (U, error)) Iterator[U] {
	return &functionMapWithErrorIterator[T, U]{i: i, f: f}
}

// Min returns the smallest element from the iterator, as by the `<` operator.
//
// If the iterator contains no elements, returns the zero value for T and `ok` is set to false.
//
//	Min([2 5 1 9 3]) → (1, true)
//	Min([]) → (0, false)
func Min[T constraints.Ordered](i Iterator[T]) (min T, ok bool) {
	for i.Next() {
		elem := i.Get()

		if !ok || elem < min {
			min = elem
			ok = true
		}
	}

	return
}

// MinFunc returns the smallest element from the iterator, using less as the comparator.
//
// If the iterator contains no elements, returns the zero value for T and `ok` is set to false.
//
//	type Person struct { Name string; Age int }
//	people := []Person{{"Alice", 30}, {"Bob", 25}, {"Charlie", 41}}
//	MinFunc(people, (p1, p2) => p1.Age < p2.Age) → (Person{"Bob", 25}, true)
//	MinFunc([], (p1, p2) => p1.Age < p2.Age) → (Person{"", 0}, false)
func MinFunc[T any](i Iterator[T], less func(T, T) bool) (min T, ok bool) {
	for i.Next() {
		elem := i.Get()

		if !ok || less(elem, min) {
			min = elem
			ok = true
		}
	}
	return
}

// Max returns the biggest element from the iterator, as by the `>` operator.
//
// If the iterator contains no elements, returns the zero value for T and `ok` is set to false.
//
//	Max([2 5 1 9 3]) → (9, true)
//	Max([]) → (0, false)
func Max[T constraints.Ordered](i Iterator[T]) (max T, ok bool) {
	for i.Next() {
		elem := i.Get()

		if !ok || elem > max {
			max = elem
			ok = true
		}
	}
	return
}

// MaxFunc returns the biggest element from the iterator, using greater as the comparator.
//
// If the iterator contains no elements, returns the zero value for T and `ok` is set to false.
//
//	type Person struct { Name string; Age int }
//	people := []Person{{"Alice", 30}, {"Bob", 25}, {"Charlie", 41}}
//	MaxFunc(people, (p1, p2) => p1.Age > p2.Age) → (Person{"Charlie", 41}, true)
//	MaxFunc([], (p1, p2) => p1.Age > p2.Age) → (Person{"", 0}, false)
func MaxFunc[T any](i Iterator[T], greater func(T, T) bool) (max T, ok bool) {
	for i.Next() {
		elem := i.Get()

		if !ok || greater(elem, max) {
			max = elem
			ok = true
		}
	}
	return
}

// None returns true if all elements from the iterator are false.
//
//	None([false true false]) → false
//	None([false false]) → true
//	None([true true]) → false
//	None([]) → true
//
// See functions Any and None; or NoneFunc which accepts objects of arbitrary type.
//
// This function short-circuits and may not exhaust the provided iterator.
func None(i Iterator[bool]) bool {
	for i.Next() {
		if i.Get() {
			return false
		}
	}
	return true
}

// NoneFunc returns true if all for all elements `f(elem)` returns false.
//
//	type Person struct { Name string; Age int }
//	NoneFunc([Person{"Alice", 30} Person{"Bob", 16}], p => p.Age >= 18) → false
//	NoneFunc([Person{"Bob", 16} Person{"Charlie", 17}], p => p.Age >= 18) → true
//	NoneFunc([Person{"Alice", 30} Person{"Deborah", 47]), p => p.Age >= 18 → false
//	NoneFunc([], p => p.Age >= 18) → true
//
// See functions AnyFunc and AllFunc; or None which accepts iterators over booleans.
//
// This function short-circuits and may not exhaust the provided iterator.
func NoneFunc[T any](i Iterator[T], f func(T) bool) bool {
	for i.Next() {
		if f(i.Get()) {
			return false
		}
	}
	return true
}

// Product returns the result of multiplying all elements in an iterator.
//
// Equivalent to `ReduceWithInitial(i, (a, b) => a * b, 1)`, but removes
// the function call for every element.
//
//	Product([1 2 3]) → 6
//	Product([3+0i 1i]) → 3i
//	Product([]) → 1
func Product[T Numeric](i Iterator[T]) T {
	r := T(1)
	for i.Next() {
		r *= i.Get()
	}
	return r
}

// Reduce returns the result of repeatedly applying a binary function.
//
// The first call to the provided function is done with the 1st and 2nd element
// of the iterable.
//
// If the iterator has only one element, this function simply returns it.
// If the iterator is empty, sets ok to false and returns a zero value for T.
//
// Assuming that the function provided is summation (`(a, b) => a + b`),
// these are the results of Reduce:
//
//	Reduce([1 2 3 4 5], sum) → (15, true)
//	Reduce([1], sum) → (1, true)
//	Reduce([], sum) → (0, false)
//
// See function ReduceWithInitial, which is accepts an initial accumulator value.
//
// See function Accumulate, which also returns partial results.
func Reduce[T any](i Iterator[T], f func(accumulator T, element T) T) (r T, ok bool) {
	for i.Next() {
		if ok {
			r = f(r, i.Get())
		} else {
			r = i.Get()
			ok = true
		}
	}
	return
}

// ReduceWithInitial returns he result of repeatedly applying a binary function,
// starting with the provided initial value.
//
// Assuming that the function provided is summation (`(a, b) => a + b`),
// these are the results of ReduceWithInitial:
//
//	ReduceWithInitial([1 2 3 4 5], sum, 0) → 15
//	ReduceWithInitial([1 2 3 4 5], sum, 5) → 20
//	ReduceWithInitial([1], sum, 5) → 6
//	ReduceWithInitial([], sum, 5) → 5
//
// See function Reduce, which is assumes the first value of an iterator
// is the initial accumulator value.
//
// See function AccumulateWithInitial, which also returns partial results.
func ReduceWithInitial[T, R any](i Iterator[T], f func(accumulator R, element T) R, initial R) R {
	r := initial
	for i.Next() {
		r = f(r, i.Get())
	}
	return r
}

// Skip returns an iterator without the first `n` elements.
//
// Panics if n is negative.
//
//	Skip([1 2 3 4 5], 2) → [3 4 5]
//	Skip([1 2 3], 5) → []
//
// See functions Limit and Slice.
func Skip[T any](i Iterator[T], n int) Iterator[T] {
	if n < 0 {
		panic("Skip count can't be negative")
	}
	for a := 0; a < n; a++ {
		if !i.Next() {
			break
		}
	}
	return i
}

// Slice returns an iterator without the first `start` elements
// and without the last `stop` elements; applying Skip and Limit.
//
// Panics if start or stop are negative; or if start is greater than stop.
//
//	Slice([1 2 3 4 5], 1, 3) → [2 3]
//	Slice([1 2 3 4 5], 0, 3) → [1 2 3] // Use Limit(..., 3)
//	Slice([1 2 3 4 5], 2, 5) → [3 4 5] // Use Skip(..., 2)
//	Slice([1 2 3 4 5], 3, 3) → [] // Use Empty[T]()
//
// See functions Limit and Skip.
//
// This function short-circuits and may not exhaust the provided iterator.
func Slice[T any](i Iterator[T], start, stop int) Iterator[T] {
	if start < 0 || stop < 0 || start > stop {
		panic(fmt.Sprintf("invalid slice: [%d:%d]", start, stop))
	}

	if start > 0 {
		Skip(i, start)
	}
	return Limit(i, stop-start)
}

// Sum returns the result of adding all elements of the iterable.
//
// Equivalent to `ReduceWithInitial(i, (a, b) => a + b, 0)`, but removes
// the function call for every element.
//
//	Sum([1 2 3]) → 6
//	Sum([3+0i 1i]) → 3+1i
//	Sum([]) → 0
func Sum[T Numeric](i Iterator[T]) T {
	var r T
	for i.Next() {
		r += i.Get()
	}
	return r
}

type takeWhileIterator[T any] struct {
	i    Iterator[T]
	pred func(T) bool
	e    T

	done bool
}

func (i *takeWhileIterator[T]) Next() bool {
	if i.done || !i.i.Next() {
		return false
	}

	i.e = i.i.Get()
	if i.pred(i.e) {
		return true
	} else {
		i.done = true
		return false
	}
}

func (i *takeWhileIterator[T]) Get() T     { return i.e }
func (i *takeWhileIterator[T]) Err() error { return i.i.Err() }

// TakeWhile returns the first elements for which `pred(elem)` is true.
// Afterwards, all elements are ignored (regardless for the result of pred).
//
//	TakeWhile([1 2 3 2 1], x => x < 3) → [1 2]
//	TakeWhile([1 2 3 2 1], x => x < 5) → [1 2 3 2 1]
//	TakeWhile([3 2 1], x => x < 3) → []
//
// This function short-circuits and may not exhaust the provided iterator.
func TakeWhile[T any](i Iterator[T], pred func(T) bool) Iterator[T] {
	return &takeWhileIterator[T]{i: i, pred: pred}
}
