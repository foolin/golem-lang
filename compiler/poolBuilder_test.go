// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	"reflect"
	"testing"

	g "github.com/mjarmy/golem-lang/core"
)

func TestPool(t *testing.T) {

	pb := newPoolBuilder()

	tassert(t, pb.constIndex(g.NewInt(4)) == 0)
	tassert(t, pb.constIndex(mustStr("a")) == 1)
	tassert(t, pb.constIndex(g.NewFloat(1.0)) == 2)
	tassert(t, pb.constIndex(mustStr("a")) == 1)
	tassert(t, pb.constIndex(g.NewInt(4)) == 0)

	constants := pb.makeConstants()
	tassert(t, reflect.DeepEqual(
		constants,
		[]g.Basic{
			g.NewInt(4),
			mustStr("a"),
			g.NewFloat(1.0)}))
}
