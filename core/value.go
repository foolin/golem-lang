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

	Freeze(Eval) (Value, Error)
	Frozen(Eval) (Bool, Error)

	Eq(Eval, Value) (Bool, Error)
	HashCode(Eval) (Int, Error)
	ToStr(Eval) (Str, Error)

	FieldNames() ([]string, Error)
	HasField(string) (bool, Error)
	GetField(string, Eval) (Value, Error)
	InvokeField(string, Eval, []Value) (Value, Error)
}

//---------------------------------------------------------------
// Shared Interfaces

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

type (
	// Basic represents the immutable types null, bool, str, int and float
	Basic interface {
		Value

		basicMarker()
	}

	// Nil is the null value.
	Nil interface {
		Basic
	}

	// Bool is the boolean value -- true or false
	Bool interface {
		Basic
		Comparable

		BoolVal() bool
		Not() Bool
	}

	// Str is a string -- defined in golem as a sequence of runes
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
	}

	// Number is a number
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

	// Float is a float64
	Float interface {
		Number
		Comparable
	}

	// Int is an int64
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
// Func

type (
	// Func is a function
	Func interface {
		Value

		Arity() Arity
		Invoke(Eval, []Value) (Value, Error)
	}
)

//---------------------------------------------------------------
// Composite

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
		Remove(Int) (List, Error)

		Map(Eval, func(Value) (Value, Error)) (List, Error)
		Reduce(Eval, Value, func(Value, Value) (Value, Error)) (Value, Error)
		Filter(Eval, func(Value) (Value, Error)) (List, Error)
	}

	// Range is an immutable, iterable representation of a sequence of integers
	Range interface {
		Composite
		Indexable
		Lenable
		Iterable

		From() Int
		To() Int
		Step() Int
		Count() Int
	}

	// Tuple is an immutable tuple of two or more values
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

	// Struct is a collection of key-value pairs
	Struct interface {
		Composite

		SetField(string, Eval, Value) Error

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
// Chan

// Chan is a goroutine channel
type Chan interface {
	Value
	chanMarker()

	Send(Value)
	Recv() Value
}
