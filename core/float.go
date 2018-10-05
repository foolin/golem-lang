// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

/*doc
## Float

Float is the set of all IEEE-754 64-bit floating-point numbers.

Valid operators for Float are:

* The equality operators `==`, `!=`
* The [`comparision`](interfaces.html#comparable) operators `>`, `>=`, `<`, `<=`, `<=>`
* The arithmetic operators `+`, `-`, `*`, `/`
* The postfix operators `++`, `--`

Applying an arithmetic operator to a Float always returns a Float.

Floats are [`hashable`](interfaces.html#hashable)

*/
type _float float64

func (f _float) ToInt() int64 {
	return int64(f)
}

func (f _float) ToFloat() float64 {
	return float64(f)
}

// NewFloat creates a new Float
func NewFloat(f float64) Float {
	return _float(f)
}

func (f _float) basicMarker() {}

func (f _float) Type() Type { return FloatType }

func (f _float) Freeze(ev Eval) (Value, Error) {
	return f, nil
}

func (f _float) Frozen(ev Eval) (Bool, Error) {
	return True, nil
}

func (f _float) ToStr(ev Eval) (Str, Error) {
	return NewStr(fmt.Sprintf("%g", f))
}

func (f _float) HashCode(ev Eval) (Int, Error) {

	writer := new(bytes.Buffer)
	err := binary.Write(writer, binary.LittleEndian, f.ToFloat())
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

	return NewInt(hashCode), nil
}

func (f _float) Eq(ev Eval, val Value) (Bool, Error) {
	if n, ok := val.(Number); ok {
		fv := f.ToFloat()
		nv := n.ToFloat()
		return NewBool(fv == nv), nil
	}
	return False, nil
}

func (f _float) Cmp(ev Eval, c Comparable) (Int, Error) {
	if n, ok := c.(Number); ok {
		fv := f.ToFloat()
		nv := n.ToFloat()
		if fv < nv {
			return NegOne, nil
		} else if fv > nv {
			return One, nil
		} else {
			return Zero, nil
		}
	}
	return nil, ComparableMismatch(FloatType, c.(Value).Type())
}

//--------------------------------------------------------------
// Number

func (f _float) Add(n Number) Number {
	fv := f.ToFloat()
	nv := n.ToFloat()
	return NewFloat(fv + nv)
}

func (f _float) Sub(n Number) Number {
	fv := f.ToFloat()
	nv := n.ToFloat()
	return NewFloat(fv - nv)
}

func (f _float) Mul(n Number) Number {
	fv := f.ToFloat()
	nv := n.ToFloat()
	return NewFloat(fv * nv)
}

func (f _float) Div(n Number) (Number, Error) {
	if n.ToFloat() == 0.0 {
		return nil, DivideByZero()
	}
	fv := f.ToFloat()
	nv := n.ToFloat()
	return NewFloat(fv / nv), nil
}

func (f _float) Negate() Number {
	return 0 - f
}

//--------------------------------------------------------------
// Float

func (f _float) Abs() Float {
	n := f.ToFloat()
	if n < 0 {
		return NewFloat(-n)
	}
	return f
}

func (f _float) Ceil() Float {
	return _float(math.Ceil(f.ToFloat()))
}

func (f _float) Floor() Float {
	return _float(math.Floor(f.ToFloat()))
}

func (f _float) Round() Float {
	return _float(math.Round(f.ToFloat()))
}

func (f _float) Format(fstr Str, prec Int) (s Str, err Error) {

	n := f.ToFloat()
	fm, e := parseFloatFormat(fstr.String())
	if e != nil {
		return nil, e
	}
	p := int(prec.ToInt())

	// catch panic from strconv
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	s = MustStr(strconv.FormatFloat(n, fm, p, 64))

	return
}

func parseFloatFormat(fstr string) (byte, error) {

	switch fstr {

	case "b", "e", "E", "f", "g", "G":
		return []byte(fstr)[0], nil

	default:
		return 0, fmt.Errorf("'%s' is an invalid format string", fstr)
	}
}

//--------------------------------------------------------------
// fields

/*doc
A Float has the following fields:

* [abs](#abs)
* [ceil](#ceil)
* [floor](#floor)
* [format](#format)
* [round](#round)

*/

var floatMethods = map[string]Method{

	/*doc
	### `abs`

	`abs` returns the absolute value of the float.

	* signature: `abs() <Float>`
	* example: `let n = -1.2; println(n.abs())`

	*/
	"abs": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return self.(Float).Abs(), nil
		}),

	/*doc
	### `ceil`

	`ceil` returns the least integer value greater than or equal to the float.

	* signature: `ceil() <Float>`
	* example: `let n = -1.2; println(n.ceil())`

	*/
	"ceil": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return self.(Float).Ceil(), nil
		}),

	/*doc
	### `floor`

	`floor` returns the greatest integer value less than or equal to the float.

	* signature: `floor() <Float>`
	* example: `let n = -1.2; println(n.floor())`

	*/
	"floor": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return self.(Float).Floor(), nil
		}),

	/*doc
	### `format`

	`format`

	* signature: `format(fmt <Str>, prec = -1 <Int>) <Str>`
	* example: `let n = 1.23; println(n.format("f"))`

	*/
	"format": NewMultipleMethod(
		[]Type{StrType},
		[]Type{IntType},
		false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			fstr := params[0].(Str)
			prec := NewInt(-1)
			if len(params) == 2 {
				prec = params[1].(Int)
			}
			return self.(Float).Format(fstr, prec)
		}),

	/*doc
	### `round`

	`round` returns the nearest integer, rounding half away from zero.

	* signature: `round() <Float>`
	* example: `let n = -1.2; println(n.round())`

	*/
	"round": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return self.(Float).Round(), nil
		}),
}

func (f _float) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(floatMethods))
	for name := range floatMethods {
		names = append(names, name)
	}
	return names, nil
}

func (f _float) HasField(name string) (bool, Error) {
	_, ok := floatMethods[name]
	return ok, nil
}

func (f _float) GetField(ev Eval, name string) (Value, Error) {
	if method, ok := floatMethods[name]; ok {
		return method.ToFunc(f, name), nil
	}
	return nil, NoSuchField(name)
}

func (f _float) InvokeField(ev Eval, name string, params []Value) (Value, Error) {
	if method, ok := floatMethods[name]; ok {
		return method.Invoke(f, ev, params)
	}
	return nil, NoSuchField(name)
}
