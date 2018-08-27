// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

type (
	// BuiltinManager manages the built-in functions (and other built-in values)
	// for a given instance of the Interpreter
	BuiltinManager interface {
		Builtins() []Value
		Contains(s string) bool
		IndexOf(s string) int
	}

	// BuiltinEntry is an entry in a BuiltinManager
	BuiltinEntry struct {
		Name  string
		Value Value
	}

	builtinManager struct {
		values []Value
		lookup map[string]int
	}
)

// NewBuiltinManager creates a new BuiltinManager
func NewBuiltinManager(entries []*BuiltinEntry) BuiltinManager {
	values := make([]Value, len(entries))
	lookup := make(map[string]int)
	for i, e := range entries {
		values[i] = e.Value
		lookup[e.Name] = i
	}
	return &builtinManager{values, lookup}
}

func (b *builtinManager) Builtins() []Value {
	return b.values
}

func (b *builtinManager) Contains(s string) bool {
	_, ok := b.lookup[s]
	return ok

}

func (b *builtinManager) IndexOf(s string) int {
	index, ok := b.lookup[s]
	if !ok {
		panic("unknown builtin")
	}
	return index
}

//-----------------------------------------------------------------

// StandardBuiltins containts the built-ins that are
// pure functions.  These functions do not do any form of I/O.
var StandardBuiltins = []*BuiltinEntry{
	{"str", BuiltinStr},
	{"len", BuiltinLen},
	{"range", BuiltinRange},
	{"assert", BuiltinAssert},
	{"merge", BuiltinMerge},
	{"chan", BuiltinChan},
	{"type", BuiltinType},
	{"freeze", BuiltinFreeze},
	{"frozen", BuiltinFrozen},
	{"fields", BuiltinFields},
	{"getval", BuiltinGetVal},
	{"setval", BuiltinSetVal},
	{"arity", BuiltinArity},
}

//-----------------------------------------------------------------

// BuiltinStr converts a single value to a Str
var BuiltinStr = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(cx Context, values []Value) (Value, Error) {
		return values[0].ToStr(cx), nil
	})

// BuiltinLen returns the length of a single Lenable
var BuiltinLen = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(cx Context, values []Value) (Value, Error) {
		if ln, ok := values[0].(Lenable); ok {
			return ln.Len(), nil
		}
		return nil, TypeMismatchError("Expected Lenable Type")
	})

// BuiltinRange creates a new Range
var BuiltinRange = NewMultipleNativeFunc(
	[]Type{IntType, IntType},
	[]Type{IntType},
	false,
	func(cx Context, values []Value) (Value, Error) {
		from := values[0].(Int)
		to := values[1].(Int)
		step := One
		if len(values) == 3 {
			step = values[2].(Int)

		}
		return NewRange(from.IntVal(), to.IntVal(), step.IntVal())
	})

// BuiltinAssert asserts that a single Bool is True
var BuiltinAssert = NewFixedNativeFunc(
	[]Type{BoolType},
	false,
	func(cx Context, values []Value) (Value, Error) {
		b := values[0].(Bool)
		if b.BoolVal() {
			return True, nil
		}
		return nil, AssertionFailedError()
	})

// BuiltinMerge merges structs together.
var BuiltinMerge = NewVariadicNativeFunc(
	[]Type{StructType, StructType},
	StructType,
	false,
	func(cx Context, values []Value) (Value, Error) {
		structs := make([]Struct, len(values))
		for i, v := range values {
			if s, ok := v.(Struct); ok {
				structs[i] = s
			} else {
				return nil, TypeMismatchError("Expected Struct")
			}
		}

		return MergeStructs(structs), nil
	})

// BuiltinChan creates a new Chan.  If an Int is passed in,
// it is used to create a buffered Chan.
var BuiltinChan = NewMultipleNativeFunc(
	[]Type{},
	[]Type{IntType},
	true,
	func(cx Context, values []Value) (Value, Error) {

		if len(values) == 0 {
			return NewChan(), nil
		}

		size := values[0].(Int)
		return NewBufferedChan(int(size.IntVal())), nil
	})

// BuiltinType returns the Str representation of the Type of a single Value
var BuiltinType = NewFixedNativeFunc(
	[]Type{AnyType},
	// Subtlety: Null has a type, but for the purposes of type()
	// we are going to pretend that it doesn't
	false,
	func(cx Context, values []Value) (Value, Error) {
		t := values[0].Type()
		return NewStr(t.String()), nil
	})

// BuiltinFreeze freezes a single Value.
var BuiltinFreeze = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(cx Context, values []Value) (Value, Error) {
		return values[0].Freeze()
	})

// BuiltinFrozen returns whether a single Value is Frozen.
var BuiltinFrozen = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(cx Context, values []Value) (Value, Error) {
		return values[0].Frozen()
	})

// BuiltinFields returns the fields of a Struct
var BuiltinFields = NewFixedNativeFunc(
	[]Type{StructType},
	false,
	func(cx Context, values []Value) (Value, Error) {
		st := values[0].(*_struct)
		fields := st.smap.fieldNames()
		result := make([]Value, len(fields))
		for i, k := range fields {
			result[i] = NewStr(k)
		}
		return NewSet(cx, result)
	})

// BuiltinGetVal gets the Value associated with a Struct's field name.
var BuiltinGetVal = NewFixedNativeFunc(
	[]Type{StructType, StrType},
	false,
	func(cx Context, values []Value) (Value, Error) {
		st := values[0].(*_struct)
		field := values[1].(Str)

		return st.GetField(cx, field)
	})

// BuiltinSetVal sets the Value associated with a Struct's field name.
var BuiltinSetVal = NewFixedNativeFunc(
	[]Type{StructType, StrType, AnyType},
	false,
	func(cx Context, values []Value) (Value, Error) {
		st := values[0].(*_struct)
		field := values[1].(Str)
		val := values[2]

		err := st.SetField(cx, field, val)
		if err != nil {
			return nil, err
		}
		return val, nil
	})

// BuiltinArity returns the arity of a function.
var BuiltinArity = NewFixedNativeFunc(
	[]Type{FuncType},
	false,
	func(cx Context, values []Value) (Value, Error) {
		//		st, err := NewStruct([]Field{
		//			NewField("min", true, NewInt(int64(f.MinArity()))),
		//			NewField("max", true, NewInt(int64(f.MaxArity())))}, true)
		//		if err != nil {
		//			panic("invalid struct")
		//		}
		//		return st, nil
		panic("TODO")
	})

//-----------------------------------------------------------------

// UnsandboxedBuiltins are builtins that are not pure functions
var UnsandboxedBuiltins = []*BuiltinEntry{
	{"print", BuiltinPrint},
	{"println", BuiltinPrintln},
}

// BuiltinPrint prints to stdout.
var BuiltinPrint = NewVariadicNativeFunc(
	[]Type{},
	AnyType,
	true,
	func(cx Context, values []Value) (Value, Error) {
		for _, v := range values {
			fmt.Print(v.ToStr(cx).String())
		}

		return Null, nil
	})

// BuiltinPrintln prints to stdout.
var BuiltinPrintln = NewVariadicNativeFunc(
	[]Type{},
	AnyType,
	true,
	func(cx Context, values []Value) (Value, Error) {
		for _, v := range values {
			fmt.Print(v.ToStr(cx).String())
		}
		fmt.Println()

		return Null, nil
	})
