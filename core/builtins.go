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
	{"assert", BuiltinAssert},

	{"freeze", BuiltinFreeze},
	{"frozen", BuiltinFrozen},
	{"iter", BuiltinIter},
	{"len", BuiltinLen},
	{"range", BuiltinRange},
	{"str", BuiltinStr},
	{"type", BuiltinType},

	{"fields", BuiltinFields},
	{"getField", BuiltinGetField},
	{"hasField", BuiltinHasField},
	{"setField", BuiltinSetField},

	{"arity", BuiltinArity},
	{"chan", BuiltinChan},
	{"merge", BuiltinMerge},
}

//-----------------------------------------------------------------

// BuiltinStr converts a single value to a Str
var BuiltinStr = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		return params[0].ToStr(ev)
	})

// BuiltinLen returns the length of a single Lenable
var BuiltinLen = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {

		if params[0].Type() == NullType {
			return nil, NullValueError()
		}

		if ln, ok := params[0].(Lenable); ok {
			return ln.Len(ev)
		}
		return nil, LenableMismatchError(params[0].Type())
	})

// BuiltinIter returns the length of a single Lenable
var BuiltinIter = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {

		if params[0].Type() == NullType {
			return nil, NullValueError()
		}

		if ibl, ok := params[0].(Iterable); ok {
			return ibl.NewIterator(ev)
		}
		return nil, IterableMismatchError(params[0].Type())
	})

// BuiltinRange creates a new Range
var BuiltinRange = NewMultipleNativeFunc(
	[]Type{IntType, IntType},
	[]Type{IntType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		from := params[0].(Int)
		to := params[1].(Int)
		step := One
		if len(params) == 3 {
			step = params[2].(Int)

		}
		return NewRange(from.IntVal(), to.IntVal(), step.IntVal())
	})

// BuiltinAssert asserts that a single Bool is True
var BuiltinAssert = NewFixedNativeFunc(
	[]Type{BoolType}, false,
	func(ev Eval, params []Value) (Value, Error) {
		b := params[0].(Bool)
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
	func(ev Eval, params []Value) (Value, Error) {
		structs := make([]Struct, len(params))
		for i, v := range params {
			structs[i] = v.(Struct)
		}
		return MergeStructs(structs)
	})

// BuiltinChan creates a new Chan.  If an Int is passed in,
// it is used to create a buffered Chan.
var BuiltinChan = NewMultipleNativeFunc(
	[]Type{},
	[]Type{IntType},
	false,
	func(ev Eval, params []Value) (Value, Error) {

		if len(params) == 0 {
			return NewChan(), nil
		}

		size := params[0].(Int)
		return NewBufferedChan(int(size.IntVal())), nil
	})

// BuiltinType returns the Str representation of the Type of a single Value
var BuiltinType = NewFixedNativeFunc(
	[]Type{AnyType},
	// Subtlety: Null has a type, but for the purposes of type()
	// we are going to pretend that it doesn't
	true,
	func(ev Eval, params []Value) (Value, Error) {

		if params[0].Type() == NullType {
			return nil, NullValueError()
		}

		return NewStr(params[0].Type().String()), nil
	})

// BuiltinFreeze freezes a single Value.
var BuiltinFreeze = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		return params[0].Freeze(ev)
	})

// BuiltinFrozen returns whether a single Value is Frozen.
var BuiltinFrozen = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		return params[0].Frozen(ev)
	})

// BuiltinFields returns the fields of a Struct
var BuiltinFields = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {

		fields, err := params[0].FieldNames()
		if err != nil {
			return nil, err
		}

		entries := make([]Value, len(fields))
		for i, k := range fields {
			entries[i] = NewStr(k)
		}
		return NewSet(ev, entries)
	})

// BuiltinGetField gets the Value associated with a Struct's field name.
var BuiltinGetField = NewFixedNativeFunc(
	[]Type{AnyType, StrType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		field := params[1].(Str)

		return params[0].GetField(ev, field.String())
	})

// BuiltinHasField gets the Value associated with a Struct's field name.
var BuiltinHasField = NewFixedNativeFunc(
	[]Type{AnyType, StrType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		field := params[1].(Str)

		b, err := params[0].HasField(field.String())
		if err != nil {
			return nil, err
		}
		return NewBool(b), nil
	})

// BuiltinSetField sets the Value associated with a Struct's field name.
var BuiltinSetField = NewFixedNativeFunc(
	[]Type{StructType, StrType, AnyType},
	true,
	func(ev Eval, params []Value) (Value, Error) {

		if params[0].Type() == NullType {
			return nil, NullValueError()
		}
		if params[1].Type() == NullType {
			return nil, NullValueError()
		}

		st := params[0].(Struct)
		fld := params[1].(Str)

		err := st.SetField(ev, fld.String(), params[2])
		if err != nil {
			return nil, err
		}
		return Null, nil
	})

// BuiltinArity returns the arity of a function.
var BuiltinArity = NewFixedNativeFunc(
	[]Type{FuncType},
	false,
	func(ev Eval, params []Value) (Value, Error) {

		fn := params[0].(Func)
		a := fn.Arity()
		k := NewStr(a.Kind.String())
		r := NewInt(int64(a.Required))

		fields := map[string]Field{
			"kind":     NewReadonlyField(k),
			"required": NewReadonlyField(r),
		}

		if a.Kind == MultipleArity {
			o := NewInt(int64(a.Optional))
			fields["optional"] = NewReadonlyField(o)
		}

		return NewFrozenFieldStruct(fields)
	})

//-----------------------------------------------------------------

// UnsandboxedBuiltins are builtins that are not pure functions
var UnsandboxedBuiltins = []*BuiltinEntry{
	{"print", BuiltinPrint},
	{"println", BuiltinPrintln},
}

// BuiltinPrint prints to stdout.
var BuiltinPrint = NewVariadicNativeFunc(
	[]Type{}, AnyType, true,
	func(ev Eval, params []Value) (Value, Error) {
		for _, v := range params {
			s, err := v.ToStr(ev)
			if err != nil {
				return nil, err
			}
			fmt.Print(s.String())
		}

		return Null, nil
	})

// BuiltinPrintln prints to stdout.
var BuiltinPrintln = NewVariadicNativeFunc(
	[]Type{}, AnyType, true,
	func(ev Eval, params []Value) (Value, Error) {
		for _, v := range params {
			s, err := v.ToStr(ev)
			if err != nil {
				return nil, err
			}
			fmt.Print(s.String())
		}
		fmt.Println()

		return Null, nil
	})
