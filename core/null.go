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

// NOTE: the type 'null' cannot be an empty struct, because empty structs have
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

func (n *null) Type() Type { return TNULL }

func (n *null) ToStr(cx Context) Str { return MakeStr("null") }

func (n *null) HashCode(cx Context) (Int, Error) { return nil, NullValueError() }

func (n *null) GetField(cx Context, key Str) (Value, Error) { return nil, NullValueError() }

func (n *null) Eq(cx Context, v Value) (Bool, Error) {
	switch v.(type) {
	case *null:
		return TRUE, nil
	default:
		return FALSE, nil
	}
}

func (n *null) Cmp(cx Context, v Value) (Int, Error) { return nil, NullValueError() }

func (n *null) Freeze() (Value, Error) {
	return nil, NullValueError()
}

func (n *null) Frozen() (Bool, Error) {
	return nil, NullValueError()
}
