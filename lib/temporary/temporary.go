// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"regexp"

	g "github.com/mjarmy/golem-lang/core"
)

// Compile compiles a regex expression
var Compile g.Value = g.NewNativeFuncStr(
	func(cx g.Context, s g.Str) (g.Value, g.Error) {

		rgx, err := regexp.Compile(s.String())
		if err != nil {
			return nil, g.NewError("RegexpError", err.Error())
		}

		return makeRegexp(rgx), nil
	})

func makeRegexp(rgx *regexp.Regexp) g.Struct {

	matchString := g.NewNativeFuncStr(
		func(cx g.Context, s g.Str) (g.Value, g.Error) {
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
