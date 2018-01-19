// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"fmt"
)

//---------------------------------------------------------------
// Value

type Value interface {
	Type() Type

	Freeze() (Value, Error)
	Frozen() (Bool, Error)

	Eq(Context, Value) (Bool, Error)
	HashCode(Context) (Int, Error)
	ToStr(Context) Str
	Cmp(Context, Value) (Int, Error)
	GetField(Context, Str) (Value, Error)
}

//---------------------------------------------------------------
// Shared Interfaces

type (
	Getable interface {
		Get(Context, Value) (Value, Error)
	}

	Indexable interface {
		Getable
		Set(Context, Value, Value) Error
	}

	Lenable interface {
		Len() Int
	}

	Sliceable interface {
		Slice(Context, Value, Value) (Value, Error)
		SliceFrom(Context, Value) (Value, Error)
		SliceTo(Context, Value) (Value, Error)
	}

	Iterable interface {
		NewIterator(Context) Iterator
	}
)

//---------------------------------------------------------------
// Basic

type (
	Basic interface {
		Value
		basicMarker()
	}

	Null interface {
		Basic
	}

	Bool interface {
		Basic
		BoolVal() bool

		Not() Bool
	}

	Str interface {
		Basic
		fmt.Stringer

		Getable
		Lenable
		Sliceable
		Iterable

		Concat(Str) Str
	}

	Number interface {
		Basic
		FloatVal() float64
		IntVal() int64

		Add(Value) (Number, Error)
		Sub(Value) (Number, Error)
		Mul(Value) (Number, Error)
		Div(Value) (Number, Error)
		Negate() Number
	}

	Float interface {
		Number
	}

	Int interface {
		Number

		Rem(Value) (Int, Error)
		BitAnd(Value) (Int, Error)
		BitOr(Value) (Int, Error)
		BitXOr(Value) (Int, Error)
		LeftShift(Value) (Int, Error)
		RightShift(Value) (Int, Error)
		Complement() Int
	}
)

//---------------------------------------------------------------
// Composite

type (
	Composite interface {
		Value
		compositeMarker()
	}

	List interface {
		Composite
		Indexable
		Lenable
		Iterable
		Sliceable

		IsEmpty() Bool
		Clear() Error

		Contains(Context, Value) (Bool, Error)
		IndexOf(Context, Value) (Int, Error)
		Join(Context, Str) Str

		Add(Context, Value) Error
		AddAll(Context, Value) Error
		Remove(Context, Int) Error

		Values() []Value

		Map(Context, func(Value) (Value, Error)) (Value, Error)
		Reduce(Context, Value, func(Value, Value) (Value, Error)) (Value, Error)
		Filter(Context, func(Value) (Value, Error)) (Value, Error)
	}

	Range interface {
		Composite
		Getable
		Lenable
		Iterable

		From() Int
		To() Int
		Step() Int
		Count() Int
	}

	Tuple interface {
		Composite
		Getable
		Lenable
	}

	Dict interface {
		Composite
		Indexable
		Lenable
		Iterable

		IsEmpty() Bool
		Clear() Error

		ContainsKey(Context, Value) (Bool, Error)
		AddAll(Context, Value) Error
		Remove(Context, Value) (Bool, Error)
	}

	Set interface {
		Composite
		Lenable
		Iterable

		IsEmpty() Bool
		Clear() Error

		Contains(Context, Value) (Bool, Error)
		Add(Context, Value) Error
		AddAll(Context, Value) Error
		Remove(Context, Value) (Bool, Error)
	}

	Struct interface {
		Composite

		FieldNames() []string
		Has(Value) (Bool, Error)

		InitField(Context, Str, Value) Error
		SetField(Context, Str, Value) Error
	}

	Iterator interface {
		Struct
		IterNext() Bool
		IterGet() (Value, Error)
	}
)

//---------------------------------------------------------------
// Func

// Func represents an instance of a function
type Func interface {
	Value
	MinArity() int
	MaxArity() int // by convention, -1 means variadic

	Invoke(Context, []Value) (Value, Error)

	funcMarker()
}

//---------------------------------------------------------------
// Chan

// Chan represents a channel
type Chan interface {
	Value
	chanMarker()
}

//---------------------------------------------------------------
// Type

type Type int

const (
	TNULL Type = iota
	TBOOL
	TSTR
	TINT
	TFLOAT
	TFUNC
	TLIST
	TRANGE
	TTUPLE
	TDICT
	TSET
	TSTRUCT
	TCHAN
)

func (t Type) String() string {
	switch t {
	case TNULL:
		return "Null"
	case TBOOL:
		return "Bool"
	case TSTR:
		return "Str"
	case TINT:
		return "Int"
	case TFLOAT:
		return "Float"
	case TFUNC:
		return "Func"
	case TLIST:
		return "List"
	case TRANGE:
		return "Range"
	case TTUPLE:
		return "Tuple"
	case TDICT:
		return "Dict"
	case TSET:
		return "Set"
	case TSTRUCT:
		return "Struct"
	case TCHAN:
		return "Chan"

	default:
		panic("unreachable")
	}
}
