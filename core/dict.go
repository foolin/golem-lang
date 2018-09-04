// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"fmt"
)

type dict struct {
	hashMap *HashMap
	frozen  bool
}

// NewDict creates a new Dict
func NewDict(ev Eval, entries []*HEntry) (Dict, Error) {
	h, err := NewHashMap(ev, entries)
	if err != nil {
		return nil, err
	}
	return &dict{h, false}, nil
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
	return NewStr(buf.String()), nil
}

func (d *dict) HashCode(ev Eval) (Int, Error) {
	return nil, HashCodeMismatchError(DictType)
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
		return ImmutableValueError()
	}

	return d.hashMap.Put(ev, key, val)
}

func (d *dict) Len(ev Eval) (Int, Error) {
	return d.hashMap.Len(), nil
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
		return nil, ImmutableValueError()
	}

	d.hashMap = EmptyHashMap()
	return d, nil
}

func (d *dict) Remove(ev Eval, key Value) (Dict, Error) {
	if d.frozen {
		return nil, ImmutableValueError()
	}

	_, err := d.hashMap.Remove(ev, key)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *dict) AddAll(ev Eval, val Value) (Dict, Error) {
	if d.frozen {
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

		if tp, ok := v.(tuple); ok {
			if len(tp) == 2 {
				err = d.hashMap.Put(ev, tp[0], tp[1])
				if err != nil {
					return nil, err
				}
			} else {
				return nil, InvalidArgumentError(
					fmt.Sprintf("Expected Tuple of length %d, not length %d",
						2, len(tp)))
			}
		} else {
			return nil, TypeMismatchError(TupleType, v.Type())
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
	return nil, NoSuchElementError()
}

//--------------------------------------------------------------

var dictMethods = map[string]Method{

	"isEmpty": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			d := self.(Dict)
			return d.IsEmpty(), nil
		}),

	"contains": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			d := self.(Dict)
			return d.Contains(ev, params[0])
		}),

	"addAll": NewFixedMethod(
		[]Type{AnyType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			d := self.(Dict)
			return d.AddAll(ev, params[0])
		}),

	"clear": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			d := self.(Dict)
			return d.Clear()
		}),

	"remove": NewFixedMethod(
		[]Type{AnyType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			d := self.(Dict)
			return d.Remove(ev, params[0])
		}),
}

func (d *dict) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(dictMethods))
	for name, _ := range dictMethods {
		names = append(names, name)
	}
	return names, nil
}

func (d *dict) HasField(name string) (bool, Error) {
	_, ok := dictMethods[name]
	return ok, nil
}

func (d *dict) GetField(name string, ev Eval) (Value, Error) {
	if method, ok := dictMethods[name]; ok {
		return method.ToFunc(d, name), nil
	}
	return nil, NoSuchFieldError(name)
}

func (d *dict) InvokeField(name string, ev Eval, params []Value) (Value, Error) {

	if method, ok := dictMethods[name]; ok {
		return method.Invoke(d, ev, params)
	}
	return nil, NoSuchFieldError(name)
}
