// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package exec

import (
	"os/exec"

	g "github.com/mjarmy/golem-lang/core"
)

/*doc

### `os.exec`

`os.exec` runs external commands.

*/

/*doc
#### `runCommand`

`runCommand` TODO replace with command().run()

	* signature: `runCommand(path <Str>, args... <Str>) <Null>`
*/

var RunCommand g.Value = g.NewVariadicNativeFunc(
	[]g.Type{g.StrType}, g.StrType, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {

		path := params[0].(g.Str).String()
		args := make([]string, len(params)-1)
		for i := 1; i < len(params); i++ {
			args[i-1] = params[i].(g.Str).String()
		}

		cmd := exec.Command(path, args...)
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
		return g.Null, nil
	})
