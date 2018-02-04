// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	g "github.com/mjarmy/golem-lang/core"
	"regexp"
)

type regexpModule struct {
	contents g.Struct
}

// NewRegexpModule creates the 'regexp' module.
func NewRegexpModule() g.Module {

	compile := g.NewNativeFunc(
		1, 1,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			s, ok := values[0].(g.Str)
			if !ok {
				return nil, g.TypeMismatchError("Expected Str")
			}

			rgx, err := regexp.Compile(s.String())
			if err != nil {
				return nil, g.NewError("RegexpError", err.Error())
			}

			return makeRegexp(rgx), nil
		})

	contents, err := g.NewStruct([]g.Field{
		g.NewField("compile", true, compile)},
		true)

	if err != nil {
		panic("NewRegexpModule")
	}
	return &regexpModule{contents}
}

func makeRegexp(rgx *regexp.Regexp) g.Struct {

	match := g.NewNativeFunc(
		1, 1,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			s, ok := values[0].(g.Str)
			if !ok {
				return nil, g.TypeMismatchError("Expected Str")
			}
			return g.NewBool(rgx.MatchString(s.String())), nil
		})

	pattern, err := g.NewStruct(
		[]g.Field{g.NewField("matchString", true, match)},
		true)

	if err != nil {
		panic("NewRegexpModule")
	}
	return pattern
}

func (m *regexpModule) GetModuleName() string {
	return "regexp"
}

func (m *regexpModule) GetContents() g.Struct {
	return m.contents
}
