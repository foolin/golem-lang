// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	"github.com/mjarmy/golem-lang/analyzer"
	g "github.com/mjarmy/golem-lang/core"
	o "github.com/mjarmy/golem-lang/core/opcodes"
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

func ok(t *testing.T, mod *g.BytecodeModule, expect *g.BytecodeModule) {

	if !reflect.DeepEqual(mod.Pool, expect.Pool) {
		t.Error(mod, " != ", expect)
	}

	if len(mod.Templates) != len(expect.Templates) {
		t.Error(mod.Templates, " != ", expect.Templates)
	}

	for i := 0; i < len(mod.Templates); i++ {

		mt := mod.Templates[i]
		et := expect.Templates[i]

		if (mt.Arity != et.Arity) || (mt.NumCaptures != et.NumCaptures) || (mt.NumLocals != et.NumLocals) {
			t.Error(mod, " != ", expect)
		}

		if !reflect.DeepEqual(mt.OpCodes, et.OpCodes) {
			t.Error("OpCodes: ", mod, " != ", expect)
		}

		// checking LineNumberTable is optional
		if et.LineNumberTable != nil {
			if !reflect.DeepEqual(mt.LineNumberTable, et.LineNumberTable) {
				t.Error("LineNumberTable: ", mod, " != ", expect)
			}
		}
	}
}

var builtInMgr = g.NewBuiltinManager(g.CommandLineBuiltins)

func newAnalyzer(source string) analyzer.Analyzer {

	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner, builtInMgr.Contains)
	mod, err := parser.ParseModule()
	if err != nil {
		panic(err)
	}

	anl := analyzer.NewAnalyzer(mod)
	errors := anl.Analyze()
	if len(errors) > 0 {
		panic(err)
	}
	return anl
}

func newCompiler(analyzer analyzer.Analyzer) Compiler {
	return NewCompiler(analyzer, builtInMgr)
}

func contents() g.Struct {
	stc, _ := g.NewStruct([]g.Field{}, true)
	return stc
}

func TestExpression(t *testing.T) {

	mod := newCompiler(newAnalyzer("-2 + -1 + -0 + 0 + 1 + 2")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(int64(-2)),
			g.MakeInt(int64(2))},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CONST, 0, 0,
					o.LOAD_NEG_ONE,
					o.PLUS,
					o.LOAD_ZERO,
					o.PLUS,
					o.LOAD_ZERO,
					o.PLUS,
					o.LOAD_ONE,
					o.PLUS,
					o.LOAD_CONST, 0, 1,
					o.PLUS,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{16, 0}},
				nil}}, contents()})

	mod = newCompiler(newAnalyzer("(2 + 3) * -4 / 10")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(int64(2)),
			g.MakeInt(int64(3)),
			g.MakeInt(int64(-4)),
			g.MakeInt(int64(10))},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CONST, 0, 0,
					o.LOAD_CONST, 0, 1,
					o.PLUS,
					o.LOAD_CONST, 0, 2,
					o.MUL,
					o.LOAD_CONST, 0, 3,
					o.DIV,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{16, 0}},
				nil}}, contents()})

	mod = newCompiler(newAnalyzer("null / true + \nfalse")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_NULL,
					o.LOAD_TRUE,
					o.DIV,
					o.LOAD_FALSE,
					o.PLUS,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{4, 2},
					{5, 1},
					{6, 0}},
				nil}}, contents()})

	mod = newCompiler(newAnalyzer("'a' * 1.23e4")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.NewStr("a"),
			g.MakeFloat(float64(12300))},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CONST, 0, 0,
					o.LOAD_CONST, 0, 1,
					o.MUL,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{8, 0}},
				nil}}, contents()})

	mod = newCompiler(newAnalyzer("'a' == true")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.NewStr("a")},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CONST, 0, 0,
					o.LOAD_TRUE,
					o.EQ,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{6, 0}},
				nil}}, contents()})

	mod = newCompiler(newAnalyzer("true != false")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_TRUE, o.LOAD_FALSE, o.NE,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{4, 0}},
				nil}}, contents()})

	mod = newCompiler(newAnalyzer("true > false; true >= false")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_TRUE, o.LOAD_FALSE, o.GT,
					o.LOAD_TRUE, o.LOAD_FALSE, o.GTE,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{7, 0}},
				nil}}, contents()})

	mod = newCompiler(newAnalyzer("true < false; true <= false; true <=> false;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_TRUE, o.LOAD_FALSE, o.LT,
					o.LOAD_TRUE, o.LOAD_FALSE, o.LTE,
					o.LOAD_TRUE, o.LOAD_FALSE, o.CMP,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{10, 0}},
				nil}}, contents()})

	mod = newCompiler(newAnalyzer("let a = 2 && 3;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(int64(2)),
			g.MakeInt(int64(3))},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 1,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CONST, 0, 0,
					o.JUMP_FALSE, 0, 17,
					o.LOAD_CONST, 0, 1,
					o.JUMP_FALSE, 0, 17,
					o.LOAD_TRUE,
					o.JUMP, 0, 18,
					o.LOAD_FALSE,
					o.STORE_LOCAL, 0, 0,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{21, 0}},
				nil}}, contents()})

	mod = newCompiler(newAnalyzer("let a = 2 || 3;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(int64(2)),
			g.MakeInt(int64(3))},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 1,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CONST, 0, 0,
					o.JUMP_TRUE, 0, 13,
					o.LOAD_CONST, 0, 1,
					o.JUMP_FALSE, 0, 17,
					o.LOAD_TRUE,
					o.JUMP, 0, 18,
					o.LOAD_FALSE,
					o.STORE_LOCAL, 0, 0,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{21, 0}},
				nil}}, contents()})
}

func TestAssignment(t *testing.T) {

	mod := newCompiler(newAnalyzer("let a = 1;\nconst b = \n2;a = 3;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(2),
			g.MakeInt(3)},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 2,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_ONE,
					o.STORE_LOCAL, 0, 0,
					o.LOAD_CONST, 0, 0,
					o.STORE_LOCAL, 0, 1,
					o.LOAD_CONST, 0, 1,
					o.DUP,
					o.STORE_LOCAL, 0, 0,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{5, 3},
					{8, 2},
					{11, 3},
					{18, 0}},
				nil}}, contents()})
}

func TestShift(t *testing.T) {

	a := 0x1234
	high, low := byte((a>>8)&0xFF), byte(a&0xFF)

	if high != 0x12 || low != 0x34 {
		panic("shift")
	}

	var b int = int(high)<<8 + int(low)
	if b != a {
		panic("shift")
	}
}

func TestIf(t *testing.T) {

	source := "if (3 == 2) { let a = 42; }"
	anl := newAnalyzer(source)
	mod := newCompiler(anl).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(3),
			g.MakeInt(2),
			g.MakeInt(42)},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 1,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CONST, 0, 0,
					o.LOAD_CONST, 0, 1,
					o.EQ,
					o.JUMP_FALSE, 0, 17,
					o.LOAD_CONST, 0, 2,
					o.STORE_LOCAL, 0, 0,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{17, 0}},
				nil}}, contents()})

	source = `let a = 1
		if (false) {
		    let b = 2
		} else {
		    let c = 3
		}
		let d = 4`

	anl = newAnalyzer(source)
	mod = newCompiler(anl).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(2),
			g.MakeInt(3),
			g.MakeInt(4)},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 4,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_ONE,
					o.STORE_LOCAL, 0, 0,
					o.LOAD_FALSE,
					o.JUMP_FALSE, 0, 18,
					o.LOAD_CONST, 0, 0,
					o.STORE_LOCAL, 0, 1,
					o.JUMP, 0, 24,
					o.LOAD_CONST, 0, 1,
					o.STORE_LOCAL, 0, 2,
					o.LOAD_CONST, 0, 2,
					o.STORE_LOCAL, 0, 3,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{5, 2},
					{9, 3},
					{15, 4},
					{18, 5},
					{24, 7},
					{30, 0}},
				nil}}, contents()})
}

func TestWhile(t *testing.T) {

	source := "let a = 1; while (0 < 1) { let b = 2; }"
	mod := newCompiler(newAnalyzer(source)).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(2)},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 2,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_ONE,
					o.STORE_LOCAL, 0, 0,
					o.LOAD_ZERO,
					o.LOAD_ONE,
					o.LT,
					o.JUMP_FALSE, 0, 20,
					o.LOAD_CONST, 0, 0,
					o.STORE_LOCAL, 0, 1,
					o.JUMP, 0, 5,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{20, 0}},
				nil}}, contents()})

	source = "let a = 'z'; while (0 < 1) \n{ break; continue; let b = 2; }; let c = 3;"
	mod = newCompiler(newAnalyzer(source)).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.NewStr("z"),
			g.MakeInt(2),
			g.MakeInt(3)},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 3,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CONST, 0, 0,
					o.STORE_LOCAL, 0, 0,
					o.LOAD_ZERO,
					o.LOAD_ONE,
					o.LT,
					o.JUMP_FALSE, 0, 28,
					o.JUMP, 0, 28,
					o.JUMP, 0, 7,
					o.LOAD_CONST, 0, 1,
					o.STORE_LOCAL, 0, 1,
					o.JUMP, 0, 7,
					o.LOAD_CONST, 0, 2,
					o.STORE_LOCAL, 0, 2,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{13, 2},
					{34, 0}},
				nil}}, contents()})
}

func TestReturn(t *testing.T) {

	source := "let a = 1; return a \n- 2; a = 3;"
	anl := newAnalyzer(source)
	mod := newCompiler(anl).Compile()

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(2),
			g.MakeInt(3)},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 1,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_ONE,
					o.STORE_LOCAL, 0, 0,
					o.LOAD_LOCAL, 0, 0,
					o.LOAD_CONST, 0, 0,
					o.SUB,
					o.RETURN,
					o.LOAD_CONST, 0, 1,
					o.DUP,
					o.STORE_LOCAL, 0, 0,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{8, 2},
					{12, 1},
					{13, 2},
					{20, 0}},
				nil}}, contents()})
}

func TestFunc(t *testing.T) {

	source := `
let a = fn() { 42; }
let b = fn(x) {
    let c = fn(y) {
        y * 7
    }
    x * x + c(x)
}
`
	anl := newAnalyzer(source)
	mod := newCompiler(anl).Compile()

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(anl.BytecodeModule()))
	//fmt.Println(mod)

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(42),
			g.MakeInt(7)},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{0, 0, 2,
				[]byte{
					o.LOAD_NULL,
					o.NEW_FUNC, 0, 1,
					o.STORE_LOCAL, 0, 0,
					o.NEW_FUNC, 0, 2,
					o.STORE_LOCAL, 0, 1,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 2},
					{7, 3},
					{13, 0}},
				nil},
			&g.Template{0, 0, 0,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CONST, 0, 0,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 2},
					{4, 0}},
				nil},
			&g.Template{1, 0, 2,
				[]byte{
					o.LOAD_NULL,
					o.NEW_FUNC, 0, 3,
					o.STORE_LOCAL, 0, 1,
					o.LOAD_LOCAL, 0, 0,
					o.LOAD_LOCAL, 0, 0,
					o.MUL,
					o.LOAD_LOCAL, 0, 1,
					o.LOAD_LOCAL, 0, 0,
					o.INVOKE, 0, 1,
					o.PLUS,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 4},
					{7, 7},
					{24, 0}},
				nil},
			&g.Template{1, 0, 1,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_LOCAL, 0, 0,
					o.LOAD_CONST, 0, 1,
					o.MUL,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 5},
					{8, 0}},
				nil}}, contents()})

	source = `
let a = fn() { }
let b = fn(x) { x; }
let c = fn(x, y) { let z = 4; x * y * z; }
a()
b(1)
c(2, 3)
`
	anl = newAnalyzer(source)
	mod = newCompiler(anl).Compile()

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(anl.BytecodeModule()))
	//fmt.Println(mod)

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(2),
			g.MakeInt(3),
			g.MakeInt(4)},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{0, 0, 3,
				[]byte{
					o.LOAD_NULL,
					o.NEW_FUNC, 0, 1,
					o.STORE_LOCAL, 0, 0,
					o.NEW_FUNC, 0, 2,
					o.STORE_LOCAL, 0, 1,
					o.NEW_FUNC, 0, 3,
					o.STORE_LOCAL, 0, 2,
					o.LOAD_LOCAL, 0, 0,
					o.INVOKE, 0, 0,
					o.LOAD_LOCAL, 0, 1,
					o.LOAD_ONE,
					o.INVOKE, 0, 1,
					o.LOAD_LOCAL, 0, 2,
					o.LOAD_CONST, 0, 0,
					o.LOAD_CONST, 0, 1,
					o.INVOKE, 0, 2,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 2},
					{7, 3},
					{13, 4},
					{19, 5},
					{25, 6},
					{32, 7},
					{44, 0}},
				nil},

			&g.Template{0, 0, 0,
				[]byte{
					o.LOAD_NULL,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0}},
				nil},

			&g.Template{1, 0, 1,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_LOCAL, 0, 0,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 3},
					{4, 0}},
				nil},

			&g.Template{2, 0, 3,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CONST, 0, 2,
					o.STORE_LOCAL, 0, 2,
					o.LOAD_LOCAL, 0, 0,
					o.LOAD_LOCAL, 0, 1,
					o.MUL,
					o.LOAD_LOCAL, 0, 2,
					o.MUL,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 4},
					{18, 0}},
				nil}}, contents()})
}

func TestCapture(t *testing.T) {

	source := `
const accumGen = fn(n) {
    return fn(i) {
        n = n + i
        return n
    }
}`

	anl := newAnalyzer(source)
	mod := newCompiler(anl).Compile()

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{0, 0, 1,
				[]byte{
					o.LOAD_NULL,
					o.NEW_FUNC, 0, 1,
					o.STORE_LOCAL, 0, 0,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 2},
					{7, 0}},
				nil},
			&g.Template{1, 0, 1,
				[]byte{
					o.LOAD_NULL,
					o.NEW_FUNC, 0, 2,
					o.FUNC_LOCAL, 0, 0,
					o.RETURN,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 3},
					{8, 0}},
				nil},
			&g.Template{1, 1, 1,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CAPTURE, 0, 0,
					o.LOAD_LOCAL, 0, 0,
					o.PLUS,
					o.DUP,
					o.STORE_CAPTURE, 0, 0,
					o.LOAD_CAPTURE, 0, 0,
					o.RETURN,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 4},
					{12, 5},
					{16, 0}},
				nil}}, contents()})

	source = `
let z = 2
const accumGen = fn(n) {
    return fn(i) {
        n = n + i + z
        return n
    }
}`

	anl = newAnalyzer(source)
	mod = newCompiler(anl).Compile()
	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(anl.BytecodeModule()))
	//fmt.Println(mod)

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(2)},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{0, 0, 2,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CONST, 0, 0,
					o.STORE_LOCAL, 0, 0,
					o.NEW_FUNC, 0, 1,
					o.FUNC_LOCAL, 0, 0,
					o.STORE_LOCAL, 0, 1,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 2},
					{7, 3},
					{16, 0}},
				nil},
			&g.Template{1, 1, 1,
				[]byte{
					o.LOAD_NULL,
					o.NEW_FUNC, 0, 2,
					o.FUNC_LOCAL, 0, 0,
					o.FUNC_CAPTURE, 0, 0,
					o.RETURN,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 4},
					{11, 0}},
				nil},
			&g.Template{1, 2, 1,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CAPTURE, 0, 0,
					o.LOAD_LOCAL, 0, 0,
					o.PLUS,
					o.LOAD_CAPTURE, 0, 1,
					o.PLUS,
					o.DUP,
					o.STORE_CAPTURE, 0, 0,
					o.LOAD_CAPTURE, 0, 0,
					o.RETURN,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 5},
					{16, 6},
					{20, 0}},
				nil}}, contents()})
}

func TestPostfix(t *testing.T) {

	source := `
let a = 10
let b = 20
let c = a++
let d = b--
`
	anl := newAnalyzer(source)
	mod := newCompiler(anl).Compile()
	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(anl.BytecodeModule()))
	//fmt.Println(mod)

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(int64(10)),
			g.MakeInt(int64(20))},
		nil,
		[][]g.Field{},
		[]*g.Template{
			&g.Template{
				0, 0, 4,
				[]byte{
					o.LOAD_NULL,
					o.LOAD_CONST, 0, 0,
					o.STORE_LOCAL, 0, 0,
					o.LOAD_CONST, 0, 1,
					o.STORE_LOCAL, 0, 1,
					o.LOAD_LOCAL, 0, 0,
					o.DUP,
					o.LOAD_ONE,
					o.PLUS,
					o.STORE_LOCAL, 0, 0,
					o.STORE_LOCAL, 0, 2,
					o.LOAD_LOCAL, 0, 1,
					o.DUP,
					o.LOAD_NEG_ONE,
					o.PLUS,
					o.STORE_LOCAL, 0, 1,
					o.STORE_LOCAL, 0, 3,
					o.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 2},
					{7, 3},
					{13, 4},
					{25, 5},
					{37, 0}},
				nil}}, contents()})
}

func TestPool(t *testing.T) {
	pool := g.EmptyHashMap()

	tassert(t, poolIndex(pool, g.MakeInt(4)) == 0)
	tassert(t, poolIndex(pool, g.NewStr("a")) == 1)
	tassert(t, poolIndex(pool, g.MakeFloat(1.0)) == 2)
	tassert(t, poolIndex(pool, g.NewStr("a")) == 1)
	tassert(t, poolIndex(pool, g.MakeInt(4)) == 0)

	slice := makePoolSlice(pool)
	tassert(t, reflect.DeepEqual(
		slice,
		[]g.Basic{
			g.MakeInt(4),
			g.NewStr("a"),
			g.MakeFloat(1.0)}))
}

func TestTry(t *testing.T) {

	source := `
let a = 1
try {
    a++
}
finally {
    a++
}
assert(a == 2)
`
	anl := newAnalyzer(source)
	mod := newCompiler(anl).Compile()
	tassert(t, mod.Templates[0].ExceptionHandlers[0] ==
		g.ExceptionHandler{5, 14, -1, 14})

	source = `
try {
    try {
        3 / 0
    } catch e2 {
        assert(1,2)
    }
} catch e {
    println(e)
}
`
	anl = newAnalyzer(source)
	mod = newCompiler(anl).Compile()
	//	fmt.Println("----------------------------")
	//	fmt.Println(source)
	//	fmt.Println("----------------------------")
	//	fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//	fmt.Println(mod)
}
