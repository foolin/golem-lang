// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
)

// A frame is an execution environment for a function.
type frame struct {
	fn     bc.Func
	locals []*bc.Ref
	btc    []byte
	pool   *bc.Pool

	stack    []g.Value
	handlers []bc.ErrorHandler
	ip       int // instruction pointer

	// isBase specifies whether this is the base frame
	// of the current Eval().
	isBase bool

	// isHandlingError specifies whether this frame is
	// currently handling and error.
	isHandlingError bool
}

func newFrame(fn bc.Func, locals []*bc.Ref, isBase bool) *frame {

	return &frame{
		fn:     fn,
		locals: locals,
		// save these so we don't have to look them up later
		btc:  fn.Template().Bytecodes,
		pool: fn.Template().Module.Pool,

		stack:    make([]g.Value, 0, 10),
		handlers: []bc.ErrorHandler{},
		ip:       0,

		isBase:          isBase,
		isHandlingError: false,
	}
}

func (f *frame) numHandlers() int {
	return len(f.handlers)
}

func (f *frame) pushHandler(h bc.ErrorHandler) {
	f.handlers = append(f.handlers, h)
}

func (f *frame) popHandler() bc.ErrorHandler {
	n := f.numHandlers() - 1
	h := f.handlers[n]
	f.handlers = f.handlers[:n]
	return h
}
