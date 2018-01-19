// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	g "github.com/mjarmy/golem-lang/core"
	"testing"
)

func TestSys(t *testing.T) {
	sys := InitSysModule()
	exit, err := sys.GetContents().GetField(nil, g.MakeStr("exit"))
	tassert(t, exit != nil)
	tassert(t, err == nil)
}
