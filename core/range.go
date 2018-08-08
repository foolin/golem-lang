// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"math"
	"reflect"
)

//---------------------------------------------------------------
// rng

type rng struct {
	from  int64
	to    int64
	step  int64
	count int64
}

// NewRange creates a new Range
func NewRange(from int64, to int64, step int64) (Range, Error) {

	switch {

	case step == 0:
		return nil, InvalidArgumentError("step cannot be 0")

	case ((step > 0) && (from > to)) || ((step < 0) && (from < to)):
		return &rng{from, to, step, 0}, nil

	default:
		count := int64(math.Ceil(float64(to-from) / float64(step)))
		return &rng{from, to, step, count}, nil
	}
}

func (r *rng) compositeMarker() {}

func (r *rng) Type() Type { return RangeType }

func (r *rng) Freeze() (Value, Error) {
	return r, nil
}

func (r *rng) Frozen() (Bool, Error) {
	return True, nil
}

func (r *rng) ToStr(cx Context) Str {
	return NewStr(fmt.Sprintf("range<%d, %d, %d>", r.from, r.to, r.step))
}

func (r *rng) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (r *rng) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *rng:
		return NewBool(reflect.DeepEqual(r, t)), nil
	default:
		return False, nil
	}
}

func (r *rng) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (r *rng) Get(cx Context, index Value) (Value, Error) {
	idx, err := boundedIndex(index, int(r.count))
	if err != nil {
		return nil, err
	}
	return NewInt(r.from + int64(idx)*r.step), nil
}

func (r *rng) Len() Int {
	return NewInt(r.count)
}

func (r *rng) From() Int  { return NewInt(r.from) }
func (r *rng) To() Int    { return NewInt(r.to) }
func (r *rng) Step() Int  { return NewInt(r.step) }
func (r *rng) Count() Int { return NewInt(r.count) }

//---------------------------------------------------------------
// Iterator

type rangeIterator struct {
	Struct
	r *rng
	n int64
}

func (r *rng) NewIterator(cx Context) Iterator {
	return initIteratorStruct(cx,
		&rangeIterator{newIteratorStruct(), r, -1})
}

func (i *rangeIterator) IterNext() Bool {
	i.n++
	return NewBool(i.n < i.r.count)
}

func (i *rangeIterator) IterGet() (Value, Error) {

	if (i.n >= 0) && (i.n < i.r.count) {
		return NewInt(i.r.from + i.n*i.r.step), nil
	}
	return nil, NoSuchElementError()
}

//--------------------------------------------------------------
// intrinsic functions

func (r *rng) GetField(cx Context, key Str) (Value, Error) {
	switch sn := key.String(); sn {

	case "from":
		return &intrinsicFunc{r, sn, NewNativeFunc0(
			func(cx Context) (Value, Error) {
				return NewInt(r.from), nil
			})}, nil

	case "to":
		return &intrinsicFunc{r, sn, NewNativeFunc0(
			func(cx Context) (Value, Error) {
				return NewInt(r.to), nil
			})}, nil

	case "step":
		return &intrinsicFunc{r, sn, NewNativeFunc0(
			func(cx Context) (Value, Error) {
				return NewInt(r.step), nil
			})}, nil

	case "count":
		return &intrinsicFunc{r, sn, NewNativeFunc0(
			func(cx Context) (Value, Error) {
				return NewInt(r.count), nil
			})}, nil

	default:
		return nil, NoSuchFieldError(key.String())
	}
}
