// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

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

//--------------------------------------------------------------
//--------------------------------------------------------------
//--------------------------------------------------------------

//func TestList(t *testing.T) {
//	ls := NewList([]Value{})
//	okType(t, ls, ListType)
//
//	var v Value
//	var err Error
//
//	v = ls.ToStr(cx)
//	ok(t, v, nil, NewStr("[ ]"))
//
//	v, err = ls.Eq(cx, NewList([]Value{}))
//	ok(t, v, err, True)
//
//	v, err = ls.Eq(cx, NewList([]Value{NewStr("a")}))
//	ok(t, v, err, False)
//
//	v, err = ls.Eq(cx, Null)
//	ok(t, v, err, False)
//
//	v = ls.Len(nil)
//	ok(t, v, nil, Zero)
//
//	err = ls.Add(cx, NewStr("a"))
//	Tassert(t, err == nil)
//
//	v, err = ls.Eq(cx, NewList([]Value{}))
//	ok(t, v, err, False)
//
//	v, err = ls.Eq(cx, NewList([]Value{NewStr("a")}))
//	ok(t, v, err, True)
//
//	v = ls.Len(nil)
//	ok(t, v, nil, One)
//
//	v, err = ls.Get(cx, Zero)
//	ok(t, v, err, NewStr("a"))
//
//	err = ls.Set(cx, Zero, NewStr("b"))
//	Tassert(t, err == nil)
//
//	v, err = ls.Get(cx, Zero)
//	ok(t, v, err, NewStr("b"))
//
//	v, err = ls.Get(cx, NegOne)
//	ok(t, v, err, NewStr("b"))
//
//	v, err = ls.Get(cx, One)
//	fail(t, v, err, "IndexOutOfBounds: 1")
//
//	err = ls.Set(cx, NegOne, True)
//	Tassert(t, err == nil)
//
//	err = ls.Set(cx, One, True)
//	fail(t, nil, err, "IndexOutOfBounds: 1")
//
//	v = ls.ToStr(cx)
//	ok(t, v, nil, NewStr("[ true ]"))
//
//	err = ls.Add(cx, NewStr("z"))
//	Tassert(t, err == nil)
//
//	v = ls.ToStr(cx)
//	ok(t, v, nil, NewStr("[ true, z ]"))
//}
//
//func newDict(cx Context, entries []*HEntry) Dict {
//	dict, err := NewDict(cx, entries)
//	if err != nil {
//		panic(err)
//	}
//	return dict
//}
//
//func TestCompositeHashCode(t *testing.T) {
//	h, err := newDict(cx, []*HEntry{}).HashCode(cx)
//	fail(t, h, err, "TypeMismatch: Expected Hashable Type")
//
//	h, err = NewList([]Value{}).HashCode(cx)
//	fail(t, h, err, "TypeMismatch: Expected Hashable Type")
//
//	h, err = newStruct([]Field{}).HashCode(cx)
//	fail(t, h, err, "TypeMismatch: Expected Hashable Type")
//}
//
//func TestDict(t *testing.T) {
//	d := newDict(cx, []*HEntry{})
//	okType(t, d, DictType)
//
//	var v Value
//	var err Error
//
//	v = d.ToStr(cx)
//	ok(t, v, err, NewStr("dict { }"))
//
//	v, err = d.Eq(cx, newDict(cx, []*HEntry{}))
//	ok(t, v, err, True)
//
//	v, err = d.Eq(cx, Null)
//	ok(t, v, err, False)
//
//	v = d.Len(nil)
//	ok(t, v, nil, Zero)
//
//	v, err = d.Get(cx, NewStr("a"))
//	ok(t, v, err, Null)
//
//	err = d.Set(cx, NewStr("a"), One)
//	Tassert(t, err == nil)
//
//	v, err = d.Get(cx, NewStr("a"))
//	ok(t, v, err, One)
//
//	v, err = d.Eq(cx, newDict(cx, []*HEntry{}))
//	ok(t, v, err, False)
//
//	v, err = d.Eq(cx, newDict(cx, []*HEntry{{NewStr("a"), One}}))
//	ok(t, v, err, True)
//
//	v = d.Len(nil)
//	ok(t, v, nil, One)
//
//	v = d.ToStr(cx)
//	ok(t, v, nil, NewStr("dict { a: 1 }"))
//
//	err = d.Set(cx, NewStr("b"), NewInt(2))
//	Tassert(t, err == nil)
//
//	v, err = d.Get(cx, NewStr("b"))
//	ok(t, v, err, NewInt(2))
//
//	v = d.ToStr(cx)
//	ok(t, v, nil, NewStr("dict { b: 2, a: 1 }"))
//
//	tp := NewTuple([]Value{One, Zero})
//	d = newDict(cx, []*HEntry{{tp, True}})
//
//	v = d.ToStr(cx)
//	ok(t, v, nil, NewStr("dict { (1, 0): true }"))
//
//	v, err = d.Get(cx, tp)
//	ok(t, v, err, True)
//
//	d, err = NewDict(cx, []*HEntry{{Null, True}})
//	fail(t, d, err, "NullValue")
//
//	d, err = NewDict(cx, []*HEntry{{NewList([]Value{}), True}})
//	fail(t, d, err, "TypeMismatch: Expected Hashable Type")
//}
//
//func newSet(cx Context, values []Value) Set {
//	set, err := NewSet(cx, values)
//	if err != nil {
//		panic(err)
//	}
//	return set
//}
//
//func TestSet(t *testing.T) {
//	s := newSet(cx, []Value{})
//	okType(t, s, SetType)
//
//	var v Value
//	var err Error
//
//	v = s.ToStr(cx)
//	ok(t, v, err, NewStr("set { }"))
//
//	v, err = s.Eq(cx, newSet(cx, []Value{}))
//	ok(t, v, err, True)
//
//	v, err = s.Eq(cx, newSet(cx, []Value{One}))
//	ok(t, v, err, False)
//
//	v, err = s.Eq(cx, Null)
//	ok(t, v, err, False)
//
//	v = s.Len(nil)
//	ok(t, v, nil, Zero)
//
//	s = newSet(cx, []Value{One})
//
//	v = s.ToStr(cx)
//	ok(t, v, err, NewStr("set { 1 }"))
//
//	v, err = s.Eq(cx, newSet(cx, []Value{}))
//	ok(t, v, err, False)
//
//	v, err = s.Eq(cx, newSet(cx, []Value{One, One, One}))
//	ok(t, v, err, True)
//
//	v, err = s.Eq(cx, Null)
//	ok(t, v, err, False)
//
//	v = s.Len(nil)
//	ok(t, v, nil, One)
//
//	s = newSet(cx, []Value{One, Zero, Zero, One})
//
//	v = s.ToStr(cx)
//	ok(t, v, err, NewStr("set { 0, 1 }"))
//
//	v = s.Len(nil)
//	ok(t, v, nil, NewInt(2))
//
//	s, err = NewSet(cx, []Value{Null})
//	fail(t, s, err, "NullValue")
//
//	s, err = NewSet(cx, []Value{NewList([]Value{})})
//	fail(t, s, err, "TypeMismatch: Expected Hashable Type")
//}
//
//func TestTuple(t *testing.T) {
//	var v Value
//	var err Error
//
//	tp := NewTuple([]Value{One, Zero})
//	okType(t, tp, TupleType)
//
//	v, err = tp.Eq(cx, NewTuple([]Value{Zero, Zero}))
//	ok(t, v, err, False)
//
//	v, err = tp.Eq(cx, NewTuple([]Value{One, Zero}))
//	ok(t, v, err, True)
//
//	v, err = tp.Eq(cx, Null)
//	ok(t, v, err, False)
//
//	v, err = tp.Get(cx, Zero)
//	ok(t, v, err, One)
//
//	v, err = tp.Get(cx, One)
//	ok(t, v, err, Zero)
//
//	v, err = tp.Get(cx, NegOne)
//	ok(t, v, err, Zero)
//
//	v, err = tp.Get(cx, NewInt(2))
//	fail(t, v, err, "IndexOutOfBounds: 2")
//
//	v = tp.ToStr(cx)
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
//	v, err = r.Eq(cx, newRange(0, 5, 2))
//	ok(t, v, err, False)
//
//	v, err = r.Eq(cx, newRange(0, 5, 1))
//	ok(t, v, err, True)
//
//	v, err = r.Eq(cx, Null)
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
//	v, err = r.Get(cx, One)
//	ok(t, v, err, NewInt(1))
//
//	r = newRange(3, 9, 2)
//	v, err = r.Get(cx, NewInt(2))
//	ok(t, v, err, NewInt(7))
//
//	r = newRange(-9, -13, -1)
//	v, err = r.Get(cx, One)
//	ok(t, v, err, NewInt(-10))
//}
