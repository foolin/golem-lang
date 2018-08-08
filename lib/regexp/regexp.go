// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"regexp"

	g "github.com/mjarmy/golem-lang/core"
)

type module struct{ contents g.Struct }

func (m *module) GetModuleName() string { return "regexp" }
func (m *module) GetContents() g.Struct { return m.contents }

// LoadModule creates the 'path' module.
func LoadModule() (g.Module, g.Error) {

	contents, err := g.NewStruct([]g.Field{
		g.NewField("compile", true, compile())},
		true)
	if err != nil {
		return nil, err
	}

	return &module{contents}, nil
}

func compile() g.NativeFunc {

	return g.NewNativeFuncStr(
		func(cx g.Context, s g.Str) (g.Value, g.Error) {

			rgx, err := regexp.Compile(s.String())
			if err != nil {
				return nil, g.NewError("RegexpError", err.Error())
			}

			return makeRegexp(rgx), nil
		})
}

func makeRegexp(rgx *regexp.Regexp) g.Struct {

	match := g.NewNativeFuncStr(
		func(cx g.Context, s g.Str) (g.Value, g.Error) {
			return g.NewBool(rgx.MatchString(s.String())), nil
		})

	pattern, err := g.NewStruct(
		[]g.Field{g.NewField("matchString", true, match)},
		true)

	if err != nil {
		panic("unreachable")
	}
	return pattern
}
