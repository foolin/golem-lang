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
		{NewStr("a"), MakeInt(1)},
		{NewStr("b"), MakeInt(2)}})

	v, err := d.Remove(cx, NewStr("z"))
	ok(t, v, err, FALSE)

	v, err = d.Remove(cx, NewStr("a"))
	ok(t, v, err, TRUE)

	e := NewHashMap(cx, []*HEntry{
		{NewStr("b"), MakeInt(2)}})

	v, err = d.Eq(cx, e)
	ok(t, v, err, TRUE)
}

func TestStrHashMap(t *testing.T) {

	hm := NewHashMap(cx, nil)

	err := hm.Put(cx, NewStr("abc"), NewStr("xyz"))
	ok(t, nil, err, nil)

	v, err := hm.Get(cx, NewStr("abc"))
	ok(t, v, err, NewStr("xyz"))

	v, err = hm.ContainsKey(cx, NewStr("abc"))
	ok(t, v, err, TRUE)

	v, err = hm.ContainsKey(cx, NewStr("bogus"))
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
			{NewStr("a"), MakeInt(1)}},
		[]*HEntry{
			{NewStr("a"), MakeInt(1)}})

	testIteratorEntries(t,
		[]*HEntry{
			{NewStr("a"), MakeInt(1)},
			{NewStr("b"), MakeInt(2)}},
		[]*HEntry{
			{NewStr("b"), MakeInt(2)},
			{NewStr("a"), MakeInt(1)}})

	testIteratorEntries(t,
		[]*HEntry{
			{NewStr("a"), MakeInt(1)},
			{NewStr("b"), MakeInt(2)},
			{NewStr("c"), MakeInt(3)}},
		[]*HEntry{
			{NewStr("b"), MakeInt(2)},
			{NewStr("a"), MakeInt(1)},
			{NewStr("c"), MakeInt(3)}})
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
