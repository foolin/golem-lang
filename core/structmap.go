// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
//"fmt"
)

// The hash map implementation that is used by Struct.
//
// TODO there is an opportunity to improve the performance of this
// data structure considerably.  Since the number of entries in a struct
// cannot be changed once it is created, we can create hash tables that
// hash on an explicit stated prime number. The go compiler should be smart
// enough to generate much more efficient code that way.

type structMap struct {
	buckets [][]*field
	size    int
}

func newStructMap() *structMap {
	return &structMap{make([][]*field, 5, 5), 0}
}

// put a field, but only if it doesn't already exist
func (s *structMap) put(field *field) {

	h := s._lookupBucket(field.name)
	n := s._indexOf(s.buckets[h], field.name)
	if n == -1 {
		if s._tooFul() {
			s._rehash()
			h = s._lookupBucket(field.name)
		}
		s.buckets[h] = append(s.buckets[h], field)
		s.size++
	}
}

func (s *structMap) get(name string) (*field, bool) {
	b := s.buckets[s._lookupBucket(name)]
	n := s._indexOf(b, name)
	if n == -1 {
		return nil, false
	} else {
		return b[n], true
	}
}

func (s *structMap) fieldNames() []string {
	fieldNames := make([]string, s.size, s.size)
	n := 0
	for _, b := range s.buckets {
		for _, f := range b {
			fieldNames[n] = f.name
			n++
		}
	}
	return fieldNames
}

//--------------------------------------------------------------
// these are internal methods -- don't call them directly

func (s *structMap) _indexOf(b []*field, name string) int {
	for i, f := range b {
		if f.name == name {
			return i
		}
	}
	return -1
}

func (s *structMap) _tooFul() bool {
	headroom := (s.size + 1) << 1
	return headroom > len(s.buckets)
}

func (s *structMap) _rehash() {

	oldBuckets := s.buckets
	capacity := len(s.buckets)<<1 + 1
	s.buckets = make([][]*field, capacity, capacity)
	for _, b := range oldBuckets {
		for _, f := range b {
			h := s._lookupBucket(f.name)
			s.buckets[h] = append(s.buckets[h], f)
		}
	}
}

func (s *structMap) _lookupBucket(name string) int {
	hv := strHash(name)
	if hv < 0 {
		hv = 0 - hv
	}
	return hv % len(s.buckets)
}
