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

// NOTE: 'null' cannot be an empty struct, because empty structs have
// unusual semantics in Go, insofar as they all point to the same address.
//
// https://golang.org/ref/spec#Size_and_alignment_guarantees
//
// To work around that, we place an arbitrary value inside the struct, so
// that it wont be empty.  This gives the singleton instance of null
// its own address
//
type null struct {
	placeholder int
}

var NULL Null = &null{0}

func (n *null) basicMarker() {}

func (n *null) TypeOf() Type { return TNULL }

func (n *null) ToStr() Str { return MakeStr("null") }

func (n *null) HashCode() (Int, Error) { return nil, NullValueError() }

func (n *null) GetField(key Str) (Value, Error) { return nil, NullValueError() }

func (n *null) Eq(v Value) Bool {
	switch v.(type) {
	case *null:
		return TRUE
	default:
		return FALSE
	}
}

func (n *null) Cmp(v Value) (Int, Error) { return nil, NullValueError() }

func (n *null) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(n, t), nil

	default:
		return nil, NullValueError()
	}
}
