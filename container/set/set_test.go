// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package set_test

import (
	"fmt"
	"testing"

	. "github.com/MKuranowski/go-extra-lib/container/set"
	"github.com/MKuranowski/go-extra-lib/iter"
	"github.com/MKuranowski/go-extra-lib/testing2/check"
	"golang.org/x/exp/slices"
)

func TestSetAddHasLenRemove(t *testing.T) {
	s := Set[int]{}

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

func TestSetClear(t *testing.T) {
	s := Set[int]{2: {}, 3: {}, 5: {}, 7: {}}

	check.EqMsg(t, s.Len(), 4, "s.Len(): after adding")

	s.Clear()

	check.EqMsg(t, s.Len(), 0, "s.Len(): after clearing")
}

func TestSetClone(t *testing.T) {
	s1 := Set[int]{1: {}, 3: {}}
	s2 := s1.Clone()

	s2.Add(5)

	check.EqMsg(t, s1.Len(), 2, "s1.Len()")
	check.EqMsg(t, s2.Len(), 3, "s2.Len()")
}

func TestBitSetEqual(t *testing.T) {
	s1 := Set[int]{2: {}, 3: {}}
	s2 := Set[int]{}

	check.FalseMsg(t, s1.Equal(s2), "{2, 3} == {}")

	s2.Add(2)
	s2.Add(3)

	check.TrueMsg(t, s1.Equal(s2), "{2, 3} == {2, 3}")

	s1.Remove(3)

	check.FalseMsg(t, s1.Equal(s2), "{2} == {2, 3}")
}

func TestSetUnion(t *testing.T) {
	s := Set[int]{1: {}, 3: {}, 5: {}}
	check.EqMsg(t, s.Len(), 3, "s.Len(): before union")

	s.Union(Set[int]{2: {}, 3: {}, 4: {}})
	check.EqMsg(t, s.Len(), 5, "s.Len(): after union")

	for i := 0; i <= 6; i++ {
		if i >= 1 && i <= 5 {
			check.TrueMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		} else {
			check.FalseMsg(t, s.Has(i), fmt.Sprintf("s.Has(%d)", i))
		}
	}
}

func TestSetIntersection(t *testing.T) {
	s := Set[int]{1: {}, 3: {}, 5: {}}
	check.EqMsg(t, s.Len(), 3, "s.Len(): before intersection")

	s.Difference(Set[int]{3: {}, 4: {}, 5: {}})
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
	s := Set[int]{1: {}, 3: {}, 5: {}}
	check.EqMsg(t, s.Len(), 3, "s.Len(): before difference")

	s.Difference(Set[int]{2: {}, 3: {}, 4: {}})
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
		Set[int]{1: {}, 2: {}, 3: {}}.IsDisjoint(Set[int]{1: {}, 2: {}, 3: {}}),
		"Set[int]{1: {}, 2: {}, 3: {}}.IsDisjoined(Set[int]{1: {}, 2: {}, 3: {}})",
	)

	check.TrueMsg(
		t,
		Set[int]{1: {}, 2: {}, 3: {}}.IsDisjoint(Set[int]{4: {}, 5: {}, 6: {}}),
		"Set[int]{1: {}, 2: {}, 3: {}}.IsDisjoined(Set[int]{4: {}, 5: {}, 6: {}})",
	)

	check.FalseMsg(
		t,
		Set[int]{1: {}, 2: {}, 3: {}}.IsDisjoint(Set[int]{3: {}, 4: {}, 5: {}}),
		"Set[int]{1: {}, 2: {}, 3: {}}.IsDisjoined(Set[int]{3: {}, 4: {}, 5: {}})",
	)
}

func TestBitSetIsSubset(t *testing.T) {
	check.TrueMsg(
		t,
		Set[int]{1: {}, 2: {}, 3: {}}.IsSubset(Set[int]{1: {}, 2: {}, 3: {}}),
		"Set[int]{1: {}, 2: {}, 3: {}}.IsSubset(Set[int]{1: {}, 2: {}, 3: {}})",
	)

	check.TrueMsg(
		t,
		Set[int]{1: {}, 2: {}}.IsSubset(Set[int]{1: {}, 2: {}, 3: {}}),
		"Set[int]{1: {}, 2: {}}.IsSubset(Set[int]{1: {}, 2: {}, 3: {}})",
	)

	check.FalseMsg(
		t,
		Set[int]{1: {}, 2: {}, 3: {}}.IsSubset(Set[int]{1: {}, 2: {}}),
		"Set[int]{1: {}, 2: {}, 3: {}}.IsSubset(Set[int]{1: {}, 2: {}})",
	)

	check.FalseMsg(
		t,
		Set[int]{1: {}, 2: {}, 3: {}}.IsSubset(Set[int]{1: {}, 2: {}, 4: {}, 5: {}}),
		"Set[int]{1: {}, 2: {}, 3: {}}.IsSubset(Set[int]{1: {}, 2: {}, 4: {}, 5: {}})",
	)

	check.TrueMsg(
		t,
		Set[int]{}.IsSubset(Set[int]{}),
		"Set[int]{}.IsSubset(Set[int]{})",
	)

	check.TrueMsg(
		t,
		Set[int]{}.IsSubset(Set[int]{1: {}, 2: {}, 3: {}}),
		"Set[int]{}.IsSubset(Set[int]{1: {}, 2: {}, 3: {}})",
	)
}

func TestBitSetIsSuperset(t *testing.T) {
	check.TrueMsg(
		t,
		Set[int]{1: {}, 2: {}, 3: {}}.IsSuperset(Set[int]{1: {}, 2: {}, 3: {}}),
		"Set[int]{1: {}, 2: {}, 3: {}}.IsSuperset(Set[int]{1: {}, 2: {}, 3: {}})",
	)

	check.FalseMsg(
		t,
		Set[int]{1: {}, 2: {}}.IsSuperset(Set[int]{1: {}, 2: {}, 3: {}}),
		"Set[int]{1: {}, 2: {}}.IsSuperset(Set[int]{1: {}, 2: {}, 3: {}})",
	)

	check.TrueMsg(
		t,
		Set[int]{1: {}, 2: {}, 3: {}}.IsSuperset(Set[int]{1: {}, 2: {}}),
		"Set[int]{1: {}, 2: {}, 3: {}}.IsSuperset(Set[int]{1: {}, 2: {}})",
	)

	check.FalseMsg(
		t,
		Set[int]{1: {}, 2: {}, 3: {}}.IsSuperset(Set[int]{1: {}, 2: {}, 4: {}, 5: {}}),
		"Set[int]{1: {}, 2: {}, 3: {}}.IsSuperset(Set[int]{1: {}, 2: {}, 4: {}, 5: {}})",
	)

	check.TrueMsg(
		t,
		Set[int]{}.IsSuperset(Set[int]{}),
		"Set[int]{}.IsSuperset(Set[int]{})",
	)

	check.FalseMsg(
		t,
		Set[int]{}.IsSuperset(Set[int]{1: {}, 2: {}, 3: {}}),
		"Set[int]{}.IsSuperset(Set[int]{1: {}, 2: {}, 3: {}})",
	)
}

func TestBitSetIter(t *testing.T) {
	s := Set[int]{1: {}, 3: {}, 11: {}, 128: {}, 1024: {}}
	sl := iter.IntoSlice(s.Iter())
	slices.Sort(sl) // sort the collected slice to avoid problems with undefined map order

	check.DeepEqMsg(
		t,
		sl,
		[]int{1, 3, 11, 128, 1024},
		"{1, 3, 11, 128, 1024}.Iter()",
	)

	check.DeepEqMsg(t, iter.IntoSlice(Set[int]{}.Iter()), []int{}, "{}.Iter()")
}
