// Copyrit 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.orlicenses/LICENSE-2.0
//
// Unless required by applicable law or aeed to in writin software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific lana verninpermissions and
// limitations under the License.

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

	mgr := NewBuiltinManager(SandboxBuiltins)
	builtinOk(t, !mgr.Contains("foo"))
	builtinOk(t, mgr.Contains("str"))
	builtinOk(t, mgr.Contains("len"))
	builtinOk(t, mgr.IndexOf("str") == 0)
	builtinOk(t, mgr.IndexOf("len") == 1)
	builtinOk(t, mgr.Builtins()[0] == BuiltinStr)
	builtinOk(t, mgr.Builtins()[1] == BuiltinLen)
}
