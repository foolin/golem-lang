// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package io

import (
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/lib/io/ioutil"
)

/*doc

## io

Module io provides basic interfaces to I/O primitives. Its primary job is to wrap
existing implementations of such primitives, such as those in package os, into shared
public interfaces that abstract the functionality, plus some other related primitives.

Because these interfaces and primitives wrap lower-level operations with various
implementations, unless otherwise informed clients should not assume they are safe
for parallel execution.

*/

/*doc
`io` has the following fields:

* [ioutil](lib_ioioutil.html)

*/

// Io is the "io" module in the standard library
var Io g.Module

func init() {

	ioutil, err := g.NewFrozenStruct(
		map[string]g.Field{
			"readDir":         g.NewField(ioutil.ReadDir),
			"readFileString":  g.NewField(ioutil.ReadFileString),
			"writeFileString": g.NewField(ioutil.WriteFileString),
		})

	g.Assert(err == nil)

	io, err := g.NewFrozenStruct(
		map[string]g.Field{
			"ioutil": g.NewField(ioutil),
		})
	g.Assert(err == nil)

	Io = g.NewNativeModule("io", io)
}
