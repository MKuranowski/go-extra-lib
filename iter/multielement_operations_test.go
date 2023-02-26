// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package iter_test

import (
	"errors"
	"testing"

	. "github.com/MKuranowski/go-extra-lib/iter"
	"github.com/MKuranowski/go-extra-lib/testing2/check"
)

func TestChain(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Chain(Over(1, 2), Over(3, 4), Over(5, 6))),
		[]int{1, 2, 3, 4, 5, 6},
		"Chain([1 2], [3 4], [5 6])",
	)
}

func TestChainFromIterator(t *testing.T) {
	its := []Iterator[int]{Over(1, 2), Over(3, 4), Over(5, 6)}
	check.DeepEqMsg(
		t,
		IntoSlice(ChainFromIterator(OverSlice(its))),
		[]int{1, 2, 3, 4, 5, 6},
		"ChainFromIterator([1 2], [3 4], [5 6])",
	)
}

func TestChainMap(t *testing.T) {
	i := ChainMap(Over(1, 5, 10), func(i int) Iterator[int] { return Over(i, i+2) })
	check.DeepEqMsg(
		t,
		IntoSlice(i),
		[]int{1, 3, 5, 7, 10, 12},
		"ChainMap([1 5 10], x => [x, x + 2])",
	)
}

func TestCompress(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Compress(Over(1, 2, 3, 4), Over(true, false, true, false))),
		[]int{1, 3},
		"Compress([1 2 3 4], [true false true false])",
	)
}

func TestCompressFunc(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(CompressFunc(
			Over("a", "b", "c", "d"),
			Over(1, 2, 3, 4),
			isEven,
		)),
		[]string{"b", "d"},
		"Compress([\"a\" \"b\" \"c\" \"d\"], [1 2 3 4], isEven)",
	)
}

func TestGroupBy(t *testing.T) {
	names := []string{"Alice", "Adam", "Amelia", "Andrew", "Bob", "Brian", "Casey", "Chloe", "Craig"}
	groups := IntoSlice(Map(
		GroupBy(OverSlice(names), func(x string) byte { return x[0] }),
		func(g Pair[byte, Iterator[string]]) Pair[byte, []string] {
			return Pair[byte, []string]{g.First, IntoSlice(g.Second)}
		},
	))

	check.DeepEqMsg(
		t,
		groups,
		[]Pair[byte, []string]{
			{'A', []string{"Alice", "Adam", "Amelia", "Andrew"}},
			{'B', []string{"Bob", "Brian"}},
			{'C', []string{"Casey", "Chloe", "Craig"}},
		},
		"GroupBy(names, name => name[0])",
	)
}

func TestGroupByKeysOnly(t *testing.T) {
	names := []string{"Alice", "Adam", "Amelia", "Andrew", "Bob", "Brian", "Casey", "Chloe", "Craig"}
	keys := Map(
		GroupBy(OverSlice(names), func(x string) byte { return x[0] }),
		func(x Pair[byte, Iterator[string]]) byte { return x.First },
	)

	check.DeepEqMsg(
		t,
		IntoSlice(keys),
		[]byte("ABC"),
		"GroupBy(names, name => name[0]): keys only",
	)
}

func TestGroupByUnordered(t *testing.T) {
	names := []string{"Alice", "Andrew", "Bob", "Casey", "Adam", "Amelia", "Chloe", "Craig", "Brian"}
	groups := IntoSlice(Map(
		GroupBy(OverSlice(names), func(x string) byte { return x[0] }),

		// Collect all group values into a slice
		func(g Pair[byte, Iterator[string]]) Pair[byte, []string] {
			return Pair[byte, []string]{g.First, IntoSlice(g.Second)}
		},
	))

	check.DeepEqMsg(
		t,
		groups,
		[]Pair[byte, []string]{
			{'A', []string{"Alice", "Andrew"}},
			{'B', []string{"Bob"}},
			{'C', []string{"Casey"}},
			{'A', []string{"Adam", "Amelia"}},
			{'C', []string{"Chloe", "Craig"}},
			{'B', []string{"Brian"}},
		},
		"GroupBy(namesUnordered, name => name[0])",
	)
}

type department struct {
	id   int
	name string
}

type employee struct {
	id   int
	name string
	dep  department
}

func TestGroupByFunc(t *testing.T) {
	deps := [...]department{
		{1, "HR"},
		{1, "Human Resources"},
		{2, "IT"},
		{2, "Computers"},
		{3, "Sales"},
	}

	employees := []employee{
		// Department 1, HR
		{1, "Alice", deps[0]},
		{2, "Andrew", deps[0]},
		{3, "Bob", deps[1]},

		// Department 2, IT
		{4, "Casey", deps[2]},
		{5, "Adam", deps[3]},
		{6, "Amelia", deps[2]},
		{7, "Chloe", deps[3]},

		// Department 3, Sales
		{8, "Craig", deps[4]},
		{9, "Brian", deps[4]},
	}

	groups := Map(
		GroupByFunc(
			OverSlice(employees),
			func(x employee) department { return x.dep },
			func(a, b department) bool { return a.id == b.id },
		),

		// Collect all group values into a slice
		func(g Pair[department, Iterator[employee]]) Pair[department, []employee] {
			return Pair[department, []employee]{g.First, IntoSlice(g.Second)}
		},
	)

	check.DeepEqMsg(
		t,
		IntoSlice(groups),
		[]Pair[department, []employee]{
			{deps[0], employees[0:3]},
			{deps[2], employees[3:7]},
			{deps[4], employees[7:9]},
		},
		"GroupByFunc(employees, e => e.department, d => d.id)",
	)
}

func TestPairwise(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Pairwise(Over(1, 2, 3), Over("a", "b", "c"))),
		[]Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}},
		"Pairwise([1 2 3], [\"a\" \"b\" \"c\"])",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Pairwise(Over(1, 2, 3), Over("a"))),
		[]Pair[int, string]{{1, "a"}},
		"Pairwise([1 2 3], [\"a\"])",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Pairwise(Empty[int](), Over("a", "b"))),
		[]Pair[int, string]{},
		"Pairwise([], [\"a\" \"b\"])",
	)
}

func TestPairwiseLongest(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(PairwiseLongest(Over(1, 2, 3), Over("a", "b", "c"), 0, "-")),
		[]Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}},
		"PairwiseLongest([1 2 3], [\"a\" \"b\" \"c\"], 0, \"-\")",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(PairwiseLongest(Over(1, 2, 3), Over("a"), 0, "-")),
		[]Pair[int, string]{{1, "a"}, {2, "-"}, {3, "-"}},
		"PairwiseLongest([1 2 3], [\"a\"], 0, \"-\"))",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(PairwiseLongest(Empty[int](), Over("a", "b"), 0, "-")),
		[]Pair[int, string]{{0, "a"}, {0, "b"}},
		"PairwiseLongest([], [\"a\" \"b\"], 0, \"-\"))",
	)
}

func TestPairwiseLongestErr(t *testing.T) {
	err := errors.New("some error")
	i := PairwiseLongest(Error[int](err), Over("a", "b"), 0, "-")

	check.DeepEqMsg(
		t,
		IntoSlice(i),
		[]Pair[int, string]{{0, "a"}, {0, "b"}},
		"PairwiseLongest(Error(someErr), [\"a\" \"b\"], 0, \"-\")",
	)
	check.SpecificErr(t, i.Err(), err)
}

func TestSort(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Sort(Over(2, 3, 1, 0))),
		[]int{0, 1, 2, 3},
		"Sort([2 3 1 0])",
	)
}

func TestSortFunc(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(SortFunc(
			Over(person{"Alice", 30},
				person{"Bob", 25},
				person{"Charlie", 41},
			),
			younger,
		)),
		[]person{{"Bob", 25}, {"Alice", 30}, {"Charlie", 41}},
		"SortFunc(people, p => p.age)",
	)
}

func TestSortStableFunc(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(SortStableFunc(
			Over(person{"Alice", 30},
				person{"Bob", 25},
				person{"Deborah", 25},
				person{"Charlie", 41},
			),
			younger,
		)),
		[]person{{"Bob", 25}, {"Deborah", 25}, {"Alice", 30}, {"Charlie", 41}},
		"SortStableFunc(people, p => p.age)",
	)
}

func TestZip(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Zip(OverString("abc"), OverString("123"), OverString("xyz")),

			// Collect zip elements ([]rune) into strings
			func(x []rune) string { return IntoString(OverSlice(x)) },
		)),
		[]string{"a1x", "b2y", "c3z"},
		"Zip(\"abc\", \"123\", \"xyz\")",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Zip(OverString("abc"), OverString("12"), OverString("x")),

			// Collect zip elements ([]rune) into strings
			func(x []rune) string { return IntoString(OverSlice(x)) },
		)),
		[]string{"a1x"},
		"Zip(\"abc\", \"12\", \"x\")",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Zip(OverString("abc"), Empty[rune]()),

			// Collect zip elements ([]rune) into strings
			func(x []rune) string { return IntoString(OverSlice(x)) },
		)),
		[]string{},
		"Zip(\"abc\", \"\")",
	)
}

func TestZipLongest(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			ZipLongest('-', OverString("abc"), OverString("123"), OverString("xyz")),

			// Collect zip elements ([]rune) into strings
			func(x []rune) string { return IntoString(OverSlice(x)) },
		)),
		[]string{"a1x", "b2y", "c3z"},
		"ZipLongest('-', \"abc\", \"123\", \"xyz\")",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			ZipLongest('-', OverString("ab"), OverString("123"), OverString("x")),

			// Collect zip elements ([]rune) into strings
			func(x []rune) string { return IntoString(OverSlice(x)) },
		)),
		[]string{"a1x", "b2-", "-3-"},
		"ZipLongest('-', \"ab\", \"123\", \"x\")",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			ZipLongest('-', OverString("ab"), Empty[rune]()),

			// Collect zip elements ([]rune) into strings
			func(x []rune) string { return IntoString(OverSlice(x)) },
		)),
		[]string{"a-", "b-"},
		"ZipLongest('-', \"ab\", \"\")",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			ZipLongest('-'),

			// Collect zip elements ([]rune) into strings
			func(x []rune) string { return IntoString(OverSlice(x)) },
		)),
		[]string{},
		"ZipLongest('-')",
	)
}

func TestZipLongestError(t *testing.T) {
	err := errors.New("some error")
	i := ZipLongest(-1, Over(1, 2), Error[int](err))

	check.DeepEqMsg(
		t,
		IntoSlice(i),
		[][]int{{1, -1}, {2, -1}},
		"ZipLongest(-1, [1 2], Error(someErr))",
	)
	check.SpecificErr(t, i.Err(), err)
}
