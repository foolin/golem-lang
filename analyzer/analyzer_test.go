// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this code code is governed by a MIT-style
// license that can be found in the LICENSE file.

package analyzer

import (
	"fmt"
	"github.com/mjarmy/golem-lang/ast"
	"github.com/mjarmy/golem-lang/parser"
	"github.com/mjarmy/golem-lang/scanner"
	"testing"
)

func ok(t *testing.T, mod *ast.Module, errors []error, dump string) {

	if len(errors) != 0 {
		t.Error(errors)
	}

	if "\n"+ast.Dump(mod.InitFunc) != dump {
		t.Error("\n"+ast.Dump(mod.InitFunc), " != ", dump)
	}

}

func fail(t *testing.T, errors []error, expect string) {

	if fmt.Sprintf("%v", errors) != expect {
		t.Error(errors, " != ", expect)
	}
}

var builtins = map[string]bool{
	"print":   true,
	"println": true,
	"str":     true,
	"len":     true,
	"range":   true,
	"assert":  true,
	"merge":   true,
	"chan":    true,
	"typeof":  true,
	"freeze":  true,
	"frozen":  true,
}
var isBuiltIn = func(s string) bool {
	_, ok := builtins[s]
	return ok
}

func mustScanner(source *scanner.Source) *scanner.Scanner {
	scn, err := scanner.NewScanner(source)
	if err != nil {
		panic(err)
	}
	return scn
}

func newModule(code string) *ast.Module {

	ast.InternalResetDebugging()

	scanner := mustScanner(
		&scanner.Source{
			Name: "foo",
			Path: "foo.glm",
			Code: code,
		})
	parser := parser.NewParser(scanner, isBuiltIn)
	mod, err := parser.ParseModule()
	if err != nil {
		panic("analyzer_test: could not parse")
	}

	return mod
}

func TestFlat(t *testing.T) {

	mod := newModule("let a = 1; const b = 2; a = b + 3;")
	errors := NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{a: v(0: a,0,false,false), b: v(1: b,1,true,false)})
.   .   LetStmt
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   BasicExpr(Int,"1")
.   .   ConstStmt
.   .   .   IdentExpr(b,v(1: b,1,true,false))
.   .   .   BasicExpr(Int,"2")
.   .   ExprStmt
.   .   .   AssignmentExpr
.   .   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   .   BinaryExpr("+")
.   .   .   .   .   IdentExpr(b,v(1: b,1,true,false))
.   .   .   .   .   BasicExpr(Int,"3")
`)

	errors = NewAnalyzer(newModule("a;")).Analyze()
	fail(t, errors, "[Symbol 'a' is not defined, at foo.glm:1:1]")

	errors = NewAnalyzer(newModule("let a = 1;const a = 1;")).Analyze()
	fail(t, errors, "[Symbol 'a' is already defined, at foo.glm:1:17]")

	errors = NewAnalyzer(newModule("const a = 1;a = 1;")).Analyze()
	fail(t, errors, "[Symbol 'a' is constant, at foo.glm:1:13]")

	errors = NewAnalyzer(newModule("a = a;")).Analyze()
	fail(t, errors, "[Symbol 'a' is not defined, at foo.glm:1:5 Symbol 'a' is not defined, at foo.glm:1:1]")
}

func TestNested(t *testing.T) {

	code := `
let a = 1
if (true) {
    a = 2
    const b = 2
} else {
    a = 3
    let b = 3
}`
	mod := newModule(code)

	errors := NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:3)
.   BlockNode(Scope defs:{a: v(0: a,0,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   BasicExpr(Int,"1")
.   .   IfStmt
.   .   .   BasicExpr(True,"true")
.   .   .   BlockNode(Scope defs:{b: v(1: b,1,true,false)})
.   .   .   .   ExprStmt
.   .   .   .   .   AssignmentExpr
.   .   .   .   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   .   .   .   BasicExpr(Int,"2")
.   .   .   .   ConstStmt
.   .   .   .   .   IdentExpr(b,v(1: b,1,true,false))
.   .   .   .   .   BasicExpr(Int,"2")
.   .   .   BlockNode(Scope defs:{b: v(2: b,2,false,false)})
.   .   .   .   ExprStmt
.   .   .   .   .   AssignmentExpr
.   .   .   .   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   .   .   .   BasicExpr(Int,"3")
.   .   .   .   LetStmt
.   .   .   .   .   IdentExpr(b,v(2: b,2,false,false))
.   .   .   .   .   BasicExpr(Int,"3")
`)
}

func TestLoop(t *testing.T) {

	mod := newModule("while true { 1 + 2; }")
	errors := NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:0)
.   BlockNode(Scope defs:{})
.   .   WhileStmt
.   .   .   BasicExpr(True,"true")
.   .   .   BlockNode(Scope defs:{})
.   .   .   .   ExprStmt
.   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   BasicExpr(Int,"1")
.   .   .   .   .   .   BasicExpr(Int,"2")
`)

	mod = newModule("while true { 1 + 2; break; continue; }")
	errors = NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:0)
.   BlockNode(Scope defs:{})
.   .   WhileStmt
.   .   .   BasicExpr(True,"true")
.   .   .   BlockNode(Scope defs:{})
.   .   .   .   ExprStmt
.   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   BasicExpr(Int,"1")
.   .   .   .   .   .   BasicExpr(Int,"2")
.   .   .   .   BreakStmt
.   .   .   .   ContinueStmt
`)

	errors = NewAnalyzer(newModule("break;")).Analyze()
	fail(t, errors, "['break' outside of loop, at foo.glm:1:1]")

	errors = NewAnalyzer(newModule("continue;")).Analyze()
	fail(t, errors, "['continue' outside of loop, at foo.glm:1:1]")

	mod = newModule("let a; for b in [] { break; continue; }")
	errors = NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:3)
.   BlockNode(Scope defs:{a: v(0: a,0,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   ForStmt(Scope defs:{#iter0: v(2: #iter0,2,false,false), b: v(1: b,1,false,false)})
.   .   .   IdentExpr(b,v(1: b,1,false,false))
.   .   .   IdentExpr(#iter0,v(2: #iter0,2,false,false))
.   .   .   ListExpr
.   .   .   BlockNode(Scope defs:{})
.   .   .   .   BreakStmt
.   .   .   .   ContinueStmt
`)

	mod = newModule("for (a, b) in [] { }")
	errors = NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:3)
.   BlockNode(Scope defs:{})
.   .   ForStmt(Scope defs:{#iter0: v(2: #iter0,2,false,false), a: v(0: a,0,false,false), b: v(1: b,1,false,false)})
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   IdentExpr(b,v(1: b,1,false,false))
.   .   .   IdentExpr(#iter0,v(2: #iter0,2,false,false))
.   .   .   ListExpr
.   .   .   BlockNode(Scope defs:{})
`)

	mod = newModule(`
for a in [] {
	for b in [] {
	}
}
`)
	errors = NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:4)
.   BlockNode(Scope defs:{})
.   .   ForStmt(Scope defs:{#iter0: v(1: #iter0,1,false,false), a: v(0: a,0,false,false)})
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   IdentExpr(#iter0,v(1: #iter0,1,false,false))
.   .   .   ListExpr
.   .   .   BlockNode(Scope defs:{})
.   .   .   .   ForStmt(Scope defs:{#iter1: v(3: #iter1,3,false,false), b: v(2: b,2,false,false)})
.   .   .   .   .   IdentExpr(b,v(2: b,2,false,false))
.   .   .   .   .   IdentExpr(#iter1,v(3: #iter1,3,false,false))
.   .   .   .   .   ListExpr
.   .   .   .   .   BlockNode(Scope defs:{})
`)

	mod = newModule(`
let a = 1
for a in [] {
}
`)
	errors = NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:3)
.   BlockNode(Scope defs:{a: v(0: a,0,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   BasicExpr(Int,"1")
.   .   ForStmt(Scope defs:{#iter0: v(2: #iter0,2,false,false), a: v(1: a,1,false,false)})
.   .   .   IdentExpr(a,v(1: a,1,false,false))
.   .   .   IdentExpr(#iter0,v(2: #iter0,2,false,false))
.   .   .   ListExpr
.   .   .   BlockNode(Scope defs:{})
`)
}

func TestAssignment(t *testing.T) {

	code := `
let x = struct { a: 0 }
let y = x.a
x.a = 3
x.a++
y--
x[y] = 42
y = x[3]
x[2]++
y.z = x[2]++
let g, h = 5
const i = 6, j
`
	mod := newModule(code)
	errors := NewAnalyzer(mod).Analyze()

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:6)
.   BlockNode(Scope defs:{g: v(2: g,2,false,false), h: v(3: h,3,false,false), i: v(4: i,4,true,false), j: v(5: j,5,true,false), x: v(0: x,0,false,false), y: v(1: y,1,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(x,v(0: x,0,false,false))
.   .   .   StructExpr([a],StructScope defs:{})
.   .   .   .   BasicExpr(Int,"0")
.   .   LetStmt
.   .   .   IdentExpr(y,v(1: y,1,false,false))
.   .   .   FieldExpr(a)
.   .   .   .   IdentExpr(x,v(0: x,0,false,false))
.   .   ExprStmt
.   .   .   AssignmentExpr
.   .   .   .   FieldExpr(a)
.   .   .   .   .   IdentExpr(x,v(0: x,0,false,false))
.   .   .   .   BasicExpr(Int,"3")
.   .   ExprStmt
.   .   .   PostfixExpr("++")
.   .   .   .   FieldExpr(a)
.   .   .   .   .   IdentExpr(x,v(0: x,0,false,false))
.   .   ExprStmt
.   .   .   PostfixExpr("--")
.   .   .   .   IdentExpr(y,v(1: y,1,false,false))
.   .   ExprStmt
.   .   .   AssignmentExpr
.   .   .   .   IndexExpr
.   .   .   .   .   IdentExpr(x,v(0: x,0,false,false))
.   .   .   .   .   IdentExpr(y,v(1: y,1,false,false))
.   .   .   .   BasicExpr(Int,"42")
.   .   ExprStmt
.   .   .   AssignmentExpr
.   .   .   .   IdentExpr(y,v(1: y,1,false,false))
.   .   .   .   IndexExpr
.   .   .   .   .   IdentExpr(x,v(0: x,0,false,false))
.   .   .   .   .   BasicExpr(Int,"3")
.   .   ExprStmt
.   .   .   PostfixExpr("++")
.   .   .   .   IndexExpr
.   .   .   .   .   IdentExpr(x,v(0: x,0,false,false))
.   .   .   .   .   BasicExpr(Int,"2")
.   .   ExprStmt
.   .   .   AssignmentExpr
.   .   .   .   FieldExpr(z)
.   .   .   .   .   IdentExpr(y,v(1: y,1,false,false))
.   .   .   .   PostfixExpr("++")
.   .   .   .   .   IndexExpr
.   .   .   .   .   .   IdentExpr(x,v(0: x,0,false,false))
.   .   .   .   .   .   BasicExpr(Int,"2")
.   .   LetStmt
.   .   .   IdentExpr(g,v(2: g,2,false,false))
.   .   .   IdentExpr(h,v(3: h,3,false,false))
.   .   .   BasicExpr(Int,"5")
.   .   ConstStmt
.   .   .   IdentExpr(i,v(4: i,4,true,false))
.   .   .   BasicExpr(Int,"6")
.   .   .   IdentExpr(j,v(5: j,5,true,false))
`)
}

func TestList(t *testing.T) {

	code := `
let a = ['x'][0]
let b = ['x']
b[0] = 3
b[0]++
`
	mod := newModule(code)
	errors := NewAnalyzer(mod).Analyze()

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{a: v(0: a,0,false,false), b: v(1: b,1,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   IndexExpr
.   .   .   .   ListExpr
.   .   .   .   .   BasicExpr(Str,"x")
.   .   .   .   BasicExpr(Int,"0")
.   .   LetStmt
.   .   .   IdentExpr(b,v(1: b,1,false,false))
.   .   .   ListExpr
.   .   .   .   BasicExpr(Str,"x")
.   .   ExprStmt
.   .   .   AssignmentExpr
.   .   .   .   IndexExpr
.   .   .   .   .   IdentExpr(b,v(1: b,1,false,false))
.   .   .   .   .   BasicExpr(Int,"0")
.   .   .   .   BasicExpr(Int,"3")
.   .   ExprStmt
.   .   .   PostfixExpr("++")
.   .   .   .   IndexExpr
.   .   .   .   .   IdentExpr(b,v(1: b,1,false,false))
.   .   .   .   .   BasicExpr(Int,"0")
`)
}

func TestTry(t *testing.T) {

	code := "let a = 1; try { } catch e { } finally { }"
	mod := newModule(code)
	errors := NewAnalyzer(mod).Analyze()

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{a: v(0: a,0,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   BasicExpr(Int,"1")
.   .   TryStmt(Scope defs:{e: v(1: e,1,true,false)})
.   .   .   BlockNode(Scope defs:{})
.   .   .   IdentExpr(e,v(1: e,1,true,false))
.   .   .   BlockNode(Scope defs:{})
.   .   .   BlockNode(Scope defs:{})
`)

	code = "let a = 1; try { } catch a { } finally { }"
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{a: v(0: a,0,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   BasicExpr(Int,"1")
.   .   TryStmt(Scope defs:{a: v(1: a,1,true,false)})
.   .   .   BlockNode(Scope defs:{})
.   .   .   IdentExpr(a,v(1: a,1,true,false))
.   .   .   BlockNode(Scope defs:{})
.   .   .   BlockNode(Scope defs:{})
`)
}

func TestFormalParams(t *testing.T) {

	errors := NewAnalyzer(newModule("fn(const a, b) { a = 1; };")).Analyze()
	fail(t, errors, "[Symbol 'a' is constant, at foo.glm:1:18]")

	errors = NewAnalyzer(newModule("fn(a, const b) { b = 1; };")).Analyze()
	fail(t, errors, "[Symbol 'b' is constant, at foo.glm:1:18]")

	errors = NewAnalyzer(newModule("fn(a, const b...) { b = 1; };")).Analyze()
	fail(t, errors, "[Symbol 'b' is constant, at foo.glm:1:21]")
}

func TestPureFunction(t *testing.T) {
	code := `
let a = 1
let b = fn(x) {
    let c = fn(y, z) {
        if (y < 33) {
            return y + z + 5
        } else {
            let b = 42
        }
    }
    return c(3)
}`

	mod := newModule(code)
	errors := NewAnalyzer(mod).Analyze()

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{a: v(0: a,0,false,false), b: v(6: b,1,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   BasicExpr(Int,"1")
.   .   LetStmt
.   .   .   IdentExpr(b,v(6: b,1,false,false))
.   .   .   FnExpr(FuncScope defs:{x: v(1: x,0,false,false)} captures:{} numLocals:2)
.   .   .   .   IdentExpr(x,v(1: x,0,false,false))
.   .   .   .   BlockNode(Scope defs:{c: v(5: c,1,false,false)})
.   .   .   .   .   LetStmt
.   .   .   .   .   .   IdentExpr(c,v(5: c,1,false,false))
.   .   .   .   .   .   FnExpr(FuncScope defs:{y: v(2: y,0,false,false), z: v(3: z,1,false,false)} captures:{} numLocals:3)
.   .   .   .   .   .   .   IdentExpr(y,v(2: y,0,false,false))
.   .   .   .   .   .   .   IdentExpr(z,v(3: z,1,false,false))
.   .   .   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   .   .   IfStmt
.   .   .   .   .   .   .   .   .   BinaryExpr("<")
.   .   .   .   .   .   .   .   .   .   IdentExpr(y,v(2: y,0,false,false))
.   .   .   .   .   .   .   .   .   .   BasicExpr(Int,"33")
.   .   .   .   .   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   .   .   .   .   ReturnStmt
.   .   .   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   .   .   .   IdentExpr(y,v(2: y,0,false,false))
.   .   .   .   .   .   .   .   .   .   .   .   .   IdentExpr(z,v(3: z,1,false,false))
.   .   .   .   .   .   .   .   .   .   .   .   BasicExpr(Int,"5")
.   .   .   .   .   .   .   .   .   BlockNode(Scope defs:{b: v(4: b,2,false,false)})
.   .   .   .   .   .   .   .   .   .   LetStmt
.   .   .   .   .   .   .   .   .   .   .   IdentExpr(b,v(4: b,2,false,false))
.   .   .   .   .   .   .   .   .   .   .   BasicExpr(Int,"42")
.   .   .   .   .   ReturnStmt
.   .   .   .   .   .   InvokeExpr
.   .   .   .   .   .   .   IdentExpr(c,v(5: c,1,false,false))
.   .   .   .   .   .   .   BasicExpr(Int,"3")
`)
}

func TestNamedFunc(t *testing.T) {

	code := `
fn a() {
    return b()
}
fn b() {
    return 42
}
`
	mod := newModule(code)
	errors := NewAnalyzer(mod).Analyze()

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{a: v(0: a,0,true,false), b: v(1: b,1,true,false)})
.   .   NamedFnStmt
.   .   .   IdentExpr(a,v(0: a,0,true,false))
.   .   .   FnExpr(FuncScope defs:{} captures:{b: (parent: v(1: b,1,true,false), child v(2: b,0,true,true))} numLocals:0)
.   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   ReturnStmt
.   .   .   .   .   .   InvokeExpr
.   .   .   .   .   .   .   IdentExpr(b,v(2: b,0,true,true))
.   .   NamedFnStmt
.   .   .   IdentExpr(b,v(1: b,1,true,false))
.   .   .   FnExpr(FuncScope defs:{} captures:{} numLocals:0)
.   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   ReturnStmt
.   .   .   .   .   .   BasicExpr(Int,"42")
`)

	errors = NewAnalyzer(newModule("fn a() {}; const a = 1;")).Analyze()
	fail(t, errors, "[Symbol 'a' is already defined, at foo.glm:1:18]")
}

func TestCaptureFunction(t *testing.T) {

	code := `
const accumGen = fn(n) {
    return fn(i) {
        n = n + i
        return n
    }
}
`
	mod := newModule(code)
	errors := NewAnalyzer(mod).Analyze()

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:1)
.   BlockNode(Scope defs:{accumGen: v(3: accumGen,0,true,false)})
.   .   ConstStmt
.   .   .   IdentExpr(accumGen,v(3: accumGen,0,true,false))
.   .   .   FnExpr(FuncScope defs:{n: v(0: n,0,false,false)} captures:{} numLocals:1)
.   .   .   .   IdentExpr(n,v(0: n,0,false,false))
.   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   ReturnStmt
.   .   .   .   .   .   FnExpr(FuncScope defs:{i: v(1: i,0,false,false)} captures:{n: (parent: v(0: n,0,false,false), child v(2: n,0,false,true))} numLocals:1)
.   .   .   .   .   .   .   IdentExpr(i,v(1: i,0,false,false))
.   .   .   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   .   .   ExprStmt
.   .   .   .   .   .   .   .   .   AssignmentExpr
.   .   .   .   .   .   .   .   .   .   IdentExpr(n,v(2: n,0,false,true))
.   .   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   .   IdentExpr(n,v(2: n,0,false,true))
.   .   .   .   .   .   .   .   .   .   .   IdentExpr(i,v(1: i,0,false,false))
.   .   .   .   .   .   .   .   ReturnStmt
.   .   .   .   .   .   .   .   .   IdentExpr(n,v(2: n,0,false,true))
`)

	code = `
let z = 2
const accumGen = fn(n) {
	return fn(i) {
		n = n + i
		n = n + z
		return n
	}
}
	`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{accumGen: v(6: accumGen,1,true,false), z: v(0: z,0,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(z,v(0: z,0,false,false))
.   .   .   BasicExpr(Int,"2")
.   .   ConstStmt
.   .   .   IdentExpr(accumGen,v(6: accumGen,1,true,false))
.   .   .   FnExpr(FuncScope defs:{n: v(1: n,0,false,false)} captures:{z: (parent: v(0: z,0,false,false), child v(4: z,0,false,true))} numLocals:1)
.   .   .   .   IdentExpr(n,v(1: n,0,false,false))
.   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   ReturnStmt
.   .   .   .   .   .   FnExpr(FuncScope defs:{i: v(2: i,0,false,false)} captures:{n: (parent: v(1: n,0,false,false), child v(3: n,0,false,true)), z: (parent: v(4: z,0,false,true), child v(5: z,1,false,true))} numLocals:1)
.   .   .   .   .   .   .   IdentExpr(i,v(2: i,0,false,false))
.   .   .   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   .   .   ExprStmt
.   .   .   .   .   .   .   .   .   AssignmentExpr
.   .   .   .   .   .   .   .   .   .   IdentExpr(n,v(3: n,0,false,true))
.   .   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   .   IdentExpr(n,v(3: n,0,false,true))
.   .   .   .   .   .   .   .   .   .   .   IdentExpr(i,v(2: i,0,false,false))
.   .   .   .   .   .   .   .   ExprStmt
.   .   .   .   .   .   .   .   .   AssignmentExpr
.   .   .   .   .   .   .   .   .   .   IdentExpr(n,v(3: n,0,false,true))
.   .   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   .   IdentExpr(n,v(3: n,0,false,true))
.   .   .   .   .   .   .   .   .   .   .   IdentExpr(z,v(5: z,1,false,true))
.   .   .   .   .   .   .   .   ReturnStmt
.   .   .   .   .   .   .   .   .   IdentExpr(n,v(3: n,0,false,true))
`)

	code = `
const a = 123
const b = 456

fn foo() {
	assert(b == 456)
	assert(a == 123)
}
foo()
`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:3)
.   BlockNode(Scope defs:{a: v(1: a,1,true,false), b: v(2: b,2,true,false), foo: v(0: foo,0,true,false)})
.   .   ConstStmt
.   .   .   IdentExpr(a,v(1: a,1,true,false))
.   .   .   BasicExpr(Int,"123")
.   .   ConstStmt
.   .   .   IdentExpr(b,v(2: b,2,true,false))
.   .   .   BasicExpr(Int,"456")
.   .   NamedFnStmt
.   .   .   IdentExpr(foo,v(0: foo,0,true,false))
.   .   .   FnExpr(FuncScope defs:{} captures:{a: (parent: v(1: a,1,true,false), child v(4: a,1,true,true)), b: (parent: v(2: b,2,true,false), child v(3: b,0,true,true))} numLocals:0)
.   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   ExprStmt
.   .   .   .   .   .   InvokeExpr
.   .   .   .   .   .   .   BuiltinExpr("assert")
.   .   .   .   .   .   .   BinaryExpr("==")
.   .   .   .   .   .   .   .   IdentExpr(b,v(3: b,0,true,true))
.   .   .   .   .   .   .   .   BasicExpr(Int,"456")
.   .   .   .   .   ExprStmt
.   .   .   .   .   .   InvokeExpr
.   .   .   .   .   .   .   BuiltinExpr("assert")
.   .   .   .   .   .   .   BinaryExpr("==")
.   .   .   .   .   .   .   .   IdentExpr(a,v(4: a,1,true,true))
.   .   .   .   .   .   .   .   BasicExpr(Int,"123")
.   .   ExprStmt
.   .   .   InvokeExpr
.   .   .   .   IdentExpr(foo,v(0: foo,0,true,false))
`)

	code = `
const foo = 1

fn bar() {
	let a = || => foo
	let b = || => foo
}
`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{bar: v(0: bar,0,true,false), foo: v(1: foo,1,true,false)})
.   .   ConstStmt
.   .   .   IdentExpr(foo,v(1: foo,1,true,false))
.   .   .   BasicExpr(Int,"1")
.   .   NamedFnStmt
.   .   .   IdentExpr(bar,v(0: bar,0,true,false))
.   .   .   FnExpr(FuncScope defs:{} captures:{foo: (parent: v(1: foo,1,true,false), child v(2: foo,0,true,true))} numLocals:2)
.   .   .   .   BlockNode(Scope defs:{a: v(4: a,0,false,false), b: v(6: b,1,false,false)})
.   .   .   .   .   LetStmt
.   .   .   .   .   .   IdentExpr(a,v(4: a,0,false,false))
.   .   .   .   .   .   FnExpr(FuncScope defs:{} captures:{foo: (parent: v(2: foo,0,true,true), child v(3: foo,0,true,true))} numLocals:0)
.   .   .   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   .   .   ExprStmt
.   .   .   .   .   .   .   .   .   IdentExpr(foo,v(3: foo,0,true,true))
.   .   .   .   .   LetStmt
.   .   .   .   .   .   IdentExpr(b,v(6: b,1,false,false))
.   .   .   .   .   .   FnExpr(FuncScope defs:{} captures:{foo: (parent: v(2: foo,0,true,true), child v(5: foo,0,true,true))} numLocals:0)
.   .   .   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   .   .   ExprStmt
.   .   .   .   .   .   .   .   .   IdentExpr(foo,v(5: foo,0,true,true))
`)
}

func TestImport(t *testing.T) {
	errors := NewAnalyzer(newModule("import foo; let foo = 2;")).Analyze()
	fail(t, errors, "[Symbol 'foo' is already defined, at foo.glm:1:17]")

	errors = NewAnalyzer(newModule("import foo; import foo;")).Analyze()
	fail(t, errors, "[Symbol 'foo' is already defined, at foo.glm:1:20]")

	errors = NewAnalyzer(newModule("import foo, zork; foo = 2;")).Analyze()
	fail(t, errors, "[Symbol 'foo' is constant, at foo.glm:1:19]")
}

func TestArity(t *testing.T) {

	code := `
fn(a...) {
}
`
	mod := newModule(code)
	errors := NewAnalyzer(mod).Analyze()

	//println(code)
	//println(ast.Dump(mod.InitFunc))
	//println(errors)

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:0)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   FnExpr(FuncScope defs:{a: v(0: a,0,false,false)} captures:{} numLocals:1)
.   .   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   .   BlockNode(Scope defs:{})
`)

}
