// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package iter_test

import (
	"testing"

	. "github.com/MKuranowski/go-extra-lib/iter"
	"github.com/MKuranowski/go-extra-lib/testing2/assert"
)

func TestIntoSlice(t *testing.T) {
	expected := []int{1, 2, 3}
	i := OverSlice(expected)
	got := IntoSlice(i)

	assert.DeepEq(t, got, expected)
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestIntoMap(t *testing.T) {
	expected := map[int]string{1: "1", 2: "2", 3: "3"}
	i := OverMap(expected)
	got := IntoMap(i)

	assert.DeepEq(t, got, expected)
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestIntoChannel(t *testing.T) {
	i := OverSlice([]int{1, 2, 3})
	ch := IntoChannel(i)

	assert.EqMsg(t, <-ch, 1, "1st element")
	assert.EqMsg(t, <-ch, 2, "2nd element")
	assert.EqMsg(t, <-ch, 3, "3rd element")

	_, isOpen := <-ch
	assert.FalseMsg(t, isOpen, "channel closed")

	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestIntoStringAscii(t *testing.T) {
	assert.Eq(t, IntoString(Over('f', 'o', 'o')), "foo")
}

func TestIntoStringUnicode(t *testing.T) {
	assert.Eq(t, IntoString(Over('ł', 'ó', 'd', 'ź')), "łódź")
}

func TestSendOver(t *testing.T) {
	i := OverSlice([]int{1, 2, 3})
	ch := make(chan int)
	go SendOver(i, ch)

	assert.EqMsg(t, <-ch, 1, "1st element")
	assert.EqMsg(t, <-ch, 2, "2nd element")
	assert.EqMsg(t, <-ch, 3, "3rd element")
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}
