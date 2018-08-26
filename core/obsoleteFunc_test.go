// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"
)

func TestObsolete(t *testing.T) {

	a := BuiltinStr
	b := BuiltinLen

	okType(t, a, FuncType)
	okType(t, b, FuncType)

	z, err := a.Eq(cx, a)
	ok(t, z, err, True)

	z, err = b.Eq(cx, b)
	ok(t, z, err, True)

	z, err = a.Eq(cx, b)
	ok(t, z, err, False)

	z, err = b.Eq(cx, a)
	ok(t, z, err, False)

	ls := NewList([]Value{One, Zero})

	v, err := a.Invoke(nil, []Value{ls})
	ok(t, v, err, NewStr("[ 1, 0 ]"))

	v, err = b.Invoke(nil, []Value{ls})
	ok(t, v, err, NewInt(2))

}
