// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"reflect"
	"testing"

	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
	"github.com/mjarmy/golem-lang/scanner"
)

// NOTE: Most of the test suite for the interpreter is over in
// 'bench_test/core_test.glm'. Testing is done there by running the Golem CLI
// application against Golem code.

var builtins []*g.BuiltinEntry = []*g.BuiltinEntry{
	{Name: "assert", Value: g.BuiltinAssert},
}

var builtinMgr = g.NewBuiltinManager(builtins)

func testCompile(t *testing.T, code string) *bc.Module {

	source := &scanner.Source{Name: "foo", Path: "foo.glm", Code: code}
	mods, errs := compiler.CompileSourceFully(builtinMgr, source, nil)
	g.Tassert(t, errs == nil)
	g.Tassert(t, len(mods) == 1)

	return mods[0]
}

func testInterpret(mods []*bc.Module) *Interpreter {
	intp := NewInterpreter(builtinMgr, mods)
	_, err := intp.InitModules()
	if err != nil {
		panic(err)
	}
	return intp
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

func fail(t *testing.T, code string, expect g.ErrorStruct) {

	mod := testCompile(t, code)
	intp := NewInterpreter(builtinMgr, []*bc.Module{mod})

	result, err := intp.InitModules()
	if result != nil {
		panic(result)
	}

	eq, e := expect.Eq(nil, err)
	g.Tassert(t, e == nil)

	if !eq.(g.Bool).BoolVal() {
		t.Error(mustStr(err), " != ", mustStr(expect))
		panic("fail")
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

	result, errStruct := itp.InitModules()
	g.Tassert(t, errStruct == nil && len(result) == 1)

	stc := mod.Contents

	val, err := stc.GetField("a", nil)
	g.Tassert(t, err == nil && val.Type() == g.IntType)
	val, err = stc.GetField("b", nil)
	g.Tassert(t, err == nil && val.Type() == g.StrType)
	val, err = stc.GetField("c", nil)
	g.Tassert(t, err == nil && val.Type() == g.FuncType)

	err = stc.SetField("a", nil, g.One)
	g.Tassert(t, err == nil)
	err = stc.SetField("b", nil, g.One)
	g.Tassert(t, err.Error() == "ReadonlyField: Field 'b' is readonly")
	err = stc.SetField("c", nil, g.One)
	g.Tassert(t, err.Error() == "ReadonlyField: Field 'c' is readonly")

	code = `
let a = 1
return a + 2
let b = 5
`
	mod = testCompile(t, code)
	itp = NewInterpreter(builtinMgr, []*bc.Module{mod})

	result, errStruct = itp.InitModules()
	g.Tassert(t, errStruct == nil && len(result) == 1)
	ok(t, result[0], nil, g.NewInt(3))

	stc = mod.Contents
	val, err = stc.GetField("a", nil)
	ok(t, val, err, g.One)
	val, err = stc.GetField("b", nil)
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
		return nil, g.UndefinedModuleError(moduleName)
	}

	mods, errs := compiler.CompileSourceFully(builtinMgr, srcMain, resolver)
	g.Tassert(t, errs == nil)
	g.Tassert(t, len(mods) == 4)

	itp := NewInterpreter(builtinMgr, mods)
	result, errStruct := itp.InitModules()
	g.Tassert(t, errStruct == nil)
	g.Tassert(t, len(result) == 4)
}

func TestFinally(t *testing.T) {

	code := `
let a = 1
try {
    3 / 0
} finally {
    a = 2
}
try {
    3 / 0
} finally {
    a = 3
}
`
	fail(t, code,
		g.NewErrorStruct(
			g.DivideByZeroError(),
			[]string{
				"    at foo.glm:4"}))

	code = `
let a = 1;
try {
	try {
		3 / 0;
	} finally {
		a++;
	}
} finally {
	a++;
}
`
	fail(t, code,
		g.NewErrorStruct(
			g.DivideByZeroError(),
			[]string{
				"    at foo.glm:5"}))

	code = `
let a = 1;
let b = fn() { a++; };
try {
	try {
		3 / 0;
	} finally {
		a++;
		b();
	}
} finally {
	a++;
}
`
	fail(t, code,
		g.NewErrorStruct(
			g.DivideByZeroError(),
			[]string{
				"    at foo.glm:6"}))

	code = `
let a = 1
let b = fn() {
	try {
		try {
			3 / 0
		} finally {
			a++
		}
	} finally {
		a++
	}
}
try {
	b()
} finally {
	a++
}
`
	//mod = testCompile(t, code)
	//fmt.Println("----------------------------")
	//fmt.Println(code)
	//fmt.Println(mod)

	fail(t, code,
		g.NewErrorStruct(
			g.DivideByZeroError(),
			[]string{
				"    at foo.glm:6",
				"    at foo.glm:15"}))

	code = `
let b = fn() {
	try {
	} finally {
		return 1;
	}
	return 2;
};
assert(b() == 1);
`
	mod := testCompile(t, code)
	testInterpret([]*bc.Module{mod})

	code = `
let a = 1;
let b = fn() {
	try {
		try {
		} finally {
			return 1;
		}
		a = 3;
	} finally {
		a = 2;
	}
};
assert(b() == 1);
assert(a == 1);
`
	mod = testCompile(t, code)
	testInterpret([]*bc.Module{mod})

	code = `
try {
	assert(1,2,3);
} finally {
}
`
	fail(t, code,
		g.NewErrorStruct(
			g.ArityError(1, 3),
			[]string{
				"    at foo.glm:3"}))

	code = `
try {
	assert(1,2,3);
} finally {
	1/0;
}
`
	fail(t, code,
		g.NewErrorStruct(
			g.DivideByZeroError(),
			[]string{
				"    at foo.glm:5"}))
}
