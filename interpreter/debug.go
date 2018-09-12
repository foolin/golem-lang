// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
)

var debugInterpreter bool

func debugString(s string) {
	if debugInterpreter {
		fmt.Printf(s)
	}
}

func debugVal(val g.Value) string {
	if val == nil {
		return "nil"
	}
	s, err := val.ToStr(nil)
	if err != nil {
		panic(err)
	}
	return s.String()
}

func (r response) debug() {
	if debugInterpreter {
		fmt.Printf("response(%s, %d, %s)\n",
			debugVal(r.result),
			r.resultIp,
			debugVal(r.es))
	}
}

func (f *frame) debug() {

	fmt.Printf("    locals:\n")
	for i, r := range f.locals {
		fmt.Printf("        %d: %s\n", i, debugVal(r.Val))
	}

	fmt.Printf("    stack:\n")
	for i, v := range f.stack {
		fmt.Printf("        %d: %s\n", i, debugVal(v))
	}

	fmt.Printf("    handlers:\n")
	for i, v := range f.handlers {
		fmt.Printf("        %d: %s\n", i, v)
	}

	fmt.Printf("    ip: %d\n", f.ip)
	fmt.Printf("    isBase: %v\n", f.isBase)
	fmt.Printf("    isHandlingError: %v\n", f.isHandlingError)
}

func (fs *frameStack) debug() {

	for i := fs.num() - 1; i >= 0; i-- {
		fmt.Printf("-----------------------------------------\n")
		fmt.Printf("frame %d\n", i)
		fs.get(i).debug()
		fmt.Printf("\n")
	}
}
