// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	"strings"
	"testing"
)

//--------------------------------------------------------------

func (s str) strContains(cx Context, params []Value) (Value, Error) {
	substr := params[0].(Str)
	return NewBool(strings.Contains(s.String(), substr.String())), nil
}

func show(val Value) {
	//println(val.ToStr(nil).String())
}

const iterate = 4 * 1000 * 1000

func TestDirectCall(t *testing.T) {

	s := NewStr("abcdef")
	substr := NewStr("cd")

	var val Value
	var err Error

	for i := 0; i < iterate; i++ {
		val, err = s.(str).strContains(nil, []Value{substr})
		ok(t, val, err, True)
		show(val)
	}
}

///////////////////////////////////////
// test hash lookup vs switch
///////////////////////////////////////

//func TestInvokeCall(t *testing.T) {
//
//	s := NewStr("abcdef")
//	substr := NewStr("cd")
//
//	var val Value
//	var err Error
//
//	for i := 0; i < iterate; i++ {
//		val, err = s.InvokeField("contains", nil, []Value{substr})
//		ok(t, val, err, True)
//		show(val)
//	}
//}
//
//func TestGetCall(t *testing.T) {
//
//	s := NewStr("abcdef")
//	substr := NewStr("cd")
//
//	var val Value
//	var err Error
//	var field Value
//
//	for i := 0; i < iterate; i++ {
//		field, err = s.GetField("contains", nil)
//		tassert(t, err == nil)
//
//		fn := field.(Func)
//		val, err = fn.Invoke(nil, []Value{substr})
//		ok(t, val, err, True)
//		show(val)
//	}
//}
