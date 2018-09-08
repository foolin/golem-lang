// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
	"github.com/mjarmy/golem-lang/scanner"
)

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
		panic("tassert")
	}
}

var builtins []*g.BuiltinEntry = []*g.BuiltinEntry{
	{Name: "assert", Value: g.BuiltinAssert},
	{Name: "println", Value: g.BuiltinPrintln},
}

var builtinMgr = g.NewBuiltinManager(builtins)

func testCompile(t *testing.T, code string) *bc.Module {

	source := &scanner.Source{Name: "foo", Path: "foo.glm", Code: code}
	mods, errs := compiler.CompileSourceFully(builtinMgr, source, nil)
	if errs != nil {
		fmt.Printf("%v\n", errs)
	}
	tassert(t, errs == nil)
	tassert(t, len(mods) == 1)

	return mods[0]
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

func TestModuleContents(t *testing.T) {

	code := `
let a = 0
const b = "xyz"
fn c() {}
`
	mod := testCompile(t, code)
	itp := NewInterpreter(builtinMgr, []*bc.Module{mod})

	result, es := itp.InitModules()
	tassert(t, es == nil && len(result) == 1)

	stc := mod.Contents

	val, err := stc.GetField(nil, "a")
	tassert(t, err == nil && val.Type() == g.IntType)
	val, err = stc.GetField(nil, "b")
	tassert(t, err == nil && val.Type() == g.StrType)
	val, err = stc.GetField(nil, "c")
	tassert(t, err == nil && val.Type() == g.FuncType)

	err = stc.SetField(nil, "a", g.One)
	tassert(t, err == nil)
	err = stc.SetField(nil, "b", g.One)
	tassert(t, err.Error() == "ReadonlyField: Field 'b' is readonly")
	err = stc.SetField(nil, "c", g.One)
	tassert(t, err.Error() == "ReadonlyField: Field 'c' is readonly")

	code = `
let a = 1
return a + 2
let b = 5
`
	mod = testCompile(t, code)
	itp = NewInterpreter(builtinMgr, []*bc.Module{mod})

	result, es = itp.InitModules()
	tassert(t, es == nil && len(result) == 1)
	ok(t, result[0], nil, g.NewInt(3))

	stc = mod.Contents
	val, err = stc.GetField(nil, "a")
	ok(t, val, err, g.One)
	val, err = stc.GetField(nil, "b")
	ok(t, val, err, g.Null)
}

func TestImport(t *testing.T) {

	srcMain := &scanner.Source{Name: "foo", Path: "foo.glm", Code: `
import a, b, c
assert([1, 2, 3] == [a.x, b.y, c.z])
`}
	sourceMap := map[string]*scanner.Source{
		"a": &scanner.Source{Name: "a", Path: "a.glm", Code: "import c; let x = 1;"},
		"b": &scanner.Source{Name: "b", Path: "b.glm", Code: "import c; let y = c.z - 1;"},
		"c": &scanner.Source{Name: "c", Path: "c.glm", Code: "let z = 3;"},
	}
	resolver := func(moduleName string) (*scanner.Source, error) {
		if src, ok := sourceMap[moduleName]; ok {
			return src, nil
		}
		return nil, g.UndefinedModule(moduleName)
	}

	mods, errs := compiler.CompileSourceFully(builtinMgr, srcMain, resolver)
	tassert(t, errs == nil)
	tassert(t, len(mods) == 4)

	itp := NewInterpreter(builtinMgr, mods)
	result, es := itp.InitModules()
	tassert(t, es == nil)
	tassert(t, len(result) == 4)
}

//--------------------------------------------------------------
//--------------------------------------------------------------
//--------------------------------------------------------------

func okInterp(t *testing.T, mods []*bc.Module) {
	intp := NewInterpreter(builtinMgr, mods)
	_, es := intp.InitModules()
	if es != nil {
		panic(es)
	}
}

func failInterp(t *testing.T, mods []*bc.Module, expect ErrorStruct) {
	intp := NewInterpreter(builtinMgr, mods)
	_, es := intp.InitModules()
	tassert(t, es != nil)

	//dumpErrorStruct("failInterp", es)
	tassert(t, reflect.DeepEqual(es, expect))
}

func TestHandleError(t *testing.T) {

	fmt.Printf("--------------------------------------------------------------\n")
	code := `
		1/0
		`
	failInterp(t, []*bc.Module{testCompile(t, code)},
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:2"}))

	fmt.Printf("--------------------------------------------------------------\n")
	code = `
		let a = (|| => 1/0)
		a()
		`
	failInterp(t, []*bc.Module{testCompile(t, code)},
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:2",
				"    at foo.glm:3"}))

	fmt.Printf("--------------------------------------------------------------\n")
	code = `
		let s = struct {
			q: prop { ||=> 1/0 }
		}
		let a = s.q
		`
	failInterp(t, []*bc.Module{testCompile(t, code)},
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:3",
				"    at foo.glm:5"}))

	fmt.Printf("--------------------------------------------------------------\n")
	code = `
		[1, 2, 3].map(
			|e| => 1/0)
		`
	failInterp(t, []*bc.Module{testCompile(t, code)},
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:3",
				"    at foo.glm:2"}))

	fmt.Printf("--------------------------------------------------------------\n")
	code = `
		let s = struct {
			q: prop { fn() {
					[1, 2, 3].map(
						|e| => 1/0)
				}
			}
		}
		s.q
		`
	failInterp(t, []*bc.Module{testCompile(t, code)},
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:5",
				"    at foo.glm:4",
				"    at foo.glm:9"}))

	fmt.Printf("--------------------------------------------------------------\n")
	code = `
		let s = struct {
			q: prop { ||=> 1/0 }
		}

		let a = [1, 2, 3].map(
			|e| => s.q)
		`
	failInterp(t, []*bc.Module{testCompile(t, code)},
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:3",
				"    at foo.glm:7",
				"    at foo.glm:6"}))
}

func TestTry(t *testing.T) {

	fmt.Printf("--------------------------------------------------------------\n")
	code := `
			let a = 0
			try {
				1/0
			} catch e {
				a = 1
			}
			assert(a == 1)
			`
	okInterp(t, []*bc.Module{testCompile(t, code)})

	fmt.Printf("--------------------------------------------------------------\n")
	code = `
			let b = 0
			try {
				1/0
			} finally {
				b = 1
			}
			assert(b == 1)
			`
	okInterp(t, []*bc.Module{testCompile(t, code)})

	fmt.Printf("--------------------------------------------------------------\n")
	code = `
			let a = 0
			let b = 0
			try {
				1/0
			} catch e {
				a = 1
			} finally {
				b = 1
			}
			assert(a == 1)
			assert(b == 1)
			`
	okInterp(t, []*bc.Module{testCompile(t, code)})

	fmt.Printf("--------------------------------------------------------------\n")
	code = `
			let a = 0
			try {
				let b = 0
				try {
					1/0
				} finally {
					b = 1
				}
				assert(b == 1)
			} catch e {
				a = 1
			}
			assert(a == 0)
			`
	okInterp(t, []*bc.Module{testCompile(t, code)})

	fmt.Printf("--------------------------------------------------------------\n")
	code = `
			let b = 0
			try {
				let a = 0
				try {
					1/0
				} catch e {
					a = 1
				}
				assert(a == 1)
			} finally {
				b = 1
			}
			assert(b == 1)
		`
	okInterp(t, []*bc.Module{testCompile(t, code)})
}
