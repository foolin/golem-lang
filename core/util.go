// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

// Assert asserts that something is true
func Assert(flag bool) {
	if !flag {
		panic("assertion failure")
	}
}

// CopyValues makes a copy of a slice of Values
func CopyValues(v []Value) []Value {
	c := make([]Value, len(v))
	copy(c, v)
	return c
}

func posIndex(v Value, length int) (int, Error) {

	i, ok := v.(Int)
	if !ok {
		return -1, TypeMismatch(IntType, v.Type())
	}

	n := int(i.IntVal())
	if n < 0 {
		n = length + n
	}
	return n, nil
}

func boundedIndex(v Value, length int) (int, Error) {
	n, err := posIndex(v, length)
	if err != nil {
		return -1, err
	}

	switch {
	case n < 0:
		return -1, IndexOutOfBounds(n)
	case n >= length:
		return -1, IndexOutOfBounds(n)
	default:
		return n, nil
	}
}

func sliceBounds(n int, length int) int {
	switch {
	case n < 0:
		return 0
	case n > length:
		return length
	default:
		return n
	}
}

func sliceIndices(from Value, to Value, length int) (int, int, Error) {
	f, err := posIndex(from, length)
	if err != nil {
		return -1, -1, err
	}
	t, err := posIndex(to, length)
	if err != nil {
		return -1, -1, err
	}

	fb := sliceBounds(f, length)
	tb := sliceBounds(t, length)
	if fb > tb {
		tb = fb
	}

	return fb, tb, nil
}

// https://en.wikipedia.org/wiki/Jenkins_hash_function
func strHash(s string) int {

	var hash int
	bytes := []byte(s)
	for _, b := range bytes {
		hash += int(b)
		hash += hash << 10
		hash ^= hash >> 6
	}
	hash += hash << 3
	hash ^= hash >> 11
	hash += hash << 15
	return hash
}

func valuesEq(ev Eval, as []Value, bs []Value) (Bool, Error) {

	if len(as) != len(bs) {
		return False, nil
	}

	for i, a := range as {
		eq, err := a.Eq(ev, bs[i])
		if err != nil {
			return nil, err
		}
		if eq == False {
			return False, nil
		}
	}

	return True, nil
}

func iteratorStruct() Struct {

	stc, err := NewFrozenFieldStruct(
		map[string]Field{
			"next": NewReadonlyField(Null),
			"get":  NewReadonlyField(Null),
		})
	if err != nil {
		panic("unreachable")
	}
	return stc
}

func iteratorFields(ev Eval, itr Iterator) (next Field, get Field) {

	next = NewReadonlyField(
		NewNullaryNativeFunc(
			func(ev Eval) (Value, Error) {
				return itr.IterNext(ev)
			}))

	get = NewReadonlyField(
		NewNullaryNativeFunc(
			func(ev Eval) (Value, Error) {
				return itr.IterGet(ev)
			}))

	return
}
