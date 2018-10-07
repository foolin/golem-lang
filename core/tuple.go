// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"fmt"
)

/*doc
## Tuple

A Tuple is an immutable sequence of two or more values.  Tuples are similar to Lists,
but they have special semantics in certain Golem expressions and statements.

Valid operators for Tuple are:

* The equality operators `==`, `!=`
* The [`index`](interfaces.html#indexable) operator `a[x]`

The index operator can return a value of any type.

Tuples are
[`lenable`](interfaces.html#lenable) and
[`iterable`](interfaces.html#iterable), and
[`hashable`](interfaces.html#hashable).

*/

type tuple []Value

// NewTuple creates a new Tuple
func NewTuple(values []Value) Tuple {
	if len(values) < 2 {
		panic(fmt.Errorf("invalid tuple size: %d", len(values)))
	}
	return tuple(values)
}

func (tp tuple) compositeMarker() {}

func (tp tuple) Type() Type { return TupleType }

func (tp tuple) Freeze(ev Eval) (Value, Error) {
	return tp, nil
}

func (tp tuple) Frozen(ev Eval) (Bool, Error) {
	return True, nil
}

func (tp tuple) ToStr(ev Eval) (Str, Error) {

	var buf bytes.Buffer
	buf.WriteString("(")
	for idx, v := range tp {
		if idx > 0 {
			buf.WriteString(", ")
		}
		s, err := v.ToStr(ev)
		if err != nil {
			return nil, err
		}
		buf.WriteString(s.String())
	}
	buf.WriteString(")")
	return NewStr(buf.String())
}

func (tp tuple) HashCode(ev Eval) (Int, Error) {

	// https://en.wikipedia.org/wiki/Jenkins_hash_function
	var hash int64
	for _, v := range tp {
		h, err := v.HashCode(ev)
		if err != nil {
			return nil, err
		}
		hash += h.ToInt()
		hash += hash << 10
		hash ^= hash >> 6
	}
	hash += hash << 3
	hash ^= hash >> 11
	hash += hash << 15
	return NewInt(hash), nil
}

func (tp tuple) Eq(ev Eval, v Value) (Bool, Error) {
	switch t := v.(type) {
	case tuple:
		return valuesEq(ev, tp, t)
	default:
		return False, nil
	}
}

func (tp tuple) Get(ev Eval, index Value) (Value, Error) {
	idx, err := boundedIndex(index, len(tp))
	if err != nil {
		return nil, err
	}
	return tp[idx], nil
}

func (tp tuple) Set(ev Eval, index Value, val Value) Error {
	return ImmutableValue()
}

func (tp tuple) Len(ev Eval) (Int, Error) {
	return NewInt(int64(len(tp))), nil
}

func (tp tuple) ToList() List {
	return NewList(CopyValues(tp))
}

func (tp tuple) Values() []Value {
	return []Value(tp)
}

//---------------------------------------------------------------
// Iterator

type tupleIterator struct {
	Struct
	tp tuple
	n  int
}

func (tp tuple) NewIterator(ev Eval) (Iterator, Error) {

	itr := &tupleIterator{iteratorStruct(), tp, -1}

	next, get := iteratorFields(ev, itr)
	itr.Internal("next", next)
	itr.Internal("get", get)

	return itr, nil
}

func (i *tupleIterator) IterNext(ev Eval) (Bool, Error) {
	i.n++
	return NewBool(i.n < len(i.tp)), nil
}

func (i *tupleIterator) IterGet(ev Eval) (Value, Error) {
	if (i.n >= 0) && (i.n < len(i.tp)) {
		return i.tp[i.n], nil
	}
	return nil, NoSuchElement()
}

//--------------------------------------------------------------
// fields

/*doc
A Tuple has the following fields:

* [toList](#tolist)

*/

var tupleMethods = map[string]Method{

	/*doc
	### `toList`

	`toList` creates a new List having the same elements as the tuple.

	* signature: `toList() <List>`
	* example: `(1,2,3).toList()`

	*/
	"toList": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			return self.(Tuple).ToList(), nil
		}),
}

func (tp tuple) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(tupleMethods))
	for name := range tupleMethods {
		names = append(names, name)
	}
	return names, nil
}

func (tp tuple) HasField(name string) (bool, Error) {
	_, ok := tupleMethods[name]
	return ok, nil
}

func (tp tuple) GetField(ev Eval, name string) (Value, Error) {
	if method, ok := tupleMethods[name]; ok {
		return method.ToFunc(tp, name), nil
	}
	return nil, NoSuchField(name)
}

func (tp tuple) InvokeField(ev Eval, name string, params []Value) (Value, Error) {
	if method, ok := tupleMethods[name]; ok {
		return method.Invoke(tp, ev, params)
	}
	return nil, NoSuchField(name)
}
