// Copyright 2018 Mike Jarmy. All rights reserved.  Use of this
// source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// +build js

package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/interpreter"
	"github.com/mjarmy/golem-lang/scanner"
)

var moduleResolver = func(name string) (*scanner.Source, error) {
	// there are no modules that can be resolved inside of WASM
	return nil, fmt.Errorf("Cannot resolve module '%s'", name)
}

// makeBuiltinMgr returns the standard builtin functions,
// plues custom 'print' and 'println' functions.
func makeBuiltinMgr(buf *bytes.Buffer) g.BuiltinManager {

	var print = g.NewVariadicNativeFunc(
		[]g.Type{}, g.AnyType, true,
		func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
			for _, v := range params {
				s, err := v.ToStr(ev)
				if err != nil {
					return nil, err
				}
				buf.WriteString(s.String())
			}

			return g.Null, nil
		})

	var println = g.NewVariadicNativeFunc(
		[]g.Type{}, g.AnyType, true,
		func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
			for _, v := range params {
				s, err := v.ToStr(ev)
				if err != nil {
					return nil, err
				}
				buf.WriteString(s.String())
			}
			buf.WriteString("\n")

			return g.Null, nil
		})

	return g.NewBuiltinManager(append(
		g.StandardBuiltins,
		[]*g.BuiltinEntry{
			{"print", print},
			{"println", println},
		}...))
}

func errorString(errors []error) string {
	var buf bytes.Buffer
	for _, e := range errors {
		buf.WriteString(fmt.Sprintf("%s\n", e.Error()))
	}
	return buf.String()
}

func getElem(id js.Value) js.Value {
	return js.Global().Get("document").Call("getElementById", id.String())
}

func setOut(out js.Value, s string) {
	out.Set("innerHTML", s)
}

func interpret(i []js.Value) {

	in := getElem(i[0])
	out := getElem(i[1])

	var buf bytes.Buffer

	// compile
	code := in.Get("value").String()
	builtinMgr := makeBuiltinMgr(&buf)
	src := &scanner.Source{Name: "<src>", Path: "<src>", Code: code}
	mods, errs := compiler.CompileSourceFully(builtinMgr, moduleResolver, src)
	if len(errs) > 0 {
		setOut(out, errorString(errs))
		return
	}

	// interpret
	itp := interpreter.NewInterpreter(builtinMgr, mods)
	result, es := itp.InitModules()
	if es != nil {
		setOut(out, es.String())
		return
	}

	// turn result into string
	if result[0] != g.Null {
		s, e := result[0].ToStr(itp)
		if e != nil {
			setOut(out, e.Error())
			return
		}
		buf.WriteString(s.String())
	}

	// done
	setOut(out, buf.String())
}

func main() {

	js.Global().Set("interpret", js.NewCallback(interpret))

	// block forever
	select {}
}
