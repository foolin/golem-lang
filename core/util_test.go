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

	v, err := valuesEq(cx, []Value{One}, []Value{One})
	ok(t, v, err, True)

	v, err = valuesEq(cx, []Value{}, []Value{})
	ok(t, v, err, True)

	v, err = valuesEq(cx, []Value{One}, []Value{Zero})
	ok(t, v, err, False)

	v, err = valuesEq(cx, []Value{One}, []Value{})
	ok(t, v, err, False)

	v, err = valuesEq(cx, []Value{}, []Value{Zero})
	ok(t, v, err, False)
}

func TestIndex(t *testing.T) {

	n, err := posIndex(Zero, 3)
	iok(t, n, err, 0)

	n, err = posIndex(NewInt(3), 3)
	iok(t, n, err, 3)

	n, err = posIndex(NegOne, 3)
	iok(t, n, err, 2)

	n, err = posIndex(NewInt(-3), 3)
	iok(t, n, err, 0)

	_, err = posIndex(NewStr(""), 3)
	ifail(t, err, "TypeMismatch: Expected Int")

	//--------------------------------------

	n, err = boundedIndex(Zero, 3)
	iok(t, n, err, 0)

	n, err = boundedIndex(NegOne, 3)
	iok(t, n, err, 2)

	n, err = boundedIndex(NewInt(-3), 3)
	iok(t, n, err, 0)

	_, err = boundedIndex(NewStr(""), 3)
	ifail(t, err, "TypeMismatch: Expected Int")

	_, err = boundedIndex(NewInt(-4), 3)
	ifail(t, err, "IndexOutOfBounds: -1")

	_, err = boundedIndex(NewInt(3), 3)
	ifail(t, err, "IndexOutOfBounds: 3")

	//--------------------------------------

	a, b, err := sliceIndices(Zero, NegOne, 3)
	iok(t, a, err, 0)
	iok(t, b, err, 2)

	a, b, err = sliceIndices(NewInt(-3), NewInt(3), 3)
	iok(t, a, err, 0)
	iok(t, b, err, 3)

	_, _, err = sliceIndices(NewStr(""), Zero, 3)
	ifail(t, err, "TypeMismatch: Expected Int")

	_, _, err = sliceIndices(Zero, NewStr(""), 3)
	ifail(t, err, "TypeMismatch: Expected Int")
}
