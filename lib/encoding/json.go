package main

import (
	"encoding/json"
	//"fmt"

	g "github.com/mjarmy/golem-lang/core"
)

func assert(flag bool) {
	if !flag {
		panic("assertion failure")
	}
}

func mustStr(val g.Value) string {
	s, e := val.ToStr(nil)
	if e != nil {
		panic("mustStr")
	}
	return s.String()
}

//---------------------------------------------------

func fromValue(v g.Value) interface{} {

	switch v.Type() {

	case g.NullType:
		return nil

	case g.BoolType:
		return v.(g.Bool).BoolVal()

	case g.IntType:
		return v.(g.Int).IntVal()

	case g.FloatType:
		return v.(g.Float).FloatVal()

	case g.StrType:
		return v.(g.Str).String()

	case g.ListType:

		vals := v.(g.List).Values()
		ifc := make([]interface{}, len(vals))
		for i, v := range vals {
			ifc[i] = fromValue(v)
		}
		return ifc

	//case g.DictType:
	//	d := v.(g.Dict)
	//	ifc := make(map[string]interface{})

	default:
		panic("unreachable")
	}
}

func marshal(val g.Value) string {

	var i = fromValue(val)

	b, e := json.Marshal(i)
	if e != nil {
		panic(e)
	}

	return string(b)
}

//---------------------------------------------------

func toValue(i interface{}) g.Value {

	if i == nil {
		return g.Null
	}

	switch t := i.(type) {

	case bool:
		return g.NewBool(t)

	case float64:

		n := int64(t)
		if t == float64(n) {
			return g.NewInt(n)
		}

		return g.NewFloat(t)

	case string:
		return g.NewStr(t)

	case []interface{}:

		vals := make([]g.Value, len(t))
		for i, x := range t {
			vals[i] = toValue(x)
		}
		return g.NewList(vals)

	//case map[string]interface{}:
	//	entries := []*g.HEntry{}
	//	for k, v := range t {
	//		entries = append(entries,
	//			&g.HEntry{g.NewStr(k), toValue(v)})
	//	}
	//	h, err := g.NewHashMap(ev, entries)
	//	g.Assert(err == nil)
	//	return g.NewDict(h)

	default:
		panic("unreachable")
	}
}

func unmarshal(s string) g.Value {

	var i interface{}

	e := json.Unmarshal([]byte(s), &i)
	if e != nil {
		panic(e)
	}

	return toValue(i)
}

//---------------------------------------------------

func roundTrip(a string) string {
	b := marshal(unmarshal(a))
	println(b)
	return b
}

func main() {
	assert(roundTrip("null") == "null")
	assert(roundTrip("true") == "true")
	assert(roundTrip("false") == "false")
	assert(roundTrip("1.1") == "1.1")
	assert(roundTrip("1") == "1")
	assert(roundTrip(`"abc"`) == `"abc"`)
	assert(roundTrip(`[1,2,3]`) == `[1,2,3]`)
	assert(roundTrip(`[]`) == `[]`)
	assert(roundTrip(`["a",[null]]`) == `["a",[null]]`)
}
