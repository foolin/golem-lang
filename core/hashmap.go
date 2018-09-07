// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

type (
	// HashMap is an associative array of Values
	HashMap struct {
		buckets [][]*HEntry
		size    int
	}

	// HEntry is an entry in a HashMap
	HEntry struct {
		Key   Value
		Value Value
	}
)

// EmptyHashMap creates an empty HashMap
func EmptyHashMap() *HashMap {
	h, _ := NewHashMap(nil, []*HEntry{})
	return h
}

// NewHashMap creates a HashMap
func NewHashMap(ev Eval, entries []*HEntry) (*HashMap, Error) {
	capacity := 5
	buckets := make([][]*HEntry, capacity)
	hm := &HashMap{buckets, 0}

	for _, e := range entries {
		err := hm.Put(ev, e.Key, e.Value)
		if err != nil {
			return nil, err
		}
	}
	return hm, nil
}

// Eq tests whether two HashMaps are equal
func (hm *HashMap) Eq(ev Eval, that *HashMap) (Bool, Error) {

	if hm.size != that.size {
		return False, nil
	}

	itr := hm.Iterator()
	for itr.Next() {
		entry := itr.Get()

		v, err := that.Get(ev, entry.Key)
		if err != nil {
			return nil, err
		}

		eq, err := entry.Value.Eq(ev, v)
		if err != nil {
			return nil, err
		}

		if !eq.BoolVal() {
			return False, nil
		}
	}

	return True, nil
}

// Get retrieves a value, or returns Null if the value is not present
func (hm *HashMap) Get(ev Eval, key Value) (value Value, err Error) {

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

	b := hm.buckets[hm.lookupBucket(ev, key)]
	n := hm.indexOf(ev, b, key)
	if n == -1 {
		return Null, nil
	}
	return b[n].Value, nil
}

// Contains returns whether the HashMap contains an Entry for the given key
func (hm *HashMap) Contains(ev Eval, key Value) (flag Bool, err Error) {

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

	b := hm.buckets[hm.lookupBucket(ev, key)]
	n := hm.indexOf(ev, b, key)
	if n == -1 {
		return False, nil
	}
	return True, nil
}

// Remove removes the value associated with the given key, if the key
// is present.  Remove returns whether or not the key was present.
func (hm *HashMap) Remove(ev Eval, key Value) (flag Bool, err Error) {

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

	h := hm.lookupBucket(ev, key)
	b := hm.buckets[h]
	n := hm.indexOf(ev, b, key)
	if n == -1 {
		return False, nil
	}
	hm.buckets[h] = append(b[:n], b[n+1:]...)
	hm.size--
	return True, nil
}

// Put adds a new key-value pair to the HashMap
func (hm *HashMap) Put(ev Eval, key Value, value Value) (err Error) {

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

	h := hm.lookupBucket(ev, key)
	n := hm.indexOf(ev, hm.buckets[h], key)
	if n == -1 {
		if hm.tooFull() {
			hm.rehash(ev)
			h = hm.lookupBucket(ev, key)
		}
		hm.buckets[h] = append(hm.buckets[h], &HEntry{key, value})
		hm.size++

	} else {
		hm.buckets[h][n].Value = value
	}

	return nil
}

// Len returns the number of entries in the HashMap
func (hm *HashMap) Len() Int {
	return NewInt(int64(hm.size))
}

//--------------------------------------------------------------

func (hm *HashMap) indexOf(ev Eval, b []*HEntry, key Value) int {
	for i, e := range b {

		eq, err := e.Key.Eq(ev, key)
		if err != nil {
			panic(err)
		}

		if eq.BoolVal() {
			return i
		}
	}
	return -1
}

func (hm *HashMap) tooFull() bool {
	headroom := (hm.size + 1) << 1
	return headroom > len(hm.buckets)
}

func (hm *HashMap) rehash(ev Eval) {
	oldBuckets := hm.buckets

	capacity := len(hm.buckets)<<1 + 1
	hm.buckets = make([][]*HEntry, capacity)
	for _, b := range oldBuckets {
		for _, e := range b {
			h := hm.lookupBucket(ev, e.Key)
			hm.buckets[h] = append(hm.buckets[h], e)
		}
	}
}

func (hm *HashMap) lookupBucket(ev Eval, key Value) int {

	// panic on an un-hashable value
	hc, err := key.HashCode(ev)
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
//
//func (hm *HashMap) dump() {
//	fmt.Println("--------------------------")
//	fmt.Printf("size: %d\n", hm.size)
//	for i, b := range hm.buckets {
//		fmt.Printf("%d, %d: [", i, len(b))
//		for j, e := range b {
//			if j > 0 {
//				fmt.Print(", ")
//			}
//			fmt.Printf("(%v:%v)", e.Key, e.Value)
//		}
//		fmt.Println("]")
//	}
//	fmt.Println("--------------------------")
//}
//
//--------------------------------------------------------------

// Iterator returns an iterator over the entries in the HashMap
func (hm *HashMap) Iterator() *HIterator {
	return &HIterator{hm, -1, -1}
}

// HIterator is an iterator over the entries in the HashMap
type HIterator struct {
	hm        *HashMap
	bucketIdx int
	entryIdx  int
}

// Next advances to the next value in the iterator, if there is one.
// Next returns whether or not there was a value to advance to.
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

// Get returns the current value in the iterator
func (h *HIterator) Get() *HEntry {
	return h.curBucket()[h.entryIdx]
}

func (h *HIterator) curBucket() []*HEntry {
	return h.hm.buckets[h.bucketIdx]
}
