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

// SandboxBuiltins containts the built-ins that are
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

// CommandLineBuiltins consists of the print(), and println()
var CommandLineBuiltins = []*BuiltinEntry{
	{"print", BuiltinPrint},
	{"println", BuiltinPrintln},
}

//-----------------------------------------------------------------

// BuiltinPrint prints to stdout.
var BuiltinPrint = NewObsoleteFunc(
	0, -1,
	func(cx Context, values []Value) (Value, Error) {
		for _, v := range values {
			fmt.Print(v.ToStr(cx).String())
		}

		return Null, nil
	})

// BuiltinPrintln prints to stdout.
var BuiltinPrintln = NewObsoleteFunc(
	0, -1,
	func(cx Context, values []Value) (Value, Error) {
		for _, v := range values {
			fmt.Print(v.ToStr(cx).String())
		}
		fmt.Println()

		return Null, nil
	})

//-----------------------------------------------------------------

// BuiltinStr converts a single value to a Str
var BuiltinStr = NewObsoleteFuncValue(
	func(cx Context, val Value) (Value, Error) {
		return val.ToStr(cx), nil
	})

// BuiltinLen returns the length of a single Lenable
var BuiltinLen = NewObsoleteFunc(
	1, 1,
	func(cx Context, values []Value) (Value, Error) {
		if ln, ok := values[0].(Lenable); ok {
			return ln.Len(), nil
		}
		return nil, TypeMismatchError("Expected Lenable Type")
	})

// BuiltinRange creates a new Range
var BuiltinRange = NewObsoleteFunc(
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
	})

// BuiltinAssert asserts that a single Bool is True
var BuiltinAssert = NewObsoleteFuncBool(
	func(cx Context, b Bool) (Value, Error) {
		if b.BoolVal() {
			return True, nil
		}
		return nil, AssertionFailedError()
	})

// BuiltinMerge merges structs together.
var BuiltinMerge = NewObsoleteFunc(
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
	})

// BuiltinChan creates a new Chan.  If an Int is passed in,
// it is used to create a buffered Chan.
var BuiltinChan = NewObsoleteFunc(
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
	})

// BuiltinType returns the Str representation of the Type of a single Value
var BuiltinType = NewObsoleteFuncValue(
	func(cx Context, val Value) (Value, Error) {
		// Null has a type, but for the purposes of type()
		// we are going to pretend that it doesn't
		if val == Null {
			return nil, NullValueError()
		}
		t := val.Type()
		return NewStr(t.String()), nil
	})

// BuiltinFreeze freezes a single Value.
var BuiltinFreeze = NewObsoleteFuncValue(
	func(cx Context, val Value) (Value, Error) {
		return val.Freeze()
	})

// BuiltinFrozen returns whether a single Value is Frozen.
var BuiltinFrozen = NewObsoleteFuncValue(
	func(cx Context, val Value) (Value, Error) {
		return val.Frozen()
	})

// BuiltinFields returns the fields of a Struct
var BuiltinFields = NewObsoleteFuncStruct(
	func(cx Context, st Struct) (Value, Error) {
		fields := st.(*_struct).smap.fieldNames()
		result := make([]Value, len(fields))
		for i, k := range fields {
			result[i] = NewStr(k)
		}
		return NewSet(cx, result)
	})

// BuiltinGetVal gets the Value associated with a Struct's field name.
var BuiltinGetVal = NewObsoleteFunc(
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
	})

// BuiltinSetVal sets the Value associated with a Struct's field name.
var BuiltinSetVal = NewObsoleteFunc(
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
	})

// BuiltinArity returns the arity of a function.
var BuiltinArity = NewObsoleteFuncFunc(
	func(cx Context, f Func) (Value, Error) {
		//		st, err := NewStruct([]Field{
		//			NewField("min", true, NewInt(int64(f.MinArity()))),
		//			NewField("max", true, NewInt(int64(f.MaxArity())))}, true)
		//		if err != nil {
		//			panic("invalid struct")
		//		}
		//		return st, nil
		return Null, nil
	})

//// BuiltinOpenPlugin wraps a plugin in a struct
//var BuiltinOpenPlugin = NewObsoleteFuncStr(
//	func(cx Context, name Str) (Value, Error) {
//
//		// open up the plugin
//		plugPath := cx.HomePath() + "/lib/" + name.String() + "/" + name.String() + ".so"
//		plug, err := plugin.Open(plugPath)
//		if err != nil {
//			return nil, PluginError(name.String(), err)
//		}
//
//		// define lookup function
//		lookup := NewObsoleteFuncStr(
//			func(cx Context, s Str) (Value, Error) {
//
//				name := s.String()
//
//				sym, e2 := plug.Lookup(name)
//				if e2 != nil {
//					return nil, PluginError(name, e2)
//				}
//
//				value, ok := sym.(*Value)
//				if !ok {
//					return nil, PluginError(name, fmt.Errorf(
//						"plugin symbol '%s' is not a Value: %s",
//						s.String(),
//						reflect.TypeOf(sym)))
//				}
//				return *value, nil
//			})
//
//		// done
//		stc, err := NewStruct([]Field{
//			NewField("lookup", true, lookup),
//		}, true)
//		if err != nil {
//			panic("unreachable")
//		}
//
//		return stc, nil
//	})
