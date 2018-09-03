// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"fmt"
)

type (
	// Error is an error
	Error error
)

func NewError(s string) Error {
	return fmt.Errorf(s)
}

// NullValueError creates a NullValue Error
func NullValueError() Error {
	return fmt.Errorf("NullValue")
}

// TempMismatchError creates a TypeMismatch Error
func TempMismatchError(msg string) Error {
	return fmt.Errorf("TypeMismatch: %s", msg)
}

// TypeMismatchError creates a TypeMismatch Error
func TypeMismatchError(typ Type /*, types ...Type*/) Error {

	var buf bytes.Buffer

	buf.WriteString(typ.String())
	//for _, t := range types {
	//	buf.WriteString(" or ")
	//	buf.WriteString(t.String())
	//}

	return fmt.Errorf("TypeMismatch: Expected %s", buf.String())
}

// NumberMismatchError creates a TypeMismatch Error
func NumberMismatchError(wrong Type) Error {
	return fmt.Errorf("TypeMismatch: Expected Int or Float, not %s", wrong)
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

// ArityMismatchError creates an ArityMismatch Error
func ArityMismatchError(expected string, actual int) Error {
	return fmt.Errorf(
		"ArityMismatch: Expected %s params, got %d",
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
	return fmt.Errorf("InvalidStructKey: '%s' is not a valid struct key.", key)
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
