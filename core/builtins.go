// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

//-----------------------------------------------------------------

type (
	BuiltinManager interface {
		Builtins() []Value
		Contains(s string) bool
		IndexOf(s string) int
	}

	BuiltinEntry struct {
		Name  string
		Value Value
	}

	builtinManager struct {
		values []Value
		lookup map[string]int
	}
)

func NewBuiltinManager(entries []*BuiltinEntry) BuiltinManager {
	values := make([]Value, len(entries), len(entries))
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

// print() and println() do IO, so they should not be included
// in a sandboxed environment.
var CommandLineBuiltins = append(
	SandboxBuiltins,
	[]*BuiltinEntry{
		{"print", BuiltinPrint},
		{"println", BuiltinPrintln}}...)

//-----------------------------------------------------------------
// print() and println() do IO, so they should not be included
// in a sandboxed environment.

var BuiltinPrint = &nativeFunc{
	0, -1,
	func(cx Context, values []Value) (Value, Error) {
		for _, v := range values {
			fmt.Print(v.ToStr(cx).String())
		}

		return NULL, nil
	}}

var BuiltinPrintln = &nativeFunc{
	0, -1,
	func(cx Context, values []Value) (Value, Error) {
		for _, v := range values {
			fmt.Print(v.ToStr(cx).String())
		}
		fmt.Println()

		return NULL, nil
	}}

//-----------------------------------------------------------------

var BuiltinStr = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		return values[0].ToStr(cx), nil
	}}

var BuiltinLen = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		if ln, ok := values[0].(Lenable); ok {
			return ln.Len(), nil
		} else {
			return nil, TypeMismatchError("Expected Lenable Type")
		}
	}}

var BuiltinRange = &nativeFunc{
	2, 3,
	func(cx Context, values []Value) (Value, Error) {
		from, ok := values[0].(Int)
		if !ok {
			return nil, TypeMismatchError("Expected 'Int'")
		}

		to, ok := values[1].(Int)
		if !ok {
			return nil, TypeMismatchError("Expected 'Int'")
		}

		step := ONE
		if len(values) == 3 {
			step, ok = values[2].(Int)
			if !ok {
				return nil, TypeMismatchError("Expected 'Int'")
			}
		}

		return NewRange(from.IntVal(), to.IntVal(), step.IntVal())
	}}

var BuiltinAssert = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		b, ok := values[0].(Bool)
		if !ok {
			return nil, TypeMismatchError("Expected 'Bool'")
		}

		if b.BoolVal() {
			return TRUE, nil
		} else {
			return nil, AssertionFailedError()
		}
	}}

var BuiltinMerge = &nativeFunc{
	2, -1,
	func(cx Context, values []Value) (Value, Error) {
		structs := make([]Struct, len(values), len(values))
		for i, v := range values {
			if s, ok := v.(Struct); ok {
				structs[i] = s
			} else {
				return nil, TypeMismatchError("Expected 'Struct'")
			}
		}

		return MergeStructs(structs), nil
	}}

var BuiltinChan = &nativeFunc{
	0, 1,
	func(cx Context, values []Value) (Value, Error) {
		switch len(values) {
		case 0:
			return NewChan(), nil
		case 1:
			size, ok := values[0].(Int)
			if !ok {
				return nil, TypeMismatchError("Expected 'Int'")
			}
			return NewBufferedChan(int(size.IntVal())), nil

		default:
			panic("arity mismatch")
		}
	}}

var BuiltinType = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		// Null has a type, but for the purposes of type()
		// we are going to pretend it doesn't
		if values[0] == NULL {
			return nil, NullValueError()
		} else {
			_type := values[0].Type()
			return MakeStr(_type.String()), nil
		}
	}}

var BuiltinFreeze = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		return values[0].Freeze()
	}}

var BuiltinFrozen = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		return values[0].Frozen()
	}}

var BuiltinFields = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		st, ok := values[0].(*_struct)
		if !ok {
			return nil, TypeMismatchError("Expected 'Struct'")
		}

		fields := st.smap.fieldNames()
		result := make([]Value, len(fields), len(fields))
		for i, k := range fields {
			result[i] = MakeStr(k)
		}
		return NewSet(cx, result), nil
	}}

var BuiltinGetVal = &nativeFunc{
	2, 2,
	func(cx Context, values []Value) (Value, Error) {
		st, ok := values[0].(*_struct)
		if !ok {
			return nil, TypeMismatchError("Expected 'Struct'")
		}
		field, ok := values[1].(Str)
		if !ok {
			return nil, TypeMismatchError("Expected 'Str'")
		}

		return st.GetField(cx, field)
	}}

var BuiltinSetVal = &nativeFunc{
	3, 3,
	func(cx Context, values []Value) (Value, Error) {
		st, ok := values[0].(*_struct)
		if !ok {
			return nil, TypeMismatchError("Expected 'Struct'")
		}
		field, ok := values[1].(Str)
		if !ok {
			return nil, TypeMismatchError("Expected 'Str'")
		}
		val := values[2]

		err := st.SetField(cx, field, val)
		if err != nil {
			return nil, err
		} else {
			return val, nil
		}
	}}

var BuiltinArity = &nativeFunc{
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		f, ok := values[0].(Func)
		if !ok {
			return nil, TypeMismatchError("Expected 'Func'")
		}

		st, err := NewStruct([]Field{
			NewField("min", true, MakeInt(int64(f.MinArity()))),
			NewField("max", true, MakeInt(int64(f.MaxArity())))}, true)
		if err != nil {
			panic("invalid struct")
		}
		return st, nil
	}}
