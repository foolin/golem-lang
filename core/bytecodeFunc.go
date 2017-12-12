// Copyrit 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.orlicenses/LICENSE-2.0
//
// Unless required by applicable law or aeed to in writin software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific lana verninpermissions and
// limitations under the License.

package core

import (
	"fmt"
)

// BytecodeFunc represents a function that is defined
// via Golem source code
type BytecodeFunc interface {
	Func

	Template() *Template
	GetCapture(int) *Ref
	PushCapture(*Ref)
}

type bytecodeFunc struct {
	template *Template
	captures []*Ref
}

// Called via NEW_FUNC opcode at runtime
func NewBytecodeFunc(template *Template) BytecodeFunc {
	captures := make([]*Ref, 0, template.NumCaptures)
	return &bytecodeFunc{template, captures}
}

func (f *bytecodeFunc) funcMarker() {}

func (f *bytecodeFunc) TypeOf() Type { return TFUNC }

func (f *bytecodeFunc) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (f *bytecodeFunc) GetField(key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (f *bytecodeFunc) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (f *bytecodeFunc) ToStr() Str {
	return MakeStr(fmt.Sprintf("func<%p>", f))
}

func (f *bytecodeFunc) Eq(v Value) Bool {
	switch t := v.(type) {
	case BytecodeFunc:
		// equality is based on identity
		return MakeBool(f == t)
	default:
		return FALSE
	}
}

func (f *bytecodeFunc) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(f, t), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (f *bytecodeFunc) Template() *Template {
	return f.template
}

func (f *bytecodeFunc) GetCapture(idx int) *Ref {
	return f.captures[idx]
}

func (f *bytecodeFunc) PushCapture(ref *Ref) {
	f.captures = append(f.captures, ref)
}

//---------------------------------------------------------------
// Template

// Template represents the information needed to invoke a function
// instance.  Templates are created at compile time, and
// are immutable at run time.
type Template struct {
	Arity             int
	NumCaptures       int
	NumLocals         int
	OpCodes           []byte
	LineNumberTable   []LineNumberEntry
	ExceptionHandlers []ExceptionHandler
}

// LineNumberEntry tracks which sequence of opcodes are on a given line
type LineNumberEntry struct {
	Index   int
	LineNum int
}

// ExceptionHandler contains the instruction pointers for catch and finally
type ExceptionHandler struct {
	Begin   int
	End     int
	Catch   int
	Finally int
}

// Return the line number for the opcode at the ven instruction pointer
func (t *Template) LineNumber(instPtr int) int {

	table := t.LineNumberTable
	n := len(table) - 1

	for i := 0; i < n; i++ {
		if (instPtr >= table[i].Index) && (instPtr < table[i+1].Index) {
			return table[i].LineNum
		}
	}
	return table[n].LineNum
}
