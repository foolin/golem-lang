// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
)

// frameStack is a stack of frames
type frameStack struct {
	frames []*frame
}

func newFrameStack() *frameStack {
	return &frameStack{frames: []*frame{}}
}

func (fs *frameStack) num() int {
	return len(fs.frames)
}

func (fs *frameStack) peek() *frame {
	return fs.frames[fs.num()-1]
}

func (fs *frameStack) get(idx int) *frame {
	return fs.frames[idx]
}

func (fs *frameStack) push(f *frame) {
	fs.frames = append(fs.frames, f)
}

func (fs *frameStack) pop() {
	fs.frames = fs.frames[:fs.num()-1]
}

func (fs *frameStack) popErrorHandler() (bc.ErrorHandler, bool) {

	f := fs.peek()
	if f.numHandlers() > 0 {
		return f.popHandler(), true
	}

	for !f.isBase {
		fs.pop()

		f = fs.peek()
		if f.numHandlers() > 0 {
			return f.popHandler(), true
		}
	}

	fs.pop()
	return bc.ErrorHandler{}, false
}

func (fs *frameStack) stackTrace() []string {

	stack := []string{}

	for i := fs.num() - 1; i >= 0; i-- {
		f := fs.get(i)
		tpl := f.fn.Template()
		lineNum := tpl.LineNumber(f.ip)
		stack = append(stack, fmt.Sprintf("    at %s:%d", tpl.Module.Path, lineNum))
	}

	return stack
}

//func (fs *frameStack) dump() {
//
//	for i := fs.num() - 1; i >= 0; i-- {
//		fmt.Printf("-----------------------------------------\n")
//		fmt.Printf("frame %d\n", i)
//		fs.get(i).dump()
//		fmt.Printf("\n")
//	}
//}
