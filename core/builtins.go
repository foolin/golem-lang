// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

type (
	// BuiltinManager manages the built-in functions (and other built-in values)
	// for a given instance of the Interpreter
	BuiltinManager interface {
		Builtins() []Value
		Contains(s string) bool
		IndexOf(s string) int
	}

	// BuiltinEntry is an entry in a BuiltinManager
	BuiltinEntry struct {
		Name  string
		Value Value
	}

	builtinManager struct {
		values []Value
		lookup map[string]int
	}
)

// NewBuiltinManager creates a new BuiltinManager
func NewBuiltinManager(entries []*BuiltinEntry) BuiltinManager {
	values := make([]Value, len(entries))
	lookup := make(map[string]int)
	for i, e := range entries {
		values[i] = e.Value
		lookup[e.Name] = i
	}
	return &builtinManager{values, lookup}
}

func (b *builtinManager) Builtins() []Value {
	return b.values
}

func (b *builtinManager) Contains(s string) bool {
	_, ok := b.lookup[s]
	return ok

}

func (b *builtinManager) IndexOf(s string) int {
	index, ok := b.lookup[s]
	if !ok {
		panic("unknown builtin")
	}
	return index
}

//-----------------------------------------------------------------

/*doc

## Standard Builtins

Golem has a collection of standard builtin functions that provide
various kinds of important functionality. All of the standard builtins
are "pure" functions that do not do any form of I/O. As such they
are suitable for use in sandboxed environments.

* [`arity()`](#arity)
* [`assert()`](#assert)
* [`chan()`](#chan)
* [`fields()`](#fields)
* [`freeze()`](#freeze)
* [`frozen()`](#frozen)
* [`has()`](#has)
* [`iter()`](#iter)
* [`len()`](#len)
* [`merge()`](#merge)
* [`range()`](#range)
* [`str()`](#str)
* [`type()`](#type)

*/

// StandardBuiltins containts the built-ins that are
// pure functions.  These functions do not do any form of I/O.
var StandardBuiltins = []*BuiltinEntry{
	{"arity", BuiltinArity},
	{"assert", BuiltinAssert},
	{"chan", BuiltinChan},
	{"fields", BuiltinFields},
	{"freeze", BuiltinFreeze},
	{"frozen", BuiltinFrozen},
	{"has", BuiltinHas},
	{"iter", BuiltinIter},
	{"len", BuiltinLen},
	{"merge", BuiltinMerge},
	{"range", BuiltinRange},
	{"str", BuiltinStr},
	{"type", BuiltinType},
}

//-----------------------------------------------------------------

/*doc
### `arity`

`arity` returns a Struct describing the [arity](https://en.wikipedia.org/wiki/Arity) of a Func.
A func's arity type is always either "Fixed", "Variadic", or "Multiple".

* signature: `arity(f <Func>) <Struct>`
* example:

```
println(arity(len))
println(arity(println))
println(arity(range))
```

*/

// BuiltinArity returns a Struct describing the arity of a Func.
var BuiltinArity = NewFixedNativeFunc(
	[]Type{FuncType},
	false,
	func(ev Eval, params []Value) (Value, Error) {

		fn := params[0].(Func)
		a := fn.Arity()
		k, err := NewStr(a.Kind.String())
		if err != nil {
			return nil, err
		}
		r := NewInt(int64(a.Required))

		fields := map[string]Field{
			"kind":     NewReadonlyField(k),
			"required": NewReadonlyField(r),
		}

		if a.Kind == MultipleArity {
			o := NewInt(int64(a.Optional))
			fields["optional"] = NewReadonlyField(o)
		}

		return NewFrozenFieldStruct(fields)
	})

/*doc
### `assert`

`assert` accepts a single boolean value, and throws an error
if the value is not equal to `true`.  `assert` returns `true`
if it does not throw an error.

* signature: `assert(b <Bool>) <Bool>`
* example: `assert(0 < 1)`

*/

// BuiltinAssert asserts that a single Bool is True
var BuiltinAssert = NewFixedNativeFunc(
	[]Type{BoolType}, false,
	func(ev Eval, params []Value) (Value, Error) {
		b := params[0].(Bool)
		if b.BoolVal() {
			return True, nil
		}
		return nil, AssertionFailed()
	})

/*doc
### `chan`

`chan` creates a [channel](chan.html) of values.  `chan` has a single optional size parameter that
defaults to 0.  If size is 0, an unbuffered channel will be created.
If the size is greater than 0, then a buffered channel of that size will be created.

* signature: `chan(size = 0 <Int>) <Chan>`
* example: `let ch = chan()`

*/

// BuiltinChan creates a new Chan.
var BuiltinChan = NewMultipleNativeFunc(
	[]Type{},
	[]Type{IntType},
	false,
	func(ev Eval, params []Value) (Value, Error) {

		if len(params) == 0 {
			return NewChan(), nil
		}

		size := params[0].(Int)
		return NewBufferedChan(int(size.IntVal())), nil
	})

/*doc
### `fields`

`fields` returns a Set of the names of a value's fields.

* signature: `fields(value <Value>) <Set>`
* example:

```
println(fields([]))
```

*/

var BuiltinFields = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {

		fields, err := params[0].FieldNames()
		if err != nil {
			return nil, err
		}

		entries := make([]Value, len(fields))
		for i, k := range fields {
			entry, err := NewStr(k)
			if err != nil {
				return nil, err
			}
			entries[i] = entry
		}
		return NewSet(ev, entries)
	})

/*doc
### `freeze`

`freeze` freezes a value, if it is not already frozen.  Its OK to call `freeze`
on values that are already frozen.  The value is returned after it is frozen.

* signature: `freeze(value <Value>) <Freeze>`
* example: `freeze([1, 2])`

*/

// BuiltinFreeze freezes a single Value.
var BuiltinFreeze = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		return params[0].Freeze(ev)
	})

/*doc
### `frozen`

`frozen` returns whether or not a value is frozen.

* signature: `frozen(value <Value>) <Bool>`
* example:

```
println(frozen('a'))
println(frozen([3, 4]))
```

*/

// BuiltinFrozen returns whether a single Value is Frozen.
var BuiltinFrozen = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		return params[0].Frozen(ev)
	})

/*doc
### `has`

`has` returns whether a value has a field with a given name.

* signature: `has(name <Str>) <Bool>`
* example:

```
let a = [1, 2]
println(has(a, 'add'))
```

*/

var BuiltinHas = NewFixedNativeFunc(
	[]Type{AnyType, StrType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		field := params[1].(Str)

		b, err := params[0].HasField(field.String())
		if err != nil {
			return nil, err
		}
		return NewBool(b), nil
	})

/*doc
### `iter`

`iter` returns an iterator for an Iterable value.  Str, List, Range, Dict,
and Set are iterable.

An iterator is a Struct that has two fields:

* A `next()` function that
returns whether there are any more values in the iterator,
and advances the iterator forwards if there is another value.

* A `get()` function that returns the currently available value.

By convention, a new iterator has to have `next()` called on it to advance
to the first available value. Calling `get()` before the first call to `next()`
throws an error.

* signature: `iter(value <Iterable>) <Struct>`
* `next` signature: `next() <Bool>`
* `get` signature: `get() <Value>`
* example:

```
let a = [1, 2, 3]
let itr = iter(a)
while itr.next() {
	println(itr.get())
}
```

*/

// BuiltinIter returns an iterator for an Iterable value.
var BuiltinIter = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {

		if params[0].Type() == NullType {
			return nil, NullValueError()
		}

		if ibl, ok := params[0].(Iterable); ok {
			return ibl.NewIterator(ev)
		}
		return nil, IterableMismatch(params[0].Type())
	})

/*doc
### `len`

`len` returns the length of a value that has a length.  Str, List, Tuple, Range, Dict,
and Set have a length

* signature: `len(value <Lenable>) <Int>`
* example: `println(len('abc'))`

*/

// BuiltinLen returns the length of a single Lenable
var BuiltinLen = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {

		if params[0].Type() == NullType {
			return nil, NullValueError()
		}

		if ln, ok := params[0].(Lenable); ok {
			return ln.Len(ev)
		}
		return nil, LenableMismatch(params[0].Type())
	})

/*doc
### `merge`

`merge` combines an arbitrary number of existing structs into a new struct.  If
there are any duplicated keys in the structs passed in to 'merge()', then the
value associated with the first such key is used.

* signature: `merge(structs... <Struct>) <Struct>`
* example:

```
let a = struct { x: 1, y: 2 }
let b = struct { y: 3, z: 4 }
let c = merge(a, b)

println(a)
println(b)
println(c)

a.x = 10

println(a)
println(b)
println(c) // x is changed here too!
```

*/

// BuiltinMerge merges structs together.
var BuiltinMerge = NewVariadicNativeFunc(
	[]Type{StructType, StructType},
	StructType,
	false,
	func(ev Eval, params []Value) (Value, Error) {
		structs := make([]Struct, len(params))
		for i, v := range params {
			structs[i] = v.(Struct)
		}
		return MergeStructs(structs)
	})

/*doc
### `range`

`range` creates a Range, starting at "from" (inclusive) and going until
"to" (exclusive).

The optional "step" parameter, which defaults to 1,
specifies the distance between succesive integers in the range.  You can
create a "backwards" range by specify a negative step value, and a "from"
that is less than "to".

* signature: `range(from <Int>, to <Int>, step = 1 <Int>) <Range>`
* example:

```
for i in range(0, 5) {
	println(i)
}
```

*/

// BuiltinRange creates a new Range
var BuiltinRange = NewMultipleNativeFunc(
	[]Type{IntType, IntType},
	[]Type{IntType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		from := params[0].(Int)
		to := params[1].(Int)
		step := One
		if len(params) == 3 {
			step = params[2].(Int)

		}
		return NewRange(from.IntVal(), to.IntVal(), step.IntVal())
	})

/*doc
### `str`

`str` returns a Str representation of a value.

* signature: `str(value <Value>) <Str>`
* example: `println(str([null, true, 1, 'abc']))`

*/

// BuiltinStr converts a single value to a Str
var BuiltinStr = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		return params[0].ToStr(ev)
	})

/*doc
### `type`

`type` returns the type of a value.

* signature: `type(value <Value>) <Str>`
* example:

```
println(type(1.23))
let a = [null, true, 1, 'xyz']
println(a.map(type))
```

*/

// BuiltinType returns the Str representation of the Type of a single Value
var BuiltinType = NewFixedNativeFunc(
	[]Type{AnyType},
	true,
	func(ev Eval, params []Value) (Value, Error) {

		return NewStr(params[0].Type().String())
	})

//-----------------------------------------------------------------

/*doc

## Unsandboxed Builtins

Golem also has "unsandboxed" builtins.  These functions
perform I/O, so they should not be included in sandboxed Golem
environments.

* [`print()`](#print)
* [`println()`](#println)

*/

// UnsandboxedBuiltins are builtins that are not pure functions
var UnsandboxedBuiltins = []*BuiltinEntry{
	{"print", BuiltinPrint},
	{"println", BuiltinPrintln},
}

/*doc
### `print`

`print` prints a sequence of values to STDOUT.

* signature: `print(values... <Value>) <Null>`

*/

// BuiltinPrint prints to stdout.
var BuiltinPrint = NewVariadicNativeFunc(
	[]Type{}, AnyType, true,
	func(ev Eval, params []Value) (Value, Error) {
		for _, v := range params {
			s, err := v.ToStr(ev)
			if err != nil {
				return nil, err
			}
			fmt.Print(s.String())
		}

		return Null, nil
	})

/*doc
### `println`

`println` prints a sequence of values to STDOUT, followed by a linefeed.

* signature: `println(values... <Value>) <Null>`

*/

// BuiltinPrintln prints to stdout.
var BuiltinPrintln = NewVariadicNativeFunc(
	[]Type{}, AnyType, true,
	func(ev Eval, params []Value) (Value, Error) {
		for _, v := range params {
			s, err := v.ToStr(ev)
			if err != nil {
				return nil, err
			}
			fmt.Print(s.String())
		}
		fmt.Println()

		return Null, nil
	})
