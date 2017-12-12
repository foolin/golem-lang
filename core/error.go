// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"fmt"
	"strings"
)

//---------------------------------------------------------------
// ErrorKind

type ErrorKind int

const (
	GENERIC ErrorKind = iota
	NULL_VALUE
	TYPE_MISMATCH
	ARITY_MISMATCH
	TUPLE_LENGTH
	DIVIDE_BY_ZERO
	INDEX_OUT_OF_BOUNDS
	NO_SUCH_FIELD
	READONLY_FIELD
	DUPLICATE_FIELD
	INVALID_ARGUMENT
	NO_SUCH_ELEMENT
	ASSERTION_FAILED
	CONST_SYMBOL
	UNDEFINIED_SYMBOL
)

func (t ErrorKind) String() string {
	switch t {

	case GENERIC:
		return "Generic"
	case NULL_VALUE:
		return "NullValue"
	case TYPE_MISMATCH:
		return "TypeMismatch"
	case ARITY_MISMATCH:
		return "ArityMismatch"
	case TUPLE_LENGTH:
		return "TupleLength"
	case DIVIDE_BY_ZERO:
		return "DivideByZero"
	case INDEX_OUT_OF_BOUNDS:
		return "IndexOutOfBounds"
	case NO_SUCH_FIELD:
		return "NoSuchField"
	case READONLY_FIELD:
		return "ReadonlyField"
	case DUPLICATE_FIELD:
		return "DuplicateField"
	case INVALID_ARGUMENT:
		return "InvalidArgument"
	case NO_SUCH_ELEMENT:
		return "NoSuchElement"
	case ASSERTION_FAILED:
		return "AssertionFailed"
	case CONST_SYMBOL:
		return "ConstSymbol"
	case UNDEFINIED_SYMBOL:
		return "UndefinedSymbol"

	default:
		panic("unreachable")
	}
}

//---------------------------------------------------------------
// Error

type Error interface {
	error
	Kind() ErrorKind
	Struct() Struct
}

type serror struct {
	kind ErrorKind
	stc  Struct
}

func (e *serror) Error() string {
	kind, kerr := e.stc.Get(MakeStr("kind"))
	msg, merr := e.stc.Get(MakeStr("msg"))

	if kerr == nil {
		if merr == nil {
			return strings.Join([]string{
				kind.ToStr().String(), ": ", msg.ToStr().String()}, "")
		} else {
			return kind.ToStr().String()
		}
	} else {
		return e.stc.ToStr().String()
	}
}

func (e *serror) Kind() ErrorKind {
	return e.kind
}

func (e *serror) Struct() Struct {
	return e.stc
}

func makeError(kind ErrorKind, msg string) Error {
	var stc Struct
	var err Error
	if msg == "" {
		stc, err = NewStruct([]*StructEntry{
			{"kind", true, false, MakeStr(kind.String())}})
	} else {
		stc, err = NewStruct([]*StructEntry{
			{"kind", true, false, MakeStr(kind.String())},
			{"msg", true, false, MakeStr(msg)}})
	}
	if err != nil {
		panic("invalid struct")
	}

	return &serror{kind, stc}
}

func MakeError(kind string, msg string) Error {
	var stc Struct
	var err Error
	if msg == "" {
		stc, err = NewStruct([]*StructEntry{
			{"kind", true, false, MakeStr(kind)}})
	} else {
		stc, err = NewStruct([]*StructEntry{
			{"kind", true, false, MakeStr(kind)},
			{"msg", true, false, MakeStr(msg)}})
	}
	if err != nil {
		panic("invalid struct")
	}

	return &serror{GENERIC, stc}
}

func GenericError(stc Struct) Error {
	return &serror{GENERIC, stc}
}

func NullValueError() Error {
	return makeError(NULL_VALUE, "")
}

func TypeMismatchError(msg string) Error {
	return makeError(TYPE_MISMATCH, msg)
}

func ArityMismatchError(expected string, actual int) Error {
	return makeError(
		ARITY_MISMATCH,
		fmt.Sprintf("Expected %s params, got %d", expected, actual))
}

func TupleLengthError(expected int, actual int) Error {
	return makeError(
		TUPLE_LENGTH,
		fmt.Sprintf("Expected Tuple of length %d, got %d", expected, actual))
}

func DivideByZeroError() Error {
	return makeError(DIVIDE_BY_ZERO, "")
}

func IndexOutOfBoundsError() Error {
	return makeError(INDEX_OUT_OF_BOUNDS, "")
}

func NoSuchFieldError(field string) Error {
	return makeError(
		NO_SUCH_FIELD,
		fmt.Sprintf("Field '%s' not found", field))
}

func ReadonlyFieldError(field string) Error {
	return makeError(
		READONLY_FIELD,
		fmt.Sprintf("Field '%s' is readonly", field))
}

func DuplicateFieldError(field string) Error {
	return makeError(
		DUPLICATE_FIELD,
		fmt.Sprintf("Field '%s' is a duplicate", field))
}

func InvalidArgumentError(msg string) Error {
	return makeError(INVALID_ARGUMENT, msg)
}

func NoSuchElementError() Error {
	return makeError(NO_SUCH_ELEMENT, "")
}

func AssertionFailedError() Error {
	return makeError(ASSERTION_FAILED, "")
}

func ConstSymbolError(name string) Error {
	return makeError(
		CONST_SYMBOL,
		fmt.Sprintf("Symbol '%s' is const", name))
}

func UndefinedSymbolError(name string) Error {
	return makeError(
		UNDEFINIED_SYMBOL,
		fmt.Sprintf("Symbol '%s' is not defined", name))
}
