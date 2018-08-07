// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"plugin"
	"reflect"
)

type _plugin struct {
	path string
	plug *plugin.Plugin
}

// NewPlugin creates a new Plugin
func NewPlugin(cx Context, path Str) (Plugin, Error) {

	// TODO make sure path doesn't have '..', etc
	plugPath := cx.HomePath() + "/lib/" + path.String() + "/" + path.String() + ".so"
	plug, err := plugin.Open(plugPath)
	if err != nil {
		return nil, PluginError(path.String(), err)
	}
	return &_plugin{path: plugPath, plug: plug}, nil
}

func (p *_plugin) pluginMarker() {}

func (p *_plugin) Type() Type { return PluginType }

func (p *_plugin) Freeze() (Value, Error) {
	return p, nil
}

func (p *_plugin) Frozen() (Bool, Error) {
	return True, nil
}

func (p *_plugin) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *_plugin:
		// equality is based on path
		return NewBool(p.path == t.path), nil
	default:
		return False, nil
	}
}

func (p *_plugin) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (p *_plugin) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (p *_plugin) ToStr(cx Context) Str {
	return NewStr(fmt.Sprintf("plugin<%s>", p.path))
}

//--------------------------------------------------------------
// intrinsic functions

func (p *_plugin) GetField(cx Context, key Str) (Value, Error) {
	switch sn := key.String(); sn {

	case "lookup":
		return &intrinsicFunc{p, sn, NewNativeFunc(
			1, 1,
			func(cx Context, values []Value) (Value, Error) {

				name, ok := values[0].(Str)
				if !ok {
					return nil, TypeMismatchError("Expected Str")
				}

				sym, err := p.plug.Lookup(name.String())
				if err != nil {
					return nil, PluginError(p.path, err)
				}

				value, ok := sym.(*Value)
				if !ok {
					return nil, PluginError(p.path, fmt.Errorf(
						"plugin symbol '%s' is not a Value: %s",
						name.String(),
						reflect.TypeOf(sym)))
				}
				return *value, nil

			})}, nil

	default:
		return nil, NoSuchFieldError(key.String())
	}
}
