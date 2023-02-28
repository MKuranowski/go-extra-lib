// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package iter

import "golang.org/x/exp/constraints"

// Pair is a utility type containing two possibly heterogenous elements.
//
// Use e.g. by [OverMap] or [Pairwise].
type Pair[T any, U any] struct {
	First  T
	Second U
}

// NumericComparable is a constraint that permits any type which supports arithmetic operators + - * /
// and comparison operators < <= >= >.
// Coincidentally, such types can be constructed from untyped integer constants and compared with != and ==.
type NumericComparable interface {
	constraints.Integer | constraints.Float
}

// Numeric is a constraint that permits any type which supports arithmetic operators + - * /.
// Coincidentally, such types can be constructed from untyped integer constants.
type Numeric interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

// IOReader is an interface satisfied by objects which incrementally
// pull elements from an IO stream.
//
// The Read() method must return (nil, io.EOF) once no more elements are available.
//
// An example of IOReader implementation is [csv.Reader].
type IOReader[T any] interface {
	Read() (T, error)
}
