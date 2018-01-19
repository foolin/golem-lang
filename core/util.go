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

package core

import ()

func posIndex(val Value, length int) (int, Error) {
	i, ok := val.(Int)
	if !ok {
		return -1, TypeMismatchError("Expected 'Int'")
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
		return 0, 0, TypeMismatchError("Expected 'Int'")
	}
	t, err := posIndex(to, length)
	if err != nil {
		return 0, 0, TypeMismatchError("Expected 'Int'")
	}

	return sliceBounds(f, length), sliceBounds(t, length), nil
}

func valuesEq(cx Context, as []Value, bs []Value) (Bool, Error) {

	if len(as) != len(bs) {
		return FALSE, nil
	}

	for i, a := range as {
		eq, err := a.Eq(cx, bs[i])
		if err != nil {
			return nil, err
		}
		if eq == FALSE {
			return FALSE, nil
		}
	}

	return TRUE, nil
}

func strHash(s string) int {

	// https://en.wikipedia.org/wiki/Jenkins_hash_function
	var hash int = 0
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

func newIteratorStruct() Struct {

	// create a struct with fields that have placeholder NULL values
	stc, err := NewStruct([]Field{
		NewField("nextValue", true, NULL),
		NewField("getValue", true, NULL)}, true)
	if err != nil {
		panic("invalid struct")
	}
	return stc
}

func initIteratorStruct(cx Context, itr Iterator) Iterator {

	// initialize the struct fields with functions that refer back to the iterator
	itr.InitField(cx, MakeStr("nextValue"), &nativeFunc{
		0, 0,
		func(cx Context, values []Value) (Value, Error) {
			return itr.IterNext(), nil
		}})
	itr.InitField(cx, MakeStr("getValue"), &nativeFunc{
		0, 0,
		func(cx Context, values []Value) (Value, Error) {
			return itr.IterGet()
		}})
	return itr
}

func assert(flag bool) {
	if !flag {
		panic("assertion failure")
	}
}
