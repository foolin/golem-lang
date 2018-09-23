// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"bytes"
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
)

type (
	// ErrorStruct is a core.Struct that is also a core.Error.
	ErrorStruct interface {
		g.Struct
		String() string
		Error() string
		StackTrace() []string
	}

	errorStruct struct {
		g.Struct
		err        g.Error
		stackTrace []string
	}
)

func newErrorStruct(err g.Error, stackTrace []string) ErrorStruct {

	// make List-of-Str
	vals := make([]g.Value, len(stackTrace))
	for i, s := range stackTrace {
		vals[i] = g.MustStr(s)
	}
	list, e := g.NewList(vals).Freeze(nil)
	g.Assert(e == nil)

	stc, e := g.NewFrozenFieldStruct(
		map[string]g.Field{
			"error":      g.NewReadonlyField(g.MustStr(err.Error())),
			"stackTrace": g.NewReadonlyField(list),
			// TODO $toStr for convenience when printing stack trace
		})
	g.Assert(e == nil)

	return &errorStruct{stc, err, stackTrace}
}

func (e *errorStruct) Error() string {
	return e.err.Error()
}

func (e *errorStruct) StackTrace() []string {
	return e.stackTrace
}

func (e *errorStruct) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Error: %s\n", e.Error()))
	for _, s := range e.StackTrace() {
		buf.WriteString(fmt.Sprintf("%s\n", s))
	}
	return buf.String()
}
