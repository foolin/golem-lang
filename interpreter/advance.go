// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"
	g "github.com/mjarmy/golem-lang/core"
	o "github.com/mjarmy/golem-lang/core/opcodes"
	"github.com/mjarmy/golem-lang/lib"
)

// Advance the interpreter forwards by one opcode.
func (i *Interpreter) advance(lastFrame int) (g.Value, g.Error) {

	pool := i.mod.Pool
	frameIndex := len(i.frames) - 1
	f := i.frames[frameIndex]
	n := len(f.stack) - 1
	opc := f.fn.Template().OpCodes

	switch opc[f.ip] {

	case o.Invoke:

		idx := index(opc, f.ip)
		params := f.stack[n-idx+1:]

		switch fn := f.stack[n-idx].(type) {
		case g.BytecodeFunc:

			// check arity
			arity := len(params)
			if arity != fn.Template().Arity {
				err := g.ArityMismatchError(
					fmt.Sprintf("%d", fn.Template().Arity), arity)
				return nil, err
			}

			// pop from stack
			f.stack = f.stack[:n-idx]

			// push a new frame
			locals := newLocals(fn.Template().NumLocals, params)
			i.frames = append(i.frames, &frame{fn, locals, []g.Value{}, 0})

		case g.NativeFunc:

			val, err := fn.Invoke(i, params)
			if err != nil {
				return nil, err
			}

			f.stack = f.stack[:n-idx]
			f.stack = append(f.stack, val)
			f.ip += 3

		default:
			return nil, g.TypeMismatchError("Expected 'Func'")
		}

	case o.Return:

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

	case o.Done:
		panic("Done cannot be executed directly")

	case o.Go:

		idx := index(opc, f.ip)
		params := f.stack[n-idx+1:]

		switch fn := f.stack[n-idx].(type) {
		case g.BytecodeFunc:
			f.stack = f.stack[:n-idx]
			f.ip += 3

			intp := &Interpreter{i.mod, i.builtInMgr, []*frame{}}
			locals := newLocals(fn.Template().NumLocals, params)
			go (func() {
				_, errTrace := intp.eval(fn, locals)
				if errTrace != nil {
					fmt.Printf("%v\n", errTrace.Error)
				}
			})()

		case g.NativeFunc:
			f.stack = f.stack[:n-idx]
			f.ip += 3

			go (func() {
				_, err := fn.Invoke(i, params)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			})()

		default:
			return nil, g.TypeMismatchError("Expected 'Func'")
		}

	case o.Throw:

		// get struct from stack
		stc, ok := f.stack[n].(g.Struct)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Struct'")
		}

		// throw an error
		return nil, g.NewErrorFromStruct(i, stc)

	case o.NewFunc:

		// push a function
		idx := index(opc, f.ip)
		tpl := i.mod.Templates[idx]
		nf := g.NewBytecodeFunc(tpl)
		f.stack = append(f.stack, nf)
		f.ip += 3

	case o.FuncLocal:

		// get function from stack
		fn, ok := f.stack[n].(g.BytecodeFunc)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'BytecodeFunc'")
		}

		// push a local onto the captures of the function
		idx := index(opc, f.ip)
		fn.PushCapture(f.locals[idx])
		f.ip += 3

	case o.FuncCapture:

		// get function from stack
		fn, ok := f.stack[n].(g.BytecodeFunc)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'BytecodeFunc'")
		}

		// push a capture onto the captures of the function
		idx := index(opc, f.ip)
		fn.PushCapture(f.fn.GetCapture(idx))
		f.ip += 3

	case o.NewStruct:

		def := i.mod.StructDefs[index(opc, f.ip)]
		stc, err := g.NewStruct(def, false)
		if err != nil {
			return nil, err
		}

		f.stack = append(f.stack, stc)
		f.ip += 3

	case o.NewList:

		size := index(opc, f.ip)
		vals := make([]g.Value, size)
		copy(vals, f.stack[n-size+1:])

		f.stack = f.stack[:n-size+1]
		f.stack = append(f.stack, g.NewList(vals))
		f.ip += 3

	case o.NewSet:

		size := index(opc, f.ip)
		vals := make([]g.Value, size)
		copy(vals, f.stack[n-size+1:])

		f.stack = f.stack[:n-size+1]
		f.stack = append(f.stack, g.NewSet(i, vals))
		f.ip += 3

	case o.NewTuple:

		size := index(opc, f.ip)
		vals := make([]g.Value, size)
		copy(vals, f.stack[n-size+1:])

		f.stack = f.stack[:n-size+1]
		f.stack = append(f.stack, g.NewTuple(vals))
		f.ip += 3

	case o.CheckTuple:

		// make sure the top of the stack is really a tuple
		tp, ok := f.stack[n].(g.Tuple)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Tuple'")
		}

		// and make sure its of the expected length
		expectedLen := index(opc, f.ip)
		tpLen := tp.Len()
		if expectedLen != int(tpLen.IntVal()) {
			return nil, g.InvalidArgumentError(
				fmt.Sprintf("Expected Tuple of length %d", expectedLen))
		}

		// do not alter stack
		f.ip += 3

	case o.CheckCast:

		// make sure the top of the stack is of the given type
		vtype := g.Type(index(opc, f.ip))
		v := f.stack[n]
		if v.Type() != vtype {
			return nil, g.TypeMismatchError(fmt.Sprintf("Expected '%s'", vtype.String()))
		}

		// do not alter stack
		f.ip += 3

	case o.NewDict:

		size := index(opc, f.ip)
		entries := make([]*g.HEntry, 0, size)

		numVals := size * 2
		for j := n - numVals + 1; j <= n; j += 2 {
			entries = append(entries, &g.HEntry{f.stack[j], f.stack[j+1]})
		}

		f.stack = f.stack[:n-numVals+1]
		f.stack = append(f.stack, g.NewDict(i, entries))
		f.ip += 3

	case o.GetField:

		idx := index(opc, f.ip)
		key, ok := pool[idx].(g.Str)
		assert(ok)

		result, err := f.stack[n].GetField(i, key)
		if err != nil {
			return nil, err
		}

		f.stack[n] = result
		f.ip += 3

	case o.InitField, o.SetField:

		idx := index(opc, f.ip)
		key, ok := pool[idx].(g.Str)
		assert(ok)

		// get struct from stack
		stc, ok := f.stack[n-1].(g.Struct)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Struct'")
		}

		// get value from stack
		value := f.stack[n]

		// init or set
		if opc[f.ip] == o.InitField {
			err := stc.InitField(i, key, value)
			if err != nil {
				return nil, err
			}
		} else {
			err := stc.SetField(i, key, value)
			if err != nil {
				return nil, err
			}
		}

		f.stack[n-1] = value
		f.stack = f.stack[:n]
		f.ip += 3

	case o.IncField:

		idx := index(opc, f.ip)
		key, ok := pool[idx].(g.Str)
		assert(ok)

		// get struct from stack
		stc, ok := f.stack[n-1].(g.Struct)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Struct'")
		}

		// get value from stack
		value := f.stack[n]

		before, err := stc.GetField(i, key)
		if err != nil {
			return nil, err
		}

		after, err := plus(i, before, value)
		if err != nil {
			return nil, err
		}

		err = stc.SetField(i, key, after)
		if err != nil {
			return nil, err
		}

		f.stack[n-1] = before
		f.stack = f.stack[:n]
		f.ip += 3

	case o.GetIndex:

		// get Getable from stack
		gtb, ok := f.stack[n-1].(g.Getable)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Getable'")
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

	case o.SetIndex:

		// get Indexable from stack
		ibl, ok := f.stack[n-2].(g.Indexable)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Getable'")
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

	case o.IncIndex:

		// get Indexable from stack
		ibl, ok := f.stack[n-2].(g.Indexable)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Getable'")
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

	case o.Slice:

		// get Sliceable from stack
		slb, ok := f.stack[n-2].(g.Sliceable)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Sliceable'")
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

	case o.SliceFrom:

		// get Sliceable from stack
		slb, ok := f.stack[n-1].(g.Sliceable)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Sliceable'")
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

	case o.SliceTo:

		// get Sliceable from stack
		slb, ok := f.stack[n-1].(g.Sliceable)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Sliceable'")
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

	case o.LoadNull:
		f.stack = append(f.stack, g.NullValue)
		f.ip++
	case o.LoadTrue:
		f.stack = append(f.stack, g.TRUE)
		f.ip++
	case o.LoadFalse:
		f.stack = append(f.stack, g.FALSE)
		f.ip++
	case o.LoadZero:
		f.stack = append(f.stack, g.Zero)
		f.ip++
	case o.LoadOne:
		f.stack = append(f.stack, g.One)
		f.ip++
	case o.LoadNegOne:
		f.stack = append(f.stack, g.NegOne)
		f.ip++

	case o.ImportModule:

		// An error here is "impossible", because we already
		// made sure the module really existed in the Analyzer.

		// get the module name from the pool
		idx := index(opc, f.ip)
		name, ok := pool[idx].(g.Str)
		assert(ok)

		// Lookup the module.
		mod, err := lib.LookupModule(name.String())
		assert(err == nil)

		// Push the module's contents onto the stack
		f.stack = append(f.stack, mod.GetContents())
		f.ip += 3

	case o.LoadBuiltin:
		idx := index(opc, f.ip)
		f.stack = append(f.stack, i.builtInMgr.Builtins()[idx])
		f.ip += 3

	case o.LoadConst:
		idx := index(opc, f.ip)
		f.stack = append(f.stack, pool[idx])
		f.ip += 3

	case o.LoadLocal:
		idx := index(opc, f.ip)
		f.stack = append(f.stack, f.locals[idx].Val)
		f.ip += 3

	case o.LoadCapture:
		idx := index(opc, f.ip)
		f.stack = append(f.stack, f.fn.GetCapture(idx).Val)
		f.ip += 3

	case o.StoreLocal:
		idx := index(opc, f.ip)
		f.locals[idx].Val = f.stack[n]
		f.stack = f.stack[:n]
		f.ip += 3

	case o.StoreCapture:
		idx := index(opc, f.ip)
		f.fn.GetCapture(idx).Val = f.stack[n]
		f.stack = f.stack[:n]
		f.ip += 3

	case o.Jump:
		f.ip = index(opc, f.ip)

	case o.JumpTrue:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Bool'")
		}

		f.stack = f.stack[:n]
		if b.BoolVal() {
			f.ip = index(opc, f.ip)
		} else {
			f.ip += 3
		}

	case o.JumpFalse:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Bool'")
		}

		f.stack = f.stack[:n]
		if b.BoolVal() {
			f.ip += 3
		} else {
			f.ip = index(opc, f.ip)
		}

	case o.Eq:
		b, err := f.stack[n-1].Eq(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = b
		f.ip++

	case o.Ne:
		b, err := f.stack[n-1].Eq(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = b.Not()
		f.ip++

	case o.Lt:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.NewBool(val.IntVal() < 0)
		f.ip++

	case o.Lte:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.NewBool(val.IntVal() <= 0)
		f.ip++

	case o.Gt:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.NewBool(val.IntVal() > 0)
		f.ip++

	case o.Gte:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.NewBool(val.IntVal() >= 0)
		f.ip++

	case o.Cmp:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case o.Has:

		// get struct from stack
		stc, ok := f.stack[n-1].(g.Struct)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Struct'")
		}

		val, err := stc.Has(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case o.Plus:
		val, err := plus(i, f.stack[n-1], f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case o.Not:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Bool'")
		}

		f.stack[n] = b.Not()
		f.ip++

	case o.Sub:
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

	case o.Mul:
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

	case o.Div:
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

	case o.Negate:
		z, ok := f.stack[n].(g.Number)
		if !ok {
			return nil, g.TypeMismatchError("Expected Number Type")
		}

		val := z.Negate()
		f.stack[n] = val
		f.ip++

	case o.Rem:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Int'")
		}

		val, err := z.Rem(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case o.BitAnd:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Int'")
		}

		val, err := z.BitAnd(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case o.BitOr:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Int'")
		}

		val, err := z.BitOr(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case o.BitXor:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Int'")
		}

		val, err := z.BitXOr(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case o.LeftShift:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Int'")
		}

		val, err := z.LeftShift(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case o.RightShift:
		z, ok := f.stack[n-1].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Int'")
		}

		val, err := z.RightShift(f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case o.Complement:
		z, ok := f.stack[n].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Int'")
		}

		val := z.Complement()
		f.stack[n] = val
		f.ip++

	case o.Iter:

		ibl, ok := f.stack[n].(g.Iterable)
		assert(ok)

		f.stack[n] = ibl.NewIterator(i)
		f.ip++

	case o.IterNext:

		itr, ok := f.stack[n].(g.Iterator)
		assert(ok)

		f.stack[n] = itr.IterNext()
		f.ip++

	case o.IterGet:

		itr, ok := f.stack[n].(g.Iterator)
		assert(ok)

		val, err := itr.IterGet()
		if err != nil {
			return nil, err
		}

		f.stack[n] = val
		f.ip++

	case o.Dup:
		f.stack = append(f.stack, f.stack[n])
		f.ip++

	case o.Pop:
		f.stack = f.stack[:n]
		f.ip++

	default:
		panic("Invalid opcode")
	}

	return nil, nil
}

func plus(cx g.Context, a g.Value, b g.Value) (g.Value, g.Error) {

	// if either is a Str, return concatenated strings
	_, ia := a.(g.Str)
	_, ib := b.(g.Str)
	if ia || ib {
		return a.ToStr(cx).Concat(b.ToStr(cx)), nil
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

func index(opcodes []byte, ip int) int {
	high := opcodes[ip+1]
	low := opcodes[ip+2]
	return int(high)<<8 + int(low)
}
