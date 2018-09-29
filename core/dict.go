// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"fmt"
)

/*doc
## Dict

A Dict is an [associative array](https://en.wikipedia.org/wiki/Associative_array),
in which the keys are all [`hashable`](interfaces.html#hashable).

Valid operators for Dict are:

* The equality operators `==`, `!=`
* The [`index`](interfaces.html#indexable) operator `a[x]`

Dicts are
[`lenable`](interfaces.html#lenable) and
[`iterable`](interfaces.html#iterable).

Each iterated element in a Dict is a 2-Tuple containing a key-value pair.

*/

type dict struct {
	hashMap *HashMap
	frozen  bool
}

// NewDict creates a new Dict
func NewDict(h *HashMap) Dict {
	return &dict{h, false}
}

func (d *dict) compositeMarker() {}

func (d *dict) Type() Type { return DictType }

func (d *dict) Freeze(ev Eval) (Value, Error) {
	d.frozen = true
	return d, nil
}

func (d *dict) Frozen(ev Eval) (Bool, Error) {
	return NewBool(d.frozen), nil
}

func (d *dict) ToStr(ev Eval) (Str, Error) {

	var buf bytes.Buffer
	buf.WriteString("dict {")
	idx := 0
	itr := d.hashMap.Iterator()

	for itr.Next() {
		entry := itr.Get()
		if idx > 0 {
			buf.WriteString(",")
		}
		idx++

		buf.WriteString(" ")
		s, err := entry.Key.ToStr(ev)
		if err != nil {
			return nil, err
		}
		buf.WriteString(s.String())

		buf.WriteString(": ")
		s, err = entry.Value.ToStr(ev)
		if err != nil {
			return nil, err
		}
		buf.WriteString(s.String())
	}

	buf.WriteString(" }")
	return NewStr(buf.String())
}

func (d *dict) HashCode(ev Eval) (Int, Error) {
	return nil, HashCodeMismatch(DictType)
}

func (d *dict) Eq(ev Eval, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *dict:
		return d.hashMap.Eq(ev, t.hashMap)
	default:
		return False, nil
	}
}

func (d *dict) Get(ev Eval, key Value) (Value, Error) {
	return d.hashMap.Get(ev, key)
}

func (d *dict) Set(ev Eval, key Value, val Value) Error {
	if d.frozen {
		return ImmutableValue()
	}

	return d.hashMap.Put(ev, key, val)
}

func (d *dict) Len(ev Eval) (Int, Error) {
	return d.hashMap.Len(), nil
}

func (d *dict) HashMap() *HashMap {
	return d.hashMap
}

//---------------------------------------------------------------

func (d *dict) IsEmpty() Bool {
	return NewBool(d.hashMap.Len().IntVal() == 0)
}

func (d *dict) Contains(ev Eval, key Value) (Bool, Error) {
	return d.hashMap.Contains(ev, key)
}

func (d *dict) Clear() (Dict, Error) {
	if d.frozen {
		return nil, ImmutableValue()
	}

	d.hashMap = EmptyHashMap()
	return d, nil
}

func (d *dict) Remove(ev Eval, key Value) (Dict, Error) {
	if d.frozen {
		return nil, ImmutableValue()
	}

	_, err := d.hashMap.Remove(ev, key)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *dict) AddAll(ev Eval, val Value) (Dict, Error) {

	if d.frozen {
		return nil, ImmutableValue()
	}

	ibl, ok := val.(Iterable)
	if !ok {
		return nil, IterableMismatch(val.Type())
	}

	itr, err := ibl.NewIterator(ev)
	if err != nil {
		return nil, err
	}

	b, err := itr.IterNext(ev)
	if err != nil {
		return nil, err
	}
	for b.BoolVal() {
		v, err := itr.IterGet(ev)
		if err != nil {
			return nil, err
		}

		tp, ok := v.(tuple)
		if !ok {
			return nil, TypeMismatch(TupleType, v.Type())
		}

		if len(tp) != 2 {
			return nil, InvalidArgument(
				fmt.Sprintf("Expected Tuple of length %d, not length %d",
					2, len(tp)))
		}

		err = d.hashMap.Put(ev, tp[0], tp[1])
		if err != nil {
			return nil, err
		}

		b, err = itr.IterNext(ev)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

//---------------------------------------------------------------
// Iterator

type dictIterator struct {
	Struct
	d       *dict
	itr     *HIterator
	hasNext bool
}

func (d *dict) NewIterator(ev Eval) (Iterator, Error) {

	itr := &dictIterator{iteratorStruct(), d, d.hashMap.Iterator(), false}

	next, get := iteratorFields(ev, itr)
	itr.Internal("next", next)
	itr.Internal("get", get)

	return itr, nil
}

func (i *dictIterator) IterNext(ev Eval) (Bool, Error) {
	i.hasNext = i.itr.Next()
	return NewBool(i.hasNext), nil
}

func (i *dictIterator) IterGet(ev Eval) (Value, Error) {

	if i.hasNext {
		entry := i.itr.Get()
		return NewTuple([]Value{entry.Key, entry.Value}), nil
	}
	return nil, NoSuchElement()
}

//--------------------------------------------------------------
// fields

/*doc
A Dict has the following fields:

* [addAll](#addall)
* [clear](#clear)
* [contains](#contains)
* [isEmpty](#isempty)
* [remove](#remove)

*/

var dictMethods = map[string]Method{

	/*doc
	### `addAll`

	`addAll` adds all of the values in the given [Iterable](interfaces.html#iterable)
	to the dict, and returns the modified dict.
	Each iterated element must be a 2-Tuple containing a key-value pair.

	* signature: `addAll(itr <Iterable>) <Dict>`
	* example:

	```
	let d = dict {'a': 1, 'b': 2}
	println(d.addAll([('b', 2), ('c', 3)]))
	```

	*/
	"addAll": NewFixedMethod(
		[]Type{AnyType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			d := self.(Dict)
			return d.AddAll(ev, params[0])
		}),

	/*doc
	### `clear`

	`clear` removes all of the entries from the dict, and returns the empty dict.

	* signature: `clear() <Dict>`
	* example:

	```
	let d = dict {'a': 1, 'b': 2}
	println(d.clear())
	```

	*/
	"clear": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			d := self.(Dict)
			return d.Clear()
		}),

	/*doc
	### `contains`

	`contains` returns whether the given key is present in the dict.

	* signature: `contains(key <Value>) <Bool>`
	* example:

	```
	let d = dict {'a': 1, 'b': 2}
	println(d.contains('b'))
	```

	*/
	"contains": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			d := self.(Dict)
			return d.Contains(ev, params[0])
		}),

	/*doc
	### `isEmpty`

	`isEmpty` returns whether the dict contains any values.

	* signature: `isEmpty() <Bool>`
	* example: `println(dict {}.isEmpty())`

	*/
	"isEmpty": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			d := self.(Dict)
			return d.IsEmpty(), nil
		}),

	/*doc
	### `remove`

	`remove` remove the entry associated with the given key from the dict,
	and returns modified dict.  If the key is not present in the dict, then
	the dict is unmodified.

	* signature: `remove(key <Value>) <Dict>`
	* example:

	```
	let d = dict {'a': 1, 'b': 2}
	println(d.remove('a'))
	```

	*/
	"remove": NewFixedMethod(
		[]Type{AnyType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			d := self.(Dict)
			return d.Remove(ev, params[0])
		}),
}

func (d *dict) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(dictMethods))
	for name := range dictMethods {
		names = append(names, name)
	}
	return names, nil
}

func (d *dict) HasField(name string) (bool, Error) {
	_, ok := dictMethods[name]
	return ok, nil
}

func (d *dict) GetField(ev Eval, name string) (Value, Error) {
	if method, ok := dictMethods[name]; ok {
		return method.ToFunc(d, name), nil
	}
	return nil, NoSuchField(name)
}

func (d *dict) InvokeField(ev Eval, name string, params []Value) (Value, Error) {

	if method, ok := dictMethods[name]; ok {
		return method.Invoke(d, ev, params)
	}
	return nil, NoSuchField(name)
}
