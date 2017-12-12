// Copyright 2017 The Golem Project Developers //
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"bytes"
	//"fmt"
)

type StructEntry struct {
	Key        string
	IsConst    bool
	IsProperty bool
	Value      Value
}

type StructEntryDef struct {
	Key        string
	IsConst    bool
	IsProperty bool
}

//--------------------------------------------------------------

func NewStruct(entries []*StructEntry) (Struct, Error) {

	smap := newStructMap()
	for _, e := range entries {
		if _, has := smap.get(e.Key); has {
			return nil, DuplicateFieldError(e.Key)
		}
		smap.put(e)
	}

	return &_struct{smap}, nil
}

func BlankStruct(def []*StructEntryDef) (Struct, Error) {

	smap := newStructMap()
	for _, d := range def {
		if _, has := smap.get(d.Key); has {
			return nil, DuplicateFieldError(d.Key)
		}
		smap.put(&StructEntry{d.Key, d.IsConst, d.IsProperty, NULL})
	}

	return &_struct{smap}, nil
}

func MergeStructs(structs []Struct) Struct {
	if len(structs) < 2 {
		panic("invalid struct merge")
	}

	smap := newStructMap()

	// subtlety: keys that are defined in more
	// than one of the structs are combined so that the value
	// is taken only from the first such struct
	for _, s := range structs {
		for _, b := range (s.(*_struct)).smap.buckets {
			for _, e := range b {
				smap.put(e)
			}
		}
	}

	return &_struct{smap}
}

type _struct struct {
	smap *structMap
}

func (stc *_struct) compositeMarker() {}

func (stc *_struct) TypeOf() Type { return TSTRUCT }

func (stc *_struct) ToStr() Str {

	var buf bytes.Buffer
	buf.WriteString("struct {")
	for i, k := range stc.Keys() {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(" ")
		buf.WriteString(k)
		buf.WriteString(": ")

		v, err := stc.GetField(str(k))
		Assert(err == nil, "invalid struct")
		buf.WriteString(v.ToStr().String())
	}
	buf.WriteString(" }")
	return MakeStr(buf.String())
}

func (stc *_struct) HashCode() (Int, Error) {
	// TODO $hash()
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (stc *_struct) Eq(v Value) Bool {

	// same type
	that, ok := v.(Struct)
	if !ok {
		return FALSE
	}

	// same number of keys
	keys := stc.Keys()
	if len(keys) != len(that.Keys()) {
		return FALSE
	}

	// all keys have same value
	for _, k := range keys {
		a, err := stc.GetField(str(k))
		Assert(err == nil, "invalid chain")

		b, err := that.GetField(str(k))
		if err != nil {
			Assert(err.Kind() == NO_SUCH_FIELD, "invalid chain")
			return FALSE
		}

		if a.Eq(b) != TRUE {
			return FALSE
		}
	}

	// done
	return TRUE
}

func (stc *_struct) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (stc *_struct) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(stc, t), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (stc *_struct) Get(index Value) (Value, Error) {
	if s, ok := index.(Str); ok {
		return stc.GetField(s)
	} else {
		return nil, TypeMismatchError("Expected 'Str'")
	}
}

func (stc *_struct) Set(index Value, val Value) Error {
	if s, ok := index.(Str); ok {
		return stc.SetField(s, val)
	} else {
		return TypeMismatchError("Expected 'Str'")
	}
}

func (stc *_struct) GetField(key Str) (Value, Error) {
	e, has := stc.smap.get(key.String())
	if has {
		if e.IsProperty {
			// The value for a property is always a tuple
			// containing two functions: the getter, and the setter.
			// TODO Add support for BytecodeFunc properties.
			fn := ((e.Value.(tuple))[0]).(NativeFunc)
			return fn.Invoke(nil)
		} else {
			return e.Value, nil
		}
	} else {
		return nil, NoSuchFieldError(key.String())
	}
}

func (stc *_struct) SetField(key Str, val Value) Error {
	e, has := stc.smap.get(key.String())
	if has {
		if e.IsConst {
			return ReadonlyFieldError(key.String())
		} else {

			if e.IsProperty {
				// The value for a property is always a tuple
				// containing two functions: the getter, and the setter.
				// TODO Add support for BytecodeFunc properties.
				fn := ((e.Value.(tuple))[1]).(NativeFunc)
				_, err := fn.Invoke([]Value{val})
				return err
			} else {
				e.Value = val
				return nil
			}
		}
	} else {
		return NoSuchFieldError(key.String())
	}
}

func (stc *_struct) Keys() []string {
	return stc.smap.keys()
}

func (stc *_struct) Has(key Value) (Bool, Error) {
	if s, ok := key.(Str); ok {
		_, has := stc.smap.get(s.String())
		return MakeBool(has), nil
	} else {
		return nil, TypeMismatchError("Expected 'Str'")
	}
}

func (stc *_struct) InitField(key Str, val Value) Error {
	e, has := stc.smap.get(key.String())
	if has {
		// We ignore IsConst here, since we are initializing the value
		e.Value = val
		return nil
	} else {
		return NoSuchFieldError(key.String())
	}
}
