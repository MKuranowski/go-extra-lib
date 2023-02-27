// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package iter

import (
	"fmt"
	"math/bits"

	"golang.org/x/exp/slices"
)

type cartesianProductIterator[T any] struct {
	items   [][]T
	indices []int
	dest    []T
	started bool
}

func (i *cartesianProductIterator[T]) Next() bool {
	if !i.started {
		i.started = true
		return true
	}

	for n := len(i.items) - 1; n >= 0; n-- {
		i.indices[n]++
		if n > 0 && i.indices[n] >= len(i.items[n]) {
			i.indices[n] = 0
		} else {
			break
		}
	}

	return i.indices[0] < len(i.items[0])
}

func (i *cartesianProductIterator[T]) Get() []T {
	for n, inner := range i.items {
		i.dest[n] = inner[i.indices[n]]
	}
	return i.dest
}

func (i *cartesianProductIterator[T]) GetCopy() []T { return slices.Clone(i.Get()) }

func (i *cartesianProductIterator[T]) Err() error { return nil }

// CartesianProduct generates the cartesian product of the inner slices.
// Equivalent to nested for loops.
//
// If the outer slice or any of the inner slices are empty, returns an empty iterator.
//
//	CartesianProduct("ABC", "xy") → ["Ax" "Ay" "Bx" "By" "Cx" "Cy"]
//	CartesianProduct("AB", "xy", "12")
//	→ ["Ax1" "Ax2" "Ay1" "Ay2" "Bx1" "Bx2" "By1" "By2" "Cx1" "Cx2" "Cy1" "Cy2"]
//	CartesianProduct("ABC", "") → []
//	CartesianProduct() → []
//
// Subsequent calls to Get() return the same slice, but mutated. See [VolatileIterator].
//
// See [CartesianProductIter], which accepts an iterator over iterators.
//
// The Err() method always returns nil.
func CartesianProduct[T any](outer ...[]T) Iterator[[]T] {
	// Special case for empty input
	if len(outer) == 0 {
		return Empty[[]T]()
	}

	// Special case for empty
	for _, inner := range outer {
		if len(inner) == 0 {
			return Empty[[]T]()
		}
	}

	return &cartesianProductIterator[T]{
		items:   outer,
		indices: make([]int, len(outer)),
		dest:    make([]T, len(outer)),
	}
}

// CartesianProductIter collects the outer iterator and all inner iterators into slices,
// and generates the cartesian product of provided inner iterables.
//
// If the outer iterator or any of the inner slices are empty, returns an empty iterator.
//
//	CartesianProductIter(["ABC", "xy"]) → ["Ax" "Ay" "Bx" "By" "Cx" "Cy"]
//	CartesianProductIter(["AB" "xy" "12"])
//	→ ["Ax1" "Ax2" "Ay1" "Ay2" "Bx1" "Bx2" "By1" "By2" "Cx1" "Cx2" "Cy1" "Cy2"]
//	CartesianProductIter(["ABC" ""]) → []
//	CartesianProductIter([]) → []
//
// Subsequent calls to Get() return the same slice, but mutated. See [VolatileIterator].
//
// See [CartesianProduct], which accepts a slice of slices.
//
// The Err() method always returns nil.
func CartesianProductIter[T any](i Iterator[Iterator[T]]) Iterator[[]T] {
	outer := make([][]T, 0)
	for i.Next() {
		outer = append(outer, IntoSlice(i.Get()))
	}
	return CartesianProduct(outer...)
}

type combinationsIterator[T any] struct {
	items, dest []T
	indices     []int
	n, r        int
	started     bool
}

func (it *combinationsIterator[T]) Next() bool {
	// Based on https://docs.python.org/3/library/itertools.html#itertools.combinations

	if !it.started {
		// initialize iterator fields
		it.n = len(it.items)
		it.indices = make([]int, it.r)
		for i := range it.indices {
			it.indices[i] = i
		}

		it.started = true
		return true
	}

	i := it.r - 1
	for ; i >= 0; i-- {
		if it.indices[i] != i+it.n-it.r {
			break
		}
	}

	if i < 0 {
		return false
	}

	it.indices[i]++
	for j := i + 1; j < it.r; j++ {
		it.indices[j] = it.indices[j-1] + 1
	}
	return true
}

func (i *combinationsIterator[T]) Get() []T {
	for n, index := range i.indices {
		i.dest[n] = i.items[index]
	}
	return i.dest
}

func (i *combinationsIterator[T]) GetCopy() []T { return slices.Clone(i.Get()) }

func (*combinationsIterator[T]) Err() error { return nil }

// Combinations generates all r-length subsequences of provided items, where
// each element may only appear once.
//
// Generated combinations are returned in lexicographical order (as defined by the input ordering).
// Elements are considered unique based on their index.
//
// Panics if r is negative.
//
//	Combinations(2, 'a', 'b', 'c', 'd') → ["ab" "ac" "ad" "bc" "bd" "cd"]
//	Combinations(3, 'a', 'b', 'c', 'd') → ["abc" "abd" "acd" "bcd"]
//	Combinations(0, 'a', 'b', 'c') → [[]]
//	Combinations(0) → [[]]
//	Combinations(3, 'a', 'b') → []
//	Combinations(3) → []
//
// Subsequent calls to Get() return the same slice, but mutated. See [VolatileIterator].
//
// See also [CombinationsIter], which accepts an iterator; or [CombinationsWithReplacement],
// which allows elements to appear multiple times in a subsequence.
// To generate all r-length combinations use [PowerSet].
//
// The Err() method always returns nil.
func Combinations[T any](r int, items ...T) Iterator[[]T] {
	if r < 0 {
		panic(fmt.Sprintf("r can't be negative - got %d", r))
	} else if r > len(items) {
		return Empty[[]T]()
	} else if r == 0 {
		return Over([]T(nil))
	}

	return &combinationsIterator[T]{items: items, dest: make([]T, r), r: r}
}

// CombinationsIter collects the items into a slice and then
// generates all r-length subsequences of provided items, where
// each element may only appear once.
//
// Generated combinations are returned in lexicographical order (as defined by the input ordering).
//
// Panics if r is negative, generates a single empty sequence if r == 0,
// returns an empty sequence if r > len(items).
//
//	CombinationsIter("abcd", 2) → ["ab" "ac" "ad" "bc" "bd" "cd"]
//	CombinationsIter("abcd", 3) → ["abc" "abd" "acd" "bcd"]
//	CombinationsIter("abc", 0) → [[]]
//	CombinationsIter("", 0) → [[]]
//	CombinationsIter("ab", 3) → []
//	CombinationsIter("", 3) → []
//
// Subsequent calls to Get() return the same slice, but mutated. See [VolatileIterator].
//
// See also [Combinations], which accepts a slice of elements; or [CombinationsWithReplacementIter],
// which allows elements to appear multiple times in a subsequence.
// To generate all r-length combinations use [PowerSetIter].
//
// The Err() method always returns nil.
func CombinationsIter[T any](items Iterator[T], r int) Iterator[[]T] {
	return Combinations(r, IntoSlice(items)...)
}

type combinationsWithReplacementIterator[T any] struct {
	items, dest []T
	indices     []int
	n, r        int
	started     bool
}

func (it *combinationsWithReplacementIterator[T]) Next() bool {
	// based on https://docs.python.org/3/library/itertools.html#itertools.combinations_with_replacement

	if !it.started {
		// initialize iterator fields
		it.n = len(it.items)
		it.indices = make([]int, it.r)

		it.started = true
		return true
	}

	i := it.r - 1
	for ; i >= 0; i-- {
		if it.indices[i] != it.n-1 {
			break
		}
	}

	if i < 0 {
		return false
	}

	newIndex := it.indices[i] + 1
	for j := i; j < it.r; j++ {
		it.indices[j] = newIndex
	}

	return true
}

func (i *combinationsWithReplacementIterator[T]) Get() []T {
	for n, index := range i.indices {
		i.dest[n] = i.items[index]
	}
	return i.dest
}

func (i *combinationsWithReplacementIterator[T]) GetCopy() []T { return slices.Clone(i.Get()) }

func (*combinationsWithReplacementIterator[T]) Err() error { return nil }

// CombinationsWithReplacement generates all r-length subsequences of provided items, where
// each element may appear multiple times in the subsequence.
//
// Generated combinations are returned in lexicographical order (as defined by the input ordering).
//
// Panics if r is negative, generates a single empty sequence if r == 0,
// returns an empty sequence if r > len(items).
//
//	CombinationsWithReplacement(2, 'a', 'b', 'c') → ["aa" "ab" "ac" "bb" "bc" "cc"]
//	CombinationsWithReplacement(3, 'a', 'b', 'c')
//	→ ["aaa" "aab" "aac" "abb" "abc" "acc" "bbb" "bbc" "bcc" "ccc"]
//	CombinationsWithReplacement(3, 'a', 'b') → ["aaa" "aab" "abb" "bbb"]
//	CombinationsWithReplacement(0, 'a', 'b') → [[]]
//	CombinationsWithReplacement(0) → [[]]
//	CombinationsWithReplacement(2) → []
//
// Subsequent calls to Get() return the same slice, but mutated. See [VolatileIterator].
//
// See also [CombinationsWithReplacementIter], which accepts an iterator; or [Combinations],
// which does not allow an element to appear multiple times in the subsequence.
//
// The Err() method always returns nil.
func CombinationsWithReplacement[T any](r int, items ...T) Iterator[[]T] {
	if r < 0 {
		panic(fmt.Sprintf("r can't be negative - got %d", r))
	} else if r == 0 {
		return Over([]T(nil))
	} else if len(items) == 0 {
		return Empty[[]T]()
	}

	return &combinationsWithReplacementIterator[T]{items: items, dest: make([]T, r), r: r}
}

// CombinationsWithReplacementIter collects all items into a slice and then
// generates all r-length subsequences of provided items, where
// each element may appear multiple times in the subsequence.
//
// Generated combinations are returned in lexicographical order (as defined by the input ordering).
//
// Panics if r is negative.
//
//	CombinationsWithReplacementIter("abc", 2) → ["aa" "ab" "ac" "bb" "bc" "cc"]
//	CombinationsWithReplacementIter("abc", 3)
//	→ ["aaa" "aab" "aac" "abb" "abc" "acc" "bbb" "bbc" "bcc" "ccc"]
//	CombinationsWithReplacementIter("ab", 3) → ["aaa" "aab" "abb" "bbb"]
//	CombinationsWithReplacementIter("ab", 0) → [[]]
//	CombinationsWithReplacementIter("", 0) → [[]]
//	CombinationsWithReplacementIter("", 2) → []
//
// Subsequent calls to Get() return the same slice, but mutated. See [VolatileIterator].
//
// See also [CombinationsWithReplacement], which accepts a slice; or [CombinationsIter],
// which does not allow an element to appear multiple times in the subsequence.
//
// The Err() method always returns nil.
func CombinationsWithReplacementIter[T any](items Iterator[T], r int) Iterator[[]T] {
	return CombinationsWithReplacement(r, IntoSlice(items)...)
}

type permutationsIterator[T any] struct {
	items, dest     []T
	indices, cycles []int
	n, r            int
	started         bool
}

func (it *permutationsIterator[T]) Next() bool {
	// based on https://docs.python.org/3/library/itertools.html#itertools.permutations

	if !it.started {
		// initialize iterator fields
		it.n = len(it.items)

		it.indices = make([]int, it.n)
		for i := range it.indices {
			it.indices[i] = i
		}

		it.cycles = make([]int, it.r)
		for i := range it.cycles {
			it.cycles[i] = it.n - i
		}

		it.started = true
		return true
	}

	for i := it.r - 1; i >= 0; i-- {
		it.cycles[i]--
		if it.cycles[i] == 0 {
			indexAtI := it.indices[i]
			copy(it.indices[i:], it.indices[i+1:])
			it.indices[it.n-1] = indexAtI

			it.cycles[i] = it.n - i
		} else {
			j := it.n - it.cycles[i]
			it.indices[i], it.indices[j] = it.indices[j], it.indices[i]
			return true
		}

	}

	return false
}

func (i *permutationsIterator[T]) Get() []T {
	for n, index := range i.indices[:i.r] {
		i.dest[n] = i.items[index]
	}
	return i.dest
}

func (i *permutationsIterator[T]) GetCopy() []T { return slices.Clone(i.Get()) }

func (i *permutationsIterator[T]) Err() error { return nil }

// Permutations generates all r-length permutations (different orderings) of the provided items.
//
// Generated combinations are returned in lexicographical order (as defined by the input ordering).
//
// Panics if r is negative, generates a single empty sequence if r == 0,
// returns an empty sequence if r > len(items).
//
//	Permutations(3, 'a', 'b', 'c') → ["abc" "acb" "bac" "bca" "cab" "cba"]
//	Permutations(2, 'a', 'b', 'c') → ["ab" "ac" "ba" "bc" "ca" "cb"]
//	Permutations(0, 'a', 'b', 'c') → [[]]
//	Permutations(0) → [[]]
//	Permutations(4, 'a', 'b', 'c') → []
//	Permutations(4) → []
//
// Subsequent calls to Get() return the same slice, but mutated. See [VolatileIterator].
//
//	See also [PermutationsIter], which accepts an iterator.
//
// The Err() method always returns nil.
func Permutations[T any](r int, items ...T) Iterator[[]T] {
	if r < 0 {
		panic(fmt.Sprintf("r can't be negative - got %d", r))
	} else if r > len(items) {
		return Empty[[]T]()
	} else if r == 0 {
		return Over([]T(nil))
	}

	return &permutationsIterator[T]{items: items, dest: make([]T, r), r: r}
}

// PermutationsIter collects all items into a slice, and then
// generates all r-length permutations (different orderings) of the provided items.
//
// Generated combinations are returned in lexicographical order (as defined by the input ordering).
//
// Panics if r is negative, generates a single empty sequence if r == 0,
// returns an empty sequence if r > len(items).
//
//	PermutationsIter("abc", 3) → ["abc" "acb" "bac" "bca" "cab" "cba"]
//	PermutationsIter("abc", 2) → ["ab" "ac" "ba" "bc" "ca" "cb"]
//	PermutationsIter("abc", 0) → [[]]
//	PermutationsIter(0) → [[]]
//	PermutationsIter("abc", 3) → []
//	PermutationsIter(3) → []
//
// Subsequent calls to Get() return the same slice, but mutated. See [VolatileIterator].
//
//	See also [Permutations], which accepts a slice of elements.
//
// The Err() method always returns nil.
func PermutationsIter[T any](i Iterator[T], r int) Iterator[[]T] {
	return Permutations(r, IntoSlice(i)...)
}

type powerSetIterator[T any] struct {
	items, dest  []T
	current, end uint64
	started      bool
}

func (i *powerSetIterator[T]) Next() bool {
	if i.started {
		i.current += 1
	} else {
		i.started = true
	}
	return i.current < i.end
}

func (i *powerSetIterator[T]) Get() []T {
	n := bits.OnesCount64(i.current)

	// Special case for the empty subset, to avoid doing unnecessary work
	if n == 0 {
		return nil
	}

	destIdx := 0
	for srcIdx, elem := range i.items {
		if i.current>>uint64(srcIdx)&1 != 0 {
			i.dest[destIdx] = elem
			destIdx++
		}
	}

	return i.dest[:n]
}

func (i *powerSetIterator[T]) GetCopy() []T { return slices.Clone(i.Get()) }

func (i *powerSetIterator[T]) Err() error { return nil }

// PowerSet generates all subsets of the provided elements.
//
// Only up to 63 elements are supported - larger power sets are deemed unpractical
// and impossible to iterate - they contain more than quintillion (10^18) elements.
//
//	PowerSet() → [[]]
//	PowerSet(1) → [[] [1]]
//	PowerSet(1, 2) → [[] [1] [2] [1 2]]
//	PowerSet(1, 2, 3) → [[] [1] [2] [1 2] [3] [1 3] [2 3] [1 2 3]]
//
// Subsequent calls to Get() return the same slice, but mutated. See [VolatileIterator].
//
// See [PowerSetIter], which accepts an iterator.
//
// The Err() method always returns nil.
func PowerSet[T any](items ...T) Iterator[[]T] {
	if len(items) > 63 {
		panic(fmt.Sprintf("PowerSet only supports up to 63 elements, got %d", len(items)))
	}

	// Special case for empty input
	if len(items) == 0 {
		return Over([]T(nil))
	}

	return &powerSetIterator[T]{items: items, dest: make([]T, len(items)), end: 1 << len(items)}
}

// PowerSetIter collects all elements into a slice, then generates all subsets of the elements.
//
// Only up to 63 elements are supported - larger power sets are deemed unpractical
// and impossible to iterate - they contain more than quintillion (10^18) elements.
//
//	PowerSetIter([]) → [[]]
//	PowerSetIter([1 2]) → [[] [1] [2] [1 2]]
//	PowerSetIter([1 2 3]) → [[] [1] [2] [1 2] [3] [1 3] [2 3] [1 2 3]]
//
// Subsequent calls to Get() return the same slice, but mutated. See [VolatileIterator].
//
// See [PowerSet], which accepts a slice directly.
//
// The Err() method always returns nil.
func PowerSetIter[T any](i Iterator[T]) Iterator[[]T] {
	return PowerSet(IntoSlice(i)...)
}
