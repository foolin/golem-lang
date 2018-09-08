// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
)

//---------------------------------------------------------------
// The Golem Interpreter
//---------------------------------------------------------------

// Interpreter interprets Golem bytecode.
type Interpreter struct {
	builtInMgr g.BuiltinManager
	modules    []*bc.Module
	modMap     map[string]*bc.Module

	frames []*frame
}

// NewInterpreter creates a new Interpreter
func NewInterpreter(
	builtInMgr g.BuiltinManager,
	modules []*bc.Module) *Interpreter {

	modMap := make(map[string]*bc.Module)
	for _, m := range modules {
		modMap[m.Name] = m
	}

	return &Interpreter{
		builtInMgr: builtInMgr,
		modules:    modules,
		modMap:     modMap,
		frames:     []*frame{},
	}
}

// InitModules initializes each of the Modules.  Note that the modules
// are initialized in reverse order.
func (itp *Interpreter) InitModules() ([]g.Value, ErrorStruct) {

	result := []g.Value{}
	for i := len(itp.modules) - 1; i >= 0; i-- {
		mod := itp.modules[i]

		// the 'init' function is always the first template in the pool
		initTpl := mod.Pool.Templates[0]

		// create empty locals
		mod.Refs = newLocals(initTpl.NumLocals, nil)

		// make init function from template
		initFn := bc.NewBytecodeFunc(initTpl)

		// invoke the "init" function
		val, err := itp.eval(initFn, mod.Refs)
		if err != nil {
			return nil, err
		}

		// prepend the value so that the result will be in the same order as itp.modules
		result = append([]g.Value{val}, result...)
	}
	return result, nil
}

// Eval evaluates a Func.
func (itp *Interpreter) Eval(fn g.Func, params []g.Value) (g.Value, g.Error) {

	switch t := fn.(type) {

	case bc.Func:
		val, es := itp.EvalBytecode(t, params)
		if es != nil {
			return nil, es.Error()
		}
		return val, nil

	case g.NativeFunc:
		return fn.Invoke(itp, params)

	default:
		panic("unreachable")
	}
}

// EvalBytecode evaluates a bytecode.Func
func (itp *Interpreter) EvalBytecode(fn bc.Func, params []g.Value) (g.Value, ErrorStruct) {

	val, err := itp.eval(fn, newLocals(fn.Template().NumLocals, params))
	if err != nil {
		return nil, err
	}
	return val, nil
}

//-------------------------------------------------------------------------

func (itp *Interpreter) numFrames() int {
	return len(itp.frames)
}

func (itp *Interpreter) peekFrame() *frame {
	return itp.frames[itp.numFrames()-1]
}

func (itp *Interpreter) pushFrame(f *frame) {
	itp.frames = append(itp.frames, f)
}

func (itp *Interpreter) popFrame() {
	itp.frames = itp.frames[:itp.numFrames()-1]
}

// advance the interpreter forwards by one opcode.
func (itp *Interpreter) advance() (g.Value, g.Error) {
	f := itp.peekFrame()
	op := ops[f.btc[f.ip]]
	return op(itp, f)
}

func (itp *Interpreter) eval(fn bc.Func, locals []*bc.Ref) (g.Value, ErrorStruct) {

	itp.frames = append(itp.frames, newFrame(fn, locals, true))

	var result g.Value
	var err g.Error
	var es ErrorStruct

	// advance until we get a result
	for result == nil {
		result, err = itp.advance()

		// an error was thrown
		if err != nil {
			es = newErrorStruct(err, itp.stackTrace())
			result, es = itp.handleError(es)
			if es != nil {
				return nil, es
			}
		}
	}

	// success
	return result, nil
}

func (itp *Interpreter) popErrorHandler() *bc.ErrorHandler {

	f := itp.peekFrame()
	if f.numHandlers() > 0 {
		return f.popHandler()
	}

	for !f.isBase {
		itp.popFrame()

		f = itp.peekFrame()
		if f.numHandlers() > 0 {
			return f.popHandler()
		}
	}

	itp.popFrame()
	return nil
}

func (itp *Interpreter) handleError(es ErrorStruct) (g.Value, ErrorStruct) {

	h := itp.popErrorHandler()
	if h == nil {
		return nil, es
	}

	f := itp.peekFrame()
	var result g.Value

	// catch
	if h.Catch != -1 {
		f.stack = append(f.stack, es)
		result, es = itp.runTryClause(h.Catch)
	}

	// finally
	if h.Finally != -1 {
		result, es = itp.runTryClause(h.Finally)
	}

	// found error
	if es != nil {
		panic("TODO")
	}

	// carry on inside the current frame
	f.ip = h.End
	return result, nil
}

func (itp *Interpreter) runTryClause(begin int) (g.Value, ErrorStruct) {

	f := itp.peekFrame()
	f.ip = begin
	for f.btc[f.ip] != bc.TryDone {

		result, err := itp.advance()
		if result != nil || err != nil {
			return result, newErrorStruct(err, itp.stackTrace())
		}

		f = itp.peekFrame()
	}

	return nil, nil
}

func (itp *Interpreter) stackTrace() []string {

	stack := []string{}

	for i := itp.numFrames() - 1; i >= 0; i-- {
		tpl := itp.frames[i].fn.Template()
		lineNum := tpl.LineNumber(itp.frames[i].ip)
		stack = append(stack, fmt.Sprintf("    at %s:%d", tpl.Module.Path, lineNum))
	}

	return stack
}

func newLocals(numLocals int, params []g.Value) []*bc.Ref {
	p := len(params)
	locals := make([]*bc.Ref, numLocals)
	for i := 0; i < numLocals; i++ {
		if i < p {
			locals[i] = &bc.Ref{Val: params[i]}
		} else {
			locals[i] = &bc.Ref{Val: g.Null}
		}
	}
	return locals
}

func (itp *Interpreter) lookupModule(name string) (*bc.Module, g.Error) {
	if mod, ok := itp.modMap[name]; ok {
		return mod, nil
	}
	return nil, g.UndefinedModule(name)
}
