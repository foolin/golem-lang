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

// Interpreter interprets Golem bc.
type Interpreter struct {
	modules    []*bc.Module
	modMap     map[string]*bc.Module
	builtInMgr g.BuiltinManager
	frames     []*frame
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
		modules:    modules,
		modMap:     modMap,
		builtInMgr: builtInMgr,
		frames:     []*frame{},
	}
}

// InitModules initializes each of the Modules.  Note that the modules
// are initialized in reverse order.
func (itp *Interpreter) InitModules() ([]g.Value, g.ErrorStruct) {

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

//-------------------------------------------------------------------------
// Eval

// Eval evaluates a Func.
func (itp *Interpreter) Eval(fn g.Func, params []g.Value) (g.Value, g.Error) {

	switch t := fn.(type) {

	case bc.Func:
		val, errStruct := itp.EvalBytecode(t, params)
		if errStruct != nil {
			return nil, errStruct.Error()
		}
		return val, nil

	case g.NativeFunc:
		return fn.Invoke(itp, params)

	default:
		panic("unreachable")
	}
}

// EvalBytecode evaluates a bytecode.Func
func (itp *Interpreter) EvalBytecode(fn bc.Func, params []g.Value) (g.Value, g.ErrorStruct) {
	val, err := itp.eval(fn, newLocals(fn.Template().NumLocals, params))
	if err != nil {
		return nil, err
	}
	return val, nil
}

//-------------------------------------------------------------------------

func (itp *Interpreter) eval(
	fn bc.Func,
	locals []*bc.Ref) (result g.Value, errStruct g.ErrorStruct) {

	lastFrame := len(itp.frames)
	itp.frames = append(itp.frames, &frame{fn, locals, []g.Value{}, 0})

	var err g.Error
	for result == nil {
		result, err = itp.advance(lastFrame)
		if err != nil {
			result, errStruct = itp.walkStack(g.NewErrorStruct(err, itp.stackTrace()))
			if errStruct != nil {
				return nil, errStruct
			}
		}
	}

	return result, nil
}

func (itp *Interpreter) walkStack(errStruct g.ErrorStruct) (g.Value, g.ErrorStruct) {

	// unwind the frames
	for len(itp.frames) > 0 {
		frameIndex := len(itp.frames) - 1
		f := itp.frames[frameIndex]
		instPtr := f.ip

		// visit exception handlers
		tpl := f.fn.Template()
		for _, eh := range tpl.ExceptionHandlers {

			// found an active handler
			if instPtr >= eh.Begin && instPtr < eh.End {

				if eh.Catch != -1 {
					f.ip = eh.Catch
					f.stack = append(f.stack, errStruct)
					cres, cerr := itp.runTryClause(f, frameIndex)
					if cerr != nil {

						// save the error
						errStruct = g.NewErrorStruct(cerr, itp.stackTrace())

						// run finally clause
						if eh.Finally != -1 {
							f.ip = eh.Finally
							fres, ferr := itp.runTryClause(f, frameIndex)
							if ferr != nil {
								// save the error
								errStruct = g.NewErrorStruct(ferr, itp.stackTrace())
							} else if fres != nil {
								// stop unwinding the stack
								return fres, nil
							}
						}

					} else {

						// run finally clause
						if eh.Finally != -1 {
							f.ip = eh.Finally
							fres, ferr := itp.runTryClause(f, frameIndex)
							if ferr == nil && fres != nil {
								// stop unwinding the stack
								return fres, nil
							}
						}

						// done!
						return cres, nil
					}
				} else {
					g.Assert(eh.Finally != -1)
					f.ip = eh.Finally
					fres, ferr := itp.runTryClause(f, frameIndex)
					if ferr != nil {
						// save the error
						errStruct = g.NewErrorStruct(ferr, itp.stackTrace())
					} else if fres != nil {
						// stop unwinding the stack
						return fres, nil
					}
				}
			}
		}

		// pop the frame
		itp.frames = itp.frames[:frameIndex]
	}

	return nil, errStruct
}

func (itp *Interpreter) runTryClause(f *frame, frameIndex int) (g.Value, g.Error) {

	btc := f.fn.Template().Bytecodes
	for btc[f.ip] != bc.Done {

		result, err := itp.advance(frameIndex)
		if result != nil || err != nil {
			return result, err
		}
	}
	f.ip++

	return nil, nil
}

func (itp *Interpreter) stackTrace() []string {

	n := len(itp.frames)
	stack := []string{}

	for i := n - 1; i >= 0; i-- {
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
