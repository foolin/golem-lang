// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"fmt"
)

// Define all the various bytecodes
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

	Plus
	Inc
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
	Go
	Return
	Done
	Throw

	NewStruct
	NewDict
	NewList
	NewSet
	NewTuple

	GetField
	InvokeField
	ReplaceField
	SetField
	IncField

	GetIndex
	SetIndex
	IncIndex

	Slice
	SliceFrom
	SliceTo

	NewIter
	IterNext
	IterGet

	//CheckCast
	CheckTuple

	Pop
	Dup

	// These are temporary values created during compilation.
	// The interpreter will panic if it encounters them.
	Break    = 0xFE
	Continue = 0xFF
)

func BytecodeString(bc byte) string {

	switch bc {

	case LoadNull:
		return "LoadNull"
	case LoadTrue:
		return "LoadTrue"
	case LoadFalse:
		return "LoadFalse"
	case LoadZero:
		return "LoadZero"
	case LoadOne:
		return "LoadOne"
	case LoadNegOne:
		return "LoadNegOne"

	case ImportModule:
		return "ImportModule"
	case LoadBuiltin:
		return "LoadBuiltin"
	case LoadConst:
		return "LoadConst"
	case LoadLocal:
		return "LoadLocal"
	case StoreLocal:
		return "StoreLocal"
	case LoadCapture:
		return "LoadCapture"
	case StoreCapture:
		return "StoreCapture"

	case Jump:
		return "Jump"
	case JumpTrue:
		return "JumpTrue"
	case JumpFalse:
		return "JumpFalse"

	case Eq:
		return "Eq"
	case Ne:
		return "Ne"
	case Gt:
		return "Gt"
	case Gte:
		return "Gte"
	case Lt:
		return "Lt"
	case Lte:
		return "Lte"
	case Cmp:
		return "Cmp"

	case Plus:
		return "Plus"
	case Inc:
		return "Inc"
	case Sub:
		return "Sub"
	case Mul:
		return "Mul"
	case Div:
		return "Div"

	case Rem:
		return "Rem"
	case BitAnd:
		return "BitAnd"
	case BitOr:
		return "BitOr"
	case BitXor:
		return "BitXor"
	case LeftShift:
		return "LeftShift"
	case RightShift:
		return "RightShift"

	case Negate:
		return "Negate"
	case Not:
		return "Not"
	case Complement:
		return "Complement"

	case NewFunc:
		return "NewFunc"
	case FuncCapture:
		return "FuncCapture"
	case FuncLocal:
		return "FuncLocal"

	case Invoke:
		return "Invoke"
	case Go:
		return "Go"
	case Return:
		return "Return"
	case Done:
		return "Done"
	case Throw:
		return "Throw"

	case NewStruct:
		return "NewStruct"
	case ReplaceField:
		return "ReplaceField"

	case NewDict:
		return "NewDict"
	case NewList:
		return "NewList"
	case NewSet:
		return "NewSet"
	case NewTuple:
		return "NewTuple"

	case GetField:
		return "GetField"
	case InvokeField:
		return "InvokeField"
	case SetField:
		return "SetField"
	case IncField:
		return "IncField"

	case GetIndex:
		return "GetIndex"
	case SetIndex:
		return "SetIndex"
	case IncIndex:
		return "IncIndex"

	case Slice:
		return "Slice"
	case SliceFrom:
		return "SliceFrom"
	case SliceTo:
		return "SliceTo"

	case NewIter:
		return "NewIter"
	case IterNext:
		return "IterNext"
	case IterGet:
		return "IterGet"

	//case CheckCast:
	//	return "CheckCast"
	case CheckTuple:
		return "CheckTuple"

	case Pop:
		return "Pop"
	case Dup:
		return "Dup"

	case Break:
		return "Break"
	case Continue:
		return "Continue"

	default:
		panic("unreachable")
	}
}

// BytecodeSize returns how 'wide' a bytecode is.  Bytecodes are always either 1, 3, or 5 bytes long.
func BytecodeSize(bc byte) int {

	switch bc {

	case
		LoadNull, LoadTrue, LoadFalse, LoadZero, LoadOne, LoadNegOne,
		Eq, Ne, Gt, Gte, Lt, Lte, Cmp,
		Plus, Inc, Sub, Mul, Div,
		Rem, BitAnd, BitOr, BitXor, LeftShift, RightShift,
		Negate, Not, Complement,
		Return, Done, Throw,
		GetIndex, SetIndex, IncIndex, Slice, SliceFrom, SliceTo,
		NewIter, IterNext, IterGet, Pop, Dup:

		return 1

	case
		ImportModule, LoadBuiltin, LoadConst,
		LoadLocal, LoadCapture, StoreLocal, StoreCapture,
		Jump, JumpTrue, JumpFalse, Break, Continue,
		NewFunc, FuncCapture, FuncLocal, Invoke, Go,
		NewStruct, GetField, ReplaceField, SetField, IncField,
		NewDict, NewList, NewSet, NewTuple /*CheckCast,*/, CheckTuple:

		return 3

	case InvokeField:

		return 5

	default:
		panic("unreachable")
	}
}

func EncodeParam(p int) (byte, byte) {

	// TODO implicit wide indexing
	if p >= (2 << 16) {
		panic("Internal Compiler Error")
	}

	return byte((p >> 8) & 0xFF), byte(p & 0xFF)
}

func EncodeWideParams(p, q int) (byte, byte, byte, byte) {

	a, b := EncodeParam(p)
	c, d := EncodeParam(q)
	return a, b, c, d
}

func DecodeParam(btc []byte, ip int) int {

	high := btc[ip+1]
	low := btc[ip+2]

	return int(high)<<8 + int(low)
}

func DecodeWideParams(btc []byte, ip int) (int, int) {

	return DecodeParam(btc, ip), DecodeParam(btc, ip+2)
}

func FmtBytecode(btc []byte, ip int) string {

	str := BytecodeString(btc[ip])

	switch BytecodeSize(btc[ip]) {

	case 1:
		return fmt.Sprintf("%d: %s", ip, str)

	case 3:

		p := DecodeParam(btc, ip)
		return fmt.Sprintf(
			"%d: %s %d %d (%d)",
			ip, padRight(str, 12, " "),
			btc[ip+1], btc[ip+2], p)

	case 5:

		p, q := DecodeWideParams(btc, ip)
		return fmt.Sprintf(
			"%d: %s %d %d (%d), %d %d (%d)",
			ip, padRight(str, 12, " "),
			btc[ip+1], btc[ip+2], p,
			btc[ip+3], btc[ip+4], q)

	default:
		panic("unreachable")
	}
}

func padRight(str string, length int, pad string) string {
	for len(str) < length {
		str = str + pad
	}
	return str
}
