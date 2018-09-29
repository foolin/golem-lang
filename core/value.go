// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

//---------------------------------------------------------------
// Value
//---------------------------------------------------------------

// Value is the base interface for every value in the Golem Language
type Value interface {
	Type() Type

	Freeze(Eval) (Value, Error)
	Frozen(Eval) (Bool, Error)

	Eq(Eval, Value) (Bool, Error)
	HashCode(Eval) (Int, Error)
	ToStr(Eval) (Str, Error)

	FieldNames() ([]string, Error)
	HasField(string) (bool, Error)
	GetField(Eval, string) (Value, Error)
	InvokeField(Eval, string, []Value) (Value, Error)
}

//---------------------------------------------------------------
// Shared Interfaces
//---------------------------------------------------------------

type (

	// Comparable is a value that can be compared to another value of the same type.
	Comparable interface {
		Cmp(Eval, Comparable) (Int, Error)
	}

	// Indexable is a value that supports the index operator
	Indexable interface {
		Get(Eval, Value) (Value, Error)
		Set(Eval, Value, Value) Error
	}

	// Lenable is a value that has a length
	Lenable interface {
		Len(Eval) (Int, Error)
	}

	// Sliceable is a value that can be sliced
	Sliceable interface {
		Slice(Eval, Value, Value) (Value, Error)
		SliceFrom(Eval, Value) (Value, Error)
		SliceTo(Eval, Value) (Value, Error)
	}

	// Iterable is a value that can be iterated
	Iterable interface {
		NewIterator(Eval) (Iterator, Error)
	}
)

//---------------------------------------------------------------
// Basic
//---------------------------------------------------------------

type (
	// Basic represents the immutable types Null, Bool, Str, Int and Float
	Basic interface {
		Value
		basicMarker()
	}

	// NullValue is the null value. The only instance of NullValue is Null.
	NullValue interface {
		Basic
	}

	// Bool is a boolean value.  The only instances of Bool are True and False.
	Bool interface {
		Basic
		Comparable

		BoolVal() bool
		Not() Bool
	}

	// Str is an indexable sequence of runes.
	Str interface {
		fmt.Stringer

		Basic
		Comparable
		Indexable
		Lenable
		Sliceable
		Iterable

		Concat(Str) Str

		Contains(Str) Bool
		Index(Str) Int
		LastIndex(Str) Int
		HasPrefix(Str) Bool
		HasSuffix(Str) Bool
		Replace(Str, Str, Int) Str
		Split(Str) List
		Trim(Str) Str
		ToChars() List
		Map(Eval, StrMapper) (Str, Error)
	}

	// StrMapper transform one string into another
	StrMapper func(Eval, Str) (Str, Error)

	// A Number is either an Int or a Float
	Number interface {
		Basic

		FloatVal() float64
		IntVal() int64

		Add(Number) Number
		Sub(Number) Number
		Mul(Number) Number
		Div(Number) (Number, Error)
		Negate() Number
	}

	// A Float is a float64
	Float interface {
		Number
		Comparable
	}

	// An Int is an int64
	Int interface {
		Number
		Comparable

		Rem(Int) Int
		BitAnd(Int) Int
		BitOr(Int) Int
		BitXOr(Int) Int
		LeftShift(Int) (Int, Error)
		RightShift(Int) (Int, Error)
		Complement() Int
	}
)

//---------------------------------------------------------------
// Composite
//---------------------------------------------------------------

type (
	// Composite is a Value that is composed of other values
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

		Values() []Value

		IsEmpty() Bool
		Contains(Eval, Value) (Bool, Error)
		IndexOf(Eval, Value) (Int, Error)
		Join(Eval, Str) (Str, Error)

		Clear() (List, Error)
		Add(Eval, Value) (List, Error)
		AddAll(Eval, Value) (List, Error)
		Remove(Eval, Int) (List, Error)

		// Sort sorts the list
		Sort(Eval, Lesser) (List, Error)

		// Map creates a new list, leaving the current list unaltered
		Map(Eval, Mapper) (List, Error)

		// Reduce creates a new list, leaving the current list unaltered
		Reduce(Eval, Value, Reducer) (Value, Error)

		// Filter creates a new list, leaving the current list unaltered
		Filter(Eval, Filterer) (List, Error)
	}

	// Lesser returns whether the first param is less than the second param
	Lesser func(Eval, Value, Value) (Bool, Error)

	// Mapper transform one value into another
	Mapper func(Eval, Value) (Value, Error)

	// Reducer combines two values into one
	Reducer func(Eval, Value, Value) (Value, Error)

	// Filterer filters a value
	Filterer func(Eval, Value) (Bool, Error)

	// Range is an immutable, iterable representation of a sequence of integers
	Range interface {
		Composite
		Lenable
		Iterable
		Indexable

		From() Int
		To() Int
		Step() Int
		Count() Int
	}

	// Tuple is an immutable sequence of two or more values
	Tuple interface {
		Composite
		Indexable
		Lenable
	}

	// Dict is an associative array
	Dict interface {
		Composite
		Lenable
		Iterable
		Indexable

		HashMap() *HashMap

		IsEmpty() Bool
		Contains(Eval, Value) (Bool, Error)

		Clear() (Dict, Error)
		AddAll(Eval, Value) (Dict, Error)
		Remove(Eval, Value) (Dict, Error)
	}

	// Set is a set of unique values
	Set interface {
		Composite
		Lenable
		Iterable

		IsEmpty() Bool
		Contains(Eval, Value) (Bool, Error)

		Clear() (Set, Error)
		Add(Eval, Value) (Set, Error)
		AddAll(Eval, Value) (Set, Error)
		Remove(Eval, Value) (Set, Error)
	}

	// Struct is a collection of key-value pairs.
	Struct interface {
		Composite

		SetField(Eval, string, Value) Error

		// Internal is for use only by the Golem Compiler
		Internal(...interface{})
	}

	// Iterator iterates over a sequence of values
	Iterator interface {
		Struct
		IterNext(Eval) (Bool, Error)
		IterGet(Eval) (Value, Error)
	}
)

//---------------------------------------------------------------
// Func
//---------------------------------------------------------------

// Func is a function
type Func interface {
	Value

	Arity() Arity
	Invoke(Eval, []Value) (Value, Error)
}

//---------------------------------------------------------------
// Chan
//---------------------------------------------------------------

// Chan is a goroutine channel
type Chan interface {
	Value
	chanMarker()

	Send(Value)
	Recv() Value
}
