// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"
)

func TestNative(t *testing.T) {

	a := BuiltinStr
	b := BuiltinLen

	okType(t, a, TFUNC)
	okType(t, b, TFUNC)

	z, err := a.Eq(cx, a)
	ok(t, z, err, TRUE)

	z, err = b.Eq(cx, b)
	ok(t, z, err, TRUE)

	z, err = a.Eq(cx, b)
	ok(t, z, err, FALSE)

	z, err = b.Eq(cx, a)
	ok(t, z, err, FALSE)

	ls := NewList([]Value{ONE, ZERO})

	v, err := a.Invoke(nil, []Value{ls})
	ok(t, v, err, NewStr("[ 1, 0 ]"))

	v, err = b.Invoke(nil, []Value{ls})
	ok(t, v, err, NewInt(2))

}
