// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	g "github.com/mjarmy/golem-lang/core"
	"regexp"
)

type regexModule struct {
	contents g.Struct
}

// InitRegexModule initializes the 'regex' module.
func InitRegexModule() g.Module {

	compile := g.NewNativeFunc(
		1, 1,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			s, ok := values[0].(g.Str)
			if !ok {
				return nil, g.TypeMismatchError("Expected Str")
			}

			rgx, err := regexp.Compile(s.String())
			if err != nil {
				return nil, g.MakeError("RegexError", err.Error())
			}

			return makePattern(rgx), nil
		})

	contents, err := g.NewStruct([]g.Field{g.NewField("compile", true, compile)}, true)
	if err != nil {
		panic("InitRegexModule")
	}

	return &regexModule{contents}
}

func makePattern(rgx *regexp.Regexp) g.Struct {

	match := g.NewNativeFunc(
		1, 1,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			s, ok := values[0].(g.Str)
			if !ok {
				return nil, g.TypeMismatchError("Expected Str")
			}

			return g.MakeBool(rgx.MatchString(s.String())), nil
		})

	pattern, err := g.NewStruct([]g.Field{g.NewField("match", true, match)}, true)
	if err != nil {
		panic("InitRegexModule")
	}

	return pattern
}

func (m *regexModule) GetModuleName() string {
	return "regex"
}

func (m *regexModule) GetContents() g.Struct {
	return m.contents
}
