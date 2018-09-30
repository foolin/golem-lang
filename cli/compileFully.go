// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

//import (
//	"github.com/mjarmy/golem-lang/compiler"
//	bc "github.com/mjarmy/golem-lang/core/bytecode"
//	"github.com/mjarmy/golem-lang/scanner"
//)
//
//// SourceResolver resolves a module name into a Source
//type SourceResolver func(sourceName string) (*scanner.Source, error)
//
//// CompileSourceFully compiles all of the Modules needed to run the program
//// that is defined in the provided Source.
//func CompileSourceFully(
//	builtinMgr compiler.BuiltinManager,
//	resolver SourceResolver,
//	source *scanner.Source) ([]*bc.Module, []error) {
//
//	sources := []*scanner.Source{source}
//	sourceSet := map[string]bool{source.Name: true}
//	result := []*bc.Module{}
//
//	idx := 0
//	for idx < len(sources) {
//
//		// compile
//		astMod, mod, errs := compiler.CompileSource(builtinMgr, sources[idx])
//		if errs != nil {
//			return nil, errs
//		}
//		result = append(result, mod)
//
//		// add imported
//		for _, impName := range astMod.Imports() {
//			if _, ok := sourceSet[impName]; !ok {
//				impSrc, err := resolver(impName)
//				if err != nil {
//					return nil, []error{err}
//				}
//				sources = append(sources, impSrc)
//				sourceSet[impName] = true
//			}
//		}
//
//		// done
//		idx++
//	}
//
//	return result, nil
//}
//
////// InitModules initializes each of the Modules.  Note that the modules
////// are initialized in reverse order.
////func (itp *Interpreter) InitModules() ([]g.Value, ErrorStruct) {
////
////	values := []g.Value{}
////	for i := len(itp.modules) - 1; i >= 0; i-- {
////		mod := itp.modules[i]
////
////		// the 'init' function is always the first template in the pool
////		initTpl := mod.Pool.Templates[0]
////
////		// create empty locals
////		mod.Refs = newLocals(initTpl.NumLocals, nil)
////
////		// make init function from template
////		initFn := bc.NewBytecodeFunc(initTpl)
////
////		// invoke the "init" function
////		itp.frameStack.push(newFrame(initFn, mod.Refs, true))
////		val, es := itp.eval()
////		if es != nil {
////			return nil, es
////		}
////
////		// prepend the value so that the values will be in the same order as itp.modules
////		values = append([]g.Value{val}, values...)
////	}
////	return values, nil
////}
