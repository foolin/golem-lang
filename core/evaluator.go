// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

// Eval evaluates Funcs.
//
// If the function is a NativeFunc, then Eval will call the native golang code.
//
// If the function is a BytecodeFunc, then Eval will evaluate the function's bytecode.
//
// In practice, this means that an Eval is actually always
// a full-fledged instance of the Golem Interpreter.
type Eval interface {
	// Eval evaluates a Func.
	Eval(Func, []Value) (Value, Error)
}
