// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
)

// FuncTemplate represents the information needed to invoke a function
// instance.  Templates are created at compile time, and
// are immutable at run time.
type FuncTemplate struct {
	Module          *Module
	Arity           g.Arity
	OptionalParams  []g.Value
	NumCaptures     int
	NumLocals       int
	Bytecodes       []byte
	LineNumberTable []LineNumberEntry
	ErrorHandlers   []ErrorHandler
}

// LineNumberEntry tracks which sequence of opcodes are on a given source code line
type LineNumberEntry struct {
	Index   int
	LineNum int
}

func (ln LineNumberEntry) String() string {
	return fmt.Sprintf(
		"LineNumberEntry(Index: %d, LineNum: %d)",
		ln.Index, ln.LineNum)
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

// ErrorHandler handles errors that are thrown for a given block of opcodes,
// by providing the instruction pointers for 'catch' and 'finally'.
type ErrorHandler struct {
	Catch   TryClause
	Finally TryClause
}

func (eh ErrorHandler) String() string {
	return fmt.Sprintf(
		"ErrorHandler(Catch: (%s), Finally: (%s)))",
		eh.Catch, eh.Finally)
}

// Try Clause contains the instruction pointers for the begin and end
// of a catch-clause or finally-clause
type TryClause struct {
	Begin int
	End   int
}

func (tc TryClause) String() string {
	return fmt.Sprintf("%d, %d", tc.Begin, tc.End)
}

func (tc TryClause) IsValid() bool {
	return tc.Begin != -1
}
