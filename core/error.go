// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

// Error is an error
type Error error

//-----------------------------------------------------------
// miscellaneous
//-----------------------------------------------------------

// NullValueError creates an Error
func NullValueError() Error {
	return fmt.Errorf("NullValue")
}

// AssertionFailed creates an Error
func AssertionFailed() Error {
	return fmt.Errorf("AssertionFailed")
}

// NoSuchElement creates an Error
func NoSuchElement() Error {
	return fmt.Errorf("NoSuchElement")
}

// ImmutableValue creates an Error
func ImmutableValue() Error {
	return fmt.Errorf("ImmutableValue")
}

// DivideByZero creates an Error
func DivideByZero() Error {
	return fmt.Errorf("DivideByZero")
}

// InvalidArgument creates an Error
func InvalidArgument(msg string) Error {
	return fmt.Errorf("InvalidArgument: %s", msg)
}

// IndexOutOfBounds creates an Error
func IndexOutOfBounds(val int) Error {
	return fmt.Errorf("IndexOutOfBounds: %d", val)
}

// ReadonlyField creates an Error
func ReadonlyField(name string) Error {
	return fmt.Errorf("ReadonlyField: Field '%s' is readonly", name)
}

// NoSuchField creates an Error
func NoSuchField(name string) Error {
	return fmt.Errorf("NoSuchField: Field '%s' not found", name)
}

// InvalidStructKey creates an Error
func InvalidStructKey(key string) Error {
	return fmt.Errorf("InvalidStructKey: '%s' is not a valid struct key", key)
}

// UndefinedModule creates an Error
func UndefinedModule(name string) Error {
	return fmt.Errorf("UndefinedModule: Module '%s' is not defined", name)
}

//--------------------------------------------------------------
// type mismatch
//--------------------------------------------------------------

// TypeMismatch creates an Error
func TypeMismatch(typ, wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Expected %s, not %s", typ, wrong)
}

// NumberMismatch creates an Error
func NumberMismatch(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Expected Int or Float, not %s", wrong)
}

// IterableMismatch creates an Error
func IterableMismatch(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s has no iter()", wrong)
}

// LenableMismatch creates an Error
func LenableMismatch(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s has no len()", wrong)
}

// IndexableMismatch creates an Error
func IndexableMismatch(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s cannot be indexed", wrong)
}

// SliceableMismatch creates an Error
func SliceableMismatch(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s cannot be sliced", wrong)
}

// ComparableMismatch creates an Error
func ComparableMismatch(a, b Type) Error {
	return fmt.Errorf("TypeMismatch: Types %s and %s cannot be compared", a, b)
}

// HashCodeMismatch creates an Error
func HashCodeMismatch(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s cannot be hashed", wrong)
}

//--------------------------------------------------------------
// arity mismatch
//--------------------------------------------------------------

// ArityMismatch creates an Error
func ArityMismatch(expected int, actual int) Error {
	return fmt.Errorf(
		"ArityMismatch: Expected %d params, got %d",
		expected, actual)
}

// ArityMismatchAtLeast creates an Error
func ArityMismatchAtLeast(expected int, actual int) Error {
	return fmt.Errorf(
		"ArityMismatch: Expected at least %d params, got %d",
		expected, actual)
}

// ArityMismatchAtMost creates an Error
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
