// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
	"github.com/mjarmy/golem-lang/scanner"
)

//---------------------------------------------------------------
// The Golem Interpreter
//---------------------------------------------------------------

// EvalCode is a convenience function for compiling Golem source code into
// a bytecode.Module, and then interpreting the module.
func EvalCode(
	code string,
	builtins []*g.Builtin,
	importer Importer) (g.Value, g.Error) {

	mod, err := CompileCode(code, builtins)
	if err != nil {
		return nil, err
	}

	itp := NewInterpreter(builtins, importer)
	return itp.EvalModule(mod)
}

// CompileCode is a convenience function for compiling Golem source code into
// a bytecode.Module.
func CompileCode(
	code string,
	builtins []*g.Builtin) (*bc.Module, g.Error) {

	source := &scanner.Source{Name: "", Path: "", Code: code}
	mod, err := compiler.CompileSource(source, builtins)
	if err != nil {
		return nil, g.Error(err)
	}

	return mod, nil
}

// Interpreter interprets Golem bytecode.
type Interpreter struct {
	builtins   []*g.Builtin
	importer   Importer
	frameStack *frameStack
}

// NewInterpreter creates a new Interpreter.  If the importer is nil,
// then no modules can be imported.
func NewInterpreter(builtins []*g.Builtin, importer Importer) *Interpreter {

	return &Interpreter{
		builtins:   builtins,
		importer:   importer,
		frameStack: newFrameStack(),
	}
}

// EvalModule evaluates a bytecode.Module by calling its wrapper "init" Func.
func (itp *Interpreter) EvalModule(mod *bc.Module) (g.Value, ErrorStruct) {

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
	return val, nil
}

// Eval evaluates a Func.  Note that this method causes Interpreter
// to implement the core.Eval interface.
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

	//-------------------------------------------
	// find an error handler

	h, ok := itp.frameStack.popErrorHandler()
	if !ok {
		return res, es
	}

	//-------------------------------------------
	// run the catch and finally clauses

	f := itp.frameStack.peek()
	f.isHandlingError = true

	responses := []*response{}
	endIP := -1

	if !h.Catch.IsEmpty() {
		endIP = h.Catch.End

		f.stack = append(f.stack, es) // put error on stack
		r := itp.runTryClause(h.Catch)
		if r != nil {
			responses = append(responses, r)
		}
	}

	if !h.Finally.IsEmpty() {
		endIP = h.Finally.End

		r := itp.runTryClause(h.Finally)
		if r != nil {
			responses = append(responses, r)
		}
	}

	f.isHandlingError = false

	//-------------------------------------------
	// figure out how to proceed

	// error recovery was successful
	if len(responses) == 0 {
		f.ip = endIP
		return itp.eval()
	}

	r := responses[len(responses)-1]

	// result: invoke the return op
	if r.result != nil {
		f.stack = append(f.stack, r.result)
		f.ip = r.resultIP
		g.Assert(f.btc[f.ip] == bc.Return)
		return itp.eval()
	}

	// error: throw an error
	return itp.handleError(nil, r.es)
}

type response struct {
	result   g.Value
	resultIP int
	es       ErrorStruct
}

// run a 'catch' or 'finally' clause
func (itp *Interpreter) runTryClause(tc bc.TryClause) *response {

	itp.frameStack.peek().ip = tc.Begin

	for {
		f := itp.frameStack.peek()

		// we've reached the end of the clause
		if f.isHandlingError && f.ip >= tc.End {
			g.Assert(f.ip == tc.End)
			return nil
		}

		// there is an explicit return inside the clause
		if f.isHandlingError && f.btc[f.ip] == bc.Return {
			n := len(f.stack) - 1
			res := f.stack[n]
			f.stack = f.stack[:n]
			return &response{res, f.ip, nil}
		}

		// advance normally
		res, err := itp.advance()
		g.Assert(res == nil)
		if err != nil {
			es := newErrorStruct(err, itp.frameStack.stackTrace())
			return &response{nil, -1, es}
		}
	}

}

// advance the interpreter forwards by one opcode.
func (itp *Interpreter) advance() (g.Value, g.Error) {

	f := itp.frameStack.peek()

	//if debugInterpreter {
	//	fmt.Printf("=========================================\n")
	//	fmt.Printf("ip: %s\n", bc.FmtBytecode(f.btc, f.ip))
	//	fmt.Printf("\n")
	//	itp.frameStack.debug()
	//}

	op := ops[f.btc[f.ip]]
	return op(itp, f)
}

//--------------------------------------------------------------
// debugging
//--------------------------------------------------------------

//var debugInterpreter bool
//
//func debugString(s string) {
//	if debugInterpreter {
//		fmt.Printf(s)
//	}
//}
//
//func debugVal(val g.Value) string {
//	if val == nil {
//		return "nil"
//	}
//	s, err := val.ToStr(nil)
//	if err != nil {
//		panic(err)
//	}
//	return s.String()
//}
//
//func (r response) debug() {
//	if debugInterpreter {
//		fmt.Printf("response(%s, %d, %s)\n",
//			debugVal(r.result),
//			r.resultIP,
//			debugVal(r.es))
//	}
//}
//
//func (f *frame) debug() {
//
//	fmt.Printf("    locals:\n")
//	for i, r := range f.locals {
//		fmt.Printf("        %d: %s\n", i, debugVal(r.Val))
//	}
//
//	fmt.Printf("    stack:\n")
//	for i, v := range f.stack {
//		fmt.Printf("        %d: %s\n", i, debugVal(v))
//	}
//
//	fmt.Printf("    handlers:\n")
//	for i, v := range f.handlers {
//		fmt.Printf("        %d: %s\n", i, v)
//	}
//
//	fmt.Printf("    ip: %d\n", f.ip)
//	fmt.Printf("    isBase: %v\n", f.isBase)
//	fmt.Printf("    isHandlingError: %v\n", f.isHandlingError)
//}
//
//func (fs *frameStack) debug() {
//
//	for i := fs.num() - 1; i >= 0; i-- {
//		fmt.Printf("-----------------------------------------\n")
//		fmt.Printf("frame %d\n", i)
//		fs.get(i).debug()
//		fmt.Printf("\n")
//	}
//}
