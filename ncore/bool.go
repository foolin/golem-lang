// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

type _bool bool

// True is true
var True Bool = _bool(true)

// False is false
var False Bool = _bool(false)

// NewBool returns True or False
func NewBool(b bool) Bool {
	if b {
		return True
	}
	return False
}

func (b _bool) BoolVal() bool {
	return bool(b)
}

func (b _bool) basicMarker() {}

func (b _bool) Type() Type { return BoolType }

func (b _bool) Freeze(cx Context) (Value, Error) {
	return b, nil
}

func (b _bool) Frozen(cx Context) (Bool, Error) {
	return True, nil
}

//func (b _bool) ToStr(cx Context) Str {
//	if b {
//		return NewStr("true")
//	}
//	return NewStr("false")
//}
//
//func (b _bool) HashCode(cx Context) (Int, Error) {
//	if b {
//		return NewInt(1009), nil
//	}
//	return NewInt(1013), nil
//}

func (b _bool) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case _bool:
		if b == t {
			return _bool(true), nil
		}
		return _bool(false), nil
	default:
		return _bool(false), nil
	}
}

//func (b _bool) GetField(cx Context, key Str) (Value, Error) {
//	return nil, NoSuchFieldError(key.String())
//}
//
//func (b _bool) Cmp(cx Context, v Value) (Int, Error) {
//	switch t := v.(type) {
//
//	case _bool:
//		if b == t {
//			return Zero, nil
//		} else if b {
//			return One, nil
//		} else {
//			return NegOne, nil
//		}
//
//	default:
//		return nil, TypeMismatchError("Expected Comparable Type")
//	}
//}

func (b _bool) Not() Bool {
	return !b
}
