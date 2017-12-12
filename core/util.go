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

// The value must be between 0 (inclusive) and max (exclusive).
func validateIndex(val Value, max int) (Int, Error) {

	if i, ok := val.(Int); ok {
		n := int(i.IntVal())
		switch {
		case n < 0:
			return nil, IndexOutOfBoundsError()
		case n >= max:
			return nil, IndexOutOfBoundsError()
		default:
			return i, nil
		}
	} else {
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

func valuesEq(as []Value, bs []Value) Bool {

	if len(as) != len(bs) {
		return FALSE
	}

	for i, a := range as {
		if a.Eq(bs[i]) == FALSE {
			return FALSE
		}
	}

	return TRUE
}

func strcat(a Value, b Value) Str {

	sa := a.ToStr().String()
	sb := b.ToStr().String()

	return str(strcpy(sa) + strcpy(sb))
}

// copy to avoid memory leaks
func strcpy(s string) string {
	c := make([]byte, len(s))
	copy(c, s)
	return string(c)
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

func Assert(flag bool, msg string) {
	if !flag {
		panic(msg)
	}
}
