// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"fmt"
)

//---------------------------------------------------------------
// tuple

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
	return NewStr(buf.String()), nil
}

func (tp tuple) HashCode(ev Eval) (Int, Error) {

	// https://en.wikipedia.org/wiki/Jenkins_hash_function
	var hash int64
	for _, v := range tp {
		h, err := v.HashCode(ev)
		if err != nil {
			return nil, err
		}
		hash += h.IntVal()
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

//--------------------------------------------------------------
// fields

func (tp tuple) FieldNames() ([]string, Error) {
	return []string{}, nil
}

func (tp tuple) HasField(name string) (bool, Error) {
	return false, nil
}

func (tp tuple) GetField(ev Eval, name string) (Value, Error) {
	return nil, NoSuchField(name)
}

func (tp tuple) InvokeField(ev Eval, name string, params []Value) (Value, Error) {
	return nil, NoSuchField(name)
}
