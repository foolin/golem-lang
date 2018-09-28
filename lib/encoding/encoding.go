// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package encoding

import (
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/lib/encoding/json"
)

/*doc

## encoding

The encoding module defines functionality that converts
data to and from byte-level and textual representations.

*/

/*doc
`encoding` has the following fields:

  * [json](lib_encodingjson.html)

*/

// Encoding is the "encoding" module in the standard library
var Encoding g.Struct

func init() {

	json, err := g.NewFrozenFieldStruct(
		map[string]g.Field{
			"marshal":       g.NewField(json.Marshal),
			"marshalIndent": g.NewField(json.MarshalIndent),
			"unmarshal":     g.NewField(json.Unmarshal),
		})
	g.Assert(err == nil)

	Encoding, err = g.NewFrozenFieldStruct(
		map[string]g.Field{
			"json": g.NewField(json),
		})
	g.Assert(err == nil)
}
