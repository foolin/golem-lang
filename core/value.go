// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

//---------------------------------------------------------------
// Value

// Value is the base interface for every value in the Golem Language
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
	// Getable is a value that supports the dot operator
	Getable interface {
		Get(Context, Value) (Value, Error)
	}

	// Indexable is a value that supports the index operator
	Indexable interface {
		Getable
		Set(Context, Value, Value) Error
	}

	// Lenable is a value that has a length
	Lenable interface {
		Len() Int
	}

	// Sliceable is a value that can be sliced
	Sliceable interface {
		Slice(Context, Value, Value) (Value, Error)
		SliceFrom(Context, Value) (Value, Error)
		SliceTo(Context, Value) (Value, Error)
	}

	// Iterable is a value that can be iterated
	Iterable interface {
		NewIterator(Context) Iterator
	}
)

//---------------------------------------------------------------
// Basic

type (
	// Basic represents the immutable types null, bool, str, int and float
	Basic interface {
		Value
		basicMarker()
	}

	// Null is the null value
	Null interface {
		Basic
	}

	// Bool is the boolean value -- true or false
	Bool interface {
		Basic
		BoolVal() bool

		Not() Bool
	}

	// Str is a string -- defined in golem as a sequence of runes
	Str interface {
		Basic
		fmt.Stringer

		Getable
		Lenable
		Sliceable
		Iterable

		Concat(Str) Str
	}

	// Number is a number
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

	// Float is a float64
	Float interface {
		Number
	}

	// Int is an int64
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
	// Composite is a Value that is composed of other values -- List, Range, Tuple,
	// Dict, Set, Struct
	Composite interface {
		Value
		compositeMarker()
	}

	// List is an indexable sequence of values
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

	// Range is an immutable, iterable representation of a  sequence of integers
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

	// Tuple is an immutable sequence of two or more values
	Tuple interface {
		Composite
		Getable
		Lenable
	}

	// Dict is an associative array, a.k.a Hash Map
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

	// Set is a set of unique values
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

	// Struct is a composite collection of names values
	Struct interface {
		Composite

		FieldNames() []string
		Has(Value) (Bool, Error)

		InitField(Context, Str, Value) Error
		SetField(Context, Str, Value) Error
	}

	// Iterator iterates over a sequence of values
	Iterator interface {
		Struct
		IterNext() Bool
		IterGet() (Value, Error)
	}
)

//---------------------------------------------------------------
// Func

// Func is a function
type Func interface {
	Value
	MinArity() int
	MaxArity() int // by convention, -1 means variadic

	Invoke(Context, []Value) (Value, Error)

	funcMarker()
}

//---------------------------------------------------------------
// Chan

// Chan is a goroutine channel
type Chan interface {
	Value
	chanMarker()
}

//---------------------------------------------------------------
// Type

// Type represents all of the possible types of a Golem Value
type Type int

// All possible kinds of Type
const (
	NullType Type = iota
	BoolType
	StrType
	IntType
	FloatType
	FuncType
	ListType
	RangeType
	TupleType
	DictType
	SetType
	StructType
	ChanType
)

func (t Type) String() string {
	switch t {
	case NullType:
		return "Null"
	case BoolType:
		return "Bool"
	case StrType:
		return "Str"
	case IntType:
		return "Int"
	case FloatType:
		return "Float"
	case FuncType:
		return "Func"
	case ListType:
		return "List"
	case RangeType:
		return "Range"
	case TupleType:
		return "Tuple"
	case DictType:
		return "Dict"
	case SetType:
		return "Set"
	case StructType:
		return "Struct"
	case ChanType:
		return "Chan"

	default:
		panic("unreachable")
	}
}
