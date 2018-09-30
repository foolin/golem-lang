// Copyright 2018 The Golem Language Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
	"github.com/mjarmy/golem-lang/interpreter"
	"github.com/mjarmy/golem-lang/lib"
	"github.com/mjarmy/golem-lang/scanner"
)

var version = "0.1.0"

func exitError(e error) {
	fmt.Printf("%s\n", e.Error())
	os.Exit(-1)
}

func exitInterpreter(es interpreter.ErrorStruct) {
	fmt.Printf(es.String())
	os.Exit(-1)
}

func readSourceFromFile(filename string) (*scanner.Source, error) {

	// code
	buf, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	code := string(buf)

	// path
	path, e := filepath.Abs(filename)
	if e != nil {
		return nil, e
	}

	// name
	ext := filepath.Ext(path)
	if ext != ".glm" {
		return nil, fmt.Errorf("Golem source file '%s' does not end in '.glm'", path)
	}
	name := strings.TrimSuffix(filepath.Base(path), ext)

	// done
	src := &scanner.Source{
		Name: name,
		Path: path,
		Code: code,
	}
	return src, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func makeResolver(localDir string, builtins []*g.Builtin) interpreter.Resolver {
	return func(name string) (*bc.Module, error) {

		// check local dir
		localPath := localDir + "/" + name + ".glm"
		if !fileExists(localPath) {
			return nil, fmt.Errorf("Cannot resolve module '%s'", name)
		}

		src, err := readSourceFromFile(localPath)
		if err != nil {
			return nil, err
		}

		return compiler.CompileSource(src, builtins)
	}
}

func commandLineArguments() g.List {

	osArgs := os.Args[2:]
	args := make([]g.Value, len(osArgs))
	for i, o := range osArgs {
		a, e := g.NewStr(o)
		if e != nil {
			exitError(e)
			return nil // can't get here
		}
		args[i] = a
	}

	return g.NewList(args)
}

func main() {

	if len(os.Args) == 1 {
		fmt.Printf("Golem %s\n", version)
		os.Exit(0)
	}

	//-------------------------------------------------------------
	// setup
	//-------------------------------------------------------------

	// local directory
	localDir, e := os.Getwd()
	if e != nil {
		exitError(e)
	}

	// use all of the available builtin functions
	builtins := append(g.StandardBuiltins, g.UnsandboxedBuiltins...)

	// resolve bytecode.Modules from the local dir
	resolver := makeResolver(localDir, builtins)

	// use the entire standard library
	library := interpreter.NewLibrary(lib.StandardLibrary, resolver)

	//-------------------------------------------------------------
	// parse, compile, interpret
	//-------------------------------------------------------------

	// read source
	src, e := readSourceFromFile(os.Args[1])
	if e != nil {
		exitError(e)
	}

	// compile
	mod, err := compiler.CompileSource(src, builtins)
	if err != nil {
		exitError(err)
	}

	// interpret
	itp := interpreter.NewInterpreter(builtins, library)
	_, es := itp.EvalModule(mod)
	if es != nil {
		exitInterpreter(es)
	}

	//-------------------------------------------------------------
	// run main() if it exists
	//-------------------------------------------------------------

	// find main
	main, err := mod.Contents().GetField(itp, "main")
	if err != nil {
		// its ok for there to be no 'main'
		if err.Error() == "NoSuchField: Field 'main' not found" {
			return
		}
		exitError(err)
	}

	// make sure that main is a function
	mainFn, ok := main.(bc.Func)
	if !ok {
		exitError(fmt.Errorf("'main' is not a function"))
	}

	// make sure that the main function takes exactly one parameter
	expected := g.Arity{Kind: g.FixedArity, Required: 1, Optional: 0}
	if mainFn.Arity() != expected {
		exitError(fmt.Errorf("ArityMismatch: main function must have 1 parameter"))
	}

	// turn the command line arguments into a List-of-Str
	argList := commandLineArguments()

	// interpret the main function
	_, es = itp.EvalBytecode(mainFn, []g.Value{argList})
	if es != nil {
		exitInterpreter(es)
	}
}
