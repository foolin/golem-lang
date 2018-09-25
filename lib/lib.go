// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/lib/encoding"
	"github.com/mjarmy/golem-lang/lib/golem"
	"github.com/mjarmy/golem-lang/lib/io"
	"github.com/mjarmy/golem-lang/lib/os"
	"github.com/mjarmy/golem-lang/lib/path"
	"github.com/mjarmy/golem-lang/lib/regexp"
)

// BuiltinLib looks up modules in the standard library
var BuiltinLib = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, values []g.Value) (g.Value, g.Error) {

		name := values[0].(g.Str)

		switch name.String() {
		case "encoding":
			return encoding.Encoding, nil
		case "golem":
			return golem.Golem, nil
		case "io":
			return io.Io, nil
		case "os":
			return os.Os, nil
		case "path":
			return path.Path, nil
		case "regexp":
			return regexp.Regexp, nil
		default:
			return nil, g.Error(fmt.Errorf(
				"LibraryNotFound: Library '%s' not found", name.String()))
		}
	})
