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

	Tassert(t, reflect.DeepEqual([]string{"a", "b"}, fm.names()))
	Tassert(t, fm.has("a"))
	Tassert(t, fm.has("b"))
	Tassert(t, !fm.has("c"))

	val, err := fm.get("a", nil)
	ok(t, val, err, Zero)
	val, err = fm.get("b", nil)
	ok(t, val, err, One)
	val, err = fm.get("c", nil)
	fail(t, val, err, "NoSuchField: Field 'c' not found")

	val, err = fm.invoke("a", nil, []Value{})
	fail(t, val, err, "TypeMismatch: Expected Func")
	val, err = fm.invoke("b", nil, []Value{})
	fail(t, val, err, "TypeMismatch: Expected Func")
	val, err = fm.invoke("c", nil, []Value{})
	fail(t, val, err, "NoSuchField: Field 'c' not found")

	err = fm.set("a", nil, One)
	Tassert(t, err == nil)
	val, err = fm.get("a", nil)
	ok(t, val, err, One)
	err = fm.set("b", nil, One)
	fail(t, nil, err, "ReadonlyField: Field 'b' is readonly")
	err = fm.set("c", nil, One)
	fail(t, nil, err, "NoSuchField: Field 'c' not found")

	fm.replace("a", NewField(NewStr("abc")))
	val, err = fm.get("a", nil)
	ok(t, val, err, NewStr("abc"))
}

func TestVirtualFieldMap(t *testing.T) {

	method := NewFixedMethod(
		[]Type{},
		false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			n := self.(Int).IntVal()
			return NewInt(int64(n * n)), nil
		})

	var fm fieldMap = &virtualFieldMap{
		NewInt(7),
		map[string]Method{
			"a": method,
		}}

	Tassert(t, reflect.DeepEqual([]string{"a"}, fm.names()))
	Tassert(t, fm.has("a"))
	Tassert(t, !fm.has("b"))

	val, err := fm.get("a", nil)
	Tassert(t, err == nil)
	fn := val.(Func)
	val, err = fn.Invoke(nil, []Value{})
	ok(t, val, err, NewInt(49))
	val, err = fm.get("b", nil)
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	val, err = fm.invoke("a", nil, []Value{})
	ok(t, val, err, NewInt(49))
	val, err = fm.invoke("b", nil, []Value{})
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	err = fm.set("a", nil, One)
	fail(t, val, err, "ReadonlyField: Field 'a' is readonly")
	err = fm.set("b", nil, One)
	fail(t, val, err, "NoSuchField: Field 'b' not found")
}

var counter int64 = 0

func next() Int {
	n := NewInt(counter)
	counter++
	return n
}

func getMap(t *testing.T, fm fieldMap) map[string]Value {

	result := make(map[string]Value)
	names := fm.names()
	for _, n := range names {
		val, err := fm.get(n, nil)
		Tassert(t, err == nil)

		if fn, ok := val.(Func); ok {
			val, err = fn.Invoke(nil, nil)
			Tassert(t, err == nil)
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
	Tassert(t, x.(*hashFieldMap).replacable)

	var y fieldMap = &virtualFieldMap{
		next(),
		map[string]Method{
			"b": NewFixedMethod(
				[]Type{}, false,
				func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
					n := self.(Int).IntVal()
					return NewInt(int64(n * 100)), nil
				}),
			"c": NewFixedMethod(
				[]Type{}, false,
				func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
					n := self.(Int).IntVal()
					return NewInt(int64(n*100 + 1)), nil
				}),
			"f": NewFixedMethod(
				[]Type{}, false,
				func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
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
	Tassert(t, z.(*hashFieldMap).replacable)

	w := mergeFieldMaps([]fieldMap{x, y, z})

	names := w.names()
	sort.Slice(names, func(i, j int) bool {
		return strings.Compare(names[i], names[j]) < 0
	})
	Tassert(t, reflect.DeepEqual([]string{"a", "b", "c", "d", "e", "f"}, names))
	Tassert(t, !w.(*hashFieldMap).replacable)

	//------------------------------------------------

	// ---------------
	// x:  a b . . e f
	// y:  . b c . . f
	// z:  . . c d e f
	// ---------------
	// w:  a b c d e f

	//fmt.Printf("%v\n", getMap(t, x))
	Tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(0),
			"b": NewInt(1),
			"e": NewInt(2),
			"f": NewInt(3),
		}))

	//fmt.Printf("%v\n", getMap(t, y))
	Tassert(t, reflect.DeepEqual(getMap(t, y),
		map[string]Value{
			"b": NewInt(400),
			"c": NewInt(401),
			"f": NewInt(402),
		}))

	//fmt.Printf("%v\n", getMap(t, z))
	Tassert(t, reflect.DeepEqual(getMap(t, z),
		map[string]Value{
			"c": NewInt(5),
			"d": NewInt(6),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	Tassert(t, reflect.DeepEqual(getMap(t, w),
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

	err := x.set("a", nil, next())
	Tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	Tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(1),
			"e": NewInt(2),
			"f": NewInt(3),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	Tassert(t, reflect.DeepEqual(getMap(t, w),
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

	err = x.set("b", nil, next())
	Tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	Tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(10),
			"e": NewInt(2),
			"f": NewInt(3),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	Tassert(t, reflect.DeepEqual(getMap(t, w),
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

	err = z.set("c", nil, next())
	Tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, z))
	Tassert(t, reflect.DeepEqual(getMap(t, z),
		map[string]Value{
			"c": NewInt(11),
			"d": NewInt(6),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	Tassert(t, reflect.DeepEqual(getMap(t, w),
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

	err = z.set("d", nil, next())
	Tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, z))
	Tassert(t, reflect.DeepEqual(getMap(t, z),
		map[string]Value{
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	Tassert(t, reflect.DeepEqual(getMap(t, w),
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

	err = x.set("e", nil, next())
	Tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	Tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(10),
			"e": NewInt(13),
			"f": NewInt(3),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	Tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(400),
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(7),
			"f": NewInt(8),
		}))

	err = z.set("e", nil, next())
	Tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	Tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(10),
			"e": NewInt(13),
			"f": NewInt(3),
		}))

	//fmt.Printf("%v\n", getMap(t, z))
	Tassert(t, reflect.DeepEqual(getMap(t, z),
		map[string]Value{
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(14),
			"f": NewInt(8),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	Tassert(t, reflect.DeepEqual(getMap(t, w),
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

	err = x.set("f", nil, next())
	Tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	Tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(10),
			"e": NewInt(13),
			"f": NewInt(15),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	Tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(400),
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(14),
			"f": NewInt(8),
		}))

	err = z.set("f", nil, next())
	Tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	Tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(9),
			"b": NewInt(10),
			"e": NewInt(13),
			"f": NewInt(15),
		}))

	//fmt.Printf("%v\n", getMap(t, z))
	Tassert(t, reflect.DeepEqual(getMap(t, z),
		map[string]Value{
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(14),
			"f": NewInt(16),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	Tassert(t, reflect.DeepEqual(getMap(t, w),
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

	err = x.set("a", nil, next())
	Tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	Tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(17),
			"b": NewInt(10),
			"e": NewInt(13),
			"f": NewInt(15),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	Tassert(t, reflect.DeepEqual(getMap(t, w),
		map[string]Value{
			"a": NewInt(17),
			"b": NewInt(400),
			"c": NewInt(11),
			"d": NewInt(12),
			"e": NewInt(14),
			"f": NewInt(16),
		}))

	err = w.set("a", nil, next())
	Tassert(t, err == nil)

	//fmt.Printf("%v\n", getMap(t, x))
	Tassert(t, reflect.DeepEqual(getMap(t, x),
		map[string]Value{
			"a": NewInt(18),
			"b": NewInt(10),
			"e": NewInt(13),
			"f": NewInt(15),
		}))

	//fmt.Printf("%v\n", getMap(t, w))
	Tassert(t, reflect.DeepEqual(getMap(t, w),
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

	Tassert(t, reflect.DeepEqual(getMap(t, y),
		map[string]Value{
			"b": NewInt(400),
			"c": NewInt(401),
			"f": NewInt(402),
		}))
}
