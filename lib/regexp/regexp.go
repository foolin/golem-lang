// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package regexp

import (
	"fmt"
	"regexp"

	g "github.com/mjarmy/golem-lang/core"
)

/*doc

# regexp

Module regexp implements regular expression search.

The syntax of the regular expressions accepted is the same general syntax used by
Perl, Python, and other languages.

*/

/*doc

`regexp` has the following fields:

* [compile](#compile)

`regexp` defines the following structs:

* [regex](#regex)

*/

// Regexp is the "regexp" module in the standard library
var Regexp g.Struct

func init() {
	var err error
	Regexp, err = g.NewFrozenFieldStruct(
		map[string]g.Field{
			"compile": g.NewField(compile),
		})
	g.Assert(err == nil)
}

/*doc

## Fields

*/

/*doc
### `compile`

`compile` parses a regular expression and returns, if successful, a
[regex](#regex) struct that can be used to match against text.

* signature: `compile(expr <Str>) <Struct>`

*/

// compile compiles a regex expression
var compile g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		s := params[0].(g.Str)

		r, err := regexp.Compile(s.String())
		if err != nil {
			return nil, g.Error(fmt.Errorf("RegexpError: %s", err.Error()))
		}

		return makeRegexp(r), nil
	})

func makeRegexp(r *regexp.Regexp) g.Struct {
	stc, err := g.NewMethodStruct(r, regexpMethods)
	g.Assert(err == nil)
	return stc
}

/*doc

## Structs

*/

/*doc
### `regex`

`regex` is a Struct that is the representation of a compiled regular expression. A
Regexp is safe for concurrent use by multiple goroutines, except for configuration
methods, such as Longest.

A `regex` struct has the fields:

* [matchString](#matchstring)

*/

var regexpMethods = map[string]g.Method{

	/*doc
	#### `matchString`

	`matchString` reports whether the string s contains any match of the regular expression.

	* signature: `matchString(<Str>) <Bool>`

	*/
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
