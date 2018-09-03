// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"
	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
)

// Advance the interpreter forwards by one opcode.
func (itp *Interpreter) advance(lastFrame int) (g.Value, g.Error) {

	frameIndex := len(itp.frames) - 1
	f := itp.frames[frameIndex]

	n := len(f.stack) - 1 // top of stack
	btc := f.fn.Template().Bytecodes
	pool := f.fn.Template().Module.Pool

	//dumpFrames(itp.frames)

	switch btc[f.ip] {

	case bc.Invoke:

		p := bc.DecodeParam(btc, f.ip)
		params := f.stack[n-p+1:]

		switch fn := f.stack[n-p].(type) {
		case bc.BytecodeFunc:

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
			itp.frames = append(itp.frames, &frame{fn, locals, []g.Value{}, 0})

		case g.NativeFunc:

			val, err := fn.Invoke(itp, params)
			if err != nil {
				return nil, err
			}

			f.stack = f.stack[:n-p]
			f.stack = append(f.stack, val)
			f.ip += 3

		default:
			return nil, g.TypeMismatchError("Expected Func")
		}

	case bc.Return:

		//// TODO once we've written a Control Flow Graph
		//// turn this sanity check on to make sure we are managing
		//// the stack properly
		//if len(f.stack) < 1 || len(f.stack) > 2 {
		//	dumpFrame(f)
		//	panic("invalid stack")
		//}

		// get result from top of stack
		result := f.stack[n]

		// pop the old frame
		itp.frames = itp.frames[:frameIndex]

		// If we are on the last frame, then we are done.
		if frameIndex == lastFrame {
			return result, nil
		}

		// push the result onto the new top frame
		f = itp.frames[len(itp.frames)-1]
		f.stack = append(f.stack, result)

		// advance the instruction pointer now that we are done invoking
		f.ip += 3

	case bc.Done:
		panic("Done cannot be executed directly")

	case bc.Go:

		p := bc.DecodeParam(btc, f.ip)
		params := f.stack[n-p+1:]

		switch fn := f.stack[n-p].(type) {
		case bc.BytecodeFunc:
			f.stack = f.stack[:n-p]
			f.ip += 3

			intp := NewInterpreter(itp.builtInMgr, itp.modules)
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
				_, err := fn.Invoke(itp, params)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			})()

		default:
			return nil, g.TypeMismatchError("Expected Func")
		}

	case bc.Throw:

		// get str from stack
		s, ok := f.stack[n].(g.Str)
		if !ok {
			return nil, g.TypeMismatchError("Expected Str")
		}

		// throw an error
		return nil, g.NewError(s.String())

	case bc.NewFunc:

		// push a function
		p := bc.DecodeParam(btc, f.ip)
		tpl := pool.Templates[p]
		nf := bc.NewBytecodeFunc(tpl)
		f.stack = append(f.stack, nf)
		f.ip += 3

	case bc.FuncLocal:

		// get function from stack
		fn, ok := f.stack[n].(bc.BytecodeFunc)
		if !ok {
			return nil, g.TypeMismatchError("Expected BytecodeFunc")
		}

		// push a local onto the captures of the function
		p := bc.DecodeParam(btc, f.ip)
		fn.PushCapture(f.locals[p])
		f.ip += 3

	case bc.FuncCapture:

		// get function from stack
		fn, ok := f.stack[n].(bc.BytecodeFunc)
		if !ok {
			return nil, g.TypeMismatchError("Expected BytecodeFunc")
		}

		// push a capture onto the captures of the function
		p := bc.DecodeParam(btc, f.ip)
		fn.PushCapture(f.fn.GetCapture(p))
		f.ip += 3

	case bc.NewList:

		size := bc.DecodeParam(btc, f.ip)
		vals := make([]g.Value, size)
		copy(vals, f.stack[n-size+1:])

		f.stack = f.stack[:n-size+1]
		f.stack = append(f.stack, g.NewList(vals))
		f.ip += 3

	//	case bc.NewSet:
	//
	//		size := bc.DecodeParam(btc, f.ip)
	//		vals := make([]g.Value, size)
	//		copy(vals, f.stack[n-size+1:])
	//		f.stack = f.stack[:n-size+1]
	//
	//		set, err := g.NewSet(itp, vals)
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		f.stack = append(f.stack, set)
	//		f.ip += 3
	//
	//	case bc.NewTuple:
	//
	//		size := bc.DecodeParam(btc, f.ip)
	//		vals := make([]g.Value, size)
	//		copy(vals, f.stack[n-size+1:])
	//
	//		f.stack = f.stack[:n-size+1]
	//		f.stack = append(f.stack, g.NewTuple(vals))
	//		f.ip += 3
	//
	//	case bc.CheckTuple:
	//
	//		// make sure the top of the stack is really a tuple
	//		tp, ok := f.stack[n].(g.Tuple)
	//		if !ok {
	//			return nil, g.TypeMismatchError("Expected Tuple")
	//		}
	//
	//		// and make sure its of the expected length
	//		expectedLen := bc.DecodeParam(btc, f.ip)
	//		tpLen := tp.Len(itp)
	//		if expectedLen != int(tpLen.IntVal()) {
	//			return nil, g.InvalidArgumentError(
	//				fmt.Sprintf("Expected Tuple of length %d", expectedLen))
	//		}
	//
	//		// do not alter stack
	//		f.ip += 3

	//case bc.CheckCast:

	//	// make sure the top of the stack is of the given type
	//	vtype := g.Type(bc.DecodeParam(btc, f.ip))
	//	v := f.stack[n]
	//	if v.Type() != vtype {
	//		return nil, g.TypeMismatchError(fmt.Sprintf("Expected %s", vtype.String()))
	//	}

	//	// do not alter stack
	//	f.ip += 3

	//	case bc.NewDict:
	//
	//		size := bc.DecodeParam(btc, f.ip)
	//		entries := make([]*g.HEntry, 0, size)
	//
	//		numVals := size * 2
	//		for j := n - numVals + 1; j <= n; j += 2 {
	//			entries = append(entries, &g.HEntry{Key: f.stack[j], Value: f.stack[j+1]})
	//		}
	//
	//		f.stack = f.stack[:n-numVals+1]
	//
	//		dict, err := g.NewDict(itp, entries)
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		f.stack = append(f.stack, dict)
	//		f.ip += 3
	//

	case bc.NewStruct:

		p := bc.DecodeParam(btc, f.ip)
		def := pool.StructDefs[p]
		fields := make(map[string]g.Field)
		for _, name := range def {
			fields[name] = g.NewField(g.Null)
		}

		stc, err := g.NewFieldStruct(fields, false)
		if err != nil {
			return nil, err
		}

		f.stack = append(f.stack, stc)
		f.ip += 3

	case bc.GetField:

		p := bc.DecodeParam(btc, f.ip)
		key, ok := pool.Constants[p].(g.Str)
		g.Assert(ok)

		result, err := f.stack[n].GetField(key.String(), itp)
		if err != nil {
			return nil, err
		}

		f.stack[n] = result
		f.ip += 3

	case bc.InvokeField:

		p, q := bc.DecodeWideParams(btc, f.ip)
		//fmt.Printf("InvokeField aaa %d %d\n", p, q)

		key, ok := pool.Constants[p].(g.Str)
		g.Assert(ok)
		//fmt.Printf("InvokeField bbb %s\n", key)

		self := f.stack[n-q]
		params := f.stack[n-q+1:]
		//fmt.Printf("InvokeField ccc %v %d %d %d\n", params, n, q, n-q+1)

		result, err := self.InvokeField(key.String(), itp, params)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("InvokeField ddd %s\n", result)

		f.stack[n-q] = result
		f.stack = f.stack[:n-q+1]
		f.ip += 5

	case bc.SetField, bc.InitField, bc.InitProperty, bc.InitReadonlyProperty:

		p := bc.DecodeParam(btc, f.ip)
		key, ok := pool.Constants[p].(g.Str)
		g.Assert(ok)

		switch btc[f.ip] {

		case bc.SetField:

			stc, ok := f.stack[n-1].(g.Struct)
			if !ok {
				return nil, g.TypeMismatchError("Expected Struct")
			}
			value := f.stack[n]

			err := stc.SetField(key.String(), itp, value)
			if err != nil {
				return nil, err
			}
			f.stack[n-1] = value
			f.stack = f.stack[:n]

		case bc.InitField:

			stc := f.stack[n-1].(g.Struct)
			value := f.stack[n]

			stc.Internal(key.String(), g.NewField(value))
			f.stack = f.stack[:n]

		case bc.InitProperty:

			stc := f.stack[n-2].(g.Struct)
			get := f.stack[n-1].(g.Func)
			set := f.stack[n].(g.Func)

			prop, err := g.NewProperty(get, set)
			g.Assert(err == nil)

			stc.Internal(key.String(), prop)
			f.stack = f.stack[:n-1]

		case bc.InitReadonlyProperty:

			stc := f.stack[n-1].(g.Struct)
			get := f.stack[n].(g.Func)

			prop, err := g.NewReadonlyProperty(get)
			g.Assert(err == nil)

			stc.Internal(key.String(), prop)
			f.stack = f.stack[:n]

		default:
			panic("unreachable")
		}

		f.ip += 3

	case bc.IncField:

		p := bc.DecodeParam(btc, f.ip)
		key, ok := pool.Constants[p].(g.Str)
		g.Assert(ok)

		// get struct from stack
		stc, ok := f.stack[n-1].(g.Struct)
		if !ok {
			return nil, g.TypeMismatchError("Expected Struct")
		}

		// get value from stack
		value := f.stack[n]

		before, err := stc.GetField(key.String(), itp)
		if err != nil {
			return nil, err
		}

		after, err := inc(itp, before, value)
		if err != nil {
			return nil, err
		}

		err = stc.SetField(key.String(), itp, after)
		if err != nil {
			return nil, err
		}

		f.stack[n-1] = before
		f.stack = f.stack[:n]
		f.ip += 3

	case bc.GetIndex:

		// get Indexable from stack
		gtb, ok := f.stack[n-1].(g.Indexable)
		if !ok {
			return nil, g.TypeMismatchError("Expected Indexable")
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

	case bc.SetIndex:

		// get IndexAssignable from stack
		ibl, ok := f.stack[n-2].(g.IndexAssignable)
		if !ok {
			return nil, g.TypeMismatchError("Expected Indexable")
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

	case bc.IncIndex:

		// get IndexAssignable from stack
		ibl, ok := f.stack[n-2].(g.IndexAssignable)
		if !ok {
			return nil, g.TypeMismatchError("Expected Indexable")
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

	case bc.Slice:

		// get Sliceable from stack
		slb, ok := f.stack[n-2].(g.Sliceable)
		if !ok {
			return nil, g.TypeMismatchError("Expected Sliceable")
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

	case bc.SliceFrom:

		// get Sliceable from stack
		slb, ok := f.stack[n-1].(g.Sliceable)
		if !ok {
			return nil, g.TypeMismatchError("Expected Sliceable")
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

	case bc.SliceTo:

		// get Sliceable from stack
		slb, ok := f.stack[n-1].(g.Sliceable)
		if !ok {
			return nil, g.TypeMismatchError("Expected Sliceable")
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

	case bc.LoadNull:
		f.stack = append(f.stack, g.Null)
		f.ip++
	case bc.LoadTrue:
		f.stack = append(f.stack, g.True)
		f.ip++
	case bc.LoadFalse:
		f.stack = append(f.stack, g.False)
		f.ip++
	case bc.LoadZero:
		f.stack = append(f.stack, g.Zero)
		f.ip++
	case bc.LoadOne:
		f.stack = append(f.stack, g.One)
		f.ip++
	case bc.LoadNegOne:
		f.stack = append(f.stack, g.NegOne)
		f.ip++

	case bc.ImportModule:

		// get the module name from the pool
		p := bc.DecodeParam(btc, f.ip)
		name, ok := pool.Constants[p].(g.Str)
		g.Assert(ok)

		// Lookup the module.
		mod, err := itp.lookupModule(name.String())
		if err != nil {
			return nil, err
		}

		// Push the module's contents onto the stack
		f.stack = append(f.stack, mod.Contents)
		f.ip += 3

	case bc.LoadBuiltin:
		p := bc.DecodeParam(btc, f.ip)
		f.stack = append(f.stack, itp.builtInMgr.Builtins()[p])
		f.ip += 3

	case bc.LoadConst:
		p := bc.DecodeParam(btc, f.ip)
		f.stack = append(f.stack, pool.Constants[p])
		f.ip += 3

	case bc.LoadLocal:
		p := bc.DecodeParam(btc, f.ip)
		f.stack = append(f.stack, f.locals[p].Val)
		f.ip += 3

	case bc.LoadCapture:
		p := bc.DecodeParam(btc, f.ip)
		f.stack = append(f.stack, f.fn.GetCapture(p).Val)
		f.ip += 3

	case bc.StoreLocal:
		p := bc.DecodeParam(btc, f.ip)
		f.locals[p].Val = f.stack[n]
		f.stack = f.stack[:n]
		f.ip += 3

	case bc.StoreCapture:
		p := bc.DecodeParam(btc, f.ip)
		f.fn.GetCapture(p).Val = f.stack[n]
		f.stack = f.stack[:n]
		f.ip += 3

	case bc.Jump:
		f.ip = bc.DecodeParam(btc, f.ip)

	case bc.JumpTrue:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return nil, g.TypeMismatchError("Expected Bool")
		}

		f.stack = f.stack[:n]
		if b.BoolVal() {
			f.ip = bc.DecodeParam(btc, f.ip)
		} else {
			f.ip += 3
		}

	case bc.JumpFalse:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return nil, g.TypeMismatchError("Expected Bool")
		}

		f.stack = f.stack[:n]
		if b.BoolVal() {
			f.ip += 3
		} else {
			f.ip = bc.DecodeParam(btc, f.ip)
		}

	case bc.Eq:
		b, err := f.stack[n-1].Eq(itp, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = b
		f.ip++

	case bc.Ne:
		b, err := f.stack[n-1].Eq(itp, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = b.Not()
		f.ip++

	case bc.Lt:
		val, err := f.stack[n-1].Cmp(itp, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.NewBool(val.IntVal() < 0)
		f.ip++

	case bc.Lte:
		val, err := f.stack[n-1].Cmp(itp, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.NewBool(val.IntVal() <= 0)
		f.ip++

	case bc.Gt:
		val, err := f.stack[n-1].Cmp(itp, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.NewBool(val.IntVal() > 0)
		f.ip++

	case bc.Gte:
		val, err := f.stack[n-1].Cmp(itp, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.NewBool(val.IntVal() >= 0)
		f.ip++

	case bc.Cmp:
		val, err := f.stack[n-1].Cmp(itp, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bc.Plus:
		val, err := plus(itp, f.stack[n-1], f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bc.Inc:
		val, err := inc(itp, f.stack[n-1], f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bc.Not:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return nil, g.TypeMismatchError("Expected Bool")
		}

		f.stack[n] = b.Not()
		f.ip++

	case bc.Sub:
		z, ok := f.stack[n-1].(g.Number)
		if !ok {
			return nil, g.TypeMismatchError("Expected Number")
		}

		val, err := z.Sub(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bc.Mul:
		z, ok := f.stack[n-1].(g.Number)
		if !ok {
			return nil, g.TypeMismatchError("Expected Number")
		}

		val, err := z.Mul(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bc.Div:
		z, ok := f.stack[n-1].(g.Number)
		if !ok {
			return nil, g.TypeMismatchError("Expected Number")
		}

		val, err := z.Div(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case bc.Negate:
		z, ok := f.stack[n].(g.Number)
		if !ok {
			return nil, g.TypeMismatchError("Expected Number")
		}

		val := z.Negate()
		f.stack[n] = val
		f.ip++

	case bc.Rem:
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

	case bc.BitAnd:
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

	case bc.BitOr:
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

	case bc.BitXor:
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

	case bc.LeftShift:
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

	case bc.RightShift:
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

	case bc.Complement:
		z, ok := f.stack[n].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected Int")
		}

		val := z.Complement()
		f.stack[n] = val
		f.ip++

	case bc.NewIter:

		ibl, ok := f.stack[n].(g.Iterable)
		g.Assert(ok)

		f.stack[n] = ibl.NewIterator(itp)
		f.ip++

	case bc.IterNext:

		itr, ok := f.stack[n].(g.Iterator)
		g.Assert(ok)

		val, err := itr.IterNext(itp)
		if err != nil {
			return nil, err
		}

		f.stack[n] = val
		f.ip++

	case bc.IterGet:

		itr, ok := f.stack[n].(g.Iterator)
		g.Assert(ok)

		val, err := itr.IterGet(itp)
		if err != nil {
			return nil, err
		}

		f.stack[n] = val
		f.ip++

	case bc.Dup:
		f.stack = append(f.stack, f.stack[n])
		f.ip++

	case bc.Pop:
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
	return nil, g.TypeMismatchError("Expected Number")
}

func inc(ev g.Evaluator, a g.Value, b g.Value) (g.Value, g.Error) {

	na, ok := a.(g.Number)
	if !ok {
		return nil, g.TypeMismatchError("Expected Number")
	}
	nb := b.(g.Number)
	return na.Add(nb)
}
