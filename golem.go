// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/mjarmy/golem-lang/analyzer"
	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/interpreter"
	"github.com/mjarmy/golem-lang/parser"
	"github.com/mjarmy/golem-lang/scanner"
	"io/ioutil"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		panic("No source file was specified")
	}

	// builtins
	builtInMgr := g.NewBuiltinManager(g.CommandLineBuiltins)

	// read source
	filename := os.Args[1]
	buf, e := ioutil.ReadFile(filename)
	if e != nil {
		panic(e)
	}
	source := string(buf)

	// parse
	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner, builtInMgr.Contains)
	exprMod, e := parser.ParseModule()
	if e != nil {
		panic(e.Error())
	}

	// analyze
	anl := analyzer.NewAnalyzer(exprMod)
	errors := anl.Analyze()
	if len(errors) > 0 {
		panic(fmt.Sprintf("%v", errors))
	}

	// compile
	cmp := compiler.NewCompiler(anl, builtInMgr)
	mod := cmp.Compile()

	// interpret
	intp := interpreter.NewInterpreter(mod, builtInMgr)
	_, err := intp.Init()
	if err != nil {
		dumpError(intp, err)
		os.Exit(-1)
	}

	// run main
	mainVal, mainErr := mod.Contents.GetField(intp, g.MakeStr("main"))
	if mainErr == nil {
		mainFn, ok := mainVal.(g.BytecodeFunc)
		if !ok {
			panic("'main' is not a function")
		}

		params := []g.Value{}
		arity := mainFn.Template().Arity
		if arity == 1 {
			osArgs := os.Args[2:]
			args := make([]g.Value, len(osArgs), len(osArgs))
			for i, a := range osArgs {
				args[i] = g.MakeStr(a)
			}
			params = append(params, g.NewList(args))
		} else if arity > 1 {
			panic("'main' has too many arguments")
		}

		_, err := intp.Eval(mainFn, params)
		if err != nil {
			dumpError(intp, err)
			os.Exit(-1)
		}
	}
}

func dumpError(cx g.Context, err g.Error) {
	fmt.Printf("Error: %s\n", err.Error())

	v, e := err.Struct().GetField(cx, g.MakeStr("stackTrace"))
	if e != nil {
		return
	}
	ls, ok := v.(g.List)
	if !ok {
		return
	}

	itr := ls.NewIterator(cx)
	for itr.IterNext().BoolVal() {
		v, e = itr.IterGet()
		if e != nil {
			return
		}
		fmt.Printf("%s\n", v.ToStr(cx))
	}
}
