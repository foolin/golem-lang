// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/core/bytecode"
)

// Advance the interpreter forwards by one opcode.
func (i *Interpreter) advance(lastFrame int) (g.Value, g.Error) {

	frameIndex := len(i.frames) - 1
	f := i.frames[frameIndex]

	n := len(f.stack) - 1 // top of stack
	btc := f.fn.Template().Bytecodes
	pool := f.fn.Template().Module.Pool

	//dumpFrames(i.frames)

	switch btc[f.ip] {

	case bytecode.Invoke:

		p := bytecode.DecodeParam(btc, f.ip)
		params := f.stack[n-p+1:]

		switch fn := f.stack[n-p].(type) {
		case bytecode.BytecodeFunc:

			// check arity
			numParams := len(params)
			arity := fn.Template().Arity
			switch arity.Kind {
			case g.FixedArity:
				if uint16(numParams) != arity.RequiredParams {
					err := g.ArityMismatchError(
						fmt.Sprintf("%d", arity.RequiredParams), numParams)
					return nil, err
				}
			default:
				panic("unreachable")
			}

			// pop from stack
			f.stack = f.stack[:n-p]

			// push a new frame
			locals := newLocals(fn.Template().NumLocals, params)
			i.frames = append(i.frames, &frame{fn, locals, []g.Value{}, 0})

		case g.NativeFunc:

			val, err := fn.Invoke(i, params)
			if err != nil {
				return nil, err
			}

			f.stack = f.stack[:n-p]
			f.stack = append(f.stack, val)
			f.ip += 3

		default:
			return nil, g.TypeMismatchError("Expected Func")
		}

	case bytecode.Return:

		// TODO once we've written a Control Flow Graph
		// turn this sanity check on to make sure we are managing
		// the stack properly

		//if len(f.stack) < 1 || len(f.stack) > 2 {
		//	for j, v := range f.stack {
		//		fmt.Printf("stack %d: %s\n", j, v.ToStr())
		//	}
		//	panic("invalid stack")
		//}

		// get result from top of stack
		result := f.stack[n]

		// pop the old frame
		i.frames = i.frames[:frameIndex]

		// If we are on the last frame, then we are done.
		if frameIndex == lastFrame {
			return result, nil
		}

		// push the result onto the new top frame
		f = i.frames[len(i.frames)-1]
		f.stack = append(f.stack, result)

		// advance the instruction pointer now that we are done invoking
		f.ip += 3

	case bytecode.Done:
		panic("Done cannot be executed directly")

	case bytecode.Go:

		p := bytecode.DecodeParam(btc, f.ip)
		params := f.stack[n-p+1:]

		switch fn := f.stack[n-p].(type) {
		case bytecode.BytecodeFunc:
			f.stack = f.stack[:n-p]
			f.ip += 3

			intp := NewInterpreter(i.builtInMgr, i.modules)
			go (func() {
				_, errStruct := intp.EvalBytecode(fn, params)
				if errStruct != nil {
					panic("TODO")
					//fmt.Printf("%v\n", errStruct.Error())
				}
			})()

		case g.NativeFunc:
			f.stack = f.stack[:n-p]
			f.ip += 3

			go (func() {
				_, err := fn.Invoke(i, params)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			})()

		default:
			return nil, g.TypeMismatchError("Expected Func")
		}

	case bytecode.Throw:

		// get str from stack
		s, ok := f.stack[n].(g.Str)
		if !ok {
			return nil, g.TypeMismatchError("Expected Str")
		}

		// throw an error
		return nil, g.NewError(s.String())

	case bytecode.NewFunc:

		// push a function
		p := bytecode.DecodeParam(btc, f.ip)
		tpl := pool.Templates[p]
		nf := bytecode.NewBytecodeFunc(tpl)
		f.stack = append(f.stack, nf)
		f.ip += 3

	case bytecode.FuncLocal:

		// get function from stack
		fn, ok := f.stack[n].(bytecode.BytecodeFunc)
		if !ok {
			return nil, g.TypeMismatchError("Expected BytecodeFunc")
		}

		// push a local onto the captures of the function
		p := bytecode.DecodeParam(btc, f.ip)
		fn.PushCapture(f.locals[p])
		f.ip += 3

	case bytecode.FuncCapture:

		// get function from stack
		fn, ok := f.stack[n].(bytecode.BytecodeFunc)
		if !ok {
			return nil, g.TypeMismatchError("Expected BytecodeFunc")
		}

		// push a capture onto the captures of the function
		p := bytecode.DecodeParam(btc, f.ip)
		fn.PushCapture(f.fn.GetCapture(p))
		f.ip += 3

	case bytecode.NewList:

		size := bytecode.DecodeParam(btc, f.ip)
		vals := make([]g.Value, size)
		copy(vals, f.stack[n-size+1:])

		f.stack = f.stack[:n-size+1]
		f.stack = append(f.stack, g.NewList(vals))
		f.ip += 3

		//	case bytecode.NewSet:
		//
		//		size := bytecode.DecodeParam(btc, f.ip)
		//		vals := make([]g.Value, size)
		//		copy(vals, f.stack[n-size+1:])
		//		f.stack = f.stack[:n-size+1]
		//
		//		set, err := g.NewSet(i, vals)
		//		if err != nil {
		//			return nil, err
		//		}
		//
		//		f.stack = append(f.stack, set)
		//		f.ip += 3
		//
		//	case bytecode.NewTuple:
		//
		//		size := bytecode.DecodeParam(btc, f.ip)
		//		vals := make([]g.Value, size)
		//		copy(vals, f.stack[n-size+1:])
		//
		//		f.stack = f.stack[:n-size+1]
		//		f.stack = append(f.stack, g.NewTuple(vals))
		//		f.ip += 3
		//
		//	case bytecode.CheckTuple:
		//
		//		// make sure the top of the stack is really a tuple
		//		tp, ok := f.stack[n].(g.Tuple)
		//		if !ok {
		//			return nil, g.TypeMismatchError("Expected Tuple")
		//		}
		//
		//		// and make sure its of the expected length
		//		expectedLen := bytecode.DecodeParam(btc, f.ip)
		//		tpLen := tp.Len(i)
		//		if expectedLen != int(tpLen.IntVal()) {
		//			return nil, g.InvalidArgumentError(
		//				fmt.Sprintf("Expected Tuple of length %d", expectedLen))
		//		}
		//
		//		// do not alter stack
		//		f.ip += 3

	case bytecode.CheckCast:

		// make sure the top of the stack is of the given type
		vtype := g.Type(bytecode.DecodeParam(btc, f.ip))
		v := f.stack[n]
		if v.Type() != vtype {
			return nil, g.TypeMismatchError(fmt.Sprintf("Expected %s", vtype.String()))
		}

		// do not alter stack
		f.ip += 3

		//	case bytecode.NewDict:
		//
		//		size := bytecode.DecodeParam(btc, f.ip)
		//		entries := make([]*g.HEntry, 0, size)
		//
		//		numVals := size * 2
		//		for j := n - numVals + 1; j <= n; j += 2 {
		//			entries = append(entries, &g.HEntry{Key: f.stack[j], Value: f.stack[j+1]})
		//		}
		//
		//		f.stack = f.stack[:n-numVals+1]
		//
		//		dict, err := g.NewDict(i, entries)
		//		if err != nil {
		//			return nil, err
		//		}
		//
		//		f.stack = append(f.stack, dict)
		//		f.ip += 3
		//

		//	case bytecode.DefineStruct:
		//
		//		defs := pool.StructDefs[bytecode.DecodeParam(btc, f.ip)]
		//		stc, err := g.DefineStruct(defs)
		//		if err != nil {
		//			return nil, err
		//		}
		//
		//		f.stack = append(f.stack, stc)
		//		f.ip += 3

	case bytecode.GetField:

		p := bytecode.DecodeParam(btc, f.ip)
		key, ok := pool.Constants[p].(g.Str)
		g.Assert(ok)

		result, err := f.stack[n].GetField(key.String(), i)
		if err != nil {
			return nil, err
		}

		f.stack[n] = result
		f.ip += 3

	case bytecode.InvokeField:

		p, q := bytecode.DecodeWideParams(btc, f.ip)
		fmt.Printf("InvokeField aaa %d %d\n", p, q)

		key, ok := pool.Constants[p].(g.Str)
		g.Assert(ok)
		fmt.Printf("InvokeField bbb %s\n", key)

		params := f.stack[n-q+1:]
		fmt.Printf("InvokeField ccc %v %d %d %d\n", params, n, q, n-q+1)

		val, err := f.stack[n].InvokeField(key.String(), i, params)
		if err != nil {
			return nil, err
		}
		fmt.Printf("InvokeField ddd %s\n", val)

		f.stack = f.stack[:n-q]
		f.stack = append(f.stack, val)
		f.ip += 5

		//	case bytecode.InitField, bytecode.SetField:
		//
		//		idx := bytecode.DecodeParam(btc, f.ip)
		//		key, ok := pool.Constants[idx].(g.Str)
		//		g.Assert(ok)
		//
		//		// get struct from stack
		//		stc, ok := f.stack[n-1].(g.Struct)
		//		if !ok {
		//			return nil, g.TypeMismatchError("Expected Struct")
		//		}
		//
		//		// get value from stack
		//		value := f.stack[n]
		//
		//		// init or set
		//		if btc[f.ip] == bytecode.InitField {
		//			err := stc.InitField(i, key, value)
		//			if err != nil {
		//				return nil, err
		//			}
		//		} else {
		//			err := stc.SetField(i, key, value)
		//			if err != nil {
		//				return nil, err
		//			}
		//		}
		//
		//		f.stack[n-1] = value
		//		f.stack = f.stack[:n]
		//		f.ip += 3
		//
		//	case bytecode.IncField:
		//
		//		idx := bytecode.DecodeParam(btc, f.ip)
		//		key, ok := pool.Constants[idx].(g.Str)
		//		g.Assert(ok)
		//
		//		// get struct from stack
		//		stc, ok := f.stack[n-1].(g.Struct)
		//		if !ok {
		//			return nil, g.TypeMismatchError("Expected Struct")
		//		}
		//
		//		// get value from stack
		//		value := f.stack[n]
		//
		//		before, err := stc.GetField(i, key)
		//		if err != nil {
		//			return nil, err
		//		}
		//
		//		after, err := plus(i, before, value)
		//		if err != nil {
		//			return nil, err
		//		}
		//
		//		err = stc.SetField(i, key, after)
		//		if err != nil {
		//			return nil, err
		//		}
		//
		//		f.stack[n-1] = before
		//		f.stack = f.stack[:n]
		//		f.ip += 3

	case bytecode.GetIndex:

		// get Indexable from stack
		gtb, ok := f.stack[n-1].(g.Indexable)
		if !ok {
			return nil, g.TypeMismatchError("Expected Indexable")
		}

		// get index from stack
		idx := f.stack[n]

		result, err := gtb.Get(i, idx)
		if err != nil {
			return nil, err
		}

		f.stack[n-1] = result
		f.stack = f.stack[:n]
		f.ip++

	case bytecode.SetIndex:

		// get IndexAssignable from stack
		ibl, ok := f.stack[n-2].(g.IndexAssignable)
		if !ok {
			return nil, g.TypeMismatchError("Expected Indexable")
		}

		// get index from stack
		idx := f.stack[n-1]

		// get value from stack
		val := f.stack[n]

		err := ibl.Set(i, idx, val)
		if err != nil {
			return nil, err
		}

		f.stack[n-2] = val
		f.stack = f.stack[:n-1]
		f.ip++

	case bytecode.IncIndex:

		// get IndexAssignable from stack
		ibl, ok := f.stack[n-2].(g.IndexAssignable)
		if !ok {
			return nil, g.TypeMismatchError("Expected Indexable")
		}

		// get index from stack
		idx := f.stack[n-1]

		// get value from stack
		val := f.stack[n]

		before, err := ibl.Get(i, idx)
		if err != nil {
			return nil, err
		}

		after, err := plus(i, before, val)
		if err != nil {
			return nil, err
		}

		err = ibl.Set(i, idx, after)
		if err != nil {
			return nil, err
		}

		f.stack[n-2] = before
		f.stack = f.stack[:n-1]
		f.ip++

	case bytecode.Slice:

		// get Sliceable from stack
		slb, ok := f.stack[n-2].(g.Sliceable)
		if !ok {
			return nil, g.TypeMismatchError("Expected Sliceable")
		}

		// get indices from stack
		from := f.stack[n-1]
		to := f.stack[n]

		result, err := slb.Slice(i, from, to)
		if err != nil {
			return nil, err
		}

		f.stack[n-2] = result
		f.stack = f.stack[:n-1]
		f.ip++

	case bytecode.SliceFrom:

		// get Sliceable from stack
		slb, ok := f.stack[n-1].(g.Sliceable)
		if !ok {
			return nil, g.TypeMismatchError("Expected Sliceable")
		}

		// get index from stack
		from := f.stack[n]

		result, err := slb.SliceFrom(i, from)
		if err != nil {
			return nil, err
		}

		f.stack[n-1] = result
		f.stack = f.stack[:n]
		f.ip++

	case bytecode.SliceTo:

		// get Sliceable from stack
		slb, ok := f.stack[n-1].(g.Sliceable)
		if !ok {
			return nil, g.TypeMismatchError("Expected Sliceable")
		}

		// get index from stack
		to := f.stack[n]

		result, err := slb.SliceTo(i, to)
		if err != nil {
			return nil, err
		}

		f.stack[n-1] = result
		f.stack = f.stack[:n]
		f.ip++

	case bytecode.LoadNull:
		f.stack = append(f.stack, g.Null)
		f.ip++
	case bytecode.LoadTrue:
		f.stack = append(f.stack, g.True)
		f.ip++
	case bytecode.LoadFalse:
		f.stack = append(f.stack, g.False)
		f.ip++
	case bytecode.LoadZero:
		f.stack = append(f.stack, g.Zero)
		f.ip++
	case bytecode.LoadOne:
		f.stack = append(f.stack, g.One)
		f.ip++
	case bytecode.LoadNegOne:
		f.stack = append(f.stack, g.NegOne)
		f.ip++

	case bytecode.ImportModule:

		// get the module name from the pool
		p := bytecode.DecodeParam(btc, f.ip)
		name, ok := pool.Constants[p].(g.Str)
		g.Assert(ok)

		// Lookup the module.
		mod, err := i.lookupModule(name.String())
		if err != nil {
			return nil, err
		}

		// Push the module's contents onto the stack
		f.stack = append(f.stack, mod.Contents)
		f.ip += 3

	case bytecode.LoadBuiltin:
		p := bytecode.DecodeParam(btc, f.ip)
		f.stack = append(f.stack, i.builtInMgr.Builtins()[p])
		f.ip += 3

	case bytecode.LoadConst:
		p := bytecode.DecodeParam(btc, f.ip)
		f.stack = append(f.stack, pool.Constants[p])
		f.ip += 3

	case bytecode.LoadLocal:
		p := bytecode.DecodeParam(btc, f.ip)
		f.stack = append(f.stack, f.locals[p].Val)
		f.ip += 3

	case bytecode.LoadCapture:
		p := bytecode.DecodeParam(btc, f.ip)
		f.stack = append(f.stack, f.fn.GetCapture(p).Val)
		f.ip += 3

	case bytecode.StoreLocal:
		p := bytecode.DecodeParam(btc, f.ip)
		f.locals[p].Val = f.stack[n]
		f.stack = f.stack[:n]
		f.ip += 3

	case bytecode.StoreCapture:
		p := bytecode.DecodeParam(btc, f.ip)
		f.fn.GetCapture(p).Val = f.stack[n]
		f.stack = f.stack[:n]
		f.ip += 3

	case bytecode.Jump:
		f.ip = bytecode.DecodeParam(btc, f.ip)

	case bytecode.JumpTrue:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return nil, g.TypeMismatchError("Expected Bool")
		}

		f.stack = f.stack[:n]
		if b.BoolVal() {
			f.ip = bytecode.DecodeParam(btc, f.ip)
		} else {
			f.ip += 3
		}

	case bytecode.JumpFalse:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return nil, g.TypeMismatchError("Expected Bool")
		}

		f.stack = f.stack[:n]
		if b.BoolVal() {
			f.ip += 3
		} else {
			f.ip = bytecode.DecodeParam(btc, f.ip)
		}

	case bytecode.Eq:
		b, err := f.stack[n-1].Eq(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = b
		f.ip++

	case bytecode.Ne:
		b, err := f.stack[n-1].Eq(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = b.Not()
		f.ip++

	case bytecode.Lt:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.NewBool(val.IntVal() < 0)
		f.ip++

	case bytecode.Lte:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.NewBool(val.IntVal() <= 0)
		f.ip++

	case bytecode.Gt:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.NewBool(val.IntVal() > 0)
		f.ip++

	case bytecode.Gte:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.NewBool(val.IntVal() >= 0)
		f.ip++

	case bytecode.Cmp:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bytecode.Has:

		panic("TODO")

		//		// get struct from stack
		//		stc, ok := f.stack[n-1].(g.Struct)
		//		if !ok {
		//			return nil, g.TypeMismatchError("Expected Struct")
		//		}
		//
		//		val, err := stc.Has(i, f.stack[n])
		//		if err != nil {
		//			return nil, err
		//		}
		//		f.stack = f.stack[:n]
		//		f.stack[n-1] = val
		//		f.ip++

	case bytecode.Plus:
		val, err := plus(i, f.stack[n-1], f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bytecode.Not:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return nil, g.TypeMismatchError("Expected Bool")
		}

		f.stack[n] = b.Not()
		f.ip++

	case bytecode.Sub:
		z, ok := f.stack[n-1].(g.Number)
		if !ok {
			return nil, g.TypeMismatchError("Expected Number Type")
		}

		val, err := z.Sub(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bytecode.Mul:
		z, ok := f.stack[n-1].(g.Number)
		if !ok {
			return nil, g.TypeMismatchError("Expected Number Type")
		}

		val, err := z.Mul(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bytecode.Div:
		z, ok := f.stack[n-1].(g.Number)
		if !ok {
			return nil, g.TypeMismatchError("Expected Number Type")
		}

		val, err := z.Div(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bytecode.Negate:
		z, ok := f.stack[n].(g.Number)
		if !ok {
			return nil, g.TypeMismatchError("Expected Number Type")
		}

		val := z.Negate()
		f.stack[n] = val
		f.ip++

	case bytecode.Rem:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected Int")
		}

		val, err := z.Rem(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bytecode.BitAnd:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected Int")
		}

		val, err := z.BitAnd(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bytecode.BitOr:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected Int")
		}

		val, err := z.BitOr(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bytecode.BitXor:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected Int")
		}

		val, err := z.BitXOr(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bytecode.LeftShift:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected Int")
		}

		val, err := z.LeftShift(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bytecode.RightShift:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected Int")
		}

		val, err := z.RightShift(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bytecode.Complement:
		z, ok := f.stack[n].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected Int")
		}

		val := z.Complement()
		f.stack[n] = val
		f.ip++

	case bytecode.NewIter:

		ibl, ok := f.stack[n].(g.Iterable)
		g.Assert(ok)

		f.stack[n] = ibl.NewIterator(i)
		f.ip++

	case bytecode.IterNext:

		itr, ok := f.stack[n].(g.Iterator)
		g.Assert(ok)

		val, err := itr.IterNext(i)
		if err != nil {
			return nil, err
		}

		f.stack[n] = val
		f.ip++

	case bytecode.IterGet:

		itr, ok := f.stack[n].(g.Iterator)
		g.Assert(ok)

		val, err := itr.IterGet(i)
		if err != nil {
			return nil, err
		}

		f.stack[n] = val
		f.ip++

	case bytecode.Dup:
		f.stack = append(f.stack, f.stack[n])
		f.ip++

	case bytecode.Pop:
		f.stack = f.stack[:n]
		f.ip++

	default:
		panic("Invalid opcode")
	}

	return nil, nil
}

func plus(ev g.Evaluator, a g.Value, b g.Value) (g.Value, g.Error) {

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
	if ia && ib {
		return na.Add(nb)
	}

	// error
	return nil, g.TypeMismatchError("Expected Number Type")
}

//func index(opcodes []byte, ip int) int {
//
//	high := opcodes[ip+1]
//	low := opcodes[ip+2]
//
//	return int(high)<<8 + int(low)
//}
//
//func wideIndex(opcodes []byte, ip int) (int, int) {
//
//	high1 := opcodes[ip+1]
//	low1 := opcodes[ip+2]
//
//	high2 := opcodes[ip+3]
//	low2 := opcodes[ip+4]
//
//	return int(high1)<<8 + int(low1), int(high2)<<8 + int(low2)
//}
