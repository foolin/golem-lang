// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"reflect"
	"sort"
	"testing"

	g "github.com/mjarmy/golem-lang/core"
)

func Tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
		panic("Tassert")
	}
}

func ok(t *testing.T, val interface{}, err g.Error, expect interface{}) {

	if err != nil {
		t.Error(err, " != ", nil)
		panic("ok")
	}

	if !reflect.DeepEqual(val, expect) {
		t.Error(val, " != ", expect)
		panic("ok")
	}
}

func okNames(t *testing.T, val []string, err g.Error, expect []string) {

	sort.Slice(val, func(i, j int) bool {
		return val[i] < val[j]
	})

	sort.Slice(expect, func(i, j int) bool {
		return expect[i] < expect[j]
	})

	ok(t, val, err, expect)
}

func okType(t *testing.T, val g.Value, expected g.Type) {
	Tassert(t, val.Type() == expected)
}

func TestBytecodeFunc(t *testing.T) {

	var a g.Func = NewBytecodeFunc(&FuncTemplate{})
	var b g.Func = NewBytecodeFunc(&FuncTemplate{})

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

	names, err := a.FieldNames()
	okNames(t, names, err, []string{})

	has, err := a.HasField("a")
	ok(t, has, err, false)
}

func TestLineNumber(t *testing.T) {

	tp := &FuncTemplate{
		Module:      nil,
		Arity:       g.Arity{},
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

	Tassert(t, tp.LineNumber(0) == 0)
	Tassert(t, tp.LineNumber(1) == 2)
	Tassert(t, tp.LineNumber(10) == 2)
	Tassert(t, tp.LineNumber(11) == 3)
	Tassert(t, tp.LineNumber(19) == 3)
	Tassert(t, tp.LineNumber(20) == 4)
	Tassert(t, tp.LineNumber(28) == 4)
	Tassert(t, tp.LineNumber(29) == 0)
}
