// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

/*doc
## List

A List is an ordered array of values.

Valid operators for List are:

* The equality operators `==`, `!=`
* The [`index`](interfaces.html#indexable) operator `a[x]`
* The [`slice`](interfaces.html#sliceable) operators `a[x:y]`, `a[x:]`, `a[:y]`

The index operator can return a value of any type.

The slice operators always return a List.

Lists are
[`lenable`](interfaces.html#lenable) and
[`iterable`](interfaces.html#iterable).

*/

type list struct {
	values []Value
	frozen bool
}

// NewList creates a new List
func NewList(values []Value) List {
	return &list{values, false}
}

func (ls *list) compositeMarker() {}

func (ls *list) Type() Type { return ListType }

func (ls *list) Freeze(ev Eval) (Value, Error) {
	ls.frozen = true
	return ls, nil
}

func (ls *list) Frozen(ev Eval) (Bool, Error) {
	return NewBool(ls.frozen), nil
}

func (ls *list) ToStr(ev Eval) (Str, Error) {

	var buf bytes.Buffer
	buf.WriteString("[")
	for idx, v := range ls.values {
		if idx > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(" ")

		s, err := v.ToStr(ev)
		if err != nil {
			return nil, err
		}

		buf.WriteString(s.String())
	}
	buf.WriteString(" ]")

	return NewStr(buf.String())
}

func (ls *list) HashCode(ev Eval) (Int, Error) {
	return nil, HashCodeMismatch(ListType)
}

func (ls *list) Eq(ev Eval, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *list:
		return valuesEq(ev, ls.values, t.values)
	default:
		return False, nil
	}
}

func (ls *list) Get(ev Eval, index Value) (Value, Error) {
	idx, err := boundedIndex(index, len(ls.values))
	if err != nil {
		return nil, err
	}
	return ls.values[idx], nil
}

func (ls *list) Set(ev Eval, index Value, val Value) Error {
	if ls.frozen {
		return ImmutableValue()
	}

	idx, err := boundedIndex(index, len(ls.values))
	if err != nil {
		return err
	}

	ls.values[idx] = val
	return nil
}

func (ls *list) Len(ev Eval) (Int, Error) {
	return NewInt(int64(len(ls.values))), nil
}

func (ls *list) Slice(ev Eval, from Value, to Value) (Value, Error) {

	f, t, err := sliceIndices(from, to, len(ls.values))
	if err != nil {
		return nil, err
	}

	result := NewList(CopyValues(ls.values[f:t]))
	if ls.frozen {
		result.(*list).frozen = true
	}
	return result, nil
}

func (ls *list) SliceFrom(ev Eval, from Value) (Value, Error) {
	return ls.Slice(ev, from, NewInt(int64(len(ls.values))))
}

func (ls *list) SliceTo(ev Eval, to Value) (Value, Error) {
	return ls.Slice(ev, Zero, to)
}

func (ls *list) Values() []Value {
	return ls.values
}

//------------------------------------------------------

func (ls *list) Contains(ev Eval, val Value) (Bool, Error) {

	idx, err := ls.Index(ev, val)
	if err != nil {
		return nil, err
	}

	eq, err := idx.Eq(ev, NegOne)
	if err != nil {
		return nil, err
	}

	return NewBool(!eq.BoolVal()), nil
}

func (ls *list) Index(ev Eval, val Value) (Int, Error) {
	for i, v := range ls.values {
		eq, err := val.Eq(ev, v)
		if err != nil {
			return nil, err
		}
		if eq.BoolVal() {
			return NewInt(int64(i)), nil
		}
	}
	return NegOne, nil
}

func (ls *list) IsEmpty() Bool {
	return NewBool(len(ls.values) == 0)
}

func (ls *list) Join(ev Eval, delim Str) (Str, Error) {

	result := make([]string, len(ls.values))
	for i, v := range ls.values {

		s, err := v.ToStr(ev)
		if err != nil {
			return nil, err
		}
		result[i] = s.String()
	}

	return NewStr(strings.Join(result, delim.String()))
}

func (ls *list) Map(ev Eval, mapper Mapper) (List, Error) {

	vals := make([]Value, len(ls.values))

	var err Error
	for i, v := range ls.values {
		vals[i], err = mapper(ev, v)
		if err != nil {
			return nil, err
		}
	}

	return NewList(vals), nil
}

func (ls *list) Reduce(ev Eval, initial Value, reducer Reducer) (Value, Error) {

	acc := initial
	var err Error

	for _, v := range ls.values {
		acc, err = reducer(ev, acc, v)
		if err != nil {
			return nil, err
		}
	}

	return acc, nil
}

func (ls *list) Filter(ev Eval, filterer Filterer) (List, Error) {

	vals := []Value{}

	for _, v := range ls.values {
		flt, err := filterer(ev, v)
		if err != nil {
			return nil, err
		}
		pred, ok := flt.(Bool)
		if !ok {
			return nil, TypeMismatch(BoolType, flt.Type())
		}

		eq, err := pred.Eq(ev, True)
		if err != nil {
			return nil, err
		}
		if eq.BoolVal() {
			vals = append(vals, v)
		}
	}

	return NewList(vals), nil
}

func (ls *list) Add(ev Eval, val Value) (List, Error) {
	if ls.frozen {
		return nil, ImmutableValue()
	}

	ls.values = append(ls.values, val)
	return ls, nil
}

func (ls *list) AddAll(ev Eval, val Value) (List, Error) {
	if ls.frozen {
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
		ls.values = append(ls.values, v)

		b, err = itr.IterNext(ev)
		if err != nil {
			return nil, err
		}
	}
	return ls, nil
}

func (ls *list) Remove(ev Eval, index Int) (List, Error) {
	if ls.frozen {
		return nil, ImmutableValue()
	}

	n := int(index.ToInt())
	if n < 0 || n >= len(ls.values) {
		return nil, IndexOutOfBounds(n)
	}
	ls.values = append(ls.values[:n], ls.values[n+1:]...)
	return ls, nil
}

func sortValues(vals []Value, cmp func(i, j int) bool) (err Error) {

	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	sort.Slice(vals, cmp)

	return
}

// DefaultLesser returns whether the first param is less than the second param
// by casting both parameters to Comparable, and callingCmp().
var DefaultLesser = func(ev Eval, a Value, b Value) (Bool, Error) {

	ca, ok := a.(Comparable)
	if !ok {
		return nil, Error(fmt.Errorf("TypeMismatch: Type %s cannot be sorted", a.Type()))
	}

	cb, ok := b.(Comparable)
	if !ok {
		return nil, Error(fmt.Errorf("TypeMismatch: Type %s cannot be sorted", b.Type()))
	}

	n, err := ca.Cmp(ev, cb)
	if err != nil {
		return nil, err
	}

	return NewBool(n.ToInt() < 0), nil
}

func (ls *list) Sort(ev Eval, lesser Lesser) (List, Error) {
	if ls.frozen {
		return nil, ImmutableValue()
	}

	err := sortValues(ls.values, func(i, j int) bool {
		b, err := lesser(ev, ls.values[i], ls.values[j])
		if err != nil {
			panic(err)
		}
		return b.BoolVal()
	})
	if err != nil {
		return nil, err
	}
	return ls, nil
}

func (ls *list) Clear() (List, Error) {
	if ls.frozen {
		return nil, ImmutableValue()
	}

	ls.values = []Value{}
	return ls, nil
}

//---------------------------------------------------------------
// Iterator

type listIterator struct {
	Struct
	ls *list
	n  int
}

func (ls *list) NewIterator(ev Eval) (Iterator, Error) {

	itr := &listIterator{iteratorStruct(), ls, -1}

	next, get := iteratorFields(ev, itr)
	itr.Internal("next", next)
	itr.Internal("get", get)

	return itr, nil
}

func (i *listIterator) IterNext(ev Eval) (Bool, Error) {
	i.n++
	return NewBool(i.n < len(i.ls.values)), nil
}

func (i *listIterator) IterGet(ev Eval) (Value, Error) {
	if (i.n >= 0) && (i.n < len(i.ls.values)) {
		return i.ls.values[i.n], nil
	}
	return nil, NoSuchElement()
}

//--------------------------------------------------------------
// fields

/*doc
A List has the following fields:

* [add](#add)
* [addAll](#addall)
* [clear](#clear)
* [contains](#contains)
* [filter](#filter)
* [index](#index)
* [isEmpty](#isempty)
* [join](#join)
* [map](#map)
* [reduce](#reduce)
* [remove](#remove)
* [sort](#sort)

*/

var listMethods = map[string]Method{

	/*doc
	### `add`

	`add` adds a value to the end of the list, and returns the modified list.

	* signature: `add(val <Value>) <List>`
	* example:

	```
	let a = [1, 2, 3]
	println(a.add(4))
	```

	*/
	"add": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			ls := self.(List)
			return ls.Add(ev, params[0])
		}),

	/*doc
	### `addAll`

	`addAll` adds all of the values in the given [Iterable](interfaces.html#iterable) to
	the end of the list, and returns the modified list.

	* signature: `addAll(itr <Iterable>) <List>`
	* example:

	```
	let a = [1, 2]
	println(a.addAll([3, 4]))
	```

	*/
	"addAll": NewFixedMethod(
		[]Type{AnyType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			ls := self.(List)
			return ls.AddAll(ev, params[0])
		}),

	/*doc
	### `clear`

	`clear` removes all of the values from the list, and returns the empty list.

	* signature: `clear() <List>`
	* example:

	```
	let a = [1, 2]
	println(a.clear())
	```

	*/
	"clear": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			ls := self.(List)
			return ls.Clear()
		}),

	/*doc
	### `contains`

	`contains` returns whether the given value is an element in the list.

	* signature: `contains(val <Value>) <Bool>`
	* example:

	```
	let a = [1, 2]
	println(a.contains(2))
	```

	*/
	"contains": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			ls := self.(List)
			return ls.Contains(ev, params[0])
		}),

	/*doc
	### `filter`

	`filter` returns a new list by passing each of the elements of the current list
	into the given predicate.  If the predicate returns `true` for an element, that
	element is added to the new list.  The original list is unmodified.

	The predicate must accept one parameter of any type, and return a Bool.

	* signature: `filter(predicate <Func>) <List>`
	* predicate signature: `fn(val <Value>) <Bool>`
	* example:

	```
	let a = [1, 2, 3, 4, 5]
	println(a.filter(|e| => e % 2 == 0))
	```

	*/
	"filter": NewFixedMethod(
		[]Type{FuncType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			ls := self.(List)

			// check arity
			fn := params[0].(Func)
			expected := Arity{FixedArity, 1, 0}
			if fn.Arity() != expected {
				return nil, fmt.Errorf(
					"ArityMismatch: filter function must have 1 parameter")
			}

			// invoke
			return ls.Filter(ev, func(ev Eval, v Value) (Bool, Error) {
				val, err := fn.Invoke(ev, []Value{v})
				if err != nil {
					return nil, err
				}

				result, ok := val.(Bool)
				if !ok {
					return nil, fmt.Errorf(
						"TypeMismatch: filter function must return Bool, not %s", val.Type())
				}
				return result, nil
			})
		}),

	/*doc
	### `index`

	`index` returns the index of the given value in the list, or -1 if the value
	is not contained in the list.

	* signature: `index(val <Value>) <Int>`
	* example:

	```
	let a = ['x', 'y', 'z']
	println(a.index('z'))
	```

	*/
	"index": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			ls := self.(List)
			return ls.Index(ev, params[0])
		}),

	/*doc
	### `isEmpty`

	`isEmpty` returns whether the list contains any values.

	* signature: `isEmpty() <Bool>`
	* example: `println([].isEmpty())`

	*/
	"isEmpty": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			ls := self.(List)
			return ls.IsEmpty(), nil
		}),

	/*doc
	### `join`

	Join concatenates the [str()](builtins.html#str) representation of the elements of the list to
	create a single string.  The separator string sep is placed between elements in
	the resulting string.  The sep parameter is optional, and defaults
	to the empty string `''`.

	* signature: `join(sep = '' <Str>) <Str>`
	* example:

	```
	println([1,2,3].join(':'))
	```

	*/
	"join": NewMultipleMethod(
		[]Type{},
		[]Type{StrType},
		false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			ls := self.(List)

			delim, err := NewStr("")
			if err != nil {
				return nil, err
			}

			if len(params) == 1 {
				delim = params[0].(Str)
			}

			return ls.Join(ev, delim)
		}),

	/*doc
	### `map`

	`map` returns a copy of the list with all its elements modified according to
	the mapping function.  The original list is unmodified.

	The mapping function must accept one value, and must return one value.

	* signature: `map(mapping <Func>) <List>`
	* mapping signature: `fn(val <Value>) <Value>`
	* example:

	```
	let a = [1,2,3]
	let b = a.map(|e| => e * e)
	println(b)
	```

	*/
	"map": NewFixedMethod(
		[]Type{FuncType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			ls := self.(List)

			// check arity
			fn := params[0].(Func)
			expected := Arity{FixedArity, 1, 0}
			if fn.Arity() != expected {
				return nil, fmt.Errorf(
					"ArityMismatch: map function must have 1 parameter")
			}

			// invoke
			return ls.Map(ev, func(ev Eval, v Value) (Value, Error) {
				return fn.Invoke(ev, []Value{v})
			})
		}),

	/*doc
	### `reduce`

	`reduce` reduces the list to a single value, by applying a "reducer" function
	to an accumulated value and each element in the list.
	Accumulation is done starting with the first element in the list,
	and ending with the last.  The original list is unmodified.

	The reducer function must accept two values, and return one value.

	* signature: `reduce(start <Value>, reducer <Func>) <List>`
	* reducer signature: `fn(accum <Value>, val <Value>) <Value>`
	* example:

	```
	let a = [1,2,3]
	let b = a.reduce(0, |acc, e| => acc + e)
	println(b)
	```

	*/
	"reduce": NewFixedMethod(
		[]Type{AnyType, FuncType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			ls := self.(List)

			initial := params[0]

			if params[1] == Null {
				return nil, NullValueError()
			}

			// check arity
			fn := params[1].(Func)
			expected := Arity{FixedArity, 2, 0}
			if fn.Arity() != expected {
				return nil, fmt.Errorf(
					"ArityMismatch: reduce function must have 2 parameters")
			}

			// invoke
			return ls.Reduce(ev, initial, func(ev Eval, acc Value, v Value) (Value, Error) {
				return fn.Invoke(ev, []Value{acc, v})
			})
		}),

	/*doc
	### `remove`

	`remove` remove the value at the given index from the list, and returns the
	modified list.

	* signature: `remove(index <Int>) <List>`
	* example: `println(['a','b','c'].remove(2))`

	*/
	"remove": NewFixedMethod(
		[]Type{IntType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			ls := self.(List)
			return ls.Remove(ev, params[0].(Int))
		}),

	/*doc

	### `sort`

	`sort` sorts the elements in the list and returns the modified list.  If the
	optional "lesser" function is provided, it is used to compare values in the list.
	If the lesser function is not provided, then the `<` operator is used.

	* signature: `sort(lesser = null <Func>) <List>`
	* lesser signature: `fn(val <Value>, val <Value>) <Bool>`
	* example:

	```
	let a = [7, 4, 11, 13, 6, 2, 9, 1]
	a.sort(|a, b| => b < a) // sort in reverse
	println(a)
	```

	*/
	"sort": NewMultipleMethod(
		[]Type{},
		[]Type{FuncType},
		false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			ls := self.(List)

			// if no function was provided, just use the default Lesser
			if len(params) == 0 {
				return ls.Sort(ev, DefaultLesser)
			}

			// check arity
			fn := params[0].(Func)
			expected := Arity{FixedArity, 2, 0}
			if fn.Arity() != expected {
				return nil, fmt.Errorf(
					"ArityMismatch: sort function must have 2 parameters")
			}

			// invoke
			return ls.Sort(ev, func(ev Eval, a Value, b Value) (Bool, Error) {
				val, err := fn.Invoke(ev, []Value{a, b})
				if err != nil {
					return nil, err
				}

				result, ok := val.(Bool)
				if !ok {
					return nil, fmt.Errorf(
						"TypeMismatch: sort function must return Bool, not %s", val.Type())
				}
				return result, nil
			})
		}),
}

func (ls *list) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(listMethods))
	for name := range listMethods {
		names = append(names, name)
	}
	return names, nil
}

func (ls *list) HasField(name string) (bool, Error) {
	_, ok := listMethods[name]
	return ok, nil
}

func (ls *list) GetField(ev Eval, name string) (Value, Error) {
	if method, ok := listMethods[name]; ok {
		return method.ToFunc(ls, name), nil
	}
	return nil, NoSuchField(name)
}

func (ls *list) InvokeField(ev Eval, name string, params []Value) (Value, Error) {

	if method, ok := listMethods[name]; ok {
		return method.Invoke(ls, ev, params)
	}
	return nil, NoSuchField(name)
}
