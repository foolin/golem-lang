// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"math"
	"reflect"
)

/*doc
## Range

A Range is a representation of an immutable sequence of integers.
Note that a Range doesn't actually contain a list of all its
Ints -- it produces them one at a time on demand.
A new Range is created by the [`range()`](#TODO) builtin function.

Valid operators for Range are:
	* The equality operators `==`, `!=`
	* The index operator `a[x]`

The index operator always return an Int.

Ranges have a [`len()`](#TODO) and are [`iterable`](#TODO).

*/

type rng struct {
	from  int64
	to    int64
	step  int64
	count int64
}

// NewRange creates a new Range
func NewRange(from int64, to int64, step int64) (Range, Error) {

	switch {

	case step == 0:
		return nil, InvalidArgument("step cannot be 0")

	case ((step > 0) && (from > to)) || ((step < 0) && (from < to)):
		return &rng{from, to, step, 0}, nil

	default:
		count := int64(math.Ceil(float64(to-from) / float64(step)))
		return &rng{from, to, step, count}, nil
	}
}

func (r *rng) compositeMarker() {}

func (r *rng) Type() Type { return RangeType }

func (r *rng) Freeze(ev Eval) (Value, Error) {
	return r, nil
}

func (r *rng) Frozen(ev Eval) (Bool, Error) {
	return True, nil
}

func (r *rng) ToStr(ev Eval) (Str, Error) {
	return NewStr(fmt.Sprintf("range<%d, %d, %d>", r.from, r.to, r.step))
}

func (r *rng) HashCode(ev Eval) (Int, Error) {
	return nil, HashCodeMismatch(RangeType)
}

func (r *rng) Eq(ev Eval, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *rng:
		return NewBool(reflect.DeepEqual(r, t)), nil
	default:
		return False, nil
	}
}

func (r *rng) Get(ev Eval, index Value) (Value, Error) {
	idx, err := boundedIndex(index, int(r.count))
	if err != nil {
		return nil, err
	}
	return NewInt(r.from + int64(idx)*r.step), nil
}

func (r *rng) Set(ev Eval, index Value, val Value) Error {
	return ImmutableValue()
}

func (r *rng) Len(ev Eval) (Int, Error) {
	return NewInt(r.count), nil
}

func (r *rng) From() Int  { return NewInt(r.from) }
func (r *rng) To() Int    { return NewInt(r.to) }
func (r *rng) Step() Int  { return NewInt(r.step) }
func (r *rng) Count() Int { return NewInt(r.count) }

//---------------------------------------------------------------
// Iterator

type rangeIterator struct {
	Struct
	r *rng
	n int64
}

func (r *rng) NewIterator(ev Eval) (Iterator, Error) {

	itr := &rangeIterator{iteratorStruct(), r, -1}

	next, get := iteratorFields(ev, itr)
	itr.Internal("next", next)
	itr.Internal("get", get)

	return itr, nil
}

func (i *rangeIterator) IterNext(ev Eval) (Bool, Error) {
	i.n++
	return NewBool(i.n < i.r.count), nil
}

func (i *rangeIterator) IterGet(ev Eval) (Value, Error) {

	if (i.n >= 0) && (i.n < i.r.count) {
		return NewInt(i.r.from + i.n*i.r.step), nil
	}
	return nil, NoSuchElement()
}

//--------------------------------------------------------------
// fields

/*doc
Range has the following fields:

*/

var rangeMethods = map[string]Method{

	/*doc
	#### `count`

	`count` is the total number of Ints in the range.

		* signature: `count() <Int>`

	*/
	"count": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			r := self.(Range)
			return r.Count(), nil
		}),

	/*doc
	#### `from`

	`from` is the first Int in the range, inclusive

		* signature: `from() <Int>`

	*/
	"from": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			r := self.(Range)
			return r.From(), nil
		}),

	/*doc
	#### `step`

	`step` is the distance between succesive Ints in the range.

		* signature: `step() <Int>`

	*/
	"step": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			r := self.(Range)
			return r.Step(), nil
		}),

	/*doc
	#### `to`

	`to` is the last Int in the range, exclusive

		* signature: `to() <Int>`

	*/
	"to": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			r := self.(Range)
			return r.To(), nil
		}),
}

func (r *rng) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(rangeMethods))
	for name := range rangeMethods {
		names = append(names, name)
	}
	return names, nil
}

func (r *rng) HasField(name string) (bool, Error) {
	_, ok := rangeMethods[name]
	return ok, nil
}

func (r *rng) GetField(ev Eval, name string) (Value, Error) {
	if method, ok := rangeMethods[name]; ok {
		return method.ToFunc(r, name), nil
	}
	return nil, NoSuchField(name)
}

func (r *rng) InvokeField(ev Eval, name string, params []Value) (Value, Error) {

	if method, ok := rangeMethods[name]; ok {
		return method.Invoke(r, ev, params)
	}
	return nil, NoSuchField(name)
}
