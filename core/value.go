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
	TypeOf() Type

	HashCode() (Int, Error)
	Eq(Value) Bool
	ToStr() Str

	Cmp(Value) (Int, Error)
	Plus(Value) (Value, Error)

	GetField(Str) (Value, Error)
}

//---------------------------------------------------------------
// Shared Interfaces

type (
	Getable interface {
		Get(Value) (Value, Error)
	}

	Indexable interface {
		Getable
		Set(Value, Value) Error
	}

	Lenable interface {
		Len() Int
	}

	Sliceable interface {
		Slice(Value, Value) (Value, Error)
		SliceFrom(Value) (Value, Error)
		SliceTo(Value) (Value, Error)
	}

	Iterable interface {
		NewIterator() Iterator
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
	}

	Number interface {
		Basic
		FloatVal() float64
		IntVal() int64

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

		Add(Value) Error
		AddAll(Value) Error
		Clear()
		Contains(Value) (Bool, Error)
		IndexOf(Value) Int
		IsEmpty() Bool
		Join(Str) Str

		Values() []Value
	}

	Range interface {
		Composite
		Getable
		Lenable
		Sliceable
		Iterable

		From() Int
		To() Int
		Step() Int
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

		AddAll(Value) Error
		Clear()
		ContainsKey(Value) (Bool, Error)
		IsEmpty() Bool
	}

	Set interface {
		Composite
		Lenable
		Iterable

		Add(Value) Error
		AddAll(Value) Error
		Clear()
		Contains(Value) (Bool, Error)
		IsEmpty() Bool
	}

	Struct interface {
		Composite
		Indexable

		Keys() []string
		Has(Value) (Bool, Error)
		InitField(Str, Value) Error
		SetField(Str, Value) Error
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
