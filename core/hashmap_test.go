// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"reflect"
	"testing"
)

func TestHashMap(t *testing.T) {
	hm := NewHashMap(cx, nil)

	ok(t, hm.Len(), nil, Zero)
	v, err := hm.Get(cx, NewInt(3))
	ok(t, v, err, NullValue)

	err = hm.Put(cx, NewInt(3), NewInt(33))
	ok(t, nil, err, nil)

	ok(t, hm.Len(), nil, One)
	v, err = hm.Get(cx, NewInt(3))
	ok(t, v, err, NewInt(33))
	v, err = hm.Get(cx, NewInt(5))
	ok(t, v, err, NullValue)

	err = hm.Put(cx, NewInt(3), NewInt(33))
	ok(t, nil, err, nil)

	ok(t, hm.Len(), nil, One)
	v, err = hm.Get(cx, NewInt(3))
	ok(t, v, err, NewInt(33))
	v, err = hm.Get(cx, NewInt(5))
	ok(t, v, err, NullValue)

	err = hm.Put(cx, NewInt(int64(2)), NewInt(int64(22)))
	ok(t, nil, err, nil)
	ok(t, hm.Len(), nil, NewInt(2))

	err = hm.Put(cx, NewInt(int64(1)), NewInt(int64(11)))
	ok(t, nil, err, nil)
	ok(t, hm.Len(), nil, NewInt(3))

	for i := 1; i <= 20; i++ {
		err = hm.Put(cx, NewInt(int64(i)), NewInt(int64(i*10+i)))
		ok(t, nil, err, nil)
	}

	for i := 1; i <= 40; i++ {
		v, err = hm.Get(cx, NewInt(int64(i)))
		if i <= 20 {
			ok(t, v, err, NewInt(int64(i*10+i)))
		} else {
			ok(t, v, err, NullValue)
		}
	}
}

func TestRemove(t *testing.T) {

	d := NewHashMap(cx, []*HEntry{
		{NewStr("a"), NewInt(1)},
		{NewStr("b"), NewInt(2)}})

	v, err := d.Remove(cx, NewStr("z"))
	ok(t, v, err, False)

	v, err = d.Remove(cx, NewStr("a"))
	ok(t, v, err, True)

	e := NewHashMap(cx, []*HEntry{
		{NewStr("b"), NewInt(2)}})

	v, err = d.Eq(cx, e)
	ok(t, v, err, True)
}

func TestStrHashMap(t *testing.T) {

	hm := NewHashMap(cx, nil)

	err := hm.Put(cx, NewStr("abc"), NewStr("xyz"))
	ok(t, nil, err, nil)

	v, err := hm.Get(cx, NewStr("abc"))
	ok(t, v, err, NewStr("xyz"))

	v, err = hm.ContainsKey(cx, NewStr("abc"))
	ok(t, v, err, True)

	v, err = hm.ContainsKey(cx, NewStr("bogus"))
	ok(t, v, err, False)
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
			{NewStr("a"), NewInt(1)}},
		[]*HEntry{
			{NewStr("a"), NewInt(1)}})

	testIteratorEntries(t,
		[]*HEntry{
			{NewStr("a"), NewInt(1)},
			{NewStr("b"), NewInt(2)}},
		[]*HEntry{
			{NewStr("b"), NewInt(2)},
			{NewStr("a"), NewInt(1)}})

	testIteratorEntries(t,
		[]*HEntry{
			{NewStr("a"), NewInt(1)},
			{NewStr("b"), NewInt(2)},
			{NewStr("c"), NewInt(3)}},
		[]*HEntry{
			{NewStr("b"), NewInt(2)},
			{NewStr("a"), NewInt(1)},
			{NewStr("c"), NewInt(3)}})
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

	err = hm.Put(cx, key, Zero)
	fail(t, nil, err, "TypeMismatch: Expected Hashable Type")
}
