// Copyright 2018 Mike Jarmy. All rights reserved.  Use of this
// source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// +build js

package main

import (
	"bytes"
	"syscall/js"

	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/interpreter"
)

func getElem(id js.Value) js.Value {
	return js.Global().Get("document").Call("getElementById", id.String())
}

func setOut(out js.Value, s string) {
	out.Set("innerHTML", s)
}

// makePrintBuiltins returns 'print' and 'println' functions that write to a bytes.Buffer
func makePrintBuiltins(buf *bytes.Buffer) []*g.Builtin {

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

	return []*g.Builtin{
		{"print", print},
		{"println", println},
	}
}

func interpret(i []js.Value) {

	in := getElem(i[0])
	out := getElem(i[1])

	var buf bytes.Buffer

	// make builtins
	builtins := append(g.StandardBuiltins, makePrintBuiltins(&buf)...)

	// compile
	code := in.Get("value").String()
	mod, err := interpreter.CompileCode(code, builtins, nil)
	if err != nil {
		setOut(out, err.Error())
		return
	}

	// evaluate
	itp := interpreter.NewInterpreter(builtins, nil)
	val, err := itp.EvalModule(mod)
	if err != nil {
		if es, ok := err.(interpreter.ErrorStruct); ok {
			setOut(out, es.String())
			return
		}
		setOut(out, err.Error())
		return
	}

	// turn result into string
	if val != g.Null {
		s, e := val.ToStr(itp)
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
