// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package iter_test

import (
	"testing"

	. "github.com/MKuranowski/go-extra-lib/iter"
	"github.com/MKuranowski/go-extra-lib/testing2/check"
)

func TestCartesianProduct(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CartesianProduct([]rune("ABC"), []rune("xy")),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"Ax", "Ay", "Bx", "By", "Cx", "Cy"},
		"CartesianProduct(\"ABC\", \"xy\")",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CartesianProduct([]rune("AB"), []rune("xy"), []rune("12")),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"Ax1", "Ax2", "Ay1", "Ay2", "Bx1", "Bx2", "By1", "By2"},
		"CartesianProduct(\"AB\", \"xy\", \"12\")",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CartesianProduct([]rune("ABC"), []rune{}),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"CartesianProduct(\"ABC\", \"\")",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CartesianProduct[rune](),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"CartesianProduct()",
	)
}

func TestCartesianProductIter(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CartesianProductIter(Over(OverString("ABC"), OverString("xy"))),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"Ax", "Ay", "Bx", "By", "Cx", "Cy"},
		"CartesianProductIter(\"ABC\", \"xy\")",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CartesianProductIter(Over(OverString("AB"), OverString("xy"), OverString("12"))),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"Ax1", "Ax2", "Ay1", "Ay2", "Bx1", "Bx2", "By1", "By2"},
		"CartesianProductIter(\"AB\", \"xy\", \"12\")",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CartesianProductIter(Over(OverString("ABC"), Empty[rune]())),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"CartesianProductIter(\"ABC\", \"\")",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CartesianProductIter(Empty[Iterator[rune]]()),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"CartesianProductIter()",
	)
}

func TestCombinations(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Combinations(2, 'a', 'b', 'c', 'd'),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"ab", "ac", "ad", "bc", "bd", "cd"},
		"Combinations(2, 'a', 'b', 'c', 'd')",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Combinations(3, 'a', 'b', 'c', 'd'),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"abc", "abd", "acd", "bcd"},
		"Combinations(3, 'a', 'b', 'c', 'd')",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Combinations(0, 'a', 'b', 'c'),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{""},
		"Combinations(0, 'a', 'b', 'c')",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Combinations[rune](0),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{""},
		"Combinations(0)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Combinations(3, 'a', 'b'),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"Combinations(3, 'a', 'b')",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Combinations[rune](3),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"Combinations(3)",
	)
}

func TestCombinationsIter(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsIter(OverString("abcd"), 2),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"ab", "ac", "ad", "bc", "bd", "cd"},
		"CombinationsIter(\"abcd\", 2)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsIter(OverString("abcd"), 3),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"abc", "abd", "acd", "bcd"},
		"CombinationsIter(\"abcd\", 3)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsIter(OverString("abc"), 0),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{""},
		"CombinationsIter(\"abc\", 0)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsIter(Empty[rune](), 0),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{""},
		"CombinationsIter(\"\", 0)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsIter(OverString("ab"), 3),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"CombinationsIter(\"ab\", 3)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsIter(Empty[rune](), 3),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"CombinationsIter(\"\", 3)",
	)
}

func TestCombinationsWithReplacement(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsWithReplacement(2, 'a', 'b', 'c'),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"aa", "ab", "ac", "bb", "bc", "cc"},
		"CombinationsWithReplacement(2, 'a', 'b', 'c')",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsWithReplacement(3, 'a', 'b', 'c'),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"aaa", "aab", "aac", "abb", "abc", "acc", "bbb", "bbc", "bcc", "ccc"},
		"CombinationsWithReplacement(3, 'a', 'b', 'c')",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsWithReplacement(3, 'a', 'b'),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"aaa", "aab", "abb", "bbb"},
		"CombinationsWithReplacement(3, 'a', 'b')",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsWithReplacement(0, 'a', 'b'),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{""},
		"CombinationsWithReplacement(0, 'a', 'b')",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsWithReplacement[rune](0),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{""},
		"CombinationsWithReplacement(0)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsWithReplacement[rune](2),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"CombinationsWithReplacement(2)",
	)
}

func TestCombinationsWithReplacementIter(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsWithReplacementIter(OverString("abc"), 2),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"aa", "ab", "ac", "bb", "bc", "cc"},
		"CombinationsWithReplacementIter(\"abc\", 2)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsWithReplacementIter(OverString("abc"), 3),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"aaa", "aab", "aac", "abb", "abc", "acc", "bbb", "bbc", "bcc", "ccc"},
		"CombinationsWithReplacementIter(\"abc\", 3)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsWithReplacementIter(OverString("ab"), 3),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"aaa", "aab", "abb", "bbb"},
		"CombinationsWithReplacementIter(\"ab\", 3)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsWithReplacementIter(OverString("ab"), 0),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{""},
		"CombinationsWithReplacementIter(\"ab\", 0)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsWithReplacementIter(Empty[rune](), 0),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{""},
		"CombinationsWithReplacementIter(\"\", 0)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			CombinationsWithReplacementIter(Empty[rune](), 2),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"CombinationsWithReplacementIter(\"\", 2)",
	)
}

func TestPermutations(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Permutations(3, 'a', 'b', 'c'),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"abc", "acb", "bac", "bca", "cab", "cba"},
		"Permutations(3, 'a', 'b', 'c')",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Permutations(2, 'a', 'b', 'c'),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"ab", "ac", "ba", "bc", "ca", "cb"},
		"Permutations(2, 'a', 'b', 'c')",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Permutations(0, 'a', 'b', 'c'),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{""},
		"Permutations(0, 'a', 'b', 'c')",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Permutations[rune](0),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{""},
		"Permutations(0)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Permutations(4, 'a', 'b', 'b'),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"Permutations(4, 'a', 'b', 'c')",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			Permutations[rune](4),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"Permutations(4)",
	)
}

func TestPermutationsIter(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			PermutationsIter(OverString("abc"), 3),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"abc", "acb", "bac", "bca", "cab", "cba"},
		"PermutationsIter(\"abc\", 3)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			PermutationsIter(OverString("abc"), 2),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{"ab", "ac", "ba", "bc", "ca", "cb"},
		"PermutationsIter(\"abc\", 2)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			PermutationsIter(OverString("abc"), 0),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{""},
		"PermutationsIter(\"abc\", 0)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			PermutationsIter(Empty[rune](), 0),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{""},
		"PermutationsIter(\"\", 0)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			PermutationsIter(OverString("abc"), 4),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"PermutationsIter(\"abc\", 4)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Map(
			PermutationsIter(Empty[rune](), 4),
			func(x []rune) string { return string(x) }, // collect into strings
		)),
		[]string{},
		"PermutationsIter(\"\", 4)",
	)
}

func TestPowerSet(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(PowerSet[int]()),
		[][]int{nil},
		"PowerSet()",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(PowerSet(1)),
		[][]int{nil, {1}},
		"PowerSet(1)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(PowerSet(1, 2)),
		[][]int{nil, {1}, {2}, {1, 2}},
		"PowerSet(1, 2)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(PowerSet(1, 2, 3)),
		[][]int{nil, {1}, {2}, {1, 2}, {3}, {1, 3}, {2, 3}, {1, 2, 3}},
		"PowerSet(1, 2, 3)",
	)
}

func TestPowerSetIter(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(PowerSetIter(Empty[int]())),
		[][]int{nil},
		"PowerSetIter()",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(PowerSetIter(Over(1))),
		[][]int{nil, {1}},
		"PowerSetIter([1])",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(PowerSetIter(Over(1, 2))),
		[][]int{nil, {1}, {2}, {1, 2}},
		"PowerSetIter([1, 2])",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(PowerSetIter(Over(1, 2, 3))),
		[][]int{nil, {1}, {2}, {1, 2}, {3}, {1, 3}, {2, 3}, {1, 2, 3}},
		"PowerSetIter([1, 2, 3])",
	)
}
