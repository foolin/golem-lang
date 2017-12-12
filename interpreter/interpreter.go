// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interpreter

import (
	"fmt"
	g "github.com/mjarmy/golem-lang/core"
)

//---------------------------------------------------------------
// The Golem Interpreter

type Interpreter struct {
	mod    *g.BytecodeModule
	frames []*frame
}

func NewInterpreter(mod *g.BytecodeModule) *Interpreter {
	return &Interpreter{mod, []*frame{}}
}

func (i *Interpreter) Init() (g.Value, *ErrorTrace) {

	// use the zeroth template
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
	return i.run(fn, i.mod.Refs)
}

func (i *Interpreter) RunBytecode(
	fn g.BytecodeFunc, params []g.Value) (result g.Value, errTrace *ErrorTrace) {

	return i.run(fn, newLocals(fn.Template().NumLocals, params))
}

func (i *Interpreter) run(
	fn g.BytecodeFunc, locals []*g.Ref) (result g.Value, errTrace *ErrorTrace) {

	i.frames = append(i.frames, &frame{fn, locals, []g.Value{}, 0})

	var err g.Error
	for result == nil {
		result, err = i.advance(0)
		if err != nil {
			result, errTrace = i.walkStack(makeErrorTrace(err, i.stackTrace()))
			if errTrace != nil {
				return nil, errTrace
			}
		}
	}

	return result, nil
}

func (i *Interpreter) walkStack(errTrace *ErrorTrace) (g.Value, *ErrorTrace) {

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
					f.stack = append(f.stack, errTrace.Struct)
					cres, cerr := i.runTryClause(f, frameIndex)
					if cerr != nil {
						// save the error
						errTrace = makeErrorTrace(cerr, i.stackTrace())

						// run finally clause
						if eh.Finally != -1 {
							f.ip = eh.Finally
							fres, ferr := i.runTryClause(f, frameIndex)
							if ferr != nil {
								// save the error
								errTrace = makeErrorTrace(ferr, i.stackTrace())
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
								errTrace = makeErrorTrace(ferr, i.stackTrace())
							} else if fres != nil {
								// stop unwinding the stack
								return fres, nil
							}
						}

						// done!
						return cres, nil
					}
				} else {
					g.Assert(eh.Finally != -1, "invalid try")
					f.ip = eh.Finally
					fres, ferr := i.runTryClause(f, frameIndex)
					if ferr != nil {
						// save the error
						errTrace = makeErrorTrace(ferr, i.stackTrace())
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
	for opc[f.ip] != g.DONE {

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
			locals[j] = &g.Ref{g.NULL}
		}
	}
	return locals
}

func (i *Interpreter) dump() {

	println("-----------------------------------------")

	f := i.frames[len(i.frames)-1]
	opc := f.fn.Template().OpCodes
	print(g.FmtOpcode(opc, f.ip))

	for j, f := range i.frames {
		fmt.Printf("frame %d\n", j)
		f.dump()
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

func (f *frame) dump() {
	fmt.Printf("    locals:\n")
	for j, r := range f.locals {
		fmt.Printf("        %d: %s\n", j, r.Val.ToStr())
	}
	fmt.Printf("    stack:\n")
	for j, v := range f.stack {
		fmt.Printf("        %d: %s\n", j, v.ToStr())
	}
	fmt.Printf("    ip: %d\n", f.ip)
}

//---------------------------------------------------------------
// A combination of an error, and a stack trace

func makeErrorTrace(err g.Error, stackTrace []string) *ErrorTrace {

	// make list-of-str
	vals := make([]g.Value, len(stackTrace), len(stackTrace))
	for i, s := range stackTrace {
		vals[i] = g.MakeStr(s)
	}
	// TODO make the list immutable
	list := g.NewList(vals)

	stc, e := g.NewStruct([]*g.StructEntry{{"stackTrace", true, false, list}})
	g.Assert(e == nil, "invalid struct")

	merge := g.MergeStructs([]g.Struct{err.Struct(), stc})
	return &ErrorTrace{err, stackTrace, merge}
}

type ErrorTrace struct {
	Error      g.Error
	StackTrace []string
	Struct     g.Struct
}
