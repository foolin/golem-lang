// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"fmt"
)

// Define all the various opcodes
const (
	LoadNull byte = iota
	LoadTrue
	LoadFalse
	LoadZero
	LoadOne
	LoadNegOne

	ImportModule
	LoadBuiltin
	LoadConst
	LoadLocal
	StoreLocal
	LoadCapture
	StoreCapture

	Jump
	JumpTrue
	JumpFalse

	Eq
	Ne
	Gt
	Gte
	Lt
	Lte
	Cmp
	Has

	Plus
	Sub
	Mul
	Div

	Rem
	BitAnd
	BitOr
	BitXor
	LeftShift
	RightShift

	Negate
	Not
	Complement

	NewFunc
	FuncCapture
	FuncLocal

	Invoke
	GoStmt
	ReturnStmt
	Done
	ThrowStmt

	DefineStruct
	NewDict
	NewList
	NewSet
	NewTuple

	GetField
	InitField
	SetField
	IncField

	GetIndex
	SetIndex
	IncIndex

	Slice
	SliceFrom
	SliceTo

	Iter
	IterNext
	IterGet

	CheckCast
	CheckTuple

	Pop
	Dup

	// These are temporary values created during compilation.
	// The interpreter will panic if it encounters them.
	BreakStmt    = 0xFD
	ContinueStmt = 0xFE
)

// OpCodeSize returns how 'wide' an opcode is.  Opcodes are always either 1 or 3 bytes.
func OpCodeSize(opc byte) int {

	switch opc {

	case ImportModule, LoadBuiltin, LoadConst,
		LoadLocal, LoadCapture, StoreLocal, StoreCapture,
		Jump, JumpTrue, JumpFalse, BreakStmt, ContinueStmt,
		NewFunc, FuncCapture, FuncLocal, Invoke, GoStmt,
		DefineStruct, GetField, InitField, SetField, IncField,
		NewDict, NewList, NewSet, NewTuple, CheckCast, CheckTuple:

		return 3

	default:
		return 1
	}
}

// FmtOpcode formats an opcode into a string
func FmtOpcode(opcodes []byte, i int) string {

	switch opcodes[i] {

	case LoadNull:
		return fmt.Sprintf("%d: LoadNull\n", i)
	case LoadTrue:
		return fmt.Sprintf("%d: LoadTrue\n", i)
	case LoadFalse:
		return fmt.Sprintf("%d: LoadFalse\n", i)
	case LoadZero:
		return fmt.Sprintf("%d: LoadZero\n", i)
	case LoadOne:
		return fmt.Sprintf("%d: LoadOne\n", i)
	case LoadNegOne:
		return fmt.Sprintf("%d: LoadNegOne\n", i)

	case ImportModule:
		return fmtIndex(opcodes, i, "ImportModule")
	case LoadBuiltin:
		return fmtIndex(opcodes, i, "LoadBuiltin")
	case LoadConst:
		return fmtIndex(opcodes, i, "LoadConst")
	case LoadLocal:
		return fmtIndex(opcodes, i, "LoadLocal")
	case StoreLocal:
		return fmtIndex(opcodes, i, "StoreLocal")
	case LoadCapture:
		return fmtIndex(opcodes, i, "LoadCapture")
	case StoreCapture:
		return fmtIndex(opcodes, i, "StoreCapture")

	case Jump:
		return fmtIndex(opcodes, i, "Jump")
	case JumpTrue:
		return fmtIndex(opcodes, i, "JumpTrue")
	case JumpFalse:
		return fmtIndex(opcodes, i, "JumpFalse")

	case Eq:
		return fmt.Sprintf("%d: Eq\n", i)
	case Ne:
		return fmt.Sprintf("%d: Ne\n", i)
	case Gt:
		return fmt.Sprintf("%d: Gt\n", i)
	case Gte:
		return fmt.Sprintf("%d: Gte\n", i)
	case Lt:
		return fmt.Sprintf("%d: Lt\n", i)
	case Lte:
		return fmt.Sprintf("%d: Lte\n", i)
	case Cmp:
		return fmt.Sprintf("%d: Cmp\n", i)
	case Has:
		return fmt.Sprintf("%d: Has\n", i)

	case Plus:
		return fmt.Sprintf("%d: Plus\n", i)
	case Sub:
		return fmt.Sprintf("%d: Sub\n", i)
	case Mul:
		return fmt.Sprintf("%d: Mul\n", i)
	case Div:
		return fmt.Sprintf("%d: Div\n", i)

	case Rem:
		return fmt.Sprintf("%d: Rem\n", i)
	case BitAnd:
		return fmt.Sprintf("%d: BitAnd\n", i)
	case BitOr:
		return fmt.Sprintf("%d: BitOr\n", i)
	case BitXor:
		return fmt.Sprintf("%d: BitXor\n", i)
	case LeftShift:
		return fmt.Sprintf("%d: LeftShift\n", i)
	case RightShift:
		return fmt.Sprintf("%d: RightShift\n", i)

	case Negate:
		return fmt.Sprintf("%d: Negate\n", i)
	case Not:
		return fmt.Sprintf("%d: Not\n", i)
	case Complement:
		return fmt.Sprintf("%d: Complement\n", i)

	case NewFunc:
		return fmtIndex(opcodes, i, "NewFunc")
	case FuncCapture:
		return fmtIndex(opcodes, i, "FuncCapture")
	case FuncLocal:
		return fmtIndex(opcodes, i, "FuncLocal")

	case Invoke:
		return fmtIndex(opcodes, i, "Invoke")
	case GoStmt:
		return fmtIndex(opcodes, i, "GoStmt")
	case ReturnStmt:
		return fmt.Sprintf("%d: ReturnStmt\n", i)
	case Done:
		return fmt.Sprintf("%d: Done\n", i)
	case ThrowStmt:
		return fmt.Sprintf("%d: ThrowStmt\n", i)

	case DefineStruct:
		return fmtIndex(opcodes, i, "DefineStruct")
	case GetField:
		return fmtIndex(opcodes, i, "GetField")
	case InitField:
		return fmtIndex(opcodes, i, "InitField")
	case SetField:
		return fmtIndex(opcodes, i, "SetField")
	case IncField:
		return fmtIndex(opcodes, i, "IncField")
	case NewDict:
		return fmtIndex(opcodes, i, "NewDict")
	case NewList:
		return fmtIndex(opcodes, i, "NewList")
	case NewSet:
		return fmtIndex(opcodes, i, "NewSet")
	case NewTuple:
		return fmtIndex(opcodes, i, "NewTuple")

	case GetIndex:
		return fmt.Sprintf("%d: GetIndex\n", i)
	case SetIndex:
		return fmt.Sprintf("%d: SetIndex\n", i)
	case IncIndex:
		return fmt.Sprintf("%d: IncIndex\n", i)

	case Slice:
		return fmt.Sprintf("%d: Slice\n", i)
	case SliceFrom:
		return fmt.Sprintf("%d: SliceFrom\n", i)
	case SliceTo:
		return fmt.Sprintf("%d: SliceTo\n", i)

	case Iter:
		return fmt.Sprintf("%d: Iter\n", i)
	case IterNext:
		return fmt.Sprintf("%d: IterNext\n", i)
	case IterGet:
		return fmt.Sprintf("%d: IterGet\n", i)

	case CheckCast:
		return fmtIndex(opcodes, i, "CheckCast")
	case CheckTuple:
		return fmtIndex(opcodes, i, "CheckTuple")

	case Pop:
		return fmt.Sprintf("%d: Pop\n", i)
	case Dup:
		return fmt.Sprintf("%d: Dup\n", i)

	case BreakStmt:
		return fmtIndex(opcodes, i, "BreakStmt")
	case ContinueStmt:
		return fmtIndex(opcodes, i, "ContinueStmt")

	default:
		panic(fmt.Sprintf("unreachable %d", opcodes[i]))
	}
}

func fmtIndex(opcodes []byte, i int, tag string) string {
	high := opcodes[i+1]
	low := opcodes[i+2]
	index := int(high)<<8 + int(low)
	return fmt.Sprintf("%d: %s %d %d (%d)\n", i, tag, high, low, index)
}
