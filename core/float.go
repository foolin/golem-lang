// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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

// NewFloat creates a new Float
func NewFloat(f float64) Float {
	return _float(f)
}

func (f _float) basicMarker() {}

func (f _float) Type() Type { return FloatType }

func (f _float) Freeze(ev Evaluator) (Value, Error) {
	return f, nil
}

func (f _float) Frozen(ev Evaluator) (Bool, Error) {
	return True, nil
}

func (f _float) ToStr(ev Evaluator) (Str, Error) {
	return NewStr(fmt.Sprintf("%g", f)), nil
}

func (f _float) HashCode(ev Evaluator) (Int, Error) {

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

	return NewInt(hashCode), nil
}

func (f _float) Eq(ev Evaluator, val Value) (Bool, Error) {
	switch t := val.(type) {

	case _float:
		return NewBool(f == t), nil

	case _int:
		return NewBool(f.FloatVal() == t.FloatVal()), nil

	default:
		return False, nil
	}
}

func (f _float) Cmp(ev Evaluator, c Comparable) (Int, Error) {
	switch t := c.(type) {

	case _float:
		if f < t {
			return NegOne, nil
		} else if f > t {
			return One, nil
		} else {
			return Zero, nil
		}

	case _int:
		g := _float(t)
		if f < g {
			return NegOne, nil
		} else if f > g {
			return One, nil
		} else {
			return Zero, nil
		}

	default:
		return nil, TempMismatchError("Expected Comparable")
	}
}

func (f _float) Add(val Value) (Number, Error) {
	switch t := val.(type) {

	case _int:
		return f + _float(t), nil

	case _float:
		return f + t, nil

	default:
		return nil, NumberMismatchError(val.Type())
	}
}

func (f _float) Sub(val Value) (Number, Error) {
	switch t := val.(type) {

	case _int:
		return f - _float(t), nil

	case _float:
		return f - t, nil

	default:
		return nil, NumberMismatchError(val.Type())
	}
}

func (f _float) Mul(val Value) (Number, Error) {
	switch t := val.(type) {

	case _int:
		return f * _float(t), nil

	case _float:
		return f * t, nil

	default:
		return nil, NumberMismatchError(val.Type())
	}
}

func (f _float) Div(val Value) (Number, Error) {
	switch t := val.(type) {

	case _int:
		if t == 0 {
			return nil, DivideByZeroError()
		}
		return f / _float(t), nil

	case _float:
		if t == 0.0 {
			return nil, DivideByZeroError()
		}
		return f / t, nil

	default:
		return nil, NumberMismatchError(val.Type())
	}
}

func (f _float) Negate() Number {
	return 0 - f
}

//--------------------------------------------------------------
// fields

func (f _float) FieldNames() ([]string, Error) {
	return []string{}, nil
}

func (f _float) HasField(name string) (bool, Error) {
	return false, nil
}

func (f _float) GetField(name string, ev Evaluator) (Value, Error) {
	return nil, NoSuchFieldError(name)
}

func (f _float) InvokeField(name string, ev Evaluator, params []Value) (Value, Error) {
	return nil, NoSuchFieldError(name)
}
