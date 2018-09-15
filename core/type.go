// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

// Type represents all of the possible types of a Golem Value
type Type int

// All possible kinds of Type
const (
	// AnyType is a special type that we use when we don't care what the actual type is.
	// It is used in conjuction with type-checking for certain kinds of function invocations.
	AnyType Type = iota

	NullType
	BoolType
	IntType
	FloatType
	StrType
	ListType
	RangeType
	TupleType
	DictType
	SetType
	StructType
	FuncType
	ChanType
)

func (t Type) String() string {

	switch t {
	case NullType:
		return "Null"
	case BoolType:
		return "Bool"
	case IntType:
		return "Int"
	case FloatType:
		return "Float"
	case StrType:
		return "Str"
	case ListType:
		return "List"
	case RangeType:
		return "Range"
	case TupleType:
		return "Tuple"
	case DictType:
		return "Dict"
	case SetType:
		return "Set"
	case StructType:
		return "Struct"
	case FuncType:
		return "Func"
	case ChanType:
		return "Chan"

	default:
		panic("unreachable")
	}
}
