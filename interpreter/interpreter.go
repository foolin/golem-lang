// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	//"fmt"

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

	values := []g.Value{}
	for i := len(itp.modules) - 1; i >= 0; i-- {
		mod := itp.modules[i]

		// the 'init' function is always the first template in the pool
		initTpl := mod.Pool.Templates[0]

		// create empty locals
		mod.Refs = newLocals(initTpl.NumLocals, nil)

		// make init function from template
		initFn := bc.NewBytecodeFunc(initTpl)

		// invoke the "init" function
		itp.frameStack.push(newFrame(initFn, mod.Refs, true))
		val, es := itp.eval()
		if es != nil {
			return nil, es
		}

		// prepend the value so that the values will be in the same order as itp.modules
		values = append([]g.Value{val}, values...)
	}
	return values, nil
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

// EvalBytecode evaluates a bytecode.Func.
func (itp *Interpreter) EvalBytecode(fn bc.Func, params []g.Value) (g.Value, ErrorStruct) {

	locals := newLocals(fn.Template().NumLocals, params)
	itp.frameStack.push(newFrame(fn, locals, true))
	val, es := itp.eval()
	if es != nil {
		return nil, es
	}
	return val, nil
}

//-------------------------------------------------------------------------

// create new local refs for use by a frame
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

// evaluate bytecode until we get a result or an error
func (itp *Interpreter) eval() (g.Value, ErrorStruct) {

	var res g.Value
	var err g.Error
	for res == nil {
		res, err = itp.advance()

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
			res, es = itp.handleError(nil, es)
			if es != nil {
				return nil, es
			}
		}
	}

	// success
	return res, nil
}

// deal with an error that was generated by a bytecode operation
func (itp *Interpreter) handleError(res g.Value, es ErrorStruct) (g.Value, ErrorStruct) {

	g.Assert(res == nil || es == nil)
	g.Assert(!(res == nil && es == nil))

	//-------------------------------------------
	// find an error handler

	h, ok := itp.frameStack.popErrorHandler()
	if !ok {
		return res, es
	}

	//-------------------------------------------
	// run the 'catch' and 'finally' clauses

	f := itp.frameStack.peek()
	f.isHandlingError = true
	endIP := -1

	// run catch
	catch := runClause{}
	if !h.Catch.IsEmpty() {
		endIP = h.Catch.End
		f.stack = append(f.stack, es) // put error on stack
		catch = itp.runTryClause(h.Catch)
	}

	// run finally
	finally := runClause{}
	if !h.Finally.IsEmpty() {
		endIP = h.Finally.End
		finally = itp.runTryClause(h.Finally)
	}

	f.isHandlingError = false

	//-------------------------------------------
	// figure out how to proceed

	// finally
	if finally.finishedEarly {
		if finally.result != nil {
			f.stack = append(f.stack, finally.result)
		}
		return itp.handleError(finally.result, finally.err)
	}

	// catch
	if catch.finishedEarly {
		if catch.result != nil {
			f.stack = append(f.stack, catch.result)
		}
		return itp.handleError(catch.result, catch.err)
	}

	// error recovery was successful
	f.ip = endIP
	return itp.eval()
}

type runClause struct {
	result        g.Value
	err           ErrorStruct
	finishedEarly bool
}

// run a 'catch' or 'finally' clause
func (itp *Interpreter) runTryClause(tc bc.TryClause) runClause {

	itp.frameStack.peek().ip = tc.Begin

	for {
		f := itp.frameStack.peek()
		if f.isHandlingError && f.ip >= tc.End {
			g.Assert(f.ip == tc.End)
			return runClause{nil, nil, false}
		}

		res, err := itp.advance()
		if err != nil {
			es := newErrorStruct(err, itp.frameStack.stackTrace())
			return runClause{nil, es, true}
		}
		if res != nil {
			return runClause{res, nil, true}
		}
	}
}

//var debugInterpreter bool
//
//func debugString(s string) {
//	if debugInterpreter {
//		fmt.Printf(s)
//	}
//}

// advance the interpreter forwards by one opcode.
func (itp *Interpreter) advance() (g.Value, g.Error) {

	f := itp.frameStack.peek()

	//if debugInterpreter {
	//	fmt.Printf("=========================================\n")
	//	fmt.Printf("ip: %s\n", bc.FmtBytecode(f.btc, f.ip))
	//	fmt.Printf("\n")
	//	itp.frameStack.dump()
	//}

	op := ops[f.btc[f.ip]]
	return op(itp, f)
}
