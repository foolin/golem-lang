// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"
)

func Assert(flag bool) {
	if !flag {
		panic("assertion failure")
	}
}

func Tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
		panic("Tassert")
	}
}

func posIndex(val Value, length int) (int, Error) {
	i, ok := val.(Int)
	if !ok {
		return -1, TypeMismatchError("Expected Int")
	}

	n := int(i.IntVal())
	if n < 0 {
		n = length + n
	}
	return n, nil
}

func boundedIndex(val Value, length int) (int, Error) {
	n, err := posIndex(val, length)
	if err != nil {
		return -1, err
	}

	switch {
	case n < 0:
		return -1, IndexOutOfBoundsError(n)
	case n >= length:
		return -1, IndexOutOfBoundsError(n)
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
		return 0, 0, TypeMismatchError("Expected Int")
	}
	t, err := posIndex(to, length)
	if err != nil {
		return 0, 0, TypeMismatchError("Expected Int")
	}

	return sliceBounds(f, length), sliceBounds(t, length), nil
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

// copy to avoid memory leaks
func strcpy(s string) string {
	c := make([]byte, len(s))
	copy(c, s)
	return string(c)
}

//func valuesEq(ev Evaluator, as []Value, bs []Value) (Bool, Error) {
//
//	if len(as) != len(bs) {
//		return False, nil
//	}
//
//	for i, a := range as {
//		eq, err := a.Eq(ev, bs[i])
//		if err != nil {
//			return nil, err
//		}
//		if eq == False {
//			return False, nil
//		}
//	}
//
//	return True, nil
//}

//func newIteratorStruct() Struct {
//
//	// create a struct with fields that have placeholder Null values
//	stc, err := NewStruct([]Field{
//		NewField("next", true, Null),
//		NewField("get", true, Null)}, true)
//	if err != nil {
//		panic("invalid iterator")
//	}
//	return stc
//}
//
//func initIteratorStruct(ev Evaluator, itr Iterator) Iterator {
//
//	// initialize the struct fields with functions that refer back to the iterator
//	err := itr.InitField(
//		ev, NewStr("next"),
//		NewFixedNativeFunc(
//			[]Type{}, false,
//			func(ev Evaluator, values []Value) (Value, Error) {
//				return itr.IterNext(), nil
//			}))
//	if err != nil {
//		panic("invalid iterator")
//	}
//
//	err = itr.InitField(
//		ev, NewStr("get"),
//		NewFixedNativeFunc(
//			[]Type{}, false,
//			func(ev Evaluator, values []Value) (Value, Error) {
//				return itr.IterGet()
//			}))
//	if err != nil {
//		panic("invalid iterator")
//	}
//
//	return itr
//}
