// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"
	"github.com/mjarmy/golem-lang/analyzer"
	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/parser"
	"github.com/mjarmy/golem-lang/scanner"
	"reflect"
	"testing"
)

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
	}
}

func okExpr(t *testing.T, source string, expect g.Value) {
	mod := newCompiler(source).Compile()
	intp := NewInterpreter(mod, builtInMgr, importResolver())

	result, err := intp.Init()
	if err != nil {
		panic(err)
	}

	b, err := result.Eq(intp, expect)
	if err != nil {
		panic("okExpr")
	}
	if !b.BoolVal() {
		t.Error(result, " != ", expect)
		panic("okExpr")
	}
}

func okRef(t *testing.T, intp *Interpreter, ref *g.Ref, expect g.Value) {
	b, err := ref.Val.Eq(intp, expect)
	if err != nil {
		panic("okRef")
	}
	if !b.BoolVal() {
		t.Error(ref.Val, " != ", expect)
	}
}

func okMod(t *testing.T, source string, expectResult g.Value, expectRefs []*g.Ref) {
	mod := newCompiler(source).Compile()
	intp := NewInterpreter(mod, builtInMgr, importResolver())

	result, err := intp.Init()
	if err != nil {
		panic(err)
	}

	b, err := result.Eq(intp, expectResult)
	if err != nil {
		panic("okMod")
	}
	if !b.BoolVal() {
		t.Error(result, " != ", expectResult)
	}

	if !reflect.DeepEqual(mod.Refs, expectRefs) {
		t.Error(mod.Refs, " != ", expectRefs)
	}
}

func failExpr(t *testing.T, source string, expect string) {

	mod := newCompiler(source).Compile()
	intp := NewInterpreter(mod, builtInMgr, importResolver())

	result, err := intp.Init()
	if result != nil {
		panic(result)
	}

	if err.Error() != expect {
		t.Error(err.Error(), " != ", expect)
	}
}

func fail(t *testing.T, source string, err g.Error, stack []string) *g.BytecodeModule {

	mod := newCompiler(source).Compile()
	intp := NewInterpreter(mod, builtInMgr, importResolver())

	expect := intp.makeErrorTrace(err, stack)

	result, err := intp.Init()
	if result != nil {
		panic(result)
	}

	if err.Error() != expect.Error() {
		t.Error(err, " != ", expect)
	}

	return mod
}

func failErr(t *testing.T, source string, expect g.Error) {

	mod := newCompiler(source).Compile()
	intp := NewInterpreter(mod, builtInMgr, importResolver())

	result, err := intp.Init()
	if result != nil {
		panic(result)
	}

	if err.Error() != expect.Error() {
		t.Error(err, " != ", expect)
	}
}

func newStruct(entries []g.Field) g.Struct {

	stc, err := g.NewStruct(entries, false)
	if err != nil {
		panic("invalid struct")
	}
	return stc
}

var builtInMgr = g.NewBuiltinManager(g.CommandLineBuiltins)

func newCompiler(source string) compiler.Compiler {
	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner, builtInMgr.Contains)
	mod, err := parser.ParseModule()
	if err != nil {
		panic(err.Error())
	}
	anl := analyzer.NewAnalyzer(mod)
	errors := anl.Analyze()
	if len(errors) > 0 {
		panic(fmt.Sprintf("%v", errors))
	}

	return compiler.NewCompiler(anl, builtInMgr)
}

type testModule struct {
	name     string
	contents g.Struct
}

func (t *testModule) GetModuleName() string { return t.name }
func (t *testModule) GetContents() g.Struct { return t.contents }

func importResolver() func(name string) (g.Module, g.Error) {
	stc, err := g.NewStruct([]g.Field{
		g.NewField("a", true, g.One)},
		false)
	if err != nil {
		panic("invalid struct")
	}

	foo := &testModule{"foo", stc}
	return func(name string) (g.Module, g.Error) {
		if name == "foo" {
			return foo, nil
		}
		return nil, g.UndefinedModuleError(name)
	}
}

func interpret(mod *g.BytecodeModule) *Interpreter {
	intp := NewInterpreter(mod, builtInMgr, importResolver())
	_, err := intp.Init()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic("interpreter failed")
	}
	return intp
}

func TestExpressions(t *testing.T) {

	okExpr(t, "(2 + 3) * -4 / 10;", g.NewInt(-2))

	okExpr(t, "(2*2*2*2 + 2*3*(8 - 1) + 2) / (17 - 2*2*2 - -1);", g.NewInt(6))

	okExpr(t, "true + 'a';", g.NewStr("truea"))
	okExpr(t, "'a' + true;", g.NewStr("atrue"))
	okExpr(t, "'a' + null;", g.NewStr("anull"))
	okExpr(t, "null + 'a';", g.NewStr("nulla"))

	failExpr(t, "true + null;", "TypeMismatch: Expected Number Type")
	failExpr(t, "1 + null;", "TypeMismatch: Expected Number Type")
	failExpr(t, "null + 1;", "TypeMismatch: Expected Number Type")

	okExpr(t, "true == 'a';", g.False)
	okExpr(t, "3 * 7 + 4 == 5 * 5;", g.True)
	okExpr(t, "1 != 1;", g.False)
	okExpr(t, "1 != 2;", g.True)

	okExpr(t, "!false;", g.True)
	okExpr(t, "!true;", g.False)
	failExpr(t, "!null;", "TypeMismatch: Expected 'Bool'")

	failExpr(t, "!'a';", "TypeMismatch: Expected 'Bool'")
	failExpr(t, "!1;", "TypeMismatch: Expected 'Bool'")
	failExpr(t, "!1.0;", "TypeMismatch: Expected 'Bool'")

	okExpr(t, "1 < 2;", g.True)
	okExpr(t, "1 <= 2;", g.True)
	okExpr(t, "1 > 2;", g.False)
	okExpr(t, "1 >= 2;", g.False)

	okExpr(t, "2 < 2;", g.False)
	okExpr(t, "2 <= 2;", g.True)
	okExpr(t, "2 > 2;", g.False)
	okExpr(t, "2 >= 2;", g.True)

	okExpr(t, "1 <=> 2;", g.NewInt(-1))
	okExpr(t, "2 <=> 2;", g.Zero)
	okExpr(t, "2 <=> 1;", g.One)

	okExpr(t, "true  && true;", g.True)
	okExpr(t, "true  && false;", g.False)
	okExpr(t, "false && true;", g.False)
	okExpr(t, "false && 12;", g.False)
	failExpr(t, "12  && false;", "TypeMismatch: Expected 'Bool'")

	okExpr(t, "true  || true;", g.True)
	okExpr(t, "true  || false;", g.True)
	okExpr(t, "false || true;", g.True)
	okExpr(t, "false || false;", g.False)
	okExpr(t, "true  || 12;", g.True)
	failExpr(t, "12  || true;", "TypeMismatch: Expected 'Bool'")

	okExpr(t, "~0;", g.NewInt(-1))

	okExpr(t, "8 % 2;", g.NewInt(8%2))
	okExpr(t, "8 & 2;", g.NewInt(int64(8)&int64(2)))
	okExpr(t, "8 | 2;", g.NewInt(8|2))
	okExpr(t, "8 ^ 2;", g.NewInt(8^2))
	okExpr(t, "8 << 2;", g.NewInt(8<<2))
	okExpr(t, "8 >> 2;", g.NewInt(8>>2))

	okExpr(t, "[true][0];", g.True)
	okExpr(t, "'abc'[1];", g.NewStr("b"))
	okExpr(t, "'abc'[-1];", g.NewStr("c"))
	failExpr(t, "[true][2];", "IndexOutOfBounds: 2")

	okExpr(t, "'abc'[1:];", g.NewStr("bc"))
	okExpr(t, "'abc'[:1];", g.NewStr("a"))
	okExpr(t, "'abcd'[1:3];", g.NewStr("bc"))
	okExpr(t, "'abcd'[1:1];", g.NewStr(""))

	okExpr(t, "[6,7,8][1:];", g.NewList([]g.Value{g.NewInt(7), g.NewInt(8)}))
	okExpr(t, "[6,7,8][:1];", g.NewList([]g.Value{g.NewInt(6)}))
	okExpr(t, "[6,7,8,9][1:3];", g.NewList([]g.Value{g.NewInt(7), g.NewInt(8)}))
	okExpr(t, "[6,7,8,9][1:1];", g.NewList([]g.Value{}))

	okExpr(t, "struct{a: 1} has 'a';", g.True)
	okExpr(t, "struct{a: 1} has 'b';", g.False)

	failExpr(t, "struct{a: 1, a: 2};", "DuplicateField: Field 'a' is a duplicate")

	okExpr(t, "struct{} == struct{};", g.True)
	okExpr(t, "struct{a:1} == struct{a:1};", g.True)
	okExpr(t, "struct{a:1,b:2} == struct{a:1,b:2};", g.True)
	okExpr(t, "struct{a:1} != struct{a:1,b:2};", g.True)
	okExpr(t, "struct{a:1,b:2} != struct{b:2};", g.True)
	okExpr(t, "struct{a:1,b:2} != struct{a:3,b:2};", g.True)
}

func TestAssignment(t *testing.T) {
	okMod(t, `
let a = 1
const B = 2
a = a + B
`,
		g.NewInt(3),
		[]*g.Ref{
			&g.Ref{g.NewInt(3)},
			&g.Ref{g.NewInt(2)}})

	okMod(t, `
let a = 1
a = a + 41
const B = a / 6
let c = B + 3
c = (c + a)/13
`,
		g.NewInt(4),
		[]*g.Ref{
			&g.Ref{g.NewInt(42)},
			&g.Ref{g.NewInt(7)},
			&g.Ref{g.NewInt(4)}})

	okMod(t, `
let a = 1
let b = a += 3
let c = ~0
c -= -2
c <<= 4
b *= 2
`,
		g.NewInt(8),
		[]*g.Ref{
			&g.Ref{g.NewInt(4)},
			&g.Ref{g.NewInt(8)},
			&g.Ref{g.NewInt(16)}})

	okMod(t, `
let a = 1
let b = 2
a = b = 11
b = a %= 4
`,
		g.NewInt(3),
		[]*g.Ref{
			&g.Ref{g.NewInt(3)},
			&g.Ref{g.NewInt(3)}})
}

func TestIf(t *testing.T) {

	okMod(t, "let a = 1; if (true) { a = 2; }",
		g.NewInt(2),
		[]*g.Ref{&g.Ref{g.NewInt(2)}})

	okMod(t, "let a = 1; if (false) { a = 2; }",
		g.NullValue,
		[]*g.Ref{&g.Ref{g.One}})

	okMod(t, "let a = 1; if (1 == 1) { a = 2; } else { a = 3; }; let b = 4;",
		g.NewInt(2),
		[]*g.Ref{
			&g.Ref{g.NewInt(2)},
			&g.Ref{g.NewInt(4)}})

	okMod(t, "let a = 1; if (1 == 2) { a = 2; } else { a = 3; }; const b = 4;",
		g.NewInt(3),
		[]*g.Ref{
			&g.Ref{g.NewInt(3)},
			&g.Ref{g.NewInt(4)}})
}

func TestWhile(t *testing.T) {

	//	source := `
	//a = 1;
	//while (a < 11) {
	//    if (a == 4) { a = a + 2; break; }
	//    a = a + 1;
	//}`
	//	mod := newCompiler(source).Compile()
	//	fmt.Println("----------------------------")
	//	fmt.Println(source)
	//	fmt.Println(mod)

	okMod(t, `
let a = 1
while (a < 3) {
    a = a + 1
}`,
		g.NewInt(3),
		[]*g.Ref{&g.Ref{g.NewInt(3)}})

	okMod(t, `
let a = 1
while (a < 11) {
    if (a == 4) { a = a + 2; break; }
    a = a + 1
}`,
		g.NewInt(6),
		[]*g.Ref{&g.Ref{g.NewInt(6)}})

	okMod(t, `
let a = 1
let b = 0
while (a < 11) {
    a = a + 1
    if (a > 5) { continue; }
    b = b + 1
}`,
		g.NewInt(11),
		[]*g.Ref{
			&g.Ref{g.NewInt(11)},
			&g.Ref{g.NewInt(4)}})

	okMod(t, `
let a = 1
return a + 2
let b = 5`,
		g.NewInt(3),
		[]*g.Ref{
			&g.Ref{g.One},
			&g.Ref{g.NullValue}})
}

func TestFunc(t *testing.T) {

	source := `
let a = fn(x) { x; }
let b = a(1)
`
	mod := newCompiler(source).Compile()

	i := interpret(mod)
	okRef(t, i, mod.Refs[1], g.One)

	source = `
let a = fn() { }
let b = fn(x) { x; }
let c = fn(x, y) { let z = 4; x * y * z; }
let d = a()
let e = b(1)
let f = c(b(2), 3)
`
	mod = newCompiler(source).Compile()

	interpret(mod)
	okRef(t, i, mod.Refs[3], g.NullValue)
	okRef(t, i, mod.Refs[4], g.One)
	okRef(t, i, mod.Refs[5], g.NewInt(24))

	source = `
let fibonacci = fn(n) {
    let x = 0
    let y = 1
    let i = 1
    while i < n {
        let z = x + y
        x = y
        y = z
        i = i + 1
    }
    return y
}
let a = fibonacci(1)
let b = fibonacci(2)
let c = fibonacci(3)
let d = fibonacci(4)
let e = fibonacci(5)
let f = fibonacci(6)
`
	mod = newCompiler(source).Compile()
	interpret(mod)
	okRef(t, i, mod.Refs[1], g.One)
	okRef(t, i, mod.Refs[2], g.One)
	okRef(t, i, mod.Refs[3], g.NewInt(2))
	okRef(t, i, mod.Refs[4], g.NewInt(3))
	okRef(t, i, mod.Refs[5], g.NewInt(5))
	okRef(t, i, mod.Refs[6], g.NewInt(8))

	source = `
let foo = fn(n) {
    let bar = fn(x) {
        return x * (x - 1)
    }
    return bar(n) + bar(n-1)
}
let a = foo(5)
`
	mod = newCompiler(source).Compile()
	interpret(mod)
	okRef(t, i, mod.Refs[1], g.NewInt(32))
}

func TestCapture(t *testing.T) {

	source := `
const accumGen = fn(n) {
    return fn(i) {
        n = n + i
        return n
    }
}
const a = accumGen(3)
let x = a(2)
let y = a(7)
`
	mod := newCompiler(source).Compile()
	i := interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[2], g.NewInt(5))
	okRef(t, i, mod.Refs[3], g.NewInt(12))

	source = `
let z = 2
const accumGen = fn(n) {
    return fn(i) {
        n = n + i
        n = n + z
        return n
    }
}
const a = accumGen(3)
let x = a(2)
z = 0
let y = a(1)
`
	mod = newCompiler(source).Compile()

	interpret(mod)

	okRef(t, i, mod.Refs[0], g.Zero)
	okRef(t, i, mod.Refs[3], g.NewInt(7))
	okRef(t, i, mod.Refs[4], g.NewInt(8))

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	source = `
const a = 123
const b = 456

fn foo() {
    assert(b == 456)
    assert(a == 123)
}
foo()
`
	mod = newCompiler(source).Compile()
	interpret(mod)
}

func TestStruct(t *testing.T) {

	source := `
let w = struct {}
let x = struct { a: 0 }
let y = struct { a: 1, b: 2 }
let z = struct { a: 3, b: 4, c: struct { d: 5 } }
`
	mod := newCompiler(source).Compile()
	i := interpret(mod)

	okRef(t, i, mod.Refs[0], newStruct([]g.Field{}))
	okRef(t, i, mod.Refs[1], newStruct([]g.Field{
		g.NewField("a", false, g.Zero)}))
	okRef(t, i, mod.Refs[2], newStruct([]g.Field{
		g.NewField("a", false, g.One),
		g.NewField("b", false, g.NewInt(2))}))
	okRef(t, i, mod.Refs[3], newStruct([]g.Field{
		g.NewField("a", false, g.NewInt(3)),
		g.NewField("b", false, g.NewInt(4)),
		g.NewField("c", false, newStruct([]g.Field{
			g.NewField("d", false, g.NewInt(5))}))}))

	source = `
let x = struct { a: 5 }
let y = x.a
x.a = 6
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[0], newStruct([]g.Field{
		g.NewField("a", false, g.NewInt(6))}))
	okRef(t, i, mod.Refs[1], g.NewInt(5))

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
	mod = newCompiler(source).Compile()
	interpret(mod)

	okRef(t, i, mod.Refs[2], g.NewInt(13))
	okRef(t, i, mod.Refs[3], g.NewInt(3))

	source = `
let a = null
a = struct { x: 8 }.x = 5
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	okRef(t, i, mod.Refs[0], g.NewInt(5))

	source = `
let a = struct { x: 8 }
assert(a has 'x')
assert(!(a has 'z'))
assert(a has 'x')
let b = struct { x: this has 'x', y: this has 'z' }
assert(b.x)
assert(!b.y)
`
	mod = newCompiler(source).Compile()
	interpret(mod)
}

func TestMerge(t *testing.T) {

	failExpr(t, "merge();", "ArityMismatch: Expected at least 2 params, got 0")
	failExpr(t, "merge(true);", "ArityMismatch: Expected at least 2 params, got 1")
	failExpr(t, "merge(struct{}, false);", "TypeMismatch: Expected 'Struct'")

	source := `
let a = struct { x: 1, y: 2}
let b = merge(struct { y: 3, z: 4}, a)
assert(b.x == 1)
assert(b.y == 3)
assert(b.z == 4)
a.x = 5
a.y = 6
assert(b.x == 5)
assert(b.y == 3)
assert(b.z == 4)
let c = merge(struct { w: 10}, b)
assert(c.w == 10)
assert(c.x == 5)
assert(c.y == 3)
assert(c.z == 4)
a.x = 7
b.z = 11
assert(c.w == 10)
assert(c.x == 7)
assert(c.y == 3)
assert(c.z == 11)
`
	mod := newCompiler(source).Compile()
	interpret(mod)
}

func TestErrStack(t *testing.T) {

	source := `
let divide = fn(x, y) {
    return x / y
}
let a = divide(3, 0)
`
	fail(t, source,
		g.DivideByZeroError(),
		[]string{
			"    at line 3",
			"    at line 5"})

	source = `
let foo = fn(n) { n + n; }
let a = foo(5, 6)
	`
	fail(t, source,
		g.ArityMismatchError("1", 2),
		[]string{
			"    at line 3"})
}

func TestPostfix(t *testing.T) {

	source := `
let a = 10
let b = 20
let c = a++
let d = b--
`
	mod := newCompiler(source).Compile()
	i := interpret(mod)

	okRef(t, i, mod.Refs[0], g.NewInt(11))
	okRef(t, i, mod.Refs[1], g.NewInt(19))
	okRef(t, i, mod.Refs[2], g.NewInt(10))
	okRef(t, i, mod.Refs[3], g.NewInt(20))

	source = `
let a = struct { x: 10 }
let b = struct { y: 20 }
let c = a.x++
let d = b.y--
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[0], newStruct([]g.Field{
		g.NewField("x", false, g.NewInt(11))}))
	okRef(t, i, mod.Refs[1], newStruct([]g.Field{
		g.NewField("y", false, g.NewInt(19))}))
	okRef(t, i, mod.Refs[2], g.NewInt(10))
	okRef(t, i, mod.Refs[3], g.NewInt(20))
}

func TestTernaryIf(t *testing.T) {

	source := `
let a = true ? 3 : 4;
let b = false ? 5 : 6;
`
	mod := newCompiler(source).Compile()
	i := interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[0], g.NewInt(3))
	okRef(t, i, mod.Refs[1], g.NewInt(6))
}

func TestList(t *testing.T) {

	source := `
let a = [];
let b = [true];
let c = [false,22];
let d = b[0];
b[0] = 33;
let e = c[1]++;
`
	mod := newCompiler(source).Compile()
	i := interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[0], g.NewList([]g.Value{}))
	okRef(t, i, mod.Refs[1], g.NewList([]g.Value{g.NewInt(33)}))
	okRef(t, i, mod.Refs[2], g.NewList([]g.Value{g.False, g.NewInt(23)}))
	okRef(t, i, mod.Refs[3], g.True)
	okRef(t, i, mod.Refs[4], g.NewInt(22))

	source = `
let a = [];
a.add(1);
assert(a == [1]);
a.add(2).add([3]);
assert(a == [1,2,[3]]);
let b = [];
b.add(4);
assert(b == [4]);
assert(a.add == a.add);
assert(b.add == b.add);
assert(a.add != b.add);
assert(b.add != a.add);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = [];
a.addAll([1,2]).addAll('bc');
assert(a == [1,2,'b','c']);
let b = [];
b.addAll(range(0,3));
b.addAll(dict { (true,false): 1, 'y': 2 });
assert(b == [ 0, 1, 2, ((true,false), 1), ('y', 2)]);
assert(a.addAll == a.addAll);
assert(b.addAll == b.addAll);
assert(a.addAll != b.addAll);
assert(b.addAll != a.addAll);
assert(a.add != a.addAll);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = "let a = []; a.addAll(false);"
	failErr(t, source, g.TypeMismatchError("Expected Iterable Type"))

	source = "let a = []; a.add(3,4);"
	failErr(t, source, g.ArityMismatchError("1", 2))

	source = `
let a = [];
assert(a.isEmpty());
a.add(1);
assert(!a.isEmpty());
a.clear();
assert(a.isEmpty());
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = [];
assert(!a.contains('x'));
assert(a.indexOf('x') == -1);
a = ['z', 'x'];
assert(a.contains('x'));
assert(a.indexOf('x') == 1);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = [];
assert(a.join() == '');
assert(a.join(',') == '');
a.add(1);
assert(a.join() == '1');
assert(a.join(',') == '1');
a.add(2);
assert(a.join() == '12');
assert(a.join(',') == '1,2');
a.add('abc');
assert(a.join() == '12abc');
assert(a.join(',') == '1,2,abc');
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let ls = [true, 0, 'abc'];
let types = ls.map(type);
assert(types == ['Bool', 'Int', 'Str']);
	`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let ls = [1, 2, 3, 4, 5];
let squares = ls.map(x => x * x);
let addedUp = ls.reduce(0, |acc, x| => acc + x);
let even = ls.filter(x => (x % 2 == 0));

assert(squares == [1, 4, 9, 16, 25]);
assert(addedUp == 15);
assert(even == [2, 4]);

ls.remove(2);
assert(ls == [1, 2, 4, 5]);
	`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = "['a'].remove(-1);"
	failErr(t, source, g.IndexOutOfBoundsError(-1))
	source = "['a'].remove(1);"
	failErr(t, source, g.IndexOutOfBoundsError(1))

	source = `
let ls = [3, 4, 5];
assert(ls[0] == 3);
assert(ls[2] == 5);
assert(ls[0:3] == [3, 4, 5]);
assert(ls[2:3] == [5]);
let a = [1, 2, 3];
let n = 0;
let b = a.map(fn(x) { n += x; x*x; });
assert(b == [1, 4, 9] && n == 6);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = "[3, 4, 5][-4];"
	failErr(t, source, g.IndexOutOfBoundsError(-1))
}

func TestDict(t *testing.T) {

	source := `
let d = dict {'x': 1, 'y': 2};
d['x'] = 0;
assert(d == dict {'y': 2, 'x': 0});
`
	mod := newCompiler(source).Compile()
	i := interpret(mod)

	source = `
let a = dict { 'x': 1, 'y': 2 };
let b = a['x'];
let c = a['z'];
a['x'] = -1;
let d = a['x'];
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[1], g.One)
	okRef(t, i, mod.Refs[2], g.NullValue)
	okRef(t, i, mod.Refs[3], g.NegOne)

	source = `
let a = dict {};
a.addAll([(1,2)]).addAll([(3,4)]);
assert(a == dict {1:2,3:4});
let b = dict {};
assert(a.addAll == a.addAll);
assert(b.addAll == b.addAll);
assert(a.addAll != b.addAll);
assert(b.addAll != a.addAll);
assert(a.clear != a.addAll);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = "let a = dict{}; a.addAll(false);"
	failErr(t, source, g.TypeMismatchError("Expected Iterable Type"))
	source = "let a = dict{}; a.addAll([false]);"
	failErr(t, source, g.TypeMismatchError("Expected Tuple"))
	source = "let a = dict{}; a.addAll([(1,2,3)]);"
	failErr(t, source, g.TupleLengthError(2, 3))

	source = "let a = dict{}; a[[1,2]];"
	failErr(t, source, g.TypeMismatchError("Expected Hashable Type"))

	source = "let a = dict{}; a[[1,2]] = 3;"
	failErr(t, source, g.TypeMismatchError("Expected Hashable Type"))

	source = "let a = dict{}; a.containsKey([1,2]);"
	failErr(t, source, g.TypeMismatchError("Expected Hashable Type"))

	source = `
let a = dict {};
assert(a.isEmpty());
a[1] = 2;
assert(!a.isEmpty());
a.clear();
assert(a.isEmpty());
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = dict {'z': 3};
assert(a.containsKey('z'));
assert(!a.containsKey('x'));
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let d = dict {'a': 1, 'b': 2};
assert(!d.remove('z'));
assert(d.remove('a'));
assert(d == dict {'b': 2});
assert(d.remove('b'));
assert(d == dict {});
assert(len(d) == 0);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let d = dict {'a': 1, 'b': 2};
assert(!d.remove('z'));
assert(d.remove('a'));
assert(d == dict {'b': 2});
assert(d.remove('b'));
assert(d == dict {});
assert(len(d) == 0);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
try {
    let d = dict {null:  'b'}
    assert(false)
} catch e {
    assert(e.kind == "NullValue")
}

try {
    let d = dict {[]:  'b'}
    assert(false)
} catch e {
    assert(e.kind == "TypeMismatch")
    assert(e.msg == "Expected Hashable Type")
}
`
	mod = newCompiler(source).Compile()
	interpret(mod)
}

func TestSet(t *testing.T) {

	source := `
let a = set {};
a.add(1);
assert(a == set {1});
a.add(2).add(3).add(2);
assert(a == set {1,2,3});
assert(set {3,1,2} == set {1,2,3});
assert(set {1,2} != set {1,2,3});
let b = set { 4 };
b.add(4);
assert(b == set { 4 });
assert(a.add == a.add);
assert(b.add == b.add);
assert(a.add != b.add);
assert(b.add != a.add);
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = set {};
a.addAll([1,2]).addAll('bc');
assert(a == set {1,2,'b','c'});
let b = set {};
b.addAll(range(0,3));
assert(b == set { 0, 1, 2 });
assert(a.addAll == a.addAll);
assert(b.addAll == b.addAll);
assert(a.addAll != b.addAll);
assert(b.addAll != a.addAll);
assert(a.add != a.addAll);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = "let a = set{}; a.addAll(false);"
	failErr(t, source, g.TypeMismatchError("Expected Iterable Type"))

	source = "let a = set{}; a.add(3,4);"
	failErr(t, source, g.ArityMismatchError("1", 2))

	source = "let a = set{}; a.add([1,2]);"
	failErr(t, source, g.TypeMismatchError("Expected Hashable Type"))

	source = "let a = set{}; a.contains([1,2]);"
	failErr(t, source, g.TypeMismatchError("Expected Hashable Type"))

	source = `
let a = set{};
assert(a.isEmpty());
a.add(1);
assert(!a.isEmpty());
a.clear();
assert(a.isEmpty());
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = set{};
assert(!a.contains('x'));
a = set {'z', 'x'};
assert(a.contains('x'));
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let s = set {'a', 'b', 'c'};
assert(!s.remove('z'));
assert(s.remove('a'));
assert(s == set {'c', 'b'});
assert(s.remove('b'));
assert(s == set {'c'});
assert(len(s) == 1);
assert(s.remove('c'));
assert(s == set {});
assert(len(s) == 0);
`
	mod = newCompiler(source).Compile()
	interpret(mod)
}

func newRange(from int64, to int64, step int64) g.Range {
	r, err := g.NewRange(from, to, step)
	if err != nil {
		panic("invalid range")
	}
	return r
}

func TestBuiltin(t *testing.T) {

	source := `
let a = len([4,5,6]);
let b = str([4,5,6]);
let c = range(0, 5);
let d = range(0, 5, 2);
print();
println();
print(a);
println(b);
print(a,b);
println(a,b);
assert(print == print);
assert(print != println);
`
	mod := newCompiler(source).Compile()
	i := interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[0], g.NewInt(3))
	okRef(t, i, mod.Refs[1], g.NewStr("[ 4, 5, 6 ]"))
	okRef(t, i, mod.Refs[2], newRange(0, 5, 1))
	okRef(t, i, mod.Refs[3], newRange(0, 5, 2))

	source = `
let a = assert(true);
`
	mod = newCompiler(source).Compile()
	interpret(mod)
	okRef(t, i, mod.Refs[0], g.True)

	fail(t, "assert(1, 2);",
		g.ArityMismatchError("1", 2),
		[]string{
			"    at line 1"})

	fail(t, "assert(1);",
		g.TypeMismatchError("Expected 'Bool'"),
		[]string{
			"    at line 1"})

	fail(t, "assert(1 == 2);",
		g.AssertionFailedError(),
		[]string{
			"    at line 1"})
}

func TestTuple(t *testing.T) {

	source := `
let a = (4,5);
let b = a[0];
let c = a[1];
`
	mod := newCompiler(source).Compile()
	i := interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[0], g.NewTuple([]g.Value{g.NewInt(4), g.NewInt(5)}))
	okRef(t, i, mod.Refs[1], g.NewInt(4))
	okRef(t, i, mod.Refs[2], g.NewInt(5))
}

func TestDecl(t *testing.T) {

	source := `
let a, b = 0;
const c = 1, d;
`
	mod := newCompiler(source).Compile()
	i := interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[0], g.NullValue)
	okRef(t, i, mod.Refs[1], g.Zero)
	okRef(t, i, mod.Refs[2], g.One)
	okRef(t, i, mod.Refs[3], g.NullValue)
}

func TestFor(t *testing.T) {

	source := `
let a = 0;
for n in [1,2,3] {
    a += n;
}
assert(a == 6);
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	source = `
let keys = '';
let values = 0;
for (k, v)  in dict {'a': 1, 'b': 2, 'c': 3} {
    keys += k;
    values += v;
}
assert(keys == 'bac');
assert(values == 6);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let entries = '';
for e in dict {'a': 1, 'b': 2, 'c': 3} {
    entries += str(e);
}
assert(entries == '(b, 2)(a, 1)(c, 3)');
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let keys = '';
let values = 0;
for (k, v)  in [('a', 1), ('b', 2), ('c', 3)] {
    keys += k;
    values += v;
}
assert(keys == 'abc');
assert(values == 6);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = "for (k, v)  in [1, 2, 3] {}"
	fail(t, source,
		g.TypeMismatchError("Expected 'Tuple'"),
		[]string{"    at line 1"})

	source = "for (a, b, c)  in [('a', 1), ('b', 2), ('c', 3)] {}"
	fail(t, source,
		g.InvalidArgumentError("Expected Tuple of length 3"),
		[]string{"    at line 1"})
}

func TestSwitch(t *testing.T) {

	source := `
let s = ''
for i in range(0, 4) {
    switch {
    case i == 0:
        s += 'a'

    case i == 1, i == 2:
        s += 'b'

    default:
        s += 'c'
    }
}
assert(s == 'abbc')
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	source = `
let s = ''
for i in range(0, 4) {
    switch {
    case i == 0, i == 1:
        s += 'a'

    case i == 2:
        s += 'b'
    }
}
assert(s == 'aab')
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let s = ''
for i in range(0, 4) {
    switch i {
    case 0, 1:
        s += 'a'

    case 2:
        s += 'b'
    }
}
assert(s == 'aab')
`
	mod = newCompiler(source).Compile()
	interpret(mod)
}

func TestGetField(t *testing.T) {

	source := "null.bogus;"
	fail(t, source,
		g.NullValueError(),
		[]string{"    at line 1"})

	err := g.NoSuchFieldError("bogus")

	failErr(t, "true.bogus;", err)
	failErr(t, "'a'.bogus;", err)
	failErr(t, "(0).bogus;", err)
	failErr(t, "(0.123).bogus;", err)

	failErr(t, "(1,2).bogus;", err)
	failErr(t, "range(1,2).bogus;", err)
	failErr(t, "[1,2].bogus;", err)
	failErr(t, "dict {'a':1}.bogus;", err)
	failErr(t, "struct {a:1}.bogus;", err)

	failErr(t, "(fn() {}).bogus;", err)
}

func TestFinally(t *testing.T) {

	source := `
let a = 1
try {
    3 / 0
} finally {
    a = 2
}
try {
    3 / 0
} finally {
    a = 3
}
`
	fail(t, source,
		g.DivideByZeroError(),
		[]string{
			"    at line 4"})

	source = `
let a = 1;
try {
    try {
        3 / 0;
    } finally {
        a++;
    }
} finally {
    a++;
}
`
	fail(t, source,
		g.DivideByZeroError(),
		[]string{
			"    at line 5"})

	source = `
let a = 1;
let b = fn() { a++; };
try {
    try {
        3 / 0;
    } finally {
        a++;
        b();
    }
} finally {
    a++;
}
`
	fail(t, source,
		g.DivideByZeroError(),
		[]string{
			"    at line 6"})

	source = `
let a = 1
let b = fn() { 
    try {
        try {
            3 / 0
        } finally {
            a++
        }
    } finally {
        a++
    }
}
try {
    b()
} finally {
    a++
}
`
	//mod = newCompiler(source).Compile()
	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	fail(t, source,
		g.DivideByZeroError(),
		[]string{
			"    at line 6",
			"    at line 15"})

	source = `
let b = fn() { 
    try {
    } finally {
        return 1;
    }
    return 2;
};
assert(b() == 1);
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = 1;
let b = fn() { 
    try {
        try {
        } finally {
            return 1;
        }
        a = 3;
    } finally {
        a = 2;
    }
};
assert(b() == 1);
assert(a == 1);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
try {
    assert(1,2,3);
} finally {
}
`
	fail(t, source,
		g.ArityMismatchError("1", 3),
		[]string{
			"    at line 3"})

	source = `
try {
    assert(1,2,3);
} finally {
    1/0;
}
`
	fail(t, source,
		g.DivideByZeroError(),
		[]string{
			"    at line 5"})
}

func TestCatch(t *testing.T) {

	source := `
try {
    3 / 0;
} catch e {
    assert(e.kind == "DivideByZero");
    assert(!(e has "msg"));
    assert(e.stackTrace == ['    at line 3']);
}
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	source = `
try {
    try {
        3 / 0;
    } catch e2 {
        assert();
    }
} catch e {
    assert(e.kind == "ArityMismatch");
    assert(e.msg == "Expected 1 params, got 0");
    assert(e.stackTrace == ['    at line 6']);
}
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = 0;
let b = 0;
try {
    3 / 0;
} catch e {
    a = 1;
}
try {
    3 / 0;
} catch e {
    b = 2;
}
assert(a == 1);
assert(b == 2);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
try {
    let s = set {'a', 'b', null}
    assert(false)
} catch e {
    assert(e.kind == "NullValue")
}

try {
    let s = set {'a', 'b', []}
    assert(false)
} catch e {
    assert(e.kind == "TypeMismatch")
    assert(e.msg == "Expected Hashable Type")
}
`
	mod = newCompiler(source).Compile()
	interpret(mod)

}

func TestCatchFinally(t *testing.T) {

	source := `
let a = 0;
try {
    3 / 0;
} catch e {
    a = 1;
} finally {
    a = 2;
}
assert(a == 2);
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = 0;
let f = fn() {
    try {
        3 / 0;
    } catch e {
        return 1;
    } finally {
        a = 2;
    }
};
let b = f();
assert(b == 1);
assert(a == 2);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = 0
let b = 0
try {
    try {
        3 / 0
    } catch e {
        assert(1,2,3)
    } finally {
        a = 1
    }
} catch e {
    b = 2
}
assert(a == 1)
assert(b == 2)
`
	mod = newCompiler(source).Compile()
	interpret(mod)
}

func TestThrow(t *testing.T) {

	source := `
try {
    throw struct { foo: 'zork' };
} catch e {
    assert(e.foo == 'zork');
    assert(e.stackTrace == ['    at line 3']);
}
`
	mod := newCompiler(source).Compile()
	interpret(mod)
}

func TestNamedFunc(t *testing.T) {

	source := `
fn a() {
    return b();
}
fn b() {
    return 42;
}
assert(a() == 42);
`
	mod := newCompiler(source).Compile()
	interpret(mod)
}

func TestLambda(t *testing.T) {

	source := `
let z = 5
let a = || => 3
let b = x => x * x
let c = |x, y| => (x + y)*z
assert(a() == 3)
assert(b(2) == 4)
assert(c(1, 2) == 15)
`
	mod := newCompiler(source).Compile()
	interpret(mod)
}

func TestGo(t *testing.T) {

	source := `
fn sum(a, c) {
	let total = 0;
	for v in a {
		total += v;
	}
    c.send(total);
}

let a = [7, 2, 8, -9, 4, 0];
let n = len(a) / 2;
let c = chan();

go sum(a[:n], c);
go sum(a[n:], c);
let x = c.recv();
let y = c.recv();
assert([x, y] == [-5, 17]);
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	source = `
let ch = chan(2);
ch.send(1);
ch.send(2);
assert([ch.recv(), ch.recv()] == [1, 2]);
`
	mod = newCompiler(source).Compile()
	interpret(mod)
}

func TestIntrinsicAssign(t *testing.T) {
	source := `
try {
    [].join = 456;
} catch e {
    assert(e.kind == 'TypeMismatch');
    assert(e.msg == "Expected 'Struct'");
}
`
	mod := newCompiler(source).Compile()
	interpret(mod)
}

func okVal(t *testing.T, val g.Value, err g.Error, expect g.Value) {

	if err != nil {
		t.Error(err, " != ", nil)
	}

	if !reflect.DeepEqual(val, expect) {
		t.Error(val, " != ", expect)
	}
}

func failVal(t *testing.T, val g.Value, err g.Error, expect string) {

	if val != nil {
		t.Error(val, " != ", nil)
	}

	if err == nil || err.Error() != expect {
		t.Error(err.Error(), " != ", expect)
	}
}

func TestModuleContents(t *testing.T) {

	source := `
let a = 0;
const b = 1;
fn main(args) {}
`
	mod := newCompiler(source).Compile()
	i := interpret(mod)
	tassert(t, reflect.DeepEqual(mod.Contents.FieldNames(), []string{"b", "a", "main"}))

	v, err := mod.Contents.GetField(i, g.NewStr("a"))
	okVal(t, v, err, g.Zero)

	v, err = mod.Contents.GetField(i, g.NewStr("b"))
	okVal(t, v, err, g.One)

	v, err = mod.Contents.GetField(i, g.NewStr("main"))
	tassert(t, err == nil)
	f, ok := v.(g.BytecodeFunc)
	tassert(t, ok)
	tassert(t, f.Template().Arity == 1)

	err = mod.Contents.SetField(i, g.NewStr("a"), g.NegOne)
	tassert(t, err == nil)
	v, err = mod.Contents.GetField(i, g.NewStr("a"))
	okVal(t, v, err, g.NegOne)

	err = mod.Contents.SetField(i, g.NewStr("b"), g.NegOne)
	failVal(t, nil, err, "ReadonlyField: Field 'b' is readonly")

	err = mod.Contents.SetField(i, g.NewStr("main"), g.NegOne)
	failVal(t, nil, err, "ReadonlyField: Field 'main' is readonly")
}

func TestTypeOf(t *testing.T) {

	source := `
assert(
    [type(true), type(""), type(0), type(0.0)] == 
    ["Bool", "Str", "Int", "Float"]);
assert(
    [type(fn(){}), type([]), type(range(0,1)), type((0,1))] == 
    ["Func", "List", "Range", "Tuple"]);
assert(
    [type(dict{}), type(set{}), type(struct{}), type(chan())] == 
    ["Dict", "Set", "Struct", "Chan"]);
`
	mod := newCompiler(source).Compile()
	interpret(mod)
}

func TestStr(t *testing.T) {

	source := `
assert('abc'.contains('b'));
assert(!'abc'.contains('z'));
assert('abc'.index('b') == 1);
assert('abc'.index('z') == -1);
assert('abc'.startsWith('a'));
assert(!'abc'.startsWith('z'));
assert('abc'.endsWith('c'));
assert(!'abc'.endsWith('z'));
assert('aaa'.replace('a', 'z') == 'zzz');
assert('aaa'.replace('a', 'z', 2) == 'zza');
assert('aaa'.replace('a', 'z', 0) == 'aaa');
assert('aaa'.replace('a', 'z', -1) == 'zzz');
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	fail(t, "'abc'.contains();", g.ArityMismatchError("1", 0), []string{"    at line 1"})
	fail(t, "'abc'.contains(1);", g.TypeMismatchError("Expected Str"), []string{"    at line 1"})

	fail(t, "'abc'.index();", g.ArityMismatchError("1", 0), []string{"    at line 1"})
	fail(t, "'abc'.index(1);", g.TypeMismatchError("Expected Str"), []string{"    at line 1"})

	fail(t, "'abc'.startsWith();", g.ArityMismatchError("1", 0), []string{"    at line 1"})
	fail(t, "'abc'.startsWith(1);", g.TypeMismatchError("Expected Str"), []string{"    at line 1"})

	fail(t, "'abc'.endsWith();", g.ArityMismatchError("1", 0), []string{"    at line 1"})
	fail(t, "'abc'.endsWith(1);", g.TypeMismatchError("Expected Str"), []string{"    at line 1"})

	fail(t, "'abc'.replace();", g.ArityMismatchError("at least 2", 0), []string{"    at line 1"})
	fail(t, "'abc'.replace(1,2,3,4);", g.ArityMismatchError("at most 3", 4), []string{"    at line 1"})
	fail(t, "'abc'.replace(0, 'a');", g.TypeMismatchError("Expected Str"), []string{"    at line 1"})
	fail(t, "'abc'.replace('a', 0);", g.TypeMismatchError("Expected Str"), []string{"    at line 1"})
	fail(t, "'abc'.replace('a', 'b', 'c');", g.TypeMismatchError("Expected Int"), []string{"    at line 1"})
}

func TestFreeze(t *testing.T) {

	fail(t, "freeze(null);", g.NullValueError(), []string{"    at line 1"})
	fail(t, "frozen(null);", g.NullValueError(), []string{"    at line 1"})

	source := `
assert(frozen(true));
assert(frozen('a'));
assert(frozen(1));
assert(frozen(1.0));
assert(frozen(range(1,2)));
assert(frozen(chan()));
assert(frozen(fn(){}));
assert(frozen((1,2)));

freeze(true);
freeze('a');
freeze(1);
freeze(1.0);
freeze(range(1,2));
freeze(chan());
freeze(fn(){});
freeze((1,2));
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	source = `
fn fail(f) {
    try {
        f();
        assert(false);
    } catch e {
        assert(e.kind == 'ImmutableValue');
    }
}

let ls = [1,2,3];
assert(!frozen(ls));

ls.clear();
ls.add('a');
ls.addAll(['b', 'c']);
ls.remove(1);
ls[1] = 'z';

assert(ls == ['a', 'z']);

freeze(ls);
assert(frozen(ls));
assert(ls == ['a', 'z']);

fail(|| => ls.clear());
fail(|| => ls.add('a'));
fail(|| => ls.addAll(['b', 'c']));
fail(|| => ls.remove(1));
fail(|| => ls[1] = 'z');
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
fn fail(f) {
    try {
        f();
        assert(false);
    } catch e {
        assert(e.kind == 'ImmutableValue');
    }
}

let d = dict {'x': 1, 'y': 2};
assert(!frozen(d));

d.clear();
d.addAll([('a', 1), ('b', 2), ('c', 3)]);
d.remove('c');
d['a'] = 0;
assert(d == dict {'a': 0, 'b': 2});

freeze(d);
assert(frozen(d));
assert(d == dict {'a': 0, 'b': 2});

fail(|| => d.clear());
fail(|| => d.addAll([('a', 1), ('b', 2), ('c', 3)]));
fail(|| => d.remove('c'));
fail(|| => d['a'] = 0);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
fn fail(f) {
    try {
        f();
        assert(false);
    } catch e {
        assert(e.kind == 'ImmutableValue');
    }
}

let s = set {'x', 'y'};
assert(!frozen(s));

s.clear();
s.addAll('a');
s.addAll(['a', 'b', 'c']);
s.remove('c');
assert(s == set {'a', 'b'});

freeze(s);
assert(frozen(s));
assert(s == set {'a', 'b'});

fail(|| => s.clear());
fail(|| => s.addAll('a'));
fail(|| => s.addAll(['a', 'b', 'c']));
fail(|| => s.remove('c'));
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
fn fail(f) {
    try {
        f();
        assert(false);
    } catch e {
        assert(e.kind == 'ImmutableValue');
    }
}

let s = struct {x: 1, y: 2};
s.y = 3;

freeze(s);
assert(frozen(s));
assert(s == struct {x: 1, y: 3});

fail(|| => s.y = 3);

assert(!frozen(merge(struct {}, struct{})));
assert(frozen(merge(s, struct{})));
assert(frozen(merge(struct {}, s)));
`
	mod = newCompiler(source).Compile()
	interpret(mod)
}

func TestFields(t *testing.T) {

	source := `
fn fail(f, kind) {
    try {
        f();
        assert(false);
    } catch e {
        assert(e.kind == kind);
    }
}

let s = struct {a: 1, b: 2};
assert(fields(s) == set {'a', 'b'});
assert(getval(s, 'a') == 1);
assert(setval(s, 'a', 3) == 3);

fail(|| => fields(0), 'TypeMismatch');
fail(|| => fields(0, 1), 'ArityMismatch');

fail(|| => getval(0, 1), 'TypeMismatch');
fail(|| => getval(s, 1), 'TypeMismatch');
fail(|| => getval(0), 'ArityMismatch');

fail(|| => setval(0, 1, 2), 'TypeMismatch');
fail(|| => setval(s, 1, 2), 'TypeMismatch');
fail(|| => setval(0), 'ArityMismatch');
fail(|| => setval(0, 1), 'ArityMismatch');
`
	mod := newCompiler(source).Compile()
	interpret(mod)
}

func TestArity(t *testing.T) {

	source := `
fn fail(f, kind) {
    try {
        f();
        assert(false);
    } catch e {
        assert(e.kind == kind);
    }
}
assert(arity(type) == struct { min: 1, max: 1 });
assert(arity(print) == struct { min: 0, max: -1 });
assert(arity(|x,y| => x + y) == struct { min: 2, max: 2 });

fail(|| => arity(0), 'TypeMismatch');
`
	mod := newCompiler(source).Compile()
	interpret(mod)
}

func TestRange(t *testing.T) {

	source := `
fn listify(r) {
    let ls = [];
    for n in r {
        ls.add(n);
    }
    return ls;
}
let a = range(0, 5);
let b = range(0, 5, 2);
let c = range(2, 14, 3);
let d = range(-1, -8, -3);
let e = range(2, 2);
let f = range(-1, -1, -1);
assert(listify(a) == [ 0, 1, 2, 3, 4 ]);
assert(listify(b) == [ 0, 2, 4 ]);
assert(listify(c) == [ 2, 5, 8, 11 ]);
assert(listify(d) == [ -1, -4, -7 ]);
assert(listify(e) == []);
assert(listify(f) == []);
assert([a.from(), a.to(), a.step(), a.count()] == [0, 5, 1, 5]);
assert([b.from(), b.to(), b.step(), b.count()] == [0, 5, 2, 3]);
assert([c.from(), c.to(), c.step(), c.count()] == [2, 14, 3, 4]);
assert([d.from(), d.to(), d.step(), d.count()] == [-1, -8, -3, 3]);
assert([e.from(), e.to(), e.step(), e.count()] == [2, 2, 1, 0]);

let i = 0;
while i < a.count() {
    assert(a[i] == i);
    i++;
} 
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)
}

func TestSlice(t *testing.T) {

	source := `
let ls = [3,4,5];
assert(ls[0:1] == [3]);
assert(ls[1:3] == [4,5]);
assert(ls[-2:-1] == [4]);
assert(ls[-5:-4] == []);
assert(ls[3:4] == []);
let s = '345';
assert(s[0:1] == '3');
assert(s[1:3] == '45');
assert(s[-2:-1] == '4');
assert(s[-5:-4] == '');
assert(s[3:4] == '');
`
	mod := newCompiler(source).Compile()
	interpret(mod)
}

func TestMultilineString(t *testing.T) {

	source :=
		"let s = `\n" +
			"ab\n" +
			"cd\n" +
			"efgh\n" +
			"`\n" +
			"\tassert(s[1:3] == 'ab')\n" +
			"\tassert(s[4:6] == 'cd')\n" +
			"\tassert(s[7:-1] == 'efgh')"

	mod := newCompiler(source).Compile()
	interpret(mod)
}

func TestUnicodeEscape(t *testing.T) {
	source := `
let s = '\u{1f496}\u{2665}\u{24}'
assert(s[0] == '💖')
assert(s[1] == '♥')
assert(s[2] == '$')
`

	mod := newCompiler(source).Compile()
	interpret(mod)
}

func TestImport(t *testing.T) {
	source := `
import foo
assert(foo.a == 1)
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	source = `
import bar
`
	fail(t, source,
		g.UndefinedModuleError("bar"),
		[]string{
			"    at line 2"})
}
