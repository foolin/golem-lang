// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

// A Builtin is a built-in value
type Builtin struct {
	Name  string
	Value Value
}

/*doc

## Sandbox Builtins

Golem has a collection of standard builtin functions that provide
various kinds of important functionality. All of the sandbox builtins
are "pure" functions that do not do any form of I/O. As such they
are suitable for use in sandboxed environments.

* [`arity()`](#arity)
* [`assert()`](#assert)
* [`chan()`](#chan)
* [`fields()`](#fields)
* [`freeze()`](#freeze)
* [`frozen()`](#frozen)
* [`has()`](#has)
* [`hashCode()`](#hashcode)
* [`iter()`](#iter)
* [`len()`](#len)
* [`merge()`](#merge)
* [`range()`](#range)
* [`stream()`](#stream)
* [`str()`](#str)
* [`type()`](#type)

*/

// SandboxBuiltins contains the built-ins that are
// pure functions.  These functions do not do any form of I/O.
var SandboxBuiltins = []*Builtin{
	{"arity", BuiltinArity},
	{"assert", BuiltinAssert},
	{"chan", BuiltinChan},
	{"fields", BuiltinFields},
	{"freeze", BuiltinFreeze},
	{"frozen", BuiltinFrozen},
	{"has", BuiltinHas},
	{"hashCode", BuiltinHashCode},
	{"iter", BuiltinIter},
	{"len", BuiltinLen},
	{"merge", BuiltinMerge},
	{"range", BuiltinRange},
	{"str", BuiltinStr},
	{"stream", BuiltinStream},
	{"type", BuiltinType},
}

//-----------------------------------------------------------------

/*doc
### `arity`

`arity` returns a [Struct](struct.html) describing the [arity](https://en.wikipedia.org/wiki/Arity) of a Func.
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

		return NewFrozenStruct(fields)
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
		return NewBufferedChan(int(size.ToInt())), nil
	})

/*doc
### `fields`

`fields` returns a [Set](set.html) of the names of a value's fields.

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

* signature: `has(val <Value>, name <Str>) <Bool>`
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
### `hashCode`

`hashCode` returns the hash code of a [`hashable`](interfaces.html#hashable) value.

* signature: `hashCode(val <Value>) <Int>`
* example:

```
println(hashCode('abc'))
```

*/

var BuiltinHashCode = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		return params[0].HashCode(ev)
	})

/*doc
### `iter`

`iter` returns an "iterator" [Struct](struct.html) for an [iterable](interfaces.html#iterable)
value.

A new iterator must have `next()` called on it to advance
to the first available value. Calling `get()` before the first call to `next()`
throws an error.

* signature: `iter(itr <Iterable>) <Struct>`
* example:

```
let a = [1, 2, 3]
let itr = iter(a)
while itr.next() {
	println(itr.get())
}
```

An iterator struct has the following fields.

* `next()` returns whether there are any more values in the iterator,
and advances the iterator forwards if there is another value.

	* signature: `next() <Bool>`

* `get()` returns the currently available value.

	* signature: `get() <Value>`

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

`len` returns the length of a value that has a [length](interfaces.html#lenable).

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

`merge` combines an arbitrary number of existing [structs](Struct.html) into a new struct.  If
there are any duplicated keys in the structs passed in to 'merge()', then the
value associated with the *last* such key is used.

* signature: `merge(structs... <Struct>) <Struct>`
* example:

```
let a = struct { x: 1, y: 2 }
let b = struct { y: 3, z: 4 }
let c = merge(a, b)

println('a: ', a)
println('b: ', b)
println('c: ', c)

a.x = 10

println()
println('a: ', a)
println('b: ', b)
println('c: ', c) // x is changed here too!
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

`range` creates a [Range](range.html), starting at "from" (inclusive) and going until
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
		return NewRange(from.ToInt(), to.ToInt(), step.ToInt())
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
### `stream`

`stream` returns a "stream" [Struct](struct.html) for an [Iterable](interfaces.html#iterable)
value.

A stream performs a series of transforms on a sequence of iterated values, and then
collects the values into a final result.

Streams have two kinds of fields:

* Transformer functions, which perform some kind of transformation on the sequence of values in the stream.

* Collector functions, which turn the sequence of values into a final result.

Streams are lazy -- calling one of the transformer
functions doesn't do any processing, it simply registers a new transformation,
and then returns the modified stream.  Processing on the sequence of values in the stream
does not start until one of the collector functions is called.

Streams are "single use" values.  Once one of the collector functions has been called,
an error will be thrown if any of the stream's functions are called.

* signature: `stream(itr <Iterable>) <Struct>`
* example:

```
// print the sum of the even squares
let a = [1, 2, 3, 4, 5]
println(stream(a)
    .map(|e| => e*e)
    .filter(|e| => e % 2 == 0)
    .reduce(0, |acc, e| => acc + e))
```

A stream has the following fields:

*/

// BuiltinStream returns a stream Struct.
var BuiltinStream = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {

		ibl, ok := params[0].(Iterable)
		if !ok {
			return nil, fmt.Errorf("TypeMismatch: stream() expected iterable value, got %s", params[0].Type())
		}
		itr, err := ibl.NewIterator(ev)
		if err != nil {
			return nil, err
		}

		s, err := NewStream(itr)

		this, err := NewMethodStruct(s, streamMethods)
		if err != nil {
			return nil, err
		}
		s.(*stream).this = this
		return this, nil
	})

var streamMethods = map[string]Method{

	/*doc
	#### transformer functions

	*/

	/*doc
	* `filter()` adds a "filter" transformation to the stream, by removing elements
	which do not match the provided predicate function.  The predicate function
	must accept one value, and must return a boolean value.

		* signature: `filter(predicate <Func>) <Stream>`
		* predicate signature: `fn(val <Value>) <Bool>`

	*/

	"filter": NewFixedMethod(
		[]Type{FuncType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(*stream)

			// check arity
			fn := params[0].(Func)
			expected := Arity{FixedArity, 1, 0}
			if fn.Arity() != expected {
				return nil, fmt.Errorf(
					"ArityMismatch: filter function must have 1 parameter")
			}

			// transform
			err := self.(Stream).Filter(func(ev Eval, v Value) (Bool, Error) {
				val, err := fn.Invoke(ev, []Value{v})
				if err != nil {
					return nil, err
				}

				result, ok := val.(Bool)
				if !ok {
					return nil, fmt.Errorf(
						"TypeMismatch: filter function must return Bool, not %s", val.Type())
				}
				return result, nil
			})
			if err != nil {
				return nil, err
			}

			return s.this, nil
		}),

	/*doc
	* `map()` adds a "map" transformation to the stream, by transforming elements
	according to the provided mapping function.  The mapping function must accept
	one value, and must return one value.

		* signature: `map(mapping <Func>) <Stream>`
		* mapping signature: `fn(val <Value>) <Value>`

	*/

	"map": NewFixedMethod(
		[]Type{FuncType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(*stream)

			// check arity
			fn := params[0].(Func)
			expected := Arity{FixedArity, 1, 0}
			if fn.Arity() != expected {
				return nil, fmt.Errorf(
					"ArityMismatch: map function must have 1 parameter")
			}

			// transform
			err := self.(Stream).Map(func(ev Eval, v Value) (Value, Error) {
				return fn.Invoke(ev, []Value{v})
			})
			if err != nil {
				return nil, err
			}

			return s.this, nil
		}),

	/*doc
	#### collector functions

	*/

	/*doc
	* `reduce()` reduces the stream to a single value, by applying a "reducer" function
	to an accumulated value and each element in the stream.
	Accumulation is done starting with the first element in the stream,
	and ending with the last.  The reducer function must accept two values, and return one value.

		* signature: `reduce(initial <Value>, reducer <Func>) <List>`
		* reducer signature: `fn(accum <Value>, val <Value>) <Value>`

	*/
	"reduce": NewFixedMethod(
		[]Type{AnyType, FuncType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(*stream)

			initial := params[0]

			if params[1] == Null {
				return nil, NullValueError()
			}

			// check arity
			fn := params[1].(Func)
			expected := Arity{FixedArity, 2, 0}
			if fn.Arity() != expected {
				return nil, fmt.Errorf(
					"ArityMismatch: reduce function must have 2 parameters")
			}

			// invoke
			return s.Reduce(ev, initial, func(ev Eval, acc Value, v Value) (Value, Error) {
				return fn.Invoke(ev, []Value{acc, v})
			})
		}),

	/*doc
	* `toList()` collects the stream's sequence of values into a [List](list.html).

		* signature: `toList() <List>`

	*/
	"toList": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(*stream)
			return s.ToList(ev)
		}),
}

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

## SideEffect Builtins

Golem also has "side-effect" builtins.  These functions
perform I/O, so they should not be included in sandboxed Golem
environments.

* [`print()`](#print)
* [`println()`](#println)

*/

// SideEffectBuiltins are builtins that are not pure functions
var SideEffectBuiltins = []*Builtin{
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
