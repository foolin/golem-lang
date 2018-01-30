// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	//	"strings"
)

type _int int64

var ZERO Int = NewInt(0)
var ONE Int = NewInt(1)
var NEG_ONE Int = NewInt(-1)

func (i _int) IntVal() int64 {
	return int64(i)
}

func (i _int) FloatVal() float64 {
	return float64(i)
}

func NewInt(i int64) Int {
	return _int(i)
}

//--------------------------------------------------------------
// Basic

func (i _int) basicMarker() {}

//--------------------------------------------------------------
// Value

func (i _int) Type() Type { return TINT }

func (i _int) Freeze() (Value, Error) {
	return i, nil
}

func (i _int) Frozen() (Bool, Error) {
	return TRUE, nil
}

func (i _int) ToStr(cx Context) Str {
	return NewStr(fmt.Sprintf("%d", i))
}

func (i _int) HashCode(cx Context) (Int, Error) {
	return i, nil
}

func (i _int) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {

	case _int:
		return NewBool(i == t), nil

	case _float:
		a := float64(i)
		b := t.FloatVal()
		return NewBool(a == b), nil

	default:
		return FALSE, nil
	}
}

func (i _int) Cmp(cx Context, v Value) (Int, Error) {
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

func (i _int) Add(v Value) (Number, Error) {
	switch t := v.(type) {

	case _int:
		return i + t, nil

	case _float:
		a := float64(i)
		b := t.FloatVal()
		return NewFloat(a + b), nil

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
		return NewFloat(a - b), nil

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
		return NewFloat(a * b), nil

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
			return NewFloat(a / b), nil
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

//--------------------------------------------------------------
// intrinsic functions

func (i _int) GetField(cx Context, key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}
