// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/ncore"
)

// BytecodeFunc represents a function that is defined
// via Golem source code
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

func (f *bytecodeFunc) Freeze(cx g.Context) (g.Value, g.Error) {
	return f, nil
}

func (f *bytecodeFunc) Frozen(cx g.Context) (g.Bool, g.Error) {
	return g.True, nil
}

func (f *bytecodeFunc) HashCode(cx g.Context) (g.Int, g.Error) {
	return nil, g.TypeMismatchError("Expected Hashable Type")
}

//func (f *bytecodeFunc) GetField(cx g.Context, key g.Str) (g.Value, g.Error) {
//	return nil, NoSuchFieldError(key.String())
//}

func (f *bytecodeFunc) Cmp(cx g.Context, v g.Value) (g.Int, g.Error) {
	return nil, g.TypeMismatchError("Expected Comparable Type")
}

func (f *bytecodeFunc) ToStr(cx g.Context) g.Str {
	return g.NewStr(fmt.Sprintf("func<%p>", f))
}

func (f *bytecodeFunc) Eq(cx g.Context, v g.Value) (g.Bool, g.Error) {
	switch t := v.(type) {
	case BytecodeFunc:
		// equality is based on identity
		return g.NewBool(f == t), nil
	default:
		return g.False, nil
	}
}

func (f *bytecodeFunc) Arity() *g.Arity { return f.template.Arity }

func (f *bytecodeFunc) Invoke(cx g.Context, values []g.Value) (g.Value, g.Error) {
	return cx.Eval(f, values)
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
