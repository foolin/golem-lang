// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"
)

func iok(t *testing.T, val int, err Error, expect int) {

	if err != nil {
		t.Error(err, " != ", nil)
	}

	if val != expect {
		t.Error(val, " != ", expect)
	}
}

func ifail(t *testing.T, err Error, expect string) {
	if err == nil || err.Error() != expect {
		t.Error(err.Error(), " != ", expect)
	}
}

func TestValuesEq(t *testing.T) {

	v, err := valuesEq(cx, []Value{ONE}, []Value{ONE})
	ok(t, v, err, TRUE)

	v, err = valuesEq(cx, []Value{}, []Value{})
	ok(t, v, err, TRUE)

	v, err = valuesEq(cx, []Value{ONE}, []Value{ZERO})
	ok(t, v, err, FALSE)

	v, err = valuesEq(cx, []Value{ONE}, []Value{})
	ok(t, v, err, FALSE)

	v, err = valuesEq(cx, []Value{}, []Value{ZERO})
	ok(t, v, err, FALSE)
}

func TestIndex(t *testing.T) {

	n, err := posIndex(ZERO, 3)
	iok(t, n, err, 0)

	n, err = posIndex(MakeInt(3), 3)
	iok(t, n, err, 3)

	n, err = posIndex(NEG_ONE, 3)
	iok(t, n, err, 2)

	n, err = posIndex(MakeInt(-3), 3)
	iok(t, n, err, 0)

	_, err = posIndex(NewStr(""), 3)
	ifail(t, err, "TypeMismatch: Expected 'Int'")

	//--------------------------------------

	n, err = boundedIndex(ZERO, 3)
	iok(t, n, err, 0)

	n, err = boundedIndex(NEG_ONE, 3)
	iok(t, n, err, 2)

	n, err = boundedIndex(MakeInt(-3), 3)
	iok(t, n, err, 0)

	_, err = boundedIndex(NewStr(""), 3)
	ifail(t, err, "TypeMismatch: Expected 'Int'")

	_, err = boundedIndex(MakeInt(-4), 3)
	ifail(t, err, "IndexOutOfBounds: -1")

	n, err = boundedIndex(MakeInt(3), 3)
	ifail(t, err, "IndexOutOfBounds: 3")

	//--------------------------------------

	a, b, err := sliceIndices(ZERO, NEG_ONE, 3)
	iok(t, a, err, 0)
	iok(t, b, err, 2)

	a, b, err = sliceIndices(MakeInt(-3), MakeInt(3), 3)
	iok(t, a, err, 0)
	iok(t, b, err, 3)

	_, _, err = sliceIndices(NewStr(""), ZERO, 3)
	ifail(t, err, "TypeMismatch: Expected 'Int'")

	_, _, err = sliceIndices(ZERO, NewStr(""), 3)
	ifail(t, err, "TypeMismatch: Expected 'Int'")
}
