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
	//"fmt"
	"testing"
)

func TestBytecodeFunc(t *testing.T) {

	a := NewBytecodeFunc(&Template{})
	b := NewBytecodeFunc(&Template{})

	okType(t, a, TFUNC)
	okType(t, b, TFUNC)

	v, err := a.Eq(cx, a)
	ok(t, v, err, TRUE)

	v, err = b.Eq(cx, b)
	ok(t, v, err, TRUE)

	v, err = a.Eq(cx, b)
	ok(t, v, err, FALSE)

	v, err = b.Eq(cx, a)
	ok(t, v, err, FALSE)
}

func TestLineNumber(t *testing.T) {

	tp := &Template{0, 0, 0, nil,
		[]LineNumberEntry{
			{0, 0},
			{1, 2},
			{11, 3},
			{20, 4},
			{29, 0}},
		nil}

	tassert(t, tp.LineNumber(0) == 0)
	tassert(t, tp.LineNumber(1) == 2)
	tassert(t, tp.LineNumber(10) == 2)
	tassert(t, tp.LineNumber(11) == 3)
	tassert(t, tp.LineNumber(19) == 3)
	tassert(t, tp.LineNumber(20) == 4)
	tassert(t, tp.LineNumber(28) == 4)
	tassert(t, tp.LineNumber(29) == 0)
}
