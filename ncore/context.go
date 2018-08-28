// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

// Context evaluates functions (including functions defined via bytecode).
// In practice, this means that a Context is actually a full-fledged instance
// of the Golem Interpreter.
type Context interface {
	Eval(Func, []Value) (Value, Error)
}
