// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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

// NewBytecodeFunc creates a new BytecodeFunc.  NewBytecodeFunc is
// called via NewFunc opcode at runtime.
func NewBytecodeFunc(template *Template) BytecodeFunc {
	captures := make([]*Ref, 0, template.NumCaptures)
	return &bytecodeFunc{template, captures}
}

func (f *bytecodeFunc) funcMarker() {}

func (f *bytecodeFunc) Type() Type { return FuncType }

func (f *bytecodeFunc) Freeze() (Value, Error) {
	return f, nil
}

func (f *bytecodeFunc) Frozen() (Bool, Error) {
	return True, nil
}

func (f *bytecodeFunc) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (f *bytecodeFunc) GetField(cx Context, key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (f *bytecodeFunc) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (f *bytecodeFunc) ToStr(cx Context) Str {
	return NewStr(fmt.Sprintf("func<%p>", f))
}

func (f *bytecodeFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case BytecodeFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *bytecodeFunc) Template() *Template {
	return f.template
}

func (f *bytecodeFunc) MinArity() int { return f.template.Arity }
func (f *bytecodeFunc) MaxArity() int { return f.template.Arity }

func (f *bytecodeFunc) GetCapture(idx int) *Ref {
	return f.captures[idx]
}

func (f *bytecodeFunc) PushCapture(ref *Ref) {
	f.captures = append(f.captures, ref)
}

func (f *bytecodeFunc) Invoke(cx Context, values []Value) (Value, Error) {
	return cx.Eval(f, values)
}

//---------------------------------------------------------------
// Template

// Template represents the information needed to invoke a function
// instance.  Templates are created at compile time, and
// are immutable at run time.
type Template struct {
	Arity             int // TODO MinArity, MaxArity
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

// ExceptionHandler handles exceptions for a given block of opcodes,
// by providing the instruction pointers for 'catch' and 'finally'
type ExceptionHandler struct {
	Begin   int
	End     int
	Catch   int
	Finally int
}

// LineNumber returns the line number for the opcode at the given instruction pointer
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
