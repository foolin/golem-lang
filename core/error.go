// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

// Error is an error
type Error error

var (
	NullValueError  = fmt.Errorf("NullValue")
	NoSuchElement   = fmt.Errorf("NoSuchElement")
	AssertionFailed = fmt.Errorf("AssertionFailed")
	ImmutableValue  = fmt.Errorf("ImmutableValue")
	DivideByZero    = fmt.Errorf("DivideByZero")
)

//-----------------------------------------------------------
// miscellaneous
//-----------------------------------------------------------

// InvalidArgument creates an InvalidArgument Error
func InvalidArgument(msg string) Error {
	return fmt.Errorf("InvalidArgument: %s", msg)
}

// IndexOutOfBounds creates an IndexOutOfBounds Error
func IndexOutOfBounds(val int) Error {
	return fmt.Errorf("IndexOutOfBounds: %d", val)
}

// ReadonlyField creates a ReadonlyField Error
func ReadonlyField(name string) Error {
	return fmt.Errorf("ReadonlyField: Field '%s' is readonly", name)
}

// NoSuchField creates a NoSuchField Error
func NoSuchField(name string) Error {
	return fmt.Errorf("NoSuchField: Field '%s' not found", name)
}

// InvalidStructKey creates a InvalidStructKey Error
func InvalidStructKey(key string) Error {
	return fmt.Errorf("InvalidStructKey: '%s' is not a valid struct key", key)
}

// UndefinedModule creates a UndefinedModule Error
func UndefinedModule(name string) Error {
	return fmt.Errorf("UndefinedModule: Module '%s' is not defined", name)
}

//--------------------------------------------------------------
// type mismatch
//--------------------------------------------------------------

// TypeMismatch creates a TypeMismatch Error
func TypeMismatch(typ, wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Expected %s, not %s", typ, wrong)
}

// NumberMismatch creates a TypeMismatch Error
func NumberMismatch(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Expected Int or Float, not %s", wrong)
}

// IterableMismatch creates a TypeMismatch Error
func IterableMismatch(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s has no iter()", wrong)
}

// LenableMismatch creates a TypeMismatch Error
func LenableMismatch(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s has no len()", wrong)
}

// IndexableMismatch creates a TypeMismatch Error
func IndexableMismatch(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s cannot be indexed", wrong)
}

// SliceableMismatch creates a TypeMismatch Error
func SliceableMismatch(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s cannot be sliced", wrong)
}

// ComparableMismatch creates a TypeMismatch Error
func ComparableMismatch(a, b Type) Error {
	return fmt.Errorf("TypeMismatch: Types %s and %s cannot be compared", a, b)
}

// HashCodeMismatch creates a TypeMismatch Error
func HashCodeMismatch(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s cannot be hashed", wrong)
}

//--------------------------------------------------------------
// arity mismatch
//--------------------------------------------------------------

// ArityMismatch creates an ArityMismatch Error
func ArityMismatch(expected int, actual int) Error {
	return fmt.Errorf(
		"ArityMismatch: Expected %d params, got %d",
		expected, actual)
}

// ArityMismatchAtLeast creates an ArityMismatch Error
func ArityMismatchAtLeast(expected int, actual int) Error {
	return fmt.Errorf(
		"ArityMismatch: Expected at least %d params, got %d",
		expected, actual)
}

// ArityMismatchAtMost creates an ArityMismatch Error
func ArityMismatchAtMost(expected int, actual int) Error {
	return fmt.Errorf(
		"ArityMismatch: Expected at most %d params, got %d",
		expected, actual)
}

//--------------------------------------------------------------
// ErrorStruct
//--------------------------------------------------------------

type (
	// ErrorStruct is a Struct that describes an Error
	ErrorStruct interface {
		Struct
		Error() Error
		StackTrace() []string
	}

	errorStruct struct {
		Struct
		err        Error
		stackTrace []string
	}
)

// NewErrorStruct creates a Struct that contains an error and a stack trace.
func NewErrorStruct(err Error, stackTrace []string) ErrorStruct {

	// make list-of-str
	vals := make([]Value, len(stackTrace))
	for i, s := range stackTrace {
		vals[i] = NewStr(s)
	}
	list, e := NewList(vals).Freeze(nil)
	Assert(e == nil)

	stc, e := NewFrozenFieldStruct(
		map[string]Field{
			"error":      NewReadonlyField(NewStr(err.Error())),
			"stackTrace": NewReadonlyField(list),
		})
	Assert(e == nil)

	return &errorStruct{stc, err, stackTrace}
}

func (e *errorStruct) Error() Error {
	return e.err
}
func (e *errorStruct) StackTrace() []string {
	return e.stackTrace
}
