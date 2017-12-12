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

import (
	"fmt"
)

type channel struct {
	ch chan Value
}

func NewChan() Chan {
	return &channel{make(chan Value)}
}

func NewBufferedChan(size int) Chan {
	return &channel{make(chan Value, size)}
}

func (ch *channel) chanMarker() {}

func (ch *channel) TypeOf() Type { return TCHAN }

func (ch *channel) Eq(v Value) Bool {
	switch t := v.(type) {
	case *channel:
		// equality is based on identity
		return MakeBool(ch == t)
	default:
		return FALSE
	}
}

func (ch *channel) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (ch *channel) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (ch *channel) ToStr() Str {
	return MakeStr(fmt.Sprintf("channel<%p>", ch))
}

func (ch *channel) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(ch, t), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

//--------------------------------------------------------------
// intrinsic functions

func (ch *channel) GetField(key Str) (Value, Error) {
	switch key.String() {

	case "send":
		return &intrinsicFunc{ch, "send", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 1 {
					return nil, ArityMismatchError("1", len(values))
				}

				ch.ch <- values[0]
				return NULL, nil
			}}}, nil

	case "recv":
		return &intrinsicFunc{ch, "recv", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 0 {
					return nil, ArityMismatchError("0", len(values))
				}

				val := <-ch.ch
				return val, nil
			}}}, nil

	default:
		return nil, NoSuchFieldError(key.String())
	}
}
