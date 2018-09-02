// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/core/bytecode"
)

//---------------------------------------------------------------
// The Golem Interpreter

// Interpreter interprets Golem bytecode.
type Interpreter struct {
	modules    []*bytecode.Module
	modMap     map[string]*bytecode.Module
	builtInMgr g.BuiltinManager
	frames     []*frame
}

// NewInterpreter creates a new Interpreter
func NewInterpreter(
	builtInMgr g.BuiltinManager,
	modules []*bytecode.Module) *Interpreter {

	modMap := make(map[string]*bytecode.Module)
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
func (i *Interpreter) InitModules() ([]g.Value, g.ErrorStruct) {

	result := []g.Value{}
	for j := len(i.modules) - 1; j >= 0; j-- {
		mod := i.modules[j]

		// the 'init' function is always the first template in the pool
		initTpl := mod.Pool.Templates[0]

		// create empty locals
		mod.Refs = newLocals(initTpl.NumLocals, nil)

		// make init function from template
		initFn := bytecode.NewBytecodeFunc(initTpl)

		// invoke the "init" function
		val, err := i.eval(initFn, mod.Refs)
		if err != nil {
			return nil, err
		}

		// prepend the value so that the result will be in the same order as i.modules
		result = append([]g.Value{val}, result...)
	}
	return result, nil
}

//-------------------------------------------------------------------------
// Evaluator

// Eval evaluates a Func.
func (i *Interpreter) Eval(fn g.Func, params []g.Value) (g.Value, g.Error) {

	switch t := fn.(type) {

	case bytecode.BytecodeFunc:
		val, errStruct := i.EvalBytecode(t, params)
		if errStruct != nil {
			return nil, errStruct.Error()
		}
		return val, nil

	case g.NativeFunc:
		return fn.Invoke(i, params)

	default:
		panic("unreachable")
	}
}

func (i *Interpreter) EvalBytecode(fn bytecode.BytecodeFunc, params []g.Value) (g.Value, g.ErrorStruct) {
	val, err := i.eval(fn, newLocals(fn.Template().NumLocals, params))
	if err != nil {
		return nil, err
	}
	return val, nil
}

//-------------------------------------------------------------------------

func (i *Interpreter) eval(
	fn bytecode.BytecodeFunc,
	locals []*bytecode.Ref) (result g.Value, errStruct g.ErrorStruct) {

	lastFrame := len(i.frames)
	i.frames = append(i.frames, &frame{fn, locals, []g.Value{}, 0})

	var err g.Error
	for result == nil {
		result, err = i.advance(lastFrame)
		if err != nil {
			result, errStruct = i.walkStack(g.NewErrorStruct(err, i.stackTrace()))
			if errStruct != nil {
				return nil, errStruct
			}
		}
	}

	return result, nil
}

func (i *Interpreter) walkStack(errStruct g.ErrorStruct) (g.Value, g.ErrorStruct) {

	// unwind the frames
	for len(i.frames) > 0 {
		frameIndex := len(i.frames) - 1
		f := i.frames[frameIndex]
		instPtr := f.ip

		// visit exception handlers
		tpl := f.fn.Template()
		for _, eh := range tpl.ExceptionHandlers {

			// found an active handler
			if instPtr >= eh.Begin && instPtr < eh.End {

				if eh.Catch != -1 {
					f.ip = eh.Catch
					f.stack = append(f.stack, errStruct)
					cres, cerr := i.runTryClause(f, frameIndex)
					if cerr != nil {

						// save the error
						errStruct = g.NewErrorStruct(cerr, i.stackTrace())

						// run finally clause
						if eh.Finally != -1 {
							f.ip = eh.Finally
							fres, ferr := i.runTryClause(f, frameIndex)
							if ferr != nil {
								// save the error
								errStruct = g.NewErrorStruct(ferr, i.stackTrace())
							} else if fres != nil {
								// stop unwinding the stack
								return fres, nil
							}
						}

					} else {

						// run finally clause
						if eh.Finally != -1 {
							f.ip = eh.Finally
							fres, ferr := i.runTryClause(f, frameIndex)
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
					fres, ferr := i.runTryClause(f, frameIndex)
					if ferr != nil {
						// save the error
						errStruct = g.NewErrorStruct(ferr, i.stackTrace())
					} else if fres != nil {
						// stop unwinding the stack
						return fres, nil
					}
				}
			}
		}

		// pop the frame
		i.frames = i.frames[:frameIndex]
	}

	return nil, errStruct
}

func (i *Interpreter) runTryClause(f *frame, frameIndex int) (g.Value, g.Error) {

	btc := f.fn.Template().Bytecodes
	for btc[f.ip] != bytecode.Done {

		result, err := i.advance(frameIndex)
		if result != nil || err != nil {
			return result, err
		}
	}
	f.ip++

	return nil, nil
}

func (i *Interpreter) stackTrace() []string {

	n := len(i.frames)
	stack := []string{}

	for j := n - 1; j >= 0; j-- {
		tpl := i.frames[j].fn.Template()
		lineNum := tpl.LineNumber(i.frames[j].ip)
		stack = append(stack, fmt.Sprintf("    at %s:%d", tpl.Module.Path, lineNum))
	}

	return stack
}

func newLocals(numLocals int, params []g.Value) []*bytecode.Ref {
	p := len(params)
	locals := make([]*bytecode.Ref, numLocals)
	for j := 0; j < numLocals; j++ {
		if j < p {
			locals[j] = &bytecode.Ref{Val: params[j]}
		} else {
			locals[j] = &bytecode.Ref{Val: g.Null}
		}
	}
	return locals
}

func (i *Interpreter) lookupModule(name string) (*bytecode.Module, g.Error) {
	if mod, ok := i.modMap[name]; ok {
		return mod, nil
	}
	return nil, g.UndefinedModuleError(name)
}
