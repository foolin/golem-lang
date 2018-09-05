// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package json

import (
	"encoding/json"
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
)

//---------------------------------------------------
// Marshal
//---------------------------------------------------

var Marshal g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.AnyType}, true,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		return marshal(params[0])
	})

func fromValue(val g.Value) (interface{}, g.Error) {

	switch val.Type() {

	case g.NullType:
		return nil, nil

	case g.BoolType:
		return val.(g.Bool).BoolVal(), nil

	case g.IntType:
		return val.(g.Int).IntVal(), nil

	case g.FloatType:
		return val.(g.Float).FloatVal(), nil

	case g.StrType:
		return val.(g.Str).String(), nil

	case g.ListType:

		vals := val.(g.List).Values()
		ifc := make([]interface{}, len(vals))

		for i, iv := range vals {
			fv, err := fromValue(iv)
			if err != nil {
				return nil, err
			}
			ifc[i] = fv
		}

		return ifc, nil

	case g.DictType:

		ifc := make(map[string]interface{})

		itr := val.(g.Dict).HashMap().Iterator()
		for itr.Next() {
			entry := itr.Get()

			s, ok := entry.Key.(g.Str)
			if !ok {
				return nil, g.Error(fmt.Errorf(
					"JsonError: %s is not a valid object key", entry.Key.Type()))
			}
			fv, err := fromValue(entry.Value)
			if err != nil {
				return nil, err
			}
			ifc[s.String()] = fv
		}

		return ifc, nil

	default:
		return nil, g.Error(fmt.Errorf(
			"JsonError: %s cannot be marshalled", val.Type()))
	}
}

func marshal(val g.Value) (g.Str, g.Error) {

	i, err := fromValue(val)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(i)
	if err != nil {
		return nil, g.Error(fmt.Errorf("JsonError: %s", err.Error()))
	}
	return g.NewStr(string(b)), nil
}

//---------------------------------------------------
// Unmarshal
//---------------------------------------------------

var Unmarshal g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		return unmarshal(ev, params[0].(g.Str))
	})

func toValue(ev g.Eval, i interface{}) (g.Value, g.Error) {

	if i == nil {
		return g.Null, nil
	}

	switch t := i.(type) {

	case bool:
		return g.NewBool(t), nil

	case float64:

		n := int64(t)
		if t == float64(n) {
			return g.NewInt(n), nil
		}

		return g.NewFloat(t), nil

	case string:
		return g.NewStr(t), nil

	case []interface{}:

		vals := make([]g.Value, len(t))
		for i, v := range t {
			val, err := toValue(ev, v)
			if err != nil {
				return nil, err
			}
			vals[i] = val
		}
		return g.NewList(vals), nil

	case map[string]interface{}:

		entries := []*g.HEntry{}
		for k, v := range t {

			val, err := toValue(ev, v)
			if err != nil {
				return nil, err
			}

			entries = append(entries,
				&g.HEntry{
					Key:   g.NewStr(k),
					Value: val,
				})
		}
		h, err := g.NewHashMap(ev, entries)
		if err != nil {
			return nil, err
		}
		return g.NewDict(h), nil

	default:
		panic("unreachable")
	}
}

func unmarshal(ev g.Eval, s g.Str) (g.Value, g.Error) {

	var i interface{}

	err := json.Unmarshal([]byte(s.String()), &i)
	if err != nil {
		return nil, g.Error(fmt.Errorf("JsonError: %s", err.Error()))
	}

	return toValue(ev, i)
}
