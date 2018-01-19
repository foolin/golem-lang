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
	"reflect"
	"testing"
)

func TestHashMap(t *testing.T) {
	hm := NewHashMap(cx, nil)

	ok(t, hm.Len(), nil, ZERO)
	v, err := hm.Get(cx, MakeInt(3))
	ok(t, v, err, NULL)

	err = hm.Put(cx, MakeInt(3), MakeInt(33))
	ok(t, nil, err, nil)

	ok(t, hm.Len(), nil, ONE)
	v, err = hm.Get(cx, MakeInt(3))
	ok(t, v, err, MakeInt(33))
	v, err = hm.Get(cx, MakeInt(5))
	ok(t, v, err, NULL)

	err = hm.Put(cx, MakeInt(3), MakeInt(33))
	ok(t, nil, err, nil)

	ok(t, hm.Len(), nil, ONE)
	v, err = hm.Get(cx, MakeInt(3))
	ok(t, v, err, MakeInt(33))
	v, err = hm.Get(cx, MakeInt(5))
	ok(t, v, err, NULL)

	err = hm.Put(cx, MakeInt(int64(2)), MakeInt(int64(22)))
	ok(t, nil, err, nil)
	ok(t, hm.Len(), nil, MakeInt(2))

	err = hm.Put(cx, MakeInt(int64(1)), MakeInt(int64(11)))
	ok(t, nil, err, nil)
	ok(t, hm.Len(), nil, MakeInt(3))

	for i := 1; i <= 20; i++ {
		err = hm.Put(cx, MakeInt(int64(i)), MakeInt(int64(i*10+i)))
		ok(t, nil, err, nil)
	}

	for i := 1; i <= 40; i++ {
		v, err = hm.Get(cx, MakeInt(int64(i)))
		if i <= 20 {
			ok(t, v, err, MakeInt(int64(i*10+i)))
		} else {
			ok(t, v, err, NULL)
		}
	}
}

func TestRemove(t *testing.T) {

	d := NewHashMap(cx, []*HEntry{
		{MakeStr("a"), MakeInt(1)},
		{MakeStr("b"), MakeInt(2)}})

	v, err := d.Remove(cx, MakeStr("z"))
	ok(t, v, err, FALSE)

	v, err = d.Remove(cx, MakeStr("a"))
	ok(t, v, err, TRUE)

	e := NewHashMap(cx, []*HEntry{
		{MakeStr("b"), MakeInt(2)}})

	v, err = d.Eq(cx, e)
	ok(t, v, err, TRUE)
}

func TestStrHashMap(t *testing.T) {

	hm := NewHashMap(cx, nil)

	err := hm.Put(cx, MakeStr("abc"), MakeStr("xyz"))
	ok(t, nil, err, nil)

	v, err := hm.Get(cx, MakeStr("abc"))
	ok(t, v, err, MakeStr("xyz"))

	v, err = hm.ContainsKey(cx, MakeStr("abc"))
	ok(t, v, err, TRUE)

	v, err = hm.ContainsKey(cx, MakeStr("bogus"))
	ok(t, v, err, FALSE)
}

func testIteratorEntries(t *testing.T, initial []*HEntry, expect []*HEntry) {

	hm := NewHashMap(cx, initial)

	entries := []*HEntry{}
	itr := hm.Iterator()
	for itr.Next() {
		entries = append(entries, itr.Get())
	}

	if !reflect.DeepEqual(entries, expect) {
		t.Error("iterator failed")
	}
}

func TestHashMapIterator(t *testing.T) {

	testIteratorEntries(t,
		[]*HEntry{},
		[]*HEntry{})

	testIteratorEntries(t,
		[]*HEntry{
			{MakeStr("a"), MakeInt(1)}},
		[]*HEntry{
			{MakeStr("a"), MakeInt(1)}})

	testIteratorEntries(t,
		[]*HEntry{
			{MakeStr("a"), MakeInt(1)},
			{MakeStr("b"), MakeInt(2)}},
		[]*HEntry{
			{MakeStr("b"), MakeInt(2)},
			{MakeStr("a"), MakeInt(1)}})

	testIteratorEntries(t,
		[]*HEntry{
			{MakeStr("a"), MakeInt(1)},
			{MakeStr("b"), MakeInt(2)},
			{MakeStr("c"), MakeInt(3)}},
		[]*HEntry{
			{MakeStr("b"), MakeInt(2)},
			{MakeStr("a"), MakeInt(1)},
			{MakeStr("c"), MakeInt(3)}})
}

func TestBogusHashCode(t *testing.T) {

	key := NewList([]Value{})
	var v Value
	var err Error

	hm := NewHashMap(cx, nil)
	v, err = hm.Get(cx, key)
	fail(t, v, err, "TypeMismatch: Expected Hashable Type")

	v, err = hm.ContainsKey(cx, key)
	fail(t, v, err, "TypeMismatch: Expected Hashable Type")

	err = hm.Put(cx, key, ZERO)
	fail(t, nil, err, "TypeMismatch: Expected Hashable Type")
}
