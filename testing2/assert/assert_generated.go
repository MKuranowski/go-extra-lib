// Code generated by generate_testing2.py. DO NOT EDIT.

package assert

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"golang.org/x/exp/constraints"
)

func True(t *testing.T, got bool) {
	if !got {
		t.Fatalf("got: %v, expected: %v", got, true)
	}
}

func TrueMsg(t *testing.T, got bool, msg string) {
	if !got {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, true)
	}
}

func False(t *testing.T, got bool) {
	if got {
		t.Fatalf("got: %v, expected: %v", got, false)
	}
}

func FalseMsg(t *testing.T, got bool, msg string) {
	if got {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, false)
	}
}

func Err(t *testing.T, got error) {
	if got == nil {
		t.Fatalf("got: %v, expected: %v", got, "non-nil")
	}
}

func ErrMsg(t *testing.T, got error, msg string) {
	if got == nil {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, "non-nil")
	}
}

func NoErr(t *testing.T, got error) {
	if got != nil {
		t.Fatalf("got: %v, expected: %v", got, nil)
	}
}

func NoErrMsg(t *testing.T, got error, msg string) {
	if got != nil {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, nil)
	}
}

func SpecificErr(t *testing.T, got, expected error) {
	if !errors.Is(got, expected) {
		t.Fatalf("got: %v, expected: %v", got, expected)
	}
}

func SpecificErrMsg(t *testing.T, got, expected error, msg string) {
	if !errors.Is(got, expected) {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, expected)
	}
}

func DeepEq(t *testing.T, got, expected any) {
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %v, expected: %v", got, expected)
	}
}

func DeepEqMsg(t *testing.T, got, expected any, msg string) {
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, expected)
	}
}

func Eq[T comparable](t *testing.T, got, expected T) {
	if !(got == expected) {
		t.Fatalf("got: %v, expected: %v", got, expected)
	}
}

func EqMsg[T comparable](t *testing.T, got, expected T, msg string) {
	if !(got == expected) {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, expected)
	}
}

func Ne[T comparable](t *testing.T, got, expected T) {
	if !(got != expected) {
		t.Fatalf("got: %v, expected: %v", got, expected)
	}
}

func NeMsg[T comparable](t *testing.T, got, expected T, msg string) {
	if !(got != expected) {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, expected)
	}
}

func Lt[T constraints.Ordered](t *testing.T, a, b T) {
	if !(a < b) {
		t.Fatalf("expected: %v < %v", a, b)
	}
}

func LtMsg[T constraints.Ordered](t *testing.T, a, b T, msg string) {
	if !(a < b) {
		t.Fatalf("%s: expected: %v < %v", msg, a, b)
	}
}

func Le[T constraints.Ordered](t *testing.T, a, b T) {
	if !(a <= b) {
		t.Fatalf("expected: %v <= %v", a, b)
	}
}

func LeMsg[T constraints.Ordered](t *testing.T, a, b T, msg string) {
	if !(a <= b) {
		t.Fatalf("%s: expected: %v <= %v", msg, a, b)
	}
}

func Gt[T constraints.Ordered](t *testing.T, a, b T) {
	if !(a > b) {
		t.Fatalf("expected: %v > %v", a, b)
	}
}

func GtMsg[T constraints.Ordered](t *testing.T, a, b T, msg string) {
	if !(a > b) {
		t.Fatalf("%s: expected: %v > %v", msg, a, b)
	}
}

func Ge[T constraints.Ordered](t *testing.T, a, b T) {
	if !(a >= b) {
		t.Fatalf("expected: %v >= %v", a, b)
	}
}

func GeMsg[T constraints.Ordered](t *testing.T, a, b T, msg string) {
	if !(a >= b) {
		t.Fatalf("%s: expected: %v >= %v", msg, a, b)
	}
}

func Close(t *testing.T, a, b, delta float64) {
	diff := math.Abs(b - a)
	if diff > delta {
		t.Fatalf("expected %v and %v to be within %v, got: %v", a, b, delta, diff)
	}
}

func CloseMsg(t *testing.T, a, b, delta float64, msg string) {
	diff := math.Abs(b - a)
	if diff > delta {
		t.Fatalf("%s: expected %v and %v to be within %v, got: %v", msg, a, b, delta, diff)
	}
}

func NaN(t *testing.T, got float64) {
	if !math.IsNaN(got) {
		t.Fatalf("got: %v, expected: %v", got, "NaN")
	}
}

func NaNMsg(t *testing.T, got float64, msg string) {
	if !math.IsNaN(got) {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, "NaN")
	}
}

func NotNaN(t *testing.T, got float64) {
	if math.IsNaN(got) {
		t.Fatalf("got: %v, expected: %v", got, "not-NaN")
	}
}

func NotNaNMsg(t *testing.T, got float64, msg string) {
	if math.IsNaN(got) {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, "not-NaN")
	}
}

func Inf(t *testing.T, got float64) {
	if !math.IsInf(got, 0) {
		t.Fatalf("got: %v, expected: %v", got, "Inf")
	}
}

func InfMsg(t *testing.T, got float64, msg string) {
	if !math.IsInf(got, 0) {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, "Inf")
	}
}

func NotInf(t *testing.T, got float64) {
	if math.IsInf(got, 0) {
		t.Fatalf("got: %v, expected: %v", got, "not-Inf")
	}
}

func NotInfMsg(t *testing.T, got float64, msg string) {
	if math.IsInf(got, 0) {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, "not-Inf")
	}
}

func Normal(t *testing.T, got float64) {
	if math.IsNaN(got) || math.IsInf(got, 0) {
		t.Fatalf("got: %v, expected: %v", got, "NaN or Inf")
	}
}

func NormalMsg(t *testing.T, got float64, msg string) {
	if math.IsNaN(got) || math.IsInf(got, 0) {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, "NaN or Inf")
	}
}

func NotNormal(t *testing.T, got float64) {
	if !math.IsNaN(got) && !math.IsInf(got, 0) {
		t.Fatalf("got: %v, expected: %v", got, "not NaN nor Inf")
	}
}

func NotNormalMsg(t *testing.T, got float64, msg string) {
	if !math.IsNaN(got) && !math.IsInf(got, 0) {
		t.Fatalf("%s: got: %v, expected: %v", msg, got, "not NaN nor Inf")
	}
}

