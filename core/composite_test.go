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

func TestChain(t *testing.T) {
	s0 := newStruct([]Field{
		NewField("a", true, ONE),
		NewField("b", true, MakeInt(2))})
	s1 := newStruct([]Field{
		NewField("b", true, MakeInt(3)),
		NewField("c", true, MakeInt(4))})

	c := MergeStructs([]Struct{s0, s1})
	//println(c.ToStr(cx).String())
	tassert(t, len(c.FieldNames()) == 3)
	//ok(t, c.ToStr(cx), nil, MakeStr("struct { b: 2, c: 4, a: 1 }"))

	v, err := c.GetField(cx, MakeStr("a"))
	ok(t, v, err, ONE)
	v, err = c.GetField(cx, MakeStr("b"))
	ok(t, v, err, MakeInt(2))
	v, err = c.GetField(cx, MakeStr("c"))
	ok(t, v, err, MakeInt(4))
}

func newStruct(fields []Field) Struct {
	stc, err := NewStruct(fields, false)
	if err != nil {
		panic("invalid struct")
	}
	return stc
}

func TestStruct(t *testing.T) {
	stc := newStruct([]Field{})
	okType(t, stc, TSTRUCT)
	tassert(t, reflect.DeepEqual(stc.FieldNames(), []string{}))

	s := stc.ToStr(cx)
	ok(t, s, nil, MakeStr("struct { }"))

	z, err := stc.Eq(cx, newStruct([]Field{}))
	ok(t, z, err, TRUE)
	z, err = stc.Eq(cx, newStruct([]Field{NewField("a", true, ONE)}))
	ok(t, z, err, FALSE)

	val, err := stc.GetField(cx, MakeStr("a"))
	fail(t, val, err, "NoSuchField: Field 'a' not found")

	//////////////////

	stc = newStruct([]Field{NewField("a", false, ONE)})
	okType(t, stc, TSTRUCT)

	s = stc.ToStr(cx)
	ok(t, s, nil, MakeStr("struct { a: 1 }"))

	z, err = stc.Eq(cx, newStruct([]Field{}))
	ok(t, z, err, FALSE)
	z, err = stc.Eq(cx, newStruct([]Field{NewField("a", true, ONE)}))
	ok(t, z, err, TRUE)

	val, err = stc.GetField(cx, MakeStr("a"))
	ok(t, val, err, ONE)

	val, err = stc.GetField(cx, MakeStr("b"))
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	err = stc.SetField(cx, MakeStr("a"), MakeInt(123))
	if err != nil {
		panic("unexpected error")
	}

	val, err = stc.GetField(cx, MakeStr("a"))
	ok(t, val, err, MakeInt(123))

	err = stc.SetField(cx, MakeStr("a"), MakeInt(456))
	if err != nil {
		panic("unexpected error")
	}

	val, err = stc.GetField(cx, MakeStr("a"))
	ok(t, val, err, MakeInt(456))

	val, err = stc.Has(MakeStr("a"))
	ok(t, val, err, TRUE)

	val, err = stc.Has(MakeStr("abc"))
	ok(t, val, err, FALSE)

	val, err = stc.Has(ZERO)
	fail(t, val, err, "TypeMismatch: Expected 'Str'")

	stc, err = NewStruct([]Field{NewField("a", true, NULL)}, false)
	if err != nil {
		panic("oops")
	}
	val, err = stc.GetField(cx, MakeStr("a"))
	ok(t, val, err, NULL)

	tassert(t, reflect.DeepEqual(stc.FieldNames(), []string{"a"}))

	stc, err = NewStruct([]Field{
		NewField("a", true, ONE),
		NewField("a", true, ZERO)}, false)
	fail(t, stc, err, "DuplicateField: Field 'a' is a duplicate")
}

func TestNativeProp(t *testing.T) {

	var propValue Value = ZERO

	getter := NewNativeFunc(0, 0,
		func(cx Context, values []Value) (Value, Error) {
			return propValue, nil
		})

	setter := NewNativeFunc(1, 1,
		func(cx Context, values []Value) (Value, Error) {
			propValue = values[0]
			return nil, nil
		})

	stc, err := NewStruct([]Field{NewProperty("a", getter, setter)}, false)
	tassert(t, err == nil)

	val, err := stc.GetField(cx, MakeStr("a"))
	ok(t, val, err, ZERO)
	tassert(t, propValue == ZERO)

	err = stc.SetField(cx, MakeStr("a"), ONE)
	tassert(t, err == nil)

	val, err = stc.GetField(cx, MakeStr("a"))
	ok(t, val, err, ONE)
	tassert(t, propValue == ONE)
}

func TestList(t *testing.T) {
	ls := NewList([]Value{})
	okType(t, ls, TLIST)

	var v Value
	var err Error

	v = ls.ToStr(cx)
	ok(t, v, nil, MakeStr("[ ]"))

	v, err = ls.Eq(cx, NewList([]Value{}))
	ok(t, v, err, TRUE)

	v, err = ls.Eq(cx, NewList([]Value{MakeStr("a")}))
	ok(t, v, err, FALSE)

	v, err = ls.Eq(cx, NULL)
	ok(t, v, err, FALSE)

	v = ls.Len()
	ok(t, v, nil, ZERO)

	ls.Add(cx, MakeStr("a"))

	v, err = ls.Eq(cx, NewList([]Value{}))
	ok(t, v, err, FALSE)

	v, err = ls.Eq(cx, NewList([]Value{MakeStr("a")}))
	ok(t, v, err, TRUE)

	v = ls.Len()
	ok(t, v, nil, ONE)

	v, err = ls.Get(cx, ZERO)
	ok(t, v, err, MakeStr("a"))

	err = ls.Set(cx, ZERO, MakeStr("b"))
	tassert(t, err == nil)

	v, err = ls.Get(cx, ZERO)
	ok(t, v, err, MakeStr("b"))

	v, err = ls.Get(cx, NEG_ONE)
	ok(t, v, err, MakeStr("b"))

	v, err = ls.Get(cx, ONE)
	fail(t, v, err, "IndexOutOfBounds: 1")

	err = ls.Set(cx, NEG_ONE, TRUE)
	tassert(t, err == nil)

	err = ls.Set(cx, ONE, TRUE)
	fail(t, nil, err, "IndexOutOfBounds: 1")

	v = ls.ToStr(cx)
	ok(t, v, nil, MakeStr("[ true ]"))

	ls.Add(cx, MakeStr("z"))

	v = ls.ToStr(cx)
	ok(t, v, nil, MakeStr("[ true, z ]"))
}

func TestCompositeHashCode(t *testing.T) {
	h, err := NewDict(cx, []*HEntry{}).HashCode(cx)
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")

	h, err = NewList([]Value{}).HashCode(cx)
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")

	h, err = newStruct([]Field{}).HashCode(cx)
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")
}

func TestDict(t *testing.T) {
	d := NewDict(cx, []*HEntry{})
	okType(t, d, TDICT)

	var v Value
	var err Error

	v = d.ToStr(cx)
	ok(t, v, err, MakeStr("dict { }"))

	v, err = d.Eq(cx, NewDict(cx, []*HEntry{}))
	ok(t, v, err, TRUE)

	v, err = d.Eq(cx, NULL)
	ok(t, v, err, FALSE)

	v = d.Len()
	ok(t, v, nil, ZERO)

	v, err = d.Get(cx, MakeStr("a"))
	ok(t, v, err, NULL)

	err = d.Set(cx, MakeStr("a"), ONE)
	tassert(t, err == nil)

	v, err = d.Get(cx, MakeStr("a"))
	ok(t, v, err, ONE)

	v, err = d.Eq(cx, NewDict(cx, []*HEntry{}))
	ok(t, v, err, FALSE)

	v, err = d.Eq(cx, NewDict(cx, []*HEntry{{MakeStr("a"), ONE}}))
	ok(t, v, err, TRUE)

	v = d.Len()
	ok(t, v, nil, ONE)

	v = d.ToStr(cx)
	ok(t, v, nil, MakeStr("dict { a: 1 }"))

	err = d.Set(cx, MakeStr("b"), MakeInt(2))
	tassert(t, err == nil)

	v, err = d.Get(cx, MakeStr("b"))
	ok(t, v, err, MakeInt(2))

	v = d.ToStr(cx)
	ok(t, v, nil, MakeStr("dict { b: 2, a: 1 }"))

	tp := NewTuple([]Value{ONE, ZERO})
	d = NewDict(cx, []*HEntry{{tp, TRUE}})

	v = d.ToStr(cx)
	ok(t, v, nil, MakeStr("dict { (1, 0): true }"))

	v, err = d.Get(cx, tp)
	ok(t, v, err, TRUE)
}

func TestSet(t *testing.T) {
	s := NewSet(cx, []Value{})
	okType(t, s, TSET)

	var v Value
	var err Error

	v = s.ToStr(cx)
	ok(t, v, err, MakeStr("set { }"))

	v, err = s.Eq(cx, NewSet(cx, []Value{}))
	ok(t, v, err, TRUE)

	v, err = s.Eq(cx, NewSet(cx, []Value{ONE}))
	ok(t, v, err, FALSE)

	v, err = s.Eq(cx, NULL)
	ok(t, v, err, FALSE)

	v = s.Len()
	ok(t, v, nil, ZERO)

	s = NewSet(cx, []Value{ONE})

	v = s.ToStr(cx)
	ok(t, v, err, MakeStr("set { 1 }"))

	v, err = s.Eq(cx, NewSet(cx, []Value{}))
	ok(t, v, err, FALSE)

	v, err = s.Eq(cx, NewSet(cx, []Value{ONE, ONE, ONE}))
	ok(t, v, err, TRUE)

	v, err = s.Eq(cx, NULL)
	ok(t, v, err, FALSE)

	v = s.Len()
	ok(t, v, nil, ONE)

	s = NewSet(cx, []Value{ONE, ZERO, ZERO, ONE})

	v = s.ToStr(cx)
	ok(t, v, err, MakeStr("set { 0, 1 }"))

	v = s.Len()
	ok(t, v, nil, MakeInt(2))
}

func TestTuple(t *testing.T) {
	var v Value
	var err Error

	tp := NewTuple([]Value{ONE, ZERO})
	okType(t, tp, TTUPLE)

	v, err = tp.Eq(cx, NewTuple([]Value{ZERO, ZERO}))
	ok(t, v, err, FALSE)

	v, err = tp.Eq(cx, NewTuple([]Value{ONE, ZERO}))
	ok(t, v, err, TRUE)

	v, err = tp.Eq(cx, NULL)
	ok(t, v, err, FALSE)

	v, err = tp.Get(cx, ZERO)
	ok(t, v, err, ONE)

	v, err = tp.Get(cx, ONE)
	ok(t, v, err, ZERO)

	v, err = tp.Get(cx, NEG_ONE)
	println("adfasfda", v.ToStr(cx))
	ok(t, v, err, ZERO)

	v, err = tp.Get(cx, MakeInt(2))
	fail(t, v, err, "IndexOutOfBounds: 2")

	v = tp.ToStr(cx)
	ok(t, v, nil, MakeStr("(1, 0)"))

	v = tp.Len()
	ok(t, v, nil, MakeInt(2))
}

func newRange(from int64, to int64, step int64) Range {
	r, err := NewRange(from, to, step)
	if err != nil {
		panic("invalid range")
	}
	return r
}

func TestRange(t *testing.T) {
	var v Value
	var err Error

	r := newRange(0, 5, 1)
	okType(t, r, TRANGE)

	v, err = r.Eq(cx, newRange(0, 5, 2))
	ok(t, v, err, FALSE)

	v, err = r.Eq(cx, newRange(0, 5, 1))
	ok(t, v, err, TRUE)

	v, err = r.Eq(cx, NULL)
	ok(t, v, err, FALSE)

	v = r.Len()
	ok(t, v, nil, MakeInt(5))

	v = newRange(0, 6, 3).Len()
	ok(t, v, nil, MakeInt(2))
	v = newRange(0, 7, 3).Len()
	ok(t, v, nil, MakeInt(3))
	v = newRange(0, 8, 3).Len()
	ok(t, v, nil, MakeInt(3))
	v = newRange(0, 9, 3).Len()
	ok(t, v, nil, MakeInt(3))

	v = newRange(0, 0, 3).Len()
	ok(t, v, nil, MakeInt(0))
	v = newRange(1, 0, 1).Len()
	ok(t, v, nil, MakeInt(0))

	v, err = NewRange(1, 0, 0)
	fail(t, v, err, "InvalidArgument: step cannot be 0")

	v = newRange(0, -5, -1).Len()
	ok(t, v, nil, MakeInt(5))
	v = newRange(-1, -8, -3).Len()
	ok(t, v, nil, MakeInt(3))

	r = newRange(0, 5, 1)
	v, err = r.Get(cx, ONE)
	ok(t, v, err, MakeInt(1))

	r = newRange(3, 9, 2)
	v, err = r.Get(cx, MakeInt(2))
	ok(t, v, err, MakeInt(7))

	r = newRange(-9, -13, -1)
	v, err = r.Get(cx, ONE)
	ok(t, v, err, MakeInt(-10))
}

func TestRangeIterator(t *testing.T) {

	var ibl Iterable = newRange(1, 5, 1)

	var itr Iterator = ibl.NewIterator(cx)
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	var n int64 = 1
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		tassert(t, err == nil)

		i, ok := v.(Int)
		tassert(t, ok)
		n *= i.IntVal()
	}
	tassert(t, n == 24)
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator(cx)
	n = 1
	for structInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, MakeStr("getValue"))

		i, ok := v.(Int)
		tassert(t, ok)
		n *= i.IntVal()
	}
	tassert(t, n == 24)
}

func TestListIterator(t *testing.T) {

	var ibl Iterable = NewList(
		[]Value{MakeInt(1), MakeInt(2), MakeInt(3), MakeInt(4)})

	var itr Iterator = ibl.NewIterator(cx)
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	var n int64 = 1
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		tassert(t, err == nil)

		i, ok := v.(Int)
		tassert(t, ok)
		n *= i.IntVal()
	}
	tassert(t, n == 24)
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator(cx)
	err = itr.SetField(cx, MakeStr("nextValue"), NULL)
	fail(t, nil, err, "ImmutableValue")
	err = itr.SetField(cx, MakeStr("getValue"), NULL)
	fail(t, nil, err, "ImmutableValue")

	n = 1
	for structInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, MakeStr("getValue"))

		i, ok := v.(Int)
		tassert(t, ok)
		n *= i.IntVal()
	}
	tassert(t, n == 24)
}

func TestDictIterator(t *testing.T) {

	var ibl Iterable = NewDict(cx,
		[]*HEntry{
			{MakeStr("a"), ONE},
			{MakeStr("b"), MakeInt(2)},
			{MakeStr("c"), MakeInt(3)}})

	var itr Iterator = ibl.NewIterator(cx)
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	s := MakeStr("")
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		tassert(t, err == nil)

		tp, ok := v.(Tuple)
		tassert(t, ok)
		s = s.Concat(tp.ToStr(cx))
	}
	ok(t, s, nil, MakeStr("(b, 2)(a, 1)(c, 3)"))
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator(cx)
	s = MakeStr("")
	for structInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, MakeStr("getValue"))

		tp, ok := v.(Tuple)
		tassert(t, ok)
		s = s.Concat(tp.ToStr(cx))
	}
	ok(t, s, nil, MakeStr("(b, 2)(a, 1)(c, 3)"))
}

func TestSetIterator(t *testing.T) {

	var ibl Iterable = NewSet(cx,
		[]Value{MakeStr("a"), MakeStr("b"), MakeStr("c")})

	var itr Iterator = ibl.NewIterator(cx)
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	s := MakeStr("")
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		tassert(t, err == nil)

		s = s.Concat(v.ToStr(cx))
	}
	ok(t, s, nil, MakeStr("bac"))
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator(cx)
	s = MakeStr("")
	for structInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, MakeStr("getValue"))

		s = s.Concat(v.ToStr(cx))
	}
	ok(t, s, nil, MakeStr("bac"))
}
