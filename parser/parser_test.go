// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package parser

import (
	"runtime"
	"testing"

	"github.com/mjarmy/golem-lang/ast"
	"github.com/mjarmy/golem-lang/scanner"
)

func ok(t *testing.T, p *Parser, expect string) {

	mod, err := p.ParseModule()
	if err != nil {
		t.Error(err, " != nil")
		panic("ok")
	} else if mod.String() != expect {
		t.Error(mod, " != ", expect)
		panic("ok")
	}
}

func fail(t *testing.T, p *Parser, expect string) {

	mod, err := p.ParseModule()
	if mod != nil {
		t.Error(mod, " != nil")
	}

	if err.Error() != expect {
		t.Error(err, " != ", expect)
		panic("fail")
	}
}

func okExpr(t *testing.T, p *Parser, expect string) {

	expr, err := parseExpression(p)
	if err != nil {
		t.Error(err, " != nil")
	}

	if expr.String() != expect {
		t.Error(expr, " != ", expect)
	}
}

func failExpr(t *testing.T, p *Parser, expect string) {

	expr, err := parseExpression(p)
	if expr != nil {
		t.Error(expr, " != nil")
	}

	if err.Error() != expect {
		t.Error(err, " != ", expect)
	}
}

func newParser(source string) *Parser {
	builtins := map[string]bool{
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
	isBuiltIn := func(s string) bool {
		_, ok := builtins[s]
		return ok
	}

	return NewParser(scanner.NewScanner(source), isBuiltIn)
}

func parseExpression(p *Parser) (expr ast.Expression, err error) {

	// In a recursive descent parser, errors can be generated deep
	// in the call stack.  We are going to use panic-recover to handle them.
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			expr = nil
			err = r.(error)
		}
	}()

	// read the first two tokens
	p.cur = p.advance()
	p.next = p.advance()

	// parse the expression
	expr = p.expression()
	p.expect(ast.EOF)

	return expr, err
}

func TestPrimary(t *testing.T) {

	p := newParser("")
	failExpr(t, p, "Unexpected EOF at (1, 1)")

	p = newParser("#")
	failExpr(t, p, "Unexpected Character '#' at (1, 1)")

	p = newParser("'")
	failExpr(t, p, "Unexpected EOF at (1, 2)")

	p = newParser("1 2")
	failExpr(t, p, "Unexpected Token '2' at (1, 3)")

	p = newParser("a == goto")
	failExpr(t, p, "Unexpected Reserved Word 'goto' at (1, 6)")

	p = newParser("1 #")
	failExpr(t, p, "Unexpected Character '#' at (1, 3)")

	p = newParser("1")
	okExpr(t, p, "1")

	p = newParser("0xa")
	okExpr(t, p, "0xa")

	p = newParser("1.2")
	okExpr(t, p, "1.2")

	p = newParser("null")
	okExpr(t, p, "null")

	p = newParser("true")
	okExpr(t, p, "true")

	p = newParser("false")
	okExpr(t, p, "false")

	p = newParser("'a'")
	okExpr(t, p, "'a'")

	p = newParser("('a')")
	okExpr(t, p, "'a'")

	p = newParser("bar")
	okExpr(t, p, "bar")

	p = newParser("str")
	okExpr(t, p, "str")
}

func TestUnary(t *testing.T) {
	p := newParser("-1")
	okExpr(t, p, "-1")

	p = newParser("- - 2")
	okExpr(t, p, "--2")

	p = newParser("!a")
	okExpr(t, p, "!a")

	p = newParser("~a")
	okExpr(t, p, "~a")
}

func TestPostfix(t *testing.T) {
	p := newParser("a++")
	okExpr(t, p, "a++")

	p = newParser("a.b++")
	okExpr(t, p, "a.b++")

	p = newParser("3++")
	failExpr(t, p, "Invalid Postfix Expression at (1, 2)")
}

func TestTernary(t *testing.T) {
	p := newParser("a ? b : c")
	okExpr(t, p, "(a ? b : c)")

	p = newParser("a || b ? b = c : d ? e : f")
	okExpr(t, p, "((a || b) ? (b = c) : (d ? e : f))")

	p = newParser("a ?")
	failExpr(t, p, "Unexpected EOF at (1, 4)")

	p = newParser("a ? b")
	failExpr(t, p, "Unexpected EOF at (1, 6)")

	p = newParser("a ? b :")
	failExpr(t, p, "Unexpected EOF at (1, 8)")
}

func TestMultiplicative(t *testing.T) {
	p := newParser("1*2")
	okExpr(t, p, "(1 * 2)")

	p = newParser("-1*-2")
	okExpr(t, p, "(-1 * -2)")

	p = newParser("1*2*3")
	okExpr(t, p, "((1 * 2) * 3)")

	p = newParser("1*2/3*4/5")
	okExpr(t, p, "((((1 * 2) / 3) * 4) / 5)")

	p = newParser("1%2&3<<4>>5")
	okExpr(t, p, "((((1 % 2) & 3) << 4) >> 5)")
}

func TestAdditive(t *testing.T) {
	p := newParser("1*2+3")
	okExpr(t, p, "((1 * 2) + 3)")

	p = newParser("1+2*3")
	okExpr(t, p, "(1 + (2 * 3))")

	p = newParser("1+2*-3")
	okExpr(t, p, "(1 + (2 * -3))")

	p = newParser("1+2+-3")
	okExpr(t, p, "((1 + 2) + -3)")

	p = newParser("1+2*3+4")
	okExpr(t, p, "((1 + (2 * 3)) + 4)")

	p = newParser("(1+2) * 3")
	okExpr(t, p, "((1 + 2) * 3)")

	p = newParser("(1*2) * 3")
	okExpr(t, p, "((1 * 2) * 3)")

	p = newParser("1 * (2 + 3)")
	okExpr(t, p, "(1 * (2 + 3))")

	p = newParser("1 ^ 2 | 3")
	okExpr(t, p, "((1 ^ 2) | 3)")

	p = newParser("1 ^ 2 % 3")
	okExpr(t, p, "(1 ^ (2 % 3))")

	p = newParser("1 +")
	failExpr(t, p, "Unexpected EOF at (1, 4)")
}

func TestAssign(t *testing.T) {

	p := newParser("a += 2")
	okExpr(t, p, "(a = (a + 2))")

	p = newParser("a -= 2")
	okExpr(t, p, "(a = (a - 2))")

	p = newParser("a *= 2")
	okExpr(t, p, "(a = (a * 2))")

	p = newParser("a /= 2")
	okExpr(t, p, "(a = (a / 2))")

	p = newParser("a %= 2")
	okExpr(t, p, "(a = (a % 2))")

	p = newParser("a |= 2")
	okExpr(t, p, "(a = (a | 2))")

	p = newParser("a &= 2")
	okExpr(t, p, "(a = (a & 2))")

	p = newParser("a ^= 2")
	okExpr(t, p, "(a = (a ^ 2))")

	p = newParser("a <<= 2")
	okExpr(t, p, "(a = (a << 2))")

	p = newParser("a >>= 2")
	okExpr(t, p, "(a = (a >> 2))")

	p = newParser("a = b = c")
	okExpr(t, p, "(a = (b = c))")

	p = newParser("a -= b += c")
	okExpr(t, p, "(a = (a - (b = (b + c))))")
}

func TestComparitive(t *testing.T) {
	p := newParser("1==3")
	okExpr(t, p, "(1 == 3)")

	p = newParser("1 ==2 +3 * - 4")
	okExpr(t, p, "(1 == (2 + (3 * -4)))")

	p = newParser("(1== 2)+ 3")
	okExpr(t, p, "((1 == 2) + 3)")

	p = newParser("1!=3")
	okExpr(t, p, "(1 != 3)")

	okExpr(t, newParser("1 < 3"), "(1 < 3)")
	okExpr(t, newParser("1 > 3"), "(1 > 3)")
	okExpr(t, newParser("1 <= 3"), "(1 <= 3)")
	okExpr(t, newParser("1 >= 3"), "(1 >= 3)")
	okExpr(t, newParser("1 <=> 3"), "(1 <=> 3)")

	okExpr(t, newParser("1 <=> 2 + 3 * 4"), "(1 <=> (2 + (3 * 4)))")

	okExpr(t, newParser("1 has 3"), "(1 has 3)")
}

func TestAndOr(t *testing.T) {

	okExpr(t, newParser("1 || 2"), "(1 || 2)")
	okExpr(t, newParser("1 || 2 || 3"), "((1 || 2) || 3)")

	okExpr(t, newParser("1 || 2 && 3"), "(1 || (2 && 3))")
	okExpr(t, newParser("1 || 2 && 3 < 4"), "(1 || (2 && (3 < 4)))")
}

func TestModule(t *testing.T) {
	p := newParser("let a =1==3; 2+ true; z =27;const a = 3;")
	ok(t, p, "fn() { let a = (1 == 3); (2 + true); (z = 27); const a = 3; }")
}

func TestStatement(t *testing.T) {
	p := newParser("if a { b;let c=12; };")
	ok(t, p, "fn() { if a { b; let c = 12; }; }")

	p = newParser("if a { b; } else { c; };")
	ok(t, p, "fn() { if a { b; } else { c; }; }")

	p = newParser("if a { b; } else { if(12 == 3) { z+5; };};")
	ok(t, p, "fn() { if a { b; } else { if (12 == 3) { (z + 5); }; }; }")

	p = newParser("if a {} else if b {} else {};;")
	ok(t, p, "fn() { if a {  } else if b {  } else {  };; }")

	p = newParser("while a { b; };")
	ok(t, p, "fn() { while a { b; }; }")

	p = newParser("break; continue; while a { b; continue; break; };")
	ok(t, p, "fn() { break; continue; while a { b; continue; break; }; }")

	p = newParser("a = b;")
	ok(t, p, "fn() { (a = b); }")

	p = newParser("let a = 3; const b = 4;")
	ok(t, p, "fn() { let a = 3; const b = 4; }")

	p = newParser("let a = 3, b; const x, y, z = 5; ")
	ok(t, p, "fn() { let a = 3, b; const x, y, z = 5; }")

	p = newParser("fn() {}")
	okExpr(t, p, "fn() {  }")

	p = newParser("fn() {};")
	ok(t, p, "fn() { fn() {  }; }")
}

func TestFor(t *testing.T) {

	p := newParser("for a in b {};")
	ok(t, p, "fn() { for a in b {  }; }")

	p = newParser("for (a,b) in c {};")
	ok(t, p, "fn() { for (a, b) in c {  }; }")

	p = newParser("for (a,b,c) in d {};")
	ok(t, p, "fn() { for (a, b, c) in d {  }; }")

	p = newParser("for a b")
	fail(t, p, "Unexpected Token 'b' at (1, 7)")

	p = newParser("for in")
	fail(t, p, "Unexpected Token 'in' at (1, 5)")

	p = newParser("for (a) in c {}")
	fail(t, p, "Invalid ForStmt Expression at (1, 5)")

	p = newParser("for () in c {}")
	fail(t, p, "Invalid ForStmt Expression at (1, 5)")
}

func TestFn(t *testing.T) {
	p := newParser("fn() { }")
	okExpr(t, p, "fn() {  }")

	p = newParser("fn() { a = 3; }")
	okExpr(t, p, "fn() { (a = 3); }")

	p = newParser("fn(x) { a = 3; }")
	okExpr(t, p, "fn(x) { (a = 3); }")

	p = newParser("fn(x,y) { a = 3; }")
	okExpr(t, p, "fn(x, y) { (a = 3); }")

	p = newParser("fn(x,y,z) { a = 3; }")
	okExpr(t, p, "fn(x, y, z) { (a = 3); }")

	p = newParser("fn(x) { let a = fn(y) { return x + y; }; }")
	okExpr(t, p, "fn(x) { let a = fn(y) { return (x + y); }; }")

	p = newParser("z = fn(x) { a = 2; return b; c = 3; };")
	ok(t, p, "fn() { (z = fn(x) { (a = 2); return b; (c = 3); }); }")

	p = newParser("fn a(x) {return x*x; }; fn b() { };")
	ok(t, p, "fn() { fn a(x) { return (x * x); }; fn b() {  }; }")

	p = newParser("fn(const x, y) {  }")
	okExpr(t, p, "fn(const x, y) {  }")

	p = newParser("fn(x, const y) {  }")
	okExpr(t, p, "fn(x, const y) {  }")

	p = newParser("fn(const x, const y) {  }")
	okExpr(t, p, "fn(const x, const y) {  }")

	p = newParser("fn(const a,b) { a=1; };")
	ok(t, p, "fn() { fn(const a, b) { (a = 1); }; }")

	p = newParser("return;")
	fail(t, p, "Unexpected Token ';' at (1, 7)")

}

func TestTry(t *testing.T) {

	p := newParser("throw a;")
	ok(t, p, "fn() { throw a; }")

	p = newParser("throw;")
	fail(t, p, "Unexpected Token ';' at (1, 6)")

	p = newParser("try { a; } catch e { b; };")
	ok(t, p, "fn() { try { a; } catch e { b; }; }")

	p = newParser("try { a; } catch e { b; } finally { c; };")
	ok(t, p, "fn() { try { a; } catch e { b; } finally { c; }; }")

	p = newParser("try { a; } finally { c; };")
	ok(t, p, "fn() { try { a; } finally { c; }; }")

	p = newParser("try;")
	fail(t, p, "Unexpected Token ';' at (1, 4)")

	p = newParser("try {}")
	fail(t, p, "Invalid Try Expression at (1, 1)")
}

func TestInvoke(t *testing.T) {
	p := newParser("a()")
	okExpr(t, p, "a()")

	p = newParser("a(1)")
	okExpr(t, p, "a(1)")

	p = newParser("a(1, 2, 3)")
	okExpr(t, p, "a(1, 2, 3)")
}

func TestStruct(t *testing.T) {
	p := newParser("struct{}")
	okExpr(t, p, "struct {  }")

	p = newParser("struct{a:1}")
	okExpr(t, p, "struct { a: 1 }")

	p = newParser("struct{a:1,b:2}")
	okExpr(t, p, "struct { a: 1, b: 2 }")

	p = newParser("struct{a:1,b:2,c:3}")
	okExpr(t, p, "struct { a: 1, b: 2, c: 3 }")

	p = newParser("struct{a:1,b:2,c:struct{d:3}}")
	okExpr(t, p, "struct { a: 1, b: 2, c: struct { d: 3 } }")

	p = newParser("struct{a:1, b: fn(x) { y + x;} }")
	okExpr(t, p, "struct { a: 1, b: fn(x) { (y + x); } }")

	p = newParser("struct{a:1, b: fn(x) { y + x;}, c: struct {d:3} }")
	okExpr(t, p, "struct { a: 1, b: fn(x) { (y + x); }, c: struct { d: 3 } }")

	p = newParser("a.b")
	okExpr(t, p, "a.b")

	p = newParser("a.b = 3")
	okExpr(t, p, "(a.b = 3)")

	p = newParser("let a.b = 3;")
	fail(t, p, "Unexpected Token '.' at (1, 6)")

	p = newParser("this")
	okExpr(t, p, "this")

	p = newParser("struct{a:this + true, b: this}")
	okExpr(t, p, "struct { a: (this + true), b: this }")

	p = newParser("a = this")
	okExpr(t, p, "(a = this)")

	p = newParser("struct{ a: this }")
	okExpr(t, p, "struct { a: this }")

	p = newParser("struct{ a: this == 2 }")
	okExpr(t, p, "struct { a: (this == 2) }")

	p = newParser("a = this.b = 3")
	okExpr(t, p, "(a = (this.b = 3))")

	p = newParser("struct { a: this.b = 3 }")
	okExpr(t, p, "struct { a: (this.b = 3) }")

	p = newParser("b = this")
	okExpr(t, p, "(b = this)")

	p = newParser("struct { a: b = this }")
	okExpr(t, p, "struct { a: (b = this) }")

	p = newParser("a = struct { x: 8 }.x = 5")
	okExpr(t, p, "(a = (struct { x: 8 }.x = 5))")

	p = newParser("this = b")
	fail(t, p, "Unexpected Token '=' at (1, 6)")

	////////////

	p = newParser("struct { a: prop { fn() { x; } } }")
	okExpr(t, p, "struct { a: prop { fn() { x; } } }")

	p = newParser("struct { a: prop { || => x } }")
	okExpr(t, p, "struct { a: prop { fn() { x; } } }")

	p = newParser("struct { a: prop { |x| => x } }")
	fail(t, p, "Invalid Property Getter at (1, 20)")

	p = newParser("struct { a: prop { |x,y| => x } }")
	fail(t, p, "Invalid Property Getter at (1, 20)")

	////////////

	p = newParser("struct { a: prop { fn() { x; }, fn(y) { y; } } }")
	okExpr(t, p, "struct { a: prop { fn() { x; }, fn(y) { y; } } }")

	p = newParser("struct { a: prop { || => x, } }")
	fail(t, p, "Unexpected Token '}' at (1, 29)")

	p = newParser("struct { a: prop { || => x, || => y} }")
	fail(t, p, "Invalid Property Setter at (1, 29)")

	p = newParser("struct { a: prop { || => x, |y| => y } }")
	okExpr(t, p, "struct { a: prop { fn() { x; }, fn(y) { y; } } }")

	p = newParser("struct { a: prop { || => x, |x,y| => y} }")
	fail(t, p, "Invalid Property Setter at (1, 29)")
}

func TestPrimarySuffix(t *testing.T) {
	p := newParser("a.b()")
	okExpr(t, p, "a.b()")

	p = newParser("a.b.c")
	okExpr(t, p, "a.b.c")

	p = newParser("a.b().c")
	okExpr(t, p, "a.b().c")

	p = newParser("['a'][0]")
	okExpr(t, p, "[ 'a' ][0]")

	p = newParser("a[[]]")
	okExpr(t, p, "a[[  ]]")

	p = newParser("a[:b]")
	okExpr(t, p, "a[:b]")
	p = newParser("a[:]")
	fail(t, p, "Unexpected Token ']' at (1, 4)")

	p = newParser("a[b:]")
	okExpr(t, p, "a[b:]")
	p = newParser("a[b:}")
	fail(t, p, "Unexpected Token '}' at (1, 5)")

	p = newParser("a[b:c]")
	okExpr(t, p, "a[b:c]")
	p = newParser("a[b:c:]")
	fail(t, p, "Unexpected Token ':' at (1, 6)")

	p = newParser("a[b][c[:x]].d[y:].e().f[g[i:j]]")
	okExpr(t, p, "a[b][c[:x]].d[y:].e().f[g[i:j]]")
}

func okExprPos(t *testing.T, p *Parser, expectBegin ast.Pos, expectEnd ast.Pos) {

	expr, err := parseExpression(p)
	if err != nil {
		t.Error(err, " != nil")
	}

	if expr.Begin() != expectBegin {
		t.Error(expr.Begin(), " != ", expectBegin)
	}

	if expr.End() != expectEnd {
		t.Error(expr.End(), " != ", expectEnd)
	}
}

func okPos(t *testing.T, p *Parser, expectBegin ast.Pos, expectEnd ast.Pos) {

	mod, err := p.ParseModule()
	if err != nil {
		t.Error(err, " != nil")
		panic("okPos")
	}

	if len(mod.Body.Statements) != 1 {
		t.Error("node count", len(mod.Body.Statements))
		panic("okPos")
	}

	n := mod.Body.Statements[0]
	if n.Begin() != expectBegin {
		t.Error(n.Begin(), " != ", expectBegin)
		panic("okPos")
	}

	if n.End() != expectEnd {
		t.Error(n.End(), " != ", expectEnd)
		panic("okPos")
	}
}

func TestPos(t *testing.T) {
	p := newParser("1.23")
	okExprPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 4})

	p = newParser("-1")
	okExprPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 2})

	p = newParser("null + true")
	okExprPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 11})

	p = newParser("a1")
	okExprPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 2})

	p = newParser("a = \n3")
	okExprPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 2, Col: 1})

	p = newParser("a(b,c)")
	okExprPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 6})

	p = newParser("struct{}")
	okExprPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 8})

	p = newParser("struct { a: 1 }")
	okExprPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 15})

	p = newParser("   this")
	okExprPos(t, p, ast.Pos{Line: 1, Col: 4}, ast.Pos{Line: 1, Col: 7})

	p = newParser("a.b")
	okExprPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 3})

	p = newParser("a.b = 2")
	okExprPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 7})

	p = newParser(`
fn() {
    return x;
}`)
	okExprPos(t, p, ast.Pos{Line: 2, Col: 1}, ast.Pos{Line: 4, Col: 1})

	p = newParser("const a = 1;")
	okPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 11})

	p = newParser("let a = 1\n;")
	okPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 9})

	p = newParser("break;")
	okPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 5})

	p = newParser("\n  continue;")
	okPos(t, p, ast.Pos{Line: 2, Col: 3}, ast.Pos{Line: 2, Col: 10})

	p = newParser("while true { 42; \n};")
	okPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 2, Col: 1})

	p = newParser("if 0 {};")
	okPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 7})

	p = newParser("if 0 {} else {};")
	okPos(t, p, ast.Pos{Line: 1, Col: 1}, ast.Pos{Line: 1, Col: 15})
}

func TestList(t *testing.T) {
	p := newParser("[]")
	okExpr(t, p, "[  ]")

	p = newParser("[a]")
	okExpr(t, p, "[ a ]")

	p = newParser("[a, b]")
	okExpr(t, p, "[ a, b ]")

	p = newParser("[a, b, [], struct{z:1   }]")
	okExpr(t, p, "[ a, b, [  ], struct { z: 1 } ]")
}

func TestSet(t *testing.T) {
	p := newParser("set {}")
	okExpr(t, p, "set {  }")

	p = newParser("set { a }")
	okExpr(t, p, "set { a }")

	p = newParser("set { a, b }")
	okExpr(t, p, "set { a, b }")

	p = newParser("set { a, b, dict {c: 1} }")
	okExpr(t, p, "set { a, b, dict { c: 1 } }")
}

func TestDict(t *testing.T) {
	p := newParser("dict{}")
	okExpr(t, p, "dict {  }")

	p = newParser("dict{'a': 1}")
	okExpr(t, p, "dict { 'a': 1 }")

	p = newParser("dict { 'a': 1, null: [  ], [  ]: dict {  } }")
	okExpr(t, p, "dict { 'a': 1, null: [  ], [  ]: dict {  } }")
}

func TestBuiltin(t *testing.T) {
	p := newParser("print(12)")
	okExpr(t, p, "print(12)")

	p = newParser("str([])")
	okExpr(t, p, "str([  ])")

	p = newParser("a = println")
	okExpr(t, p, "(a = println)")

	p = newParser("len - null")
	okExpr(t, p, "(len - null)")

	p = newParser("ch = chan()")
	okExpr(t, p, "(ch = chan())")
}

func TestTuple(t *testing.T) {
	p := newParser("(a, b)")
	okExpr(t, p, "(a, b)")

	p = newParser("(a, b, struct { z: 1 })[2]")
	okExpr(t, p, "(a, b, struct { z: 1 })[2]")
}

func TestSwitch(t *testing.T) {

	p := newParser("switch { case a: x; };")
	ok(t, p, "fn() { switch { case a: x; }; }")

	p = newParser("switch { case a, b: x; y; };")
	ok(t, p, "fn() { switch { case a, b: x; y; }; }")

	p = newParser("switch { case a: x; case b: y; };")
	ok(t, p, "fn() { switch { case a: x; case b: y; }; }")

	p = newParser("switch true { case a: x; default: false; y; };")
	ok(t, p, "fn() { switch true { case a: x; default: false; y; }; }")

	p = newParser("switch { case a: x; case b: y; default: z; };")
	ok(t, p, "fn() { switch { case a: x; case b: y; default: z; }; }")

	p = newParser("switch { }")
	fail(t, p, "Unexpected Token '}' at (1, 10)")

	p = newParser("switch { case a: x;")
	fail(t, p, "Unexpected EOF at (1, 20)")

	p = newParser("switch { default: x; }")
	fail(t, p, "Unexpected Token 'default' at (1, 10)")

	p = newParser("switch { case case a: x; }")
	fail(t, p, "Unexpected Token 'case' at (1, 15)")

	p = newParser("switch { case z, x; }")
	fail(t, p, "Unexpected Token ';' at (1, 19)")

	p = newParser("switch { case a, b, c: }")
	fail(t, p, "Invalid SwitchStmt Expression at (1, 22)")

	p = newParser("switch { case a: b; default: }")
	fail(t, p, "Invalid SwitchStmt Expression at (1, 28)")
}

func TestLambda(t *testing.T) {

	p := newParser("|| => true")
	okExpr(t, p, "fn() { true; }")

	p = newParser("| | => true")
	okExpr(t, p, "fn() { true; }")

	p = newParser("|x| => true")
	okExpr(t, p, "fn(x) { true; }")

	p = newParser("|x, y| => true")
	okExpr(t, p, "fn(x, y) { true; }")

	p = newParser("|x, y, z| => true")
	okExpr(t, p, "fn(x, y, z) { true; }")
}

func TestSpawn(t *testing.T) {

	p := newParser("go foo();")
	ok(t, p, "fn() { go foo(); }")

	p = newParser("go false(a,b,c);")
	ok(t, p, "fn() { go false(a, b, c); }")

	p = newParser("go foo;")
	fail(t, p, "Unexpected Token ';' at (1, 7)")
}

func TestImport(t *testing.T) {
	p := newParser("import a;")
	ok(t, p, "fn() { import a; }")

	p = newParser("import a; import b;let z = 3; ")
	ok(t, p, "fn() { import a; import b; let z = 3; }")

	p = newParser("import a, b,c")
	ok(t, p, "fn() { import a, b, c; }")

	p = newParser("let z = 3; import a;")
	fail(t, p, "Unexpected Token 'import' at (1, 12)")
}

func TestLookaheadLF(t *testing.T) {
	p := newParser(`
fn
a() {}`)
	ok(t, p, "fn() { fn a() {  }; }")

	p = newParser(`
|a

| => a*a`)
	ok(t, p, "fn() { fn(a) { (a * a); }; }")
}

func parse(t *testing.T, source string) {
	p := newParser(source)
	_, err := p.ParseModule()
	if err != nil {
		t.Error(err, " != nil")
		panic("ok")
	}
}

func TestSemicolons(t *testing.T) {
	parse(t, `
let a
const b`)

	parse(t, `
let a = struct {
    x: 8,
    y: 5,
    plus:  fn() { return this.x + this.y; },
    minus: || => this.x - this.y, 
    plus:  fn() { return this.x + this.y 
	}
}
let b = a.plus()
let c = a.minus()
`)

	parse(t, `
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
`)

	parse(t, `
let fibonacciGenerator = fn() {
    let x = 1
    let y = 1
    return fn() {
        let z = x
        x = y
        y = x + z
        return z
    }
}

println("Fibonacci series:")
let nextFib = fibonacciGenerator()
for i in range(0, 10) {
    println(i, " == ", nextFib())
}
`)

	p := newParser(`
let a
let b
= 1;
const c = 
1; const d
`)
	ok(t, p, "fn() { let a; let b = 1; const c = 1; const d; }")

	p = newParser(`
a
++`)
	ok(t, p, "fn() { a++; }")

	p = newParser(`
break
throw
z
continue`)
	ok(t, p, "fn() { break; throw z; continue; }")

	p = newParser(`
return
z
`)
	ok(t, p, "fn() { return z; }")

	p = newParser(`
fn a() {
    return b();
}
fn b() {
    return 42;
}
`)
	ok(t, p, "fn() { fn a() { return b(); }; fn b() { return 42; }; }")

	parse(t, `
let a = fn() { 42; }
let b = fn(x) {
    let c = fn(y) {
        y * 7
    }
    x * x + c(x)
}
`)
}
