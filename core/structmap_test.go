// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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

	sm.put(newField("a", true, Zero))
	tassert(t, sm.size == 1)
	tassert(t, len(sm.buckets) == 5)
	tassert(t, reflect.DeepEqual(sm.fieldNames(), []string{"a"}))

	f, has := sm.get("a")
	tassert(t, has)
	ok(t, f.value, nil, Zero)
	_, has = sm.get("b")
	tassert(t, !has)

	sm.put(newField("b", true, One))
	tassert(t, sm.size == 2)
	tassert(t, len(sm.buckets) == 5)
	tassert(t, reflect.DeepEqual(sm.fieldNames(), []string{"b", "a"}))

	f, has = sm.get("a")
	tassert(t, has)
	ok(t, f.value, nil, Zero)
	f, has = sm.get("b")
	tassert(t, has)
	ok(t, f.value, nil, One)

	sm.put(newField("c", true, NegOne))
	tassert(t, sm.size == 3)
	tassert(t, len(sm.buckets) == 11)
	tassert(t, reflect.DeepEqual(sm.fieldNames(), []string{"b", "a", "c"}))

	f, has = sm.get("c")
	tassert(t, has)
	ok(t, f.value, nil, NegOne)

	sm.put(newField("c", true, Zero))

	f, has = sm.get("c")
	tassert(t, has)
	ok(t, f.value, nil, NegOne)
}
