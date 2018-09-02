// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"reflect"
	"sort"
	"strings"

	"github.com/mjarmy/golem-lang/scanner"
)

//--------------------------------------------------------------
// Struct

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

func NewVirtualStruct(methods map[string]Method, frozen bool) (Struct, Error) {

	for key, _ := range methods {
		if !scanner.IsIdentifier(key) {
			return nil, InvalidStructKeyError(key)
		}
	}

	return &_struct{
		fieldMap: &virtualFieldMap{
			methods: methods,
		},
		frozen: frozen,
	}, nil
}

func (st *_struct) compositeMarker() {}

func (st *_struct) Type() Type { return StructType }

func (st *_struct) Freeze(ev Evaluator) (Value, Error) {
	st.frozen = true
	return st, nil
}

func (st *_struct) Frozen(ev Evaluator) (Bool, Error) {
	return NewBool(st.frozen), nil
}

func (st *_struct) ToStr(ev Evaluator) (Str, Error) {

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

func (st *_struct) HashCode(ev Evaluator) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (st *_struct) Eq(ev Evaluator, val Value) (Bool, Error) {

	// same type
	that, ok := val.(Struct)
	if !ok {
		return False, nil
	}

	// same fields
	n1, err := st.FieldNames()
	if err != nil {
		return nil, err
	}
	n2, err := that.FieldNames()
	if err != nil {
		return nil, err
	}
	// Unfortunately we have to sort the keys, because the underlying
	// golang map in the fieldMap iterates over its keys unpredictably.
	sort.Slice(n1, func(i, j int) bool {
		return strings.Compare(n1[i], n1[j]) < 0
	})
	sort.Slice(n2, func(i, j int) bool {
		return strings.Compare(n2[i], n2[j]) < 0
	})
	if !reflect.DeepEqual(n1, n2) {
		return False, nil
	}

	// all fields have same value
	for _, name := range n1 {

		a, err := st.GetField(name, ev)
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

func (st *_struct) Cmp(ev Evaluator, val Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

//--------------------------------------------------------------
// fields

func (st *_struct) FieldNames() ([]string, Error) {
	return st.fieldMap.names(), nil
}

func (st *_struct) HasField(name string) (bool, Error) {
	return st.fieldMap.has(name), nil
}

func (st *_struct) GetField(name string, ev Evaluator) (Value, Error) {
	return st.fieldMap.get(name, ev)
}

func (st *_struct) InvokeField(name string, ev Evaluator, params []Value) (Value, Error) {
	return st.fieldMap.invoke(name, ev, params)
}

func (st *_struct) SetField(name string, ev Evaluator, val Value) Error {

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
