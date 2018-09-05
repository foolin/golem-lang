// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

type (
	// Error is an error
	Error error
)

// NewError creates a new Error
func NewError(s string) Error {
	return fmt.Errorf(s)
}

// NullValueError creates a NullValue Error
func NullValueError() Error {
	return fmt.Errorf("NullValue")
}

// TypeMismatchError creates a TypeMismatch Error
func TypeMismatchError(typ, wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Expected %s, not %s", typ, wrong)
}

// NumberMismatchError creates a TypeMismatch Error
func NumberMismatchError(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Expected Int or Float, not %s", wrong)
}

// IterableMismatchError creates a TypeMismatch Error
func IterableMismatchError(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s has no iter()", wrong)
}

// LenableMismatchError creates a TypeMismatch Error
func LenableMismatchError(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s has no len()", wrong)
}

// IndexableMismatchError creates a TypeMismatch Error
func IndexableMismatchError(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s cannot be indexed", wrong)
}

// SliceableMismatchError creates a TypeMismatch Error
func SliceableMismatchError(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s cannot be sliced", wrong)
}

// ComparableMismatchError creates a TypeMismatch Error
func ComparableMismatchError(a, b Type) Error {
	return fmt.Errorf("TypeMismatch: Types %s and %s cannot be compared", a, b)
}

// HashCodeMismatchError creates a TypeMismatch Error
func HashCodeMismatchError(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Type %s cannot be hashed", wrong)
}

// DivideByZeroError creates a DivideByZero Error
func DivideByZeroError() Error {
	return fmt.Errorf("DivideByZero")
}

// InvalidArgumentError creates an InvalidArgument Error
func InvalidArgumentError(msg string) Error {
	return fmt.Errorf("InvalidArgument: %s", msg)
}

// IndexOutOfBoundsError creates an IndexOutOfBounds Error
func IndexOutOfBoundsError(val int) Error {
	return fmt.Errorf("IndexOutOfBounds: %d", val)
}

// ArityError creates an ArityMismatch Error
func ArityError(expected int, actual int) Error {
	return fmt.Errorf(
		"ArityMismatch: Expected %d params, got %d",
		expected, actual)
}

// ArityAtLeastError creates an ArityMismatch Error
func ArityAtLeastError(expected int, actual int) Error {
	return fmt.Errorf(
		"ArityMismatch: Expected at least %d params, got %d",
		expected, actual)
}

// ArityAtMostError creates an ArityMismatch Error
func ArityAtMostError(expected int, actual int) Error {
	return fmt.Errorf(
		"ArityMismatch: Expected at most %d params, got %d",
		expected, actual)
}

// ReadonlyFieldError creates a ReadonlyField Error
func ReadonlyFieldError(name string) Error {
	return fmt.Errorf("ReadonlyField: Field '%s' is readonly", name)
}

// NoSuchFieldError creates a NoSuchField Error
func NoSuchFieldError(name string) Error {
	return fmt.Errorf("NoSuchField: Field '%s' not found", name)
}

// ImmutableValueError creates an ImmutableValue Error
func ImmutableValueError() Error {
	return fmt.Errorf("ImmutableValue")
}

// InvalidStructKeyError creates a ReadonlyField Error
func InvalidStructKeyError(key string) Error {
	return fmt.Errorf("InvalidStructKey: '%s' is not a valid struct key", key)
}

// UndefinedModuleError creates a UndefinedModule Error
func UndefinedModuleError(name string) Error {
	return fmt.Errorf("UndefinedModule: Module '%s' is not defined", name)
}

// NoSuchElementError creates a NoSuchElement Error
func NoSuchElementError() Error {
	return fmt.Errorf("NoSuchElement")
}

// AssertionFailedError creates a NoSuchElement Error
func AssertionFailedError() Error {
	return fmt.Errorf("AssertionFailed")
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

	stc, e := NewFieldStruct(
		map[string]Field{
			"error":      NewReadonlyField(NewStr(err.Error())),
			"stackTrace": NewReadonlyField(list),
		}, true)
	Assert(e == nil)

	return &errorStruct{stc, err, stackTrace}
}

func (e *errorStruct) Error() Error {
	return e.err
}
func (e *errorStruct) StackTrace() []string {
	return e.stackTrace
}
