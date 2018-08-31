// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	"fmt"
)

// Error is an error
type Error error

func NewError(s string) Error {
	return fmt.Errorf(s)
}

// NullValueError creates a NullValue Error
func NullValueError() Error {
	return fmt.Errorf("NullValue")
}

// TypeMismatchError creates a TypeMismatch Error
func TypeMismatchError(msg string) Error {
	return fmt.Errorf("TypeMismatch: %s", msg)
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
