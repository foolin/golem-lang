// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this code code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"reflect"
	"testing"

	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
	"github.com/mjarmy/golem-lang/scanner"
)

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

var builtins []*g.BuiltinEntry = nil
var builtinMgr = g.NewBuiltinManager(builtins)

func testCompile(t *testing.T, code string) *bc.Module {

	source := &scanner.Source{Name: "foo", Path: "foo.glm", Code: code}
	mods, errs := CompileSourceFully(builtinMgr, source, nil)
	g.Tassert(t, errs == nil)
	g.Tassert(t, len(mods) == 1)

	return mods[0]
}

func fixedArity(numParams int) g.Arity {
	return g.Arity{
		Kind:           g.FixedArity,
		RequiredParams: uint16(numParams),
		OptionalParams: 0,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
			}},
	})

	mod = testCompile(t, "'a' * 1.23e4")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewStr("a"), g.NewFloat(float64(12300))},
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
				ExceptionHandlers: nil,
			}},
	})

	mod = testCompile(t, "'a' == true")
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewStr("a")},
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
			}},
	})

	code = "let a = 'z'; while (0 < 1) \n{ break; continue; let b = 2; }; let c = 3;"
	mod = testCompile(t, code)
	ok(t, mod.Pool, &bc.Pool{
		Constants:  []g.Basic{g.NewStr("z"), g.NewInt(2), g.NewInt(3)},
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
				ExceptionHandlers: nil,
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
			ExceptionHandlers: nil,
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
			ExceptionHandlers: nil,
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
			ExceptionHandlers: nil,
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
			ExceptionHandlers: nil,
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
			ExceptionHandlers: nil,
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
			ExceptionHandlers: nil,
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
			ExceptionHandlers: nil,
		}},
	})
}

func TestImport(t *testing.T) {

	srcMain := &scanner.Source{Name: "foo", Path: "foo.glm", Code: "import a, b;"}
	sourceMap := map[string]*scanner.Source{
		"a": &scanner.Source{Name: "a", Path: "a.glm", Code: "import c;"},
		"b": &scanner.Source{Name: "b", Path: "b.glm", Code: "import c;"},
		"c": &scanner.Source{Name: "c", Path: "c.glm", Code: ""},
	}
	resolver := func(moduleName string) (*scanner.Source, error) {
		if src, ok := sourceMap[moduleName]; ok {
			return src, nil
		}
		panic("unreachable")
	}

	mods, errs := CompileSourceFully(builtinMgr, srcMain, resolver)
	g.Tassert(t, errs == nil)
	g.Tassert(t, len(mods) == 4)
	g.Tassert(t, mods[0].Name == "foo")
	g.Tassert(t, mods[1].Name == "a")
	g.Tassert(t, mods[2].Name == "b")
	g.Tassert(t, mods[3].Name == "c")
}

func TestTry(t *testing.T) {

	code := `
let a = 1
try {
    a++
}
finally {
    a++
}
`
	mod := testCompile(t, code)

	g.Tassert(t, mod.Pool.Templates[0].ExceptionHandlers[0] ==
		bc.ExceptionHandler{
			Begin:   5,
			End:     14,
			Catch:   -1,
			Finally: 14,
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
			g.NewStr("abc"),
			g.NewStr("contains"),
			g.NewStr("b"),
			g.NewStr("z"),
			g.NewStr("iter"),
			g.NewStr("next"),
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
			ExceptionHandlers: nil,
		}},
	})
}

func TestStruct(t *testing.T) {

	code := `
let x = struct {}
let y = struct {a: 1}
let z = struct {a: 1, b: 2}
`
	mod := testCompile(t, code)

	fmt.Println("----------------------------")
	fmt.Println(code)
	fmt.Println("----------------------------")
	fmt.Println(mod.Pool)
}
