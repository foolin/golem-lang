// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"strings"
)

//---------------------------------------------------------------
// list

type list struct {
	array  []Value
	frozen bool
}

// NewList creates a new List
func NewList(values []Value) List {
	return &list{values, false}
}

func (ls *list) compositeMarker() {}

func (ls *list) Type() Type { return ListType }

func (ls *list) Freeze(ev Evaluator) (Value, Error) {
	ls.frozen = true
	return ls, nil
}

func (ls *list) Frozen(ev Evaluator) (Bool, Error) {
	return NewBool(ls.frozen), nil
}

func (ls *list) ToStr(ev Evaluator) (Str, Error) {

	var buf bytes.Buffer
	buf.WriteString("[")
	for idx, v := range ls.array {
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

	return NewStr(buf.String()), nil
}

func (ls *list) HashCode(ev Evaluator) (Int, Error) {
	return nil, HashCodeMismatchError(ListType)
}

func (ls *list) Eq(ev Evaluator, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *list:
		return valuesEq(ev, ls.array, t.array)
	default:
		return False, nil
	}
}

func (ls *list) Get(ev Evaluator, index Value) (Value, Error) {
	idx, err := boundedIndex(index, len(ls.array))
	if err != nil {
		return nil, err
	}
	return ls.array[idx], nil
}

func (ls *list) Set(ev Evaluator, index Value, val Value) Error {
	if ls.frozen {
		return ImmutableValueError()
	}

	idx, err := boundedIndex(index, len(ls.array))
	if err != nil {
		return err
	}

	ls.array[idx] = val
	return nil
}

func (ls *list) Len(ev Evaluator) (Int, Error) {
	return NewInt(int64(len(ls.array))), nil
}

func (ls *list) Slice(ev Evaluator, from Value, to Value) (Value, Error) {

	f, t, err := sliceIndices(from, to, len(ls.array))
	if err != nil {
		return nil, err
	}

	a := ls.array[f:t]
	b := make([]Value, len(a))
	copy(b, a)
	return NewList(b), nil
}

func (ls *list) SliceFrom(ev Evaluator, from Value) (Value, Error) {
	return ls.Slice(ev, from, NewInt(int64(len(ls.array))))
}

func (ls *list) SliceTo(ev Evaluator, to Value) (Value, Error) {
	return ls.Slice(ev, Zero, to)
}

func (ls *list) Values() []Value {
	return ls.array
}

//------------------------------------------------------

func (ls *list) Contains(ev Evaluator, val Value) (Bool, Error) {

	idx, err := ls.IndexOf(ev, val)
	if err != nil {
		return nil, err
	}

	eq, err := idx.Eq(ev, NegOne)
	if err != nil {
		return nil, err
	}

	return NewBool(!eq.BoolVal()), nil
}

func (ls *list) IndexOf(ev Evaluator, val Value) (Int, Error) {
	for i, v := range ls.array {
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
	return NewBool(len(ls.array) == 0)
}

func (ls *list) Join(ev Evaluator, delim Str) (Str, Error) {

	result := make([]string, len(ls.array))
	for i, v := range ls.array {

		s, err := v.ToStr(ev)
		if err != nil {
			return nil, err
		}
		result[i] = s.String()
	}

	return NewStr(strings.Join(result, delim.String())), nil
}

func (ls *list) Map(ev Evaluator, mapper func(Value) (Value, Error)) (List, Error) {

	vals := make([]Value, len(ls.array))

	var err Error
	for i, v := range ls.array {
		vals[i], err = mapper(v)
		if err != nil {
			return nil, err
		}
	}

	return NewList(vals), nil
}

func (ls *list) Reduce(
	ev Evaluator, initial Value, reducer func(Value, Value) (Value, Error)) (Value, Error) {

	acc := initial
	var err Error

	for _, v := range ls.array {
		acc, err = reducer(acc, v)
		if err != nil {
			return nil, err
		}
	}

	return acc, nil
}

func (ls *list) Filter(ev Evaluator, filterer func(Value) (Value, Error)) (List, Error) {

	vals := []Value{}

	for _, v := range ls.array {
		flt, err := filterer(v)
		if err != nil {
			return nil, err
		}
		pred, ok := flt.(Bool)
		if !ok {
			return nil, TypeMismatchError(BoolType, flt.Type())
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

func (ls *list) Add(ev Evaluator, val Value) (List, Error) {
	if ls.frozen {
		return nil, ImmutableValueError()
	}

	ls.array = append(ls.array, val)
	return ls, nil
}

func (ls *list) AddAll(ev Evaluator, val Value) (List, Error) {
	if ls.frozen {
		return nil, ImmutableValueError()
	}

	ibl, ok := val.(Iterable)
	if !ok {
		return nil, IterableMismatchError(val.Type())
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
		ls.array = append(ls.array, v)

		b, err = itr.IterNext(ev)
		if err != nil {
			return nil, err
		}
	}
	return ls, nil
}

func (ls *list) Remove(index Int) (List, Error) {
	if ls.frozen {
		return nil, ImmutableValueError()
	}

	n := int(index.IntVal())
	if n < 0 || n >= len(ls.array) {
		return nil, IndexOutOfBoundsError(n)
	}
	ls.array = append(ls.array[:n], ls.array[n+1:]...)
	return ls, nil
}

func (ls *list) Clear() (List, Error) {
	if ls.frozen {
		return nil, ImmutableValueError()
	}

	ls.array = []Value{}
	return ls, nil
}

//---------------------------------------------------------------
// Iterator

type listIterator struct {
	Struct
	ls *list
	n  int
}

func (ls *list) NewIterator(ev Evaluator) (Iterator, Error) {

	itr := &listIterator{iteratorStruct(), ls, -1}

	next, get := iteratorFields(ev, itr)
	itr.Internal("next", next)
	itr.Internal("get", get)

	return itr, nil
}

func (i *listIterator) IterNext(ev Evaluator) (Bool, Error) {
	i.n++
	return NewBool(i.n < len(i.ls.array)), nil
}

func (i *listIterator) IterGet(ev Evaluator) (Value, Error) {
	if (i.n >= 0) && (i.n < len(i.ls.array)) {
		return i.ls.array[i.n], nil
	}
	return nil, NoSuchElementError()
}

//--------------------------------------------------------------
// fields

var listMethods = map[string]Method{

	"add": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			ls := self.(List)
			return ls.Add(ev, params[0])
		}),

	"addAll": NewFixedMethod(
		[]Type{AnyType}, false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			ls := self.(List)
			return ls.AddAll(ev, params[0])
		}),

	"clear": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			ls := self.(List)
			return ls.Clear()
		}),

	"isEmpty": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			ls := self.(List)
			return ls.IsEmpty(), nil
		}),

	"contains": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			ls := self.(List)
			return ls.Contains(ev, params[0])
		}),

	"indexOf": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			ls := self.(List)
			return ls.IndexOf(ev, params[0])
		}),

	"remove": NewFixedMethod(
		[]Type{IntType}, true,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			ls := self.(List)
			return ls.Remove(params[0].(Int))
		}),

	"join": NewMultipleMethod(
		[]Type{},
		[]Type{StrType},
		false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			ls := self.(List)
			delim := NewStr("")
			if len(params) == 1 {
				delim = params[0].(Str)
			}

			return ls.Join(ev, delim)
		}),

	"map": NewFixedMethod(
		[]Type{FuncType}, false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			ls := self.(List)

			fn := params[0].(Func)
			return ls.Map(ev, func(v Value) (Value, Error) {
				return fn.Invoke(ev, []Value{v})
			})
		}),

	"reduce": NewFixedMethod(
		[]Type{AnyType, FuncType}, true,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			ls := self.(List)

			if params[1] == Null {
				return nil, NullValueError()
			}

			initial := params[0]
			fn := params[1].(Func)
			return ls.Reduce(ev, initial, func(acc Value, v Value) (Value, Error) {
				return fn.Invoke(ev, []Value{acc, v})
			})
		}),

	"filter": NewFixedMethod(
		[]Type{FuncType}, false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			ls := self.(List)

			fn := params[0].(Func)
			return ls.Filter(ev, func(v Value) (Value, Error) {
				return fn.Invoke(ev, []Value{v})
			})
		}),
}

func (ls *list) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(listMethods))
	for name, _ := range listMethods {
		names = append(names, name)
	}
	return names, nil
}

func (ls *list) HasField(name string) (bool, Error) {
	_, ok := listMethods[name]
	return ok, nil
}

func (ls *list) GetField(name string, ev Evaluator) (Value, Error) {
	if method, ok := listMethods[name]; ok {
		return method.ToFunc(ls, name), nil
	}
	return nil, NoSuchFieldError(name)
}

func (ls *list) InvokeField(name string, ev Evaluator, params []Value) (Value, Error) {

	if method, ok := listMethods[name]; ok {
		return method.Invoke(ls, ev, params)
	}
	return nil, NoSuchFieldError(name)
}
