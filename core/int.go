// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

type _int int64

// Zero is the integer 0
var Zero = NewInt(0)

// One is the integer 1
var One = NewInt(1)

// NegOne is the integer -1
var NegOne = NewInt(-1)

// NewInt creates a new Int
func NewInt(i int64) Int {
	return _int(i)
}

func (i _int) IntVal() int64 {
	return int64(i)
}

func (i _int) FloatVal() float64 {
	return float64(i)
}

//--------------------------------------------------------------
// Basic

func (i _int) basicMarker() {}

//--------------------------------------------------------------
// Value

func (i _int) Type() Type { return IntType }

func (i _int) Freeze(ev Evaluator) (Value, Error) {
	return i, nil
}

func (i _int) Frozen(ev Evaluator) (Bool, Error) {
	return True, nil
}

func (i _int) ToStr(ev Evaluator) (Str, Error) {
	return NewStr(fmt.Sprintf("%d", i)), nil
}

func (i _int) HashCode(ev Evaluator) (Int, Error) {
	return i, nil
}

func (i _int) Eq(ev Evaluator, val Value) (Bool, Error) {
	switch t := val.(type) {

	case _int:
		return NewBool(i == t), nil

	case _float:
		a := float64(i)
		b := t.FloatVal()
		return NewBool(a == b), nil

	default:
		return False, nil
	}
}

func (i _int) Cmp(ev Evaluator, c Comparable) (Int, Error) {
	switch t := c.(type) {

	case _int:
		if i < t {
			return NegOne, nil
		} else if i > t {
			return One, nil
		} else {
			return Zero, nil
		}

	case _float:
		a := float64(i)
		b := t.FloatVal()
		if a < b {
			return NegOne, nil
		} else if a > b {
			return One, nil
		} else {
			return Zero, nil
		}

	default:
		return nil, ComparableMismatchError(IntType, c.(Value).Type())
	}
}

//--------------------------------------------------------------
// Number

func (i _int) Add(n Number) Number {
	switch t := n.(type) {

	case _int:
		return i + t

	case _float:
		a := float64(i)
		b := t.FloatVal()
		return NewFloat(a + b)

	default:
		panic("unreachable")
	}
}

func (i _int) Sub(n Number) Number {
	switch t := n.(type) {

	case _int:
		return i - t

	case _float:
		a := float64(i)
		b := t.FloatVal()
		return NewFloat(a - b)

	default:
		panic("unreachable")
	}
}

func (i _int) Mul(n Number) Number {
	switch t := n.(type) {

	case _int:
		return i * t

	case _float:
		a := float64(i)
		b := t.FloatVal()
		return NewFloat(a * b)

	default:
		panic("unreachable")
	}
}

func (i _int) Div(n Number) (Number, Error) {
	switch t := n.(type) {

	case _int:
		if t == 0 {
			return nil, DivideByZeroError()
		}
		return i / t, nil

	case _float:
		a := float64(i)
		b := t.FloatVal()
		if b == 0.0 {
			return nil, DivideByZeroError()
		}
		return NewFloat(a / b), nil

	default:
		panic("unreachable")
	}
}

func (i _int) Negate() Number {
	return 0 - i
}

//--------------------------------------------------------------
// Int

func (i _int) Rem(n Int) Int {
	switch t := n.(type) {
	case _int:
		return i % t
	default:
		panic("unreachable")
	}
}

func (i _int) BitAnd(n Int) Int {
	switch t := n.(type) {
	case _int:
		return i & t
	default:
		panic("unreachable")
	}
}

func (i _int) BitOr(n Int) Int {
	switch t := n.(type) {
	case _int:
		return i | t
	default:
		panic("unreachable")
	}
}

func (i _int) BitXOr(n Int) Int {
	switch t := n.(type) {
	case _int:
		return i ^ t
	default:
		panic("unreachable")
	}
}

func (i _int) LeftShift(n Int) (Int, Error) {
	switch t := n.(type) {
	case _int:
		if t < 0 {
			return nil, InvalidArgumentError("Shift count cannot be less than zero")
		}
		return i << uint(t), nil
	default:
		panic("unreachable")
	}
}

func (i _int) RightShift(n Int) (Int, Error) {
	switch t := n.(type) {
	case _int:
		if t < 0 {
			return nil, InvalidArgumentError("Shift count cannot be less than zero")
		}
		return i >> uint(t), nil
	default:
		panic("unreachable")
	}
}

func (i _int) Complement() Int {
	return ^i
}

//--------------------------------------------------------------
// fields

func (i _int) FieldNames() ([]string, Error) {
	return []string{}, nil
}

func (i _int) HasField(name string) (bool, Error) {
	return false, nil
}

func (i _int) GetField(name string, ev Evaluator) (Value, Error) {
	return nil, NoSuchFieldError(name)
}

func (i _int) InvokeField(name string, ev Evaluator, params []Value) (Value, Error) {
	return nil, NoSuchFieldError(name)
}
