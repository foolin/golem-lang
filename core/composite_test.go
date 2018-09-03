// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"
)

func TestStruct(t *testing.T) {

	var val Value
	var err Error

	fs, err := NewFieldStruct(
		map[string]Field{
			"goto": NewField(NewInt(0)),
		}, true)
	fail(t, fs, err, "InvalidStructKey: 'goto' is not a valid struct key.")

	fs, err = NewFieldStruct(
		map[string]Field{
			"a": NewField(NewInt(0)),
			"b": NewField(NewInt(1)),
		}, true)
	Tassert(t, err == nil)

	val, err = fs.ToStr(nil)
	Tassert(t, err == nil)
	println(val.(Str).String())

	method := NewFixedMethod(
		[]Type{},
		false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			return Null, nil
		})

	vs, err := NewVirtualStruct(
		map[string]Method{
			" ": method,
		}, true)
	fail(t, vs, err, "InvalidStructKey: ' ' is not a valid struct key.")

	vs, err = NewVirtualStruct(
		map[string]Method{
			"a": method,
		}, true)
	Tassert(t, err == nil)

	val, err = vs.ToStr(nil)
	Tassert(t, err == nil)
	println(val.(Str).String())
}

func TestList(t *testing.T) {
	ls := NewList([]Value{})
	okType(t, ls, ListType)

	var v Value
	var err Error

	v, err = ls.ToStr(nil)
	ok(t, v, err, NewStr("[ ]"))

	v, err = ls.Eq(nil, NewList([]Value{}))
	ok(t, v, err, True)

	v, err = ls.Eq(nil, NewList([]Value{NewStr("a")}))
	ok(t, v, err, False)

	v, err = ls.Eq(nil, Null)
	ok(t, v, err, False)

	v, err = ls.Len(nil)
	ok(t, v, err, Zero)

	err = ls.Add(nil, NewStr("a"))
	Tassert(t, err == nil)

	v, err = ls.Eq(nil, NewList([]Value{}))
	ok(t, v, err, False)

	v, err = ls.Eq(nil, NewList([]Value{NewStr("a")}))
	ok(t, v, err, True)

	v, err = ls.Len(nil)
	ok(t, v, err, One)

	v, err = ls.Get(nil, Zero)
	ok(t, v, err, NewStr("a"))

	err = ls.Set(nil, Zero, NewStr("b"))
	Tassert(t, err == nil)

	v, err = ls.Get(nil, Zero)
	ok(t, v, err, NewStr("b"))

	v, err = ls.Get(nil, NegOne)
	ok(t, v, err, NewStr("b"))

	v, err = ls.Get(nil, One)
	fail(t, v, err, "IndexOutOfBounds: 1")

	err = ls.Set(nil, NegOne, True)
	Tassert(t, err == nil)

	err = ls.Set(nil, One, True)
	fail(t, nil, err, "IndexOutOfBounds: 1")

	v, err = ls.ToStr(nil)
	ok(t, v, err, NewStr("[ true ]"))

	err = ls.Add(nil, NewStr("z"))
	Tassert(t, err == nil)

	v, err = ls.ToStr(nil)
	ok(t, v, err, NewStr("[ true, z ]"))
}

//--------------------------------------------------------------
//--------------------------------------------------------------
//--------------------------------------------------------------

//
//func newDict(nil Context, entries []*HEntry) Dict {
//	dict, err := NewDict(nil, entries)
//	if err != nil {
//		panic(err)
//	}
//	return dict
//}
//
//func TestCompositeHashCode(t *testing.T) {
//	h, err := newDict(nil, []*HEntry{}).HashCode(nil)
//	fail(t, h, err, "TypeMismatch: Expected Hashable")
//
//	h, err = NewList([]Value{}).HashCode(nil)
//	fail(t, h, err, "TypeMismatch: Expected Hashable")
//
//	h, err = newStruct([]Field{}).HashCode(nil)
//	fail(t, h, err, "TypeMismatch: Expected Hashable")
//}
//
//func TestDict(t *testing.T) {
//	d := newDict(nil, []*HEntry{})
//	okType(t, d, DictType)
//
//	var v Value
//	var err Error
//
//	v = d.ToStr(nil)
//	ok(t, v, err, NewStr("dict { }"))
//
//	v, err = d.Eq(nil, newDict(nil, []*HEntry{}))
//	ok(t, v, err, True)
//
//	v, err = d.Eq(nil, Null)
//	ok(t, v, err, False)
//
//	v = d.Len(nil)
//	ok(t, v, nil, Zero)
//
//	v, err = d.Get(nil, NewStr("a"))
//	ok(t, v, err, Null)
//
//	err = d.Set(nil, NewStr("a"), One)
//	Tassert(t, err == nil)
//
//	v, err = d.Get(nil, NewStr("a"))
//	ok(t, v, err, One)
//
//	v, err = d.Eq(nil, newDict(nil, []*HEntry{}))
//	ok(t, v, err, False)
//
//	v, err = d.Eq(nil, newDict(nil, []*HEntry{{NewStr("a"), One}}))
//	ok(t, v, err, True)
//
//	v = d.Len(nil)
//	ok(t, v, nil, One)
//
//	v = d.ToStr(nil)
//	ok(t, v, nil, NewStr("dict { a: 1 }"))
//
//	err = d.Set(nil, NewStr("b"), NewInt(2))
//	Tassert(t, err == nil)
//
//	v, err = d.Get(nil, NewStr("b"))
//	ok(t, v, err, NewInt(2))
//
//	v = d.ToStr(nil)
//	ok(t, v, nil, NewStr("dict { b: 2, a: 1 }"))
//
//	tp := NewTuple([]Value{One, Zero})
//	d = newDict(nil, []*HEntry{{tp, True}})
//
//	v = d.ToStr(nil)
//	ok(t, v, nil, NewStr("dict { (1, 0): true }"))
//
//	v, err = d.Get(nil, tp)
//	ok(t, v, err, True)
//
//	d, err = NewDict(nil, []*HEntry{{Null, True}})
//	fail(t, d, err, "NullValue")
//
//	d, err = NewDict(nil, []*HEntry{{NewList([]Value{}), True}})
//	fail(t, d, err, "TypeMismatch: Expected Hashable")
//}
//
//func newSet(nil Context, values []Value) Set {
//	set, err := NewSet(nil, values)
//	if err != nil {
//		panic(err)
//	}
//	return set
//}
//
//func TestSet(t *testing.T) {
//	s := newSet(nil, []Value{})
//	okType(t, s, SetType)
//
//	var v Value
//	var err Error
//
//	v = s.ToStr(nil)
//	ok(t, v, err, NewStr("set { }"))
//
//	v, err = s.Eq(nil, newSet(nil, []Value{}))
//	ok(t, v, err, True)
//
//	v, err = s.Eq(nil, newSet(nil, []Value{One}))
//	ok(t, v, err, False)
//
//	v, err = s.Eq(nil, Null)
//	ok(t, v, err, False)
//
//	v = s.Len(nil)
//	ok(t, v, nil, Zero)
//
//	s = newSet(nil, []Value{One})
//
//	v = s.ToStr(nil)
//	ok(t, v, err, NewStr("set { 1 }"))
//
//	v, err = s.Eq(nil, newSet(nil, []Value{}))
//	ok(t, v, err, False)
//
//	v, err = s.Eq(nil, newSet(nil, []Value{One, One, One}))
//	ok(t, v, err, True)
//
//	v, err = s.Eq(nil, Null)
//	ok(t, v, err, False)
//
//	v = s.Len(nil)
//	ok(t, v, nil, One)
//
//	s = newSet(nil, []Value{One, Zero, Zero, One})
//
//	v = s.ToStr(nil)
//	ok(t, v, err, NewStr("set { 0, 1 }"))
//
//	v = s.Len(nil)
//	ok(t, v, nil, NewInt(2))
//
//	s, err = NewSet(nil, []Value{Null})
//	fail(t, s, err, "NullValue")
//
//	s, err = NewSet(nil, []Value{NewList([]Value{})})
//	fail(t, s, err, "TypeMismatch: Expected Hashable")
//}
//
//func TestTuple(t *testing.T) {
//	var v Value
//	var err Error
//
//	tp := NewTuple([]Value{One, Zero})
//	okType(t, tp, TupleType)
//
//	v, err = tp.Eq(nil, NewTuple([]Value{Zero, Zero}))
//	ok(t, v, err, False)
//
//	v, err = tp.Eq(nil, NewTuple([]Value{One, Zero}))
//	ok(t, v, err, True)
//
//	v, err = tp.Eq(nil, Null)
//	ok(t, v, err, False)
//
//	v, err = tp.Get(nil, Zero)
//	ok(t, v, err, One)
//
//	v, err = tp.Get(nil, One)
//	ok(t, v, err, Zero)
//
//	v, err = tp.Get(nil, NegOne)
//	ok(t, v, err, Zero)
//
//	v, err = tp.Get(nil, NewInt(2))
//	fail(t, v, err, "IndexOutOfBounds: 2")
//
//	v = tp.ToStr(nil)
//	ok(t, v, nil, NewStr("(1, 0)"))
//
//	v = tp.Len(nil)
//	ok(t, v, nil, NewInt(2))
//}
//
//func newRange(from int64, to int64, step int64) Range {
//	r, err := NewRange(from, to, step)
//	if err != nil {
//		panic("invalid range")
//	}
//	return r
//}
//
//func TestRange(t *testing.T) {
//	var v Value
//	var err Error
//
//	r := newRange(0, 5, 1)
//	okType(t, r, RangeType)
//
//	v, err = r.Eq(nil, newRange(0, 5, 2))
//	ok(t, v, err, False)
//
//	v, err = r.Eq(nil, newRange(0, 5, 1))
//	ok(t, v, err, True)
//
//	v, err = r.Eq(nil, Null)
//	ok(t, v, err, False)
//
//	v = r.Len(nil)
//	ok(t, v, nil, NewInt(5))
//
//	v = newRange(0, 6, 3).Len(nil)
//	ok(t, v, nil, NewInt(2))
//	v = newRange(0, 7, 3).Len(nil)
//	ok(t, v, nil, NewInt(3))
//	v = newRange(0, 8, 3).Len(nil)
//	ok(t, v, nil, NewInt(3))
//	v = newRange(0, 9, 3).Len(nil)
//	ok(t, v, nil, NewInt(3))
//
//	v = newRange(0, 0, 3).Len(nil)
//	ok(t, v, nil, NewInt(0))
//	v = newRange(1, 0, 1).Len(nil)
//	ok(t, v, nil, NewInt(0))
//
//	v, err = NewRange(1, 0, 0)
//	fail(t, v, err, "InvalidArgument: step cannot be 0")
//
//	v = newRange(0, -5, -1).Len(nil)
//	ok(t, v, nil, NewInt(5))
//	v = newRange(-1, -8, -3).Len(nil)
//	ok(t, v, nil, NewInt(3))
//
//	r = newRange(0, 5, 1)
//	v, err = r.Get(nil, One)
//	ok(t, v, err, NewInt(1))
//
//	r = newRange(3, 9, 2)
//	v, err = r.Get(nil, NewInt(2))
//	ok(t, v, err, NewInt(7))
//
//	r = newRange(-9, -13, -1)
//	v, err = r.Get(nil, One)
//	ok(t, v, err, NewInt(-10))
//}
