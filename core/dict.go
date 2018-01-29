// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
)

type dict struct {
	hashMap *HashMap
	frozen  bool
}

func NewDict(cx Context, entries []*HEntry) Dict {
	return &dict{NewHashMap(cx, entries), false}
}

func (d *dict) compositeMarker() {}

func (d *dict) Type() Type { return TDICT }

func (d *dict) Freeze() (Value, Error) {
	d.frozen = true
	return d, nil
}

func (d *dict) Frozen() (Bool, Error) {
	return MakeBool(d.frozen), nil
}

func (d *dict) ToStr(cx Context) Str {

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
		s := entry.Key.ToStr(cx)
		buf.WriteString(s.String())

		buf.WriteString(": ")
		s = entry.Value.ToStr(cx)
		buf.WriteString(s.String())
	}

	buf.WriteString(" }")
	return NewStr(buf.String())
}

func (d *dict) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (d *dict) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *dict:
		return d.hashMap.Eq(cx, t.hashMap)
	default:
		return FALSE, nil
	}
}

func (d *dict) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (d *dict) Get(cx Context, key Value) (Value, Error) {
	return d.hashMap.Get(cx, key)
}

func (d *dict) Len() Int {
	return d.hashMap.Len()
}

func (d *dict) IsEmpty() Bool {
	return MakeBool(d.hashMap.Len().IntVal() == 0)
}

func (d *dict) ContainsKey(cx Context, key Value) (Bool, Error) {
	return d.hashMap.ContainsKey(cx, key)
}

//---------------------------------------------------------------
// Mutation

func (d *dict) Set(cx Context, key Value, val Value) Error {
	if d.frozen {
		return ImmutableValueError()
	}

	return d.hashMap.Put(cx, key, val)
}

func (d *dict) Clear() Error {
	if d.frozen {
		return ImmutableValueError()
	}

	d.hashMap = EmptyHashMap()
	return nil
}

func (d *dict) Remove(cx Context, key Value) (Bool, Error) {
	if d.frozen {
		return nil, ImmutableValueError()
	}

	return d.hashMap.Remove(cx, key)
}

func (d *dict) AddAll(cx Context, val Value) Error {
	if d.frozen {
		return ImmutableValueError()
	}

	if ibl, ok := val.(Iterable); ok {
		itr := ibl.NewIterator(cx)
		for itr.IterNext().BoolVal() {
			v, err := itr.IterGet()
			if err != nil {
				return err
			}
			if tp, ok := v.(tuple); ok {
				if len(tp) == 2 {
					d.hashMap.Put(cx, tp[0], tp[1])
				} else {
					return TupleLengthError(2, len(tp))
				}
			} else {
				return TypeMismatchError("Expected Tuple")
			}
		}
		return nil
	} else {
		return TypeMismatchError("Expected Iterable Type")
	}
}

//---------------------------------------------------------------
// Iterator

type dictIterator struct {
	Struct
	d       *dict
	itr     *HIterator
	hasNext bool
}

func (d *dict) NewIterator(cx Context) Iterator {
	return initIteratorStruct(cx,
		&dictIterator{newIteratorStruct(), d, d.hashMap.Iterator(), false})
}

func (i *dictIterator) IterNext() Bool {
	i.hasNext = i.itr.Next()
	return MakeBool(i.hasNext)
}

func (i *dictIterator) IterGet() (Value, Error) {

	if i.hasNext {
		entry := i.itr.Get()
		return NewTuple([]Value{entry.Key, entry.Value}), nil
	} else {
		return nil, NoSuchElementError()
	}
}

//--------------------------------------------------------------
// intrinsic functions

func (d *dict) GetField(cx Context, key Str) (Value, Error) {
	switch sn := key.String(); sn {

	case "addAll":
		return &intrinsicFunc{d, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				err := d.AddAll(cx, values[0])
				if err != nil {
					return nil, err
				} else {
					return d, nil
				}
			}}}, nil

	case "clear":
		return &intrinsicFunc{d, sn, &nativeFunc{
			0, 0,
			func(cx Context, values []Value) (Value, Error) {
				err := d.Clear()
				if err != nil {
					return nil, err
				} else {
					return d, nil
				}
			}}}, nil

	case "isEmpty":
		return &intrinsicFunc{d, sn, &nativeFunc{
			0, 0,
			func(cx Context, values []Value) (Value, Error) {
				return d.IsEmpty(), nil
			}}}, nil

	case "containsKey":
		return &intrinsicFunc{d, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				return d.ContainsKey(cx, values[0])
			}}}, nil

	case "remove":
		return &intrinsicFunc{d, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				return d.Remove(cx, values[0])
			}}}, nil

	default:
		return nil, NoSuchFieldError(key.String())
	}
}
