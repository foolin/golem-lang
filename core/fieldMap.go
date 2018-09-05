// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

type (
	fieldMap interface {
		names() []string
		has(string) bool
		get(string, Eval) (Value, Error)
		invoke(string, Eval, []Value) (Value, Error)
		set(string, Eval, Value) Error

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
		case *methodFieldMap:
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

// TODO: Using a golang map is a placeholder implementation.
// Since we know that no keys will be added or removed once the fieldMap
// is instantiated, there are probably more efficient implementations
// available that we can substitute in at some point.

type hashFieldMap struct {
	fields     map[string]Field
	replacable bool
}

func (fm *hashFieldMap) names() []string {

	names := make([]string, 0, len(fm.fields))
	for name := range fm.fields {
		names = append(names, name)
	}
	return names
}

func (fm *hashFieldMap) has(name string) bool {

	_, ok := fm.fields[name]
	return ok
}

func (fm *hashFieldMap) get(name string, ev Eval) (Value, Error) {

	if f, ok := fm.fields[name]; ok {
		return f.Get(ev)
	}
	return nil, NoSuchFieldError(name)
}

func (fm *hashFieldMap) invoke(name string, ev Eval, params []Value) (Value, Error) {

	if f, ok := fm.fields[name]; ok {
		return f.Invoke(ev, params)
	}
	return nil, NoSuchFieldError(name)
}

func (fm *hashFieldMap) set(name string, ev Eval, val Value) Error {
	if f, ok := fm.fields[name]; ok {
		if f.IsReadonly() {
			return ReadonlyFieldError(name)
		}
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
// methodFieldMap
//--------------------------------------------------------------

// TODO: Using a golang map is a placeholder implementation.
// Since we know that no keys will be added or removed once the fieldMap
// is instantiated, there are probably more efficient implementations
// available that we can substitute in at some point.

type methodFieldMap struct {
	self    interface{}
	methods map[string]Method
}

func (fm *methodFieldMap) names() []string {

	names := make([]string, 0, len(fm.methods))
	for name := range fm.methods {
		names = append(names, name)
	}
	return names
}

func (fm *methodFieldMap) has(name string) bool {

	_, ok := fm.methods[name]
	return ok
}

func (fm *methodFieldMap) get(name string, ev Eval) (Value, Error) {

	if m, ok := fm.methods[name]; ok {
		return m.ToFunc(fm.self, name), nil
	}
	return nil, NoSuchFieldError(name)
}

func (fm *methodFieldMap) invoke(name string, ev Eval, params []Value) (Value, Error) {

	if m, ok := fm.methods[name]; ok {
		return m.Invoke(fm.self, ev, params)
	}
	return nil, NoSuchFieldError(name)
}

func (fm *methodFieldMap) set(name string, ev Eval, val Value) Error {

	if _, ok := fm.methods[name]; ok {
		return ReadonlyFieldError(name)
	}
	return NoSuchFieldError(name)
}

func (fm *methodFieldMap) replace(name string, field Field) {
	panic("Internal Error")
}
