// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
//"fmt"
)

type (
	fieldMap interface {
		names() []string
		has(string) bool
		get(string, Evaluator) (Value, Error)
		invoke(string, Evaluator, []Value) (Value, Error)
		set(string, Evaluator, Value) Error
		internalReplace(string, Field)
	}
)

//--------------------------------------------------------------

// baseHashFieldMap
// hashFieldMap (implements InternalReplace)
// mergedHashFieldMap (does not implement InternalReplace)

//--------------------------------------------------------------

// virtualFieldMap (does not implement InternalReplace)

// VirtualFieldMap will always only have (Arity, Invoker) tuples
// func (s str) lookupFunc(name string) (Arity, Invoke, Error) {

// When merging a VirtualFieldMap, it will be necessary to iterate over
// all its tuples and create the actual Funcs.
