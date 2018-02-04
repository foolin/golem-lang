// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	//"reflect"
	"testing"
	//	//"github.com/mjarmy/golem-lang/analyzer"
	//	"github.com/mjarmy/golem-lang/compiler"
	//	g "github.com/mjarmy/golem-lang/core"
	//	"github.com/mjarmy/golem-lang/parser"
	//	"github.com/mjarmy/golem-lang/scanner"
)

//func newCompiler(source string) compiler.Compiler {
//	scanner := scanner.NewScanner(source)
//	parser := parser.NewParser(scanner, builtInMgr.Contains)
//	mod, err := parser.ParseModule()
//	if err != nil {
//		panic(err.Error())
//	}
//	anl := analyzer.NewAnalyzer(mod)
//	errors := anl.Analyze()
//	if len(errors) > 0 {
//		panic(fmt.Sprintf("%v", errors))
//	}
//
//	return compiler.NewCompiler(anl, builtInMgr)
//}
//
//func interpret(mod *g.BytecodeModule) {
//	intp := NewInterpreter(mod, builtInMgr)
//	_, err := intp.Init()
//	if err != nil {
//		fmt.Printf("%v\n", err)
//		panic("interpret failed")
//	}
//}

func TestRegexp(t *testing.T) {

	//	source := `
	//import regexp;
	//
	//let pattern = regexp.compile("abc");
	//assert([pattern.match("xyzabcdef"), pattern.match("123")] == [true, false]);
	//`
	//	mod := newCompiler(source).Compile()
	//	interpret(mod)
}

//	regex := InitRegexpModule()
//
//	compile, err := regex.GetContents().GetField(nil, g.NewStr("compile"))
//	tassert(t, compile != nil && err == nil)
//	fnCompile := compile.(g.NativeFunc)
//
//	pattern, err := fnCompile.Invoke(nil, []g.Value{g.NewStr(`^[a-z]+\[[0-9]+\]$`)})
//	tassert(t, pattern != nil && err == nil)
//
//	match, err := pattern.GetField(nil, g.NewStr("match"))
//	tassert(t, match != nil && err == nil)
//	fnMatch := match.(g.NativeFunc)
//
//	result, err := fnMatch.Invoke(nil, []g.Value{g.NewStr("foo[123]")})
//	ok(t, result, err, g.True)
//
//	result, err = fnMatch.Invoke(nil, []g.Value{g.NewStr("456")})
//	ok(t, result, err, g.False)
//
//	pattern, err = fnCompile.Invoke(nil, []g.Value{g.NewStr("\\")})
//	tassert(t, pattern == nil && err.Error() ==
//		"RegexpError: error parsing regexp: trailing backslash at end of expression: ``")
//}

////////////////////////////////////////////////////////////////
// TODO put in interpreter_test

//func TestImport(t *testing.T) {
//
//	anl := newAnalyzer("import sys; let b = 2;")
//	errors := anl.Analyze()
//
//	//fmt.Println(ast.Dump(anl.Module()))
//	//fmt.Println(errors)
//
//	ok(t, anl, errors, `
//FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
//.   BlockNode
//.   .   ImportStmt
//.   .   .   IdentExpr(sys,(0,true,false))
//.   .   LetStmt
//.   .   .   IdentExpr(b,(1,false,false))
//.   .   .   BasicExpr(Int,"2")
//`)
//
//	errors = newAnalyzer("import sys; let sys = 2;").Analyze()
//	fail(t, errors, "[Symbol 'sys' is already defined]")
//
//	errors = newAnalyzer("import sys; sys = 2;").Analyze()
//	fail(t, errors, "[Symbol 'sys' is constant]")
//
//	errors = newAnalyzer("import foo;").Analyze()
//	fail(t, errors, "[Module 'foo' is not defined]")
//}
//
