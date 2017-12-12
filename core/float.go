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
	"encoding/binary"
	"fmt"
)

type _float float64

func (f _float) IntVal() int64 {
	return int64(f)
}

func (f _float) FloatVal() float64 {
	return float64(f)
}

func MakeFloat(f float64) Float {
	return _float(f)
}

func (f _float) basicMarker() {}

func (f _float) TypeOf() Type { return TFLOAT }

func (f _float) ToStr() Str {
	return MakeStr(fmt.Sprintf("%g", f))
}

func (f _float) HashCode() (Int, Error) {

	writer := new(bytes.Buffer)
	err := binary.Write(writer, binary.LittleEndian, f.FloatVal())
	if err != nil {
		panic("Float.HashCode() write failed")
	}
	b := writer.Bytes()

	var hashCode int64
	reader := bytes.NewReader(b)
	err = binary.Read(reader, binary.LittleEndian, &hashCode)
	if err != nil {
		panic("Float.HashCode() read failed")
	}

	return MakeInt(hashCode), nil
}

func (f _float) Eq(v Value) Bool {
	switch t := v.(type) {

	case _float:
		return MakeBool(f == t)

	case _int:
		return MakeBool(f.FloatVal() == t.FloatVal())

	default:
		return FALSE
	}
}

func (f _float) GetField(key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (f _float) Cmp(v Value) (Int, Error) {
	switch t := v.(type) {

	case _float:
		if f < t {
			return NEG_ONE, nil
		} else if f > t {
			return ONE, nil
		} else {
			return ZERO, nil
		}

	case _int:
		g := _float(t)
		if f < g {
			return NEG_ONE, nil
		} else if f > g {
			return ONE, nil
		} else {
			return ZERO, nil
		}

	default:
		return nil, TypeMismatchError("Expected Comparable Type")
	}
}

func (f _float) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(f, t), nil

	case _int:
		return f + _float(t), nil

	case _float:
		return f + t, nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (f _float) Sub(v Value) (Number, Error) {
	switch t := v.(type) {

	case _int:
		return f - _float(t), nil

	case _float:
		return f - t, nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (f _float) Mul(v Value) (Number, Error) {
	switch t := v.(type) {

	case _int:
		return f * _float(t), nil

	case _float:
		return f * t, nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (f _float) Div(v Value) (Number, Error) {
	switch t := v.(type) {

	case _int:
		if t == 0 {
			return nil, DivideByZeroError()
		} else {
			return f / _float(t), nil
		}

	case _float:
		if t == 0.0 {
			return nil, DivideByZeroError()
		} else {
			return f / t, nil
		}

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (f _float) Negate() Number {
	return 0 - f
}
