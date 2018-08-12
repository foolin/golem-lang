// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this code code is governed by a MIT-style
// license that can be found in the LICENSE file.

package analyzer

import (
	"testing"
	//"github.com/mjarmy/golem-lang/ast"
)

func TestStruct(t *testing.T) {

	mod := newModule("this")
	errors := NewAnalyzer(mod).Analyze()
	fail(t, errors, "['this' outside of struct, at foo.glm:1:1]")

	code := `
struct{ }
`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:0)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([],StructScope defs:{})
`)

	code = `
struct{ a: 1 }
`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:0)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([a],StructScope defs:{})
.   .   .   .   BasicExpr(Int,"1")
`)

	code = `
struct{ a: this }
`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:1)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([a],StructScope defs:{this: v(0: this,0,true,false)})
.   .   .   .   ThisExpr(v(0: this,0,true,false))
`)

	code = `
struct{ a: struct { b: this } }
`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:1)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([a],StructScope defs:{})
.   .   .   .   StructExpr([b],StructScope defs:{this: v(0: this,0,true,false)})
.   .   .   .   .   ThisExpr(v(0: this,0,true,false))
`)

	code = `
struct{ a: struct { b: 1 }, c: this.a }
`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:1)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([a, c],StructScope defs:{this: v(0: this,0,true,false)})
.   .   .   .   StructExpr([b],StructScope defs:{})
.   .   .   .   .   BasicExpr(Int,"1")
.   .   .   .   FieldExpr(a)
.   .   .   .   .   ThisExpr(v(0: this,0,true,false))
`)

	code = `
struct{ a: struct { b: this }, c: this }
`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([a, c],StructScope defs:{this: v(1: this,1,true,false)})
.   .   .   .   StructExpr([b],StructScope defs:{this: v(0: this,0,true,false)})
.   .   .   .   .   ThisExpr(v(0: this,0,true,false))
.   .   .   .   ThisExpr(v(1: this,1,true,false))
`)

	code = `
struct{ a: this, b: struct { c: this } }
`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()
	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{})
.   .   ExprStmt
.   .   .   StructExpr([a, b],StructScope defs:{this: v(0: this,0,true,false)})
.   .   .   .   ThisExpr(v(0: this,0,true,false))
.   .   .   .   StructExpr([c],StructScope defs:{this: v(1: this,1,true,false)})
.   .   .   .   .   ThisExpr(v(1: this,1,true,false))
`)

	code = `
let a = struct {
	x: 8,
	y: 5,
	plus:  fn() { return this.x + this.y; },
	minus: fn() { return this.x - this.y; }
}
let b = a.plus()
let c = a.minus()
`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:4)
.   BlockNode(Scope defs:{a: v(3: a,1,false,false), b: v(4: b,2,false,false), c: v(5: c,3,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(a,v(3: a,1,false,false))
.   .   .   StructExpr([x, y, plus, minus],StructScope defs:{this: v(0: this,0,true,false)})
.   .   .   .   BasicExpr(Int,"8")
.   .   .   .   BasicExpr(Int,"5")
.   .   .   .   FnExpr(FuncScope defs:{} captures:{this: (parent: v(0: this,0,true,false), child v(1: this,0,true,true))} numLocals:0)
.   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   ReturnStmt
.   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   FieldExpr(x)
.   .   .   .   .   .   .   .   .   ThisExpr(v(1: this,0,true,true))
.   .   .   .   .   .   .   .   FieldExpr(y)
.   .   .   .   .   .   .   .   .   ThisExpr(v(1: this,0,true,true))
.   .   .   .   FnExpr(FuncScope defs:{} captures:{this: (parent: v(0: this,0,true,false), child v(2: this,0,true,true))} numLocals:0)
.   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   ReturnStmt
.   .   .   .   .   .   .   BinaryExpr("-")
.   .   .   .   .   .   .   .   FieldExpr(x)
.   .   .   .   .   .   .   .   .   ThisExpr(v(2: this,0,true,true))
.   .   .   .   .   .   .   .   FieldExpr(y)
.   .   .   .   .   .   .   .   .   ThisExpr(v(2: this,0,true,true))
.   .   LetStmt
.   .   .   IdentExpr(b,v(4: b,2,false,false))
.   .   .   InvokeExpr
.   .   .   .   FieldExpr(plus)
.   .   .   .   .   IdentExpr(a,v(3: a,1,false,false))
.   .   LetStmt
.   .   .   IdentExpr(c,v(5: c,3,false,false))
.   .   .   InvokeExpr
.   .   .   .   FieldExpr(minus)
.   .   .   .   .   IdentExpr(a,v(3: a,1,false,false))
`)

	code = `
let x = 1;
let y = 2;
struct {
	a: prop { || => x },
	b: prop { || => y, |v| => y = v }
}
`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:2)
.   BlockNode(Scope defs:{x: v(0: x,0,false,false), y: v(1: y,1,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(x,v(0: x,0,false,false))
.   .   .   BasicExpr(Int,"1")
.   .   LetStmt
.   .   .   IdentExpr(y,v(1: y,1,false,false))
.   .   .   BasicExpr(Int,"2")
.   .   ExprStmt
.   .   .   StructExpr([a, b],StructScope defs:{})
.   .   .   .   PropNode
.   .   .   .   .   FnExpr(FuncScope defs:{} captures:{x: (parent: v(0: x,0,false,false), child v(2: x,0,false,true))} numLocals:0)
.   .   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   .   ExprStmt
.   .   .   .   .   .   .   .   IdentExpr(x,v(2: x,0,false,true))
.   .   .   .   PropNode
.   .   .   .   .   FnExpr(FuncScope defs:{} captures:{y: (parent: v(1: y,1,false,false), child v(3: y,0,false,true))} numLocals:0)
.   .   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   .   ExprStmt
.   .   .   .   .   .   .   .   IdentExpr(y,v(3: y,0,false,true))
.   .   .   .   .   FnExpr(FuncScope defs:{v: v(4: v,0,false,false)} captures:{y: (parent: v(1: y,1,false,false), child v(5: y,0,false,true))} numLocals:1)
.   .   .   .   .   .   IdentExpr(v,v(4: v,0,false,false))
.   .   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   .   ExprStmt
.   .   .   .   .   .   .   .   AssignmentExpr
.   .   .   .   .   .   .   .   .   IdentExpr(y,v(5: y,0,false,true))
.   .   .   .   .   .   .   .   .   IdentExpr(v,v(4: v,0,false,false))
`)

	code = `
let x = 1
let u = struct {
	a: prop { || => x },
	b: || => this.a
}
println(u)
`
	mod = newModule(code)
	errors = NewAnalyzer(mod).Analyze()

	//println(code)
	//println(ast.Dump(anl.Module()))
	//println(errors)

	ok(t, mod, errors, `
FnExpr(FuncScope defs:{} captures:{} numLocals:3)
.   BlockNode(Scope defs:{u: v(4: u,2,false,false), x: v(0: x,0,false,false)})
.   .   LetStmt
.   .   .   IdentExpr(x,v(0: x,0,false,false))
.   .   .   BasicExpr(Int,"1")
.   .   LetStmt
.   .   .   IdentExpr(u,v(4: u,2,false,false))
.   .   .   StructExpr([a, b],StructScope defs:{this: v(2: this,1,true,false)})
.   .   .   .   PropNode
.   .   .   .   .   FnExpr(FuncScope defs:{} captures:{x: (parent: v(0: x,0,false,false), child v(1: x,0,false,true))} numLocals:0)
.   .   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   .   ExprStmt
.   .   .   .   .   .   .   .   IdentExpr(x,v(1: x,0,false,true))
.   .   .   .   FnExpr(FuncScope defs:{} captures:{this: (parent: v(2: this,1,true,false), child v(3: this,0,true,true))} numLocals:0)
.   .   .   .   .   BlockNode(Scope defs:{})
.   .   .   .   .   .   ExprStmt
.   .   .   .   .   .   .   FieldExpr(a)
.   .   .   .   .   .   .   .   ThisExpr(v(3: this,0,true,true))
.   .   ExprStmt
.   .   .   InvokeExpr
.   .   .   .   BuiltinExpr("println")
.   .   .   .   IdentExpr(u,v(4: u,2,false,false))
`)
}
