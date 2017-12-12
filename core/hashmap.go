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

// The hash map implementation that is used by Dict.

type (
	HashMap struct {
		buckets [][]*HEntry
		size    int
	}

	HEntry struct {
		Key   Value
		Value Value
	}
)

func EmptyHashMap() *HashMap {
	return NewHashMap([]*HEntry{})
}

func NewHashMap(entries []*HEntry) *HashMap {
	capacity := 5
	buckets := make([][]*HEntry, capacity, capacity)
	hm := &HashMap{buckets, 0}

	for _, e := range entries {
		hm.Put(e.Key, e.Value)
	}
	return hm
}

func (hm *HashMap) Get(key Value) (value Value, err Error) {

	// recover from an un-hashable value
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(Error); ok {
				value = nil
				err = e
			} else {
				panic(r)
			}
		}
	}()

	b := hm.buckets[hm.lookupBucket(key)]
	n := hm.indexOf(b, key)
	if n == -1 {
		return NULL, nil
	} else {
		return b[n].Value, nil
	}
}

func (hm *HashMap) ContainsKey(key Value) (flag Bool, err Error) {

	// recover from an un-hashable value
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(Error); ok {
				flag = nil
				err = e
			} else {
				panic(r)
			}
		}
	}()

	b := hm.buckets[hm.lookupBucket(key)]
	n := hm.indexOf(b, key)
	if n == -1 {
		return FALSE, nil
	} else {
		return TRUE, nil
	}
}

func (hm *HashMap) Put(key Value, value Value) (err Error) {

	// recover from an un-hashable value
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	h := hm.lookupBucket(key)
	n := hm.indexOf(hm.buckets[h], key)
	if n == -1 {
		if hm.tooFull() {
			hm.rehash()
			h = hm.lookupBucket(key)
		}
		hm.buckets[h] = append(hm.buckets[h], &HEntry{key, value})
		hm.size++

	} else {
		hm.buckets[h][n].Value = value
	}

	return nil
}

func (hm *HashMap) Len() Int {
	return MakeInt(int64(hm.size))
}

//--------------------------------------------------------------

func (hm *HashMap) indexOf(b []*HEntry, key Value) int {
	for i, e := range b {

		if e.Key.Eq(key).BoolVal() {
			return i
		}
	}
	return -1
}

func (hm *HashMap) tooFull() bool {
	headroom := (hm.size + 1) << 1
	return headroom > len(hm.buckets)
}

func (hm *HashMap) rehash() {
	oldBuckets := hm.buckets

	capacity := len(hm.buckets)<<1 + 1
	hm.buckets = make([][]*HEntry, capacity, capacity)
	for _, b := range oldBuckets {
		for _, e := range b {
			h := hm.lookupBucket(e.Key)
			hm.buckets[h] = append(hm.buckets[h], e)
		}
	}
}

func (hm *HashMap) lookupBucket(key Value) int {

	// panic on an un-hashable value
	hc, err := key.HashCode()
	if err != nil {
		panic(err)
	}

	hv := int(hc.IntVal())
	if hv < 0 {
		hv = 0 - hv
	}

	return hv % len(hm.buckets)
}

//--------------------------------------------------------------

func (hm *HashMap) Iterator() *HIterator {
	return &HIterator{hm, -1, -1}
}

type HIterator struct {
	hm        *HashMap
	bucketIdx int
	entryIdx  int
}

func (h *HIterator) Next() bool {

	// advance to next entry in current []*HEntry
	h.entryIdx++

	// if we are not pointing at a valid entry
	if (h.bucketIdx == -1) || (h.entryIdx >= len(h.curBucket())) {

		// then advance to next non-empty []*HEntry
		h.bucketIdx++
		for (h.bucketIdx < len(h.hm.buckets)) && (len(h.curBucket()) == 0) {
			h.bucketIdx++
		}
		if !(h.bucketIdx < len(h.hm.buckets)) {
			return false
		}

		// and point at first entry of the new []*HEntry
		h.entryIdx = 0
	}

	return true
}

func (h *HIterator) Get() *HEntry {
	return h.curBucket()[h.entryIdx]
}

func (h *HIterator) curBucket() []*HEntry {
	return h.hm.buckets[h.bucketIdx]
}
