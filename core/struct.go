// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	//"fmt"
)

//--------------------------------------------------------------
// Struct

type _struct struct {
	smap   *structMap
	frozen bool
}

// NewStruct creates a new Struct
func NewStruct(fields []Field, frozen bool) (Struct, Error) {

	smap := newStructMap()
	for _, f := range fields {
		ff := f.(*field)
		if _, has := smap.get(ff.name); has {
			return nil, DuplicateFieldError(ff.name)
		}
		smap.put(ff)
	}

	return &_struct{smap, frozen}, nil
}

// DefineStruct defines an un-initialized Struct.
// This function is called by the Golem Interpreter -- you shouldn't use it yourself
// unless you are completely sure you know what you are doing.
func DefineStruct(defs []*FieldDef) (Struct, Error) {

	smap := newStructMap()
	for _, d := range defs {
		if _, has := smap.get(d.Name); has {
			return nil, DuplicateFieldError(d.Name)
		}
		smap.put(&field{d.Name, d.IsReadonly, d.IsProperty, Null})
	}

	return &_struct{smap, false}, nil
}

// MergeStructs merges Structs together into one Struct.
// Field name that are defined in more than one of the structs are combined so
// that the value of the field is taken only from the first such Struct.
// IfStmt any of the structs are frozen, then the resulting struct is also frozen.
func MergeStructs(structs []Struct) Struct {
	if len(structs) < 2 {
		panic("invalid struct merge")
	}

	smap := newStructMap()
	frozen := false

	for _, s := range structs {
		st := s.(*_struct)
		if st.frozen {
			frozen = true
		}
		for _, b := range st.smap.buckets {
			for _, f := range b {
				smap.put(f)
			}
		}
	}

	return &_struct{smap, frozen}
}

func (st *_struct) compositeMarker() {}

func (st *_struct) Type() Type { return StructType }

func (st *_struct) Freeze() (Value, Error) {
	st.frozen = true
	return st, nil
}

func (st *_struct) Frozen() (Bool, Error) {
	return NewBool(st.frozen), nil
}

func (st *_struct) ToStr(cx Context) Str {

	var buf bytes.Buffer
	buf.WriteString("struct {")
	for i, n := range st.FieldNames() {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(" ")
		buf.WriteString(n)
		buf.WriteString(": ")

		v, err := st.GetField(cx, str(n))
		assert(err == nil)
		buf.WriteString(v.ToStr(cx).String())
	}
	buf.WriteString(" }")
	return NewStr(buf.String())
}

func (st *_struct) HashCode(cx Context) (Int, Error) {
	// TODO $hash()
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (st *_struct) Eq(cx Context, v Value) (Bool, Error) {

	// same type
	that, ok := v.(Struct)
	if !ok {
		return False, nil
	}

	// same number of fields
	fields := st.FieldNames()
	if len(fields) != len(that.FieldNames()) {
		return False, nil
	}

	// all fields have same value
	for _, n := range fields {
		a, err := st.GetField(cx, str(n))
		assert(err == nil)

		b, err := that.GetField(cx, str(n))
		if err != nil {
			return False, nil
		}

		eq, err := a.Eq(cx, b)
		if err != nil {
			return nil, err
		}
		if eq != True {
			return False, nil
		}
	}

	// done
	return True, nil
}

func (st *_struct) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (st *_struct) GetField(cx Context, name Str) (Value, Error) {
	f, has := st.smap.get(name.String())
	if has {
		if f.isProperty {
			// The value for a property is always a tuple
			// containing two functions: the getter, and the setter.
			fn := ((f.value.(tuple))[0]).(Func)
			return fn.Invoke(cx, nil)
		}
		return f.value, nil
	}
	return nil, NoSuchFieldError(name.String())
}

func (st *_struct) FieldNames() []string {
	return st.smap.fieldNames()
}

func (st *_struct) Has(name Value) (Bool, Error) {
	if s, ok := name.(Str); ok {
		_, has := st.smap.get(s.String())
		return NewBool(has), nil
	}
	return nil, TypeMismatchError("Expected Str")
}

//---------------------------------------------------------------
// Mutation

// InitField initializes the Value of a Field in a Struct.
// This function is called by the Golem Interpreter -- you shouldn't use it yourself
// unless you are completely sure you know what you are doing.
func (st *_struct) InitField(cx Context, name Str, val Value) Error {

	// We ignore 'frozen' and isReadonly here, since we are initializing the value

	f, has := st.smap.get(name.String())
	if !has {
		return NoSuchFieldError(name.String())
	}
	f.value = val
	return nil
}

// SetField sets the Value of a Field in a Struct.
func (st *_struct) SetField(cx Context, name Str, val Value) Error {

	if st.frozen {
		return ImmutableValueError()
	}

	f, has := st.smap.get(name.String())
	if has {
		switch {
		case f.isReadonly:
			return ReadonlyFieldError(name.String())
		case f.isProperty:
			// The value for a property is always a tuple
			// containing two functions: the getter, and the setter.
			fn := ((f.value.(tuple))[1]).(Func)
			_, err := fn.Invoke(cx, []Value{val})
			return err
		default:

			f.value = val
			return nil
		}
	} else {
		return NoSuchFieldError(name.String())
	}
}
