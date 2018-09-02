// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	//"reflect"
	"testing"

	"github.com/mjarmy/golem-lang/compiler"
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/core/bytecode"
	"github.com/mjarmy/golem-lang/scanner"
)

// NOTE: Almost all of the test suite for the interpreter is over in the
// 'bench_test' directory. Testing is done there by running the Golem CLI
// application against Golem code.

var builtins []*g.BuiltinEntry = nil
var builtinMgr = g.NewBuiltinManager(builtins)

func testCompile(t *testing.T, code string) *bytecode.Module {

	source := &scanner.Source{Name: "foo", Path: "foo.glm", Code: code}
	mods, errs := compiler.CompileSourceFully(builtinMgr, source, nil)
	g.Tassert(t, errs == nil)
	g.Tassert(t, len(mods) == 1)

	return mods[0]
}

func TestModuleContents(t *testing.T) {

	code := `
let a = 0
const b = "xyz"
fn c() {}
`
	mod := testCompile(t, code)

	intp := NewInterpreter(builtinMgr, []*bytecode.Module{mod})

	result, errStruct := intp.InitModules()
	g.Tassert(t, errStruct == nil && len(result) == 1)

	stc := mod.Contents

	val, err := stc.GetField("a", nil)
	g.Tassert(t, err == nil && val.Type() == g.IntType)

	val, err = stc.GetField("b", nil)
	g.Tassert(t, err == nil && val.Type() == g.StrType)

	val, err = stc.GetField("c", nil)
	g.Tassert(t, err == nil && val.Type() == g.FuncType)

	err = stc.SetField("a", nil, g.One)
	g.Tassert(t, err == nil)

	err = stc.SetField("b", nil, g.One)
	g.Tassert(t, err.Error() == "ReadonlyField: Field 'b' is readonly")

	err = stc.SetField("c", nil, g.One)
	g.Tassert(t, err.Error() == "ReadonlyField: Field 'c' is readonly")
}
