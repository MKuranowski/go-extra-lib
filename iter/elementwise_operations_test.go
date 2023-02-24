// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package iter_test

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	. "github.com/MKuranowski/go-extra-lib/iter"
	"github.com/MKuranowski/go-extra-lib/testing2/check"
)

type person struct {
	name string
	age  int
}

func isOver18(p person) bool     { return p.age >= 18 }
func younger(p1, p2 person) bool { return p1.age < p2.age }
func older(p1, p2 person) bool   { return p1.age > p2.age }

func add(a, b int) int { return a + b }
func isOdd(x int) bool { return x%2 == 1 }

func TestAccumulate(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Accumulate(Over(1, 2, 3, 4, 5), add)),
		[]int{1, 3, 6, 10, 15},
		"Accumulate([1 2 3 4 5], add)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Accumulate(Over(1), add)),
		[]int{1},
		"Accumulate([1], add)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Accumulate(Empty[int](), add)),
		[]int{},
		"Accumulate([], add)",
	)
}

func TestAccumulateWithInitial(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(AccumulateWithInitial(Over(1, 2, 3, 4, 5), add, 0)),
		[]int{0, 1, 3, 6, 10, 15},
		"AccumulateWithInitial([1 2 3 4 5], add, 0)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(AccumulateWithInitial(Over(1, 2, 3, 4, 5), add, 5)),
		[]int{5, 6, 8, 11, 15, 20},
		"AccumulateWithInitial([1 2 3 4 5], add, 5)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(AccumulateWithInitial(Over(1), add, 5)),
		[]int{5, 6},
		"AccumulateWithInitial([1], add, 5)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(AccumulateWithInitial(Empty[int](), add, 5)),
		[]int{5},
		"AccumulateWithInitial([], add, 0)",
	)
}

func TestAny(t *testing.T) {
	check.TrueMsg(
		t,
		Any(Over(false, true, false)),
		"Any([false true false])",
	)

	check.FalseMsg(
		t,
		Any(Over(false, false)),
		"Any([false false])",
	)

	check.TrueMsg(
		t,
		Any(Over(true, true)),
		"Any([true true])",
	)

	check.FalseMsg(
		t,
		Any(Empty[bool]()),
		"Any([])",
	)
}

func TestAnyFunc(t *testing.T) {
	check.TrueMsg(
		t,
		AnyFunc(
			Over(person{"Alice", 30}, person{"Bob", 16}),
			isOver18,
		),
		"AnyFunc([{\"Alice\", 30}, {\"Bob\", 16}], isOver18)",
	)

	check.FalseMsg(
		t,
		AnyFunc(
			Over(person{"Bob", 16}, person{"Charlie", 17}),
			isOver18,
		),
		"AnyFunc([{\"Bob\", 16}, {\"Charlie\", 17}], isOver18)",
	)

	check.TrueMsg(
		t,
		AnyFunc(
			Over(person{"Alice", 30}, person{"Deborah", 47}),
			isOver18,
		),
		"AnyFunc([{\"Alice\", 30}, {\"Deborah\", 47}], isOver18)",
	)

	check.FalseMsg(
		t,
		AnyFunc(Empty[person](), isOver18),
		"AnyFunc([], isOver18)",
	)
}

func TestAll(t *testing.T) {
	check.FalseMsg(
		t,
		All(Over(false, true, false)),
		"All([false true false])",
	)

	check.FalseMsg(
		t,
		All(Over(false, false)),
		"All([false false])",
	)

	check.TrueMsg(
		t,
		All(Over(true, true)),
		"All([true true])",
	)

	check.TrueMsg(
		t,
		All(Empty[bool]()),
		"All([])",
	)
}

func TestAllFunc(t *testing.T) {
	check.FalseMsg(
		t,
		AllFunc(
			Over(person{"Alice", 30}, person{"Bob", 16}),
			isOver18,
		),
		"AllFunc([{\"Alice\", 30}, {\"Bob\", 16}], isOver18)",
	)

	check.FalseMsg(
		t,
		AllFunc(
			Over(person{name: "Bob", age: 16}, person{"Charlie", 17}),
			isOver18,
		),
		"AllFunc([{\"Bob\", 16}, {\"Charlie\", 17}], isOver18)",
	)

	check.TrueMsg(
		t,
		AllFunc(
			Over(person{"Alice", 30}, person{"Deborah", 47}),
			isOver18,
		),
		"AllFunc([{\"Alice\", 30}, {\"Deborah\", 47}], isOver18)",
	)

	check.TrueMsg(
		t,
		AllFunc(Empty[person](), isOver18),
		"AllFunc([], isOver18)",
	)
}

func TestCount(t *testing.T) {
	check.EqMsg(t, Count(Over(1, 2, 3)), 3, "Count([1 2 3])")
	check.EqMsg(t, Count(Empty[int]()), 0, "Count([])")
}

func TestDropWhile(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(DropWhile(Over(1, 2, 3, 2, 1), func(x int) bool { return x < 3 })),
		[]int{3, 2, 1},
		"DropWhile([1 2 3 2 1], x => x < 3)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(DropWhile(Over(1, 2, 3, 2, 1), func(x int) bool { return x < 5 })),
		[]int{},
		"DropWhile([1 2 3 2 1], x => x < 5)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(DropWhile(Over(1, 2, 3, 2, 1), func(x int) bool { return x < 0 })),
		[]int{1, 2, 3, 2, 1},
		"DropWhile([1 2 3 2 1], x => x < 0)",
	)
}

func TestEnumerate(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Enumerate(Over("a", "b", "c"), 0)),
		[]Pair[int, string]{{0, "a"}, {1, "b"}, {2, "c"}},
		"Enumerate([a b c], 0)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Enumerate(Over("a", "b", "c"), 42)),
		[]Pair[int, string]{{42, "a"}, {43, "b"}, {44, "c"}},
		"Enumerate([a b c], 42)",
	)
}

func TestExhaust(t *testing.T) {
	// NOTE: this test relies on the fact that sliceIterator
	//       allows multiple calls to Next() after exhaustion.
	i := Over(1, 2, 3)
	Exhaust(i)
	check.False(t, i.Next())
}

func TestFilter(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Filter(Over(1, 2, 3, 4, 5, 6), isOdd)),
		[]int{1, 3, 5},
		"Filter([1 2 3 4 5 6], isOdd)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Filter(Over(2, 4, 6), isOdd)),
		[]int{},
		"Filter([2 4 6], isOdd)",
	)
}

func TestForEach(t *testing.T) {
	b := strings.Builder{}
	ForEach(Over(1, 2, 3), func(i int) { b.WriteString(strconv.Itoa(i)) })
	check.Eq(t, b.String(), "123")
}

func TestForEachWithError(t *testing.T) {
	b := strings.Builder{}
	expected := errors.New("i can't be negative")

	got := ForEachWithError(
		Over(1, -1, 2),
		func(i int) error {
			if i < 0 {
				return expected
			}
			b.WriteString(strconv.Itoa(i))
			return nil
		},
	)

	check.Eq(t, b.String(), "1")
	check.SpecificErr(t, got, expected)
}

func TestLimit(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Limit(Over(1, 2, 3, 4, 5), 3)),
		[]int{1, 2, 3},
		"Limit([1 2 3 4 5], 3)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Limit(Over(1, 2, 3), 5)),
		[]int{1, 2, 3},
		"Limit([1 2 3], 5)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Limit(Empty[int](), 5)),
		[]int{},
		"Limit([], 5)",
	)
}

func TestMap(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Map(Over(1, 2, 3), func(x int) int { return x + 5 })),
		[]int{6, 7, 8},
		"Map([1 2 3], x => x + 5)",
	)
}

func TestMapWithError(t *testing.T) {
	expectedError := errors.New("i can't be negative")
	it := MapWithError(
		Over(1, 2, -1, -2, 3, 4),
		func(x int) (int, error) {
			if x < 0 {
				return 0, expectedError
			}
			return x + 5, nil
		},
	)

	check.DeepEq(t, IntoSlice(it), []int{6, 7})
	check.SpecificErr(t, it.Err(), expectedError)
}

func TestMin(t *testing.T) {
	min, ok := Min(Over(2, 5, 1, 9, 3))
	check.EqMsg(t, min, 1, "Min([2 5 1 9 3])")
	check.TrueMsg(t, ok, "Min([2 5 1 9 3]): ok")

	min, ok = Min(Empty[int]())
	check.EqMsg(t, min, 0, "Min([])")
	check.FalseMsg(t, ok, "Min([]): ok")
}

func TestMinFunc(t *testing.T) {
	min, ok := MinFunc(
		Over(
			person{"Alice", 30},
			person{"Bob", 25},
			person{"Charlie", 41},
		),
		younger,
	)
	check.DeepEqMsg(t, min, person{"Bob", 25}, "MinFunc([{\"Alice\", 30}, {\"Bob\", 25}, {\"Charlie\", 41}], younger)")
	check.TrueMsg(t, ok, "MinFunc([{\"Alice\", 30}, {\"Bob\", 25}, {\"Charlie\", 41}], younger): ok")

	min, ok = MinFunc(Empty[person](), younger)
	check.DeepEqMsg(t, min, person{"", 0}, "MinFunc([], younger)")
	check.FalseMsg(t, ok, "MinFunc([], younger): ok")
}

func TestMax(t *testing.T) {
	max, ok := Max(Over(2, 5, 1, 9, 3))
	check.EqMsg(t, max, 9, "Max([2 5 1 9 3])")
	check.TrueMsg(t, ok, "Max([2 5 1 9 3]): ok")

	max, ok = Max(Empty[int]())
	check.EqMsg(t, max, 0, "Max([])")
	check.FalseMsg(t, ok, "Max([]): ok")
}

func TestMaxFunc(t *testing.T) {
	max, ok := MaxFunc(
		Over(
			person{"Alice", 30},
			person{"Bob", 25},
			person{"Charlie", 41},
		),
		older,
	)
	check.DeepEqMsg(t, max, person{"Charlie", 41}, "MaxFunc([{\"Alice\", 30}, {\"Bob\", 25}, {\"Charlie\", 41}], older)")
	check.TrueMsg(t, ok, "MaxFunc([{\"Alice\", 30}, {\"Bob\", 25}, {\"Charlie\", 41}], older): ok")

	max, ok = MaxFunc(Empty[person](), older)
	check.DeepEqMsg(t, max, person{"", 0}, "MaxFunc([], older)")
	check.FalseMsg(t, ok, "MaxFunc([], older): ok")
}

func TestNone(t *testing.T) {
	check.FalseMsg(
		t,
		None(Over(false, true, false)),
		"None([false true false])",
	)

	check.TrueMsg(
		t,
		None(Over(false, false)),
		"None([false false])",
	)

	check.FalseMsg(
		t,
		None(Over(true, true)),
		"None([true true])",
	)

	check.TrueMsg(
		t,
		None(Empty[bool]()),
		"None([])",
	)
}

func TestNoneFunc(t *testing.T) {
	check.FalseMsg(
		t,
		NoneFunc(
			Over(person{"Alice", 30}, person{"Bob", 16}),
			isOver18,
		),
		"NoneFunc([{\"Alice\", 30}, {\"Bob\", 16}], isOver18)",
	)

	check.TrueMsg(
		t,
		NoneFunc(
			Over(person{name: "Bob", age: 16}, person{"Charlie", 17}),
			isOver18,
		),
		"NoneFunc([{\"Bob\", 16}, {\"Charlie\", 17}], isOver18)",
	)

	check.FalseMsg(
		t,
		NoneFunc(
			Over(person{"Alice", 30}, person{"Deborah", 47}),
			isOver18,
		),
		"NoneFunc([{\"Alice\", 30}, {\"Deborah\", 47}], isOver18)",
	)

	check.TrueMsg(
		t,
		NoneFunc(Empty[person](), isOver18),
		"NoneFunc([], isOver18)",
	)
}

func TestProduct(t *testing.T) {
	check.EqMsg(t, Product(Over(1, 2, 3)), 6, "Product([1 2 3])")
	check.EqMsg(t, Product(Over(3+0i, 1i)), 3i, "Product([3+0i 1i])")
	check.EqMsg(t, Product(Empty[int]()), 1, "Product([])")
}

func TestReduce(t *testing.T) {
	r, ok := Reduce(Over(1, 2, 3, 4, 5), add)
	check.EqMsg(t, r, 15, "Reduce([1 2 3 4 5], add)")
	check.TrueMsg(t, ok, "Reduce([1 2 3 4 5], add): ok")

	r, ok = Reduce(Over(1), add)
	check.EqMsg(t, r, 1, "Reduce([1], add)")
	check.TrueMsg(t, ok, "Reduce([1], add): ok")

	r, ok = Reduce(Empty[int](), add)
	check.EqMsg(t, r, 0, "Reduce([], add)")
	check.FalseMsg(t, ok, "Reduce([], add): ok")
}

func TestReduceWithInitial(t *testing.T) {
	check.EqMsg(
		t,
		ReduceWithInitial(Over(1, 2, 3, 4, 5), add, 0),
		15,
		"ReduceWithInitial([1 2 3 4 5], add, 0)",
	)

	check.EqMsg(
		t,
		ReduceWithInitial(Over(1, 2, 3, 4, 5), add, 5),
		20,
		"ReduceWithInitial([1 2 3 4 5], add, 5)",
	)

	check.EqMsg(
		t,
		ReduceWithInitial(Over(1), add, 5),
		6,
		"ReduceWithInitial([1], add, 5)",
	)

	check.EqMsg(
		t,
		ReduceWithInitial(Empty[int](), add, 5),
		5,
		"ReduceWithInitial([], add, 5)",
	)
}

func TestSkip(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Skip(Over(1, 2, 3, 4, 5), 2)),
		[]int{3, 4, 5},
		"Skip([1 2 3 4 5], 2)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Skip(Over(1, 2, 3), 5)),
		[]int{},
		"Skip([1 2 3], 5)",
	)
}

func TestSlice(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Slice(Over(1, 2, 3, 4, 5), 1, 3)),
		[]int{2, 3},
		"Slice([1 2 3 4 5], 1, 3)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Slice(Over(1, 2, 3, 4, 5), 0, 3)),
		[]int{1, 2, 3},
		"Slice([1 2 3 4 5], 0, 3)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Slice(Over(1, 2, 3, 4, 5), 2, 5)),
		[]int{3, 4, 5},
		"Slice([1 2 3 4 5], 2, 5)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Slice(Over(1, 2, 3, 4, 5), 3, 3)),
		[]int{},
		"Slice([1 2 3 4 5], 3, 3)",
	)
}

func TestSum(t *testing.T) {
	check.EqMsg(t, Sum(Over(1, 2, 3)), 6, "Sum([1 2 3])")
	check.EqMsg(t, Sum(Over(3+0i, 1i)), 3+1i, "Sum([3+0i 1i])")
	check.EqMsg(t, Sum(Empty[int]()), 0, "Sum([])")
}

func TestTakeWhile(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(TakeWhile(Over(1, 2, 3, 2, 1), func(x int) bool { return x < 3 })),
		[]int{1, 2},
		"TakeWhile([1 2 3 2 1], x => x < 3)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(TakeWhile(Over(1, 2, 3, 2, 1), func(x int) bool { return x < 5 })),
		[]int{1, 2, 3, 2, 1},
		"TakeWhile([1 2 3 2 1], x => x < 5)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(TakeWhile(Over(3, 2, 1), func(x int) bool { return x < 3 })),
		[]int{},
		"TakeWhile([3 2 1], x => x < 3)",
	)
}
