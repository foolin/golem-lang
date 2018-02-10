// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

//-----------------------------------------------------------------

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

// SandboxBuiltins is the 'standard' list of built-ins that are
// pure functions.  These functions do not do any form of I/O.
var SandboxBuiltins = []*BuiltinEntry{
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

// CommandLineBuiltins consists of the SandboxBuiltins, plus print() and println().
var CommandLineBuiltins = append(
	SandboxBuiltins,
	[]*BuiltinEntry{
		{"print", BuiltinPrint},
		{"println", BuiltinPrintln}}...)

//-----------------------------------------------------------------

// BuiltinPrint prints to stdout.
var BuiltinPrint = &nativeFunc{
	0, -1,
	func(cx Context, values []Value) (Value, Error) {
		for _, v := range values {
			fmt.Print(v.ToStr(cx).String())
		}

		return NullValue, nil
	}}

// BuiltinPrintln prints to stdout.
var BuiltinPrintln = &nativeFunc{
	0, -1,
	func(cx Context, values []Value) (Value, Error) {
		for _, v := range values {
			fmt.Print(v.ToStr(cx).String())
		}
		fmt.Println()

		return NullValue, nil
	}}

//-----------------------------------------------------------------

// BuiltinStr converts a single value to a Str
var BuiltinStr = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		return values[0].ToStr(cx), nil
	}}

// BuiltinLen returns the length of a single Lenable
var BuiltinLen = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		if ln, ok := values[0].(Lenable); ok {
			return ln.Len(), nil
		}
		return nil, TypeMismatchError("Expected Lenable Type")
	}}

// BuiltinRange creates a new Range
var BuiltinRange = &nativeFunc{
	2, 3,
	func(cx Context, values []Value) (Value, Error) {
		from, ok := values[0].(Int)
		if !ok {
			return nil, TypeMismatchError("Expected Int")
		}

		to, ok := values[1].(Int)
		if !ok {
			return nil, TypeMismatchError("Expected Int")
		}

		step := One
		if len(values) == 3 {
			step, ok = values[2].(Int)
			if !ok {
				return nil, TypeMismatchError("Expected Int")
			}
		}

		return NewRange(from.IntVal(), to.IntVal(), step.IntVal())
	}}

// BuiltinAssert asserts that a single Bool is True
var BuiltinAssert = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		b, ok := values[0].(Bool)
		if !ok {
			return nil, TypeMismatchError("Expected Bool")
		}

		if b.BoolVal() {
			return True, nil
		}
		return nil, AssertionFailedError()
	}}

// BuiltinMerge merges structs together.
var BuiltinMerge = &nativeFunc{
	2, -1,
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
	}}

// BuiltinChan creates a new Chan.  IfStmt an Int is passed in,
// it is used to create a buffered Chan.
var BuiltinChan = &nativeFunc{
	0, 1,
	func(cx Context, values []Value) (Value, Error) {
		switch len(values) {
		case 0:
			return NewChan(), nil
		case 1:
			size, ok := values[0].(Int)
			if !ok {
				return nil, TypeMismatchError("Expected Int")
			}
			return NewBufferedChan(int(size.IntVal())), nil

		default:
			panic("arity mismatch")
		}
	}}

// BuiltinType returns the Str representation of the Type of a single Value
var BuiltinType = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		// Null has a type, but for the purposes of type()
		// we are going to pretend it doesn't
		if values[0] == NullValue {
			return nil, NullValueError()
		}
		t := values[0].Type()
		return NewStr(t.String()), nil
	}}

// BuiltinFreeze freezes a single Value.
var BuiltinFreeze = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		return values[0].Freeze()
	}}

// BuiltinFrozen returns whether a single Value is Frozen.
var BuiltinFrozen = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		return values[0].Frozen()
	}}

// BuiltinFields returns the fields of a Struct
var BuiltinFields = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		st, ok := values[0].(*_struct)
		if !ok {
			return nil, TypeMismatchError("Expected Struct")
		}

		fields := st.smap.fieldNames()
		result := make([]Value, len(fields))
		for i, k := range fields {
			result[i] = NewStr(k)
		}
		return NewSet(cx, result)
	}}

// BuiltinGetVal gets the Value associated with a Struct's field name.
var BuiltinGetVal = &nativeFunc{
	2, 2,
	func(cx Context, values []Value) (Value, Error) {
		st, ok := values[0].(*_struct)
		if !ok {
			return nil, TypeMismatchError("Expected Struct")
		}
		field, ok := values[1].(Str)
		if !ok {
			return nil, TypeMismatchError("Expected Str")
		}

		return st.GetField(cx, field)
	}}

// BuiltinSetVal sets the Value associated with a Struct's field name.
var BuiltinSetVal = &nativeFunc{
	3, 3,
	func(cx Context, values []Value) (Value, Error) {
		st, ok := values[0].(*_struct)
		if !ok {
			return nil, TypeMismatchError("Expected Struct")
		}
		field, ok := values[1].(Str)
		if !ok {
			return nil, TypeMismatchError("Expected Str")
		}
		val := values[2]

		err := st.SetField(cx, field, val)
		if err != nil {
			return nil, err
		}
		return val, nil
	}}

// BuiltinArity returns the arity of a function.
var BuiltinArity = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		f, ok := values[0].(Func)
		if !ok {
			return nil, TypeMismatchError("Expected Func")
		}

		st, err := NewStruct([]Field{
			NewField("min", true, NewInt(int64(f.MinArity()))),
			NewField("max", true, NewInt(int64(f.MaxArity())))}, true)
		if err != nil {
			panic("invalid struct")
		}
		return st, nil
	}}
