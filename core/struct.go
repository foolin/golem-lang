// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"fmt"

	"github.com/mjarmy/golem-lang/scanner"
)

//--------------------------------------------------------------
// Field Struct
//--------------------------------------------------------------

type _struct struct {
	fieldMap fieldMap
	frozen   bool
}

// NewFieldStruct create a new Struct backed by Fields
func NewFieldStruct(fields map[string]Field) (Struct, Error) {

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

// NewFrozenFieldStruct create a new frozen Struct backed by Fields
func NewFrozenFieldStruct(fields map[string]Field) (Struct, Error) {

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

func (st *_struct) HashCode(ev Eval) (Int, Error) {
	return nil, HashCodeMismatch(StructType)
}

func (st *_struct) Eq(ev Eval, val Value) (Bool, Error) {

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

func (st *_struct) Internal(args ...interface{}) {

	name := args[0].(string)
	field := args[1].(Field)
	st.fieldMap.replace(name, field)
}
