// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	g "github.com/mjarmy/golem-lang/core"
	"reflect"
	"testing"
)

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
	}
}
func ok(t *testing.T, val g.Value, err g.Error, expect g.Value) {
	if err != nil {
		panic("ok")
		t.Error(err, " != ", nil)
	}

	if !reflect.DeepEqual(val, expect) {
		t.Error(val, " != ", expect)
	}
}

func TestRegex(t *testing.T) {
	regex := InitRegexModule()

	compile, err := regex.GetContents().GetField(nil, g.MakeStr("compile"))
	tassert(t, compile != nil && err == nil)
	fnCompile := compile.(g.NativeFunc)

	pattern, err := fnCompile.Invoke(nil, []g.Value{g.MakeStr(`^[a-z]+\[[0-9]+\]$`)})
	tassert(t, pattern != nil && err == nil)

	match, err := pattern.GetField(nil, g.MakeStr("match"))
	tassert(t, match != nil && err == nil)
	fnMatch := match.(g.NativeFunc)

	result, err := fnMatch.Invoke(nil, []g.Value{g.MakeStr("foo[123]")})
	ok(t, result, err, g.TRUE)

	result, err = fnMatch.Invoke(nil, []g.Value{g.MakeStr("456")})
	ok(t, result, err, g.FALSE)

	pattern, err = fnCompile.Invoke(nil, []g.Value{g.MakeStr("\\")})
	tassert(t, pattern == nil && err.Error() ==
		"RegexError: error parsing regexp: trailing backslash at end of expression: ``")
}
