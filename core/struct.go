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

func MergeStructs(structs []Struct) Struct {
	if len(structs) < 2 {
		panic("invalid struct merge")
	}

	smap := newStructMap()
	frozen := false

	// FieldNames that are defined in more
	// than one of the structs are combined so that the value
	// is taken only from the first such struct

	// If any of the structs are frozen,
	// then the resulting struct is also frozen.

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

func (st *_struct) Type() Type { return TSTRUCT }

func (st *_struct) Freeze() (Value, Error) {
	st.frozen = true
	return st, nil
}

func (st *_struct) Frozen() (Bool, Error) {
	return MakeBool(st.frozen), nil
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
	return MakeStr(buf.String())
}

func (st *_struct) HashCode(cx Context) (Int, Error) {
	// TODO $hash()
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (st *_struct) Eq(cx Context, v Value) (Bool, Error) {

	// same type
	that, ok := v.(Struct)
	if !ok {
		return FALSE, nil
	}

	// same number of fields
	fields := st.FieldNames()
	if len(fields) != len(that.FieldNames()) {
		return FALSE, nil
	}

	// all fields have same value
	for _, n := range fields {
		a, err := st.GetField(cx, str(n))
		assert(err == nil)

		b, err := that.GetField(cx, str(n))
		if err != nil {
			assert(err.Kind() == NO_SUCH_FIELD)
		}

		eq, err := a.Eq(cx, b)
		if err != nil {
			return nil, err
		}
		if eq != TRUE {
			return FALSE, nil
		}
	}

	// done
	return TRUE, nil
}

func (st *_struct) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (st *_struct) GetField(cx Context, name Str) (Value, Error) {
	f, has := st.smap.get(name.String())
	if has {
		if f.isProperty {
			fn := ((f.value.(tuple))[0]).(Func)
			return fn.Invoke(cx, nil)
		} else {
			return f.value, nil
		}
	} else {
		return nil, NoSuchFieldError(name.String())
	}
}

func (st *_struct) FieldNames() []string {
	return st.smap.fieldNames()
}

func (st *_struct) Has(name Value) (Bool, Error) {
	if s, ok := name.(Str); ok {
		_, has := st.smap.get(s.String())
		return MakeBool(has), nil
	} else {
		return nil, TypeMismatchError("Expected 'Str'")
	}
}

//---------------------------------------------------------------
// Mutation

func (st *_struct) InitField(cx Context, name Str, val Value) Error {
	// We ignore 'frozen' here, since we are initializing the value

	f, has := st.smap.get(name.String())
	if has {
		// We ignore IsConst here, since we are initializing the value
		if f.isProperty {
			// TODO what are the semantics here?
			// Do we really want to invoke this func?

			//// The value for a property is always a tuple
			//// containing two functions: the getter, and the setter.
			//fn := ((f.value.(tuple))[1]).(Func)
			//_, err := fn.Invoke(cx, []Value{val})
			//return err

			return nil
		} else {
			f.value = val
			return nil
		}
	} else {
		return NoSuchFieldError(name.String())
	}
}

func (st *_struct) SetField(cx Context, name Str, val Value) Error {

	if st.frozen {
		return ImmutableValueError()
	}

	f, has := st.smap.get(name.String())
	if has {
		if f.isConst {
			return ReadonlyFieldError(name.String())
		} else {
			if f.isProperty {
				// The value for a property is always a tuple
				// containing two functions: the getter, and the setter.
				fn := ((f.value.(tuple))[1]).(Func)
				_, err := fn.Invoke(cx, []Value{val})
				return err
			} else {
				f.value = val
				return nil
			}
		}
	} else {
		return NoSuchFieldError(name.String())
	}
}
