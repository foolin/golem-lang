// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

//---------------------------------------------------------------
// An intrinsic function is a function that is an intrinsic
// part of a given Type. These functions are created on the
// fly.

type obsoleteIntrinsicFunc struct {
	owner Value
	name  string
	ObsoleteFunc
}

func (f *obsoleteIntrinsicFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *obsoleteIntrinsicFunc:
		// equality for intrinsic functions is based on whether
		// they have the same owner, and the same name
		ownerEq, err := f.owner.Eq(cx, t.owner)
		if err != nil {
			return nil, err
		}
		return NewBool(ownerEq.BoolVal() && (f.name == t.name)), nil
	default:
		return False, nil
	}
}
