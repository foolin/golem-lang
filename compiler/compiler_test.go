// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	//"fmt"
	"reflect"
	"testing"

	"github.com/mjarmy/golem-lang/analyzer"
	g "github.com/mjarmy/golem-lang/core"
	o "github.com/mjarmy/golem-lang/core/opcodes"
	"github.com/mjarmy/golem-lang/parser"
	"github.com/mjarmy/golem-lang/scanner"
)

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
	}
}

func ok(t *testing.T, mod *g.Module, expect *g.Module) {

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

	scanner := scanner.NewScanner("", "", source)
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
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewInt(int64(-2)), g.NewInt(int64(2))},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})

	mod = newCompiler(newAnalyzer("(2 + 3) * -4 / 10")).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewInt(int64(2)), g.NewInt(int64(3)), g.NewInt(int64(-4)), g.NewInt(int64(10))},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})

	mod = newCompiler(newAnalyzer("null / true + \nfalse")).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})

	mod = newCompiler(newAnalyzer("'a' * 1.23e4")).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewStr("a"), g.NewFloat(float64(12300))},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})

	mod = newCompiler(newAnalyzer("'a' == true")).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewStr("a")},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})

	mod = newCompiler(newAnalyzer("true != false")).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})

	mod = newCompiler(newAnalyzer("true > false; true >= false")).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})

	mod = newCompiler(newAnalyzer("true < false; true <= false; true <=> false;")).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})

	mod = newCompiler(newAnalyzer("let a = 2 && 3;")).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewInt(int64(2)), g.NewInt(int64(3))},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})

	mod = newCompiler(newAnalyzer("let a = 2 || 3;")).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewInt(int64(2)), g.NewInt(int64(3))},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})
}

func TestAssignment(t *testing.T) {

	mod := newCompiler(newAnalyzer("let a = 1;\nconst b = \n2;a = 3;")).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewInt(2), g.NewInt(3)},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
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

	source := "if (3 == 2) { let a = 42; }"
	anl := newAnalyzer(source)
	mod := newCompiler(anl).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewInt(3), g.NewInt(2), g.NewInt(42)},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})

	source = `let a = 1
		if (false) {
		    let b = 2
		} else {
		    let c = 3
		}
		let d = 4`

	anl = newAnalyzer(source)
	mod = newCompiler(anl).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewInt(2), g.NewInt(3), g.NewInt(4)},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})
}

func TestWhile(t *testing.T) {

	source := "let a = 1; while (0 < 1) { let b = 2; }"
	mod := newCompiler(newAnalyzer(source)).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewInt(2)},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})

	source = "let a = 'z'; while (0 < 1) \n{ break; continue; let b = 2; }; let c = 3;"
	mod = newCompiler(newAnalyzer(source)).Compile()
	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewStr("z"), g.NewInt(2), g.NewInt(3)},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})
}

func TestReturn(t *testing.T) {

	source := "let a = 1; return a \n- 2; a = 3;"
	anl := newAnalyzer(source)
	mod := newCompiler(anl).Compile()

	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewInt(2), g.NewInt(3)},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
		Contents: contents(),
	})
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
	//fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//fmt.Println(mod)

	ok(t, mod, &g.Module{
		Pool: []g.Basic{
			g.NewInt(42),
			g.NewInt(7)},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
			&g.Template{
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
			&g.Template{
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
			&g.Template{
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
		Contents: contents()})

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
	//fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//fmt.Println(mod)

	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewInt(2), g.NewInt(3), g.NewInt(4)},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{
			&g.Template{
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
			&g.Template{
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
			&g.Template{
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
			&g.Template{
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
		Contents: contents(),
	})
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

	ok(t, mod, &g.Module{
		Pool:       []g.Basic{},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{&g.Template{
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
		}, &g.Template{
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
		}, &g.Template{
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
		Contents: contents(),
	})

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
	//fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//fmt.Println(mod)

	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewInt(2)},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{&g.Template{
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
		}, &g.Template{
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
		}, &g.Template{
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
		Contents: contents(),
	})
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
	//fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//fmt.Println(mod)

	ok(t, mod, &g.Module{
		Pool:       []g.Basic{g.NewInt(int64(10)), g.NewInt(int64(20))},
		Refs:       nil,
		StructDefs: [][]*g.FieldDef{},
		Templates: []*g.Template{&g.Template{
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
		Contents: contents(),
	})
}

func TestPool(t *testing.T) {
	pool := g.EmptyHashMap()

	tassert(t, poolIndex(pool, g.NewInt(4)) == 0)
	tassert(t, poolIndex(pool, g.NewStr("a")) == 1)
	tassert(t, poolIndex(pool, g.NewFloat(1.0)) == 2)
	tassert(t, poolIndex(pool, g.NewStr("a")) == 1)
	tassert(t, poolIndex(pool, g.NewInt(4)) == 0)

	slice := makePoolSlice(pool)
	tassert(t, reflect.DeepEqual(
		slice,
		[]g.Basic{
			g.NewInt(4),
			g.NewStr("a"),
			g.NewFloat(1.0)}))
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
		g.ExceptionHandler{
			Begin:   5,
			End:     14,
			Catch:   -1,
			Finally: 14,
		})
}
