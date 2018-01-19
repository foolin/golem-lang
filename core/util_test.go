// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	_, err = posIndex(MakeStr(""), 3)
	ifail(t, err, "TypeMismatch: Expected 'Int'")

	//--------------------------------------

	n, err = boundedIndex(ZERO, 3)
	iok(t, n, err, 0)

	n, err = boundedIndex(NEG_ONE, 3)
	iok(t, n, err, 2)

	n, err = boundedIndex(MakeInt(-3), 3)
	iok(t, n, err, 0)

	_, err = boundedIndex(MakeStr(""), 3)
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

	_, _, err = sliceIndices(MakeStr(""), ZERO, 3)
	ifail(t, err, "TypeMismatch: Expected 'Int'")

	_, _, err = sliceIndices(ZERO, MakeStr(""), 3)
	ifail(t, err, "TypeMismatch: Expected 'Int'")
}
