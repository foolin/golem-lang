// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	"fmt"
	"testing"

	"github.com/mjarmy/golem-lang/analyzer"
	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/interpreter"
	"github.com/mjarmy/golem-lang/parser"
	"github.com/mjarmy/golem-lang/scanner"
)

func interpret(source string) {

	// scan
	scanner := scanner.NewScanner(source)

	// parse
	builtInMgr := g.NewBuiltinManager(g.CommandLineBuiltins)
	parser := parser.NewParser(scanner, builtInMgr.Contains)
	astMod, err := parser.ParseModule()
	if err != nil {
		panic(err.Error())
	}

	// analyze
	anl := analyzer.NewAnalyzer(astMod)
	errors := anl.Analyze()
	if len(errors) > 0 {
		panic(fmt.Sprintf("%v", errors))
	}

	// compile
	cmp := compiler.NewCompiler(anl, builtInMgr)
	mod := cmp.Compile()

	// interpret
	intp := interpreter.NewInterpreter(mod, builtInMgr, LookupModule)
	_, err = intp.Init()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic("interpret failed")
	}
}

func TestRegexp(t *testing.T) {

	interpret(`
import regexp
let rgx = regexp.compile("abc")
assert(
	[rgx.matchString("xyzabcdef"), rgx.matchString("123")] == 
	[true, false])

rgx = regexp.compile("^[a-z]+\\[[0-9]+\\]$")
assert(
	[rgx.matchString("foo[123]"), rgx.matchString("456")] == 
	[true, false])

try {
	rgx = regexp.compile("\\")
	assert(false)
} catch e {
	assert(e.kind == 'RegexpError')
	assert(e.msg == 'error parsing regexp: trailing backslash at end of expression: ` + "``" + `')
}
`)
}
