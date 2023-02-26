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

// Numeric is a constraint that permits any type which supports arithmetic operators + - * /.
// Coincidentally, such types can be constructed from untyped integer constants.
type Numeric interface {
	constraints.Integer | constraints.Float | constraints.Complex
}
