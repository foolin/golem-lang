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

func (b _bool) Type() Type { return TBOOL }

func (b _bool) Freeze() (Value, Error) {
	return b, nil
}

func (b _bool) Frozen() (Bool, Error) {
	return TRUE, nil
}

func (b _bool) ToStr(cx Context) Str {
	if b {
		return MakeStr("true")
	} else {
		return MakeStr("false")
	}
}

func (b _bool) HashCode(cx Context) (Int, Error) {
	if b {
		return MakeInt(1009), nil
	} else {
		return MakeInt(1013), nil
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
