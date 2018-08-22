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
	"sync"

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

// homeDir looks up the path of the golem executable
func homeDir() string {
	ex, err := os.Executable()
	if err != nil {
		exitError(err)
	}
	return filepath.Dir(ex)
}

var sourceMap = make(map[string]*scanner.Source)
var sourceMutex = &sync.Mutex{}

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
	sourceMap[src.Name] = src
	return src, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func makeModuleResolver(homeDir, localDir string) compiler.ModuleResolver {
	return func(name string) (*scanner.Source, error) {

		sourceMutex.Lock()
		defer sourceMutex.Unlock()
		if src, ok := sourceMap[name]; ok {
			return src, nil
		}

		// check std lib first
		libPath := homeDir + "/lib/" + name + "/" + name + ".glm"
		if fileExists(libPath) {
			return readSourceFromFile(libPath)
		}

		// check local dir
		localPath := localDir + "/" + name + ".glm"
		if fileExists(localPath) {
			return readSourceFromFile(localPath)
		}

		// can't find file
		return nil, fmt.Errorf("Cannot resolve module '%s'", name)
	}
}

func main() {

	if len(os.Args) == 1 {
		fmt.Printf("Golem %s\n", version)
		os.Exit(0)
	}

	homeDir := homeDir()
	localDir, e := os.Getwd()
	if e != nil {
		exitError(e)
	}

	// command line builtins
	builtinMgr := g.NewBuiltinManager(g.CommandLineBuiltins)

	// read source
	src, e := readSourceFromFile(os.Args[1])
	if e != nil {
		exitError(e)
	}

	// compile
	resolver := makeModuleResolver(homeDir, localDir)
	mods, errs := compiler.CompileSourceFully(builtinMgr, src, resolver)
	if len(errs) > 0 {
		exitErrors(errs)
	}

	// interpret
	intp := interpreter.NewInterpreter(homeDir, builtinMgr, mods)
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
