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

func (f _float) Freeze(ev Eval) (Value, Error) {
	return f, nil
}

func (f _float) Frozen(ev Eval) (Bool, Error) {
	return True, nil
}

func (f _float) ToStr(ev Eval) (Str, Error) {
	return NewStr(fmt.Sprintf("%g", f)), nil
}

func (f _float) HashCode(ev Eval) (Int, Error) {

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

func (f _float) Eq(ev Eval, val Value) (Bool, Error) {
	if n, ok := val.(Number); ok {
		fv := f.FloatVal()
		nv := n.FloatVal()
		return NewBool(fv == nv), nil
	}
	return False, nil
}

func (f _float) Cmp(ev Eval, c Comparable) (Int, Error) {
	if n, ok := c.(Number); ok {
		fv := f.FloatVal()
		nv := n.FloatVal()
		if fv < nv {
			return NegOne, nil
		} else if fv > nv {
			return One, nil
		} else {
			return Zero, nil
		}
	}
	return nil, ComparableMismatchError(FloatType, c.(Value).Type())
}

//--------------------------------------------------------------
// Number

func (f _float) Add(n Number) Number {
	fv := f.FloatVal()
	nv := n.FloatVal()
	return NewFloat(fv + nv)
}

func (f _float) Sub(n Number) Number {
	fv := f.FloatVal()
	nv := n.FloatVal()
	return NewFloat(fv - nv)
}

func (f _float) Mul(n Number) Number {
	fv := f.FloatVal()
	nv := n.FloatVal()
	return NewFloat(fv * nv)
}

func (f _float) Div(n Number) (Number, Error) {
	if n.FloatVal() == 0.0 {
		return nil, DivideByZeroError()
	}
	fv := f.FloatVal()
	nv := n.FloatVal()
	return NewFloat(fv / nv), nil
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

func (f _float) GetField(name string, ev Eval) (Value, Error) {
	return nil, NoSuchFieldError(name)
}

func (f _float) InvokeField(name string, ev Eval, params []Value) (Value, Error) {
	return nil, NoSuchFieldError(name)
}
