// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
)

//---------------------------------------------------------------
// tuple

type tuple []Value

// NewTuple creates a new Tuple
func NewTuple(values []Value) Tuple {
	if len(values) < 2 {
		panic("invalid tuple size")
	}
	return tuple(values)
}

func (tp tuple) compositeMarker() {}

func (tp tuple) Type() Type { return TupleType }

func (tp tuple) Freeze() (Value, Error) {
	return tp, nil
}

func (tp tuple) Frozen() (Bool, Error) {
	return True, nil
}

func (tp tuple) ToStr(cx Context) Str {
	var buf bytes.Buffer
	buf.WriteString("(")
	for idx, v := range tp {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.ToStr(cx).String())
	}
	buf.WriteString(")")
	return NewStr(buf.String())
}

func (tp tuple) HashCode(cx Context) (Int, Error) {

	// https://en.wikipedia.org/wiki/Jenkins_hash_function
	var hash int64
	for _, v := range tp {
		h, err := v.HashCode(cx)
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

func (tp tuple) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case tuple:
		return valuesEq(cx, tp, t)
	default:
		return False, nil
	}
}

func (tp tuple) GetField(cx Context, key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (tp tuple) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (tp tuple) Get(cx Context, index Value) (Value, Error) {
	idx, err := boundedIndex(index, len(tp))
	if err != nil {
		return nil, err
	}
	return tp[idx], nil
}

func (tp tuple) Len() Int {
	return NewInt(int64(len(tp)))
}
