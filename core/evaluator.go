// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

// Evaluator evaluates Funcs. If the function is a BytecodeFunc, then the
// Evaluator will evaluate the function's bytecode.  In practice, this means that an
// Evaluator is actually a full-fledged instance of the Golem Interpreter.
type Evaluator interface {
	// Eval evaluates a Func.
	Eval(Func, []Value) (Value, Error)
}
