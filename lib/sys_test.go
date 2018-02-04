// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	"testing"

	g "github.com/mjarmy/golem-lang/core"
)

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
	}
}

func TestSys(t *testing.T) {
	sys := NewSysModule()
	exit, err := sys.GetContents().GetField(nil, g.NewStr("exit"))
	tassert(t, exit != nil)
	tassert(t, err == nil)
}
