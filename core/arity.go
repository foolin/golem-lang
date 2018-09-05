// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

type (
	// ArityKind defines the various kinds of Arity for a Func
	ArityKind uint16

	// Arity defines the arity of a function
	Arity struct {
		Kind     ArityKind
		Required uint16
		// For FixedArity and VariadicArity, this value is ignored and should be
		// set to 0.  For MultipleArity, it must be set to a non-zero integer.
		Optional uint16
	}
)

// The various kinds of arity
const (
	// FixedArity means a function always takes a fixed number of parameters
	FixedArity ArityKind = iota

	// VariadicArity means that any extra parameters supplied upon invocation will
	// be collected together into a list.
	VariadicArity

	// MultipleArity means that some of the parameters can be omitted, in which case
	// predifined optional values will be substituted.
	MultipleArity
)

func (a Arity) String() string {

	return fmt.Sprintf(
		"Arity(%s,%d,%d)",
		a.Kind.String(),
		a.Required,
		a.Optional)
}

func (k ArityKind) String() string {

	switch k {
	case FixedArity:
		return "Fixed"
	case VariadicArity:
		return "Variadic"
	case MultipleArity:
		return "Multiple"

	default:
		panic("unreachable")
	}
}
