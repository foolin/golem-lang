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

	Freeze(Evaluator) (Value, Error)
	Frozen(Evaluator) (Bool, Error)

	Eq(Evaluator, Value) (Bool, Error)
	HashCode(Evaluator) (Int, Error)
	ToStr(Evaluator) (Str, Error)

	FieldNames() ([]string, Error)
	HasField(string) (bool, Error)
	GetField(string, Evaluator) (Value, Error)
	InvokeField(string, Evaluator, []Value) (Value, Error)
}

//---------------------------------------------------------------
// Shared Interfaces

type (

	// Comparable is a value that can be compared to another value of the same type.
	Comparable interface {
		Cmp(Evaluator, Comparable) (Int, Error)
	}

	// Indexable is a value that supports the index operator
	Indexable interface {
		Get(Evaluator, Value) (Value, Error)
		Set(Evaluator, Value, Value) Error
	}

	// Lenable is a value that has a length
	Lenable interface {
		Len(Evaluator) (Int, Error)
	}

	// Sliceable is a value that can be sliced
	Sliceable interface {
		Slice(Evaluator, Value, Value) (Value, Error)
		SliceFrom(Evaluator, Value) (Value, Error)
		SliceTo(Evaluator, Value) (Value, Error)
	}

	// Iterable is a value that can be iterated
	Iterable interface {
		NewIterator(Evaluator) (Iterator, Error)
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
		Comparable
	}

	// Int is an int64
	Int interface {
		Number
		Comparable

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
// Func

type (
	// Func is a function
	Func interface {
		Value

		Arity() Arity
		Invoke(Evaluator, []Value) (Value, Error)
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

		Values() []Value

		IsEmpty() Bool
		Clear() Error

		Contains(Evaluator, Value) (Bool, Error)
		IndexOf(Evaluator, Value) (Int, Error)
		Join(Evaluator, Str) (Str, Error)

		Add(Evaluator, Value) Error
		AddAll(Evaluator, Value) Error
		Remove(Evaluator, Int) Error

		Map(Evaluator, func(Value) (Value, Error)) (Value, Error)
		Reduce(Evaluator, Value, func(Value, Value) (Value, Error)) (Value, Error)
		Filter(Evaluator, func(Value) (Value, Error)) (Value, Error)
	}

	// Range is an immutable, iterable representation of a  sequence of integers
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

	// Tuple is an immutable sequence of two or more values
	Tuple interface {
		Composite
		Indexable
		Lenable
	}

	//	// Dict is an associative array, a.k.a Hash Map
	//	Dict interface {
	//		Composite
	//		Lenable
	//		Iterable
	//
	//		IsEmpty() Bool
	//		Clear() Error
	//
	//		ContainsKey(Evaluator, Value) (Bool, Error)
	//		AddAll(Evaluator, Value) Error
	//		Remove(Evaluator, Value) (Bool, Error)
	//	}

	// Set is a set of unique values
	Set interface {
		Composite
		Lenable
		Iterable

		IsEmpty() Bool
		Clear() Error

		Contains(Evaluator, Value) (Bool, Error)
		Add(Evaluator, Value) Error
		AddAll(Evaluator, Value) Error
		Remove(Evaluator, Value) (Bool, Error)
	}

	// Struct is a collection of key-Field pairs
	Struct interface {
		Composite

		SetField(string, Evaluator, Value) Error

		// Internal is for use only by the Golem Compiler
		Internal(...interface{})
	}

	// Iterator iterates over a sequence of values
	Iterator interface {
		Struct
		IterNext(Evaluator) (Bool, Error)
		IterGet(Evaluator) (Value, Error)
	}
)

////---------------------------------------------------------------
//// Chan
//
//// Chan is a goroutine channel
//type Chan interface {
//	Value
//	chanMarker()
//}
