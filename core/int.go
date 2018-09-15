// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

/*doc
## Int

Int is the set of all signed 64-bit integers (-9223372036854775808 to 9223372036854775807).

Valid operators for Int are:
	* The equality operators `==`, `!=`
	* The comparison operators `>`, `>=`, `<`, `<=`, `<=>`
	* The arithmetic operators `+`, `-`, `*`, `/`
	* The integer arithmetic operators <code>&#124;</code>, `^`, `%`, `&`, `<<`, `>>`
	* The unary integer complement operator `~`
	* The postfix operators `++`, `--`

When applying an arithmetic operator `+`, `-`, `*`, `/`to an Int, if the other
operand is a Float, then the result will be a Float, otherwise the result will be an Int.

Ints are [`hashable`](#TODO)

Int has no fields.

*/
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

func (i _int) Freeze(ev Eval) (Value, Error) {
	return i, nil
}

func (i _int) Frozen(ev Eval) (Bool, Error) {
	return True, nil
}

func (i _int) ToStr(ev Eval) (Str, Error) {
	return NewStr(fmt.Sprintf("%d", i))
}

func (i _int) HashCode(ev Eval) (Int, Error) {
	return i, nil
}

func (i _int) Eq(ev Eval, val Value) (Bool, Error) {
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

func (i _int) Cmp(ev Eval, c Comparable) (Int, Error) {
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
		return nil, ComparableMismatch(IntType, c.(Value).Type())
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
			return nil, DivideByZero()
		}
		return i / t, nil

	case _float:
		a := float64(i)
		b := t.FloatVal()
		if b == 0.0 {
			return nil, DivideByZero()
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
			return nil, InvalidArgument("Shift count cannot be less than zero")
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
			return nil, InvalidArgument("Shift count cannot be less than zero")
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

func (i _int) GetField(ev Eval, name string) (Value, Error) {
	return nil, NoSuchField(name)
}

func (i _int) InvokeField(ev Eval, name string, params []Value) (Value, Error) {
	return nil, NoSuchField(name)
}
