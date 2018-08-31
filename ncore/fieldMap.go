// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
//"fmt"
)

type (
	fieldMap interface {
		names() []string
		has(string) bool
		get(string, Evaluator) (Value, Error)
		invoke(string, Evaluator, []Value) (Value, Error)
		set(string, Evaluator, Value) Error

		replace(string, Field)
	}
)

func mergeFieldMaps(fieldMaps []fieldMap) fieldMap {

	fields := make(map[string]Field)

	for _, fm := range fieldMaps {

		switch t := fm.(type) {
		case *hashFieldMap:
			for k, v := range t.fields {
				fields[k] = v
			}
		case *virtualFieldMap:
			for k, v := range t.methods {
				fn := v.ToFunc(t.self, k)
				fields[k] = NewReadonlyField(fn)
			}
		default:
			panic("unreachable")
		}

	}

	return &hashFieldMap{fields, false}
}

//--------------------------------------------------------------
// hashFieldMap
//--------------------------------------------------------------

type hashFieldMap struct {
	fields     map[string]Field
	replacable bool
}

func (fm *hashFieldMap) names() []string {

	names := make([]string, 0, len(fm.fields))
	for name, _ := range fm.fields {
		names = append(names, name)
	}
	return names
}

func (fm *hashFieldMap) has(name string) bool {

	_, ok := fm.fields[name]
	return ok
}

func (fm *hashFieldMap) get(name string, ev Evaluator) (Value, Error) {

	if f, ok := fm.fields[name]; ok {
		return f.Get(ev)
	}
	return nil, NoSuchFieldError(name)
}

func (fm *hashFieldMap) invoke(name string, ev Evaluator, params []Value) (Value, Error) {

	if f, ok := fm.fields[name]; ok {
		return f.Invoke(ev, params)
	}
	return nil, NoSuchFieldError(name)
}

func (fm *hashFieldMap) set(name string, ev Evaluator, val Value) Error {
	if f, ok := fm.fields[name]; ok {
		return f.Set(ev, val)
	}
	return NoSuchFieldError(name)
}

func (fm *hashFieldMap) replace(name string, field Field) {

	if !fm.replacable {
		panic("Internal Error")
	}

	_, ok := fm.fields[name]
	if !ok {
		panic("Internal Error")
	}

	fm.fields[name] = field
}

//--------------------------------------------------------------
// virtualFieldMap
//--------------------------------------------------------------

type virtualFieldMap struct {
	self    interface{}
	methods map[string]Method
}

func (fm *virtualFieldMap) names() []string {

	names := make([]string, 0, len(fm.methods))
	for name, _ := range fm.methods {
		names = append(names, name)
	}
	return names
}

func (fm *virtualFieldMap) has(name string) bool {

	_, ok := fm.methods[name]
	return ok
}

func (fm *virtualFieldMap) get(name string, ev Evaluator) (Value, Error) {

	if m, ok := fm.methods[name]; ok {
		return m.ToFunc(fm.self, name), nil
	}
	return nil, NoSuchFieldError(name)
}

func (fm *virtualFieldMap) invoke(name string, ev Evaluator, params []Value) (Value, Error) {

	if m, ok := fm.methods[name]; ok {
		return m.Invoke(fm.self, ev, params)
	}
	return nil, NoSuchFieldError(name)
}

func (fm *virtualFieldMap) set(name string, ev Evaluator, val Value) Error {

	if _, ok := fm.methods[name]; ok {
		return ReadonlyFieldError(name)
	}
	return NoSuchFieldError(name)
}

func (fm *virtualFieldMap) replace(name string, field Field) {
	panic("Internal Error")
}
