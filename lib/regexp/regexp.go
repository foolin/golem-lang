// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package regexp

import (
	"regexp"

	g "github.com/mjarmy/golem-lang/core"
)

// Regexp is the "regexp" module in the standard library
var Regexp g.Struct

func init() {
	var err error
	Regexp, err = g.NewStruct([]g.Field{
		g.NewField("compile", true, compile),
	}, true)
	if err != nil {
		panic("unreachable")
	}
}

// compile compiles a regex expression
var compile g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(cx g.Context, values []g.Value) (g.Value, g.Error) {
		s := values[0].(g.Str)

		rgx, err := regexp.Compile(s.String())
		if err != nil {
			return nil, g.NewError("RegexpError", err.Error())
		}

		return makeRegexp(rgx), nil
	})

func makeRegexp(rgx *regexp.Regexp) g.Struct {

	matchString := g.NewFixedNativeFunc(
		[]g.Type{g.StrType}, false,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			s := values[0].(g.Str)
			return g.NewBool(rgx.MatchString(s.String())), nil
		})

	stc, err := g.NewStruct(
		[]g.Field{g.NewField("matchString", true, matchString)},
		true)
	if err != nil {
		panic("unreachable")
	}
	return stc
}
