// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package analyzer

import (
	"bytes"
	"fmt"
	"github.com/mjarmy/golem-lang/ast"
	"reflect"
	"testing"
)

func testGetOk(test *testing.T, s *scope, symbol string, v *ast.Variable) {

	entry, ok := s.get(symbol)

	if !ok {
		test.Error("not ok")
	}

	if !reflect.DeepEqual(entry, v) {
		test.Error(entry, " != ", v)
	}
}

func testGetMissing(test *testing.T, s *scope, symbol string) {
	_, ok := s.get(symbol)
	if ok {
		test.Error("not missing")
	}
}

func testScopeOk(test *testing.T, s *scope, expect string) {
	ds := dumpScope(s)
	if ("\n" + ds) != expect {
		fmt.Println("--------------------------------------------------------------")
		fmt.Println(ds)
		fmt.Println("--------------------------------------------------------------")
		fmt.Println(expect)
		test.Error("scope not ok")
	}
}

func dumpScope(s *scope) string {
	var buf bytes.Buffer

	for s != nil {
		buf.WriteString(s.String())
		buf.WriteString("\n")
		s = s.parent
	}

	return buf.String()
}

func TestGetPut(test *testing.T) {

	s := newFuncScope(nil)

	testGetMissing(test, s, "a")
	s.put("a", true)
	testGetOk(test, s, "a", &ast.Variable{"a", 0, true, false})

	t := newBlockScope(s)
	testGetOk(test, t, "a", &ast.Variable{"a", 0, true, false})

	testGetMissing(test, t, "b")
	t.put("b", false)
	testGetOk(test, t, "b", &ast.Variable{"b", 1, false, false})

	testGetMissing(test, s, "b")
}

func TestCaptureScope(test *testing.T) {

	s0 := newFuncScope(nil)
	s1 := newBlockScope(s0)
	s2 := newFuncScope(s1)
	s3 := newBlockScope(s2)
	s4 := newFuncScope(s3)
	s5 := newBlockScope(s4)

	s0.put("a", false)
	s1.put("b", false)
	s2.put("c", false)
	s3.put("d", false)
	s4.put("e", false)
	s5.put("f", false)

	s5.get("a")
	s5.get("c")

	testScopeOk(test, s5, `
BlockNode defs:{f: (1,false,false)}
Func  defs:{e: (0,false,false)} captures:{a: (0,false,true), c: (1,false,true)} parentCaptures:{a: (0,false,true), c: (0,false,false)} numLocals:2
BlockNode defs:{d: (1,false,false)}
Func  defs:{c: (0,false,false)} captures:{a: (0,false,true)} parentCaptures:{a: (0,false,false)} numLocals:2
BlockNode defs:{b: (1,false,false)}
Func  defs:{a: (0,false,false)} captures:{} parentCaptures:{} numLocals:2
`)
}

func TestPlainStructScope(test *testing.T) {

	stc := &ast.StructExpr{nil, nil, nil, nil, nil, -1}

	s0 := newFuncScope(nil)
	s1 := newBlockScope(s0)
	s2 := newStructScope(s1, stc)

	testScopeOk(test, s2, `
Struct defs:{}
BlockNode defs:{}
Func  defs:{} captures:{} parentCaptures:{} numLocals:0
`)

	if stc.LocalThisIndex != -1 {
		test.Error("LocalThisIndex is wrong", stc.LocalThisIndex, -1)
	}
}

func TestThisStructScope(test *testing.T) {

	struct2 := &ast.StructExpr{nil, nil, nil, nil, nil, -1}
	struct3 := &ast.StructExpr{nil, nil, nil, nil, nil, -1}

	s0 := newFuncScope(nil)
	s1 := newBlockScope(s0)
	s2 := newStructScope(s1, struct2)
	s3 := newStructScope(s2, struct3)

	s0.put("a", false)
	s1.put("b", false)
	s3.this()

	testScopeOk(test, s3, `
Struct defs:{this: (2,true,false)}
Struct defs:{}
BlockNode defs:{b: (1,false,false)}
Func  defs:{a: (0,false,false)} captures:{} parentCaptures:{} numLocals:3
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

	s0 := newFuncScope(nil)
	s1 := newBlockScope(s0)
	s2 := newStructScope(s1, struct2)
	s3 := newFuncScope(s2)
	s4 := newBlockScope(s3)

	s4.this()
	// simulate encountering 'this' again within the s4 block
	s4.this()

	testScopeOk(test, s4, `
BlockNode defs:{}
Func  defs:{} captures:{this: (0,true,true)} parentCaptures:{this: (0,true,false)} numLocals:0
Struct defs:{this: (0,true,false)}
BlockNode defs:{}
Func  defs:{} captures:{} parentCaptures:{} numLocals:1
`)

	if struct2.LocalThisIndex != 0 {
		test.Error("LocalThisIndex is wrong", struct2.LocalThisIndex, 0)
	}
}
