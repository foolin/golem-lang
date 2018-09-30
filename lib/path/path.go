// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package path

import (
	g "github.com/mjarmy/golem-lang/core"
	fpath "github.com/mjarmy/golem-lang/lib/path/filepath"
)

/*doc

## path

Module path implements utility routines for manipulating slash-separated paths.

*/

/*doc
`encoding` has the following fields:

  * [filepath](lib_pathfilepath.html)

*/

// Path is the "path" module in the standard library
var Path g.Struct

func init() {

	filepath, err := g.NewFrozenStruct(
		map[string]g.Field{
			"ext":  g.NewField(fpath.Ext),
			"walk": g.NewField(fpath.Walk),
		})
	g.Assert(err == nil)

	Path, err = g.NewFrozenStruct(
		map[string]g.Field{
			"filepath": g.NewField(filepath),
		})
	g.Assert(err == nil)
}
