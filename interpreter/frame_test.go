// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"reflect"
	"testing"

	bc "github.com/mjarmy/golem-lang/core/bytecode"
)

func TestFrames(t *testing.T) {

	h1 := bc.ErrorHandler{CatchBegin: 1}
	h2 := bc.ErrorHandler{CatchBegin: 2}

	//---------------------------------------

	itp := &Interpreter{frames: []*frame{
		&frame{
			isBase:   true,
			handlers: []bc.ErrorHandler{},
		},
	}}
	dumpFrames(itp.frames)

	_, ok := itp.popErrorHandler()

	tassert(t, !ok)
	tassert(t, itp.numFrames() == 0)

	//---------------------------------------

	itp = &Interpreter{frames: []*frame{
		&frame{
			isBase:   true,
			handlers: []bc.ErrorHandler{h1},
		},
	}}
	dumpFrames(itp.frames)

	h, ok := itp.popErrorHandler()

	tassert(t, reflect.DeepEqual(h, h1))
	tassert(t, itp.numFrames() == 1)
	tassert(t, reflect.DeepEqual(itp.peekFrame().handlers, []bc.ErrorHandler{}))

	//---------------------------------------

	itp = &Interpreter{frames: []*frame{
		&frame{
			isBase:   true,
			handlers: []bc.ErrorHandler{h1, h2},
		},
	}}
	dumpFrames(itp.frames)

	h, ok = itp.popErrorHandler()

	tassert(t, reflect.DeepEqual(h, h2))
	tassert(t, itp.numFrames() == 1)
	tassert(t, reflect.DeepEqual(itp.peekFrame().handlers, []bc.ErrorHandler{h1}))

	//---------------------------------------

	itp = &Interpreter{frames: []*frame{
		&frame{
			isBase:   true,
			handlers: []bc.ErrorHandler{h1, h2},
		},
		&frame{
			isBase:   false,
			handlers: []bc.ErrorHandler{},
		},
	}}
	dumpFrames(itp.frames)

	h, ok = itp.popErrorHandler()

	tassert(t, reflect.DeepEqual(h, h2))
	tassert(t, itp.numFrames() == 1)
	tassert(t, reflect.DeepEqual(itp.peekFrame().handlers, []bc.ErrorHandler{h1}))

	//---------------------------------------

	itp = &Interpreter{frames: []*frame{
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
	dumpFrames(itp.frames)

	h, ok = itp.popErrorHandler()

	tassert(t, reflect.DeepEqual(h, h1))
	tassert(t, itp.numFrames() == 2)
	tassert(t, reflect.DeepEqual(itp.peekFrame().handlers, []bc.ErrorHandler{}))

	//---------------------------------------

	itp = &Interpreter{frames: []*frame{
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
	dumpFrames(itp.frames)

	_, ok = itp.popErrorHandler()

	tassert(t, !ok)
	tassert(t, itp.numFrames() == 0)
}
