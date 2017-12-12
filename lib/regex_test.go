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

package lib

import (
	g "github.com/mjarmy/golem-lang/core"
	"reflect"
	"testing"
)

func assert(t *testing.T, flag bool) {
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

	compile, err := regex.GetContents().GetField(g.MakeStr("compile"))
	assert(t, compile != nil && err == nil)
	fnCompile := compile.(g.NativeFunc)

	pattern, err := fnCompile.Invoke([]g.Value{g.MakeStr(`^[a-z]+\[[0-9]+\]$`)})
	assert(t, pattern != nil && err == nil)

	match, err := pattern.GetField(g.MakeStr("match"))
	assert(t, match != nil && err == nil)
	fnMatch := match.(g.NativeFunc)

	result, err := fnMatch.Invoke([]g.Value{g.MakeStr("foo[123]")})
	ok(t, result, err, g.TRUE)

	result, err = fnMatch.Invoke([]g.Value{g.MakeStr("456")})
	ok(t, result, err, g.FALSE)

	pattern, err = fnCompile.Invoke([]g.Value{g.MakeStr("\\")})
	assert(t, pattern == nil && err.Error() ==
		"RegexError: error parsing regexp: trailing backslash at end of expression: ``")
}
