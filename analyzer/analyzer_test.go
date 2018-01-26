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

func dump(source string) {
	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner, isBuiltIn)
	mod, err := parser.ParseModule()
	if err != nil {
		panic("analyzer_test: could not parse")
	}
	fmt.Println(ast.Dump(mod))
}

func newAnalyzer(source string) Analyzer {

	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner, isBuiltIn)
	mod, err := parser.ParseModule()
	if err != nil {
		panic("analyzer_test: could not parse")
	}
	return NewAnalyzer(mod)
}

func TestFlat(t *testing.T) {

	anl := newAnalyzer("let a = 1; const b = 2; a = b + 3;;")
	errors := anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(a,(0,false,false))
.   .   .   BasicExpr(INT,"1")
.   .   Const
.   .   .   IdentExpr(b,(1,true,false))
.   .   .   BasicExpr(INT,"2")
.   .   AssignmentExpr
.   .   .   IdentExpr(a,(0,false,false))
.   .   .   BinaryExpr("+")
.   .   .   .   IdentExpr(b,(1,true,false))
.   .   .   .   BasicExpr(INT,"3")
.   .   Empty
`)

	errors = newAnalyzer("a;").Analyze()
	fail(t, errors, "[Symbol 'a' is not defined]")

	errors = newAnalyzer("let a = 1;const a = 1;").Analyze()
	fail(t, errors, "[Symbol 'a' is already defined]")

	errors = newAnalyzer("const a = 1;a = 1;").Analyze()
	fail(t, errors, "[Symbol 'a' is constant]")

	errors = newAnalyzer("a = a;").Analyze()
	fail(t, errors, "[Symbol 'a' is not defined Symbol 'a' is not defined]")
}

func TestNested(t *testing.T) {

	source := `
let a = 1;
if (true) {
    a = 2;
    const b = 2;
} else {
    a = 3;
    let b = 3;
}`
	anl := newAnalyzer(source)
	//errors := anl.Analyze()
	//fmt.Println(source)
	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)

	errors := anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:3 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(a,(0,false,false))
.   .   .   BasicExpr(INT,"1")
.   .   If
.   .   .   BasicExpr(TRUE,"true")
.   .   .   Block
.   .   .   .   AssignmentExpr
.   .   .   .   .   IdentExpr(a,(0,false,false))
.   .   .   .   .   BasicExpr(INT,"2")
.   .   .   .   Const
.   .   .   .   .   IdentExpr(b,(1,true,false))
.   .   .   .   .   BasicExpr(INT,"2")
.   .   .   Block
.   .   .   .   AssignmentExpr
.   .   .   .   .   IdentExpr(a,(0,false,false))
.   .   .   .   .   BasicExpr(INT,"3")
.   .   .   .   Let
.   .   .   .   .   IdentExpr(b,(2,false,false))
.   .   .   .   .   BasicExpr(INT,"3")
`)
}

func TestLoop(t *testing.T) {

	anl := newAnalyzer("while true { 1 + 2; }")
	errors := anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:0 numCaptures:0 parentCaptures:[])
.   Block
.   .   While
.   .   .   BasicExpr(TRUE,"true")
.   .   .   Block
.   .   .   .   BinaryExpr("+")
.   .   .   .   .   BasicExpr(INT,"1")
.   .   .   .   .   BasicExpr(INT,"2")
`)

	anl = newAnalyzer("while true { 1 + 2; break; continue; }")
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:0 numCaptures:0 parentCaptures:[])
.   Block
.   .   While
.   .   .   BasicExpr(TRUE,"true")
.   .   .   Block
.   .   .   .   BinaryExpr("+")
.   .   .   .   .   BasicExpr(INT,"1")
.   .   .   .   .   BasicExpr(INT,"2")
.   .   .   .   Break
.   .   .   .   Continue
`)

	errors = newAnalyzer("break;").Analyze()
	fail(t, errors, "['break' outside of loop]")

	errors = newAnalyzer("continue;").Analyze()
	fail(t, errors, "['continue' outside of loop]")

	anl = newAnalyzer("let a; for b in [] { break; continue; }")
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:3 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(a,(0,false,false))
.   .   For
.   .   .   IdentExpr(b,(1,false,false))
.   .   .   IdentExpr(#synthetic0,(2,false,false))
.   .   .   ListExpr
.   .   .   Block
.   .   .   .   Break
.   .   .   .   Continue
`)

	anl = newAnalyzer("for (a, b) in [] { }")
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:3 numCaptures:0 parentCaptures:[])
.   Block
.   .   For
.   .   .   IdentExpr(a,(0,false,false))
.   .   .   IdentExpr(b,(1,false,false))
.   .   .   IdentExpr(#synthetic0,(2,false,false))
.   .   .   ListExpr
.   .   .   Block
`)

	anl = newAnalyzer(`
for a in [] {
    for b in [] {
    }
}
`)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:4 numCaptures:0 parentCaptures:[])
.   Block
.   .   For
.   .   .   IdentExpr(a,(0,false,false))
.   .   .   IdentExpr(#synthetic0,(1,false,false))
.   .   .   ListExpr
.   .   .   Block
.   .   .   .   For
.   .   .   .   .   IdentExpr(b,(2,false,false))
.   .   .   .   .   IdentExpr(#synthetic1,(3,false,false))
.   .   .   .   .   ListExpr
.   .   .   .   .   Block
`)

	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)
}

func TestPureFunction(t *testing.T) {
	source := `
let a = 1;
let b = fn(x) {
    let c = fn(y, z) {
        if (y < 33) {
            return y + z + 5;
        } else {
            let b = 42;
        }
    };
    return c(3);
};`

	anl := newAnalyzer(source)
	errors := anl.Analyze()

	//fmt.Println(source)
	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)

	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(a,(0,false,false))
.   .   .   BasicExpr(INT,"1")
.   .   Let
.   .   .   IdentExpr(b,(1,false,false))
.   .   .   FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   .   .   .   IdentExpr(x,(0,false,false))
.   .   .   .   Block
.   .   .   .   .   Let
.   .   .   .   .   .   IdentExpr(c,(1,false,false))
.   .   .   .   .   .   FnExpr(numLocals:3 numCaptures:0 parentCaptures:[])
.   .   .   .   .   .   .   IdentExpr(y,(0,false,false))
.   .   .   .   .   .   .   IdentExpr(z,(1,false,false))
.   .   .   .   .   .   .   Block
.   .   .   .   .   .   .   .   If
.   .   .   .   .   .   .   .   .   BinaryExpr("<")
.   .   .   .   .   .   .   .   .   .   IdentExpr(y,(0,false,false))
.   .   .   .   .   .   .   .   .   .   BasicExpr(INT,"33")
.   .   .   .   .   .   .   .   .   Block
.   .   .   .   .   .   .   .   .   .   Return
.   .   .   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   .   .   .   IdentExpr(y,(0,false,false))
.   .   .   .   .   .   .   .   .   .   .   .   .   IdentExpr(z,(1,false,false))
.   .   .   .   .   .   .   .   .   .   .   .   BasicExpr(INT,"5")
.   .   .   .   .   .   .   .   .   Block
.   .   .   .   .   .   .   .   .   .   Let
.   .   .   .   .   .   .   .   .   .   .   IdentExpr(b,(2,false,false))
.   .   .   .   .   .   .   .   .   .   .   BasicExpr(INT,"42")
.   .   .   .   .   Return
.   .   .   .   .   .   InvokeExpr
.   .   .   .   .   .   .   IdentExpr(c,(1,false,false))
.   .   .   .   .   .   .   BasicExpr(INT,"3")
`)
}

func TestCaptureFunction(t *testing.T) {

	source := `
const accumGen = fn(n) {
    return fn(i) {
        n = n + i;
        return n;
    };
};
`
	anl := newAnalyzer(source)
	errors := anl.Analyze()

	//fmt.Println(source)
	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)

	ok(t, anl, errors, `
FnExpr(numLocals:1 numCaptures:0 parentCaptures:[])
.   Block
.   .   Const
.   .   .   IdentExpr(accumGen,(0,true,false))
.   .   .   FnExpr(numLocals:1 numCaptures:0 parentCaptures:[])
.   .   .   .   IdentExpr(n,(0,false,false))
.   .   .   .   Block
.   .   .   .   .   Return
.   .   .   .   .   .   FnExpr(numLocals:1 numCaptures:1 parentCaptures:[(0,false,false)])
.   .   .   .   .   .   .   IdentExpr(i,(0,false,false))
.   .   .   .   .   .   .   Block
.   .   .   .   .   .   .   .   AssignmentExpr
.   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
.   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
.   .   .   .   .   .   .   .   .   .   IdentExpr(i,(0,false,false))
.   .   .   .   .   .   .   .   Return
.   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
`)

	source = `
let z = 2;
const accumGen = fn(n) {
    return fn(i) {
        n = n + i;
        n = n + z;
        return n;
    };
};
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()

	//fmt.Println(source)
	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)

	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(z,(0,false,false))
.   .   .   BasicExpr(INT,"2")
.   .   Const
.   .   .   IdentExpr(accumGen,(1,true,false))
.   .   .   FnExpr(numLocals:1 numCaptures:1 parentCaptures:[(0,false,false)])
.   .   .   .   IdentExpr(n,(0,false,false))
.   .   .   .   Block
.   .   .   .   .   Return
.   .   .   .   .   .   FnExpr(numLocals:1 numCaptures:2 parentCaptures:[(0,false,false), (0,false,true)])
.   .   .   .   .   .   .   IdentExpr(i,(0,false,false))
.   .   .   .   .   .   .   Block
.   .   .   .   .   .   .   .   AssignmentExpr
.   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
.   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
.   .   .   .   .   .   .   .   .   .   IdentExpr(i,(0,false,false))
.   .   .   .   .   .   .   .   AssignmentExpr
.   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
.   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
.   .   .   .   .   .   .   .   .   .   IdentExpr(z,(1,false,true))
.   .   .   .   .   .   .   .   Return
.   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
`)

	source = `
const a = 123;
const b = 456;

fn foo() {
    assert(b == 456);
    assert(a == 123);
}
foo();
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()

	ok(t, anl, errors, `
FnExpr(numLocals:3 numCaptures:0 parentCaptures:[])
.   Block
.   .   Const
.   .   .   IdentExpr(a,(1,true,false))
.   .   .   BasicExpr(INT,"123")
.   .   Const
.   .   .   IdentExpr(b,(2,true,false))
.   .   .   BasicExpr(INT,"456")
.   .   NamedFn
.   .   .   IdentExpr(foo,(0,true,false))
.   .   .   FnExpr(numLocals:0 numCaptures:2 parentCaptures:[(2,true,false), (1,true,false)])
.   .   .   .   Block
.   .   .   .   .   InvokeExpr
.   .   .   .   .   .   BuiltinExpr("assert")
.   .   .   .   .   .   BinaryExpr("==")
.   .   .   .   .   .   .   IdentExpr(b,(0,true,true))
.   .   .   .   .   .   .   BasicExpr(INT,"456")
.   .   .   .   .   InvokeExpr
.   .   .   .   .   .   BuiltinExpr("assert")
.   .   .   .   .   .   BinaryExpr("==")
.   .   .   .   .   .   .   IdentExpr(a,(1,true,true))
.   .   .   .   .   .   .   BasicExpr(INT,"123")
.   .   InvokeExpr
.   .   .   IdentExpr(foo,(0,true,false))
`)

	//println(source)
	//println(ast.Dump(anl.Module()))
	//println(errors)
}

func TestStruct(t *testing.T) {

	errors := newAnalyzer("this;").Analyze()
	fail(t, errors, "['this' outside of loop]")

	source := `
struct{ };
`
	anl := newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:0 numCaptures:0 parentCaptures:[])
.   Block
.   .   StructExpr([],-1)
`)

	source = `
struct{ a: 1 };
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:0 numCaptures:0 parentCaptures:[])
.   Block
.   .   StructExpr([a],-1)
.   .   .   BasicExpr(INT,"1")
`)

	source = `
struct{ a: this };
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:1 numCaptures:0 parentCaptures:[])
.   Block
.   .   StructExpr([a],0)
.   .   .   ThisExpr((0,true,false))
`)

	source = `
struct{ a: struct { b: this } };
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:1 numCaptures:0 parentCaptures:[])
.   Block
.   .   StructExpr([a],-1)
.   .   .   StructExpr([b],0)
.   .   .   .   ThisExpr((0,true,false))
`)

	source = `
struct{ a: struct { b: 1 }, c: this.a };
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:1 numCaptures:0 parentCaptures:[])
.   Block
.   .   StructExpr([a, c],0)
.   .   .   StructExpr([b],-1)
.   .   .   .   BasicExpr(INT,"1")
.   .   .   FieldExpr(a)
.   .   .   .   ThisExpr((0,true,false))
`)

	source = `
struct{ a: struct { b: this }, c: this };
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   StructExpr([a, c],1)
.   .   .   StructExpr([b],0)
.   .   .   .   ThisExpr((0,true,false))
.   .   .   ThisExpr((1,true,false))
`)

	source = `
struct{ a: this, b: struct { c: this } };
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   StructExpr([a, b],0)
.   .   .   ThisExpr((0,true,false))
.   .   .   StructExpr([c],1)
.   .   .   .   ThisExpr((1,true,false))
`)

	source = `
let a = struct {
    x: 8,
    y: 5,
    plus:  fn() { return this.x + this.y; },
    minus: fn() { return this.x - this.y; }
};
let b = a.plus();
let c = a.minus();
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()

	ok(t, anl, errors, `
FnExpr(numLocals:4 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(a,(1,false,false))
.   .   .   StructExpr([x, y, plus, minus],0)
.   .   .   .   BasicExpr(INT,"8")
.   .   .   .   BasicExpr(INT,"5")
.   .   .   .   FnExpr(numLocals:0 numCaptures:1 parentCaptures:[(0,true,false)])
.   .   .   .   .   Block
.   .   .   .   .   .   Return
.   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   FieldExpr(x)
.   .   .   .   .   .   .   .   .   ThisExpr((0,true,true))
.   .   .   .   .   .   .   .   FieldExpr(y)
.   .   .   .   .   .   .   .   .   ThisExpr((0,true,true))
.   .   .   .   FnExpr(numLocals:0 numCaptures:1 parentCaptures:[(0,true,false)])
.   .   .   .   .   Block
.   .   .   .   .   .   Return
.   .   .   .   .   .   .   BinaryExpr("-")
.   .   .   .   .   .   .   .   FieldExpr(x)
.   .   .   .   .   .   .   .   .   ThisExpr((0,true,true))
.   .   .   .   .   .   .   .   FieldExpr(y)
.   .   .   .   .   .   .   .   .   ThisExpr((0,true,true))
.   .   Let
.   .   .   IdentExpr(b,(2,false,false))
.   .   .   InvokeExpr
.   .   .   .   FieldExpr(plus)
.   .   .   .   .   IdentExpr(a,(1,false,false))
.   .   Let
.   .   .   IdentExpr(c,(3,false,false))
.   .   .   InvokeExpr
.   .   .   .   FieldExpr(minus)
.   .   .   .   .   IdentExpr(a,(1,false,false))
`)
}

func TestAssignment(t *testing.T) {

	source := `
let x = struct { a: 0 };
let y = x.a;
x.a = 3;
x.a++;
y--;
x[y] = 42;
y = x[3];
x[2]++;
y.z = x[2]++;
let g, h = 5;
const i = 6, j;
`
	anl := newAnalyzer(source)
	errors := anl.Analyze()

	//fmt.Println(source)
	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)

	ok(t, anl, errors, `
FnExpr(numLocals:6 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(x,(0,false,false))
.   .   .   StructExpr([a],-1)
.   .   .   .   BasicExpr(INT,"0")
.   .   Let
.   .   .   IdentExpr(y,(1,false,false))
.   .   .   FieldExpr(a)
.   .   .   .   IdentExpr(x,(0,false,false))
.   .   AssignmentExpr
.   .   .   FieldExpr(a)
.   .   .   .   IdentExpr(x,(0,false,false))
.   .   .   BasicExpr(INT,"3")
.   .   PostfixExpr("++")
.   .   .   FieldExpr(a)
.   .   .   .   IdentExpr(x,(0,false,false))
.   .   PostfixExpr("--")
.   .   .   IdentExpr(y,(1,false,false))
.   .   AssignmentExpr
.   .   .   IndexExpr
.   .   .   .   IdentExpr(x,(0,false,false))
.   .   .   .   IdentExpr(y,(1,false,false))
.   .   .   BasicExpr(INT,"42")
.   .   AssignmentExpr
.   .   .   IdentExpr(y,(1,false,false))
.   .   .   IndexExpr
.   .   .   .   IdentExpr(x,(0,false,false))
.   .   .   .   BasicExpr(INT,"3")
.   .   PostfixExpr("++")
.   .   .   IndexExpr
.   .   .   .   IdentExpr(x,(0,false,false))
.   .   .   .   BasicExpr(INT,"2")
.   .   AssignmentExpr
.   .   .   FieldExpr(z)
.   .   .   .   IdentExpr(y,(1,false,false))
.   .   .   PostfixExpr("++")
.   .   .   .   IndexExpr
.   .   .   .   .   IdentExpr(x,(0,false,false))
.   .   .   .   .   BasicExpr(INT,"2")
.   .   Let
.   .   .   IdentExpr(g,(2,false,false))
.   .   .   IdentExpr(h,(3,false,false))
.   .   .   BasicExpr(INT,"5")
.   .   Const
.   .   .   IdentExpr(i,(4,true,false))
.   .   .   BasicExpr(INT,"6")
.   .   .   IdentExpr(j,(5,true,false))
`)
}

func TestList(t *testing.T) {

	source := `
let a = ['x'][0];
let b = ['x'];
b[0] = 3;
b[0]++;
`
	anl := newAnalyzer(source)
	errors := anl.Analyze()

	//fmt.Println(source)
	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)

	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(a,(0,false,false))
.   .   .   IndexExpr
.   .   .   .   ListExpr
.   .   .   .   .   BasicExpr(STR,"x")
.   .   .   .   BasicExpr(INT,"0")
.   .   Let
.   .   .   IdentExpr(b,(1,false,false))
.   .   .   ListExpr
.   .   .   .   BasicExpr(STR,"x")
.   .   AssignmentExpr
.   .   .   IndexExpr
.   .   .   .   IdentExpr(b,(1,false,false))
.   .   .   .   BasicExpr(INT,"0")
.   .   .   BasicExpr(INT,"3")
.   .   PostfixExpr("++")
.   .   .   IndexExpr
.   .   .   .   IdentExpr(b,(1,false,false))
.   .   .   .   BasicExpr(INT,"0")
`)
}

func TestTry(t *testing.T) {

	source := "let a = 1; try { } catch e { } finally { }"
	anl := newAnalyzer(source)
	errors := anl.Analyze()

	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(a,(0,false,false))
.   .   .   BasicExpr(INT,"1")
.   .   Try
.   .   .   Block
.   .   .   IdentExpr(e,(1,true,false))
.   .   .   Block
.   .   .   Block
`)

	source = "let a = 1; try { } catch a { } finally { }"
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	fail(t, errors, "[Symbol 'a' is already defined]")
}

func TestNamedFunc(t *testing.T) {

	source := `
fn a() {
    return b();
}
fn b() {
    return 42;
}
`
	anl := newAnalyzer(source)
	errors := anl.Analyze()

	//fmt.Println(source)
	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)

	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   NamedFn
.   .   .   IdentExpr(a,(0,true,false))
.   .   .   FnExpr(numLocals:0 numCaptures:1 parentCaptures:[(1,true,false)])
.   .   .   .   Block
.   .   .   .   .   Return
.   .   .   .   .   .   InvokeExpr
.   .   .   .   .   .   .   IdentExpr(b,(0,true,true))
.   .   NamedFn
.   .   .   IdentExpr(b,(1,true,false))
.   .   .   FnExpr(numLocals:0 numCaptures:0 parentCaptures:[])
.   .   .   .   Block
.   .   .   .   .   Return
.   .   .   .   .   .   BasicExpr(INT,"42")
`)

	errors = newAnalyzer("fn a() {} const a = 1;").Analyze()
	fail(t, errors, "[Symbol 'a' is already defined]")
}

func TestImport(t *testing.T) {

	anl := newAnalyzer("import sys; let b = 2;")
	errors := anl.Analyze()

	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)

	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   Import
.   .   .   IdentExpr(sys,(0,true,false))
.   .   Let
.   .   .   IdentExpr(b,(1,false,false))
.   .   .   BasicExpr(INT,"2")
`)

	errors = newAnalyzer("import sys; let sys = 2;").Analyze()
	fail(t, errors, "[Symbol 'sys' is already defined]")

	errors = newAnalyzer("import sys; sys = 2;").Analyze()
	fail(t, errors, "[Symbol 'sys' is constant]")

	errors = newAnalyzer("import foo;").Analyze()
	fail(t, errors, "[Module 'foo' is not defined]")
}

func TestFormalParams(t *testing.T) {

	errors := newAnalyzer("fn(const a, b) { a = 1; };").Analyze()
	fail(t, errors, "[Symbol 'a' is constant]")
}
