// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golem

import (
	g "github.com/mjarmy/golem-lang/core"
)

/*doc

## golem

The golem module defines functions that perform introspection and transformation
on Golem values and code, as well as providing a few useful utility functions.

*/

// Golem is the "golem" module in the standard library
var Golem g.Module

func init() {
	golem, err := g.NewFrozenStruct(
		map[string]g.Field{
			"getField":     g.NewField(getField),
			"makeHashCode": g.NewField(makeHashCode),
			"setField":     g.NewField(setField),
			"toDict":       g.NewField(toDict),
		})
	g.Assert(err == nil)

	Golem = g.NewNativeModule("golem", golem)
}

/*doc
`golem` has the following fields:

* [getField](#getfield)
* [makeHashCode](#makehashcode)
* [setField](#setfield)
* [toDict](#todict)

*/

/*doc
### `getField`

`getField` returns the value associated with a field name.

* signature: `getField(value <Value>, name <Str>) <Value>`
* example:

```
import golem
let a = [1, 2]
let f = golem.getField(a, 'add')
f(3)
println(a)
```

*/

var getField = g.NewFixedNativeFunc(
	[]g.Type{g.AnyType, g.StrType},
	false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		field := params[1].(g.Str)

		return params[0].GetField(ev, field.String())
	})

/*doc
### `makeHashCode`

`makeHashCode` generates a hashCode for a sequence of
[`hashable`](interfaces.html#hashable) values. `makeHashCode` uses the
[Jenkins hash function](https://en.wikipedia.org/wiki/Jenkins_hash_function)
algorithm.

* signature: `makeHashCode(values... <Value>) <Int>`
* example:

```
import golem
println(golem.makeHashCode(1, 2, 3))
```

*/

var makeHashCode = g.NewVariadicNativeFunc(
	[]g.Type{}, g.AnyType, true,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {

		// https://en.wikipedia.org/wiki/Jenkins_hash_function
		var hash int64
		for _, v := range params {
			h, err := v.HashCode(ev)
			if err != nil {
				return nil, err
			}

			hash += h.ToInt()
			hash += hash << 10
			hash ^= hash >> 6
		}
		hash += hash << 3
		hash ^= hash >> 11
		hash += hash << 15
		return g.NewInt(hash), nil
	})

/*doc
### `setField`

`setField` sets the value associated with a field name. `setField` only works
on Structs -- you cannot set the fields of other types. `setField` returns `null`
if it was successful.

* signature: `setField(s <Struct>, name <Str>, value <Value>) <Null>`
* example:

```
import golem
let s = struct { a: 1, b: 2 }
golem.setField(s, 'a', 3)
println(s)
```

*/

var setField = g.NewFixedNativeFunc(
	[]g.Type{g.StructType, g.StrType, g.AnyType},
	true,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {

		if params[0].Type() == g.NullType {
			return nil, g.NullValueError()
		}
		if params[1].Type() == g.NullType {
			return nil, g.NullValueError()
		}

		st := params[0].(g.Struct)
		fld := params[1].(g.Str)

		err := st.SetField(ev, fld.String(), params[2])
		if err != nil {
			return nil, err
		}
		return g.Null, nil
	})

/*doc
### `toDict`

`toDict` converts a Struct into a Dict.

* signature: `toDict(s <Struct>) <Dict>`
* example:

```
import golem
println(golem.toDict(struct { a: 1, b: 2 }))
```

*/

var toDict = g.NewFixedNativeFunc(
	[]g.Type{g.StructType},
	true,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		st := params[0].(g.Struct)
		return st.ToDict(ev)
	})
