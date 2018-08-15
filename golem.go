// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	//"sync"

	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/interpreter"
	"github.com/mjarmy/golem-lang/scanner"
)

var version = "0.8.2"

func exitError(e error) {
	fmt.Printf("%s\n", e.Error())
	os.Exit(-1)
}

func exitErrors(errors []error) {
	for _, e := range errors {
		fmt.Printf("%s\n", e.Error())
	}
	os.Exit(-1)
}

// homePath looks up the path of the golem executable
func homePath() string {
	ex, err := os.Executable()
	if err != nil {
		exitError(err)
	}
	return filepath.Dir(ex)
}

//var libModules = make(map[string]g.Module)
//var libMutex = &sync.Mutex{}

//// lookupModule looks up a module by to loadng a plugin from the '$HOME/lib' directory.
//func lookupModule(homePath, name string) (g.Module, g.Error) {
//
//	libMutex.Lock()
//	defer libMutex.Unlock()
//
//	if mod, ok := libModules[name]; ok {
//		return mod, nil
//	}
//
//	// open the plugin
//	p, err := plugin.Open(homePath + "/lib/" + name + "/" + name + ".so")
//	if err != nil {
//		return nil, g.CouldNotLoadModuleError(name, err)
//	}
//
//	// lookup the 'LoadModule' function
//	f, err := p.Lookup("LoadModule")
//	if err != nil {
//		return nil, g.CouldNotLoadModuleError(name, err)
//	}
//	loader := f.(func() (g.Module, g.Error))
//
//	// load the module
//	mod, gerr := loader()
//	if gerr != nil {
//		return nil, gerr
//	}
//	if mod.GetModuleName() != name {
//		return nil, g.CouldNotLoadModuleError(
//			name,
//			fmt.Errorf("Module name mismatch %s != %s", mod.GetModuleName(), name))
//	}
//
//	libModules[name] = mod
//	return mod, nil
//}

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
	return &scanner.Source{
		Name: name,
		Path: path,
		Code: code,
	}, nil
}

func main() {

	if len(os.Args) == 1 {
		fmt.Printf("Golem %s\n", version)
		os.Exit(0)
	}

	homePath := homePath()

	// command line builtins
	builtinMgr := g.NewBuiltinManager(g.CommandLineBuiltins)

	// read source
	source, e := readSourceFromFile(os.Args[1])
	if e != nil {
		exitError(e)
	}

	// compile
	mods, errs := compiler.CompileSourceFully(builtinMgr, source, nil)
	if len(errs) > 0 {
		exitErrors(errs)
	}

	// interpret
	intp := interpreter.NewInterpreter(homePath, builtinMgr, mods)
	_, err := intp.InitModules()
	if err != nil {
		dumpError(intp, err)
		os.Exit(-1)
	}

	// run main, if it exists
	mainVal, mainErr := mods[0].Contents.GetField(intp, g.NewStr("main"))
	if mainErr == nil {
		mainFn, ok := mainVal.(g.BytecodeFunc)
		if !ok {
			exitError(fmt.Errorf("'main' is not a function"))
		}

		// gather up the command line arguments into a single list
		// that will be passed into the main function
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
			exitError(fmt.Errorf("'main' has too many arguments"))
		}

		// evaluate the main function
		_, err := intp.Eval(mainFn, params)
		if err != nil {
			dumpError(intp, err)
			os.Exit(-1)
		}
	}
}
