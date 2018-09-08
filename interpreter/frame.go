// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"

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

	// isBase specifies whether this is the base frame
	// of the current Eval().
	isBase bool

	// instruction pointer
	ip int
}

func newFrame(fn bc.Func, locals []*bc.Ref, isBase bool) *frame {

	// save these so we don't have to look them up later
	btc := fn.Template().Bytecodes
	pool := fn.Template().Module.Pool

	return &frame{
		fn:     fn,
		locals: locals,
		btc:    btc,
		pool:   pool,

		stack:    make([]g.Value, 0, 10),
		handlers: []bc.ErrorHandler{},
		isBase:   isBase,
		ip:       0,
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

//-------------------------------------------------------------------

func toStr(val g.Value) string {
	s, err := val.ToStr(nil)
	if err != nil {
		panic(err)
	}
	return s.String()
}

func dumpFrames(frames []*frame) {

	fmt.Printf("-----------------------------------------\n")

	for i, f := range frames {
		fmt.Printf("frame %d\n", i)
		dumpFrame(f)
	}
}

func dumpFrame(f *frame) {

	fmt.Printf("    locals:\n")
	for i, r := range f.locals {
		fmt.Printf("        %d: %s\n", i, toStr(r.Val))
	}

	fmt.Printf("    stack:\n")
	for i, v := range f.stack {
		fmt.Printf("        %d: %s\n", i, toStr(v))
	}

	fmt.Printf("    handlers:\n")
	for i, v := range f.handlers {
		fmt.Printf("        %d: %s\n", i, v)
	}

	fmt.Printf("    isBase: %v\n", f.isBase)

	fmt.Printf("    ip: %d\n", f.ip)
}
