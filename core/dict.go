// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"bytes"
	"reflect"
)

type dict struct {
	hashMap *HashMap
}

func NewDict(entries []*HEntry) Dict {
	return &dict{NewHashMap(entries)}
}

func (d *dict) compositeMarker() {}

func (d *dict) TypeOf() Type { return TDICT }

func (d *dict) ToStr() Str {

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
		s := entry.Key.ToStr()
		buf.WriteString(s.String())

		buf.WriteString(": ")
		s = entry.Value.ToStr()
		buf.WriteString(s.String())
	}

	buf.WriteString(" }")
	return MakeStr(buf.String())
}

func (d *dict) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (d *dict) Eq(v Value) Bool {
	switch t := v.(type) {
	case *dict:
		return MakeBool(reflect.DeepEqual(d.hashMap, t.hashMap))
	default:
		return FALSE
	}
}

func (d *dict) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (d *dict) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(d, t), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (d *dict) Get(key Value) (Value, Error) {
	return d.hashMap.Get(key)
}

func (d *dict) Set(key Value, val Value) Error {
	return d.hashMap.Put(key, val)
}

func (d *dict) Len() Int {
	return d.hashMap.Len()
}

func (d *dict) Clear() {
	d.hashMap = EmptyHashMap()
}

func (d *dict) IsEmpty() Bool {
	return MakeBool(d.hashMap.Len().IntVal() == 0)
}

func (d *dict) ContainsKey(key Value) (Bool, Error) {
	return d.hashMap.ContainsKey(key)
}

func (d *dict) AddAll(val Value) Error {
	if ibl, ok := val.(Iterable); ok {
		itr := ibl.NewIterator()
		for itr.IterNext().BoolVal() {
			v, err := itr.IterGet()
			if err != nil {
				return err
			}
			if tp, ok := v.(tuple); ok {
				if len(tp) == 2 {
					d.hashMap.Put(tp[0], tp[1])
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

func (d *dict) NewIterator() Iterator {

	stc, err := NewStruct([]*StructEntry{
		{"nextValue", true, false, NULL},
		{"getValue", true, false, NULL}})
	if err != nil {
		panic("invalid struct")
	}

	itr := &dictIterator{stc, d, d.hashMap.Iterator(), false}

	stc.InitField(MakeStr("nextValue"), &nativeFunc{
		func(values []Value) (Value, Error) {
			return itr.IterNext(), nil
		}})
	stc.InitField(MakeStr("getValue"), &nativeFunc{
		func(values []Value) (Value, Error) {
			return itr.IterGet()
		}})

	return itr
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

func (d *dict) GetField(key Str) (Value, Error) {
	switch key.String() {

	case "addAll":
		return &intrinsicFunc{d, "addAll", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 1 {
					return nil, ArityMismatchError("1", len(values))
				}
				err := d.AddAll(values[0])
				if err != nil {
					return nil, err
				} else {
					return d, nil
				}
			}}}, nil

	case "clear":
		return &intrinsicFunc{d, "clear", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 0 {
					return nil, ArityMismatchError("0", len(values))
				}
				d.Clear()
				return d, nil
			}}}, nil

	case "isEmpty":
		return &intrinsicFunc{d, "isEmpty", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 0 {
					return nil, ArityMismatchError("0", len(values))
				}
				return d.IsEmpty(), nil
			}}}, nil

	case "containsKey":
		return &intrinsicFunc{d, "containsKey", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 1 {
					return nil, ArityMismatchError("1", len(values))
				}
				return d.ContainsKey(values[0])
			}}}, nil

	default:
		return nil, NoSuchFieldError(key.String())
	}
}
