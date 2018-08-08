// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"strings"
)

// Error represent an error
type Error interface {
	error
	Struct() Struct
}

type serror struct {
	st  Struct
	str string
}

func (e *serror) Error() string {
	return e.str
}

func (e *serror) Struct() Struct {
	return e.st
}

// NewError creates a new error
func NewError(kind string, msg string) Error {

	var st Struct
	var err Error
	var str string

	if msg == "" {
		st, err = NewStruct([]Field{
			NewField("kind", true, NewStr(kind))}, true)
		str = kind
	} else {
		st, err = NewStruct([]Field{
			NewField("kind", true, NewStr(kind)),
			NewField("msg", true, NewStr(msg))}, true)

		str = strings.Join([]string{kind, ": ", msg}, "")
	}
	if err != nil {
		panic("invalid struct")
	}

	return &serror{st, str}
}

// NewErrorFromStruct creates an error from a Struct.  This
// is used in Golem code to do a 'throw'.
func NewErrorFromStruct(cx Context, st Struct) Error {

	var str string
	kind, kerr := st.GetField(cx, NewStr("kind"))
	msg, merr := st.GetField(cx, NewStr("msg"))

	if kerr == nil {
		if merr == nil {
			str = strings.Join([]string{
				kind.ToStr(cx).String(), ": ", msg.ToStr(cx).String()}, "")
		} else {
			str = kind.ToStr(cx).String()
		}
	} else {
		str = st.ToStr(cx).String()
	}

	return &serror{st, str}
}

// NullValueError creates a NullValue Error
func NullValueError() Error {
	return NewError("NullValue", "")
}

// TypeMismatchError creates a TypeMismatch Error
func TypeMismatchError(msg string) Error {
	return NewError("TypeMismatch", msg)
}

// ArityMismatchError creates an ArityMismatch Error
func ArityMismatchError(expected string, actual int) Error {
	return NewError(
		"ArityMismatch",
		fmt.Sprintf("Expected %s params, got %d", expected, actual))
}

// TupleLengthError creates a TupleLength Error
func TupleLengthError(expected int, actual int) Error {
	return NewError(
		"TupleLength",
		fmt.Sprintf("Expected Tuple of length %d, got %d", expected, actual))
}

// DivideByZeroError creates a DivideByZero Error
func DivideByZeroError() Error {
	return NewError("DivideByZero", "")
}

// IndexOutOfBoundsError creates an IndexOutOfBounds Error
func IndexOutOfBoundsError(val int) Error {
	return NewError(
		"IndexOutOfBounds",
		fmt.Sprintf("%d", val))
}

// NoSuchFieldError creates a NoSuchField Error
func NoSuchFieldError(field string) Error {
	return NewError(
		"NoSuchField",
		fmt.Sprintf("Field '%s' not found", field))
}

// ReadonlyFieldError creates a ReadonlyField Error
func ReadonlyFieldError(field string) Error {
	return NewError(
		"ReadonlyField",
		fmt.Sprintf("Field '%s' is readonly", field))
}

// DuplicateFieldError creates a DuplicateField Error
func DuplicateFieldError(field string) Error {
	return NewError(
		"DuplicateField",
		fmt.Sprintf("Field '%s' is a duplicate", field))
}

// InvalidArgumentError creates an InvalidArgument Error
func InvalidArgumentError(msg string) Error {
	return NewError("InvalidArgument", msg)
}

// NoSuchElementError creates a NoSuchElement Error
func NoSuchElementError() Error {
	return NewError("NoSuchElement", "")
}

// AssertionFailedError creates an AssertionFailed Error
func AssertionFailedError() Error {
	return NewError("AssertionFailed", "")
}

// ConstSymbolError creates a ConstSymbol Error
func ConstSymbolError(name string) Error {
	return NewError(
		"ConstSymbol",
		fmt.Sprintf("Symbol '%s' is const", name))
}

// UndefinedSymbolError creates a UndefinedSymbol Error
func UndefinedSymbolError(name string) Error {
	return NewError(
		"UndefinedSymbol",
		fmt.Sprintf("Symbol '%s' is not defined", name))
}

// ImmutableValueError creates an ImmutableValue Error
func ImmutableValueError() Error {
	return NewError("ImmutableValue", "")
}

// UndefinedModuleError creates a UndefinedModule Error
func UndefinedModuleError(name string) Error {
	return NewError(
		"UndefinedModule",
		fmt.Sprintf("Module '%s' is not defined", name))
}

// CouldNotLoadModuleError creates a CouldNotLoadModule Error
func CouldNotLoadModuleError(name string, err error) Error {
	return NewError(
		"CouldNotLoadModule",
		fmt.Sprintf("Could not load Module '%s': %s", name, err.Error()))
}

// PluginError creates a Plugin Error
func PluginError(name string, err error) Error {
	return NewError(
		"Plugin",
		fmt.Sprintf("Plugin '%s': %s", name, err.Error()))
}
