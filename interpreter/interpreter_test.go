// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"reflect"
	"testing"

	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/scanner"
)

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
		panic("tassert")
	}
}

var builtins = []*g.Builtin{
	{"assert", g.BuiltinAssert},
	{"println", g.BuiltinPrintln},
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

func TestEvalCode(t *testing.T) {

	val, err := EvalCode("1 + 2", nil, nil)
	ok(t, val, err, g.NewInt(3))

	blt := []*g.Builtin{{"a", g.NewInt(2)}}
	val, err = EvalCode("1 + a", blt, nil)
	ok(t, val, err, g.NewInt(3))

	stc, err := g.NewStruct(map[string]g.Field{
		"b": g.NewField(g.NewInt(1)),
	})
	tassert(t, err == nil)
	var mod g.Module = g.NewNativeModule("foo", stc)

	val, err = EvalCode("import foo; foo.b + a", blt, []g.Module{mod})
	ok(t, val, err, g.NewInt(3))
}

func TestCompileCode(t *testing.T) {

	blt := []*g.Builtin{{"a", g.NewInt(2)}}

	mod, err := CompileCode("1 + a", blt, nil)
	tassert(t, err == nil)

	val, err := NewInterpreter(blt, nil).EvalModule(mod)
	ok(t, val, err, g.NewInt(3))

	blt[0].Value = g.NewInt(4)
	val, err = NewInterpreter(blt, nil).EvalModule(mod)
	ok(t, val, err, g.NewInt(5))
}

func TestModuleContents(t *testing.T) {

	code := `
let a = 0
const b = "xyz"
fn c() {}
`
	mod, err := CompileCode(code, builtins, nil)
	_, es := NewInterpreter(builtins, nil).EvalModule(mod)
	tassert(t, es == nil)

	stc := mod.Contents()

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
	mod, err = CompileCode(code, builtins, nil)
	val, es = NewInterpreter(builtins, nil).EvalModule(mod)
	tassert(t, es == nil)
	ok(t, val, nil, g.NewInt(3))

	stc = mod.Contents()
	val, err = stc.GetField(nil, "a")
	ok(t, val, err, g.One)
	val, err = stc.GetField(nil, "b")
	ok(t, val, err, g.Null)
}

//func TestImport(t *testing.T) {
//
//	srcMain := &scanner.Source{Name: "foo", Path: "foo.glm", Code: `
//import a, b, c
//assert([1, 2, 3] == [a.x, b.y, c.z])
//`}
//	sourceMap := map[string]*scanner.Source{
//		"a": &scanner.Source{Name: "a", Path: "a.glm", Code: "import c; let x = 1;"},
//		"b": &scanner.Source{Name: "b", Path: "b.glm", Code: "import c; let y = c.z - 1;"},
//		"c": &scanner.Source{Name: "c", Path: "c.glm", Code: "let z = 3;"},
//	}
//	resolver := func(moduleName string) (*scanner.Source, error) {
//		if src, ok := sourceMap[moduleName]; ok {
//			return src, nil
//		}
//		return nil, g.UndefinedModule(moduleName)
//	}
//
//	mods, errs := compiler.CompileSourceFully(builtins, resolver, srcMain)
//	tassert(t, errs == nil)
//	tassert(t, len(mods) == 4)
//	tassert(t, mods[0].Name() == "foo")
//	tassert(t, mods[1].Name() == "a")
//	tassert(t, mods[2].Name() == "b")
//	tassert(t, mods[3].Name() == "c")
//
//	itp := NewInterpreter(builtins, mods)
//	result, es := itp.InitModules()
//	tassert(t, es == nil)
//	tassert(t, len(result) == 4)
//}

//--------------------------------------------------------------
//--------------------------------------------------------------
//--------------------------------------------------------------

func failInterp(t *testing.T, code string, expect ErrorStruct) {

	source := &scanner.Source{Name: "foo", Path: "foo.glm", Code: code}
	_, mod, err := compiler.CompileSource(source, builtins)
	tassert(t, err == nil)

	intp := NewInterpreter(builtins, nil)
	_, es := intp.EvalModule(mod)
	tassert(t, es != nil)

	//dumpErrorStruct("failInterp", es)
	tassert(t, reflect.DeepEqual(es, expect))
}

func TestStackTrace(t *testing.T) {

	code := `
		1/0
		`
	failInterp(t, code,
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:2"}))

	code = `
		let a = (|| => 1/0)
		a()
		`
	failInterp(t, code,
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:2",
				"    at foo.glm:3"}))

	code = `
		let s = struct {
			q: prop { ||=> 1/0 }
		}
		let a = s.q
		`
	failInterp(t, code,
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:3",
				"    at foo.glm:5"}))

	code = `
		[1, 2, 3].map(
			|e| => 1/0)
		`
	failInterp(t, code,
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:3",
				"    at foo.glm:2"}))

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
	failInterp(t, code,
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:5",
				"    at foo.glm:4",
				"    at foo.glm:9"}))

	code = `
		let s = struct {
			q: prop { ||=> 1/0 }
		}

		let a = [1, 2, 3].map(
			|e| => s.q)
		`
	failInterp(t, code,
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:3",
				"    at foo.glm:7",
				"    at foo.glm:6"}))

	code = `
		fn b() {
			a()
		}

		fn a() {
			let s = struct {
				q: prop { ||=> 1/0 }
			}

			let ls = [1, 2, 3].map(
				|e| => s.q)
		}

		fn c() {
			b()
		}

		c()
		`
	failInterp(t, code,
		newErrorStruct(
			g.DivideByZero(),
			[]string{
				"    at foo.glm:8",
				"    at foo.glm:12",
				"    at foo.glm:11",
				"    at foo.glm:3",
				"    at foo.glm:16",
				"    at foo.glm:19"}))
}

//func okInterp(t *testing.T, mods []*bc.Module) {
//	intp := NewInterpreter(builtins, mods)
//	_, es := intp.InitModules()
//	if es != nil {
//		panic(es)
//	}
//}
//
//func TestDebug(t *testing.T) {
//
//	debugInterpreter = true
//
//	code := `
//fn fail(func, err) {
//    try {
//        func()
//        assert(false)
//    } catch e {
//        //println(e.error)
//        assert(err == e.error)
//    }
//}
//            let p = 0
//            let q = 0
//
//            fn a() {
//                try {
//                    1/0
//                } catch e {
//                    throw e.error
//                } finally {
//                    q++
//                }
//                assert(false)
//            }
//
//            fn b() {
//                a()
//            }
//
//            fn c() {
//                try {
//                    b()
//                } catch e {
//                    assert(e.error == 'DivideByZero')
//                    p++
//                } finally {
//                    throw 'TestError'
//                }
//                assert(false)
//            }
//
//            fn d() {
//                c()
//            }
//
//            fail(a, 'DivideByZero')
//            assert([p,q] == [0,1])
//
//            //fail(b, 'DivideByZero')
//            //assert([p,q] == [0,2])
//
//            //fail(c, 'TestError')
//            //assert([p,q] == [1,3])
//
//            //fail(d, 'TestError')
//            //assert([p,q] == [2,4])
//`
//	okInterp(t, testCompile(t, code))
//}
