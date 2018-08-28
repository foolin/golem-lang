// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"reflect"
	"testing"

	g "github.com/mjarmy/golem-lang/ncore"
)

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
		panic("tassert")
	}
}

func ok(t *testing.T, val g.Value, err g.Error, expect g.Value) {

	if err != nil {
		t.Error(err, " != ", nil)
	}

	if !reflect.DeepEqual(val, expect) {
		t.Error(val, " != ", expect)
		panic("ok")
	}
}

func okType(t *testing.T, val g.Value, expected g.Type) {
	tassert(t, val.Type() == expected)
}

func TestBytecodeFunc(t *testing.T) {

	a := NewBytecodeFunc(&FuncTemplate{})
	b := NewBytecodeFunc(&FuncTemplate{})

	okType(t, a, g.FuncType)
	okType(t, b, g.FuncType)

	v, err := a.Eq(nil, a)
	ok(t, v, err, g.True)

	v, err = b.Eq(nil, b)
	ok(t, v, err, g.True)

	v, err = a.Eq(nil, b)
	ok(t, v, err, g.False)

	v, err = b.Eq(nil, a)
	ok(t, v, err, g.False)
}

func TestLineNumber(t *testing.T) {

	tp := &FuncTemplate{
		Module:      nil,
		Arity:       nil,
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
