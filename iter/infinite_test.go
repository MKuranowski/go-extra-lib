// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package iter_test

import (
	"math"
	"testing"

	. "github.com/MKuranowski/go-extra-lib/iter"
	"github.com/MKuranowski/go-extra-lib/testing2/check"
)

func TestCycle(t *testing.T) {
	check.EqMsg(
		t,
		IntoString(Cycle(3, 'a', 'b', 'c')),
		"abcabcabc",
		"Cycle(3, 'a', 'b', 'c')",
	)

	check.EqMsg(
		t,
		IntoString(Cycle(0, 'a', 'b', 'c')),
		"",
		"Cycle(0, 'a', 'b', 'c')",
	)

	check.EqMsg(
		t,
		IntoString(Cycle[rune](2)),
		"",
		"Cycle(2)",
	)
}

func TestCycleIter(t *testing.T) {
	check.EqMsg(
		t,
		IntoString(CycleIter(OverString("abc"), 3)),
		"abcabcabc",
		"CycleIter(\"abc\", 3)",
	)

	check.EqMsg(
		t,
		IntoString(CycleIter(OverString("abc"), 0)),
		"",
		"CycleIter(\"abc\", 0)",
	)

	check.EqMsg(
		t,
		IntoString(CycleIter(Empty[rune](), 2)),
		"",
		"CycleIter(\"\", 2)",
	)
}

func TestInfiniteRange(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			InfiniteRange[int](),
			10, // Only take the first 10 elements
		)),
		[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		"InfiniteRange[int]()[:10]",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			InfiniteRange[uint](),
			10, // Only take the first 10 elements
		)),
		[]uint{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		"InfiniteRange[uint]()[:10]",
	)
}

func TestInfiniteRangeFrom(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			InfiniteRangeFrom(int(10)),
			5, // Only take the first 5 elements
		)),
		[]int{10, 11, 12, 13, 14},
		"InfiniteRangeFrom(10)[:5]",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			InfiniteRangeFrom(uint(10)),
			5, // Only take the first 5 elements
		)),
		[]uint{10, 11, 12, 13, 14},
		"InfiniteRangeFrom(uint(10))[:5]",
	)
}

func TestInfiniteRangeFromOverflow(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			InfiniteRangeFrom(int(math.MaxInt)-2),
			5, // Only take the first 5 elements
		)),
		[]int{math.MaxInt - 2, math.MaxInt - 1, math.MaxInt, math.MinInt, math.MinInt + 1},
		"InfiniteRangeFrom(math.MaxInt-2)[:5]",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			InfiniteRangeFrom(uint(math.MaxUint)-2),
			5, // Only take the first 5 elements
		)),
		[]uint{math.MaxUint - 2, math.MaxUint - 1, math.MaxUint, 0, 1},
		"InfiniteRangeFrom(math.MaxUint-2)[:5]",
	)
}

func TestInfiniteRangeWithStep(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			InfiniteRangeWithStep(10, 2),
			5, // Only take the first 5 elements
		)),
		[]int{10, 12, 14, 16, 18},
		"InfiniteRangeWithStep(10, 2)[:5]",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			InfiniteRangeWithStep(10, -2),
			5, // Only take the first 5 elements
		)),
		[]int{10, 8, 6, 4, 2},
		"InfiniteRangeWithStep(10, -2)[:5]",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			InfiniteRangeWithStep(uint(10), 2),
			5, // Only take the first 5 elements
		)),
		[]uint{10, 12, 14, 16, 18},
		"InfiniteRangeWithStep(uint(10), 2)[:5]",
	)
}

func TestInfiniteRangeWithStepOverflow(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			InfiniteRangeWithStep(int(math.MaxInt)-3, 2),
			4, // Only take the first 4 elements
		)),
		[]int{math.MaxInt - 3, math.MaxInt - 1, math.MinInt, math.MinInt + 2},
		"InfiniteRangeWithStep(math.MaxInt-3, 2)[:4]",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			InfiniteRangeWithStep(int(math.MinInt)+3, -2),
			4, // Only take the first 4 elements
		)),
		[]int{math.MinInt + 3, math.MinInt + 1, math.MaxInt, math.MaxInt - 2},
		"InfiniteRangeWithStep(math.MinInt+3, -2)[:4]",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			InfiniteRangeWithStep(uint(math.MaxUint)-3, 2),
			4, // Only take the first 5 elements
		)),
		[]uint{math.MaxUint - 3, math.MaxUint - 1, 0, 2},
		"InfiniteRangeWithStep(math.MaxUint-3, 2)[:4]",
	)
}

func TestRange(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Range(5)),
		[]int{0, 1, 2, 3, 4},
		"Range(5)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Range(0)),
		[]int{},
		"Range(0)",
	)
}

func TestRangeFrom(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(RangeFrom(5, 10)),
		[]int{5, 6, 7, 8, 9},
		"RangeFrom(5, 10)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(RangeFrom(5, 5)),
		[]int{},
		"RangeFrom(5, 5)",
	)
}

func TestRangeWithStep(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(RangeWithStep(5, 11, 2)),
		[]int{5, 7, 9},
		"RangeWithStep(5, 11, 2)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(RangeWithStep(5, 12, 2)),
		[]int{5, 7, 9, 11},
		"RangeWithStep(5, 12, 2)",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(RangeWithStep(5, 5, 2)),
		[]int{},
		"RangeWithStep(5, 5, 2)",
	)
}

func TestRepeatedlyApply(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			RepeatedlyApply(func(x int) int { return x + 5 }, 0),
			5, // Take only the first 5 values
		)),
		[]int{0, 5, 10, 15, 20},
		"RepeatedlyApply(x => x + 5, 0)[:5]",
	)
}

func TestRepeat(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			Repeat(1, 2, 3),
			9, // Only take the first 9 elements (3 repeats)
		)),
		[]int{1, 2, 3, 1, 2, 3, 1, 2, 3},
		"Repeat(1, 2, 3)[:9]",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			Repeat(1),
			3, // Only take the first 3 elements
		)),
		[]int{1, 1, 1},
		"Repeat(1)[:3]",
	)
}

func TestRepeatIter(t *testing.T) {
	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			RepeatIter(Over(1, 2, 3)),
			9, // Only take the first 9 elements (3 repeats)
		)),
		[]int{1, 2, 3, 1, 2, 3, 1, 2, 3},
		"RepeatIter([1 2 3])[:9]",
	)

	check.DeepEqMsg(
		t,
		IntoSlice(Limit(
			RepeatIter(Over(1)),
			3, // Only take the first 3 elements
		)),
		[]int{1, 1, 1},
		"RepeatIter([1])[:3]",
	)
}
