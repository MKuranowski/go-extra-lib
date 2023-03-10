// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package slices2_test

import (
	"testing"

	"github.com/MKuranowski/go-extra-lib/slices2"
	"github.com/MKuranowski/go-extra-lib/testing2/check"
)

func TestBatchesEven(t *testing.T) {
	check.DeepEqMsg(
		t,
		slices2.Batches([]int{1, 2, 3, 4}, 2),
		[][]int{{1, 2}, {3, 4}},
		"after partitioning",
	)
}

func TestBatchesUneven(t *testing.T) {
	check.DeepEqMsg(
		t,
		slices2.Batches([]int{1, 2, 3, 4, 5}, 2),
		[][]int{{1, 2}, {3, 4}, {5}},
		"after partitioning",
	)
}

func TestDeleteAndSetToZero(t *testing.T) {
	old := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	new := slices2.DeleteAndSetToZero(old, 3, 6)
	check.DeepEqMsg(t, new, []int{1, 2, 3, 7, 8, 9}, "new slice")
	check.DeepEqMsg(t, old, []int{1, 2, 3, 7, 8, 9, 0, 0, 0}, "original slice")
}

func TestExpand(t *testing.T) {
	s := []int{1, 2, 3, 4}
	s = slices2.Expand(s, 2, 2)
	check.DeepEqMsg(t, s, []int{1, 2, 0, 0, 3, 4}, "after expand")
}

func TestExtend(t *testing.T) {
	s := []int{1, 2, 3, 4}
	s = slices2.Extend(s, 2)
	check.DeepEqMsg(t, s, []int{1, 2, 3, 4, 0, 0}, "after extend")
}

func TestFilter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8}
	s = slices2.Filter(s, func(x int) bool { return x%2 == 0 })
	check.DeepEqMsg(t, s, []int{2, 4, 6, 8}, "after filter")
}

func TestFilterAndSetToZero(t *testing.T) {
	old := []int{1, 2, 3, 4, 5, 6, 7, 8}
	new := slices2.FilterAndSetToZero(old, func(x int) bool { return x%2 == 0 })
	check.DeepEqMsg(t, new, []int{2, 4, 6, 8}, "new slice")
	check.DeepEqMsg(t, old, []int{2, 4, 6, 8, 0, 0, 0, 0}, "original slice")
}

func TestReverseEven(t *testing.T) {
	s := []int{1, 2, 3, 4}
	slices2.Reverse(s)
	check.DeepEqMsg(t, s, []int{4, 3, 2, 1}, "after reverse")
}

func TestReverseOdd(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	slices2.Reverse(s)
	check.DeepEqMsg(t, s, []int{5, 4, 3, 2, 1}, "after reverse")
}

func TestSlidingWindow(t *testing.T) {
	check.DeepEqMsg(
		t,
		slices2.SlidingWindow([]int{1, 2, 3, 4, 5}, 2),
		[][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}},
		"windows",
	)
}

func TestSlidingWindowSmall(t *testing.T) {
	check.DeepEqMsg(
		t,
		slices2.SlidingWindow([]int{1, 2}, 3),
		[][]int{{1, 2}},
		"windows",
	)
}
