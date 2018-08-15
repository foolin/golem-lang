// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	"github.com/mjarmy/golem-lang/analyzer"
	"github.com/mjarmy/golem-lang/ast"
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/parser"
	"github.com/mjarmy/golem-lang/scanner"
)

type ImportResolverFunc func(moduleName string) (*scanner.Source, error)

// CompileSourceFully compiles all of the Modules needed to run the program
// that is defined in the provided Source.
func CompileSourceFully(
	builtinMgr g.BuiltinManager,
	source *scanner.Source,
	resolver ImportResolverFunc) ([]*g.Module, []error) {

	sources := []*scanner.Source{source}
	sourceSet := map[string]bool{source.Name: true}
	result := []*g.Module{}

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
	builtinMgr g.BuiltinManager,
	source *scanner.Source) (*ast.Module, *g.Module, []error) {

	// scan
	scanner := scanner.NewScanner(source)

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
