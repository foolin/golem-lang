// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this code code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"
	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/scanner"
	"reflect"
	"testing"
)

var builtins []*g.BuiltinEntry = append(
	g.SandboxBuiltins,
	g.CommandLineBuiltins...)
var builtinMgr = g.NewBuiltinManager(builtins)

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
		panic("assertion failure")
	}
}

func testCompile(t *testing.T, code string) *g.Module {

	source := &scanner.Source{Name: "foo", Path: "foo.glm", Code: code}
	mods, errs := compiler.CompileSourceFully(builtinMgr, source, nil)
	tassert(t, errs == nil)
	tassert(t, len(mods) == 1)

	return mods[0]
}

func okExpr(t *testing.T, code string, expect g.Value) {

	mod := testCompile(t, code)
	intp := NewInterpreter(".", builtinMgr, []*g.Module{mod})

	result, err := intp.InitModules()
	tassert(t, err == nil)
	tassert(t, len(result) == 1)

	b, err := result[0].Eq(intp, expect)
	tassert(t, err == nil)

	if !b.BoolVal() {
		t.Error(result, " != ", expect)
	}
}

func failExpr(t *testing.T, code string, expect string) {

	mod := testCompile(t, code)
	intp := NewInterpreter(".", builtinMgr, []*g.Module{mod})

	result, err := intp.InitModules()
	tassert(t, result == nil)

	tassert(t, err != nil)
	if err.Error() != expect {
		t.Error(err.Error(), " != ", expect)
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

func okMod(t *testing.T, code string, expect g.Value, expectRefs []*g.Ref) {

	mod := testCompile(t, code)
	intp := NewInterpreter(".", builtinMgr, []*g.Module{mod})

	result, err := intp.InitModules()
	tassert(t, err == nil)
	tassert(t, len(result) == 1)

	b, err := result[0].Eq(intp, expect)
	tassert(t, err == nil)

	if !b.BoolVal() {
		t.Error(result, " != ", expect)
	}

	if !reflect.DeepEqual(mod.Refs, expectRefs) {
		t.Error(mod.Refs, " != ", expectRefs)
	}
}

func interpret(mods []*g.Module) *Interpreter {
	intp := NewInterpreter(".", builtinMgr, mods)
	_, err := intp.InitModules()
	if err != nil {
		panic(err)
	}
	return intp
}

func fail(t *testing.T, code string, err g.Error, stack []string) *g.Module {

	mod := testCompile(t, code)
	intp := NewInterpreter(".", builtinMgr, []*g.Module{mod})

	expect := intp.makeErrorTrace(err, stack)

	result, err := intp.InitModules()
	if result != nil {
		panic(result)
	}

	if err.Error() != expect.Error() {
		t.Error(err, " != ", expect)
	}

	return mod
}

func failErr(t *testing.T, code string, expect g.Error) {

	mod := testCompile(t, code)
	intp := NewInterpreter(".", builtinMgr, []*g.Module{mod})

	result, err := intp.InitModules()
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
	failExpr(t, "!null;", "TypeMismatch: Expected Bool")

	failExpr(t, "!'a';", "TypeMismatch: Expected Bool")
	failExpr(t, "!1;", "TypeMismatch: Expected Bool")
	failExpr(t, "!1.0;", "TypeMismatch: Expected Bool")

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
	failExpr(t, "12  && false;", "TypeMismatch: Expected Bool")

	okExpr(t, "true  || true;", g.True)
	okExpr(t, "true  || false;", g.True)
	okExpr(t, "false || true;", g.True)
	okExpr(t, "false || false;", g.False)
	okExpr(t, "true  || 12;", g.True)
	failExpr(t, "12  || true;", "TypeMismatch: Expected Bool")

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
			&g.Ref{Val: g.NewInt(3)},
			&g.Ref{Val: g.NewInt(2)}})

	okMod(t, `
let a = 1
a = a + 41
const B = a / 6
let c = B + 3
c = (c + a)/13
`,
		g.NewInt(4),
		[]*g.Ref{
			&g.Ref{Val: g.NewInt(42)},
			&g.Ref{Val: g.NewInt(7)},
			&g.Ref{Val: g.NewInt(4)}})

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
			&g.Ref{Val: g.NewInt(4)},
			&g.Ref{Val: g.NewInt(8)},
			&g.Ref{Val: g.NewInt(16)}})

	okMod(t, `
let a = 1
let b = 2
a = b = 11
b = a %= 4
`,
		g.NewInt(3),
		[]*g.Ref{
			&g.Ref{Val: g.NewInt(3)},
			&g.Ref{Val: g.NewInt(3)}})
}

func TestIf(t *testing.T) {

	okMod(t, "let a = 1; if (true) { a = 2; }",
		g.NewInt(2),
		[]*g.Ref{&g.Ref{Val: g.NewInt(2)}})

	okMod(t, "let a = 1; if (false) { a = 2; }",
		g.Null,
		[]*g.Ref{&g.Ref{Val: g.One}})

	okMod(t, "let a = 1; if (1 == 1) { a = 2; } else { a = 3; }; let b = 4;",
		g.NewInt(2),
		[]*g.Ref{
			&g.Ref{Val: g.NewInt(2)},
			&g.Ref{Val: g.NewInt(4)}})

	okMod(t, "let a = 1; if (1 == 2) { a = 2; } else { a = 3; }; const b = 4;",
		g.NewInt(3),
		[]*g.Ref{
			&g.Ref{Val: g.NewInt(3)},
			&g.Ref{Val: g.NewInt(4)}})
}

func TestWhile(t *testing.T) {

	okMod(t, `
let a = 1
while (a < 3) {
    a = a + 1
}`,
		g.NewInt(3),
		[]*g.Ref{&g.Ref{Val: g.NewInt(3)}})

	okMod(t, `
let a = 1
while (a < 11) {
    if (a == 4) { a = a + 2; break; }
    a = a + 1
}`,
		g.NewInt(6),
		[]*g.Ref{&g.Ref{Val: g.NewInt(6)}})

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
			&g.Ref{Val: g.NewInt(11)},
			&g.Ref{Val: g.NewInt(4)}})

	okMod(t, `
let a = 1
return a + 2
let b = 5`,
		g.NewInt(3),
		[]*g.Ref{
			&g.Ref{Val: g.One},
			&g.Ref{Val: g.Null}})
}

func TestStruct(t *testing.T) {

	code := `
let w = struct {}
let x = struct { a: 0 }
let y = struct { a: 1, b: 2 }
let z = struct { a: 3, b: 4, c: struct { d: 5 } }
`
	mod := testCompile(t, code)
	i := interpret([]*g.Module{mod})

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

	code = `
let x = struct { a: 5 }
let y = x.a
x.a = 6
`
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})

	okRef(t, i, mod.Refs[0], newStruct([]g.Field{
		g.NewField("a", false, g.NewInt(6))}))
	okRef(t, i, mod.Refs[1], g.NewInt(5))

	code = `
let a = struct {
	x: 8,
	y: 5,
	plus:  fn() { return this.x + this.y; },
	minus: fn() { return this.x - this.y; }
}
let b = a.plus()
let c = a.minus()
	`
	mod = testCompile(t, code)
	fmt.Println("----------------------------")
	fmt.Println(code)
	fmt.Println(mod)

	interpret([]*g.Module{mod})
	okRef(t, i, mod.Refs[2], g.NewInt(13))
	okRef(t, i, mod.Refs[3], g.NewInt(3))

	code = `
let a = null
a = struct { x: 8 }.x = 5
`
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})

	okRef(t, i, mod.Refs[0], g.NewInt(5))

	code = `
let a = struct { x: 8 }
assert(a has 'x')
assert(!(a has 'z'))
assert(a has 'x')
let b = struct { x: this has 'x', y: this has 'z' }
assert(b.x)
assert(!b.y)
`
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})
}

func TestErrStack(t *testing.T) {

	code := `
let divide = fn(x, y) {
    return x / y
}
let a = divide(3, 0)
`
	fail(t, code,
		g.DivideByZeroError(),
		[]string{
			"    at foo.glm:3",
			"    at foo.glm:5"})

	code = `
let foo = fn(n) { n + n; }
let a = foo(5, 6)
	`
	fail(t, code,
		g.ArityMismatchError("1", 2),
		[]string{
			"    at foo.glm:3"})
}

func TestPostfix(t *testing.T) {

	code := `
let a = 10
let b = 20
let c = a++
let d = b--
`
	mod := testCompile(t, code)
	i := interpret([]*g.Module{mod})

	okRef(t, i, mod.Refs[0], g.NewInt(11))
	okRef(t, i, mod.Refs[1], g.NewInt(19))
	okRef(t, i, mod.Refs[2], g.NewInt(10))
	okRef(t, i, mod.Refs[3], g.NewInt(20))

	code = `
let a = struct { x: 10 }
let b = struct { y: 20 }
let c = a.x++
let d = b.y--
`
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})

	//fmt.Println("----------------------------")
	//fmt.Println(code)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[0], newStruct([]g.Field{
		g.NewField("x", false, g.NewInt(11))}))
	okRef(t, i, mod.Refs[1], newStruct([]g.Field{
		g.NewField("y", false, g.NewInt(19))}))
	okRef(t, i, mod.Refs[2], g.NewInt(10))
	okRef(t, i, mod.Refs[3], g.NewInt(20))
}

func TestTernaryIf(t *testing.T) {

	code := `
let a = true ? 3 : 4;
let b = false ? 5 : 6;
`
	mod := testCompile(t, code)
	i := interpret([]*g.Module{mod})

	//fmt.Println("----------------------------")
	//fmt.Println(code)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[0], g.NewInt(3))
	okRef(t, i, mod.Refs[1], g.NewInt(6))
}

func newRange(from int64, to int64, step int64) g.Range {
	r, err := g.NewRange(from, to, step)
	if err != nil {
		panic("invalid range")
	}
	return r
}

func TestBuiltin(t *testing.T) {

	code := `
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
	mod := testCompile(t, code)
	i := interpret([]*g.Module{mod})

	//fmt.Println("----------------------------")
	//fmt.Println(code)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[0], g.NewInt(3))
	okRef(t, i, mod.Refs[1], g.NewStr("[ 4, 5, 6 ]"))
	okRef(t, i, mod.Refs[2], newRange(0, 5, 1))
	okRef(t, i, mod.Refs[3], newRange(0, 5, 2))

	code = `
let a = assert(true);
`
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})
	okRef(t, i, mod.Refs[0], g.True)

	fail(t, "assert(1, 2);",
		g.ArityMismatchError("1", 2),
		[]string{
			"    at foo.glm:1"})

	fail(t, "assert(1);",
		g.TypeMismatchError("Expected Bool"),
		[]string{
			"    at foo.glm:1"})

	fail(t, "assert(1 == 2);",
		g.AssertionFailedError(),
		[]string{
			"    at foo.glm:1"})
}

func TestDecl(t *testing.T) {

	code := `
let a, b = 0;
const c = 1, d;
`
	mod := testCompile(t, code)
	i := interpret([]*g.Module{mod})

	//fmt.Println("----------------------------")
	//fmt.Println(code)
	//fmt.Println(mod)

	okRef(t, i, mod.Refs[0], g.Null)
	okRef(t, i, mod.Refs[1], g.Zero)
	okRef(t, i, mod.Refs[2], g.One)
	okRef(t, i, mod.Refs[3], g.Null)
}

func TestFor(t *testing.T) {

	code := `
let a = 0;
for n in [1,2,3] {
    a += n;
}
assert(a == 6);
`
	mod := testCompile(t, code)
	interpret([]*g.Module{mod})

	code = `
let keys = '';
let values = 0;
for (k, v)  in dict {'a': 1, 'b': 2, 'c': 3} {
    keys += k;
    values += v;
}
assert(keys == 'bac');
assert(values == 6);
`
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})

	code = `
let entries = '';
for e in dict {'a': 1, 'b': 2, 'c': 3} {
    entries += str(e);
}
assert(entries == '(b, 2)(a, 1)(c, 3)');
`
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})

	code = `
let keys = '';
let values = 0;
for (k, v)  in [('a', 1), ('b', 2), ('c', 3)] {
    keys += k;
    values += v;
}
assert(keys == 'abc');
assert(values == 6);
`
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})

	code = "for (k, v)  in [1, 2, 3] {}"
	fail(t, code,
		g.TypeMismatchError("Expected Tuple"),
		[]string{"    at foo.glm:1"})

	code = "for (a, b, c)  in [('a', 1), ('b', 2), ('c', 3)] {}"
	fail(t, code,
		g.InvalidArgumentError("Expected Tuple of length 3"),
		[]string{"    at foo.glm:1"})
}

func TestSwitch(t *testing.T) {

	code := `
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
	mod := testCompile(t, code)
	interpret([]*g.Module{mod})

	code = `
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
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})

	code = `
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
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})
}

func TestGetField(t *testing.T) {

	code := "null.bogus;"
	fail(t, code,
		g.NullValueError(),
		[]string{"    at foo.glm:1"})

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

	code := `
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
	fail(t, code,
		g.DivideByZeroError(),
		[]string{
			"    at foo.glm:4"})

	code = `
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
	fail(t, code,
		g.DivideByZeroError(),
		[]string{
			"    at foo.glm:5"})

	code = `
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
	fail(t, code,
		g.DivideByZeroError(),
		[]string{
			"    at foo.glm:6"})

	code = `
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
	//mod = testCompile(t, code)
	//fmt.Println("----------------------------")
	//fmt.Println(code)
	//fmt.Println(mod)

	fail(t, code,
		g.DivideByZeroError(),
		[]string{
			"    at foo.glm:6",
			"    at foo.glm:15"})

	code = `
let b = fn() {
    try {
    } finally {
        return 1;
    }
    return 2;
};
assert(b() == 1);
`
	mod := testCompile(t, code)
	interpret([]*g.Module{mod})

	code = `
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
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})

	code = `
try {
    assert(1,2,3);
} finally {
}
`
	fail(t, code,
		g.ArityMismatchError("1", 3),
		[]string{
			"    at foo.glm:3"})

	code = `
try {
    assert(1,2,3);
} finally {
    1/0;
}
`
	fail(t, code,
		g.DivideByZeroError(),
		[]string{
			"    at foo.glm:5"})
}

func TestCatch(t *testing.T) {

	code := `
try {
    3 / 0;
} catch e {
    assert(e.kind == "DivideByZero");
    assert(!(e has "msg"));
    assert(e.stackTrace == ['    at foo.glm:3']);
}
`
	mod := testCompile(t, code)
	interpret([]*g.Module{mod})

	code = `
try {
    try {
        3 / 0;
    } catch e2 {
        assert();
    }
} catch e {
    assert(e.kind == "ArityMismatch");
    assert(e.msg == "Expected 1 params, got 0");
    assert(e.stackTrace == ['    at foo.glm:6']);
}
`
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})

	code = `
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
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})

	code = `
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
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})

}

func TestCatchFinally(t *testing.T) {

	code := `
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
	mod := testCompile(t, code)
	interpret([]*g.Module{mod})

	code = `
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
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})

	code = `
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
	mod = testCompile(t, code)
	interpret([]*g.Module{mod})
}

func TestThrow(t *testing.T) {

	code := `
try {
    throw struct { foo: 'zork' };
} catch e {
    assert(e.foo == 'zork');
    assert(e.stackTrace == ['    at foo.glm:3']);
}
`
	mod := testCompile(t, code)
	interpret([]*g.Module{mod})
}

func TestIntrinsicAssign(t *testing.T) {
	code := `
try {
    [].join = 456;
} catch e {
    assert(e.kind == 'TypeMismatch');
    assert(e.msg == "Expected Struct");
}
`
	mod := testCompile(t, code)
	interpret([]*g.Module{mod})
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

	code := `
let a = 0;
const b = 1;
fn main(args) {}
`
	mod := testCompile(t, code)
	i := interpret([]*g.Module{mod})
	tassert(t, reflect.DeepEqual(mod.Contents.FieldNames(), []string{"b", "a", "main"}))

	v, err := mod.Contents.GetField(i, g.NewStr("a"))
	okVal(t, v, err, g.Zero)

	v, err = mod.Contents.GetField(i, g.NewStr("b"))
	okVal(t, v, err, g.One)

	v, err = mod.Contents.GetField(i, g.NewStr("main"))
	tassert(t, err == nil)
	f, ok := v.(g.BytecodeFunc)
	tassert(t, ok)
	tassert(t, f.Template().Arity.Kind == g.FixedArity)
	tassert(t, f.Template().Arity.RequiredParams == 1)

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

	code := `
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
	mod := testCompile(t, code)
	interpret([]*g.Module{mod})
}

func TestRawString(t *testing.T) {

	code :=
		"let s = `\n" +
			"ab\n" +
			"cd\n" +
			"efgh\n" +
			"`\n" +
			"\tassert(s[1:3] == 'ab')\n" +
			"\tassert(s[4:6] == 'cd')\n" +
			"\tassert(s[7:-1] == 'efgh')"

	mod := testCompile(t, code)
	interpret([]*g.Module{mod})
}

func TestImport(t *testing.T) {

	srcMain := &scanner.Source{Name: "foo", Path: "foo.glm", Code: `
import a, b, c
assert([1, 2, 3] == [a.x, b.y, c.z])
`}
	sourceMap := map[string]*scanner.Source{
		"a": &scanner.Source{Name: "a", Path: "a.glm", Code: "import c; let x = 1;"},
		"b": &scanner.Source{Name: "b", Path: "b.glm", Code: "import c; let y = c.z - 1;"},
		"c": &scanner.Source{Name: "c", Path: "c.glm", Code: "let z = 3;"},
	}
	resolver := func(moduleName string) (*scanner.Source, error) {
		if src, ok := sourceMap[moduleName]; ok {
			return src, nil
		}
		return nil, g.UndefinedModuleError(moduleName)
	}

	mods, errs := compiler.CompileSourceFully(builtinMgr, srcMain, resolver)
	tassert(t, errs == nil)
	tassert(t, len(mods) == 4)
	interpret(mods)
}
