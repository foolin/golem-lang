// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this code code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	//"fmt"
	"reflect"
	//"strings"
	"testing"

	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
	"github.com/mjarmy/golem-lang/scanner"
)

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
		panic("tassert")
	}
}

func ok(t *testing.T, pool *bc.Pool, expect *bc.Pool) {

	if !reflect.DeepEqual(pool.Constants, expect.Constants) {
		t.Error(pool, " != ", expect)
	}

	if len(pool.Templates) != len(expect.Templates) {
		t.Error(pool.Templates, " != ", expect.Templates)
	}

	for i := 0; i < len(pool.Templates); i++ {

		mt := pool.Templates[i]
		et := expect.Templates[i]

		if !reflect.DeepEqual(mt.Arity, et.Arity) ||
			(mt.NumCaptures != et.NumCaptures) ||
			(mt.NumLocals != et.NumLocals) {
			t.Error(pool, " != ", expect)
		}

		if !reflect.DeepEqual(mt.Bytecodes, et.Bytecodes) {
			t.Error("Bytecodes: ", pool, " != ", expect)
		}

		// checking LineNumberTable is optional
		if et.LineNumberTable != nil {
			if !reflect.DeepEqual(mt.LineNumberTable, et.LineNumberTable) {
				t.Error("LineNumberTable: ", pool, " != ", expect)
			}
		}
	}
}

var builtins = []*g.Builtin{
	{"assert", g.BuiltinAssert},
	{"println", g.BuiltinPrintln},
}

func testCompile(t *testing.T, code string) *bc.Module {

	source := &scanner.Source{Name: "foo", Path: "foo.glm", Code: code}
	mod, err := CompileSource(source, builtins)
	tassert(t, err == nil)

	return mod
}

func fixedArity(numParams int) g.Arity {
	return g.Arity{
		Kind:     g.FixedArity,
		Required: uint16(numParams),
		Optional: 0,
	}
}

func TestExpression(t *testing.T) {

	mod := testCompile(t, "-2 + -1 + -0 + 0 + 1 + 2")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewInt(int64(-2)), g.NewInt(int64(2))},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   0,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadConst, 0, 0,
					bc.LoadNegOne,
					bc.Plus,
					bc.LoadZero,
					bc.Plus,
					bc.LoadZero,
					bc.Plus,
					bc.LoadOne,
					bc.Plus,
					bc.LoadConst, 0, 1,
					bc.Plus,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 16, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})

	mod = testCompile(t, "(2 + 3) * -4 / 10")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewInt(int64(2)), g.NewInt(int64(3)), g.NewInt(int64(-4)), g.NewInt(int64(10))},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   0,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadConst, 0, 0,
					bc.LoadConst, 0, 1,
					bc.Plus,
					bc.LoadConst, 0, 2,
					bc.Mul,
					bc.LoadConst, 0, 3,
					bc.Div,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 16, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})

	mod = testCompile(t, "null / true + \nfalse")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   0,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadNull,
					bc.LoadTrue,
					bc.Div,
					bc.LoadFalse,
					bc.Plus,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 4, LineNum: 2},
					{Index: 5, LineNum: 1},
					{Index: 6, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})

	mod = testCompile(t, "'a' * 1.23e4")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.MustStr("a"), g.NewFloat(float64(12300))},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   0,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadConst, 0, 0,
					bc.LoadConst, 0, 1,
					bc.Mul,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 8, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})

	mod = testCompile(t, "'a' == true")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.MustStr("a")},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   0,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadConst, 0, 0,
					bc.LoadTrue,
					bc.Eq,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 6, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})

	mod = testCompile(t, "true != false")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   0,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadTrue,
					bc.LoadFalse,
					bc.Ne,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 4, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})

	mod = testCompile(t, "true > false; true >= false")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   0,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadTrue,
					bc.LoadFalse,
					bc.Gt,
					bc.LoadTrue,
					bc.LoadFalse,
					bc.Gte,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 7, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})

	mod = testCompile(t, "true < false; true <= false; true <=> false;")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   0,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadTrue,
					bc.LoadFalse,
					bc.Lt,
					bc.LoadTrue,
					bc.LoadFalse,
					bc.Lte,
					bc.LoadTrue,
					bc.LoadFalse,
					bc.Cmp,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 10, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})

	mod = testCompile(t, "let a = 2 && 3;")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewInt(int64(2)), g.NewInt(int64(3))},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   1,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadConst, 0, 0,
					bc.JumpFalse, 0, 17,
					bc.LoadConst, 0, 1,
					bc.JumpFalse, 0, 17,
					bc.LoadTrue,
					bc.Jump, 0, 18,
					bc.LoadFalse,
					bc.StoreLocal, 0, 0,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 21, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})

	mod = testCompile(t, "let a = 2 || 3;")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewInt(int64(2)), g.NewInt(int64(3))},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   1,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadConst, 0, 0,
					bc.JumpTrue, 0, 13,
					bc.LoadConst, 0, 1,
					bc.JumpFalse, 0, 17,
					bc.LoadTrue,
					bc.Jump, 0, 18,
					bc.LoadFalse,
					bc.StoreLocal, 0, 0,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 21, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})
}

func TestAssignment(t *testing.T) {

	mod := testCompile(t, "let a = 1;\nconst b = \n2;a = 3;")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewInt(2), g.NewInt(3)},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   2,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadOne,
					bc.StoreLocal, 0, 0,
					bc.LoadConst, 0, 0,
					bc.StoreLocal, 0, 1,
					bc.LoadConst, 0, 1,
					bc.Dup,
					bc.StoreLocal, 0, 0,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 5, LineNum: 3},
					{Index: 8, LineNum: 2},
					{Index: 11, LineNum: 3},
					{Index: 18, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})
}

func TestShift(t *testing.T) {

	a := 0x1234
	high, low := byte((a>>8)&0xFF), byte(a&0xFF)

	if high != 0x12 || low != 0x34 {
		panic("shift")
	}

	var b = int(high)<<8 + int(low)
	if b != a {
		panic("shift")
	}
}

func TestIf(t *testing.T) {

	code := "if (3 == 2) { let a = 42; }"
	mod := testCompile(t, code)
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewInt(3), g.NewInt(2), g.NewInt(42)},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   1,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadConst, 0, 0,
					bc.LoadConst, 0, 1,
					bc.Eq,
					bc.JumpFalse, 0, 17,
					bc.LoadConst, 0, 2,
					bc.StoreLocal, 0, 0,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 17, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})

	code = `let a = 1
		if (false) {
		    let b = 2
		} else {
		    let c = 3
		}
		let d = 4`

	mod = testCompile(t, code)
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewInt(2), g.NewInt(3), g.NewInt(4)},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   4,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadOne,
					bc.StoreLocal, 0, 0,
					bc.LoadFalse,
					bc.JumpFalse, 0, 18,
					bc.LoadConst, 0, 0,
					bc.StoreLocal, 0, 1,
					bc.Jump, 0, 24,
					bc.LoadConst, 0, 1,
					bc.StoreLocal, 0, 2,
					bc.LoadConst, 0, 2,
					bc.StoreLocal, 0, 3,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 5, LineNum: 2},
					{Index: 9, LineNum: 3},
					{Index: 15, LineNum: 4},
					{Index: 18, LineNum: 5},
					{Index: 24, LineNum: 7},
					{Index: 30, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})
}

func TestWhile(t *testing.T) {

	code := "let a = 1; while (0 < 1) { let b = 2; }"
	mod := testCompile(t, code)
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewInt(2)},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   2,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadOne,
					bc.StoreLocal, 0, 0,
					bc.LoadZero,
					bc.LoadOne,
					bc.Lt,
					bc.JumpFalse, 0, 20,
					bc.LoadConst, 0, 0,
					bc.StoreLocal, 0, 1,
					bc.Jump, 0, 5,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 20, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})

	code = "let a = 'z'; while (0 < 1) \n{ break; continue; let b = 2; }; let c = 3;"
	mod = testCompile(t, code)
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.MustStr("z"), g.NewInt(2), g.NewInt(3)},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   3,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadConst, 0, 0,
					bc.StoreLocal, 0, 0,
					bc.LoadZero,
					bc.LoadOne,
					bc.Lt,
					bc.JumpFalse, 0, 28,
					bc.Jump, 0, 28,
					bc.Jump, 0, 7,
					bc.LoadConst, 0, 1,
					bc.StoreLocal, 0, 1,
					bc.Jump, 0, 7,
					bc.LoadConst, 0, 2,
					bc.StoreLocal, 0, 2,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 13, LineNum: 2},
					{Index: 34, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})
}

func TestReturn(t *testing.T) {

	code := "let a = 1; return a \n- 2; a = 3;"
	mod := testCompile(t, code)

	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewInt(2), g.NewInt(3)},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   1,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadOne,
					bc.StoreLocal, 0, 0,
					bc.LoadLocal, 0, 0,
					bc.LoadConst, 0, 0,
					bc.Sub,
					bc.Return,
					bc.LoadConst, 0, 1,
					bc.Dup,
					bc.StoreLocal, 0, 0,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 8, LineNum: 2},
					{Index: 12, LineNum: 1},
					{Index: 13, LineNum: 2},
					{Index: 20, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})
}

func TestFunc(t *testing.T) {

	code := `
let a = fn() { 42; }
let b = fn(x) {
    let c = fn(y) {
        y * 7
    }
    x * x + c(x)
}
`
	mod := testCompile(t, code)

	//fmt.Println("----------------------------")
	//fmt.Println(code)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(astMod.Module()))
	//fmt.Println(mod)

	ok(t, mod.Pool, &bc.Pool{
		Constants: []g.Basic{
			g.NewInt(42),
			g.NewInt(7)},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   2,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.NewFunc, 0, 1,
					bc.StoreLocal, 0, 0,
					bc.NewFunc, 0, 2,
					bc.StoreLocal, 0, 1,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 2},
					{Index: 7, LineNum: 3},
					{Index: 13, LineNum: 0}},
				ErrorHandlers: nil,
			},
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   0,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadConst, 0, 0,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 2},
					{Index: 4, LineNum: 0}},
				ErrorHandlers: nil,
			},
			&bc.FuncTemplate{
				Arity:       fixedArity(1),
				NumCaptures: 0,
				NumLocals:   2,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.NewFunc, 0, 3,
					bc.StoreLocal, 0, 1,
					bc.LoadLocal, 0, 0,
					bc.LoadLocal, 0, 0,
					bc.Mul,
					bc.LoadLocal, 0, 1,
					bc.LoadLocal, 0, 0,
					bc.Invoke, 0, 1,
					bc.Plus,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 4},
					{Index: 7, LineNum: 7},
					{Index: 24, LineNum: 0}},
				ErrorHandlers: nil,
			},
			&bc.FuncTemplate{
				Arity:       fixedArity(1),
				NumCaptures: 0,
				NumLocals:   1,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadLocal, 0, 0,
					bc.LoadConst, 0, 1,
					bc.Mul,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 5},
					{Index: 8, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})

	code = `
let a = fn() { }
let b = fn(x) { x; }
let c = fn(x, y) { let z = 4; x * y * z; }
a()
b(1)
c(2, 3)
`
	mod = testCompile(t, code)

	//fmt.Println("----------------------------")
	//fmt.Println(code)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(astMod.Module()))
	//fmt.Println(mod)

	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewInt(2), g.NewInt(3), g.NewInt(4)},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   3,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.NewFunc, 0, 1,
					bc.StoreLocal, 0, 0,
					bc.NewFunc, 0, 2,
					bc.StoreLocal, 0, 1,
					bc.NewFunc, 0, 3,
					bc.StoreLocal, 0, 2,
					bc.LoadLocal, 0, 0,
					bc.Invoke, 0, 0,
					bc.LoadLocal, 0, 1,
					bc.LoadOne,
					bc.Invoke, 0, 1,
					bc.LoadLocal, 0, 2,
					bc.LoadConst, 0, 0,
					bc.LoadConst, 0, 1,
					bc.Invoke, 0, 2,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 2},
					{Index: 7, LineNum: 3},
					{Index: 13, LineNum: 4},
					{Index: 19, LineNum: 5},
					{Index: 25, LineNum: 6},
					{Index: 32, LineNum: 7},
					{Index: 44, LineNum: 0}},
				ErrorHandlers: nil,
			},
			&bc.FuncTemplate{
				Arity:       fixedArity(0),
				NumCaptures: 0,
				NumLocals:   0,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0}},
				ErrorHandlers: nil,
			},
			&bc.FuncTemplate{
				Arity:       fixedArity(1),
				NumCaptures: 0,
				NumLocals:   1,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadLocal, 0, 0,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 3},
					{Index: 4, LineNum: 0}},
				ErrorHandlers: nil,
			},
			&bc.FuncTemplate{
				Arity:       fixedArity(2),
				NumCaptures: 0,
				NumLocals:   3,
				Bytecodes: []byte{
					bc.LoadNull,
					bc.LoadConst, 0, 2,
					bc.StoreLocal, 0, 2,
					bc.LoadLocal, 0, 0,
					bc.LoadLocal, 0, 1,
					bc.Mul,
					bc.LoadLocal, 0, 2,
					bc.Mul,
					bc.Return},
				LineNumberTable: []bc.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 4},
					{Index: 18, LineNum: 0}},
				ErrorHandlers: nil,
			}},
	})
}

func TestCapture(t *testing.T) {

	code := `
const accumGen = fn(n) {
    return fn(i) {
        n = n + i
        return n
    }
}`

	mod := testCompile(t, code)

	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{&bc.FuncTemplate{
			Arity:       fixedArity(0),
			NumCaptures: 0,
			NumLocals:   1,
			Bytecodes: []byte{
				bc.LoadNull,
				bc.NewFunc, 0, 1,
				bc.StoreLocal, 0, 0,
				bc.Return},
			LineNumberTable: []bc.LineNumberEntry{
				{Index: 0, LineNum: 0},
				{Index: 1, LineNum: 2},
				{Index: 7, LineNum: 0}},
			ErrorHandlers: nil,
		}, &bc.FuncTemplate{
			Arity:       fixedArity(1),
			NumCaptures: 0,
			NumLocals:   1,
			Bytecodes: []byte{
				bc.LoadNull,
				bc.NewFunc, 0, 2,
				bc.FuncLocal, 0, 0,
				bc.Return,
				bc.Return},
			LineNumberTable: []bc.LineNumberEntry{
				{Index: 0, LineNum: 0},
				{Index: 1, LineNum: 3},
				{Index: 8, LineNum: 0}},
			ErrorHandlers: nil,
		}, &bc.FuncTemplate{
			Arity:       fixedArity(1),
			NumCaptures: 1,
			NumLocals:   1,
			Bytecodes: []byte{
				bc.LoadNull,
				bc.LoadCapture, 0, 0,
				bc.LoadLocal, 0, 0,
				bc.Plus,
				bc.Dup,
				bc.StoreCapture, 0, 0,
				bc.LoadCapture, 0, 0,
				bc.Return,
				bc.Return},
			LineNumberTable: []bc.LineNumberEntry{
				{Index: 0, LineNum: 0},
				{Index: 1, LineNum: 4},
				{Index: 12, LineNum: 5},
				{Index: 16, LineNum: 0}},
			ErrorHandlers: nil,
		}},
	})

	code = `
let z = 2
const accumGen = fn(n) {
    return fn(i) {
        n = n + i + z
        return n
    }
}`

	mod = testCompile(t, code)
	//fmt.Println("----------------------------")
	//fmt.Println(code)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(astMod.Module()))
	//fmt.Println(mod)

	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewInt(2)},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{&bc.FuncTemplate{
			Arity:       fixedArity(0),
			NumCaptures: 0,
			NumLocals:   2,
			Bytecodes: []byte{
				bc.LoadNull,
				bc.LoadConst, 0, 0,
				bc.StoreLocal, 0, 0,
				bc.NewFunc, 0, 1,
				bc.FuncLocal, 0, 0,
				bc.StoreLocal, 0, 1,
				bc.Return},
			LineNumberTable: []bc.LineNumberEntry{
				{Index: 0, LineNum: 0},
				{Index: 1, LineNum: 2},
				{Index: 7, LineNum: 3},
				{Index: 16, LineNum: 0}},
			ErrorHandlers: nil,
		}, &bc.FuncTemplate{
			Arity:       fixedArity(1),
			NumCaptures: 1,
			NumLocals:   1,
			Bytecodes: []byte{
				bc.LoadNull,
				bc.NewFunc, 0, 2,
				bc.FuncLocal, 0, 0,
				bc.FuncCapture, 0, 0,
				bc.Return,
				bc.Return},
			LineNumberTable: []bc.LineNumberEntry{
				{Index: 0, LineNum: 0},
				{Index: 1, LineNum: 4},
				{Index: 11, LineNum: 0}},
			ErrorHandlers: nil,
		}, &bc.FuncTemplate{
			Arity:       fixedArity(1),
			NumCaptures: 2,
			NumLocals:   1,
			Bytecodes: []byte{
				bc.LoadNull,
				bc.LoadCapture, 0, 0,
				bc.LoadLocal, 0, 0,
				bc.Plus,
				bc.LoadCapture, 0, 1,
				bc.Plus,
				bc.Dup,
				bc.StoreCapture, 0, 0,
				bc.LoadCapture, 0, 0,
				bc.Return,
				bc.Return},
			LineNumberTable: []bc.LineNumberEntry{
				{Index: 0, LineNum: 0},
				{Index: 1, LineNum: 5},
				{Index: 16, LineNum: 6},
				{Index: 20, LineNum: 0}},
			ErrorHandlers: nil,
		}},
	})
}

func TestPostfix(t *testing.T) {

	code := `
let a = 10
let b = 20
let c = a++
let d = b--
`
	mod := testCompile(t, code)
	//fmt.Println("----------------------------")
	//fmt.Println(code)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(astMod.Module()))
	//fmt.Println(mod)

	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewInt(int64(10)), g.NewInt(int64(20))},
		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{&bc.FuncTemplate{
			Arity:       fixedArity(0),
			NumCaptures: 0,
			NumLocals:   4,
			Bytecodes: []byte{
				bc.LoadNull,
				bc.LoadConst, 0, 0,
				bc.StoreLocal, 0, 0,
				bc.LoadConst, 0, 1,
				bc.StoreLocal, 0, 1,
				bc.LoadLocal, 0, 0,
				bc.Dup,
				bc.LoadOne,
				bc.Inc,
				bc.StoreLocal, 0, 0,
				bc.StoreLocal, 0, 2, bc.LoadLocal, 0, 1,
				bc.Dup,
				bc.LoadNegOne,
				bc.Inc,
				bc.StoreLocal, 0, 1,
				bc.StoreLocal, 0, 3,
				bc.Return},
			LineNumberTable: []bc.LineNumberEntry{
				{Index: 0, LineNum: 0},
				{Index: 1, LineNum: 2},
				{Index: 7, LineNum: 3},
				{Index: 13, LineNum: 4},
				{Index: 25, LineNum: 5},
				{Index: 37, LineNum: 0}},
			ErrorHandlers: nil,
		}},
	})
}

func TestInvokeField(t *testing.T) {

	code := `
let s = 'abc'
let c = s.contains
let x = c('b')
let y = s.contains('z')

let ls = [1]
let p = ls.iter().next()
`
	mod := testCompile(t, code)

	//fmt.Println("----------------------------")
	//fmt.Println(code)
	//fmt.Println("----------------------------")
	//fmt.Println(mod.Pool)

	ok(t, mod.Pool, &bc.Pool{
		Constants: []g.Basic{
			g.MustStr("abc"),
			g.MustStr("contains"),
			g.MustStr("b"),
			g.MustStr("z"),
			g.MustStr("iter"),
			g.MustStr("next"),
		},

		StructDefs: [][]string{},
		Templates: []*bc.FuncTemplate{&bc.FuncTemplate{
			Arity:       fixedArity(0),
			NumCaptures: 0,
			NumLocals:   6,
			Bytecodes: []byte{
				bc.LoadNull,
				bc.LoadConst, 0, 0,
				bc.StoreLocal, 0, 0,
				bc.LoadLocal, 0, 0,
				bc.GetField, 0, 1,
				bc.StoreLocal, 0, 1,
				bc.LoadLocal, 0, 1,
				bc.LoadConst, 0, 2,
				bc.Invoke, 0, 1,
				bc.StoreLocal, 0, 2,
				bc.LoadLocal, 0, 0,
				bc.LoadConst, 0, 3,
				bc.InvokeField, 0, 1, 0, 1,
				bc.StoreLocal, 0, 3,
				bc.LoadOne,
				bc.NewList, 0, 1,
				bc.StoreLocal, 0, 4,
				bc.LoadLocal, 0, 4,
				bc.InvokeField, 0, 4, 0, 0,
				bc.InvokeField, 0, 5, 0, 0,
				bc.StoreLocal, 0, 5,
				bc.Return,
			},
			LineNumberTable: []bc.LineNumberEntry{
				{Index: 0, LineNum: 0},
				{Index: 1, LineNum: 2},
				{Index: 7, LineNum: 3},
				{Index: 16, LineNum: 4},
				{Index: 28, LineNum: 5},
				{Index: 42, LineNum: 7},
				{Index: 49, LineNum: 8},
				{Index: 65, LineNum: 0},
			},
			ErrorHandlers: nil,
		}},
	})
}

//func TestDebug(t *testing.T) {
//
//	code := `
//fn a() {
//    println('---- begin')
//    try {
//        println('---- try')
//        1/0
//    } catch e {
//        println('---- catch')
//        return 2
//    } finally {
//        println('---- finally')
//        return 3
//    }
//    println('---- end')
//}
//println(a())
////assert(a() == 3)
//`
//	lines := strings.Split(code, "\n")
//
//	mod := testCompile(t, code)
//
//	f := func(curLine int) string {
//		if curLine == 0 {
//			return "// --------"
//		}
//		return "// " + lines[curLine-1]
//	}
//
//	fmt.Println("----------------------------")
//	fmt.Println(code)
//	fmt.Println("----------------------------")
//	fmt.Println(mod.Pool.DebugString(f))
//}
