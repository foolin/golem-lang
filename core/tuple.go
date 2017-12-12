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
)

//---------------------------------------------------------------
// tuple

type tuple []Value

func NewTuple(values []Value) Tuple {
	if len(values) < 2 {
		panic("invalid tuple size")
	}
	return tuple(values)
}

func (tp tuple) compositeMarker() {}

func (tp tuple) TypeOf() Type { return TTUPLE }

func (tp tuple) ToStr() Str {
	var buf bytes.Buffer
	buf.WriteString("(")
	for idx, v := range tp {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.ToStr().String())
	}
	buf.WriteString(")")
	return MakeStr(buf.String())
}

func (tp tuple) HashCode() (Int, Error) {

	// https://en.wikipedia.org/wiki/Jenkins_hash_function
	var hash int64 = 0
	for _, v := range tp {
		h, err := v.HashCode()
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
	return MakeInt(hash), nil
}

func (tp tuple) Eq(v Value) Bool {
	switch t := v.(type) {
	case tuple:
		return valuesEq(tp, t)
	default:
		return FALSE
	}
}

func (tp tuple) GetField(key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (tp tuple) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (tp tuple) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(tp, t), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (tp tuple) Get(index Value) (Value, Error) {
	idx, err := validateIndex(index, len(tp))
	if err != nil {
		return nil, err
	}
	return tp[idx.IntVal()], nil
}

func (tp tuple) Len() Int {
	return MakeInt(int64(len(tp)))
}
