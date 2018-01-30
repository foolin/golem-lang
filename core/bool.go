// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

type _bool bool

var TRUE Bool = _bool(true)
var FALSE Bool = _bool(false)

func NewBool(b bool) Bool {
	if b {
		return TRUE
	} else {
		return FALSE
	}
}

func (b _bool) BoolVal() bool {
	return bool(b)
}

func (b _bool) basicMarker() {}

func (b _bool) Type() Type { return TBOOL }

func (b _bool) Freeze() (Value, Error) {
	return b, nil
}

func (b _bool) Frozen() (Bool, Error) {
	return TRUE, nil
}

func (b _bool) ToStr(cx Context) Str {
	if b {
		return NewStr("true")
	} else {
		return NewStr("false")
	}
}

func (b _bool) HashCode(cx Context) (Int, Error) {
	if b {
		return NewInt(1009), nil
	} else {
		return NewInt(1013), nil
	}
}

func (b _bool) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case _bool:
		if b == t {
			return _bool(true), nil
		} else {
			return _bool(false), nil
		}
	default:
		return _bool(false), nil
	}
}

func (b _bool) GetField(cx Context, key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (b _bool) Cmp(cx Context, v Value) (Int, Error) {
	switch t := v.(type) {

	case _bool:
		if b == t {
			return ZERO, nil
		} else if b {
			return ONE, nil
		} else {
			return NEG_ONE, nil
		}

	default:
		return nil, TypeMismatchError("Expected Comparable Type")
	}
}

func (b _bool) Not() Bool {
	return !b
}
