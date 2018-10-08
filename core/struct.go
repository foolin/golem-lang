// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"fmt"
	//"sync"

	"github.com/mjarmy/golem-lang/scanner"
)

/*doc
## Struct

A Struct is a collection of fields that are defined in golem code.

Valid operators for Struct are:

* The equality operators `==`, `!=`

Structs do not have any pre-defined fields.

Structs can have the following magic fields:

* `__eq__` overrides the `==` operator

	* signature: `__eq__(x <Value>) <Bool>`

* `__hashCode__` causes a struct to be [hashable](interfaces.html#hashable), so it can be
used as a key in a dict, or an entry in a set.  Note that if you define `__hashCode__`,
you *must* also always define `__eq__`.  Values that are equal must have the same hashCode.

	* signature: `__hashCode__() <Int>`

* `__str__` overrides the value returned by the builtin function [`str`](builtins.html#str)

	* signature: `__str__() <Str>`

*/

//--------------------------------------------------------------
// Field Struct
//--------------------------------------------------------------

type _struct struct {
	fieldMap fieldMap
	frozen   bool
}

// NewStruct create a new Struct backed by Fields
func NewStruct(fields map[string]Field) (Struct, Error) {

	for key := range fields {
		if !scanner.IsIdentifier(key) {
			return nil, InvalidStructKey(key)
		}
	}

	return &_struct{
		fieldMap: &hashFieldMap{
			fields:     fields,
			replacable: true,
		},
		frozen: false,
	}, nil
}

// NewFrozenStruct create a new frozen Struct backed by Fields
func NewFrozenStruct(fields map[string]Field) (Struct, Error) {

	for key := range fields {
		if !scanner.IsIdentifier(key) {
			return nil, InvalidStructKey(key)
		}
	}

	return &_struct{
		fieldMap: &hashFieldMap{
			fields:     fields,
			replacable: true,
		},
		frozen: true,
	}, nil
}

// NewMethodStruct create a new Struct backed by Methods.
func NewMethodStruct(self interface{}, methods map[string]Method) (Struct, Error) {

	for key := range methods {
		if !scanner.IsIdentifier(key) {
			return nil, InvalidStructKey(key)
		}
	}

	return &_struct{
		fieldMap: &methodFieldMap{
			self:    self,
			methods: methods,
			//funcs:   map[string]NativeFunc{},
			//mx:      sync.Mutex{},
		},
		frozen: true,
	}, nil
}

// MergeStructs creates a new Struct by merging together some existing structs.
func MergeStructs(structs []Struct) (Struct, Error) {

	if len(structs) < 2 {
		panic(fmt.Errorf("invalid struct merge size: %d", len(structs)))
	}

	frozen := structs[0].(*_struct).frozen
	fieldMaps := make([]fieldMap, len(structs))
	for i, st := range structs {
		s := st.(*_struct)

		if i > 0 {
			if frozen != s.frozen {
				return nil, InvalidArgument(
					"Cannot merge structs unless they are all frozen, or all unfrozen")
			}
		}

		fieldMaps[i] = s.fieldMap
	}

	return &_struct{
		fieldMap: mergeFieldMaps(fieldMaps),
		frozen:   frozen,
	}, nil

}

func (st *_struct) compositeMarker() {}

func (st *_struct) Type() Type { return StructType }

func (st *_struct) Freeze(ev Eval) (Value, Error) {
	st.frozen = true
	return st, nil
}

func (st *_struct) Frozen(ev Eval) (Bool, Error) {
	return NewBool(st.frozen), nil
}

func (st *_struct) ToStr(ev Eval) (Str, Error) {

	magic, err := st.HasField("__str__")
	if err != nil {
		return nil, err
	}
	if magic {
		return st.magicStr(ev)
	}

	//---------------------------------------

	var buf bytes.Buffer
	buf.WriteString("struct {")

	names := st.fieldMap.names()
	for i, name := range names {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(" ")
		buf.WriteString(name)
		buf.WriteString(": ")

		v, err := st.GetField(ev, name)
		if err != nil {
			return nil, err
		}

		s, err := v.ToStr(ev)
		if err != nil {
			return nil, err
		}

		buf.WriteString(s.String())
	}

	buf.WriteString(" }")
	return NewStr(buf.String())
}

func (st *_struct) magicStr(ev Eval) (Str, Error) {

	fv, err := st.GetField(ev, "__str__")
	if err != nil {
		return nil, err
	}

	fn, ok := fv.(Func)
	if !ok {
		return nil, fmt.Errorf(
			"TypeMismatch: __str__ must be a Func, not %s", fv.Type())
	}

	// check arity
	expected := Arity{FixedArity, 0, 0}
	if fn.Arity() != expected {
		return nil, fmt.Errorf(
			"ArityMismatch: __str__ function must have 0 parameters")
	}

	result, err := fn.Invoke(ev, nil)
	if err != nil {
		return nil, err
	}

	s, ok := result.(Str)
	if !ok {
		return nil, fmt.Errorf(
			"TypeMismatch: __str__ must return a Str, not %s", result.Type())
	}

	return s, nil
}

func (st *_struct) HashCode(ev Eval) (Int, Error) {

	magic, err := st.HasField("__hashCode__")
	if err != nil {
		return nil, err
	}
	if magic {
		return st.magicHashCode(ev)
	}

	//---------------------------------------

	return nil, HashCodeMismatch(StructType)
}

func (st *_struct) magicHashCode(ev Eval) (Int, Error) {

	fv, err := st.GetField(ev, "__hashCode__")
	if err != nil {
		return nil, err
	}

	fn, ok := fv.(Func)
	if !ok {
		return nil, fmt.Errorf(
			"TypeMismatch: __hashCode__ must be a Func, not %s", fv.Type())
	}

	// check arity
	expected := Arity{FixedArity, 0, 0}
	if fn.Arity() != expected {
		return nil, fmt.Errorf(
			"ArityMismatch: __hashCode__ function must have 0 parameters")
	}

	result, err := fn.Invoke(ev, nil)
	if err != nil {
		return nil, err
	}

	i, ok := result.(Int)
	if !ok {
		return nil, fmt.Errorf(
			"TypeMismatch: __hashCode__ must return an Int, not %s", result.Type())
	}

	return i, nil
}

func (st *_struct) Eq(ev Eval, val Value) (Bool, Error) {

	magic, err := st.HasField("__eq__")
	if err != nil {
		return nil, err
	}
	if magic {
		return st.magicEq(ev, val)
	}

	//---------------------------------------

	// same type
	that, ok := val.(Struct)
	if !ok {
		return False, nil
	}

	// same number of fields
	n1, err := st.FieldNames()
	if err != nil {
		return nil, err
	}
	n2, err := that.FieldNames()
	if err != nil {
		return nil, err
	}
	if len(n1) != len(n2) {
		return False, nil
	}

	// same fields, with same values
	for _, name := range n1 {

		has, err := that.HasField(name)
		if err != nil {
			return nil, err
		}
		if !has {
			return False, nil
		}

		a, err := st.GetField(ev, name)
		if err != nil {
			return nil, err
		}

		b, err := that.GetField(ev, name)
		if err != nil {
			return nil, err
		}

		eq, err := a.Eq(ev, b)
		if err != nil {
			return nil, err
		}
		if !eq.BoolVal() {
			return False, nil
		}
	}

	// done
	return True, nil
}

func (st *_struct) magicEq(ev Eval, val Value) (Bool, Error) {

	fv, err := st.GetField(ev, "__eq__")
	if err != nil {
		return nil, err
	}

	fn, ok := fv.(Func)
	if !ok {
		return nil, fmt.Errorf(
			"TypeMismatch: __eq__ must be a Func, not %s", fv.Type())
	}

	// check arity
	expected := Arity{FixedArity, 1, 0}
	if fn.Arity() != expected {
		return nil, fmt.Errorf(
			"ArityMismatch: __eq__ function must have 1 parameter")
	}

	result, err := fn.Invoke(ev, []Value{val})
	if err != nil {
		return nil, err
	}

	b, ok := result.(Bool)
	if !ok {
		return nil, fmt.Errorf(
			"TypeMismatch: __eq__ must return a Bool, not %s", result.Type())
	}

	return b, nil
}

func (st *_struct) ToDict(ev Eval) (Dict, Error) {

	names := st.fieldMap.names()
	entries := make([]*HEntry, len(names))
	for i, n := range names {
		key := MustStr(n)
		val, err := st.GetField(ev, n)
		if err != nil {
			return nil, err
		}
		entries[i] = &HEntry{key, val}
	}

	hm, err := NewHashMap(ev, entries)
	if err != nil {
		return nil, err
	}
	return NewDict(hm), nil
}

func (st *_struct) Internal(args ...interface{}) {

	name := args[0].(string)
	field := args[1].(Field)
	st.fieldMap.replace(name, field)
}

//--------------------------------------------------------------
// fields

func (st *_struct) FieldNames() ([]string, Error) {
	return st.fieldMap.names(), nil
}

func (st *_struct) HasField(name string) (bool, Error) {
	return st.fieldMap.has(name), nil
}

func (st *_struct) GetField(ev Eval, name string) (Value, Error) {
	return st.fieldMap.get(ev, name)
}

func (st *_struct) InvokeField(ev Eval, name string, params []Value) (Value, Error) {
	return st.fieldMap.invoke(ev, name, params)
}

func (st *_struct) SetField(ev Eval, name string, val Value) Error {

	if st.frozen {
		return ImmutableValue()
	}

	return st.fieldMap.set(ev, name, val)
}
