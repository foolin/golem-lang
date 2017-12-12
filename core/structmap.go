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
//"fmt"
)

// The hash map implementation that is used by Struct.

type structMap struct {
	buckets [][]*StructEntry
	size    int
}

func newStructMap() *structMap {
	return &structMap{make([][]*StructEntry, 5, 5), 0}
}

// put an entry, but only if it doesn't already exist
func (s *structMap) put(entry *StructEntry) {

	h := s.lookupBucket(entry.Key)
	n := s.indexOf(s.buckets[h], entry.Key)
	if n == -1 {
		if s.tooFull() {
			s.rehash()
			h = s.lookupBucket(entry.Key)
		}
		s.buckets[h] = append(s.buckets[h], entry)
		s.size++
	}
}

func (s *structMap) get(key string) (*StructEntry, bool) {
	b := s.buckets[s.lookupBucket(key)]
	n := s.indexOf(b, key)
	if n == -1 {
		return nil, false
	} else {
		return b[n], true
	}
}

func (s *structMap) keys() []string {
	keys := make([]string, s.size, s.size)
	n := 0
	for _, b := range s.buckets {
		for _, e := range b {
			keys[n] = e.Key
			n++
		}
	}
	return keys
}

//--------------------------------------------------------------

func (s *structMap) indexOf(b []*StructEntry, key string) int {
	for i, e := range b {
		if e.Key == key {
			return i
		}
	}
	return -1
}

func (s *structMap) tooFull() bool {
	headroom := (s.size + 1) << 1
	return headroom > len(s.buckets)
}

func (s *structMap) rehash() {

	oldBuckets := s.buckets
	capacity := len(s.buckets)<<1 + 1
	s.buckets = make([][]*StructEntry, capacity, capacity)
	for _, b := range oldBuckets {
		for _, e := range b {
			h := s.lookupBucket(e.Key)
			s.buckets[h] = append(s.buckets[h], e)
		}
	}
}

func (s *structMap) lookupBucket(key string) int {
	hv := strHash(key)
	if hv < 0 {
		hv = 0 - hv
	}
	return hv % len(s.buckets)
}
