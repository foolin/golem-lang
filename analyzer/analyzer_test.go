// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package analyzer

import (
	"fmt"
	"github.com/mjarmy/golem-lang/ast"
	"github.com/mjarmy/golem-lang/parser"
	"github.com/mjarmy/golem-lang/scanner"
	"testing"
)

func ok(t *testing.T, anl Analyzer, errors []error, dump string) {

	if len(errors) != 0 {
		t.Error(errors)
	}

	if "\n"+ast.Dump(anl.Module()) != dump {
		t.Error("\n"+ast.Dump(anl.Module()), " != ", dump)
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

func newAnalyzer(source string) Analyzer {

	ast.InternalResetDebugging()

	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner, isBuiltIn)
	mod, err := parser.ParseModule()
	if err != nil {
		panic("analyzer_test: could not parse")
	}
	return NewAnalyzer(mod)
}

func TestFlat(t *testing.T) {

	anl := newAnalyzer("let a = 1; const b = 2; a = b + 3;")
	errors := anl.Analyze()
	ok(t, anl, errors, `
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

	errors = newAnalyzer("a;").Analyze()
	fail(t, errors, "[Symbol 'a' is not defined, at (1, 1)]")

	errors = newAnalyzer("let a = 1;const a = 1;").Analyze()
	fail(t, errors, "[Symbol 'a' is already defined, at (1, 17)]")

	errors = newAnalyzer("const a = 1;a = 1;").Analyze()
	fail(t, errors, "[Symbol 'a' is constant, at (1, 13)]")

	errors = newAnalyzer("a = a;").Analyze()
	fail(t, errors, "[Symbol 'a' is not defined, at (1, 5) Symbol 'a' is not defined, at (1, 1)]")
}

func TestNested(t *testing.T) {

	source := `
let a = 1
if (true) {
    a = 2
    const b = 2
} else {
    a = 3
    let b = 3
}`
	anl := newAnalyzer(source)

	errors := anl.Analyze()
	ok(t, anl, errors, `
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

	anl := newAnalyzer("while true { 1 + 2; }")
	errors := anl.Analyze()
	ok(t, anl, errors, `
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

	anl = newAnalyzer("while true { 1 + 2; break; continue; }")
	errors = anl.Analyze()
	ok(t, anl, errors, `
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

	errors = newAnalyzer("break;").Analyze()
	fail(t, errors, "['break' outside of loop, at (1, 1)]")

	errors = newAnalyzer("continue;").Analyze()
	fail(t, errors, "['continue' outside of loop, at (1, 1)]")

	anl = newAnalyzer("let a; for b in [] { break; continue; }")
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:3)
.   BlockNode(Scope defs:{a: v(0: a,0,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   ForStmt(Scope defs:{#synthetic0: v(2: #synthetic0,2,false,false), b: v(1: b,1,false,false)})
.   .   .   IdentExpr(b,v(1: b,1,false,false))
.   .   .   IdentExpr(#synthetic0,v(2: #synthetic0,2,false,false))
.   .   .   ListExpr
.   .   .   BlockNode(Scope defs:{})
.   .   .   .   BreakStmt
.   .   .   .   ContinueStmt
`)

	anl = newAnalyzer("for (a, b) in [] { }")
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:3)
.   BlockNode(Scope defs:{})
.   .   ForStmt(Scope defs:{#synthetic0: v(2: #synthetic0,2,false,false), a: v(0: a,0,false,false), b: v(1: b,1,false,false)})
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   IdentExpr(b,v(1: b,1,false,false))
.   .   .   IdentExpr(#synthetic0,v(2: #synthetic0,2,false,false))
.   .   .   ListExpr
.   .   .   BlockNode(Scope defs:{})
`)

	anl = newAnalyzer(`
for a in [] {
	for b in [] {
	}
}
`)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:4)
.   BlockNode(Scope defs:{})
.   .   ForStmt(Scope defs:{#synthetic0: v(1: #synthetic0,1,false,false), a: v(0: a,0,false,false)})
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   IdentExpr(#synthetic0,v(1: #synthetic0,1,false,false))
.   .   .   ListExpr
.   .   .   BlockNode(Scope defs:{})
.   .   .   .   ForStmt(Scope defs:{#synthetic1: v(3: #synthetic1,3,false,false), b: v(2: b,2,false,false)})
.   .   .   .   .   IdentExpr(b,v(2: b,2,false,false))
.   .   .   .   .   IdentExpr(#synthetic1,v(3: #synthetic1,3,false,false))
.   .   .   .   .   ListExpr
.   .   .   .   .   BlockNode(Scope defs:{})
`)

	anl = newAnalyzer(`
let a = 1
for a in [] {
}
`)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:3)
.   BlockNode(Scope defs:{a: v(0: a,0,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(a,v(0: a,0,false,false))
.   .   .   BasicExpr(Int,"1")
.   .   ForStmt(Scope defs:{#synthetic0: v(2: #synthetic0,2,false,false), a: v(1: a,1,false,false)})
.   .   .   IdentExpr(a,v(1: a,1,false,false))
.   .   .   IdentExpr(#synthetic0,v(2: #synthetic0,2,false,false))
.   .   .   ListExpr
.   .   .   BlockNode(Scope defs:{})
`)
}

func TestStruct(t *testing.T) {

	errors := newAnalyzer("this").Analyze()
	fail(t, errors, "['this' outside of struct, at (1, 1)]")

	source := `
struct{ }
`
	anl := newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:0)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([],StructScope defs:{})
`)

	source = `
struct{ a: 1 }
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:0)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([a],StructScope defs:{})
.   .   .   .   BasicExpr(Int,"1")
`)

	source = `
struct{ a: this }
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:1)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([a],StructScope defs:{this: v(0: this,0,true,false)})
.   .   .   .   ThisExpr(v(0: this,0,true,false))
`)

	source = `
struct{ a: struct { b: this } }
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:1)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([a],StructScope defs:{})
.   .   .   .   StructExpr([b],StructScope defs:{this: v(0: this,0,true,false)})
.   .   .   .   .   ThisExpr(v(0: this,0,true,false))
`)

	source = `
struct{ a: struct { b: 1 }, c: this.a }
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:1)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([a, c],StructScope defs:{this: v(0: this,0,true,false)})
.   .   .   .   StructExpr([b],StructScope defs:{})
.   .   .   .   .   BasicExpr(Int,"1")
.   .   .   .   FieldExpr(a)
.   .   .   .   .   ThisExpr(v(0: this,0,true,false))
`)

	source = `
struct{ a: struct { b: this }, c: this }
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([a, c],StructScope defs:{this: v(1: this,1,true,false)})
.   .   .   .   StructExpr([b],StructScope defs:{this: v(0: this,0,true,false)})
.   .   .   .   .   ThisExpr(v(0: this,0,true,false))
.   .   .   .   ThisExpr(v(1: this,1,true,false))
`)

	source = `
struct{ a: this, b: struct { c: this } }
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([a, b],StructScope defs:{this: v(0: this,0,true,false)})
.   .   .   .   ThisExpr(v(0: this,0,true,false))
.   .   .   .   StructExpr([c],StructScope defs:{this: v(1: this,1,true,false)})
.   .   .   .   .   ThisExpr(v(1: this,1,true,false))
`)

	source = `
let a = struct {
	x: 8,
	y: 5,
	plus:  fn() { return this.x + this.y; },
	minus: fn() { return this.x - this.y; }
}
let b = a.plus()
let c = a.minus()
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()

	ok(t, anl, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:3)
.   BlockNode(Scope defs:{a: v(3: a,0,false,false), b: v(4: b,1,false,false), c: v(5: c,2,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(a,v(3: a,0,false,false))
.   .   .   StructExpr([x, y, plus, minus],StructScope defs:{this: v(0: this,0,true,false)})
.   .   .   .   BasicExpr(Int,"8")
.   .   .   .   BasicExpr(Int,"5")
.   .   .   .   FnExpr(FuncScope defs:{} captures:{this: (parent: v(0: this,0,true,false), child v(1: this,0,true,true))} numLocals:1)
.   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   ReturnStmt
.   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   FieldExpr(x)
.   .   .   .   .   .   .   .   .   ThisExpr(v(1: this,0,true,true))
.   .   .   .   .   .   .   .   FieldExpr(y)
.   .   .   .   .   .   .   .   .   ThisExpr(v(1: this,0,true,true))
.   .   .   .   FnExpr(FuncScope defs:{} captures:{this: (parent: v(0: this,0,true,false), child v(2: this,0,true,true))} numLocals:0)
.   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   ReturnStmt
.   .   .   .   .   .   .   BinaryExpr("-")
.   .   .   .   .   .   .   .   FieldExpr(x)
.   .   .   .   .   .   .   .   .   ThisExpr(v(2: this,0,true,true))
.   .   .   .   .   .   .   .   FieldExpr(y)
.   .   .   .   .   .   .   .   .   ThisExpr(v(2: this,0,true,true))
.   .   LetStmt
.   .   .   IdentExpr(b,v(4: b,1,false,false))
.   .   .   InvokeExpr
.   .   .   .   FieldExpr(plus)
.   .   .   .   .   IdentExpr(a,v(3: a,0,false,false))
.   .   LetStmt
.   .   .   IdentExpr(c,v(5: c,2,false,false))
.   .   .   InvokeExpr
.   .   .   .   FieldExpr(minus)
.   .   .   .   .   IdentExpr(a,v(3: a,0,false,false))
`)
}

func TestAssignment(t *testing.T) {

	source := `
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
	anl := newAnalyzer(source)
	errors := anl.Analyze()

	ok(t, anl, errors, `
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

	source := `
let a = ['x'][0]
let b = ['x']
b[0] = 3
b[0]++
`
	anl := newAnalyzer(source)
	errors := anl.Analyze()

	ok(t, anl, errors, `
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

	source := "let a = 1; try { } catch e { } finally { }"
	anl := newAnalyzer(source)
	errors := anl.Analyze()

	ok(t, anl, errors, `
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

	source = "let a = 1; try { } catch a { } finally { }"
	anl = newAnalyzer(source)
	errors = anl.Analyze()

	ok(t, anl, errors, `
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

	errors := newAnalyzer("fn(const a, b) { a = 1; };").Analyze()
	fail(t, errors, "[Symbol 'a' is constant, at (1, 18)]")
}

func TestImport(t *testing.T) {
	errors := newAnalyzer("import foo; let foo = 2;").Analyze()
	fail(t, errors, "[Symbol 'foo' is already defined, at (1, 17)]")

	errors = newAnalyzer("import foo, zork; foo = 2;").Analyze()
	fail(t, errors, "[Symbol 'foo' is constant, at (1, 19)]")
}

func TestPureFunction(t *testing.T) {
	source := `
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

	anl := newAnalyzer(source)
	errors := anl.Analyze()

	ok(t, anl, errors, `
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

	source := `
fn a() {
    return b()
}
fn b() {
    return 42
}
`
	anl := newAnalyzer(source)
	errors := anl.Analyze()

	ok(t, anl, errors, `
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

	errors = newAnalyzer("fn a() {}; const a = 1;").Analyze()
	fail(t, errors, "[Symbol 'a' is already defined, at (1, 18)]")
}

func TestCaptureFunction(t *testing.T) {

	source := `
const accumGen = fn(n) {
    return fn(i) {
        n = n + i
        return n
    }
}
`
	anl := newAnalyzer(source)
	errors := anl.Analyze()

	ok(t, anl, errors, `
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

	source = `
let z = 2
const accumGen = fn(n) {
	return fn(i) {
		n = n + i
		n = n + z
		return n
	}
}
	`
	anl = newAnalyzer(source)
	errors = anl.Analyze()

	ok(t, anl, errors, `
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

	source = `
const a = 123
const b = 456

fn foo() {
	assert(b == 456)
	assert(a == 123)
}
foo()
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()

	ok(t, anl, errors, `
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

	source = `
const foo = 1

fn bar() {
	let a = || => foo
	let b = || => foo
}
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()

	//println(source)
	//println(ast.Dump(anl.Module()))
	//println(errors)

	ok(t, anl, errors, `
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
