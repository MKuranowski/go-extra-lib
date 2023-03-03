// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package bitset_test

import (
	"fmt"
	"testing"

	. "github.com/MKuranowski/go-extra-lib/container/bitset"
	"github.com/MKuranowski/go-extra-lib/iter"
	"github.com/MKuranowski/go-extra-lib/testing2/check"
)

// BitSet

func TestBitSetAddHasLenRemove(t *testing.T) {
	s := &BitSet{}

	check.EqMsg(t, s.Len(), 0, "s.Len(): empty set")

	s.Add(2)
	s.Add(3)
	s.Add(5)
	s.Add(7)

	check.EqMsg(t, s.Len(), 4, "s.Len(): after adding")
	check.FalseMsg(t, s.Has(0), "s.Has(0): after adding")
	check.FalseMsg(t, s.Has(1), "s.Has(1): after adding")
	check.TrueMsg(t, s.Has(2), "s.Has(2): after adding")
	check.TrueMsg(t, s.Has(3), "s.Has(3): after adding")
	check.FalseMsg(t, s.Has(4), "s.Has(4): after adding")
	check.TrueMsg(t, s.Has(5), "s.Has(5): after adding")
	check.FalseMsg(t, s.Has(6), "s.Has(6): after adding")
	check.TrueMsg(t, s.Has(7), "s.Has(7): after adding")

	s.Add(2)

	check.EqMsg(t, s.Len(), 4, "s.Len(): after adding duplicate")
	check.TrueMsg(t, s.Has(2), "s.Has(2): after adding duplicate")

	s.Remove(5)
	s.Remove(7)

	check.EqMsg(t, s.Len(), 2, "s.Len(): after removing")
	check.FalseMsg(t, s.Has(0), "s.Has(0): after removing")
	check.FalseMsg(t, s.Has(1), "s.Has(1): after removing")
	check.TrueMsg(t, s.Has(2), "s.Has(2): after removing")
	check.TrueMsg(t, s.Has(3), "s.Has(3): after removing")
	check.FalseMsg(t, s.Has(4), "s.Has(4): after removing")
	check.FalseMsg(t, s.Has(5), "s.Has(5): after removing")
	check.FalseMsg(t, s.Has(6), "s.Has(6): after removing")
	check.FalseMsg(t, s.Has(7), "s.Has(7): after removing")
}

func TestBitSetOf(t *testing.T) {
	s := Of(1, 3, 5)

	check.EqMsg(t, s.Len(), 3, "s.Len()")
	check.TrueMsg(t, s.Has(1), "s.Has(1)")
	check.TrueMsg(t, s.Has(3), "s.Has(3)")
	check.TrueMsg(t, s.Has(5), "s.Has(5)")
}

func TestBitSetClear(t *testing.T) {
	s := &BitSet{}

	s.Add(2)
	s.Add(3)
	s.Add(5)
	s.Add(7)

	check.EqMsg(t, s.Len(), 4, "s.Len(): after adding")

	s.Clear()

	check.EqMsg(t, s.Len(), 0, "s.Len(): after clearing")
}

func TestBitSetClone(t *testing.T) {
	s1 := &BitSet{}

	s1.Add(1)
	s1.Add(3)

	s2 := s1.Clone()

	s2.Add(5)

	check.EqMsg(t, s1.Len(), 2, "s1.Len()")
	check.EqMsg(t, s2.Len(), 3, "s2.Len()")
}

func TestBitSetEqual(t *testing.T) {
	s1, s2 := new(BitSet), new(BitSet)

	s1.Add(2)
	s1.Add(3)

	check.FalseMsg(t, s1.Equal(s2), "{2, 3} == {}")

	s2.Add(2)
	s2.Add(3)

	check.TrueMsg(t, s1.Equal(s2), "{2, 3} == {2, 3}")

	s1.Remove(3)

	check.FalseMsg(t, s1.Equal(s2), "{2} == {2, 3}")
}

func TestBitSetUnion(t *testing.T) {
	s := Of(1, 3, 5)
	check.EqMsg(t, s.Len(), 3, "s.Len(): before union")

	s.Union(Of(2, 3, 4))
	check.EqMsg(t, s.Len(), 5, "s.Len(): after union")

	for i := 0; i <= 6; i++ {
		if i >= 1 && i <= 5 {
			check.TrueMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		} else {
			check.FalseMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		}
	}
}

func TestBitSetIntersection(t *testing.T) {
	s := Of(1, 3, 5)
	check.EqMsg(t, s.Len(), 3, "s.Len(): before intersection")

	s.Difference(Of(3, 4, 5))
	check.EqMsg(t, s.Len(), 1, "s.Len(): after intersection")

	for i := 0; i <= 6; i++ {
		if i == 1 {
			check.TrueMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		} else {
			check.FalseMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		}
	}
}

func TestBitSetDifference(t *testing.T) {
	s := Of(1, 3, 5)
	check.EqMsg(t, s.Len(), 3, "s.Len(): before difference")

	s.Difference(Of(2, 3, 4))
	check.EqMsg(t, s.Len(), 2, "s.Len(): after difference")

	for i := 0; i <= 6; i++ {
		if i == 1 || i == 5 {
			check.TrueMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		} else {
			check.FalseMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		}
	}
}

func TestBitSetIsDisjoint(t *testing.T) {
	check.FalseMsg(
		t,
		Of(1, 2, 3).IsDisjoint(Of(1, 2, 3)),
		"Of(1, 2, 3).IsDisjoined(Of(1, 2, 3))",
	)

	check.TrueMsg(
		t,
		Of(1, 2, 3).IsDisjoint(Of(4, 5, 6)),
		"Of(1, 2, 3).IsDisjoined(Of(4, 5, 6))",
	)

	check.FalseMsg(
		t,
		Of(1, 2, 3).IsDisjoint(Of(3, 4, 5)),
		"Of(1, 2, 3).IsDisjoined(Of(3, 5, 4))",
	)
}

func TestBitSetIsSubset(t *testing.T) {
	check.TrueMsg(
		t,
		Of(1, 2, 3).IsSubset(Of(1, 2, 3)),
		"Of(1, 2, 3).IsSubset(Of(1, 2, 3))",
	)

	check.TrueMsg(
		t,
		Of(1, 2).IsSubset(Of(1, 2, 3)),
		"Of(1, 2).IsSubset(Of(1, 2, 3))",
	)

	check.FalseMsg(
		t,
		Of(1, 2, 3).IsSubset(Of(1, 2)),
		"Of(1, 2, 3).IsSubset(Of(1, 2))",
	)

	check.FalseMsg(
		t,
		Of(1, 2, 3).IsSubset(Of(1, 2, 4, 5)),
		"Of(1, 2, 3).IsSubset(Of(1, 2, 4, 5))",
	)

	check.TrueMsg(
		t,
		Of().IsSubset(Of()),
		"Of().IsSubset(Of())",
	)

	check.TrueMsg(
		t,
		Of().IsSubset(Of(1, 2, 3)),
		"Of().IsSubset(Of(1, 2, 3))",
	)
}

func TestBitSetIsSuperset(t *testing.T) {
	check.TrueMsg(
		t,
		Of(1, 2, 3).IsSuperset(Of(1, 2, 3)),
		"Of(1, 2, 3).IsSuperset(Of(1, 2, 3))",
	)

	check.FalseMsg(
		t,
		Of(1, 2).IsSuperset(Of(1, 2, 3)),
		"Of(1, 2).IsSuperset(Of(1, 2, 3))",
	)

	check.TrueMsg(
		t,
		Of(1, 2, 3).IsSuperset(Of(1, 2)),
		"Of(1, 2, 3).IsSuperset(Of(1, 2))",
	)

	check.FalseMsg(
		t,
		Of(1, 2, 3).IsSuperset(Of(1, 2, 4, 5)),
		"Of(1, 2, 3).IsSuperset(Of(1, 2, 4, 5))",
	)

	check.TrueMsg(
		t,
		Of().IsSuperset(Of()),
		"Of().IsSuperset(Of())",
	)

	check.FalseMsg(
		t,
		Of().IsSuperset(Of(1, 2, 3)),
		"Of().IsSuperset(Of(1, 2, 3))",
	)
}

func TestBitSetIter(t *testing.T) {
	check.DeepEqMsg(
		t,
		iter.IntoSlice(Of(1, 3, 11, 128, 1024).Iter()),
		[]int{1, 3, 11, 128, 1024},
		"Of(1, 3, 11, 128, 1024).Iter()",
	)

	check.DeepEqMsg(t, iter.IntoSlice((&BitSet{}).Iter()), []int{}, "Of().Iter()")
}

// Small

func TestSmallAddHasLenRemove(t *testing.T) {
	s := Small(0)

	check.EqMsg(t, s.Len(), 0, "s.Len(): empty set")

	s.Add(2)
	s.Add(3)
	s.Add(5)
	s.Add(7)

	check.EqMsg(t, s.Len(), 4, "s.Len(): after adding")
	check.FalseMsg(t, s.Has(0), "s.Has(0): after adding")
	check.FalseMsg(t, s.Has(1), "s.Has(1): after adding")
	check.TrueMsg(t, s.Has(2), "s.Has(2): after adding")
	check.TrueMsg(t, s.Has(3), "s.Has(3): after adding")
	check.FalseMsg(t, s.Has(4), "s.Has(4): after adding")
	check.TrueMsg(t, s.Has(5), "s.Has(5): after adding")
	check.FalseMsg(t, s.Has(6), "s.Has(6): after adding")
	check.TrueMsg(t, s.Has(7), "s.Has(7): after adding")

	s.Add(2)

	check.EqMsg(t, s.Len(), 4, "s.Len(): after adding duplicate")
	check.TrueMsg(t, s.Has(2), "s.Has(2): after adding duplicate")

	s.Remove(5)
	s.Remove(7)

	check.EqMsg(t, s.Len(), 2, "s.Len(): after removing")
	check.FalseMsg(t, s.Has(0), "s.Has(0): after removing")
	check.FalseMsg(t, s.Has(1), "s.Has(1): after removing")
	check.TrueMsg(t, s.Has(2), "s.Has(2): after removing")
	check.TrueMsg(t, s.Has(3), "s.Has(3): after removing")
	check.FalseMsg(t, s.Has(4), "s.Has(4): after removing")
	check.FalseMsg(t, s.Has(5), "s.Has(5): after removing")
	check.FalseMsg(t, s.Has(6), "s.Has(6): after removing")
	check.FalseMsg(t, s.Has(7), "s.Has(7): after removing")
}

func TestSmallOf(t *testing.T) {
	s := SmallOf(1, 3, 5)

	check.EqMsg(t, s.Len(), 3, "s.Len()")
	check.TrueMsg(t, s.Has(1), "s.Has(1)")
	check.TrueMsg(t, s.Has(3), "s.Has(3)")
	check.TrueMsg(t, s.Has(5), "s.Has(5)")
}

func TestSmallClear(t *testing.T) {
	s := Small(0)

	s.Add(2)
	s.Add(3)
	s.Add(5)
	s.Add(7)

	check.EqMsg(t, s.Len(), 4, "s.Len(): after adding")

	s.Clear()

	check.EqMsg(t, s.Len(), 0, "s.Len(): after clearing")
}

func TestSmallClone(t *testing.T) {
	s1 := Small(0)

	s1.Add(1)
	s1.Add(3)

	s2 := s1.Clone()

	s2.Add(5)

	check.EqMsg(t, s1.Len(), 2, "s1.Len()")
	check.EqMsg(t, s2.Len(), 3, "s2.Len()")
}

func TestSmallEqual(t *testing.T) {
	s1, s2 := new(Small), new(Small)

	s1.Add(2)
	s1.Add(3)

	check.FalseMsg(t, s1.Equal(*s2), "{2, 3} == {}")

	s2.Add(2)
	s2.Add(3)

	check.TrueMsg(t, s1.Equal(*s2), "{2, 3} == {2, 3}")

	s1.Remove(3)

	check.FalseMsg(t, s1.Equal(*s2), "{2} == {2, 3}")
}

func TestSmallUnion(t *testing.T) {
	s := SmallOf(1, 3, 5)
	check.EqMsg(t, s.Len(), 3, "s.Len(): before union")

	s.Union(SmallOf(2, 3, 4))
	check.EqMsg(t, s.Len(), 5, "s.Len(): after union")

	for i := 0; i <= 6; i++ {
		if i >= 1 && i <= 5 {
			check.TrueMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		} else {
			check.FalseMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		}
	}
}

func TestSmallIntersection(t *testing.T) {
	s := SmallOf(1, 3, 5)
	check.EqMsg(t, s.Len(), 3, "s.Len(): before intersection")

	s.Difference(SmallOf(3, 4, 5))
	check.EqMsg(t, s.Len(), 1, "s.Len(): after intersection")

	for i := 0; i <= 6; i++ {
		if i == 1 {
			check.TrueMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		} else {
			check.FalseMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		}
	}
}

func TestSmallDifference(t *testing.T) {
	s := SmallOf(1, 3, 5)
	check.EqMsg(t, s.Len(), 3, "s.Len(): before difference")

	s.Difference(SmallOf(2, 3, 4))
	check.EqMsg(t, s.Len(), 2, "s.Len(): after difference")

	for i := 0; i <= 6; i++ {
		if i == 1 || i == 5 {
			check.TrueMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		} else {
			check.FalseMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		}
	}
}

func TestSmallIsDisjoint(t *testing.T) {
	check.FalseMsg(
		t,
		SmallOf(1, 2, 3).IsDisjoint(SmallOf(1, 2, 3)),
		"SmallOf(1, 2, 3).IsDisjoined(SmallOf(1, 2, 3))",
	)

	check.TrueMsg(
		t,
		SmallOf(1, 2, 3).IsDisjoint(SmallOf(4, 5, 6)),
		"SmallOf(1, 2, 3).IsDisjoined(SmallOf(4, 5, 6))",
	)

	check.FalseMsg(
		t,
		SmallOf(1, 2, 3).IsDisjoint(SmallOf(3, 4, 5)),
		"SmallOf(1, 2, 3).IsDisjoined(SmallOf(3, 5, 4))",
	)
}

func TestSmallIsSubset(t *testing.T) {
	check.TrueMsg(
		t,
		SmallOf(1, 2, 3).IsSubset(SmallOf(1, 2, 3)),
		"SmallOf(1, 2, 3).IsSubset(SmallOf(1, 2, 3))",
	)

	check.TrueMsg(
		t,
		SmallOf(1, 2).IsSubset(SmallOf(1, 2, 3)),
		"SmallOf(1, 2).IsSubset(SmallOf(1, 2, 3))",
	)

	check.FalseMsg(
		t,
		SmallOf(1, 2, 3).IsSubset(SmallOf(1, 2)),
		"SmallOf(1, 2, 3).IsSubset(SmallOf(1, 2))",
	)

	check.FalseMsg(
		t,
		SmallOf(1, 2, 3).IsSubset(SmallOf(1, 2, 4, 5)),
		"SmallOf(1, 2, 3).IsSubset(SmallOf(1, 2, 4, 5))",
	)

	check.TrueMsg(
		t,
		SmallOf().IsSubset(SmallOf()),
		"SmallOf().IsSubset(SmallOf())",
	)

	check.TrueMsg(
		t,
		SmallOf().IsSubset(SmallOf(1, 2, 3)),
		"SmallOf().IsSubset(SmallOf(1, 2, 3))",
	)
}

func TestSmallIsSuperset(t *testing.T) {
	check.TrueMsg(
		t,
		SmallOf(1, 2, 3).IsSuperset(SmallOf(1, 2, 3)),
		"SmallOf(1, 2, 3).IsSuperset(SmallOf(1, 2, 3))",
	)

	check.FalseMsg(
		t,
		SmallOf(1, 2).IsSuperset(SmallOf(1, 2, 3)),
		"SmallOf(1, 2).IsSuperset(SmallOf(1, 2, 3))",
	)

	check.TrueMsg(
		t,
		SmallOf(1, 2, 3).IsSuperset(SmallOf(1, 2)),
		"SmallOf(1, 2, 3).IsSuperset(SmallOf(1, 2))",
	)

	check.FalseMsg(
		t,
		SmallOf(1, 2, 3).IsSuperset(SmallOf(1, 2, 4, 5)),
		"SmallOf(1, 2, 3).IsSuperset(SmallOf(1, 2, 4, 5))",
	)

	check.TrueMsg(
		t,
		SmallOf().IsSuperset(SmallOf()),
		"SmallOf().IsSuperset(SmallOf())",
	)

	check.FalseMsg(
		t,
		SmallOf().IsSuperset(SmallOf(1, 2, 3)),
		"SmallOf().IsSuperset(SmallOf(1, 2, 3))",
	)
}

func TestSmallIter(t *testing.T) {
	check.DeepEqMsg(
		t,
		iter.IntoSlice(SmallOf(1, 3, 11, 60).Iter()),
		[]int{1, 3, 11, 60},
		"SmallOf(1, 3, 11, 60).Iter()",
	)

	check.DeepEqMsg(t, iter.IntoSlice(Small(0).Iter()), []int{}, "Of().Iter()")
}
