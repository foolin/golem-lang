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

func (ls *list) Freeze() (Value, Error) {
	ls.frozen = true
	return ls, nil
}

func (ls *list) Frozen() (Bool, Error) {
	return NewBool(ls.frozen), nil
}

func (ls *list) ToStr(cx Context) Str {

	var buf bytes.Buffer
	buf.WriteString("[")
	for idx, v := range ls.array {
		if idx > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(" ")
		buf.WriteString(v.ToStr(cx).String())
	}
	buf.WriteString(" ]")
	return NewStr(buf.String())
}

func (ls *list) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (ls *list) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *list:
		return valuesEq(cx, ls.array, t.array)
	default:
		return False, nil
	}
}

func (ls *list) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (ls *list) Get(cx Context, index Value) (Value, Error) {
	idx, err := boundedIndex(index, len(ls.array))
	if err != nil {
		return nil, err
	}
	return ls.array[idx], nil
}

func (ls *list) Contains(cx Context, val Value) (Bool, Error) {

	idx, err := ls.IndexOf(cx, val)
	if err != nil {
		return nil, err
	}

	eq, err := idx.Eq(cx, NegOne)
	if err != nil {
		return nil, err
	}

	return NewBool(!eq.BoolVal()), nil
}

func (ls *list) IndexOf(cx Context, val Value) (Int, Error) {
	for i, v := range ls.array {
		eq, err := val.Eq(cx, v)
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

func (ls *list) Len() Int {
	return NewInt(int64(len(ls.array)))
}

func (ls *list) Slice(cx Context, from Value, to Value) (Value, Error) {

	f, t, err := sliceIndices(from, to, len(ls.array))
	if err != nil {
		return nil, err
	}

	a := ls.array[f:t]
	b := make([]Value, len(a))
	copy(b, a)
	return NewList(b), nil
}

func (ls *list) SliceFrom(cx Context, from Value) (Value, Error) {
	return ls.Slice(cx, from, NewInt(int64(len(ls.array))))
}

func (ls *list) SliceTo(cx Context, to Value) (Value, Error) {
	return ls.Slice(cx, Zero, to)
}

func (ls *list) Values() []Value {
	return ls.array
}

func (ls *list) Join(cx Context, delim Str) Str {

	s := make([]string, len(ls.array))
	for i, v := range ls.array {
		s[i] = v.ToStr(cx).String()
	}

	return NewStr(strings.Join(s, delim.ToStr(cx).String()))
}

func (ls *list) Map(cx Context, mapper func(Value) (Value, Error)) (Value, Error) {

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

func (ls *list) Reduce(cx Context, initial Value, reducer func(Value, Value) (Value, Error)) (Value, Error) {

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

func (ls *list) Filter(cx Context, filterer func(Value) (Value, Error)) (Value, Error) {

	vals := []Value{}

	for _, v := range ls.array {
		flt, err := filterer(v)
		if err != nil {
			return nil, err
		}
		pred, ok := flt.(Bool)
		if !ok {
			return nil, TypeMismatchError("Expected Bool")
		}

		eq, err := pred.Eq(cx, True)
		if err != nil {
			return nil, err
		}
		if eq.BoolVal() {
			vals = append(vals, v)
		}
	}

	return NewList(vals), nil
}

//---------------------------------------------------------------
// Mutation

func (ls *list) Set(cx Context, index Value, val Value) Error {
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

func (ls *list) Add(cx Context, val Value) Error {
	if ls.frozen {
		return ImmutableValueError()
	}

	ls.array = append(ls.array, val)
	return nil
}

func (ls *list) AddAll(cx Context, val Value) Error {
	if ls.frozen {
		return ImmutableValueError()
	}

	if ibl, ok := val.(Iterable); ok {
		itr := ibl.NewIterator(cx)
		for itr.IterNext().BoolVal() {
			v, err := itr.IterGet()
			if err != nil {
				return err
			}
			ls.array = append(ls.array, v)
		}
		return nil
	}
	return TypeMismatchError("Expected Iterable Type")
}

func (ls *list) Remove(cx Context, index Int) Error {
	if ls.frozen {
		return ImmutableValueError()
	}

	n := int(index.IntVal())
	if n < 0 || n >= len(ls.array) {
		return IndexOutOfBoundsError(n)
	}
	ls.array = append(ls.array[:n], ls.array[n+1:]...)
	return nil
}

func (ls *list) Clear() Error {
	if ls.frozen {
		return ImmutableValueError()
	}

	ls.array = []Value{}
	return nil
}

//---------------------------------------------------------------
// Iterator

type listIterator struct {
	Struct
	ls *list
	n  int
}

func (ls *list) NewIterator(cx Context) Iterator {
	return initIteratorStruct(cx,
		&listIterator{newIteratorStruct(), ls, -1})
}

func (i *listIterator) IterNext() Bool {
	i.n++
	return NewBool(i.n < len(i.ls.array))
}

func (i *listIterator) IterGet() (Value, Error) {
	if (i.n >= 0) && (i.n < len(i.ls.array)) {
		return i.ls.array[i.n], nil
	}
	return nil, NoSuchElementError()
}

//--------------------------------------------------------------
// intrinsic functions

func (ls *list) GetField(cx Context, key Str) (Value, Error) {
	switch sn := key.String(); sn {

	case "add":
		return &intrinsicFunc{ls, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				err := ls.Add(cx, values[0])
				if err != nil {
					return nil, err
				}
				return ls, nil
			}}}, nil

	case "addAll":
		return &intrinsicFunc{ls, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				err := ls.AddAll(cx, values[0])
				if err != nil {
					return nil, err
				}
				return ls, nil
			}}}, nil

	case "remove":
		return &intrinsicFunc{ls, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				index, ok := values[0].(Int)
				if !ok {
					return nil, TypeMismatchError("Expected Iterable Type")
				}

				err := ls.Remove(cx, index)
				if err != nil {
					return nil, err
				}
				return ls, nil
			}}}, nil

	case "clear":
		return &intrinsicFunc{ls, sn, &nativeFunc{
			0, 0,
			func(cx Context, values []Value) (Value, Error) {
				err := ls.Clear()
				if err != nil {
					return nil, err
				}
				return ls, nil
			}}}, nil

	case "isEmpty":
		return &intrinsicFunc{ls, sn, &nativeFunc{
			0, 0,
			func(cx Context, values []Value) (Value, Error) {
				return ls.IsEmpty(), nil
			}}}, nil

	case "contains":
		return &intrinsicFunc{ls, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				return ls.Contains(cx, values[0])
			}}}, nil

	case "indexOf":
		return &intrinsicFunc{ls, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				return ls.IndexOf(cx, values[0])
			}}}, nil

	case "join":
		return &intrinsicFunc{ls, sn, &nativeFunc{
			0, 1,
			func(cx Context, values []Value) (Value, Error) {
				var delim Str
				switch len(values) {
				case 0:
					delim = NewStr("")
				case 1:
					if s, ok := values[0].(Str); ok {
						delim = s
					} else {
						return nil, TypeMismatchError("Expected Str")
					}
				default:
					panic("arity mismatch")
				}

				return ls.Join(cx, delim), nil
			}}}, nil

	case "map":
		return &intrinsicFunc{ls, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {

				if f, ok := values[0].(Func); ok {
					return ls.Map(cx, func(v Value) (Value, Error) {
						return f.Invoke(cx, []Value{v})
					})
				}
				return nil, TypeMismatchError("Expected Func")

			}}}, nil

	case "reduce":
		return &intrinsicFunc{ls, sn, &nativeFunc{
			2, 2,
			func(cx Context, values []Value) (Value, Error) {

				initial := values[0]
				if f, ok := values[1].(Func); ok {
					return ls.Reduce(cx, initial, func(acc Value, v Value) (Value, Error) {
						return f.Invoke(cx, []Value{acc, v})
					})
				}
				return nil, TypeMismatchError("Expected Func")
			}}}, nil

	case "filter":
		return &intrinsicFunc{ls, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {

				if f, ok := values[0].(Func); ok {
					return ls.Filter(cx, func(v Value) (Value, Error) {
						return f.Invoke(cx, []Value{v})
					})
				}
				return nil, TypeMismatchError("Expected Func")

			}}}, nil

	default:
		return nil, NoSuchFieldError(key.String())
	}
}
