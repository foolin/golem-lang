// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
	os "github.com/mjarmy/golem-lang/lib/os"
	path "github.com/mjarmy/golem-lang/lib/path"
	regexp "github.com/mjarmy/golem-lang/lib/regexp"
)

// BuiltinLib looks up modules in the standard library
var BuiltinLib = g.NewNativeFuncStr(
	func(cx g.Context, name g.Str) (g.Value, g.Error) {

		switch name.String() {
		case "os":
			return os.Os, nil
		case "path":
			return path.Path, nil
		case "regexp":
			return regexp.Regexp, nil
		default:
			return nil, g.NewError(
				"Library",
				fmt.Sprintf("Library '%s' not found", name.String()))
		}
	})
