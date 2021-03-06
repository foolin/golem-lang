// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package json

import (
	"encoding/json"
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
)

/*doc

## `encoding.json`

`encoding.json` implements encoding and decoding of JSON as defined in RFC 7159

*/

/*doc
`encoding.json` has the following fields:

* [marshal](#marshal)
* [marshalIndent](#marshalIndent)
* [unmarshal](#unmarshal)

*/

//---------------------------------------------------
// Marshal
//---------------------------------------------------

/*doc
### `marshal`

`marshal` returns the JSON encoding of a value.  Null, Bool, Float Int, Str, and List
are marshalled as their corresponding JSON elements.  Structs and Dicts are marshalled
as JSON objects. Other golem types cannot be marshalled.

* signature: `marshal(value <Value>) <Str>`
* example:

```
import encoding
let s = struct { a: [1, 2, 3]}
println(encoding.json.marshal(s))
```

*/

// Marshal marshals a Value into a JSON string
var Marshal g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.AnyType}, true,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		return marshal(ev, params[0])
	})

/*doc
### `marshalIndent`

`marshalIndent` is like `marshal` but applies indent to format the output.
Each JSON element in the output will begin on a new line beginning with prefix
followed by one or more copies of indent according to the indentation nesting.

* signature: `marshalIndent(value <Value>, prefix <Str>, indent <Str>) <Str>`
* example:

```
import encoding
let s = struct { a: [1, 2, 3]}
println(encoding.json.marshalIndent(s, '', '  '))
```

*/

// MarshalIndent marshals a Value into a JSON string
var MarshalIndent g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.AnyType, g.StrType, g.StrType}, true,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		return marshalIndent(
			ev, params[0], params[1].(g.Str), params[2].(g.Str))
	})

func fromList(ev g.Eval, list g.List) (interface{}, g.Error) {

	vals := list.Values()
	ifc := make([]interface{}, len(vals))

	for i, iv := range vals {
		fv, err := fromValue(ev, iv)
		if err != nil {
			return nil, err
		}
		ifc[i] = fv
	}

	return ifc, nil
}

func fromDict(ev g.Eval, dict g.Dict) (interface{}, g.Error) {

	ifc := make(map[string]interface{})

	itr := dict.HashMap().Iterator()
	for itr.Next() {
		entry := itr.Get()

		s, ok := entry.Key.(g.Str)
		if !ok {
			return nil, g.Error(fmt.Errorf(
				"JsonError: %s is not a valid object key", entry.Key.Type()))
		}

		fv, err := fromValue(ev, entry.Value)
		if err != nil {
			return nil, err
		}
		ifc[s.String()] = fv
	}

	return ifc, nil
}

func fromStruct(ev g.Eval, st g.Struct) (interface{}, g.Error) {

	ifc := make(map[string]interface{})

	names, err := st.FieldNames()
	if err != nil {
		return nil, err
	}

	for _, k := range names {

		val, err := st.GetField(ev, k)
		if err != nil {
			return nil, err
		}

		fv, err := fromValue(ev, val)
		if err != nil {
			return nil, err
		}

		ifc[k] = fv
	}

	return ifc, nil
}

func fromValue(ev g.Eval, val g.Value) (interface{}, g.Error) {

	switch val.Type() {

	case g.NullType:
		return nil, nil

	case g.BoolType:
		return val.(g.Bool).BoolVal(), nil

	case g.IntType:
		return val.(g.Int).ToInt(), nil

	case g.FloatType:
		return val.(g.Float).ToFloat(), nil

	case g.StrType:
		return val.(g.Str).String(), nil

	case g.ListType:
		return fromList(ev, val.(g.List))

	case g.DictType:
		return fromDict(ev, val.(g.Dict))

	case g.StructType:
		return fromStruct(ev, val.(g.Struct))

	default:
		return nil, g.Error(fmt.Errorf(
			"JsonError: %s cannot be marshalled", val.Type()))
	}
}

func marshal(ev g.Eval, val g.Value) (g.Str, g.Error) {

	ifc, err := fromValue(ev, val)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(ifc)
	if err != nil {
		return nil, g.Error(fmt.Errorf("JsonError: %s", err.Error()))
	}
	return g.NewStr(string(b))
}

func marshalIndent(ev g.Eval, val g.Value, prefix, indent g.Str) (g.Str, g.Error) {

	ifc, err := fromValue(ev, val)
	if err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(ifc, prefix.String(), indent.String())
	if err != nil {
		return nil, g.Error(fmt.Errorf("JsonError: %s", err.Error()))
	}
	return g.NewStr(string(b))
}

//---------------------------------------------------
// Unmarshal
//---------------------------------------------------

/*doc
### `unmarshal`

`unmarshal` parses JSON-encoded data.
The optional useStructs parameter, which defaults to false, specifies whether
the data should be marshalled into structs rather than dicts.  If this parameter
is set to true, then the keys of the JSON objects in the data must all be valid
Golem identifiers.

* signature: `unmarshal(text <Str>, useStructs = false <Bool>) <Value>`
* example:

```
import encoding
let text = `{
  "a": [ 1, 2, 3 ]
}`
println(encoding.json.unmarshal(text))
println(encoding.json.unmarshal(text, true))
```

*/

// Unmarshal unmarshals a JSON string into a Value
var Unmarshal g.Value = g.NewMultipleNativeFunc(
	[]g.Type{g.StrType},
	[]g.Type{g.BoolType},
	false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		s := params[0].(g.Str)
		useStructs := false
		if len(params) == 2 {
			useStructs = params[1].(g.Bool).BoolVal()
		}

		return unmarshal(ev, s, useStructs)
	})

func toList(ev g.Eval, ifc []interface{}, useStructs bool) (g.Value, g.Error) {

	vals := make([]g.Value, len(ifc))
	for i, v := range ifc {
		val, err := toValue(ev, v, useStructs)
		if err != nil {
			return nil, err
		}
		vals[i] = val
	}
	return g.NewList(vals), nil
}

func toStruct(ev g.Eval, ifc map[string]interface{}) (g.Value, g.Error) {

	fields := make(map[string]g.Field)
	for k, v := range ifc {

		val, err := toValue(ev, v, true)
		if err != nil {
			return nil, err
		}

		fields[k] = g.NewField(val)
	}
	return g.NewStruct(fields)
}

func toDict(ev g.Eval, ifc map[string]interface{}) (g.Value, g.Error) {

	entries := []*g.HEntry{}
	for k, v := range ifc {

		val, err := toValue(ev, v, false)
		if err != nil {
			return nil, err
		}

		ks, err := g.NewStr(k)
		if err != nil {
			return nil, err
		}

		entries = append(entries,
			&g.HEntry{Key: ks, Value: val})
	}
	h, err := g.NewHashMap(ev, entries)
	if err != nil {
		return nil, err
	}
	return g.NewDict(h), nil
}

func toValue(ev g.Eval, ifc interface{}, useStructs bool) (g.Value, g.Error) {

	if ifc == nil {
		return g.Null, nil
	}

	switch t := ifc.(type) {

	case bool:
		return g.NewBool(t), nil

	case float64:
		n := int64(t)
		if t == float64(n) {
			return g.NewInt(n), nil
		}
		return g.NewFloat(t), nil

	case string:
		return g.NewStr(t)

	case []interface{}:
		return toList(ev, t, useStructs)

	case map[string]interface{}:
		if useStructs {
			return toStruct(ev, t)
		}
		return toDict(ev, t)

	default:
		panic("unreachable")
	}
}

func unmarshal(ev g.Eval, s g.Str, useStructs bool) (g.Value, g.Error) {

	var ifc interface{}

	err := json.Unmarshal([]byte(s.String()), &ifc)
	if err != nil {
		return nil, g.Error(fmt.Errorf("JsonError: %s", err.Error()))
	}

	return toValue(ev, ifc, useStructs)
}
