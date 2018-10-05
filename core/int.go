// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

/*doc
## Int

Int is the set of all signed 64-bit integers (-9223372036854775808 to 9223372036854775807).

Integer literals can either be in decimal format, e.g. `123`, or hexidecimal format,
e.g. `0xabcd`.

Valid operators for an Int are:

* The equality operators `==`, `!=`
* The [`comparision`](interfaces.html#comparable) operators `>`, `>=`, `<`, `<=`, `<=>`
* The arithmetic operators `+`, `-`, `*`, `/`
* The integer arithmetic operators <code>&#124;</code>, `^`, `%`, `&`, `<<`, `>>`
* The unary integer complement operator `~`
* The postfix operators `++`, `--`

When applying an arithmetic operator `+`, `-`, `*`, `/` to an Int, if the other
operand is a Float, then the result will be a Float, otherwise the result will be an Int.

Ints are [`hashable`](interfaces.html#hashable)

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

func (i _int) ToInt() int64 {
	return int64(i)
}

func (i _int) ToFloat() float64 {
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
		b := t.ToFloat()
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
		b := t.ToFloat()
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
		b := t.ToFloat()
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
		b := t.ToFloat()
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
		b := t.ToFloat()
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
		b := t.ToFloat()
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

func (i _int) Abs() Int {
	n := i.ToInt()
	if n < 0 {
		return NewInt(-n)
	}
	return i
}

func (i _int) Format(base Int) (s Str, err Error) {

	n := i.ToInt()
	b := int(base.ToInt())

	// catch panic from strconv
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	s = MustStr(strconv.FormatInt(n, b))

	return
}

func (i _int) ToChar() (Str, Error) {
	r := rune(i.ToInt())

	if !utf8.ValidRune(r) {
		return nil, fmt.Errorf("InvalidUtf8Rune: %d is not a valid rune", r)
	}

	buf := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(buf, r)
	return NewStr(string(buf))
}

//--------------------------------------------------------------
// fields

/*doc
An Int has the following fields:

* [abs](#abs)
* [format](#format)
* [toChar](#tochar)
* [toFloat](#tofloat)

*/

var intMethods = map[string]Method{

	/*doc
	### `abs`

	`abs` returns the absolute value of the int.

	* signature: `abs() <Int>`
	* example: `let n = -1; println(n.abs())`

	*/
	"abs": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return self.(Int).Abs(), nil
		}),

	/*doc
	### `format`

	`format` returns the string representation of int in the given base,
	for 2 <= base <= 36. The result uses the lower-case letters 'a' to 'z'
	for digit values >= 10.  If the base is omitted, it defaults to 10.

	* signature: `format(base = 10 <Int>) <Str>`
	* example: `let n = 11259375; println(n.format(16))`

	*/
	"format": NewMultipleMethod(
		[]Type{},
		[]Type{IntType},
		false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			base := NewInt(10)
			if len(params) == 1 {
				base = params[0].(Int)
			}
			return self.(Int).Format(base)
		}),

	/*doc
	### `toChar`

	`toChar` converts an int that is a valid rune into a string with a single
	unicode character.

	* signature: `toChar() <Str>`
	* example: `let n = 19990; println(n.toChar())`

	*/
	"toChar": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return self.(Int).ToChar()
		}),

	/*doc
	### `toFloat`

	`toFloat` converts an int to a float

	* signature: `toFloat() <Float>`
	* example: `let n = 123; println(n.toFloat())`

	*/
	"toFloat": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return NewFloat(self.(Int).ToFloat()), nil
		}),
}

func (i _int) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(intMethods))
	for name := range intMethods {
		names = append(names, name)
	}
	return names, nil
}

func (i _int) HasField(name string) (bool, Error) {
	_, ok := intMethods[name]
	return ok, nil
}

func (i _int) GetField(ev Eval, name string) (Value, Error) {
	if method, ok := intMethods[name]; ok {
		return method.ToFunc(i, name), nil
	}
	return nil, NoSuchField(name)
}

func (i _int) InvokeField(ev Eval, name string, params []Value) (Value, Error) {
	if method, ok := intMethods[name]; ok {
		return method.Invoke(i, ev, params)
	}
	return nil, NoSuchField(name)
}
