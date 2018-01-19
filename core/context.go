// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

//---------------------------------------------------------------
// A Context can evaluate bytecode.

type Context interface {
	Eval(BytecodeFunc, []Value) (Value, Error)
}
