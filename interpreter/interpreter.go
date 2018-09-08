// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
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

	frameStack *frameStack
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
		frameStack: newFrameStack(),
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
		val, es := itp.eval(initFn, mod.Refs)
		if es != nil {
			return nil, es
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
			return nil, es
		}
		return val, nil

	case g.NativeFunc:
		return fn.Invoke(itp, params)

	default:
		panic("unreachable")
	}
}

func (itp *Interpreter) EvalBytecode(fn bc.Func, params []g.Value) (g.Value, ErrorStruct) {

	val, es := itp.eval(fn, newLocals(fn.Template().NumLocals, params))
	if es != nil {
		return nil, es
	}
	return val, nil
}

//-------------------------------------------------------------------------

// advance the interpreter forwards by one opcode.
func (itp *Interpreter) advance() (g.Value, g.Error) {

	f := itp.frameStack.peek()
	op := ops[f.btc[f.ip]]
	return op(itp, f)
}

func (itp *Interpreter) eval(fn bc.Func, locals []*bc.Ref) (g.Value, ErrorStruct) {

	itp.frameStack.push(newFrame(fn, locals, true))

	// advance until we get a result
	var result g.Value
	var err g.Error
	for result == nil {
		result, err = itp.advance()

		// an error was thrown
		if err != nil {

			// If the error is already an ErrorStruct, that means it is being
			// propagated back up by unwinding recursive calls to Eval().  In that
			// case we should preserve the ErrorStruct, since we want to pass its
			// stack trace all the way back up.
			es, ok := err.(ErrorStruct)
			if !ok {
				// If the error not already an ErrorStruct, then we create one.
				es = newErrorStruct(err, itp.frameStack.stackTrace())
			}

			// handle the error
			result, es = itp.handleError(es)
			if es != nil {
				return nil, es
			}
		}
	}

	// success
	return result, nil
}

func (itp *Interpreter) handleError(es ErrorStruct) (g.Value, ErrorStruct) {

	h, ok := itp.frameStack.popErrorHandler()
	if !ok {
		return nil, es
	}

	//-------------------------------------------

	f := itp.frameStack.peek()
	f.isHandlingError = true
	endIp := -1
	var result g.Value

	// catch
	if h.CatchBegin != -1 {
		endIp = h.CatchEnd
		f.stack = append(f.stack, es)
		result, es = itp.runTryClause(h.CatchBegin, h.CatchEnd)
	}

	// finally
	if h.FinallyBegin != -1 {
		endIp = h.FinallyEnd
		result, es = itp.runTryClause(h.FinallyBegin, h.FinallyEnd)
	}

	g.Assert(endIp != -1)
	f.isHandlingError = false

	//-------------------------------------------

	if es != nil {
		panic("TODO")
	}
	if result != nil {
		panic("TODO")
	}

	// carry on inside the current frame
	f.ip = endIp
	return nil, nil
}

func (itp *Interpreter) runTryClause(beginIp, endIp int) (g.Value, ErrorStruct) {

	itp.frameStack.peek().ip = beginIp

	var result g.Value
	var err g.Error
	for {
		f := itp.frameStack.peek()
		if f.isHandlingError && f.ip >= endIp {
			g.Assert(f.ip == endIp)
			return nil, nil
		}

		result, err = itp.advance()
		if err != nil {
			return nil, newErrorStruct(err, itp.frameStack.stackTrace())
		}
		if result != nil {
			return result, nil
		}
	}
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
