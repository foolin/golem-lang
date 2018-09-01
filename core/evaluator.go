// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

// Evaluator evaluates functions that are defined via bytecode.
// In practice, this means that an Evaluator is actually a full-fledged instance
// of the Golem Interpreter.
type Evaluator interface {
	Eval(Func, []Value) (Value, Error)
}
