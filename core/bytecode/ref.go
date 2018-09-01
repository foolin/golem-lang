// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
)

// Ref is a container for a Value.  Refs are used by the interpreter
// as a place to store the value of a variable.
type Ref struct {
	Val g.Value
}

// NewRef creates a new Ref
func NewRef(val g.Value) *Ref {
	return &Ref{val}
}

func (r *Ref) String() string {
	return fmt.Sprintf("Ref(%v)", r.Val)
}
