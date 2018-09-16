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
on Golem values and code.

*/

// Golem is the "golem" module in the standard library
var Golem g.Struct

func init() {
	var err error
	Golem, err = g.NewFrozenFieldStruct(
		map[string]g.Field{
			"fields":   g.NewField(fields),
			"getField": g.NewField(getField),
			"hasField": g.NewField(hasField),
			"setField": g.NewField(setField),
		})
	g.Assert(err == nil)
}

/*doc
#### `fields`

`fields` returns a Set of the names of a value's fields.

	* signature: `fields(value <Value>) <Set>`
	* example:

	```
    import golem

    println(golem.fields([]))
    println(golem.fields(struct { a: 1, b: 2}))
	```

*/

var fields = g.NewFixedNativeFunc(
	[]g.Type{g.AnyType},
	false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {

		fields, err := params[0].FieldNames()
		if err != nil {
			return nil, err
		}

		entries := make([]g.Value, len(fields))
		for i, k := range fields {
			entry, err := g.NewStr(k)
			if err != nil {
				return nil, err
			}
			entries[i] = entry
		}
		return g.NewSet(ev, entries)
	})

/*doc
#### `getField`

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
#### `hasField`

`hasField` returns whether a value has a field with a given name.

	* signature: `getField(value <Value>, name <Str>) <Bool>`
	* example:

	```
    import golem

    let a = [1, 2]
    println(golem.hasField(a, 'add'))
	```

*/

var hasField = g.NewFixedNativeFunc(
	[]g.Type{g.AnyType, g.StrType},
	false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		field := params[1].(g.Str)

		b, err := params[0].HasField(field.String())
		if err != nil {
			return nil, err
		}
		return g.NewBool(b), nil
	})

/*doc
#### `setField`

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
