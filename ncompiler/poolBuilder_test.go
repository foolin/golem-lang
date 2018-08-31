// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncompiler

import (
	"reflect"
	"testing"

	g "github.com/mjarmy/golem-lang/ncore"
)

func TestPool(t *testing.T) {

	pb := newPoolBuilder()

	g.Tassert(t, pb.constIndex(g.NewInt(4)) == 0)
	g.Tassert(t, pb.constIndex(g.NewStr("a")) == 1)
	g.Tassert(t, pb.constIndex(g.NewFloat(1.0)) == 2)
	g.Tassert(t, pb.constIndex(g.NewStr("a")) == 1)
	g.Tassert(t, pb.constIndex(g.NewInt(4)) == 0)

	constants := pb.makeConstants()
	g.Tassert(t, reflect.DeepEqual(
		constants,
		[]g.Basic{
			g.NewInt(4),
			g.NewStr("a"),
			g.NewFloat(1.0)}))
}
