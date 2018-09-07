// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
)

type op func(*Interpreter, *frame) (g.Value, g.Error)

var ops []op

func init() {
	ops = []op{
		opLoadNull,
		opLoadTrue,
		opLoadFalse,
		opLoadZero,
		opLoadOne,
		opLoadNegOne,

		opImportModule,
		opLoadBuiltin,
		opLoadConst,
		opLoadLocal,
		opStoreLocal,
		opLoadCapture,
		opStoreCapture,

		opJump,
		opJumpTrue,
		opJumpFalse,

		opEq,
		opNe,
		opGt,
		opGte,
		opLt,
		opLte,
		opCmp,

		opPlus,
		opInc,
		opSub,
		opMul,
		opDiv,

		opRem,
		opBitAnd,
		opBitOr,
		opBitXor,
		opLeftShift,
		opRightShift,

		opNegate,
		opNot,
		opComplement,

		opPop,
		opDup,

		opNewFunc,
		opFuncCapture,
		opFuncLocal,

		opInvoke,
		opGo,
		opReturn,

		opPushTry,
		opPopTry,
		opThrow,

		opNewStruct,
		opNewDict,
		opNewList,
		opNewSet,
		opNewTuple,
		opCheckTuple,

		opGetField,
		opInvokeField,
		opInitField,
		opInitProperty,
		opInitReadonlyProperty,
		opSetField,
		opIncField,

		opGetIndex,
		opSetIndex,
		opIncIndex,

		opSlice,
		opSliceFrom,
		opSliceTo,

		opNewIter,
		opIterNext,
		opIterGet,
	}
}

// advance the interpreter forwards by one opcode.
func (itp *Interpreter) advance() (g.Value, g.Error) {
	f := itp.peekFrame()
	op := ops[f.btc[f.ip]]
	return op(itp, f)
}

//--------------------------------------------------------------

func opInvoke(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	p := bc.DecodeParam(f.btc, f.ip)
	params := f.stack[n-p+1:]

	switch fn := f.stack[n-p].(type) {
	case bc.Func:

		// check arity, and modify params if necessary
		arity := fn.Template().Arity
		numParams := len(params)
		numReq := int(arity.Required)

		switch arity.Kind {

		case g.FixedArity:

			if numParams != numReq {
				return nil, g.ArityMismatch(numReq, numParams)
			}

		case g.VariadicArity:

			if numParams < numReq {
				return nil, g.ArityMismatchAtLeast(numReq, numParams)
			}

			// turn any extra params into a list
			req := params[:numReq]
			vrd := params[numReq:]
			params = append(req, g.NewList(g.CopyValues(vrd)))

		case g.MultipleArity:

			if numParams < numReq {
				return nil, g.ArityMismatchAtLeast(numReq, numParams)
			}

			numOpt := int(arity.Optional)
			if numParams > (numReq + numOpt) {
				return nil, g.ArityMismatchAtMost(numReq+numOpt, numParams)
			}

			// add any missing optional params
			if numParams < (numReq + numOpt) {
				opt := fn.Template().OptionalParams
				missing := g.CopyValues(opt[numParams-numReq:])
				params = append(params, missing...)
			}

		default:
			panic("unreachable")
		}

		// Pop from stack. Note this is the only opcode that pops
		// the stack before doing the actual operation.  Also note that
		// we do not actually advance the instruction pointer here.
		// That is done in bc.Return.
		f.stack = f.stack[:n-p]

		// push a new frame
		locals := newLocals(fn.Template().NumLocals, params)
		itp.pushFrame(newFrame(fn, locals, false))

	case g.NativeFunc:

		val, err := fn.Invoke(itp, params)
		if err != nil {
			return nil, err
		}

		// pop from stack
		f.stack = f.stack[:n-p]

		// push result
		f.stack = append(f.stack, val)

		f.ip += 3

	default:
		return nil, g.TypeMismatch(g.FuncType, f.stack[n-p].Type())
	}

	return nil, nil
}

func opReturn(itp *Interpreter, f *frame) (g.Value, g.Error) {

	//// TODO once we've written a Control Flow Graph
	//// turn this sanity check on to make sure we are managing
	//// the stack properly
	//if len(f.stack) < 1 || len(f.stack) > 2 {
	//	dumpFrame(f)
	//	panic("invalid stack")
	//}

	// get result from top of stack
	n := len(f.stack) - 1
	result := f.stack[n]

	// if we are handling an error then 'return' has different semantics.
	if f.isHandlingError {
		// If we are on the last frame, then we are done.
		if f.isLast {
			return result, nil
		}
		return nil, nil
	}

	// discard the frame
	itp.popFrame()

	// If we are on the last frame, then we are done.
	if f.isLast {
		return result, nil
	}

	// push the result onto the new top frame
	f = itp.peekFrame()
	f.stack = append(f.stack, result)

	// advance the instruction pointer now that we are done invoking
	f.ip += 3

	return nil, nil
}

func opPushTry(itp *Interpreter, f *frame) (g.Value, g.Error) {

	p := bc.DecodeParam(f.btc, f.ip)
	f.pushHandler(f.handlerPool[p])
	f.ip += 3

	return nil, nil
}

func opPopTry(itp *Interpreter, f *frame) (g.Value, g.Error) {

	f.popHandler()
	f.ip++

	return nil, nil
}

func opThrow(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	s, ok := f.stack[n].(g.Str)
	if !ok {
		return nil, g.TypeMismatch(g.StrType, f.stack[n].Type())
	}

	return nil, g.Error(fmt.Errorf("%s", s.String()))
}

func opGo(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	p := bc.DecodeParam(f.btc, f.ip)
	params := f.stack[n-p+1:]

	switch fn := f.stack[n-p].(type) {
	case bc.Func:
		f.stack = f.stack[:n-p]
		f.ip += 3

		intp := NewInterpreter(itp.builtInMgr, itp.modules)
		go (func() {
			_, es := intp.EvalBytecode(fn, params)
			if es != nil {
				panic("TODO how to handle exceptions in goroutines")
			}
		})()

	case g.NativeFunc:
		f.stack = f.stack[:n-p]
		f.ip += 3

		go (func() {
			_, err := fn.Invoke(itp, params)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		})()

	default:
		return nil, g.TypeMismatch(g.FuncType, f.stack[n-p].Type())
	}

	return nil, nil
}

func opNewFunc(itp *Interpreter, f *frame) (g.Value, g.Error) {

	// push a function
	p := bc.DecodeParam(f.btc, f.ip)
	tpl := f.pool.Templates[p]
	nf := bc.NewBytecodeFunc(tpl)
	f.stack = append(f.stack, nf)
	f.ip += 3

	return nil, nil
}

func opFuncLocal(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	// get function from stack
	fn := f.stack[n].(bc.Func)

	// push a local onto the captures of the function
	p := bc.DecodeParam(f.btc, f.ip)
	fn.PushCapture(f.locals[p])
	f.ip += 3

	return nil, nil
}

func opFuncCapture(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	// get function from stack
	fn := f.stack[n].(bc.Func)

	// push a capture onto the captures of the function
	p := bc.DecodeParam(f.btc, f.ip)
	fn.PushCapture(f.fn.GetCapture(p))
	f.ip += 3

	return nil, nil
}

func opNewList(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	size := bc.DecodeParam(f.btc, f.ip)
	ns := n - size + 1
	vals := g.CopyValues(f.stack[ns:])

	f.stack = f.stack[:ns]
	f.stack = append(f.stack, g.NewList(vals))
	f.ip += 3

	return nil, nil
}

func opNewSet(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	size := bc.DecodeParam(f.btc, f.ip)
	ns := n - size + 1
	vals := g.CopyValues(f.stack[ns:])

	set, err := g.NewSet(itp, vals)
	if err != nil {
		return nil, err
	}

	f.stack = f.stack[:ns]
	f.stack = append(f.stack, set)
	f.ip += 3

	return nil, nil
}

func opNewTuple(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	size := bc.DecodeParam(f.btc, f.ip)
	ns := n - size + 1
	vals := g.CopyValues(f.stack[ns:])

	f.stack = f.stack[:ns]
	f.stack = append(f.stack, g.NewTuple(vals))
	f.ip += 3

	return nil, nil
}

func opCheckTuple(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	// make sure the top of the stack is really a tuple
	tp, ok := f.stack[n].(g.Tuple)
	if !ok {
		return nil, g.TypeMismatch(g.TupleType, f.stack[n].Type())
	}

	// and make sure its of the expected length
	expectedLen := bc.DecodeParam(f.btc, f.ip)
	tpLen, err := tp.Len(itp)
	if err != nil {
		return nil, err
	}
	if expectedLen != int(tpLen.IntVal()) {
		return nil, g.InvalidArgument(
			fmt.Sprintf(
				"Expected Tuple of length %d, not length %d",
				expectedLen, int(tpLen.IntVal())))
	}

	// do not alter stack
	f.ip += 3

	return nil, nil
}

func opNewDict(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	size := bc.DecodeParam(f.btc, f.ip)
	entries := make([]*g.HEntry, 0, size)

	numVals := size * 2
	for j := n - numVals + 1; j <= n; j += 2 {
		entries = append(entries, &g.HEntry{Key: f.stack[j], Value: f.stack[j+1]})
	}

	f.stack = f.stack[:n-numVals+1]

	hashmap, err := g.NewHashMap(itp, entries)
	if err != nil {
		return nil, err
	}

	f.stack = append(f.stack, g.NewDict(hashmap))
	f.ip += 3

	return nil, nil
}

func opNewStruct(itp *Interpreter, f *frame) (g.Value, g.Error) {

	p := bc.DecodeParam(f.btc, f.ip)
	def := f.pool.StructDefs[p]
	fields := make(map[string]g.Field)
	for _, name := range def {
		fields[name] = g.NewField(g.Null)
	}

	stc, err := g.NewFieldStruct(fields)
	if err != nil {
		return nil, err
	}

	f.stack = append(f.stack, stc)
	f.ip += 3

	return nil, nil
}

func opNewIter(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	ibl, ok := f.stack[n].(g.Iterable)
	g.Assert(ok)

	itr, err := ibl.NewIterator(itp)
	if err != nil {
		return nil, err
	}

	f.stack[n] = itr
	f.ip++

	return nil, nil
}

func opIterNext(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	itr, ok := f.stack[n].(g.Iterator)
	g.Assert(ok)

	val, err := itr.IterNext(itp)
	if err != nil {
		return nil, err
	}

	f.stack[n] = val
	f.ip++

	return nil, nil
}

func opIterGet(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	itr, ok := f.stack[n].(g.Iterator)
	g.Assert(ok)

	val, err := itr.IterGet(itp)
	if err != nil {
		return nil, err
	}

	f.stack[n] = val
	f.ip++

	return nil, nil
}

func opGetField(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	p := bc.DecodeParam(f.btc, f.ip)
	key, ok := f.pool.Constants[p].(g.Str)
	g.Assert(ok)

	result, err := f.stack[n].GetField(itp, key.String())
	if err != nil {
		return nil, err
	}

	f.stack[n] = result
	f.ip += 3

	return nil, nil
}

func opInvokeField(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	p, q := bc.DecodeWideParams(f.btc, f.ip)

	key, ok := f.pool.Constants[p].(g.Str)
	g.Assert(ok)

	self := f.stack[n-q]
	params := f.stack[n-q+1:]

	result, err := self.InvokeField(itp, key.String(), params)
	if err != nil {
		return nil, err
	}

	f.stack[n-q] = result
	f.stack = f.stack[:n-q+1]
	f.ip += 5

	return nil, nil
}

func opSetField(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	p := bc.DecodeParam(f.btc, f.ip)
	key, ok := f.pool.Constants[p].(g.Str)
	g.Assert(ok)

	if f.stack[n-1].Type() == g.NullType {
		return nil, g.NullValueError()
	}

	stc, ok := f.stack[n-1].(g.Struct)
	if !ok {
		return nil, g.TypeMismatch(g.StructType, f.stack[n-1].Type())
	}
	value := f.stack[n]

	err := stc.SetField(itp, key.String(), value)
	if err != nil {
		return nil, err
	}
	f.stack[n-1] = value
	f.stack = f.stack[:n]

	f.ip += 3

	return nil, nil
}

func opInitField(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	p := bc.DecodeParam(f.btc, f.ip)
	key, ok := f.pool.Constants[p].(g.Str)
	g.Assert(ok)

	stc := f.stack[n-1].(g.Struct)
	value := f.stack[n]

	stc.Internal(key.String(), g.NewField(value))
	f.stack = f.stack[:n]

	f.ip += 3

	return nil, nil
}

func opInitProperty(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	p := bc.DecodeParam(f.btc, f.ip)
	key, ok := f.pool.Constants[p].(g.Str)
	g.Assert(ok)

	stc := f.stack[n-2].(g.Struct)
	get := f.stack[n-1].(g.Func)
	set := f.stack[n].(g.Func)

	prop, err := g.NewProperty(get, set)
	g.Assert(err == nil)

	stc.Internal(key.String(), prop)
	f.stack = f.stack[:n-1]

	f.ip += 3

	return nil, nil
}

func opInitReadonlyProperty(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	p := bc.DecodeParam(f.btc, f.ip)
	key, ok := f.pool.Constants[p].(g.Str)
	g.Assert(ok)

	stc := f.stack[n-1].(g.Struct)
	get := f.stack[n].(g.Func)

	prop, err := g.NewReadonlyProperty(get)
	g.Assert(err == nil)

	stc.Internal(key.String(), prop)
	f.stack = f.stack[:n]

	f.ip += 3

	return nil, nil
}

func opIncField(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	p := bc.DecodeParam(f.btc, f.ip)
	key, ok := f.pool.Constants[p].(g.Str)
	g.Assert(ok)

	// get struct from stack
	stc, ok := f.stack[n-1].(g.Struct)
	if !ok {
		return nil, g.TypeMismatch(g.StructType, f.stack[n-1].Type())
	}

	// get value from stack
	value := f.stack[n]

	before, err := stc.GetField(itp, key.String())
	if err != nil {
		return nil, err
	}

	after, err := inc(itp, before, value)
	if err != nil {
		return nil, err
	}

	err = stc.SetField(itp, key.String(), after)
	if err != nil {
		return nil, err
	}

	f.stack[n-1] = before
	f.stack = f.stack[:n]
	f.ip += 3

	return nil, nil
}

func opGetIndex(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	// get Indexable from stack
	gtb, ok := f.stack[n-1].(g.Indexable)
	if !ok {
		return nil, g.IndexableMismatch(f.stack[n-1].Type())
	}

	// get index from stack
	idx := f.stack[n]

	result, err := gtb.Get(itp, idx)
	if err != nil {
		return nil, err
	}

	f.stack[n-1] = result
	f.stack = f.stack[:n]
	f.ip++

	return nil, nil
}

func opSetIndex(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	// get Indexable from stack
	ibl, ok := f.stack[n-2].(g.Indexable)
	if !ok {
		return nil, g.IndexableMismatch(f.stack[n-2].Type())
	}

	// get index from stack
	idx := f.stack[n-1]

	// get value from stack
	val := f.stack[n]

	err := ibl.Set(itp, idx, val)
	if err != nil {
		return nil, err
	}

	f.stack[n-2] = val
	f.stack = f.stack[:n-1]
	f.ip++

	return nil, nil
}

func opIncIndex(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	// get Indexable from stack
	ibl, ok := f.stack[n-2].(g.Indexable)
	if !ok {
		return nil, g.IndexableMismatch(f.stack[n-2].Type())
	}

	// get index from stack
	idx := f.stack[n-1]

	// get value from stack
	val := f.stack[n]

	before, err := ibl.Get(itp, idx)
	if err != nil {
		return nil, err
	}

	after, err := inc(itp, before, val)
	if err != nil {
		return nil, err
	}

	err = ibl.Set(itp, idx, after)
	if err != nil {
		return nil, err
	}

	f.stack[n-2] = before
	f.stack = f.stack[:n-1]
	f.ip++

	return nil, nil
}

func opSlice(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	// get Sliceable from stack
	slb, ok := f.stack[n-2].(g.Sliceable)
	if !ok {
		return nil, g.SliceableMismatch(f.stack[n-2].Type())
	}

	// get indices from stack
	from := f.stack[n-1]
	to := f.stack[n]

	result, err := slb.Slice(itp, from, to)
	if err != nil {
		return nil, err
	}

	f.stack[n-2] = result
	f.stack = f.stack[:n-1]
	f.ip++

	return nil, nil
}

func opSliceFrom(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	// get Sliceable from stack
	slb, ok := f.stack[n-1].(g.Sliceable)
	if !ok {
		return nil, g.SliceableMismatch(f.stack[n-1].Type())
	}

	// get index from stack
	from := f.stack[n]

	result, err := slb.SliceFrom(itp, from)
	if err != nil {
		return nil, err
	}

	f.stack[n-1] = result
	f.stack = f.stack[:n]
	f.ip++

	return nil, nil
}

func opSliceTo(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	// get Sliceable from stack
	slb, ok := f.stack[n-1].(g.Sliceable)
	if !ok {
		return nil, g.SliceableMismatch(f.stack[n-1].Type())
	}

	// get index from stack
	to := f.stack[n]

	result, err := slb.SliceTo(itp, to)
	if err != nil {
		return nil, err
	}

	f.stack[n-1] = result
	f.stack = f.stack[:n]
	f.ip++

	return nil, nil
}

func opLoadNull(itp *Interpreter, f *frame) (g.Value, g.Error) {

	f.stack = append(f.stack, g.Null)
	f.ip++

	return nil, nil
}

func opLoadTrue(itp *Interpreter, f *frame) (g.Value, g.Error) {

	f.stack = append(f.stack, g.True)
	f.ip++

	return nil, nil
}

func opLoadFalse(itp *Interpreter, f *frame) (g.Value, g.Error) {

	f.stack = append(f.stack, g.False)
	f.ip++

	return nil, nil
}

func opLoadZero(itp *Interpreter, f *frame) (g.Value, g.Error) {

	f.stack = append(f.stack, g.Zero)
	f.ip++

	return nil, nil
}

func opLoadOne(itp *Interpreter, f *frame) (g.Value, g.Error) {

	f.stack = append(f.stack, g.One)
	f.ip++

	return nil, nil
}

func opLoadNegOne(itp *Interpreter, f *frame) (g.Value, g.Error) {

	f.stack = append(f.stack, g.NegOne)
	f.ip++

	return nil, nil
}

func opImportModule(itp *Interpreter, f *frame) (g.Value, g.Error) {

	// get the module name from the f.pool
	p := bc.DecodeParam(f.btc, f.ip)
	name, ok := f.pool.Constants[p].(g.Str)
	g.Assert(ok)

	// Lookup the module.
	mod, err := itp.lookupModule(name.String())
	if err != nil {
		return nil, err
	}

	// Push the module's contents onto the stack
	f.stack = append(f.stack, mod.Contents)
	f.ip += 3

	return nil, nil
}

func opLoadBuiltin(itp *Interpreter, f *frame) (g.Value, g.Error) {

	p := bc.DecodeParam(f.btc, f.ip)
	f.stack = append(f.stack, itp.builtInMgr.Builtins()[p])
	f.ip += 3

	return nil, nil
}

func opLoadConst(itp *Interpreter, f *frame) (g.Value, g.Error) {

	p := bc.DecodeParam(f.btc, f.ip)
	f.stack = append(f.stack, f.pool.Constants[p])
	f.ip += 3

	return nil, nil
}

func opLoadLocal(itp *Interpreter, f *frame) (g.Value, g.Error) {

	p := bc.DecodeParam(f.btc, f.ip)
	f.stack = append(f.stack, f.locals[p].Val)
	f.ip += 3

	return nil, nil
}

func opLoadCapture(itp *Interpreter, f *frame) (g.Value, g.Error) {

	p := bc.DecodeParam(f.btc, f.ip)
	f.stack = append(f.stack, f.fn.GetCapture(p).Val)
	f.ip += 3

	return nil, nil
}

func opStoreLocal(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	p := bc.DecodeParam(f.btc, f.ip)
	f.locals[p].Val = f.stack[n]
	f.stack = f.stack[:n]
	f.ip += 3

	return nil, nil
}

func opStoreCapture(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	p := bc.DecodeParam(f.btc, f.ip)
	f.fn.GetCapture(p).Val = f.stack[n]
	f.stack = f.stack[:n]
	f.ip += 3

	return nil, nil
}

func opJump(itp *Interpreter, f *frame) (g.Value, g.Error) {

	f.ip = bc.DecodeParam(f.btc, f.ip)

	return nil, nil
}

func opJumpTrue(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	b, ok := f.stack[n].(g.Bool)
	if !ok {
		return nil, g.TypeMismatch(g.BoolType, f.stack[n].Type())
	}

	f.stack = f.stack[:n]
	if b.BoolVal() {
		f.ip = bc.DecodeParam(f.btc, f.ip)
	} else {
		f.ip += 3
	}

	return nil, nil
}

func opJumpFalse(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	b, ok := f.stack[n].(g.Bool)
	if !ok {
		return nil, g.TypeMismatch(g.BoolType, f.stack[n].Type())
	}

	f.stack = f.stack[:n]
	if b.BoolVal() {
		f.ip += 3
	} else {
		f.ip = bc.DecodeParam(f.btc, f.ip)
	}

	return nil, nil
}

func opEq(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	b, err := f.stack[n-1].Eq(itp, f.stack[n])
	if err != nil {
		return nil, err
	}
	f.stack = f.stack[:n]
	f.stack[n-1] = b
	f.ip++

	return nil, nil
}

func opNe(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	b, err := f.stack[n-1].Eq(itp, f.stack[n])
	if err != nil {
		return nil, err
	}
	f.stack = f.stack[:n]
	f.stack[n-1] = b.Not()
	f.ip++

	return nil, nil
}

func opLt(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Comparable)
	rhs, rhsOk := f.stack[n].(g.Comparable)
	if !lhsOk || !rhsOk {
		return nil, g.ComparableMismatch(f.stack[n-1].Type(), f.stack[n].Type())
	}

	val, err := lhs.Cmp(itp, rhs)
	if err != nil {
		return nil, err
	}
	f.stack = f.stack[:n]
	f.stack[n-1] = g.NewBool(val.IntVal() < 0)
	f.ip++

	return nil, nil
}

func opLte(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Comparable)
	rhs, rhsOk := f.stack[n].(g.Comparable)
	if !lhsOk || !rhsOk {
		return nil, g.ComparableMismatch(f.stack[n-1].Type(), f.stack[n].Type())
	}

	val, err := lhs.Cmp(itp, rhs)
	if err != nil {
		return nil, err
	}
	f.stack = f.stack[:n]
	f.stack[n-1] = g.NewBool(val.IntVal() <= 0)
	f.ip++

	return nil, nil
}

func opGt(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Comparable)
	rhs, rhsOk := f.stack[n].(g.Comparable)
	if !lhsOk || !rhsOk {
		return nil, g.ComparableMismatch(f.stack[n-1].Type(), f.stack[n].Type())
	}

	val, err := lhs.Cmp(itp, rhs)
	if err != nil {
		return nil, err
	}
	f.stack = f.stack[:n]
	f.stack[n-1] = g.NewBool(val.IntVal() > 0)
	f.ip++

	return nil, nil
}

func opGte(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Comparable)
	rhs, rhsOk := f.stack[n].(g.Comparable)
	if !lhsOk || !rhsOk {
		return nil, g.ComparableMismatch(f.stack[n-1].Type(), f.stack[n].Type())
	}

	val, err := lhs.Cmp(itp, rhs)
	if err != nil {
		return nil, err
	}
	f.stack = f.stack[:n]
	f.stack[n-1] = g.NewBool(val.IntVal() >= 0)
	f.ip++

	return nil, nil
}

func opCmp(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Comparable)
	rhs, rhsOk := f.stack[n].(g.Comparable)
	if !lhsOk || !rhsOk {
		return nil, g.ComparableMismatch(f.stack[n-1].Type(), f.stack[n].Type())
	}

	val, err := lhs.Cmp(itp, rhs)
	if err != nil {
		return nil, err
	}
	f.stack = f.stack[:n]
	f.stack[n-1] = val
	f.ip++

	return nil, nil
}

func opPlus(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	val, err := plus(itp, f.stack[n-1], f.stack[n])
	if err != nil {
		return nil, err
	}
	f.stack = f.stack[:n]
	f.stack[n-1] = val
	f.ip++

	return nil, nil
}

func opInc(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	val, err := inc(itp, f.stack[n-1], f.stack[n])
	if err != nil {
		return nil, err
	}
	f.stack = f.stack[:n]
	f.stack[n-1] = val
	f.ip++

	return nil, nil
}

func opNot(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	b, ok := f.stack[n].(g.Bool)
	if !ok {
		return nil, g.TypeMismatch(g.BoolType, f.stack[n].Type())
	}

	f.stack[n] = b.Not()
	f.ip++

	return nil, nil
}

func opSub(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Number)
	rhs, rhsOk := f.stack[n].(g.Number)
	if !lhsOk {
		return nil, g.NumberMismatch(f.stack[n-1].Type())
	}
	if !rhsOk {
		return nil, g.NumberMismatch(f.stack[n].Type())
	}

	val := lhs.Sub(rhs)

	f.stack = f.stack[:n]
	f.stack[n-1] = val
	f.ip++

	return nil, nil
}

func opMul(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Number)
	rhs, rhsOk := f.stack[n].(g.Number)
	if !lhsOk {
		return nil, g.NumberMismatch(f.stack[n-1].Type())
	}
	if !rhsOk {
		return nil, g.NumberMismatch(f.stack[n].Type())
	}

	val := lhs.Mul(rhs)

	f.stack = f.stack[:n]
	f.stack[n-1] = val
	f.ip++

	return nil, nil
}

func opDiv(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Number)
	rhs, rhsOk := f.stack[n].(g.Number)
	if !lhsOk {
		return nil, g.NumberMismatch(f.stack[n-1].Type())
	}
	if !rhsOk {
		return nil, g.NumberMismatch(f.stack[n].Type())
	}

	val, err := lhs.Div(rhs)
	if err != nil {
		return nil, err
	}

	f.stack = f.stack[:n]
	f.stack[n-1] = val
	f.ip++

	return nil, nil
}

func opNegate(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n].(g.Number)
	if !lhsOk {
		return nil, g.NumberMismatch(f.stack[n].Type())
	}

	val := lhs.Negate()
	f.stack[n] = val
	f.ip++

	return nil, nil
}

func opRem(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Int)
	rhs, rhsOk := f.stack[n].(g.Int)
	if !lhsOk {
		return nil, g.TypeMismatch(g.IntType, f.stack[n-1].Type())
	}
	if !rhsOk {
		return nil, g.TypeMismatch(g.IntType, f.stack[n].Type())
	}

	val := lhs.Rem(rhs)

	f.stack = f.stack[:n]
	f.stack[n-1] = val
	f.ip++

	return nil, nil
}

func opBitAnd(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Int)
	rhs, rhsOk := f.stack[n].(g.Int)
	if !lhsOk {
		return nil, g.TypeMismatch(g.IntType, f.stack[n-1].Type())
	}
	if !rhsOk {
		return nil, g.TypeMismatch(g.IntType, f.stack[n].Type())
	}

	val := lhs.BitAnd(rhs)

	f.stack = f.stack[:n]
	f.stack[n-1] = val
	f.ip++

	return nil, nil
}

func opBitOr(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Int)
	rhs, rhsOk := f.stack[n].(g.Int)
	if !lhsOk {
		return nil, g.TypeMismatch(g.IntType, f.stack[n-1].Type())
	}
	if !rhsOk {
		return nil, g.TypeMismatch(g.IntType, f.stack[n].Type())
	}

	val := lhs.BitOr(rhs)

	f.stack = f.stack[:n]
	f.stack[n-1] = val
	f.ip++

	return nil, nil
}

func opBitXor(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Int)
	rhs, rhsOk := f.stack[n].(g.Int)
	if !lhsOk {
		return nil, g.TypeMismatch(g.IntType, f.stack[n-1].Type())
	}
	if !rhsOk {
		return nil, g.TypeMismatch(g.IntType, f.stack[n].Type())
	}

	val := lhs.BitXOr(rhs)

	f.stack = f.stack[:n]
	f.stack[n-1] = val
	f.ip++

	return nil, nil
}

func opLeftShift(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Int)
	rhs, rhsOk := f.stack[n].(g.Int)
	if !lhsOk {
		return nil, g.TypeMismatch(g.IntType, f.stack[n-1].Type())
	}
	if !rhsOk {
		return nil, g.TypeMismatch(g.IntType, f.stack[n].Type())
	}

	val, err := lhs.LeftShift(rhs)
	if err != nil {
		return nil, err
	}

	f.stack = f.stack[:n]
	f.stack[n-1] = val
	f.ip++

	return nil, nil
}

func opRightShift(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, lhsOk := f.stack[n-1].(g.Int)
	rhs, rhsOk := f.stack[n].(g.Int)
	if !lhsOk {
		return nil, g.TypeMismatch(g.IntType, f.stack[n-1].Type())
	}
	if !rhsOk {
		return nil, g.TypeMismatch(g.IntType, f.stack[n].Type())
	}

	val, err := lhs.RightShift(rhs)
	if err != nil {
		return nil, err
	}

	f.stack = f.stack[:n]
	f.stack[n-1] = val
	f.ip++

	return nil, nil
}

func opComplement(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	lhs, ok := f.stack[n].(g.Int)
	if !ok {
		return nil, g.TypeMismatch(g.IntType, f.stack[n-1].Type())
	}

	val := lhs.Complement()
	f.stack[n] = val
	f.ip++

	return nil, nil
}

func opDup(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	f.stack = append(f.stack, f.stack[n])
	f.ip++

	return nil, nil
}

func opPop(itp *Interpreter, f *frame) (g.Value, g.Error) {

	n := len(f.stack) - 1

	f.stack = f.stack[:n]
	f.ip++

	return nil, nil
}

//--------------------------------------------------------------

func plus(ev g.Eval, a g.Value, b g.Value) (g.Value, g.Error) {

	// if either is a Str, return concatenated strings
	_, ia := a.(g.Str)
	_, ib := b.(g.Str)
	if ia || ib {

		as, err := a.ToStr(ev)
		if err != nil {
			return nil, err
		}

		bs, err := b.ToStr(ev)
		if err != nil {
			return nil, err
		}

		return as.Concat(bs), nil
	}

	// if both are Numbers, add them together
	na, ia := a.(g.Number)
	nb, ib := b.(g.Number)
	if !ia {
		return nil, g.NumberMismatch(a.Type())
	}
	if !ib {
		return nil, g.NumberMismatch(b.Type())
	}

	return na.Add(nb), nil
}

func inc(ev g.Eval, a g.Value, b g.Value) (g.Value, g.Error) {

	na, ok := a.(g.Number)
	if !ok {
		return nil, g.NumberMismatch(a.Type())
	}
	nb := b.(g.Number) // this cast always succeeds
	return na.Add(nb), nil
}
