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

func TestStructMap(t *testing.T) {
	sm := newStructMap()
	assert(t, sm.size == 0)
	assert(t, reflect.DeepEqual(sm.keys(), []string{}))

	_, has := sm.get("a")
	assert(t, !has)
	_, has = sm.get("b")
	assert(t, !has)

	sm.put(&StructEntry{"a", true, false, ZERO})
	assert(t, sm.size == 1)
	assert(t, len(sm.buckets) == 5)
	assert(t, reflect.DeepEqual(sm.keys(), []string{"a"}))

	e, has := sm.get("a")
	assert(t, has)
	ok(t, e.Value, nil, ZERO)
	_, has = sm.get("b")
	assert(t, !has)

	sm.put(&StructEntry{"b", true, false, ONE})
	assert(t, sm.size == 2)
	assert(t, len(sm.buckets) == 5)
	assert(t, reflect.DeepEqual(sm.keys(), []string{"b", "a"}))

	e, has = sm.get("a")
	assert(t, has)
	ok(t, e.Value, nil, ZERO)
	e, has = sm.get("b")
	assert(t, has)
	ok(t, e.Value, nil, ONE)

	sm.put(&StructEntry{"c", true, false, NEG_ONE})
	assert(t, sm.size == 3)
	assert(t, len(sm.buckets) == 11)
	assert(t, reflect.DeepEqual(sm.keys(), []string{"b", "a", "c"}))

	e, has = sm.get("c")
	assert(t, has)
	ok(t, e.Value, nil, NEG_ONE)

	sm.put(&StructEntry{"c", true, false, ZERO})

	e, has = sm.get("c")
	assert(t, has)
	ok(t, e.Value, nil, NEG_ONE)
}
