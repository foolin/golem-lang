// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

/*doc
## Chan

A Chan is a conduit through which you can send and receive values.
A new Chan is created by the [`chan()`](builtins.html#chan) builtin function.

Valid operators for Chan are:

* The equality operators `==`, `!=`

*/

type channel struct {
	ch chan Value
}

// NewChan creates a new Chan
func NewChan() Chan {
	return &channel{make(chan Value)}
}

// NewBufferedChan creates a new buffered Chan
func NewBufferedChan(size int) Chan {
	return &channel{make(chan Value, size)}
}

func (ch *channel) chanMarker() {}

func (ch *channel) Type() Type { return ChanType }

func (ch *channel) Freeze(ev Eval) (Value, Error) {
	return ch, nil
}

func (ch *channel) Frozen(ev Eval) (Bool, Error) {
	return True, nil
}

func (ch *channel) Eq(ev Eval, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *channel:
		// equality is based on identity
		return NewBool(ch == t), nil
	default:
		return False, nil
	}
}

func (ch *channel) HashCode(ev Eval) (Int, Error) {
	return nil, HashCodeMismatch(ChanType)
}

func (ch *channel) ToStr(ev Eval) (Str, Error) {
	return NewStr(fmt.Sprintf("chan<%p>", ch))
}

func (ch *channel) Send(val Value) {
	ch.ch <- val
}

func (ch *channel) Recv() Value {
	return <-ch.ch
}

//--------------------------------------------------------------
// fields

/*doc
Chan has the following fields:

* [send](#send)
* [recv](#recv)

*/

var chanMethods = map[string]Method{

	/*doc
	### `send`

	`send` sends a value to the chan.

	* signature: `send(val <Value>)`

	*/
	"send": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			ch := self.(Chan)
			ch.Send(params[0])
			return Null, nil
		}),

	/*doc
	### `recv`

	`recv` receives a value from the chan.

	* signature: `recv() <Value>`

	*/
	"recv": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			ch := self.(Chan)
			val := ch.Recv()
			return val, nil
		}),
}

func (ch *channel) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(chanMethods))
	for name := range chanMethods {
		names = append(names, name)
	}
	return names, nil
}

func (ch *channel) HasField(name string) (bool, Error) {
	_, ok := chanMethods[name]
	return ok, nil
}

func (ch *channel) GetField(ev Eval, name string) (Value, Error) {
	if method, ok := chanMethods[name]; ok {
		return method.ToFunc(ch, name), nil
	}
	return nil, NoSuchField(name)
}

func (ch *channel) InvokeField(ev Eval, name string, params []Value) (Value, Error) {

	if method, ok := chanMethods[name]; ok {
		return method.Invoke(ch, ev, params)
	}
	return nil, NoSuchField(name)
}
