// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package iter_test

import (
	"errors"
	"testing"

	. "github.com/MKuranowski/go-extra-lib/iter"
	"github.com/MKuranowski/go-extra-lib/testing2/assert"
)

func TestOverSlice(t *testing.T) {
	i := OverSlice([]int{1, 2, 3})

	assert.TrueMsg(t, i.Next(), "i.Next(): 1st call")
	assert.EqMsg(t, i.Get(), 1, "i.Get(): 1st call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 2nd call")
	assert.EqMsg(t, i.Get(), 2, "i.Get(): 2nd call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 3rd call")
	assert.EqMsg(t, i.Get(), 3, "i.Get(): 3rd call")

	assert.FalseMsg(t, i.Next(), "i.Next(): 4th call")
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestOver(t *testing.T) {
	i := Over(1, 2, 3)

	assert.TrueMsg(t, i.Next(), "i.Next(): 1st call")
	assert.EqMsg(t, i.Get(), 1, "i.Get(): 1st call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 2nd call")
	assert.EqMsg(t, i.Get(), 2, "i.Get(): 2nd call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 3rd call")
	assert.EqMsg(t, i.Get(), 3, "i.Get(): 3rd call")

	assert.FalseMsg(t, i.Next(), "i.Next(): 4th call")
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestOverChannel(t *testing.T) {
	ch := make(chan int)
	go func() {
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)
	}()

	i := OverChannel(ch)

	assert.TrueMsg(t, i.Next(), "i.Next(): 1st call")
	assert.EqMsg(t, i.Get(), 1, "i.Get(): 1st call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 2nd call")
	assert.EqMsg(t, i.Get(), 2, "i.Get(): 2nd call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 3rd call")
	assert.EqMsg(t, i.Get(), 3, "i.Get(): 3rd call")

	assert.FalseMsg(t, i.Next(), "i.Next(): 4th call")
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestOverMap(t *testing.T) {
	i := OverMap(map[int]string{1: "1", 2: "2", 3: "3"})
	seenOne := false
	seenTwo := false
	seenThree := false

	for n := 0; n < 3; n++ {
		assert.TrueMsg(t, i.Next(), "i.Next()")
		elem := i.Get()

		if elem.First == 1 && elem.Second == "1" {
			seenOne = true
		}

		if elem.First == 2 && elem.Second == "2" {
			seenTwo = true
		}

		if elem.First == 3 && elem.Second == "3" {
			seenThree = true
		}
	}

	assert.TrueMsg(t, seenOne, "i.Get() generated Pair{1, \"1\"}")
	assert.TrueMsg(t, seenTwo, "i.Get() generated Pair{2, \"2\"}")
	assert.TrueMsg(t, seenThree, "i.Get() generated Pair{3, \"3\"}")

	assert.FalseMsg(t, i.Next(), "i.Next(): last call")
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestOverMapKeys(t *testing.T) {
	i := OverMapKeys(map[int]string{1: "1", 2: "2", 3: "3"})
	seenOne := false
	seenTwo := false
	seenThree := false

	for n := 0; n < 3; n++ {
		assert.TrueMsg(t, i.Next(), "i.Next()")
		switch got := i.Get(); got {
		case 1:
			seenOne = true

		case 2:
			seenTwo = true

		case 3:
			seenThree = true

		default:
			t.Fatalf("i.Get(): unexpected key: %v", got)
		}

	}

	assert.TrueMsg(t, seenOne, "i.Get() generated 1")
	assert.TrueMsg(t, seenTwo, "i.Get() generated 2")
	assert.TrueMsg(t, seenThree, "i.Get() generated 3")

	assert.FalseMsg(t, i.Next(), "i.Next(): last call")
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestOverMapValues(t *testing.T) {
	i := OverMapValues(map[int]string{1: "1", 2: "2", 3: "3"})
	seenOne := false
	seenTwo := false
	seenThree := false

	for n := 0; n < 3; n++ {
		assert.TrueMsg(t, i.Next(), "i.Next()")
		switch got := i.Get(); got {
		case "1":
			seenOne = true

		case "2":
			seenTwo = true

		case "3":
			seenThree = true

		default:
			t.Fatalf("i.Get(): unexpected key: %v", got)
		}

	}

	assert.TrueMsg(t, seenOne, "i.Get() generated 1")
	assert.TrueMsg(t, seenTwo, "i.Get() generated 2")
	assert.TrueMsg(t, seenThree, "i.Get() generated 3")

	assert.FalseMsg(t, i.Next(), "i.Next(): last call")
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestOverStringAscii(t *testing.T) {
	i := OverString("foo")

	assert.TrueMsg(t, i.Next(), "i.Next(): 1st call")
	assert.EqMsg(t, i.Get(), 'f', "i.Get(): 1st call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 2nd call")
	assert.EqMsg(t, i.Get(), 'o', "i.Get(): 2nd call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 3rd call")
	assert.EqMsg(t, i.Get(), 'o', "i.Get(): 3rd call")

	assert.FalseMsg(t, i.Next(), "i.Next(): 4th call")
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestOverStringUnicode(t *testing.T) {
	i := OverString("łódź")

	assert.TrueMsg(t, i.Next(), "i.Next(): 1st call")
	assert.EqMsg(t, i.Get(), 'ł', "i.Get(): 1st call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 2nd call")
	assert.EqMsg(t, i.Get(), 'ó', "i.Get(): 2nd call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 3rd call")
	assert.EqMsg(t, i.Get(), 'd', "i.Get(): 3rd call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 4th call")
	assert.EqMsg(t, i.Get(), 'ź', "i.Get(): 4th call")

	assert.FalseMsg(t, i.Next(), "i.Next(): 5th call")
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestOverStringInvalidUnicode(t *testing.T) {
	i := OverString("\xB3\xF3d\x9F") // "łódź" in CP-1250

	assert.TrueMsg(t, i.Next(), "i.Next(): 1st call")
	assert.EqMsg(t, i.Get(), '\uFFFD', "i.Get(): 1st call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 2nd call")
	assert.EqMsg(t, i.Get(), '\uFFFD', "i.Get(): 2nd call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 3rd call")
	assert.EqMsg(t, i.Get(), 'd', "i.Get(): 3rd call")

	assert.TrueMsg(t, i.Next(), "i.Next(): 4th call")
	assert.EqMsg(t, i.Get(), '\uFFFD', "i.Get(): 4th call")

	assert.FalseMsg(t, i.Next(), "i.Next(): 5th call")
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestEmpty(t *testing.T) {
	i := Empty[int]()
	assert.FalseMsg(t, i.Next(), "i.Next()")
	assert.NoErrMsg(t, i.Err(), "i.Err()")
}

func TestError(t *testing.T) {
	dummyErr := errors.New("dummy error")
	i := Error[int](dummyErr)
	assert.FalseMsg(t, i.Next(), "i.Next()")
	assert.SpecificErrMsg(t, i.Err(), dummyErr, "i.Err()")
}
