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

	// read source
	filename := os.Args[1]
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	source := string(buf)

	// parse
	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner)
	exprMod, err := parser.ParseModule()
	if err != nil {
		panic(err.Error())
	}

	// analyze
	anl := analyzer.NewAnalyzer(exprMod)
	errors := anl.Analyze()
	if len(errors) > 0 {
		panic(fmt.Sprintf("%v", errors))
	}

	// compile
	cmp := compiler.NewCompiler(anl)
	mod := cmp.Compile()

	// interpret
	intp := interpreter.NewInterpreter(mod)
	_, errTrace := intp.Init()
	if errTrace != nil {
		fmt.Printf("%v\n", errTrace.Error)
		fmt.Printf("%v\n", errTrace.StackTrace)
	}

	// run main
	mainVal, mainErr := mod.Contents.GetField(g.MakeStr("main"))
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

		intp = interpreter.NewInterpreter(mod)
		_, errTrace := intp.RunBytecode(mainFn, params)
		if errTrace != nil {
			fmt.Printf("%v\n", errTrace.Error)
			fmt.Printf("%v\n", errTrace.StackTrace)
		}
	}
}
