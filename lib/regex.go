// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lib

import (
	g "github.com/mjarmy/golem-lang/core"
	"regexp"
)

type regexModule struct {
	contents g.Struct
}

func InitRegexModule() g.Module {

	compile := g.NewNativeFunc(
		func(values []g.Value) (g.Value, g.Error) {
			if len(values) != 1 {
				return nil, g.ArityMismatchError("1", len(values))
			}
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

	contents, err := g.NewStruct([]*g.StructEntry{
		{"compile", true, false, compile}})
	g.Assert(err == nil, "InitRegexModule")
	return &regexModule{contents}
}

func makePattern(rgx *regexp.Regexp) g.Struct {

	match := g.NewNativeFunc(
		func(values []g.Value) (g.Value, g.Error) {
			if len(values) != 1 {
				return nil, g.ArityMismatchError("1", len(values))
			}
			s, ok := values[0].(g.Str)
			if !ok {
				return nil, g.TypeMismatchError("Expected Str")
			}

			return g.MakeBool(rgx.MatchString(s.String())), nil
		})

	pattern, err := g.NewStruct([]*g.StructEntry{
		{"match", true, false, match}})
	g.Assert(err == nil, "InitSysModule")
	return pattern
}

func (m *regexModule) GetModuleName() string {
	return "regex"
}

func (m *regexModule) GetContents() g.Struct {
	return m.contents
}
