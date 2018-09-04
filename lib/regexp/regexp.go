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
	Regexp, err = g.NewFieldStruct(
		map[string]g.Field{
			"compile": g.NewField(compile),
		}, true)
	g.Assert(err == nil)
}

// compile compiles a regex expression
var compile g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		s := params[0].(g.Str)

		r, err := regexp.Compile(s.String())
		if err != nil {
			return nil, g.NewError("RegexpError: " + err.Error())
		}

		return makeRegexp(r), nil
	})

func makeRegexp(r *regexp.Regexp) g.Struct {
	stc, err := g.NewMethodStruct(r, regexpMethods)
	g.Assert(err == nil)
	return stc
}

var regexpMethods = map[string]g.Method{

	"matchString": g.NewFixedMethod(
		[]g.Type{g.StrType}, false,
		func(self interface{}, ev g.Eval, params []g.Value) (g.Value, g.Error) {
			r := self.(*regexp.Regexp)
			return matchString(r, params[0].(g.Str)), nil
		}),
}

func matchString(r *regexp.Regexp, s g.Str) g.Bool {
	return g.NewBool(r.MatchString(s.String()))
}
