// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

/*doc
# Null

Null represents the absence of a value. The only instance of Null is `null`.

Null has no valid operators, and no fields.

*/
type null struct {
	// https://golang.org/ref/spec#Size_and_alignment_guarantees
	// Zero-size variables share the same address, so we use a placeholder
	// to give Null a size.
	placeholder int
}

// Null represents the null value
var Null NullValue = &null{1019}

func (n *null) basicMarker() {}

func (n *null) Type() Type {
	return NullType
}

func (n *null) ToStr(ev Eval) (Str, Error) {
	return NewStr("null")
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

func (n *null) GetField(ev Eval, name string) (Value, Error) {
	return nil, NullValueError()
}

func (n *null) InvokeField(ev Eval, name string, params []Value) (Value, Error) {
	return nil, NullValueError()
}
