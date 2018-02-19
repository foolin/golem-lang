// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/mjarmy/golem-lang/analyzer"
	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/interpreter"
	"github.com/mjarmy/golem-lang/lib"
	"github.com/mjarmy/golem-lang/parser"
	"github.com/mjarmy/golem-lang/scanner"
	"io/ioutil"
	"os"
)

var version = "0.8.2"

func exit(msg string) {
	fmt.Printf("%s\n", msg)
	os.Exit(-1)
}

func main() {

	if len(os.Args) == 1 {
		fmt.Printf("Golem %s\n", version)
		os.Exit(0)
	}

	// read source
	filename := os.Args[1]
	buf, e := ioutil.ReadFile(filename)
	if e != nil {
		panic(e)
	}
	source := string(buf)

	// scan
	scanner := scanner.NewScanner(source)

	// command line builtins
	builtInMgr := g.NewBuiltinManager(g.CommandLineBuiltins)

	// parse
	parser := parser.NewParser(scanner, builtInMgr.Contains)
	astMod, e := parser.ParseModule()
	if e != nil {
		exit(e.Error())
	}

	// analyze
	anl := analyzer.NewAnalyzer(astMod)
	errors := anl.Analyze()
	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Printf("%s\n", e)
			os.Exit(-1)
		}
	}

	// compile
	cmp := compiler.NewCompiler(anl, builtInMgr)
	mod := cmp.Compile()

	// interpret with modules from standard library
	intp := interpreter.NewInterpreter(mod, builtInMgr, lib.LookupModule)
	_, err := intp.Init()
	if err != nil {
		dumpError(intp, err)
		os.Exit(-1)
	}

	// run main
	mainVal, mainErr := mod.Contents.GetField(intp, g.NewStr("main"))
	if mainErr == nil {
		mainFn, ok := mainVal.(g.BytecodeFunc)
		if !ok {
			exit("'main' is not a function")
		}

		params := []g.Value{}
		arity := mainFn.Template().Arity
		if arity == 1 {
			osArgs := os.Args[2:]
			args := make([]g.Value, len(osArgs))
			for i, a := range osArgs {
				args[i] = g.NewStr(a)
			}
			params = append(params, g.NewList(args))
		} else if arity > 1 {
			exit("'main' has too many arguments")
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

	v, e := err.Struct().GetField(cx, g.NewStr("stackTrace"))
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
