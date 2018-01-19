// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"
	g "github.com/mjarmy/golem-lang/core"
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

	case g.INVOKE:

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

	case g.RETURN:

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
		} else {

			// push the result onto the new top frame
			f = i.frames[len(i.frames)-1]
			f.stack = append(f.stack, result)

			// advance the instruction pointer now that we are done invoking
			f.ip += 3
		}

	case g.DONE:
		panic("DONE cannot be executed directly")

	case g.GO:

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

	case g.THROW:

		// get struct from stack
		stc, ok := f.stack[n].(g.Struct)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Struct'")
		}

		// throw an error
		return nil, g.MakeErrorFromStruct(i, stc)

	case g.NEW_FUNC:

		// push a function
		idx := index(opc, f.ip)
		tpl := i.mod.Templates[idx]
		nf := g.NewBytecodeFunc(tpl)
		f.stack = append(f.stack, nf)
		f.ip += 3

	case g.FUNC_LOCAL:

		// get function from stack
		fn, ok := f.stack[n].(g.BytecodeFunc)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'BytecodeFunc'")
		}

		// push a local onto the captures of the function
		idx := index(opc, f.ip)
		fn.PushCapture(f.locals[idx])
		f.ip += 3

	case g.FUNC_CAPTURE:

		// get function from stack
		fn, ok := f.stack[n].(g.BytecodeFunc)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'BytecodeFunc'")
		}

		// push a capture onto the captures of the function
		idx := index(opc, f.ip)
		fn.PushCapture(f.fn.GetCapture(idx))
		f.ip += 3

	case g.NEW_STRUCT:

		def := i.mod.StructDefs[index(opc, f.ip)]
		stc, err := g.NewStruct(def, false)
		if err != nil {
			return nil, err
		}

		f.stack = append(f.stack, stc)
		f.ip += 3

	case g.NEW_LIST:

		size := index(opc, f.ip)
		vals := make([]g.Value, size)
		copy(vals, f.stack[n-size+1:])

		f.stack = f.stack[:n-size+1]
		f.stack = append(f.stack, g.NewList(vals))
		f.ip += 3

	case g.NEW_SET:

		size := index(opc, f.ip)
		vals := make([]g.Value, size)
		copy(vals, f.stack[n-size+1:])

		f.stack = f.stack[:n-size+1]
		f.stack = append(f.stack, g.NewSet(i, vals))
		f.ip += 3

	case g.NEW_TUPLE:

		size := index(opc, f.ip)
		vals := make([]g.Value, size)
		copy(vals, f.stack[n-size+1:])

		f.stack = f.stack[:n-size+1]
		f.stack = append(f.stack, g.NewTuple(vals))
		f.ip += 3

	case g.CHECK_TUPLE:

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

	case g.CHECK_CAST:

		// make sure the top of the stack is of the given type
		vtype := g.Type(index(opc, f.ip))
		v := f.stack[n]
		if v.Type() != vtype {
			return nil, g.TypeMismatchError(fmt.Sprintf("Expected '%s'", vtype.String()))
		}

		// do not alter stack
		f.ip += 3

	case g.NEW_DICT:

		size := index(opc, f.ip)
		entries := make([]*g.HEntry, 0, size)

		numVals := size * 2
		for j := n - numVals + 1; j <= n; j += 2 {
			entries = append(entries, &g.HEntry{f.stack[j], f.stack[j+1]})
		}

		f.stack = f.stack[:n-numVals+1]
		f.stack = append(f.stack, g.NewDict(i, entries))
		f.ip += 3

	case g.GET_FIELD:

		idx := index(opc, f.ip)
		key, ok := pool[idx].(g.Str)
		assert(ok)

		result, err := f.stack[n].GetField(i, key)
		if err != nil {
			return nil, err
		}

		f.stack[n] = result
		f.ip += 3

	case g.INIT_FIELD, g.SET_FIELD:

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
		if opc[f.ip] == g.INIT_FIELD {
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

	case g.INC_FIELD:

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

	case g.GET_INDEX:

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

	case g.SET_INDEX:

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

	case g.INC_INDEX:

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

	case g.SLICE:

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

	case g.SLICE_FROM:

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

	case g.SLICE_TO:

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

	case g.LOAD_NULL:
		f.stack = append(f.stack, g.NULL)
		f.ip++
	case g.LOAD_TRUE:
		f.stack = append(f.stack, g.TRUE)
		f.ip++
	case g.LOAD_FALSE:
		f.stack = append(f.stack, g.FALSE)
		f.ip++
	case g.LOAD_ZERO:
		f.stack = append(f.stack, g.ZERO)
		f.ip++
	case g.LOAD_ONE:
		f.stack = append(f.stack, g.ONE)
		f.ip++
	case g.LOAD_NEG_ONE:
		f.stack = append(f.stack, g.NEG_ONE)
		f.ip++

	case g.IMPORT_MODULE:

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

	case g.LOAD_BUILTIN:
		idx := index(opc, f.ip)
		f.stack = append(f.stack, i.builtInMgr.Builtins()[idx])
		f.ip += 3

	case g.LOAD_CONST:
		idx := index(opc, f.ip)
		f.stack = append(f.stack, pool[idx])
		f.ip += 3

	case g.LOAD_LOCAL:
		idx := index(opc, f.ip)
		f.stack = append(f.stack, f.locals[idx].Val)
		f.ip += 3

	case g.LOAD_CAPTURE:
		idx := index(opc, f.ip)
		f.stack = append(f.stack, f.fn.GetCapture(idx).Val)
		f.ip += 3

	case g.STORE_LOCAL:
		idx := index(opc, f.ip)
		f.locals[idx].Val = f.stack[n]
		f.stack = f.stack[:n]
		f.ip += 3

	case g.STORE_CAPTURE:
		idx := index(opc, f.ip)
		f.fn.GetCapture(idx).Val = f.stack[n]
		f.stack = f.stack[:n]
		f.ip += 3

	case g.JUMP:
		f.ip = index(opc, f.ip)

	case g.JUMP_TRUE:
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

	case g.JUMP_FALSE:
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

	case g.EQ:
		b, err := f.stack[n-1].Eq(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = b
		f.ip++

	case g.NE:
		b, err := f.stack[n-1].Eq(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = b.Not()
		f.ip++

	case g.LT:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.MakeBool(val.IntVal() < 0)
		f.ip++

	case g.LTE:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.MakeBool(val.IntVal() <= 0)
		f.ip++

	case g.GT:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.MakeBool(val.IntVal() > 0)
		f.ip++

	case g.GTE:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = g.MakeBool(val.IntVal() >= 0)
		f.ip++

	case g.CMP:
		val, err := f.stack[n-1].Cmp(i, f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.HAS:

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

	case g.PLUS:
		val, err := plus(i, f.stack[n-1], f.stack[n])
		if err != nil {
			return nil, err
		}
		f.stack = f.stack[:n]
		f.stack[n-1] = val
		f.ip++

	case g.NOT:
		b, ok := f.stack[n].(g.Bool)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Bool'")
		}

		f.stack[n] = b.Not()
		f.ip++

	case g.SUB:
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

	case g.MUL:
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

	case g.DIV:
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

	case g.NEGATE:
		z, ok := f.stack[n].(g.Number)
		if !ok {
			return nil, g.TypeMismatchError("Expected Number Type")
		}

		val := z.Negate()
		f.stack[n] = val
		f.ip++

	case g.REM:
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

	case g.BIT_AND:
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

	case g.BIT_OR:
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

	case g.BIT_XOR:
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

	case g.LEFT_SHIFT:
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

	case g.RIGHT_SHIFT:
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

	case g.COMPLEMENT:
		z, ok := f.stack[n].(g.Int)
		if !ok {
			return nil, g.TypeMismatchError("Expected 'Int'")
		}

		val := z.Complement()
		f.stack[n] = val
		f.ip++

	case g.ITER:

		ibl, ok := f.stack[n].(g.Iterable)
		assert(ok)

		f.stack[n] = ibl.NewIterator(i)
		f.ip++

	case g.ITER_NEXT:

		itr, ok := f.stack[n].(g.Iterator)
		assert(ok)

		f.stack[n] = itr.IterNext()
		f.ip++

	case g.ITER_GET:

		itr, ok := f.stack[n].(g.Iterator)
		assert(ok)

		val, err := itr.IterGet()
		if err != nil {
			return nil, err
		}

		f.stack[n] = val
		f.ip++

	case g.DUP:
		f.stack = append(f.stack, f.stack[n])
		f.ip++

	case g.POP:
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
