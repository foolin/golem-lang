// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	//"fmt"
	"testing"
)

func TestBytecodeFunc(t *testing.T) {

	a := NewBytecodeFunc(&FuncTemplate{})
	b := NewBytecodeFunc(&FuncTemplate{})

	okType(t, a, FuncType)
	okType(t, b, FuncType)

	v, err := a.Eq(cx, a)
	ok(t, v, err, True)

	v, err = b.Eq(cx, b)
	ok(t, v, err, True)

	v, err = a.Eq(cx, b)
	ok(t, v, err, False)

	v, err = b.Eq(cx, a)
	ok(t, v, err, False)
}

func TestLineNumber(t *testing.T) {

	tp := &FuncTemplate{
		ModuleName:  "",
		ModulePath:  "",
		Arity:       0,
		NumCaptures: 0,
		NumLocals:   0,
		OpCodes:     nil,
		LineNumberTable: []LineNumberEntry{
			{0, 0},
			{1, 2},
			{11, 3},
			{20, 4},
			{29, 0}},
		ExceptionHandlers: nil,
	}

	tassert(t, tp.LineNumber(0) == 0)
	tassert(t, tp.LineNumber(1) == 2)
	tassert(t, tp.LineNumber(10) == 2)
	tassert(t, tp.LineNumber(11) == 3)
	tassert(t, tp.LineNumber(19) == 3)
	tassert(t, tp.LineNumber(20) == 4)
	tassert(t, tp.LineNumber(28) == 4)
	tassert(t, tp.LineNumber(29) == 0)
}
