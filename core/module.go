// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

// A Module is a namespace that contains a Struct.
type Module interface {
	Name() string
	Contents() Struct
}

// A NativeModule is a Module that is implemented in Go
type NativeModule struct {
	name     string
	contents Struct
}

// NewNativeModule creates a new NativeModule
func NewNativeModule(name string, contents Struct) *NativeModule {
	return &NativeModule{name, contents}
}

func (m *NativeModule) Name() string {
	return m.name
}

func (m *NativeModule) Contents() Struct {
	return m.contents
}
