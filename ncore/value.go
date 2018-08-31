// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

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
	ToStr(Evaluator) Str
	Cmp(Evaluator, Value) (Int, Error)

	FieldNames() ([]string, Error)
	HasField(string) (bool, Error)
	GetField(string, Evaluator) (Value, Error)
	InvokeField(string, Evaluator, []Value) (Value, Error)
}

//---------------------------------------------------------------
// Shared Interfaces

type (
	// Indexable is a value that supports the index operator
	Indexable interface {
		Get(Evaluator, Value) (Value, Error)
	}

	// IndexAssignable is a value that supports index assignment
	IndexAssignable interface {
		Indexable
		Set(Evaluator, Value, Value) Error
	}

	// Lenable is a value that has a length
	Lenable interface {
		Len(Evaluator) Int
	}

	// Sliceable is a value that can be sliced
	Sliceable interface {
		Slice(Evaluator, Value, Value) (Value, Error)
		SliceFrom(Evaluator, Value) (Value, Error)
		SliceTo(Evaluator, Value) (Value, Error)
	}

	//// Iterable is a value that can be iterated
	//Iterable interface {
	//	NewIterator(Evaluator) Iterator
	//}
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

		BoolVal() bool
		Not() Bool
	}

	// Str is a string -- defined in golem as a sequence of runes
	Str interface {
		Basic
		fmt.Stringer

		Indexable
		Lenable
		Sliceable
		//Iterable

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
// Func

type (
	// ArityKind defines the various kinds of Arity for a Func
	ArityKind uint16

	// Arity defines the arity of a function
	Arity struct {
		Kind           ArityKind
		RequiredParams uint16
		// For FixedArity and VariadicArity, this value is ignored and should be
		// set to 0.  For MultipleArity, it must be set to a non-zero integer.
		OptionalParams uint16
	}

	// Func is a function
	Func interface {
		Value

		Arity() Arity
		Invoke(Evaluator, []Value) (Value, Error)
	}
)

// The various types of arity
const (
	// FixedArity means a function always takes a fixed number of parameters
	FixedArity ArityKind = iota

	// VariadicArity means that any extra parameters supplied upon invocation will
	// be collected together into a list.
	VariadicArity

	// MultipleArity means that some of the parameters can be omitted, in which case
	// predifined optional values will be substituted.
	MultipleArity
)

func (a Arity) String() string {

	return fmt.Sprintf(
		"Arity(%s,%d,%d)",
		a.Kind.String(),
		a.RequiredParams,
		a.OptionalParams)
}

func (k ArityKind) String() string {

	switch k {
	case FixedArity:
		return "Fixed"
	case VariadicArity:
		return "Variadic"
	case MultipleArity:
		return "Multiple"

	default:
		panic("unreachable")
	}
}

////---------------------------------------------------------------
//// Composite
//
//type (
//	// Composite is a Value that is composed of other values -- List, Range, Tuple,
//	// Dict, Set, Struct
//	Composite interface {
//		Value
//		compositeMarker()
//	}
//
//	// List is an indexable sequence of values
//	List interface {
//		Composite
//		IndexAssignable
//		Lenable
//		Iterable
//		Sliceable
//
//		IsEmpty() Bool
//		Clear() Error
//
//		Contains(Evaluator, Value) (Bool, Error)
//		IndexOf(Evaluator, Value) (Int, Error)
//		Join(Evaluator, Str) Str
//
//		Add(Evaluator, Value) Error
//		AddAll(Evaluator, Value) Error
//		Remove(Evaluator, Int) Error
//
//		Values() []Value
//
//		Map(Evaluator, func(Value) (Value, Error)) (Value, Error)
//		Reduce(Evaluator, Value, func(Value, Value) (Value, Error)) (Value, Error)
//		Filter(Evaluator, func(Value) (Value, Error)) (Value, Error)
//	}
//
//	// Range is an immutable, iterable representation of a  sequence of integers
//	Range interface {
//		Composite
//		Indexable
//		Lenable
//		Iterable
//
//		From() Int
//		To() Int
//		Step() Int
//		Count() Int
//	}
//
//	// Tuple is an immutable sequence of two or more values
//	Tuple interface {
//		Composite
//		Indexable
//		Lenable
//	}
//
//	// Dict is an associative array, a.k.a Hash Map
//	Dict interface {
//		Composite
//		IndexAssignable
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
//
//	// Set is a set of unique values
//	Set interface {
//		Composite
//		Lenable
//		Iterable
//
//		IsEmpty() Bool
//		Clear() Error
//
//		Contains(Evaluator, Value) (Bool, Error)
//		Add(Evaluator, Value) Error
//		AddAll(Evaluator, Value) Error
//		Remove(Evaluator, Value) (Bool, Error)
//	}
//
//	// Struct is a composite collection of names values
//	Struct interface {
//		Composite
//
//		SetField(string, Evaluator, Value) Error
//
// InternalReplace is a 'secret' internal function that is used
// in certain situations by the Interpreter.
// You should never use it yourself, since it will panic if not used properly.
//		InternalReplaceField(string, Field) Error
//	}
//
//	// Iterator iterates over a sequence of values
//	Iterator interface {
//		Struct
//		IterNext() Bool
//		IterGet() (Value, Error)
//	}
//)
//
////---------------------------------------------------------------
//// Chan
//
//// Chan is a goroutine channel
//type Chan interface {
//	Value
//	chanMarker()
//}
