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

	g "github.com/mjarmy/golem-lang/core"
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
	mod := g.NewNativeModule("foo", st)

	// create an importer that can import the module
	importer := interpreter.NewImporter([]g.Module{mod})

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
// main
//--------------------------------------------------------------

func main() {
	example1()
	example2()
	example3()
	example4()
	example5()
}
