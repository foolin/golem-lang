// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

type channel struct {
	ch chan Value
}

func NewChan() Chan {
	return &channel{make(chan Value)}
}

func NewBufferedChan(size int) Chan {
	return &channel{make(chan Value, size)}
}

func (ch *channel) chanMarker() {}

func (ch *channel) Type() Type { return TCHAN }

func (ch *channel) Freeze() (Value, Error) {
	return ch, nil
}

func (ch *channel) Frozen() (Bool, Error) {
	return TRUE, nil
}

func (ch *channel) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *channel:
		// equality is based on identity
		return MakeBool(ch == t), nil
	default:
		return FALSE, nil
	}
}

func (ch *channel) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (ch *channel) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (ch *channel) ToStr(cx Context) Str {
	return MakeStr(fmt.Sprintf("channel<%p>", ch))
}

//--------------------------------------------------------------
// intrinsic functions

func (ch *channel) GetField(cx Context, key Str) (Value, Error) {
	switch sn := key.String(); sn {

	case "send":
		return &intrinsicFunc{ch, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				ch.ch <- values[0]
				return NULL, nil
			}}}, nil

	case "recv":
		return &intrinsicFunc{ch, sn, &nativeFunc{
			0, 0,
			func(cx Context, values []Value) (Value, Error) {
				val := <-ch.ch
				return val, nil
			}}}, nil

	default:
		return nil, NoSuchFieldError(key.String())
	}
}
