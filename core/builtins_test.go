// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"
)

func builtinOk(t *testing.T, flag bool) {
	if !flag {
		t.Error("builtins")
	}
}

func TestBuiltins(t *testing.T) {

	mgr := NewBuiltinManager(StandardBuiltins)
	builtinOk(t, !mgr.Contains("foo"))
	builtinOk(t, mgr.Contains("str"))
	builtinOk(t, mgr.Contains("len"))
	builtinOk(t, mgr.IndexOf("str") == 0)
	builtinOk(t, mgr.IndexOf("len") == 1)
	builtinOk(t, mgr.Builtins()[0] == BuiltinStr)
	builtinOk(t, mgr.Builtins()[1] == BuiltinLen)
}
