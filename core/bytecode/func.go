// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
)

/*doc
## Func

A function is a sequence of [expressions](syntax.html#expressions)
and [statements](syntax.html#statements) that can be invoked to perform a task.

Valid operators for Func are:

* The equality operators `==`, `!=`
* The invocation operator `a(x)`

Funcs have an [`arity()`](builtins.html#arity)

Funcs have no fields.

*/

// Func is a core.Func that is implemented in Golem
type Func interface {
	g.Func

	Template() *FuncTemplate
	GetCapture(int) *Ref
	PushCapture(*Ref)
}

type bytecodeFunc struct {
	template *FuncTemplate
	captures []*Ref
}

// NewBytecodeFunc creates a new Func.  NewBytecodeFunc is
// called via NewFunc opcode at runtime.
func NewBytecodeFunc(template *FuncTemplate) Func {
	captures := make([]*Ref, 0, template.NumCaptures)
	return &bytecodeFunc{template, captures}
}

func (f *bytecodeFunc) Type() g.Type { return g.FuncType }

func (f *bytecodeFunc) Freeze(ev g.Eval) (g.Value, g.Error) {
	return f, nil
}

func (f *bytecodeFunc) Frozen(ev g.Eval) (g.Bool, g.Error) {
	return g.True, nil
}

func (f *bytecodeFunc) HashCode(ev g.Eval) (g.Int, g.Error) {
	return nil, g.HashCodeMismatch(g.FuncType)
}

func (f *bytecodeFunc) ToStr(ev g.Eval) (g.Str, g.Error) {
	return g.NewStr(fmt.Sprintf("func<%p>", f))
}

func (f *bytecodeFunc) Eq(ev g.Eval, val g.Value) (g.Bool, g.Error) {
	switch t := val.(type) {
	case Func:
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

func (f *bytecodeFunc) GetField(ev g.Eval, name string) (g.Value, g.Error) {
	return nil, g.NoSuchField(name)
}

func (f *bytecodeFunc) InvokeField(ev g.Eval, name string, params []g.Value) (g.Value, g.Error) {
	return nil, g.NoSuchField(name)
}

//--------------------------------------------------------------
// func

func (f *bytecodeFunc) Arity() g.Arity { return f.template.Arity }

func (f *bytecodeFunc) Invoke(ev g.Eval, params []g.Value) (g.Value, g.Error) {
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
