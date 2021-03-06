// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

// Eval evaluates Funcs.  In practice, an Eval is actually always a full-fledged instance
// of the Golem Interpreter.
type Eval interface {
	// Eval evaluates a Func.
	// If the function is a NativeFunc, then Eval will call native golang code.
	// If the function is a bytecode.Func, then Eval will evaluate the function's bytecode.
	Eval(Func, []Value) (Value, Error)
}
