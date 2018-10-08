// Copyright 2018 The Golem Language Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

//--------------------------------------------------------------
// This program contains examples of how to sandbox a Golem
// interpreter in a Go program.
//--------------------------------------------------------------

import (
	"fmt"

	// Its easier to read the code if we alias the 'core' package.
	// By convention we use 'g', as in 'golem'.
	g "github.com/mjarmy/golem-lang/core"

	"github.com/mjarmy/golem-lang/core/bytecode"
	"github.com/mjarmy/golem-lang/interpreter"
	"github.com/mjarmy/golem-lang/lib"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//--------------------------------------------------------------
// evaluate an expression
//--------------------------------------------------------------

func example1() {

	// evaluate
	val, err := interpreter.EvalCode("1 + 0", nil, nil)
	check(err)

	// print result
	n := val.(g.Int)
	fmt.Printf("example 1: %d\n", n.ToInt())
}

//--------------------------------------------------------------
// evaluate an expression, using a builtin value
//--------------------------------------------------------------

func example2() {

	// create a builtin value named "a"
	a := &g.Builtin{Name: "a", Value: g.NewInt(1)}
	builtins := []*g.Builtin{a}

	// evaluate
	val, err := interpreter.EvalCode("a + 1", builtins, nil)
	check(err)

	// print result
	n := val.(g.Int)
	fmt.Printf("example 2: %d\n", n.ToInt())
}

//--------------------------------------------------------------
// evaluate an expression, using a builtin value and
// an imported module
//--------------------------------------------------------------

func example3() {

	// create a builtin value named "a"
	a := &g.Builtin{Name: "a", Value: g.NewInt(2)}
	builtins := []*g.Builtin{a}

	// create a module called "foo", with a single field named "b"
	st, err := g.NewStruct(map[string]g.Field{
		"b": g.NewField(g.NewInt(1)),
	})
	check(err)
	foo := g.NewNativeModule("foo", st)

	// create an importer that can import the module
	importer := interpreter.NewImporter([]g.Module{foo})

	// evaluate
	code := `
import foo
a + foo.b
`
	val, err := interpreter.EvalCode(code, builtins, importer)
	check(err)

	// print result
	n := val.(g.Int)
	fmt.Printf("example 3: %d\n", n.ToInt())
}

//--------------------------------------------------------------
// compile an expression, and intepret it twice
//--------------------------------------------------------------

func example4() {

	// create a builtin value named "a"
	a := &g.Builtin{Name: "a", Value: g.NewInt(2)}
	builtins := []*g.Builtin{a}

	// compile some code into a module
	mod, err := interpreter.CompileCode("2 + a", builtins)
	check(err)

	// create an interpreter
	itp := interpreter.NewInterpreter(builtins, nil)

	// interpret
	v1, err := itp.EvalModule(mod)
	check(err)

	// Change the builtin value.
	// Note: its OK to change the value of a builtin in between runs
	// of the interpreter,  but you cannot add or remove builtins
	// from a compiled module.
	a.Value = g.MustStr("xyz")

	// interpret again
	v2, err := itp.EvalModule(mod)
	check(err)

	// print result
	n := v1.(g.Int)
	s := v2.(g.Str)
	fmt.Printf("example 4: %d %s\n", n.ToInt(), s.String())
}

//--------------------------------------------------------------
// Create an interpreter that uses the sandboxed builtins and
// sandboxed library
//--------------------------------------------------------------

func example5() {

	// compile some code
	code := `
import encoding

let fibonacciGenerator = fn() {
    let x = 1
    let y = 1
    return fn() {
        let z = x
        x = y
        y = x + z
        return z
    }
}

let fg = fibonacciGenerator()

let list = []
for i in range(0, 10) {
    list.add(fg())
}
return encoding.json.marshal(list)
`
	mod, err := interpreter.CompileCode(code, g.SandboxBuiltins)
	check(err)

	// Create an interpreter.  Note that we must use the same builtins
	// to run code that were used to compile it.
	itp := interpreter.NewInterpreter(
		g.SandboxBuiltins,
		interpreter.NewImporter(lib.SandboxLibrary))

	// run the code by evaluating the module
	val, err := itp.EvalModule(mod)
	check(err)

	// print result
	s := val.(g.Str)
	fmt.Printf("example 5: %s\n", s.String())
}

//--------------------------------------------------------------
// Create an interpreter that calls a function in a
// compiled module.
//--------------------------------------------------------------

func example6() {

	// compile some code
	code := `
fn foo(x, y) {
	return (x+y, x*y)
}
`
	mod, err := interpreter.CompileCode(code, g.SandboxBuiltins)
	check(err)

	// Create an interpreter.
	itp := interpreter.NewInterpreter(g.SandboxBuiltins, nil)

	// evaluate the module (this must be done exactly once)
	_, err = itp.EvalModule(mod)
	check(err)

	// get the compiled 'foo' function
	fooVal, err := mod.Contents().GetField(itp, "foo")
	check(err)
	foo := fooVal.(bytecode.Func)

	// run the 'foo' function twice
	tuple1, err := itp.EvalBytecode(foo, []g.Value{g.NewInt(2), g.NewInt(3)})
	check(err)
	tuple2, err := itp.EvalBytecode(foo, []g.Value{g.NewInt(4), g.NewInt(5)})
	check(err)

	// print the string value of a list of the two tuples
	list := g.NewList([]g.Value{tuple1, tuple2})
	s, err := list.ToStr(itp)
	check(err)
	fmt.Printf("example 6: %s\n", s.String())
}

//--------------------------------------------------------------
// Create an interpreter that calls a function in a
// compiled module.  The function modifys a captured top-level
// variable in the module.
//--------------------------------------------------------------

func example7() {

	// compile some code
	code := `
let tuples = []
fn foo(x, y) {
	let t = (x+y, x*y)
	tuples.add(t)
}
`
	mod, err := interpreter.CompileCode(code, g.SandboxBuiltins)
	check(err)

	// Create an interpreter.
	itp := interpreter.NewInterpreter(g.SandboxBuiltins, nil)

	// evaluate the module (this must be done once)
	_, err = itp.EvalModule(mod)
	check(err)

	// get the compiled 'foo' function
	fooVal, err := mod.Contents().GetField(itp, "foo")
	check(err)
	foo := fooVal.(bytecode.Func)

	// run the 'foo' function twice
	_, err = itp.EvalBytecode(foo, []g.Value{g.NewInt(2), g.NewInt(3)})
	check(err)
	_, err = itp.EvalBytecode(foo, []g.Value{g.NewInt(4), g.NewInt(5)})
	check(err)

	// get the list of tuples
	tuplesVal, err := mod.Contents().GetField(itp, "tuples")
	check(err)
	tuples := tuplesVal.(g.List)

	// print the tuples
	s, err := tuples.ToStr(itp)
	check(err)
	fmt.Printf("example 7: %s\n", s.String())
}

//--------------------------------------------------------------
// Create a module that has a native Go function.
//--------------------------------------------------------------

func example8() {

	// create a module called "foo", with a single field named "cube",
	// that is a g.NativeFunc which accepts a single int and returns its cube.
	cube := g.NewFixedNativeFunc(
		[]g.Type{g.IntType}, false,
		func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
			n := params[0].(g.Int).ToInt()
			return g.NewInt(n * n * n), nil
		})
	st, err := g.NewStruct(map[string]g.Field{
		"cube": g.NewReadonlyField(cube),
	})
	check(err)
	foo := g.NewNativeModule("foo", st)

	// create an importer that can import the module
	importer := interpreter.NewImporter([]g.Module{foo})

	// evaluate
	code := `
import foo
return foo.cube(2)
`
	val, err := interpreter.EvalCode(code, nil, importer)
	check(err)

	// print result
	n := val.(g.Int)
	fmt.Printf("example 8: %d\n", n.ToInt())
}

//--------------------------------------------------------------
// Use methods to efficiently wrap a native Go struct inside
// a module.
//--------------------------------------------------------------

func example9() {

	// Golem has first-class functions.
	//
	// However, sometimes when we create a Golem struct that wraps some native Go code,
	// it would be much more efficient to avoid having to create all of the functions
	// on the struct unless we really need them. We can use "Methods" to provide this
	// optimization.

	// define a simple struct
	type zork struct {
		a int64
		b int64
	}

	// Define some methods on the struct.
	methods := map[string]g.Method{

		"a": g.NewWrapperMethod(
			func(self interface{}) g.Value {
				z := self.(*zork)
				return g.NewInt(z.a)
			}),

		"b": g.NewWrapperMethod(
			func(self interface{}) g.Value {
				z := self.(*zork)
				return g.NewInt(z.b)
			}),

		"add": g.NewFixedMethod(
			[]g.Type{}, false,
			func(self interface{}, ev g.Eval, params []g.Value) (g.Value, g.Error) {
				z := self.(*zork)
				return g.NewInt(z.a + z.b), nil
			}),

		"mul": g.NewFixedMethod(
			[]g.Type{}, false,
			func(self interface{}, ev g.Eval, params []g.Value) (g.Value, g.Error) {
				z := self.(*zork)
				return g.NewInt(z.a * z.b), nil
			}),
	}

	// Create a g.NativeFunc which accepts two ints and returns a "zork" struct.
	newZork := g.NewFixedNativeFunc(
		[]g.Type{g.IntType, g.IntType}, false,
		func(ev g.Eval, params []g.Value) (g.Value, g.Error) {

			a := params[0].(g.Int).ToInt()
			b := params[1].(g.Int).ToInt()

			return g.NewMethodStruct(&zork{a, b}, methods)
		})

	// create a builtin value named "newZork" for our function.
	builtins := []*g.Builtin{{"newZork", newZork}}

	// compile some code
	code := `
let z1 = newZork(2, 3)
let z2 = newZork(4, 5)
return [
	(z1.a(), z1.b(), z1.add(), z1.mul()),
	(z2.a(), z2.b(), z2.add(), z2.mul())
]`
	mod, err := interpreter.CompileCode(code, builtins)
	check(err)

	// create an interpreter
	itp := interpreter.NewInterpreter(builtins, nil)

	// interpret
	val, err := itp.EvalModule(mod)
	check(err)
	list := val.(g.List)

	// print result
	fmt.Printf("example 9: ")
	for _, v := range list.Values() {
		s, err := v.ToStr(itp)
		check(err)
		fmt.Printf("%s ", s.String())
	}
	fmt.Printf("\n")
}

//--------------------------------------------------------------
// main
//--------------------------------------------------------------

func main() {
	example1()
	example2()
	example3()
	example4()
	example5()
	example6()
	example7()
	example8()
	example9()
}
