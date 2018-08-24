// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this code code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"reflect"
	"testing"

	g "github.com/mjarmy/golem-lang/core"
	o "github.com/mjarmy/golem-lang/core/opcodes"
	"github.com/mjarmy/golem-lang/scanner"
)

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
	}
}

func ok(t *testing.T, pool *g.Pool, expect *g.Pool) {

	if !reflect.DeepEqual(pool.Constants, expect.Constants) {
		t.Error(pool, " != ", expect)
	}

	if len(pool.Templates) != len(expect.Templates) {
		t.Error(pool.Templates, " != ", expect.Templates)
	}

	for i := 0; i < len(pool.Templates); i++ {

		mt := pool.Templates[i]
		et := expect.Templates[i]

		if (mt.Arity != et.Arity) || (mt.NumCaptures != et.NumCaptures) || (mt.NumLocals != et.NumLocals) {
			t.Error(pool, " != ", expect)
		}

		if !reflect.DeepEqual(mt.OpCodes, et.OpCodes) {
			t.Error("OpCodes: ", pool, " != ", expect)
		}

		// checking LineNumberTable is optional
		if et.LineNumberTable != nil {
			if !reflect.DeepEqual(mt.LineNumberTable, et.LineNumberTable) {
				t.Error("LineNumberTable: ", pool, " != ", expect)
			}
		}
	}
}

var builtins []*g.BuiltinEntry = append(
	g.SandboxBuiltins,
	g.CommandLineBuiltins...)
var builtinMgr = g.NewBuiltinManager(builtins)

func testCompile(t *testing.T, code string) *g.Module {

	source := &scanner.Source{Name: "foo", Path: "foo.glm", Code: code}
	mods, errs := CompileSourceFully(builtinMgr, source, nil)
	tassert(t, errs == nil)
	tassert(t, len(mods) == 1)

	fmt.Printf("%v\n", mods)

	return mods[0]
}

func TestExpression(t *testing.T) {

	mod := testCompile(t, "-2 + -1 + -0 + 0 + 1 + 2")
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewInt(int64(-2)), g.NewInt(int64(2))},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   0,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadConst, 0, 0,
					o.LoadNegOne,
					o.Plus,
					o.LoadZero,
					o.Plus,
					o.LoadZero,
					o.Plus,
					o.LoadOne,
					o.Plus,
					o.LoadConst, 0, 1,
					o.Plus,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 16, LineNum: 0}},
				ExceptionHandlers: nil,
			}},
	})

	mod = testCompile(t, "(2 + 3) * -4 / 10")
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewInt(int64(2)), g.NewInt(int64(3)), g.NewInt(int64(-4)), g.NewInt(int64(10))},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   0,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadConst, 0, 0,
					o.LoadConst, 0, 1,
					o.Plus,
					o.LoadConst, 0, 2,
					o.Mul,
					o.LoadConst, 0, 3,
					o.Div,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 16, LineNum: 0}},
				ExceptionHandlers: nil,
			}},
	})

	mod = testCompile(t, "null / true + \nfalse")
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   0,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadNull,
					o.LoadTrue,
					o.Div,
					o.LoadFalse,
					o.Plus,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 4, LineNum: 2},
					{Index: 5, LineNum: 1},
					{Index: 6, LineNum: 0}},
				ExceptionHandlers: nil,
			}},
	})

	mod = testCompile(t, "'a' * 1.23e4")
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewStr("a"), g.NewFloat(float64(12300))},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   0,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadConst, 0, 0,
					o.LoadConst, 0, 1,
					o.Mul,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 8, LineNum: 0}},
				ExceptionHandlers: nil,
			}},
	})

	mod = testCompile(t, "'a' == true")
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewStr("a")},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   0,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadConst, 0, 0,
					o.LoadTrue,
					o.Eq,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 6, LineNum: 0}},
				ExceptionHandlers: nil,
			}},
	})

	mod = testCompile(t, "true != false")
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   0,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadTrue,
					o.LoadFalse,
					o.Ne,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 4, LineNum: 0}},
				ExceptionHandlers: nil,
			}},
	})

	mod = testCompile(t, "true > false; true >= false")
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   0,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadTrue,
					o.LoadFalse,
					o.Gt,
					o.LoadTrue,
					o.LoadFalse,
					o.Gte,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 7, LineNum: 0}},
				ExceptionHandlers: nil,
			}},
	})

	mod = testCompile(t, "true < false; true <= false; true <=> false;")
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   0,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadTrue,
					o.LoadFalse,
					o.Lt,
					o.LoadTrue,
					o.LoadFalse,
					o.Lte,
					o.LoadTrue,
					o.LoadFalse,
					o.Cmp,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 10, LineNum: 0}},
				ExceptionHandlers: nil,
			}},
	})

	mod = testCompile(t, "let a = 2 && 3;")
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewInt(int64(2)), g.NewInt(int64(3))},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   1,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadConst, 0, 0,
					o.JumpFalse, 0, 17,
					o.LoadConst, 0, 1,
					o.JumpFalse, 0, 17,
					o.LoadTrue,
					o.Jump, 0, 18,
					o.LoadFalse,
					o.StoreLocal, 0, 0,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 21, LineNum: 0}},
				ExceptionHandlers: nil,
			}},
	})

	mod = testCompile(t, "let a = 2 || 3;")
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewInt(int64(2)), g.NewInt(int64(3))},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   1,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadConst, 0, 0,
					o.JumpTrue, 0, 13,
					o.LoadConst, 0, 1,
					o.JumpFalse, 0, 17,
					o.LoadTrue,
					o.Jump, 0, 18,
					o.LoadFalse,
					o.StoreLocal, 0, 0,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 21, LineNum: 0}},
				ExceptionHandlers: nil,
			}},
	})
}

func TestAssignment(t *testing.T) {

	mod := testCompile(t, "let a = 1;\nconst b = \n2;a = 3;")
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewInt(2), g.NewInt(3)},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   2,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadOne,
					o.StoreLocal, 0, 0,
					o.LoadConst, 0, 0,
					o.StoreLocal, 0, 1,
					o.LoadConst, 0, 1,
					o.Dup,
					o.StoreLocal, 0, 0,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
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
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewInt(3), g.NewInt(2), g.NewInt(42)},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   1,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadConst, 0, 0,
					o.LoadConst, 0, 1,
					o.Eq,
					o.JumpFalse, 0, 17,
					o.LoadConst, 0, 2,
					o.StoreLocal, 0, 0,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
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
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewInt(2), g.NewInt(3), g.NewInt(4)},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   4,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadOne,
					o.StoreLocal, 0, 0,
					o.LoadFalse,
					o.JumpFalse, 0, 18,
					o.LoadConst, 0, 0,
					o.StoreLocal, 0, 1,
					o.Jump, 0, 24,
					o.LoadConst, 0, 1,
					o.StoreLocal, 0, 2,
					o.LoadConst, 0, 2,
					o.StoreLocal, 0, 3,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
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
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewInt(2)},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   2,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadOne,
					o.StoreLocal, 0, 0,
					o.LoadZero,
					o.LoadOne,
					o.Lt,
					o.JumpFalse, 0, 20,
					o.LoadConst, 0, 0,
					o.StoreLocal, 0, 1,
					o.Jump, 0, 5,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 1},
					{Index: 20, LineNum: 0}},
				ExceptionHandlers: nil,
			}},
	})

	code = "let a = 'z'; while (0 < 1) \n{ break; continue; let b = 2; }; let c = 3;"
	mod = testCompile(t, code)
	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewStr("z"), g.NewInt(2), g.NewInt(3)},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   3,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadConst, 0, 0,
					o.StoreLocal, 0, 0,
					o.LoadZero,
					o.LoadOne,
					o.Lt,
					o.JumpFalse, 0, 28,
					o.Jump, 0, 28,
					o.Jump, 0, 7,
					o.LoadConst, 0, 1,
					o.StoreLocal, 0, 1,
					o.Jump, 0, 7,
					o.LoadConst, 0, 2,
					o.StoreLocal, 0, 2,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
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

	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewInt(2), g.NewInt(3)},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   1,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadOne,
					o.StoreLocal, 0, 0,
					o.LoadLocal, 0, 0,
					o.LoadConst, 0, 0,
					o.Sub,
					o.ReturnStmt,
					o.LoadConst, 0, 1,
					o.Dup,
					o.StoreLocal, 0, 0,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
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

	ok(t, mod.Pool, &g.Pool{
		Constants: []g.Basic{
			g.NewInt(42),
			g.NewInt(7)},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   2,
				OpCodes: []byte{
					o.LoadNull,
					o.NewFunc, 0, 1,
					o.StoreLocal, 0, 0,
					o.NewFunc, 0, 2,
					o.StoreLocal, 0, 1,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 2},
					{Index: 7, LineNum: 3},
					{Index: 13, LineNum: 0}},
				ExceptionHandlers: nil,
			},
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   0,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadConst, 0, 0,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 2},
					{Index: 4, LineNum: 0}},
				ExceptionHandlers: nil,
			},
			&g.FuncTemplate{
				Arity:       1,
				NumCaptures: 0,
				NumLocals:   2,
				OpCodes: []byte{
					o.LoadNull,
					o.NewFunc, 0, 3,
					o.StoreLocal, 0, 1,
					o.LoadLocal, 0, 0,
					o.LoadLocal, 0, 0,
					o.Mul,
					o.LoadLocal, 0, 1,
					o.LoadLocal, 0, 0,
					o.Invoke, 0, 1,
					o.Plus,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 4},
					{Index: 7, LineNum: 7},
					{Index: 24, LineNum: 0}},
				ExceptionHandlers: nil,
			},
			&g.FuncTemplate{
				Arity:       1,
				NumCaptures: 0,
				NumLocals:   1,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadLocal, 0, 0,
					o.LoadConst, 0, 1,
					o.Mul,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
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

	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewInt(2), g.NewInt(3), g.NewInt(4)},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   3,
				OpCodes: []byte{
					o.LoadNull,
					o.NewFunc, 0, 1,
					o.StoreLocal, 0, 0,
					o.NewFunc, 0, 2,
					o.StoreLocal, 0, 1,
					o.NewFunc, 0, 3,
					o.StoreLocal, 0, 2,
					o.LoadLocal, 0, 0,
					o.Invoke, 0, 0,
					o.LoadLocal, 0, 1,
					o.LoadOne,
					o.Invoke, 0, 1,
					o.LoadLocal, 0, 2,
					o.LoadConst, 0, 0,
					o.LoadConst, 0, 1,
					o.Invoke, 0, 2,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
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
			&g.FuncTemplate{
				Arity:       0,
				NumCaptures: 0,
				NumLocals:   0,
				OpCodes: []byte{
					o.LoadNull,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0}},
				ExceptionHandlers: nil,
			},
			&g.FuncTemplate{
				Arity:       1,
				NumCaptures: 0,
				NumLocals:   1,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadLocal, 0, 0,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
					{Index: 0, LineNum: 0},
					{Index: 1, LineNum: 3},
					{Index: 4, LineNum: 0}},
				ExceptionHandlers: nil,
			},
			&g.FuncTemplate{
				Arity:       2,
				NumCaptures: 0,
				NumLocals:   3,
				OpCodes: []byte{
					o.LoadNull,
					o.LoadConst, 0, 2,
					o.StoreLocal, 0, 2,
					o.LoadLocal, 0, 0,
					o.LoadLocal, 0, 1,
					o.Mul,
					o.LoadLocal, 0, 2,
					o.Mul,
					o.ReturnStmt},
				LineNumberTable: []g.LineNumberEntry{
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

	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{&g.FuncTemplate{
			Arity:       0,
			NumCaptures: 0,
			NumLocals:   1,
			OpCodes: []byte{
				o.LoadNull,
				o.NewFunc, 0, 1,
				o.StoreLocal, 0, 0,
				o.ReturnStmt},
			LineNumberTable: []g.LineNumberEntry{
				{Index: 0, LineNum: 0},
				{Index: 1, LineNum: 2},
				{Index: 7, LineNum: 0}},
			ExceptionHandlers: nil,
		}, &g.FuncTemplate{
			Arity:       1,
			NumCaptures: 0,
			NumLocals:   1,
			OpCodes: []byte{
				o.LoadNull,
				o.NewFunc, 0, 2,
				o.FuncLocal, 0, 0,
				o.ReturnStmt,
				o.ReturnStmt},
			LineNumberTable: []g.LineNumberEntry{
				{Index: 0, LineNum: 0},
				{Index: 1, LineNum: 3},
				{Index: 8, LineNum: 0}},
			ExceptionHandlers: nil,
		}, &g.FuncTemplate{
			Arity:       1,
			NumCaptures: 1,
			NumLocals:   1,
			OpCodes: []byte{
				o.LoadNull,
				o.LoadCapture, 0, 0,
				o.LoadLocal, 0, 0,
				o.Plus,
				o.Dup,
				o.StoreCapture, 0, 0,
				o.LoadCapture, 0, 0,
				o.ReturnStmt,
				o.ReturnStmt},
			LineNumberTable: []g.LineNumberEntry{
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

	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewInt(2)},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{&g.FuncTemplate{
			Arity:       0,
			NumCaptures: 0,
			NumLocals:   2,
			OpCodes: []byte{
				o.LoadNull,
				o.LoadConst, 0, 0,
				o.StoreLocal, 0, 0,
				o.NewFunc, 0, 1,
				o.FuncLocal, 0, 0,
				o.StoreLocal, 0, 1,
				o.ReturnStmt},
			LineNumberTable: []g.LineNumberEntry{
				{Index: 0, LineNum: 0},
				{Index: 1, LineNum: 2},
				{Index: 7, LineNum: 3},
				{Index: 16, LineNum: 0}},
			ExceptionHandlers: nil,
		}, &g.FuncTemplate{
			Arity:       1,
			NumCaptures: 1,
			NumLocals:   1,
			OpCodes: []byte{
				o.LoadNull,
				o.NewFunc, 0, 2,
				o.FuncLocal, 0, 0,
				o.FuncCapture, 0, 0,
				o.ReturnStmt,
				o.ReturnStmt},
			LineNumberTable: []g.LineNumberEntry{
				{Index: 0, LineNum: 0},
				{Index: 1, LineNum: 4},
				{Index: 11, LineNum: 0}},
			ExceptionHandlers: nil,
		}, &g.FuncTemplate{
			Arity:       1,
			NumCaptures: 2,
			NumLocals:   1,
			OpCodes: []byte{
				o.LoadNull,
				o.LoadCapture, 0, 0,
				o.LoadLocal, 0, 0,
				o.Plus,
				o.LoadCapture, 0, 1,
				o.Plus,
				o.Dup,
				o.StoreCapture, 0, 0,
				o.LoadCapture, 0, 0,
				o.ReturnStmt,
				o.ReturnStmt},
			LineNumberTable: []g.LineNumberEntry{
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

	ok(t, mod.Pool, &g.Pool{
		Constants:  []g.Basic{g.NewInt(int64(10)), g.NewInt(int64(20))},
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.FuncTemplate{&g.FuncTemplate{
			Arity:       0,
			NumCaptures: 0,
			NumLocals:   4,
			OpCodes: []byte{
				o.LoadNull,
				o.LoadConst, 0, 0,
				o.StoreLocal, 0, 0,
				o.LoadConst, 0, 1,
				o.StoreLocal, 0, 1,
				o.LoadLocal, 0, 0,
				o.Dup,
				o.LoadOne,
				o.Plus,
				o.StoreLocal, 0, 0,
				o.StoreLocal, 0, 2,
				o.LoadLocal, 0, 1,
				o.Dup,
				o.LoadNegOne,
				o.Plus,
				o.StoreLocal, 0, 1,
				o.StoreLocal, 0, 3,
				o.ReturnStmt},
			LineNumberTable: []g.LineNumberEntry{
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

func TestTry(t *testing.T) {

	code := `
let a = 1
try {
    a++
}
finally {
    a++
}
assert(a == 2)
`
	mod := testCompile(t, code)
	tassert(t, mod.Pool.Templates[0].ExceptionHandlers[0] ==
		g.ExceptionHandler{
			Begin:   5,
			End:     14,
			Catch:   -1,
			Finally: 14,
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
		return nil, g.UndefinedModuleError(moduleName)
	}

	mods, errs := CompileSourceFully(builtinMgr, srcMain, resolver)
	tassert(t, errs == nil)
	tassert(t, len(mods) == 4)
	tassert(t, mods[0].Name == "foo")
	tassert(t, mods[1].Name == "a")
	tassert(t, mods[2].Name == "b")
	tassert(t, mods[3].Name == "c")
}
