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
// The interpreter manages a stack of frames.
type frame struct {
	fn     bc.Func
	locals []*bc.Ref

	// whether this is the last frame that will be processed during the current Eval()
	isLast bool

	stack        []g.Value
	handlerStack []bc.ErrorHandler

	// current instruction point
	ip int

	// useful things from the function template
	btc         []byte
	pool        *bc.Pool
	handlerPool []bc.ErrorHandler
}

func newFrame(fn bc.Func, locals []*bc.Ref, isLast bool) *frame {

	// save this stuff so we don't have to look it up later
	btc := fn.Template().Bytecodes
	pool := fn.Template().Module.Pool
	handlerPool := fn.Template().ErrorHandlers

	return &frame{
		fn:     fn,
		locals: locals,
		isLast: isLast,

		stack:        []g.Value{},
		handlerStack: []bc.ErrorHandler{},
		ip:           0,

		btc:         btc,
		pool:        pool,
		handlerPool: handlerPool,
	}
}

func (f *frame) numHandlers() int {
	return len(f.handlerStack)
}

func (f *frame) pushHandler(h bc.ErrorHandler) {
	f.handlerStack = append(f.handlerStack, h)
}

func (f *frame) popHandler() bc.ErrorHandler {
	n := f.numHandlers() - 1
	h := f.handlerStack[n]
	f.handlerStack = f.handlerStack[:n]
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

	fmt.Printf("    handlerStack:\n")
	for i, v := range f.handlerStack {
		fmt.Printf("        %d: %s\n", i, v)
	}

	fmt.Printf("    ip: %d\n", f.ip)
}
