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

// Null represents the null value
var Null Nil = &null{0}

func (n *null) basicMarker() {}

func (n *null) Type() Type {
	return NullType
}

func (n *null) ToStr(ev Eval) (Str, Error) {
	return NewStr("null"), nil
}

func (n *null) HashCode(ev Eval) (Int, Error) {
	return nil, NullValueError()
}

func (n *null) Eq(ev Eval, val Value) (Bool, Error) {
	switch val.(type) {
	case *null:
		return True, nil
	default:
		return False, nil
	}
}

func (n *null) Freeze(ev Eval) (Value, Error) {
	return nil, NullValueError()
}

func (n *null) Frozen(ev Eval) (Bool, Error) {
	return nil, NullValueError()
}

//--------------------------------------------------------------
// fields

func (n *null) FieldNames() ([]string, Error) {
	return nil, NullValueError()
}

func (n *null) HasField(name string) (bool, Error) {
	return false, NullValueError()
}

func (n *null) GetField(name string, ev Eval) (Value, Error) {
	return nil, NullValueError()
}

func (n *null) InvokeField(name string, ev Eval, params []Value) (Value, Error) {
	return nil, NullValueError()
}
