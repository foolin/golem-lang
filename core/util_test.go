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

func TestValidateIndex(t *testing.T) {
	i, err := validateIndex(MakeInt(0), 2)
	ok(t, i, err, MakeInt(0))

	i, err = validateIndex(MakeInt(1), 2)
	ok(t, i, err, MakeInt(1))

	i, err = validateIndex(MakeStr(""), 2)
	fail(t, i, err, "TypeMismatch: Expected 'Int'")

	i, err = validateIndex(MakeInt(-1), 2)
	fail(t, i, err, "IndexOutOfBounds")

	i, err = validateIndex(MakeInt(2), 2)
	fail(t, i, err, "IndexOutOfBounds")
}

func TestValuesEq(t *testing.T) {

	v := valuesEq([]Value{ONE}, []Value{ONE})
	ok(t, v, nil, TRUE)

	v = valuesEq([]Value{}, []Value{})
	ok(t, v, nil, TRUE)

	v = valuesEq([]Value{ONE}, []Value{ZERO})
	ok(t, v, nil, FALSE)

	v = valuesEq([]Value{ONE}, []Value{})
	ok(t, v, nil, FALSE)

	v = valuesEq([]Value{}, []Value{ZERO})
	ok(t, v, nil, FALSE)
}
