// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	"github.com/mjarmy/golem-lang/analyzer"
	"github.com/mjarmy/golem-lang/ast"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
	"github.com/mjarmy/golem-lang/parser"
	"github.com/mjarmy/golem-lang/scanner"
)

// ModuleResolver resolves a module name into a Source
type ModuleResolver func(moduleName string) (*scanner.Source, error)

// CompileSourceFully compiles all of the Modules needed to run the program
// that is defined in the provided Source.
func CompileSourceFully(
	builtinMgr BuiltinManager,
	resolver ModuleResolver,
	source *scanner.Source) ([]*bc.Module, []error) {

	sources := []*scanner.Source{source}
	sourceSet := map[string]bool{source.Name: true}
	result := []*bc.Module{}

	idx := 0
	for idx < len(sources) {

		// compile
		astMod, mod, errs := CompileSource(builtinMgr, sources[idx])
		if errs != nil {
			return nil, errs
		}
		result = append(result, mod)

		// add imported
		for _, impName := range astMod.Imports() {
			if _, ok := sourceSet[impName]; !ok {
				impSrc, err := resolver(impName)
				if err != nil {
					return nil, []error{err}
				}
				sources = append(sources, impSrc)
				sourceSet[impName] = true
			}
		}

		// done
		idx++
	}

	return result, nil
}

// CompileSource compiles a Module from Source
func CompileSource(
	builtinMgr BuiltinManager,
	source *scanner.Source) (*ast.Module, *bc.Module, []error) {

	// scan
	scanner, e := scanner.NewScanner(source)
	if e != nil {
		return nil, nil, []error{e}
	}

	// parse
	parser := parser.NewParser(scanner, builtinMgr.Contains)
	astMod, e := parser.ParseModule()
	if e != nil {
		return nil, nil, []error{e}
	}

	// analyze
	anl := analyzer.NewAnalyzer(astMod)
	errors := anl.Analyze()
	if len(errors) > 0 {
		return nil, nil, errors
	}

	// compile
	cmp := NewCompiler(builtinMgr, astMod)
	return astMod, cmp.Compile(), nil
}
