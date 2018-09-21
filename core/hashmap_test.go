// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"reflect"
	"testing"
)

func newHashMap(t *testing.T, entries []*HEntry) *HashMap {
	h, err := NewHashMap(nil, entries)
	tassert(t, err == nil)
	return h
}

func TestHashMap(t *testing.T) {
	hm := newHashMap(t, nil)

	ok(t, hm.Len(), nil, Zero)
	v, err := hm.Get(nil, NewInt(3))
	ok(t, v, err, Null)

	err = hm.Put(nil, NewInt(3), NewInt(33))
	ok(t, nil, err, nil)

	ok(t, hm.Len(), nil, One)
	v, err = hm.Get(nil, NewInt(3))
	ok(t, v, err, NewInt(33))
	v, err = hm.Get(nil, NewInt(5))
	ok(t, v, err, Null)

	err = hm.Put(nil, NewInt(3), NewInt(33))
	ok(t, nil, err, nil)

	ok(t, hm.Len(), nil, One)
	v, err = hm.Get(nil, NewInt(3))
	ok(t, v, err, NewInt(33))
	v, err = hm.Get(nil, NewInt(5))
	ok(t, v, err, Null)

	err = hm.Put(nil, NewInt(int64(2)), NewInt(int64(22)))
	ok(t, nil, err, nil)
	ok(t, hm.Len(), nil, NewInt(2))

	err = hm.Put(nil, NewInt(int64(1)), NewInt(int64(11)))
	ok(t, nil, err, nil)
	ok(t, hm.Len(), nil, NewInt(3))

	for i := 1; i <= 20; i++ {
		err = hm.Put(nil, NewInt(int64(i)), NewInt(int64(i*10+i)))
		ok(t, nil, err, nil)
	}

	for i := 1; i <= 40; i++ {
		v, err = hm.Get(nil, NewInt(int64(i)))
		if i <= 20 {
			ok(t, v, err, NewInt(int64(i*10+i)))
		} else {
			ok(t, v, err, Null)
		}
	}
}

func TestRemove(t *testing.T) {

	d := newHashMap(t, []*HEntry{
		{MustStr("a"), NewInt(1)},
		{MustStr("b"), NewInt(2)}})

	v, err := d.Remove(nil, MustStr("z"))
	ok(t, v, err, False)

	v, err = d.Remove(nil, MustStr("a"))
	ok(t, v, err, True)

	e := newHashMap(t, []*HEntry{
		{MustStr("b"), NewInt(2)}})

	v, err = d.Eq(nil, e)
	ok(t, v, err, True)
}

func TestStrHashMap(t *testing.T) {

	hm := newHashMap(t, nil)

	err := hm.Put(nil, MustStr("abc"), MustStr("xyz"))
	ok(t, nil, err, nil)

	v, err := hm.Get(nil, MustStr("abc"))
	ok(t, v, err, MustStr("xyz"))

	v, err = hm.Contains(nil, MustStr("abc"))
	ok(t, v, err, True)

	v, err = hm.Contains(nil, MustStr("bogus"))
	ok(t, v, err, False)
}

func testIteratorEntries(t *testing.T, initial []*HEntry, expect []*HEntry) {

	hm := newHashMap(t, initial)

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
			{MustStr("a"), NewInt(1)}},
		[]*HEntry{
			{MustStr("a"), NewInt(1)}})

	testIteratorEntries(t,
		[]*HEntry{
			{MustStr("a"), NewInt(1)},
			{MustStr("b"), NewInt(2)}},
		[]*HEntry{
			{MustStr("b"), NewInt(2)},
			{MustStr("a"), NewInt(1)}})

	testIteratorEntries(t,
		[]*HEntry{
			{MustStr("a"), NewInt(1)},
			{MustStr("b"), NewInt(2)},
			{MustStr("c"), NewInt(3)}},
		[]*HEntry{
			{MustStr("b"), NewInt(2)},
			{MustStr("a"), NewInt(1)},
			{MustStr("c"), NewInt(3)}})
}
