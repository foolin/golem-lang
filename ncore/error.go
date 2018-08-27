// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	"fmt"
)

// Error is an error
type Error error

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

func IndexOutOfBoundsError(val int) Error {
	return fmt.Errorf("IndexOutOfBounds: %d", val)
}
