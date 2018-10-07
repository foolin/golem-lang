// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/lib/encoding"
	"github.com/mjarmy/golem-lang/lib/golem"
	"github.com/mjarmy/golem-lang/lib/io"
	"github.com/mjarmy/golem-lang/lib/os"
	"github.com/mjarmy/golem-lang/lib/path"
	"github.com/mjarmy/golem-lang/lib/regexp"
)

var SandboxLibrary = []g.Module{
	encoding.Encoding,
	golem.Golem,
	regexp.Regexp,
}

var SideEffectLibrary = []g.Module{
	io.Io,
	os.Os,
	path.Path,
}
