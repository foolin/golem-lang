// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	//"sync"

	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
)

// A Library is a collection of core.Modules.
type Library interface {
	// Get returns a core.Module (possibly  by resolving a bytecode.Module,
	// compiling it, and initializing it)
	Get(itp *Interpreter, name string) (g.Module, error)
}

// A Resolver can resolve a bytecode.Module by name
type Resolver func(name string) (*bc.Module, error)

type library struct {
	moduleMap map[string]g.Module
	resolver  Resolver
	//mx        *sync.Mutex
}

// NewLibrary creates a new Library.  If resolver is nil, then
// no new modules will be added to the library at run time.
func NewLibrary(modules []g.Module, resolver Resolver) Library {

	var moduleMap = map[string]g.Module{}
	for _, m := range modules {
		moduleMap[m.Name()] = m
	}

	return &library{moduleMap, resolver /*, &sync.Mutex{}*/}
}

func (lib *library) Get(itp *Interpreter, name string) (g.Module, error) {

	// TODO we can end up in a race if a goroutine tries to compile a module
	//lib.mx.Lock()
	//defer lib.mx.Unlock()

	// check if we already know about the module
	m, ok := lib.moduleMap[name]
	if ok {
		return m, nil
	}

	// if we don't have a resolver, the module is undefined
	if lib.resolver == nil {
		return nil, g.UndefinedModule(name)
	}

	// try to resolve the module
	mb, err := lib.resolver(name)
	if err != nil {
		return nil, err
	}

	// initialze the new module
	_, err = itp.EvalModule(mb)
	if err != nil {
		return nil, err
	}

	// add the new module to the library, and return
	lib.moduleMap[name] = mb
	return mb, nil
}
