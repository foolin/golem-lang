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
