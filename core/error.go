// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"strings"
)

//---------------------------------------------------------------
// ErrorKind

type ErrorKind int

const (
	ERROR ErrorKind = iota
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
	IMMUTABLE_VALUE
)

func (t ErrorKind) String() string {
	switch t {

	case ERROR:
		return "Error"
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
	case IMMUTABLE_VALUE:
		return "ImmutableValue"

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
	st   Struct
	str  string
}

func (e *serror) Error() string {
	return e.str
}

func (e *serror) Kind() ErrorKind {
	return e.kind
}

func (e *serror) Struct() Struct {
	return e.st
}

func makeError(kind ErrorKind, msg string) Error {

	var st Struct
	var err Error
	var str string

	if msg == "" {
		st, err = NewStruct([]Field{
			NewField("kind", true, NewStr(kind.String()))}, true)
		str = kind.String()
	} else {
		st, err = NewStruct([]Field{
			NewField("kind", true, NewStr(kind.String())),
			NewField("msg", true, NewStr(msg))}, true)

		str = strings.Join([]string{kind.String(), ": ", msg}, "")
	}
	if err != nil {
		panic("invalid struct")
	}

	return &serror{kind, st, str}
}

func MakeError(kind string, msg string) Error {

	var st Struct
	var err Error
	var str string

	if msg == "" {
		st, err = NewStruct([]Field{
			NewField("kind", true, NewStr(kind))}, true)
		str = kind
	} else {
		st, err = NewStruct([]Field{
			NewField("kind", true, NewStr(kind)),
			NewField("msg", true, NewStr(msg))}, true)

		str = strings.Join([]string{kind, ": ", msg}, "")
	}
	if err != nil {
		panic("invalid struct")
	}

	return &serror{ERROR, st, str}
}

func MakeErrorFromStruct(cx Context, st Struct) Error {

	var str string
	kind, kerr := st.GetField(cx, NewStr("kind"))
	msg, merr := st.GetField(cx, NewStr("msg"))

	if kerr == nil {
		if merr == nil {
			str = strings.Join([]string{
				kind.ToStr(cx).String(), ": ", msg.ToStr(cx).String()}, "")
		} else {
			str = kind.ToStr(cx).String()
		}
	} else {
		str = st.ToStr(cx).String()
	}

	return &serror{ERROR, st, str}
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

func IndexOutOfBoundsError(val int) Error {
	return makeError(
		INDEX_OUT_OF_BOUNDS,
		fmt.Sprintf("%d", val))
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

func ImmutableValueError() Error {
	return makeError(IMMUTABLE_VALUE, "")
}
