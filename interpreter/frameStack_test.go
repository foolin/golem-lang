// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"reflect"
	"testing"

	bc "github.com/mjarmy/golem-lang/core/bytecode"
)

func TestFrameStack(t *testing.T) {

	h1 := bc.ErrorHandler{Catch: bc.TryClause{Begin: 1}, Finally: bc.TryClause{}}
	h2 := bc.ErrorHandler{Catch: bc.TryClause{Begin: 2}, Finally: bc.TryClause{}}
	h3 := bc.ErrorHandler{Catch: bc.TryClause{Begin: 3}, Finally: bc.TryClause{}}
	h4 := bc.ErrorHandler{Catch: bc.TryClause{Begin: 4}, Finally: bc.TryClause{}}

	//---------------------------------------

	fs := &frameStack{frames: []*frame{
		&frame{
			isBase:   true,
			handlers: []bc.ErrorHandler{},
		},
	}}
	//dumpFrames(fs.frames)

	_, ok := fs.popErrorHandler()
	tassert(t, !ok)
	tassert(t, fs.num() == 0)

	//---------------------------------------

	fs = &frameStack{frames: []*frame{
		&frame{
			isBase:   true,
			handlers: []bc.ErrorHandler{h1},
		},
	}}
	//dumpFrames(fs.frames)

	h, _ := fs.popErrorHandler()
	tassert(t, reflect.DeepEqual(h, h1))
	tassert(t, fs.num() == 1)
	tassert(t, reflect.DeepEqual(fs.peek().handlers, []bc.ErrorHandler{}))

	//---------------------------------------

	fs = &frameStack{frames: []*frame{
		&frame{
			isBase:   true,
			handlers: []bc.ErrorHandler{h1, h2},
		},
	}}
	//dumpFrames(fs.frames)

	h, _ = fs.popErrorHandler()
	tassert(t, reflect.DeepEqual(h, h2))
	tassert(t, fs.num() == 1)
	tassert(t, reflect.DeepEqual(fs.peek().handlers, []bc.ErrorHandler{h1}))

	//---------------------------------------

	fs = &frameStack{frames: []*frame{
		&frame{
			isBase:   true,
			handlers: []bc.ErrorHandler{h1, h2},
		},
		&frame{
			isBase:   false,
			handlers: []bc.ErrorHandler{},
		},
	}}
	//dumpFrames(fs.frames)

	h, _ = fs.popErrorHandler()
	tassert(t, reflect.DeepEqual(h, h2))
	tassert(t, fs.num() == 1)
	tassert(t, reflect.DeepEqual(fs.peek().handlers, []bc.ErrorHandler{h1}))

	//---------------------------------------

	fs = &frameStack{frames: []*frame{
		&frame{
			isBase:   true,
			handlers: []bc.ErrorHandler{},
		},
		&frame{
			isBase:   false,
			handlers: []bc.ErrorHandler{h1},
		},
		&frame{
			isBase:   false,
			handlers: []bc.ErrorHandler{},
		},
	}}
	//dumpFrames(fs.frames)

	h, _ = fs.popErrorHandler()
	tassert(t, reflect.DeepEqual(h, h1))
	tassert(t, fs.num() == 2)
	tassert(t, reflect.DeepEqual(fs.peek().handlers, []bc.ErrorHandler{}))

	//---------------------------------------

	fs = &frameStack{frames: []*frame{
		&frame{
			isBase:   true,
			handlers: []bc.ErrorHandler{},
		},
		&frame{
			isBase:   false,
			handlers: []bc.ErrorHandler{},
		},
		&frame{
			isBase:   false,
			handlers: []bc.ErrorHandler{},
		},
	}}
	//dumpFrames(fs.frames)

	_, ok = fs.popErrorHandler()

	tassert(t, !ok)
	tassert(t, fs.num() == 0)

	//---------------------------------------

	fs = &frameStack{frames: []*frame{
		&frame{
			isBase:   true,
			handlers: []bc.ErrorHandler{},
		},
		&frame{
			isBase:   false,
			handlers: []bc.ErrorHandler{h1, h2},
		},
		&frame{
			isBase:   false,
			handlers: []bc.ErrorHandler{h3},
		},
		&frame{
			isBase:   false,
			handlers: []bc.ErrorHandler{h4},
		},
		&frame{
			isBase:   false,
			handlers: []bc.ErrorHandler{},
		},
	}}
	//dumpFrames(fs.frames)

	h, _ = fs.popErrorHandler()
	tassert(t, reflect.DeepEqual(h, h4))
	tassert(t, fs.num() == 4)
	tassert(t, reflect.DeepEqual(fs.peek().handlers, []bc.ErrorHandler{}))

	h, _ = fs.popErrorHandler()
	tassert(t, reflect.DeepEqual(h, h3))
	tassert(t, fs.num() == 3)
	tassert(t, reflect.DeepEqual(fs.peek().handlers, []bc.ErrorHandler{}))

	h, _ = fs.popErrorHandler()
	tassert(t, reflect.DeepEqual(h, h2))
	tassert(t, fs.num() == 2)
	tassert(t, reflect.DeepEqual(fs.peek().handlers, []bc.ErrorHandler{h1}))

	h, _ = fs.popErrorHandler()
	tassert(t, reflect.DeepEqual(h, h1))
	tassert(t, fs.num() == 2)
	tassert(t, reflect.DeepEqual(fs.peek().handlers, []bc.ErrorHandler{}))

	_, ok = fs.popErrorHandler()
	tassert(t, !ok)
	tassert(t, fs.num() == 0)
}
