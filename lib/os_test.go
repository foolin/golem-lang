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

func TestOs(t *testing.T) {
	os := NewOsModule()
	exit, err := os.GetContents().GetField(nil, g.NewStr("exit"))
	tassert(t, exit != nil)
	tassert(t, err == nil)

	interpret(`
import os

// lets not actually call 'exit' succesfully :-)

try {
	os.exit('foo')
	assert(false)
} catch e {
	assert(e.kind == 'TypeMismatch')
	assert(e.msg == 'Expected Int')
}

try {
	os.exit(1, 2)
	assert(false)
} catch e {
	assert(e.kind == 'ArityMismatch')
	assert(e.msg == 'Expected at most 1 params, got 2')
}
`)
}
