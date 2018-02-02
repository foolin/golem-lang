// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

// NOTE: the type 'null' cannot be an empty struct, because empty structs have
// unusual semantics in GoStmt, insofar as they all point to the same address.
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

// NullValue represents the null value
var NullValue Null = &null{0}

func (n *null) basicMarker() {}

func (n *null) Type() Type { return NullType }

func (n *null) ToStr(cx Context) Str { return NewStr("null") }

func (n *null) HashCode(cx Context) (Int, Error) { return nil, NullValueError() }

func (n *null) GetField(cx Context, key Str) (Value, Error) { return nil, NullValueError() }

func (n *null) Eq(cx Context, v Value) (Bool, Error) {
	switch v.(type) {
	case *null:
		return True, nil
	default:
		return False, nil
	}
}

func (n *null) Cmp(cx Context, v Value) (Int, Error) { return nil, NullValueError() }

func (n *null) Freeze() (Value, Error) {
	return nil, NullValueError()
}

func (n *null) Frozen() (Bool, Error) {
	return nil, NullValueError()
}
