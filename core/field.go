// Copyright 2017 The Golem Project Developers
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

//--------------------------------------------------------------
// A Field is a name-value pair inside a Struct

type Field interface {
	Name() string
}

type field struct {
	name       string
	isConst    bool
	isProperty bool
	value      Value
}

func (f *field) Name() string {
	return f.name
}

// Create a name-value pair.
func NewField(name string, isConst bool, value Value) Field {
	return &field{name, isConst, false, value}
}

// Create a Property using 'getter' and 'setter' functions
// If the setter is nil, then the property is const.
func NewProperty(name string, getter Func, setter Func) Field {

	if getter.MinArity() != 0 || getter.MaxArity() != 0 {
		panic("Property getter does not have arity 0")
	}

	if setter != nil {
		if setter.MinArity() != 1 || setter.MaxArity() != 1 {
			panic("Property setter does not have arity 1")
		}
	}

	prop := NewTuple([]Value{getter, setter})
	return &field{name, setter == nil, true, prop}
}

//fnGetter := NewNativeFunc(0, 0,
//	func(cx Context, values []Value) (Value, Error) {
//		return getter()
//	})

//var fnSetter NativeFunc = nil
//if setter != nil {
//	fnSetter = NewNativeFunc(1, 1,
//		func(cx Context, values []Value) (Value, Error) {
//			return nil, setter(values[0])
//		})
//}
