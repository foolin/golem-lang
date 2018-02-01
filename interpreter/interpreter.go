// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
	o "github.com/mjarmy/golem-lang/core/opcodes"
)

//---------------------------------------------------------------
// The Golem Interpreter

// Interpreter interprets Golem bytecode.
type Interpreter struct {
	mod        *g.BytecodeModule
	builtInMgr g.BuiltinManager
	frames     []*frame
}

// NewInterpreter creates a new Interpreter
func NewInterpreter(mod *g.BytecodeModule, builtInMgr g.BuiltinManager) *Interpreter {
	return &Interpreter{mod, builtInMgr, []*frame{}}
}

// Init initializes an interpreter, by interpreting its "init" function.
func (i *Interpreter) Init() (g.Value, g.Error) {

	// the init function is always the zeroth template
	tpl := i.mod.Templates[0]
	//tpl := mod.Templates[0]
	if tpl.Arity != 0 || tpl.NumCaptures != 0 {
		panic("TODO")
	}

	// create empty locals
	i.mod.Refs = newLocals(tpl.NumLocals, nil)

	// make func
	fn := g.NewBytecodeFunc(tpl)

	// go
	return i.eval(fn, i.mod.Refs)
}

// Eval evaluates a given BytecodeFunc.  Note that this function has the same
// signature as core.Context.Eval()
func (i *Interpreter) Eval(fn g.BytecodeFunc, params []g.Value) (g.Value, g.Error) {
	return i.eval(fn, newLocals(fn.Template().NumLocals, params))
}

func (i *Interpreter) eval(fn g.BytecodeFunc, locals []*g.Ref) (result g.Value, errTrace g.Error) {

	lastFrame := len(i.frames)
	i.frames = append(i.frames, &frame{fn, locals, []g.Value{}, 0})

	var err g.Error
	for result == nil {
		result, err = i.advance(lastFrame)
		if err != nil {
			result, errTrace = i.walkStack(i.makeErrorTrace(err, i.stackTrace()))
			if errTrace != nil {
				return nil, errTrace
			}
		}
	}

	return result, nil
}

func (i *Interpreter) walkStack(errTrace g.Error) (g.Value, g.Error) {

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
					f.stack = append(f.stack, errTrace.Struct())
					cres, cerr := i.runTryClause(f, frameIndex)
					if cerr != nil {
						// save the error
						errTrace = i.makeErrorTrace(cerr, i.stackTrace())

						// run finally clause
						if eh.Finally != -1 {
							f.ip = eh.Finally
							fres, ferr := i.runTryClause(f, frameIndex)
							if ferr != nil {
								// save the error
								errTrace = i.makeErrorTrace(ferr, i.stackTrace())
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
							if ferr != nil {
								// save the error
								errTrace = i.makeErrorTrace(ferr, i.stackTrace())
							} else if fres != nil {
								// stop unwinding the stack
								return fres, nil
							}
						}

						// done!
						return cres, nil
					}
				} else {
					assert(eh.Finally != -1)
					f.ip = eh.Finally
					fres, ferr := i.runTryClause(f, frameIndex)
					if ferr != nil {
						// save the error
						errTrace = i.makeErrorTrace(ferr, i.stackTrace())
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

	return nil, errTrace
}

func (i *Interpreter) runTryClause(f *frame, frameIndex int) (g.Value, g.Error) {

	opc := f.fn.Template().OpCodes
	for opc[f.ip] != o.Done {

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
		stack = append(stack, fmt.Sprintf("    at line %d", lineNum))
	}

	return stack
}

func newLocals(numLocals int, params []g.Value) []*g.Ref {
	p := len(params)
	locals := make([]*g.Ref, numLocals, numLocals)
	for j := 0; j < numLocals; j++ {
		if j < p {
			locals[j] = &g.Ref{params[j]}
		} else {
			locals[j] = &g.Ref{g.NullValue}
		}
	}
	return locals
}

func (i *Interpreter) dump() {

	println("-----------------------------------------")

	f := i.frames[len(i.frames)-1]
	opc := f.fn.Template().OpCodes
	print(o.FmtOpcode(opc, f.ip))

	for j, f := range i.frames {
		fmt.Printf("frame %d\n", j)
		i.dumpFrame(f)
	}
}

//---------------------------------------------------------------
// An execution environment, a.k.a 'stack frame'.

type frame struct {
	fn     g.BytecodeFunc
	locals []*g.Ref
	stack  []g.Value
	ip     int
}

func (i *Interpreter) dumpFrame(f *frame) {
	fmt.Printf("    locals:\n")
	for j, r := range f.locals {
		fmt.Printf("        %d: %s\n", j, r.Val.ToStr(i))
	}
	fmt.Printf("    stack:\n")
	for j, v := range f.stack {
		fmt.Printf("        %d: %s\n", j, v.ToStr(i))
	}
	fmt.Printf("    ip: %d\n", f.ip)
}

//---------------------------------------------------------------

func (i *Interpreter) makeErrorTrace(err g.Error, stackTrace []string) g.Error {

	// make list-of-str
	vals := make([]g.Value, len(stackTrace), len(stackTrace))
	for i, s := range stackTrace {
		vals[i] = g.NewStr(s)
	}
	list := g.NewList(vals)
	list.Freeze()

	stc, e := g.NewStruct([]g.Field{g.NewField("stackTrace", true, list)}, true)
	assert(e == nil)

	merge := g.MergeStructs([]g.Struct{err.Struct(), stc})
	return g.NewErrorFromStruct(i, merge)
}

func assert(flag bool) {
	if !flag {
		panic("assertion failure")
	}
}
