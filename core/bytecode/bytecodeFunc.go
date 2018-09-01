// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
)

// BytecodeFunc is a Func that is implemented in Golem
type BytecodeFunc interface {
	g.Func

	Template() *FuncTemplate
	GetCapture(int) *Ref
	PushCapture(*Ref)
}

type bytecodeFunc struct {
	template *FuncTemplate
	captures []*Ref
}

// NewBytecodeFunc creates a new BytecodeFunc.  NewBytecodeFunc is
// called via NewFunc opcode at runtime.
func NewBytecodeFunc(template *FuncTemplate) BytecodeFunc {
	captures := make([]*Ref, 0, template.NumCaptures)
	return &bytecodeFunc{template, captures}
}

func (f *bytecodeFunc) Type() g.Type { return g.FuncType }

func (f *bytecodeFunc) Freeze(ev g.Evaluator) (g.Value, g.Error) {
	return f, nil
}

func (f *bytecodeFunc) Frozen(ev g.Evaluator) (g.Bool, g.Error) {
	return g.True, nil
}

func (f *bytecodeFunc) HashCode(ev g.Evaluator) (g.Int, g.Error) {
	return nil, g.TypeMismatchError("Expected Hashable Type")
}

func (f *bytecodeFunc) Cmp(ev g.Evaluator, val g.Value) (g.Int, g.Error) {
	return nil, g.TypeMismatchError("Expected Comparable Type")
}

func (f *bytecodeFunc) ToStr(ev g.Evaluator) (g.Str, g.Error) {
	return g.NewStr(fmt.Sprintf("func<%p>", f)), nil
}

func (f *bytecodeFunc) Eq(ev g.Evaluator, val g.Value) (g.Bool, g.Error) {
	switch t := val.(type) {
	case BytecodeFunc:
		// equality is based on identity
		return g.NewBool(f == t), nil
	default:
		return g.False, nil
	}
}

//--------------------------------------------------------------
// fields

func (f *bytecodeFunc) FieldNames() ([]string, g.Error) {
	return []string{}, nil
}

func (f *bytecodeFunc) HasField(name string) (bool, g.Error) {
	return false, nil
}

func (f *bytecodeFunc) GetField(name string, ev g.Evaluator) (g.Value, g.Error) {
	return nil, g.NoSuchFieldError(name)
}

func (f *bytecodeFunc) InvokeField(name string, ev g.Evaluator, params []g.Value) (g.Value, g.Error) {
	return nil, g.NoSuchFieldError(name)
}

//--------------------------------------------------------------
// func

func (f *bytecodeFunc) Arity() g.Arity { return f.template.Arity }

func (f *bytecodeFunc) Invoke(ev g.Evaluator, params []g.Value) (g.Value, g.Error) {
	return ev.Eval(f, params)
}

//---------------------------------------------------------------

func (f *bytecodeFunc) Template() *FuncTemplate {
	return f.template
}

func (f *bytecodeFunc) GetCapture(idx int) *Ref {
	return f.captures[idx]
}

func (f *bytecodeFunc) PushCapture(ref *Ref) {
	f.captures = append(f.captures, ref)
}
