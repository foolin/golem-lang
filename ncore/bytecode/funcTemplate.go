// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	g "github.com/mjarmy/golem-lang/ncore"
)

// FuncTemplate represents the information needed to invoke a function
// instance.  Templates are created at compile time, and
// are immutable at run time.
type FuncTemplate struct {
	Module            *Module
	Arity             *g.Arity
	NumCaptures       int
	NumLocals         int
	OpCodes           []byte
	LineNumberTable   []LineNumberEntry
	ExceptionHandlers []ExceptionHandler
}

// LineNumberEntry tracks which sequence of opcodes are on a given line
type LineNumberEntry struct {
	Index   int
	LineNum int
}

// ExceptionHandler handles exceptions for a given block of opcodes,
// by providing the instruction pointers for 'catch' and 'finally'
type ExceptionHandler struct {
	Begin   int
	End     int
	Catch   int
	Finally int
}

// LineNumber returns the line number for the opcode at the given instruction pointer
func (t *FuncTemplate) LineNumber(instPtr int) int {

	table := t.LineNumberTable
	n := len(table) - 1

	for i := 0; i < n; i++ {
		if (instPtr >= table[i].Index) && (instPtr < table[i+1].Index) {
			return table[i].LineNum
		}
	}
	return table[n].LineNum
}
