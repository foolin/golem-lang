// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	//"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestHashFieldMap(t *testing.T) {

	var fm fieldMap = &hashFieldMap{
		map[string]Field{
			"a": NewField(Zero),
			"b": NewReadonlyField(One),
		},
		true}

	//ab := reflect.DeepEqual([]string{"a", "b"}, fm.names())
	//ba := reflect.DeepEqual([]string{"b", "a"}, fm.names())
	//tassert(t, ab || ba)

	tassert(t, fm.has("a"))
	tassert(t, fm.has("b"))
	tassert(t, !fm.has("c"))

	val, err := fm.get(nil, "a")
	ok(t, val, err, Zero)
	val, err = fm.get(nil, "b")
	ok(t, val, err, One)
	val, err = fm.get(nil, "c")
	fail(t, val, err, "NoSuchField: Field 'c' not found")

	val, err = fm.invoke(nil, "a", []Value{})
	fail(t, val, err, "TypeMismatch: Expected Func, not Int")
	val, err = fm.invoke(nil, "b", []Value{})
	fail(t, val, err, "TypeMismatch: Expected Func, not Int")
	val, err = fm.invoke(nil, "c", []Value{})
	fail(t, val, err, "NoSuchField: Field 'c' not found")

	err = fm.set(nil, "a", One)
	tassert(t, err == nil)
	val, err = fm.get(nil, "a")
	ok(t, val, err, One)
	err = fm.set(nil, "b", One)
	fail(t, nil, err, "ReadonlyField: Field 'b' is readonly")
	err = fm.set(nil, "c", One)
	fail(t, nil, err, "NoSuchField: Field 'c' not found")

	fm.replace("a", NewField(NewStr("abc")))
	val, err = fm.get(nil, "a")
	ok(t, val, err, NewStr("abc"))
}

func TestMethodFieldMap(t *testing.T) {

	method := NewFixedMethod(
		[]Type{},
		false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			n := self.(Int).IntVal()
			return NewInt(int64(n * n)), nil
		})

	var fm fieldMap = &methodFieldMap{
		NewInt(7),
		map[string]Method{
			"a": method,
		}}

	tassert(t, reflect.DeepEqual([]string{"a"}, fm.names()))
	tassert(t, fm.has("a"))
	tassert(t, !fm.has("b"))

	val, err := fm.get(nil, "a")
	tassert(t, err == nil)
	fn := val.(Func)
	val, err = fn.Invoke(nil, []Value{})
	ok(t, val, err, NewInt(49))
	val, err = fm.get(nil, "b")
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	val, err = fm.invoke(nil, "a", []Value{})
	ok(t, val, err, NewInt(49))
	val, err = fm.invoke(nil, "b", []Value{})
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	err = fm.set(nil, "a", One)
	fail(t, val, err, "ReadonlyField: Field 'a' is readonly")
	err = fm.set(nil, "b", One)
	fail(t, val, err, "NoSuchField: Field 'b' not found")
}

var counter int64

func next() Int {
	n := NewInt(counter)
	counter++
	return n
}

func getMap(t *testing.T, fm fieldMap) map[string]Value {

	result := make(map[string]Value)
	names := fm.names()
	for _, n := range names {
		val, err := fm.get(nil, n)
		tassert(t, err == nil)

		if fn, ok := val.(Func); ok {
			val, err = fn.Invoke(nil, nil)
			tassert(t, err == nil)
		}

		result[n] = val
	}

	return result
}

func TestMergeFieldMaps(t *testing.T) {

	var x fieldMap = &hashFieldMap{
		map[string]Field{
			"a": NewField(next()),
			"b": NewField(next()),
			"e": NewField(next()),
			"f": NewField(next()),
		},
		true}
	tassert(t, x.(*hashFieldMap).replacable)

	var y fieldMap = &methodFieldMap{
		next(),
		map[string]Method{
			"b": NewFixedMethod(
				[]Type{}, false,
				func(self interface{}, ev Eval, params []Value) (Value, Error) {
					n := self.(Int).IntVal()
					return NewInt(int64(n * 100)), nil
				}),
			"c": NewFixedMethod(
				[]Type{}, false,
				func(self interface{}, ev Eval, params []Value) (Value, Error) {
					n := self.(Int).IntVal()
					return NewInt(int64(n*100 + 1)), nil
				}),
			"f": NewFixedMethod(
				[]Type{}, false,
				func(self interface{}, ev Eval, params []Value) (Value, Error) {
					n := self.(Int).IntVal()
					return NewInt(int64(n*100 + 2)), nil
				}),
		}}

	var z fieldMap = &hashFieldMap{
		map[string]Field{
			"c": NewField(next()),
			"d": NewField(next()),
			"e": NewField(next()),
			"f": NewField(next()),
		},
		true}
	tassert(t, z.(*hashFieldMap).replacable)

	w := mergeFieldMaps([]fieldMap{x, y, z})

	names := w.names()
	sort.Slice(names, func(i, j int) bool {
		return strings.Compare(names[i], names[j]) < 0
	})
	tassert(t, reflect.DeepEqual([]string{"a", "b", "c", "d", "e", "f"}, names))
	tassert(t, !w.(*hashFieldMap).replacable)

	//------------------------------------------------

	// ---------------
	// x:  a b . . e f
	// y:  . b c . . f
	// z:  . . c d e f
	// ---------------
	// w:  a b c d e f

	//fmt.Printf("%v\n", getMap(t, x))
	tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(0),
			"b": NewInt(1),
			"e": NewInt(2),
			"f": NewInt(3),
		}))

	//fmt.Printf("%v\n", getMap(t, y))
	tassert(t, reflect.DeepEqual(getMap(t, y),
		map[string]Value{
			"b": NewInt(400),
			"c": NewInt(401),
			"f": NewInt(402),
		}))

	//fmt.Printf("%v\n", getMap(t, z))
	tassert(t, reflect.DeepEqual(getMap(t, z),
		map[string]Value{
			"c": NewInt(5),
			"d": NewInt(6),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(0),
			"b": NewInt(400),
			"c": NewInt(5),
			"d": NewInt(6),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	//------------------------------------------------
	// a

	err := x.set(nil, "a", next())
	tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(1),
			"e": NewInt(2),
			"f": NewInt(3),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(400),
			"c": NewInt(5),
			"d": NewInt(6),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	//------------------------------------------------
	// b

	err = x.set(nil, "b", next())
	tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(10),
			"e": NewInt(2),
			"f": NewInt(3),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(400),
			"c": NewInt(5),
			"d": NewInt(6),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	//------------------------------------------------
	// c

	err = z.set(nil, "c", next())
	tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, z))
	tassert(t, reflect.DeepEqual(getMap(t, z),
		map[string]Value{
			"c": NewInt(11),
			"d": NewInt(6),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(400),
			"c": NewInt(11),
			"d": NewInt(6),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	//------------------------------------------------
	// d

	err = z.set(nil, "d", next())
	tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, z))
	tassert(t, reflect.DeepEqual(getMap(t, z),
		map[string]Value{
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(400),
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	//------------------------------------------------
	// e

	err = x.set(nil, "e", next())
	tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(10),
			"e": NewInt(13),
			"f": NewInt(3),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(400),
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	err = z.set(nil, "e", next())
	tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(10),
			"e": NewInt(13),
			"f": NewInt(3),
		}))

	//fmt.Printf("%v\n", getMap(t, z))
	tassert(t, reflect.DeepEqual(getMap(t, z),
		map[string]Value{
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(14),
			"f": NewInt(8),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(400),
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(14),
			"f": NewInt(8),
		}))

	//------------------------------------------------
	// f

	err = x.set(nil, "f", next())
	tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(10),
			"e": NewInt(13),
			"f": NewInt(15),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(400),
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(14),
			"f": NewInt(8),
		}))

	err = z.set(nil, "f", next())
	tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(10),
			"e": NewInt(13),
			"f": NewInt(15),
		}))

	//fmt.Printf("%v\n", getMap(t, z))
	tassert(t, reflect.DeepEqual(getMap(t, z),
		map[string]Value{
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(14),
			"f": NewInt(16),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(400),
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(14),
			"f": NewInt(16),
		}))

	//------------------------------------------------
	// bidrectional

	err = x.set(nil, "a", next())
	tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(17),
			"b": NewInt(10),
			"e": NewInt(13),
			"f": NewInt(15),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(17),
			"b": NewInt(400),
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(14),
			"f": NewInt(16),
		}))

	err = w.set(nil, "a", next())
	tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(18),
			"b": NewInt(10),
			"e": NewInt(13),
			"f": NewInt(15),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(18),
			"b": NewInt(400),
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(14),
			"f": NewInt(16),
		}))

	//------------------------------------------------
	// 'y' is immutable

	tassert(t, reflect.DeepEqual(getMap(t, y),
		map[string]Value{
			"b": NewInt(400),
			"c": NewInt(401),
			"f": NewInt(402),
		}))
}
