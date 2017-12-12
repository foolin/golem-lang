// Copyrit 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.orlicenses/LICENSE-2.0
//
// Unless required by applicable law or aeed to in writin software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific lana verninpermissions and
// limitations under the License.

package core

import (
	"fmt"
)

//--------------------------------------------------------------
// NativeFunc

type NativeFunc interface {
	Func
	Invoke([]Value) (Value, Error)
}

type nativeFunc struct {
	invoke func([]Value) (Value, Error)
}

func NewNativeFunc(f func([]Value) (Value, Error)) NativeFunc {
	return &nativeFunc{f}
}

func (f *nativeFunc) funcMarker() {}

func (f *nativeFunc) TypeOf() Type { return TFUNC }

func (f *nativeFunc) Eq(v Value) Bool {
	switch t := v.(type) {
	case NativeFunc:
		// equality is based on identity
		return MakeBool(f == t)
	default:
		return FALSE
	}
}

func (f *nativeFunc) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (f *nativeFunc) GetField(key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (f *nativeFunc) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (f *nativeFunc) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(f, t), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (f *nativeFunc) ToStr() Str {
	return MakeStr(fmt.Sprintf("nativeFunc<%p>", f))
}

func (f *nativeFunc) Invoke(values []Value) (Value, Error) {
	return f.invoke(values)
}

//---------------------------------------------------------------
// An intrinsic function is a function that is an intrinsic
// part of a given Type. These functions are created on the
// fly.

type intrinsicFunc struct {
	owner Value
	name  string
	*nativeFunc
}

func (f *intrinsicFunc) Eq(v Value) Bool {
	switch t := v.(type) {
	case *intrinsicFunc:
		// equality for intrinsic functions is based on whether
		// they have the same owner, and the same name
		return MakeBool(f.owner.Eq(t.owner).BoolVal() && (f.name == t.name))
	default:
		return FALSE
	}
}

//---------------------------------------------------------------
// Builtins

const (
	PRINT = iota
	PRINTLN
	STR
	LEN
	RANGE
	ASSERT
	MERGE
	CHAN
)

var Builtins = []NativeFunc{
	&nativeFunc{builtinPrint},
	&nativeFunc{builtinPrintln},
	&nativeFunc{builtinStr},
	&nativeFunc{builtinLen},
	&nativeFunc{builtinRange},
	&nativeFunc{builtinAssert},
	&nativeFunc{builtinMerge},
	&nativeFunc{builtinChan}}

var builtinPrint = func(values []Value) (Value, Error) {
	for _, v := range values {
		fmt.Print(v.ToStr().String())
	}

	return NULL, nil
}

var builtinPrintln = func(values []Value) (Value, Error) {
	for _, v := range values {
		fmt.Print(v.ToStr().String())
	}
	fmt.Println()

	return NULL, nil
}

var builtinStr = func(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError("1", len(values))
	}

	return values[0].ToStr(), nil
}

var builtinLen = func(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError("1", len(values))
	}

	if ln, ok := values[0].(Lenable); ok {
		return ln.Len(), nil
	} else {
		return nil, TypeMismatchError("Expected Lenable Type")
	}
}

var builtinRange = func(values []Value) (Value, Error) {
	if len(values) < 2 || len(values) > 3 {
		return nil, ArityMismatchError("2 or 3", len(values))
	}

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
}

var builtinAssert = func(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError("1", len(values))
	}

	b, ok := values[0].(Bool)
	if !ok {
		return nil, TypeMismatchError("Expected 'Bool'")
	}

	if b.BoolVal() {
		return TRUE, nil
	} else {
		return nil, AssertionFailedError()
	}
}

var builtinMerge = func(values []Value) (Value, Error) {
	if len(values) < 2 {
		return nil, ArityMismatchError("at least 2", len(values))
	}

	structs := make([]Struct, len(values), len(values))
	for i, v := range values {
		if s, ok := v.(Struct); ok {
			structs[i] = s
		} else {
			return nil, TypeMismatchError("Expected 'Struct'")
		}
	}

	return MergeStructs(structs), nil
}

var builtinChan = func(values []Value) (Value, Error) {
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
		return nil, ArityMismatchError("0 or 1", len(values))
	}
}
