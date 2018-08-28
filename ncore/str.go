// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	"strings"
	"unicode/utf8"
)

type str string

func (s str) String() string {
	return string(s)
}

// NewStr creates a new String
func NewStr(s string) Str {
	return str(s)
}

func (s str) basicMarker() {}

func (s str) Type() Type { return StrType }

func (s str) Freeze(cx Context) (Value, Error) {
	return s, nil
}

func (s str) Frozen(cx Context) (Bool, Error) {
	return True, nil
}

func (s str) ToStr(cx Context) Str { return s }

func (s str) HashCode(cx Context) (Int, Error) {
	h := strHash(string(s))
	return NewInt(int64(h)), nil
}

func (s str) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {

	case str:
		return NewBool(s == t), nil

	default:
		return False, nil
	}
}

func (s str) Cmp(cx Context, v Value) (Int, Error) {
	switch t := v.(type) {

	case str:
		cmp := strings.Compare(string(s), string(t))
		return NewInt(int64(cmp)), nil

	default:
		return nil, TypeMismatchError("Expected Comparable Type")
	}
}

func (s str) Get(cx Context, index Value) (Value, Error) {
	// TODO implement this more efficiently
	runes := []rune(string(s))

	idx, err := boundedIndex(index, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[idx])), nil
}

func (s str) Len(cx Context) Int {
	n := utf8.RuneCountInString(string(s))
	return NewInt(int64(n))
}

func (s str) Slice(cx Context, from Value, to Value) (Value, Error) {
	runes := []rune(string(s))

	f, t, err := sliceIndices(from, to, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[f:t])), nil
}

func (s str) SliceFrom(cx Context, from Value) (Value, Error) {
	runes := []rune(string(s))

	f, _, err := sliceIndices(from, NegOne, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[f:])), nil
}

func (s str) SliceTo(cx Context, to Value) (Value, Error) {
	runes := []rune(string(s))

	_, t, err := sliceIndices(Zero, to, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[:t])), nil
}

func (s str) Concat(that Str) Str {
	a := string(s)
	b := string(that.(str))
	return str(strcpy(a) + strcpy(b))
}

////---------------------------------------------------------------
//// Iterator
//
//type strIterator struct {
//	Struct
//	runes []rune
//	n     int
//}
//
//func (s str) NewIterator(cx Context) Iterator {
//	return initIteratorStruct(cx,
//		&strIterator{newIteratorStruct(), []rune(string(s)), -1})
//}
//
//func (i *strIterator) IterNext() Bool {
//	i.n++
//	return NewBool(i.n < len(i.runes))
//}
//
//func (i *strIterator) IterGet() (Value, Error) {
//
//	if (i.n >= 0) && (i.n < len(i.runes)) {
//		return str([]rune{i.runes[i.n]}), nil
//	}
//	return nil, NoSuchElementError()
//}

//--------------------------------------------------------------
// fields

var strFields = map[string]bool{
	"contains":   true,
	"index":      true,
	"lastIndex":  true,
	"startsWith": true,
	"endsWith":   true,
	"replace":    true,
	"split":      true,
	//"iterator":   true,
}

func (s str) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(strFields))
	for k, _ := range strFields {
		names = append(names, k)
	}
	return names, nil
}

func (s str) HasField(cx Context, val Value) (Bool, Error) {
	if s, ok := val.(Str); ok {
		_, has := strFields[s.String()]
		return NewBool(has), nil
	}
	return nil, TypeMismatchError("Expected Str")
}

func (s str) GetField(cx Context, name Str) (Value, Error) {

	arity, inv, err := s.lookupFunc(name)
	if err != nil {
		return nil, err
	}

	return NewNativeFunc(arity, inv), nil
}

func (s str) InvokeField(cx Context, name Str, params []Value) (Value, Error) {

	_, inv, err := s.lookupFunc(name)
	if err != nil {
		return nil, err
	}

	return inv(cx, params)
}

func (s str) lookupFunc(name Str) (Arity, Invoker, Error) {
	switch name.String() {

	case "contains":
		return Arity{FixedArity, 1, 0}, s.Contains, nil

	default:
		return Arity{}, nil, NoSuchFieldError(name.String())
	}
}

func (s str) Contains(cx Context, params []Value) (Value, Error) {

	err := VetFixedFuncParams([]Type{StrType}, false, params)
	if err != nil {
		return nil, err
	}

	z := params[0].(Str)
	return NewBool(strings.Contains(s.String(), z.String())), nil
}

//func (s str) GetField(cx Context, key Str) (Value, Error) {
//	switch sn := key.String(); sn {
//
//	case "contains":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{StrType}, false,
//			func(cx Context, params []Value) (Value, Error) {
//				z := params[0].(Str)
//				return NewBool(strings.Contains(string(s), z.String())), nil
//			})}, nil
//
//	case "index":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{StrType}, false,
//			func(cx Context, params []Value) (Value, Error) {
//				z := params[0].(Str)
//				return NewInt(int64(strings.Index(string(s), z.String()))), nil
//			})}, nil
//
//	case "lastIndex":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{StrType}, false,
//			func(cx Context, params []Value) (Value, Error) {
//				z := params[0].(Str)
//				return NewInt(int64(strings.LastIndex(string(s), z.String()))), nil
//			})}, nil
//
//	case "startsWith":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{StrType}, false,
//			func(cx Context, params []Value) (Value, Error) {
//				z := params[0].(Str)
//				return NewBool(strings.HasPrefix(string(s), z.String())), nil
//			})}, nil
//
//	case "endsWith":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{StrType}, false,
//			func(cx Context, params []Value) (Value, Error) {
//				z := params[0].(Str)
//				return NewBool(strings.HasSuffix(string(s), z.String())), nil
//			})}, nil
//
//	case "replace":
//		return &virtualFunc{s, sn, NewMultipleNativeFunc(
//			[]Type{StrType, StrType},
//			[]Type{IntType},
//			false,
//			func(cx Context, params []Value) (Value, Error) {
//				a := params[0].(Str)
//				b := params[1].(Str)
//				n := NegOne
//				if len(params) == 3 {
//					n = params[2].(Int)
//				}
//				return NewStr(strings.Replace(
//					string(s),
//					a.String(),
//					b.String(),
//					int(n.IntVal()))), nil
//			})}, nil
//
//	case "split":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{StrType}, false,
//			func(cx Context, params []Value) (Value, Error) {
//				z := params[0].(Str)
//				tokens := strings.Split(string(s), z.String())
//				result := make([]Value, len(tokens))
//				for i, t := range tokens {
//					result[i] = NewStr(t)
//				}
//				return NewList(result), nil
//
//			})}, nil
//
//	case "iterator":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{}, false,
//			func(cx Context, params []Value) (Value, Error) {
//				return s.NewIterator(cx), nil
//			})}, nil
//
//	default:
//		return nil, NoSuchFieldError(key.String())
//	}
//}
