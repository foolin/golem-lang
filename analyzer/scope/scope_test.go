// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package scope

import (
	"bytes"
	"fmt"
	"github.com/mjarmy/golem-lang/ast"
	"reflect"
	"testing"
)

func testGetOk(test *testing.T, s *Scope, symbol string, v *ast.Variable) {

	entry, ok := s.Get(symbol)

	if !ok {
		test.Error("not ok")
	}

	if !reflect.DeepEqual(entry, v) {
		test.Error(entry, " != ", v)
	}
}

func testGetMissing(test *testing.T, s *Scope, symbol string) {
	_, ok := s.Get(symbol)
	if ok {
		test.Error("not missing")
	}
}

func testScopeOk(test *testing.T, s *Scope, expect string) {
	ds := dumpScope(s)
	if ("\n" + ds) != expect {
		fmt.Println("--------------------------------------------------------------")
		fmt.Println(ds)
		fmt.Println("--------------------------------------------------------------")
		fmt.Println(expect)
		test.Error("Scope not ok")
	}
}

func dumpScope(s *Scope) string {
	var buf bytes.Buffer

	for s != nil {
		buf.WriteString(s.String())
		buf.WriteString("\n")
		s = s.Parent
	}

	return buf.String()
}

func TestGetPut(test *testing.T) {

	s := NewFuncScope(nil)

	testGetMissing(test, s, "a")
	s.Put("a", true)
	testGetOk(test, s, "a", &ast.Variable{"a", 0, true, false})

	t := NewBlockScope(s)
	testGetOk(test, t, "a", &ast.Variable{"a", 0, true, false})

	testGetMissing(test, t, "b")
	t.Put("b", false)
	testGetOk(test, t, "b", &ast.Variable{"b", 1, false, false})

	testGetMissing(test, s, "b")
}

func TestCaptureScope(test *testing.T) {

	s0 := NewFuncScope(nil)
	s1 := NewBlockScope(s0)
	s2 := NewFuncScope(s1)
	s3 := NewBlockScope(s2)
	s4 := NewFuncScope(s3)
	s5 := NewBlockScope(s4)

	s0.Put("a", false)
	s1.Put("b", false)
	s2.Put("c", false)
	s3.Put("d", false)
	s4.Put("e", false)
	s5.Put("f", false)

	s5.Get("a")
	s5.Get("c")

	testScopeOk(test, s5, `
Block defs:{f: (1,false,false)}
Func defs:{e: (0,false,false)} captures:{a: (0,false,true), c: (1,false,true)} parentCaptures:{a: (0,false,true), c: (0,false,false)} numLocals:2
Block defs:{d: (1,false,false)}
Func defs:{c: (0,false,false)} captures:{a: (0,false,true)} parentCaptures:{a: (0,false,false)} numLocals:2
Block defs:{b: (1,false,false)}
Func defs:{a: (0,false,false)} captures:{} parentCaptures:{} numLocals:2
`)
}

func TestPlainStructScope(test *testing.T) {

	stc := &ast.StructExpr{nil, nil, nil, nil, nil, -1}

	s0 := NewFuncScope(nil)
	s1 := NewBlockScope(s0)
	s2 := NewStructScope(s1, stc)

	testScopeOk(test, s2, `
Struct defs:{}
Block defs:{}
Func defs:{} captures:{} parentCaptures:{} numLocals:0
`)

	if stc.LocalThisIndex != -1 {
		test.Error("LocalThisIndex is wrong", stc.LocalThisIndex, -1)
	}
}

func TestThisStructScope(test *testing.T) {

	struct2 := &ast.StructExpr{nil, nil, nil, nil, nil, -1}
	struct3 := &ast.StructExpr{nil, nil, nil, nil, nil, -1}

	s0 := NewFuncScope(nil)
	s1 := NewBlockScope(s0)
	s2 := NewStructScope(s1, struct2)
	s3 := NewStructScope(s2, struct3)

	s0.Put("a", false)
	s1.Put("b", false)
	s3.This()

	testScopeOk(test, s3, `
Struct defs:{this: (2,true,false)}
Struct defs:{}
Block defs:{b: (1,false,false)}
Func defs:{a: (0,false,false)} captures:{} parentCaptures:{} numLocals:3
`)

	if struct2.LocalThisIndex != -1 {
		test.Error("LocalThisIndex is wrong", struct2.LocalThisIndex, -1)
	}
	if struct3.LocalThisIndex != 2 {
		test.Error("LocalThisIndex is wrong", struct3.LocalThisIndex, 2)
	}
}

func TestMethodScope(test *testing.T) {

	struct2 := &ast.StructExpr{nil, nil, nil, nil, nil, -1}

	s0 := NewFuncScope(nil)
	s1 := NewBlockScope(s0)
	s2 := NewStructScope(s1, struct2)
	s3 := NewFuncScope(s2)
	s4 := NewBlockScope(s3)

	s4.This()
	// simulate encountering 'this' again within the s4 block
	s4.This()

	testScopeOk(test, s4, `
Block defs:{}
Func defs:{} captures:{this: (0,true,true)} parentCaptures:{this: (0,true,false)} numLocals:0
Struct defs:{this: (0,true,false)}
Block defs:{}
Func defs:{} captures:{} parentCaptures:{} numLocals:1
`)

	if struct2.LocalThisIndex != 0 {
		test.Error("LocalThisIndex is wrong", struct2.LocalThisIndex, 0)
	}
}
