// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

// set contains an implementation of an unordered collection of elements
// in a `map[T]struct{}`
package set

import "github.com/MKuranowski/go-extra-lib/iter"

// Set is a type implementing an unordered collection of elements
// in a map[T]struct{}, for fast membership checking.
//
// Sets, just as maps, can be nil - and this state new elements can't be added to it.\
//
// To make an empty, non-nil set use make(Set[T]).
//
// To make a non-empty set, use map literals:
//
//	numbers := Set[int]{1: {}, 2: {}, 3: {}}
//
// Sets can be iterated with a range loop,
//
//	for elem := range set
//
// Given operation complexity assumes that element access, insertion and removal
// of a map is on average constant.
type Set[T comparable] map[T]struct{}

// Has returns true if the provided element is in the set.
//
// Average complexity: constant
func (s Set[T]) Has(x T) bool {
	_, has := s[x]
	return has
}

// Add ensures given element is in the set.
//
// Average complexity: constant
func (s Set[T]) Add(x T) { s[x] = struct{}{} }

// Remove ensures given element is not in the set.
//
// Average complexity: constant
func (s Set[T]) Remove(x T) { delete(s, x) }

// Len returns the number of elements in the set.
// Shorthand for len(s).
func (s Set[T]) Len() int { return len(s) }

// Clear ensures no elements are presents in the set.
//
// Average complexity: linear
func (s Set[T]) Clear() {
	for elem := range s {
		delete(s, elem)
	}
}

// Clone returns a shallow copy of the set.
//
// Average complexity: linear
func (s Set[T]) Clone() Set[T] {
	if s == nil {
		return nil
	}

	n := make(Set[T], len(s))
	for elem := range s {
		n[elem] = struct{}{}
	}
	return n
}

// Equal returns true if s1 and s2 contain the same elements.
//
// Average complexity: constant if len(s1) != len(s2),
// otherwise linear in therms of len(s1).
func (s1 Set[T]) Equal(s2 Set[T]) bool {
	if len(s1) != len(s2) {
		return false
	}

	for elem := range s1 {
		if _, has := s2[elem]; !has {
			return false
		}
	}
	return true
}

// Union ensures every element from s2 is also present in s1.
//
// Average complexity: linear in terms of len(s2).
func (s1 Set[T]) Union(s2 Set[T]) {
	for elem := range s2 {
		s1[elem] = struct{}{}
	}
}

// Intersection ensures s1 only contains elements also present in s2.
//
// Average complexity: linear in terms of len(s1).
func (s1 Set[T]) Intersection(s2 Set[T]) {
	for elem := range s1 {
		if _, has := s2[elem]; !has {
			delete(s1, elem)
		}
	}
}

// Difference ensures s1 only contains elements not present in s2.
//
// Average complexity: linear in terms of len(s2).
func (s1 Set[T]) Difference(s2 Set[T]) {
	for elem := range s2 {
		delete(s1, elem)
	}
}

// IsDisjoint returns true if s1 and s2 have no elements in common.
//
// Average complexity; linear in terms of len(s1).
func (s1 Set[T]) IsDisjoint(s2 Set[T]) bool {
	for elem := range s1 {
		if _, has := s2[elem]; has {
			return false
		}
	}
	return true
}

// IsSubset returns true if s2 contains every element from s1.
//
// Average complexity: constant if len(s1) > len(s2),
// otherwise linear in terms of len(s1).
func (s1 Set[T]) IsSubset(s2 Set[T]) bool {
	// short-circuit by the pigeonhole principle
	if len(s1) > len(s2) {
		return false
	}

	for elem := range s1 {
		if _, has := s2[elem]; !has {
			return false
		}
	}
	return true
}

// IsSuperset returns true if s1 contains every element of s2.
//
// Average complexity: constant if len(s2) > len(s1),
// otherwise linear in therms of len(s2).
func (s1 Set[T]) IsSuperset(s2 Set[T]) bool {
	// short-circuit by the pigeonhole principle
	if len(s2) > len(s1) {
		return false
	}

	for elem := range s2 {
		if _, has := s1[elem]; !has {
			return false
		}
	}
	return true
}

// Iter returns an [iter.Iterator] over the elements of the set.
func (s Set[T]) Iter() iter.Iterator[T] { return iter.OverMapKeys(s) }
