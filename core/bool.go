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

type _bool bool

var TRUE Bool = _bool(true)
var FALSE Bool = _bool(false)

func MakeBool(b bool) Bool {
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

func (b _bool) TypeOf() Type { return TBOOL }

func (b _bool) ToStr() Str {
	if b {
		return MakeStr("true")
	} else {
		return MakeStr("false")
	}
}

func (b _bool) HashCode() (Int, Error) {
	if b {
		return MakeInt(1009), nil
	} else {
		return MakeInt(1013), nil
	}
}

func (b _bool) Eq(v Value) Bool {
	switch t := v.(type) {
	case _bool:
		if b == t {
			return _bool(true)
		} else {
			return _bool(false)
		}
	default:
		return _bool(false)
	}
}

func (b _bool) GetField(key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (b _bool) Cmp(v Value) (Int, Error) {
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

func (b _bool) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(b, t), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (b _bool) Not() Bool {
	return !b
}
