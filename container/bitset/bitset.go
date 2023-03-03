// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

// bitset contains an efficient implementation of a set of unsigned numbers.
package bitset

import (
	"math/big"
	"math/bits"

	"github.com/MKuranowski/go-extra-lib/iter"
)

var (
	bigZero = big.Int{}
)

// BitSet is a set of (almost) arbitrary-sized integers.
//
// The zero value (`&BitSet{}`) is a BitSet containing no elements.
//
// The representation uses [big.Int] to check whether a number is included in the set,
// so a map-based set may be a better use-case for sparse sets without any upper-bound.
//
// Even tough most operations accept `int` as an argument,
// those functions will panic if the provided number is negative.
//
// See also [Small] - a more efficient implementation if elements are known
// to be in range [0, 63] inclusive - and which can be used as a key in a map
// (by fulfilling the comparable protocol).
type BitSet struct {
	n big.Int
}

// Of returns a BitSet containing all the provided elements
func Of(is ...int) *BitSet {
	b := &BitSet{}
	for _, i := range is {
		b.Add(i)
	}
	return b
}

// Has returns true if the provided number is in the set.
func (s *BitSet) Has(i int) bool { return s.n.Bit(i) != 0 }

// Add ensures that the provided number is in the set.
func (s *BitSet) Add(i int) { s.n.SetBit(&s.n, i, 1) }

// Remove ensures that the provided number is not in the set.
func (s *BitSet) Remove(i int) { s.n.SetBit(&s.n, i, 0) }

// Len returns the number of elements in the set.
func (s *BitSet) Len() int {
	n := 0
	for _, word := range s.n.Bits() {
		n += bits.OnesCount(uint(word))
	}
	return n
}

// Clear ensures that no numbers are present in the set.
func (s *BitSet) Clear() { s.n.SetUint64(0) }

// Clone returns a new set with the same elements.
func (s *BitSet) Clone() *BitSet {
	n := &BitSet{}
	n.n.Set(&s.n)
	return n
}

// Equal returns true if s1 contains the same elements as s2.
func (s1 *BitSet) Equal(s2 *BitSet) bool { return s1.n.Cmp(&s2.n) == 0 }

// Union ensures s1 contains all elements from s2.
func (s1 *BitSet) Union(s2 *BitSet) { s1.n.Or(&s1.n, &s2.n) }

// Intersection ensures s1 only contains elements that are present in both s1 and s2.
func (s1 *BitSet) Intersection(s2 *BitSet) { s1.n.And(&s1.n, &s2.n) }

// Difference ensures s1 does not contain any elements from s1.
func (s1 *BitSet) Difference(s2 *BitSet) { s1.n.AndNot(&s1.n, &s2.n) }

// IsDisjoint returns true if s1 and s2 have no elements in common.
func (s1 *BitSet) IsDisjoint(s2 *BitSet) bool {
	return (&big.Int{}).And(&s1.n, &s2.n).Cmp(&bigZero) == 0
}

// IsSubset returns true if every element of s1 is also present in s2.
func (s1 *BitSet) IsSubset(s2 *BitSet) bool {
	return (&big.Int{}).And(&s1.n, &s2.n).Cmp(&s1.n) == 0
}

// IsSuperset returns true if every element of s2 is also present in s1.
func (s1 *BitSet) IsSuperset(s2 *BitSet) bool {
	return (&big.Int{}).And(&s2.n, &s1.n).Cmp(&s2.n) == 0
}

type bitsetIterator struct {
	s       BitSet
	n       int
	started bool
}

func (i *bitsetIterator) Next() bool {
	// Shift out the last-generated element, except if there was no such element
	if i.started {
		i.s.n.Rsh(&i.s.n, 1)
		i.n++
	} else {
		i.started = true
	}

	if i.s.n.Cmp(&bigZero) == 0 {
		return false
	}

	// Calculate the offset to the next number
	offset := i.s.n.TrailingZeroBits()
	i.s.n.Rsh(&i.s.n, offset)
	i.n += int(offset)

	if i.s.n.Bit(0) == 0 {
		panic("big.Int.TrailingZeroBits() has lied")
	}
	return true
}

func (i bitsetIterator) Get() int { return i.n }
func (bitsetIterator) Err() error { return nil }

func (s *BitSet) Iter() iter.Iterator[int] {
	return &bitsetIterator{s: *s}
}

// Small is a set of integers between 0 and 63 (inclusive),
// which also fulfills the `comparable` interface and can be used e.g. as a map key.
//
// The zero value (`Small(0)`) is a set containing no elements.
//
// The representation is a simple wrapper around uint64.
//
// Trying to provide elements outside of the <0, 63> range is undefined behavior.
//
// See also [BitSet] - an implementation for (almost) arbitrary sized elements.
type Small uint64

// SmallOf returns a bitset containing the provided numbers.
func SmallOf(is ...int) Small {
	s := Small(0)
	for _, i := range is {
		s.Add(i)
	}
	return s
}

// Has returns true if the provided number is in the set.
func (s Small) Has(i int) bool { return (s>>Small(i))&1 != 0 }

// Add ensures that the provided number is in the set.
func (s *Small) Add(i int) { *s |= 1 << Small(i) }

// Remove ensures that the provided number is not in the set.
func (s *Small) Remove(i int) { *s &^= 1 << Small(i) }

// Len returns the number of elements in the set.
func (s Small) Len() int { return bits.OnesCount64(uint64(s)) }

// Clear ensures that no numbers are present in the set.
func (s *Small) Clear() { *s = 0 }

// Clone returns a new set with the same elements.
func (s Small) Clone() Small { return s }

// Equal returns true if b1 contains the same elements as b2.
func (s1 Small) Equal(s2 Small) bool { return s1 == s2 }

// Union ensures b1 contains all elements from b2.
func (s1 *Small) Union(s2 Small) { *s1 |= s2 }

// Intersection ensures b1 only contains elements that are present in both b1 and b2.
func (s1 *Small) Intersection(s2 Small) { *s1 &= s2 }

// Difference ensures b1 does not contain any elements from b2.
func (s1 *Small) Difference(s2 Small) { *s1 &^= s2 }

// IsDisjoint returns true if s1 and s2 have no elements in common.
func (s1 Small) IsDisjoint(s2 Small) bool { return s1&s2 == 0 }

// IsSubset returns true if every element of s1 is also present in s2.
func (s1 Small) IsSubset(s2 Small) bool { return s1&s2 == s1 }

// IsSuperset returns true if every element of s2 is also present in s1.
func (s1 Small) IsSuperset(s2 Small) bool { return s2&s1 == s2 }

type smallIterator struct {
	s       uint64
	n       int
	started bool
}

func (i *smallIterator) Next() bool {
	// Shift out the last-generated element, except if there was no such element
	if i.started {
		i.s, i.n = i.s>>1, i.n+1
	} else {
		i.started = true
	}

	if i.s == 0 {
		return false
	}

	// Calculate the offset to the next number
	offset := bits.TrailingZeros64(i.s)
	i.s, i.n = i.s>>offset, i.n+offset

	if i.s&1 == 0 {
		panic("TrailingZeroBits64 has lied")
	}
	return true
}

func (i smallIterator) Get() int { return i.n }
func (smallIterator) Err() error { return nil }

// Iter returns an iterator over the elements in the set.
//
// Any changes made during iteration are not reflected in the iterator;
// iteration is actually performed on a copy of the set.
func (s Small) Iter() iter.Iterator[int] { return &smallIterator{s: uint64(s)} }
