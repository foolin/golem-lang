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
	"reflect"
	"testing"
)

func newField(name string, isConst bool, value Value) *field {
	f, _ := NewField(name, isConst, value).(*field)
	return f
}

func TestStructMap(t *testing.T) {
	sm := newStructMap()
	tassert(t, sm.size == 0)
	tassert(t, reflect.DeepEqual(sm.fieldNames(), []string{}))

	_, has := sm.get("a")
	tassert(t, !has)
	_, has = sm.get("b")
	tassert(t, !has)

	sm.put(newField("a", true, ZERO))
	tassert(t, sm.size == 1)
	tassert(t, len(sm.buckets) == 5)
	tassert(t, reflect.DeepEqual(sm.fieldNames(), []string{"a"}))

	f, has := sm.get("a")
	tassert(t, has)
	ok(t, f.value, nil, ZERO)
	_, has = sm.get("b")
	tassert(t, !has)

	sm.put(newField("b", true, ONE))
	tassert(t, sm.size == 2)
	tassert(t, len(sm.buckets) == 5)
	tassert(t, reflect.DeepEqual(sm.fieldNames(), []string{"b", "a"}))

	f, has = sm.get("a")
	tassert(t, has)
	ok(t, f.value, nil, ZERO)
	f, has = sm.get("b")
	tassert(t, has)
	ok(t, f.value, nil, ONE)

	sm.put(newField("c", true, NEG_ONE))
	tassert(t, sm.size == 3)
	tassert(t, len(sm.buckets) == 11)
	tassert(t, reflect.DeepEqual(sm.fieldNames(), []string{"b", "a", "c"}))

	f, has = sm.get("c")
	tassert(t, has)
	ok(t, f.value, nil, NEG_ONE)

	sm.put(newField("c", true, ZERO))

	f, has = sm.get("c")
	tassert(t, has)
	ok(t, f.value, nil, NEG_ONE)
}
