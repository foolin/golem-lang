// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
)

/*doc
## Set

A Set is an un-ordered collection of unique, [`hashable`](interfaces.html#hashable) values.

Valid operators for Set are:

* The equality operators `==`, `!=`

Sets are
[`lenable`](interfaces.html#lenable) and
[`iterable`](interfaces.html#iterable).

*/

type set struct {
	hashMap *HashMap
	frozen  bool
}

// NewSet creates a new Set
func NewSet(ev Eval, values []Value) (Set, Error) {

	hashMap := EmptyHashMap()
	for _, v := range values {
		err := hashMap.Put(ev, v, True)
		if err != nil {
			return nil, err
		}
	}

	return &set{hashMap, false}, nil
}

func (s *set) compositeMarker() {}

func (s *set) Type() Type { return SetType }

func (s *set) Freeze(ev Eval) (Value, Error) {
	s.frozen = true
	return s, nil
}

func (s *set) Frozen(ev Eval) (Bool, Error) {
	return NewBool(s.frozen), nil
}

func (s *set) ToStr(ev Eval) (Str, Error) {

	var buf bytes.Buffer
	buf.WriteString("set {")
	idx := 0
	itr := s.hashMap.Iterator()

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
	}

	buf.WriteString(" }")
	return NewStr(buf.String())
}

func (s *set) HashCode(ev Eval) (Int, Error) {
	return nil, HashCodeMismatch(SetType)
}

func (s *set) Eq(ev Eval, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *set:
		return s.hashMap.Eq(ev, t.hashMap)
	default:
		return False, nil
	}
}

func (s *set) Len(ev Eval) (Int, Error) {
	return s.hashMap.Len(), nil
}

//---------------------------------------------------------------

func (s *set) IsEmpty() Bool {
	return NewBool(s.hashMap.Len().ToInt() == 0)
}

func (s *set) Contains(ev Eval, key Value) (Bool, Error) {
	return s.hashMap.Contains(ev, key)
}

func (s *set) Add(ev Eval, val Value) (Set, Error) {
	if s.frozen {
		return nil, ImmutableValue()
	}

	err := s.hashMap.Put(ev, val, True)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *set) AddAll(ev Eval, ibl Iterable) (Set, Error) {
	if s.frozen {
		return nil, ImmutableValue()
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

		err = s.hashMap.Put(ev, v, True)
		if err != nil {
			return nil, err
		}

		b, err = itr.IterNext(ev)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *set) Clear() (Set, Error) {
	if s.frozen {
		return nil, ImmutableValue()
	}

	s.hashMap = EmptyHashMap()
	return s, nil
}

func (s *set) Remove(ev Eval, key Value) (Set, Error) {
	if s.frozen {
		return nil, ImmutableValue()
	}

	_, err := s.hashMap.Remove(ev, key)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *set) Copy(ev Eval) (Set, Error) {

	values := []Value{}

	itr := s.hashMap.Iterator()
	for itr.Next() {
		entry := itr.Get()
		values = append(values, entry.Key)
	}

	return NewSet(ev, values)
}

func (s *set) ContainsAll(ev Eval, ibl Iterable) (Bool, Error) {

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

		//--------------------------
		c, err := s.hashMap.Contains(ev, v)
		if err != nil {
			return nil, err
		}
		if !c.BoolVal() {
			return False, nil
		}
		//--------------------------

		b, err = itr.IterNext(ev)
		if err != nil {
			return nil, err
		}
	}

	//--------------------------
	return True, nil
	//--------------------------
}

func (s *set) ContainsAny(ev Eval, ibl Iterable) (Bool, Error) {

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

		//--------------------------
		c, err := s.hashMap.Contains(ev, v)
		if err != nil {
			return nil, err
		}
		if c.BoolVal() {
			return True, nil
		}
		//--------------------------

		b, err = itr.IterNext(ev)
		if err != nil {
			return nil, err
		}
	}

	//--------------------------
	return False, nil
	//--------------------------
}

//---------------------------------------------------------------
// Iterator

type setIterator struct {
	Struct
	s       *set
	itr     *HIterator
	hasNext bool
}

func (s *set) NewIterator(ev Eval) (Iterator, Error) {

	itr := &setIterator{iteratorStruct(), s, s.hashMap.Iterator(), false}

	next, get := iteratorFields(ev, itr)
	itr.Internal("next", next)
	itr.Internal("get", get)

	return itr, nil
}

func (i *setIterator) IterNext(ev Eval) (Bool, Error) {
	i.hasNext = i.itr.Next()
	return NewBool(i.hasNext), nil
}

func (i *setIterator) IterGet(ev Eval) (Value, Error) {

	if i.hasNext {
		entry := i.itr.Get()
		return entry.Key, nil
	}
	return nil, NoSuchElement()
}

//--------------------------------------------------------------
// fields

/*doc
A Set has the following fields:

* [add](#add)
* [addAll](#addall)
* [clear](#clear)
* [contains](#contains)
* [containsAll](#containsall)
* [containsAny](#containsany)
* [copy](#copy)
* [isEmpty](#isempty)
* [remove](#remove)

*/

var setMethods = map[string]Method{

	/*doc
	### `add`

	`add` adds a value to the set, and returns the modified set.

	* signature: `add(val <Value>) <Set>`
	* example:

	```
	let a = set {1, 2, 3}
	println(a.add(4))
	```

	*/
	"add": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(Set)
			return s.Add(ev, params[0])
		}),

	/*doc
	### `addAll`

	`addAll` adds all of the values in the given [Iterable](interfaces.html#iterable) to
	the set, and returns the modified set.

	* signature: `addAll(itr <Iterable>) <Set>`
	* example:

	```
	let a = set {1, 2}
	println(a.addAll([3, 4]))
	```

	*/
	"addAll": NewFixedMethod(
		[]Type{AnyType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(Set)

			ibl, ok := params[0].(Iterable)
			if !ok {
				return nil, IterableMismatch(params[0].Type())
			}

			return s.AddAll(ev, ibl)
		}),

	/*doc
	### `clear`

	`clear` removes all of the values from the set, and returns the empty set.

	* signature: `clear() <Set>`
	* example:

	```
	let a = set {1, 2}
	println(a.clear())
	```

	*/
	"clear": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			s := self.(Set)
			return s.Clear()
		}),

	/*doc
	### `contains`

	`contains` returns whether the given value is an element in the set.

	* signature: `contains(val <Value>) <Bool>`
	* example:

	```
	let a = set {1, 2}
	println(a.contains(2))
	```

	*/
	"contains": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(Set)
			return s.Contains(ev, params[0])
		}),

	/*doc
	### `containsAll`

	`containsAll` returns whether the set contains all of the values
	in the given [Iterable](interfaces.html#iterable).

	* signature: `containsAll(itr <Iterable>) <Bool>`
	* example:

	```
	let a = set {1, 2}
	println(a.containsAll([1, 2]))
	```

	*/
	"containsAll": NewFixedMethod(
		[]Type{AnyType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(Set)

			ibl, ok := params[0].(Iterable)
			if !ok {
				return nil, IterableMismatch(params[0].Type())
			}

			return s.ContainsAll(ev, ibl)
		}),

	/*doc
	### `containsAny`

	`containsAny` returns whether the set contains all of the values
	in the given [Iterable](interfaces.html#iterable).

	* signature: `containsAny(itr <Iterable>) <Bool>`
	* example:

	```
	let a = set {1, 2}
	println(a.containsAny([1, 2]))
	```

	*/
	"containsAny": NewFixedMethod(
		[]Type{AnyType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(Set)

			ibl, ok := params[0].(Iterable)
			if !ok {
				return nil, IterableMismatch(params[0].Type())
			}

			return s.ContainsAny(ev, ibl)
		}),

	/*doc
	### `copy`

	`copy` returns a shallow copy of the set

	* signature: `copy() <Set>`
	* example:

	```
	println(set {1,2}.copy())
	```

	*/
	"copy": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			s := self.(Set)
			return s.Copy(ev)
		}),

	/*doc
	### `isEmpty`

	`isEmpty` returns whether the set contains any values.

	* signature: `isEmpty() <Bool>`
	* example: `println(set{}.isEmpty())`

	*/
	"isEmpty": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			s := self.(Set)
			return s.IsEmpty(), nil
		}),

	/*doc
	### `remove`

	`remove` remove the value from the set, and returns the
	modified set.  If the value is not present in the set, then
	the set is unmodified.

	* signature: `remove(value <Value>) <Set>`
	* example:

	```
	let a = set {1, 2, 3}
	println(a.remove(2))
	```

	*/
	"remove": NewFixedMethod(
		[]Type{AnyType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(Set)
			return s.Remove(ev, params[0])
		}),
}

func (s *set) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(setMethods))
	for name := range setMethods {
		names = append(names, name)
	}
	return names, nil
}

func (s *set) HasField(name string) (bool, Error) {
	_, ok := setMethods[name]
	return ok, nil
}

func (s *set) GetField(ev Eval, name string) (Value, Error) {
	if method, ok := setMethods[name]; ok {
		return method.ToFunc(s, name), nil
	}
	return nil, NoSuchField(name)
}

func (s *set) InvokeField(ev Eval, name string, params []Value) (Value, Error) {

	if method, ok := setMethods[name]; ok {
		return method.Invoke(s, ev, params)
	}
	return nil, NoSuchField(name)
}
