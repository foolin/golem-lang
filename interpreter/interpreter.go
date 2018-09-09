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

var debugInterpreter bool

func debugString(s string) {
	if debugInterpreter {
		fmt.Printf(s)
	}
}

// advance the interpreter forwards by one opcode.
func (itp *Interpreter) advance() (g.Value, g.Error) {

	f := itp.frameStack.peek()

	if debugInterpreter {
		fmt.Printf("=========================================\n")
		fmt.Printf("ip: %s\n", bc.FmtBytecode(f.btc, f.ip))
		fmt.Printf("\n")
		itp.frameStack.dump()
	}

	op := ops[f.btc[f.ip]]
	return op(itp, f)
}

func (itp *Interpreter) eval() (g.Value, ErrorStruct) {

	// advance until we get a res
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

func (itp *Interpreter) handleError(res g.Value, es ErrorStruct) (g.Value, ErrorStruct) {

	debugString(fmt.Sprintf("handleError es %v\n", es))

	// find an error handler
	h, ok := itp.frameStack.popErrorHandler()
	if !ok {
		return res, es
	}

	debugString(fmt.Sprintf("handleError popped %v\n", h))

	//-------------------------------------------

	f := itp.frameStack.peek()
	f.isHandlingError = true
	endIP := -1

	// catch
	var catch = clause{nil, nil, true}
	if h.Catch.IsValid() {
		endIP = h.Catch.End
		f.stack = append(f.stack, es)
		catch = itp.runTryClause(h.Catch)
		//debugString(fmt.Sprintf("handleError catch (%v, %v)\n", catchRes, catchEs))
	}

	// finally
	var finally = clause{nil, nil, true}
	if h.Finally.IsValid() {
		endIP = h.Finally.End
		finally = itp.runTryClause(h.Finally)
		//debugString(fmt.Sprintf("handleError finally (%v, %v)\n", finRes, finEs))
	}

	g.Assert(endIP != -1)
	f.isHandlingError = false

	//-------------------------------------------

	if !finally.ended {
		debugString(fmt.Sprintf("handleError fin did not end\n"))
		if finally.res != nil {
			f.stack = append(f.stack, finally.res)
		}
		return itp.handleError(finally.res, finally.es)
	}

	if !catch.ended {
		debugString(fmt.Sprintf("handleError catch did not end\n"))
		if catch.res != nil {
			f.stack = append(f.stack, catch.res)
		}
		return itp.handleError(catch.res, catch.es)
	}

	// carry on inside the current frame
	f.ip = endIP
	debugString(fmt.Sprintf("handleError carry on %d\n", f.ip))
	return itp.eval()
}

type clause struct {
	res   g.Value
	es    ErrorStruct
	ended bool
}

func (itp *Interpreter) runTryClause(tc bc.TryClause) clause {

	debugString(fmt.Sprintf("runTryClause %d, %d\n", tc.Begin, tc.End))

	itp.frameStack.peek().ip = tc.Begin

	for {
		f := itp.frameStack.peek()
		if f.isHandlingError && f.ip >= tc.End {
			g.Assert(f.ip == tc.End)
			return clause{nil, nil, true}
		}

		res, err := itp.advance()
		if err != nil {
			es := newErrorStruct(err, itp.frameStack.stackTrace())
			return clause{nil, es, false}
		}
		if res != nil {
			return clause{res, nil, false}
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
