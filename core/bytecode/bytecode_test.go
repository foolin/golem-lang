// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"testing"

	g "github.com/mjarmy/golem-lang/core"
)

func push(btc []byte, bytes ...byte) []byte {
	return append(btc, bytes...)
}

func TestBytecodes(t *testing.T) {

	btc := []byte{}

	btc = push(btc, LoadNull)

	a, b := EncodeParam(0)
	btc = push(btc, LoadConst, a, b)

	a, b = EncodeParam(100)
	btc = push(btc, LoadLocal, a, b)

	a, b, c, d := EncodeWideParams(3, 1234)
	btc = push(btc, InvokeField, a, b, c, d)

	g.Tassert(t, 12 == len(btc))

	g.Tassert(t, LoadNull == btc[0])

	g.Tassert(t, LoadConst == btc[1])
	g.Tassert(t, 0 == DecodeParam(btc, 1))

	g.Tassert(t, LoadLocal == btc[4])
	g.Tassert(t, 100 == DecodeParam(btc, 4))

	g.Tassert(t, InvokeField == btc[7])
	p, q := DecodeWideParams(btc, 7)
	g.Tassert(t, 3 == p && 1234 == q)

	g.Tassert(t, FmtBytecode(btc, 0) == "0: LoadNull")
	g.Tassert(t, FmtBytecode(btc, 1) == "1: LoadConst    0 0 (0)")
	g.Tassert(t, FmtBytecode(btc, 4) == "4: LoadLocal    0 100 (100)")
	g.Tassert(t, FmtBytecode(btc, 7) == "7: InvokeField  0 3 (3), 4 210 (1234)")
}
