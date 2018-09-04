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

func NewFieldStruct(fields map[string]Field, frozen bool) (Struct, Error) {

	for key, _ := range fields {
		if !scanner.IsIdentifier(key) {
			return nil, InvalidStructKeyError(key)
		}
	}

	return &_struct{
		fieldMap: &hashFieldMap{
			fields:     fields,
			replacable: true,
		},
		frozen: frozen,
	}, nil
}

func NewMethodStruct(self interface{}, methods map[string]Method) (Struct, Error) {

	for key, _ := range methods {
		if !scanner.IsIdentifier(key) {
			return nil, InvalidStructKeyError(key)
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
				return nil, InvalidArgumentError(
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
		i++
		buf.WriteString(" ")
		buf.WriteString(name)
		buf.WriteString(": ")

		v, err := st.GetField(name, ev)
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
	return NewStr(buf.String()), nil
}

func (st *_struct) HashCode(ev Eval) (Int, Error) {
	return nil, HashCodeMismatchError(StructType)
}

func (this *_struct) Eq(ev Eval, val Value) (Bool, Error) {

	// same type
	that, ok := val.(Struct)
	if !ok {
		return False, nil
	}

	// same number of fields
	n1, err := this.FieldNames()
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

		a, err := this.GetField(name, ev)
		if err != nil {
			return nil, err
		}

		b, err := that.GetField(name, ev)
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

func (st *_struct) GetField(name string, ev Eval) (Value, Error) {
	return st.fieldMap.get(name, ev)
}

func (st *_struct) InvokeField(name string, ev Eval, params []Value) (Value, Error) {
	return st.fieldMap.invoke(name, ev, params)
}

func (st *_struct) SetField(name string, ev Eval, val Value) Error {

	if st.frozen {
		return ImmutableValueError()
	}

	return st.fieldMap.set(name, ev, val)
}

func (st *_struct) Internal(args ...interface{}) {

	name := args[0].(string)
	field := args[1].(Field)
	st.fieldMap.replace(name, field)
}
