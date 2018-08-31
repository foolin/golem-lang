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

func (s str) Freeze(ev Evaluator) (Value, Error) {
	return s, nil
}

func (s str) Frozen(ev Evaluator) (Bool, Error) {
	return True, nil
}

func (s str) ToStr(ev Evaluator) Str { return s }

func (s str) HashCode(ev Evaluator) (Int, Error) {
	h := strHash(string(s))
	return NewInt(int64(h)), nil
}

func (s str) Eq(ev Evaluator, v Value) (Bool, Error) {
	switch t := v.(type) {

	case str:
		return NewBool(s == t), nil

	default:
		return False, nil
	}
}

func (s str) Cmp(ev Evaluator, v Value) (Int, Error) {
	switch t := v.(type) {

	case str:
		cmp := strings.Compare(string(s), string(t))
		return NewInt(int64(cmp)), nil

	default:
		return nil, TypeMismatchError("Expected Comparable Type")
	}
}

func (s str) Get(ev Evaluator, index Value) (Value, Error) {
	// TODO implement this more efficiently
	runes := []rune(string(s))

	idx, err := boundedIndex(index, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[idx])), nil
}

func (s str) Len(ev Evaluator) Int {
	n := utf8.RuneCountInString(string(s))
	return NewInt(int64(n))
}

func (s str) Slice(ev Evaluator, from Value, to Value) (Value, Error) {
	runes := []rune(string(s))

	f, t, err := sliceIndices(from, to, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[f:t])), nil
}

func (s str) SliceFrom(ev Evaluator, from Value) (Value, Error) {
	runes := []rune(string(s))

	f, _, err := sliceIndices(from, NegOne, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[f:])), nil
}

func (s str) SliceTo(ev Evaluator, to Value) (Value, Error) {
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
//func (s str) NewIterator(ev Evaluator) Iterator {
//	return initIteratorStruct(ev,
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

var strMethods = map[string]Method{

	//"contains",
	//"index",

	//"lastIndex",
	//"startsWith",
	//"endsWith",
	//"replace",
	//"split",
	//"iterator",
}

func (s str) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(strMethods))
	for name, _ := range strMethods {
		names = append(names, name)
	}
	return names, nil
}

func (s str) HasField(name string) (bool, Error) {
	_, ok := strMethods[name]
	return ok, nil
}

func (s str) GetField(name string, ev Evaluator) (Value, Error) {
	if method, ok := strMethods[name]; ok {
		return method.ToFunc(s, name), nil
	}
	return nil, NoSuchFieldError(name)
}

func (s str) InvokeField(name string, ev Evaluator, params []Value) (Value, Error) {

	if method, ok := strMethods[name]; ok {
		return method.Invoke(s, ev, params)
	}
	return nil, NoSuchFieldError(name)
}

//var containsMethod Method = nil //NewFixedMethod(
//	StrType, []Type{StrType}, false,
//	func(ev Evaluator, self Value, params []Value) (Value, Error) {
//		s := self.(Str)
//		substr := params[0].(Str)
//		return NewBool(strings.Contains(s.String(), substr.String())), nil
//	}), true

//var containsMethod = NewFixedMethod(
//	[]Type{StrType}, false,
//	func(ev Evaluator, self Value, params []Value) Invoke {
//		return func(ev Evaluator, params []Value) (Value, Error) {
//		return self.Contains(params[0].(Str)), nil
//	})

//	case "index":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{StrType}, false,
//			func(ev Evaluator, params []Value) (Value, Error) {
//				z := params[0].(Str)
//				return NewInt(int64(strings.Index(string(s), z.String()))), nil
//			})}, nil

//--------------------------------------------------------------

//func (s str) GetField(ev Evaluator, key Str) (Value, Error) {
//	switch sn := key.String(); sn {
//
//	case "contains":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{StrType}, false,
//			func(ev Evaluator, params []Value) (Value, Error) {
//				z := params[0].(Str)
//				return NewBool(strings.Contains(string(s), z.String())), nil
//			})}, nil
//
//	case "index":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{StrType}, false,
//			func(ev Evaluator, params []Value) (Value, Error) {
//				z := params[0].(Str)
//				return NewInt(int64(strings.Index(string(s), z.String()))), nil
//			})}, nil
//
//	case "lastIndex":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{StrType}, false,
//			func(ev Evaluator, params []Value) (Value, Error) {
//				z := params[0].(Str)
//				return NewInt(int64(strings.LastIndex(string(s), z.String()))), nil
//			})}, nil
//
//	case "startsWith":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{StrType}, false,
//			func(ev Evaluator, params []Value) (Value, Error) {
//				z := params[0].(Str)
//				return NewBool(strings.HasPrefix(string(s), z.String())), nil
//			})}, nil
//
//	case "endsWith":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{StrType}, false,
//			func(ev Evaluator, params []Value) (Value, Error) {
//				z := params[0].(Str)
//				return NewBool(strings.HasSuffix(string(s), z.String())), nil
//			})}, nil
//
//	case "replace":
//		return &virtualFunc{s, sn, NewMultipleNativeFunc(
//			[]Type{StrType, StrType},
//			[]Type{IntType},
//			false,
//			func(ev Evaluator, params []Value) (Value, Error) {
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
//			func(ev Evaluator, params []Value) (Value, Error) {
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
//			func(ev Evaluator, params []Value) (Value, Error) {
//				return s.NewIterator(ev), nil
//			})}, nil
//
//	default:
//		return nil, NoSuchFieldError(key.String())
//	}
//}
