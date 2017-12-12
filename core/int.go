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
	"fmt"
	//	"strings"
)

type _int int64

var ZERO Int = MakeInt(0)
var ONE Int = MakeInt(1)
var NEG_ONE Int = MakeInt(-1)

func (i _int) IntVal() int64 {
	return int64(i)
}

func (i _int) FloatVal() float64 {
	return float64(i)
}

func MakeInt(i int64) Int {
	return _int(i)
}

//--------------------------------------------------------------
// Basic

func (i _int) basicMarker() {}

//--------------------------------------------------------------
// Value

func (i _int) TypeOf() Type { return TINT }

func (i _int) ToStr() Str {
	return MakeStr(fmt.Sprintf("%d", i))
}

func (i _int) HashCode() (Int, Error) {
	return i, nil
}

func (i _int) Eq(v Value) Bool {
	switch t := v.(type) {

	case _int:
		return MakeBool(i == t)

	case _float:
		a := float64(i)
		b := t.FloatVal()
		return MakeBool(a == b)

	default:
		return FALSE
	}
}

func (i _int) GetField(key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (i _int) Cmp(v Value) (Int, Error) {
	switch t := v.(type) {

	case _int:
		if i < t {
			return NEG_ONE, nil
		} else if i > t {
			return ONE, nil
		} else {
			return ZERO, nil
		}

	case _float:
		a := float64(i)
		b := t.FloatVal()
		if a < b {
			return NEG_ONE, nil
		} else if a > b {
			return ONE, nil
		} else {
			return ZERO, nil
		}

	default:
		return nil, TypeMismatchError("Expected Comparable Type")
	}
}

func (i _int) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(i, t), nil

	case _int:
		return i + t, nil

	case _float:
		a := float64(i)
		b := t.FloatVal()
		return MakeFloat(a + b), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

//--------------------------------------------------------------
// Number

func (i _int) Sub(v Value) (Number, Error) {
	switch t := v.(type) {

	case _int:
		return i - t, nil

	case _float:
		a := float64(i)
		b := t.FloatVal()
		return MakeFloat(a - b), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (i _int) Mul(v Value) (Number, Error) {
	switch t := v.(type) {

	case _int:
		return i * t, nil

	case _float:
		a := float64(i)
		b := t.FloatVal()
		return MakeFloat(a * b), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (i _int) Div(v Value) (Number, Error) {
	switch t := v.(type) {

	case _int:
		if t == 0 {
			return nil, DivideByZeroError()
		} else {
			return i / t, nil
		}

	case _float:
		a := float64(i)
		b := t.FloatVal()
		if b == 0.0 {
			return nil, DivideByZeroError()
		} else {
			return MakeFloat(a / b), nil
		}

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (i _int) Negate() Number {
	return 0 - i
}

//--------------------------------------------------------------
// Int

func (i _int) Rem(v Value) (Int, Error) {
	switch t := v.(type) {
	case _int:
		return i % t, nil
	default:
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

func (i _int) BitAnd(v Value) (Int, Error) {
	switch t := v.(type) {
	case _int:
		return i & t, nil
	default:
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

func (i _int) BitOr(v Value) (Int, Error) {
	switch t := v.(type) {
	case _int:
		return i | t, nil
	default:
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

func (i _int) BitXOr(v Value) (Int, Error) {
	switch t := v.(type) {
	case _int:
		return i ^ t, nil
	default:
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

func (i _int) LeftShift(v Value) (Int, Error) {
	switch t := v.(type) {
	case _int:
		if t < 0 {
			return nil, InvalidArgumentError("Shift count cannot be less than zero")
		} else {
			return i << uint(t), nil
		}
	default:
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

func (i _int) RightShift(v Value) (Int, Error) {
	switch t := v.(type) {
	case _int:
		if t < 0 {
			return nil, InvalidArgumentError("Shift count cannot be less than zero")
		} else {
			return i >> uint(t), nil
		}
	default:
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

func (i _int) Complement() Int {
	return ^i
}
