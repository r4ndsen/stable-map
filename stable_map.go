package stable_map

import (
	"iter"
	"slices"
	"sync"
)

type entry[K comparable, V any] struct {
	key     K
	value   V
	deleted bool
}

func newEntry[K comparable, V any](key K, value V) *entry[K, V] {
	return &entry[K, V]{
		key:   key,
		value: value,
	}
}

func NewStableMap[K comparable, V any]() *StableMap[K, V] {
	sm := new(StableMap[K, V])
	sm.init()

	return sm
}

type StableMap[K comparable, V any] struct {
	sync.RWMutex

	entries []*entry[K, V]
	refMap  map[K]*entry[K, V]

	len     int
	deleted int
}

func (sm *StableMap[K, V]) init() {
	sm.entries = make([]*entry[K, V], 0)
	sm.refMap = make(map[K]*entry[K, V])
	sm.len = 0
	sm.deleted = 0
}

func (sm *StableMap[K, V]) delete() {
	sm.len--
	sm.deleted++

	if sm.deleted >= 100_000 {
		entries := make([]*entry[K, V], 0, len(sm.entries)-sm.deleted)

		sm.deleted = 0

		for _, e := range sm.entries {
			if e.deleted {
				continue
			}
			entries = append(entries, e)
		}

		sm.entries = entries
	}
}

func (sm *StableMap[K, V]) Delete(key K) {
	sm.Lock()
	defer sm.Unlock()

	if e, ok := sm.refMap[key]; ok {
		var zeroValue V

		e.deleted = true
		e.value = zeroValue
		delete(sm.refMap, key)
		sm.delete()
	}
}

func (sm *StableMap[K, V]) Len() int {
	return sm.len
}

func (sm *StableMap[K, V]) Reset() {
	sm.Lock()
	defer sm.Unlock()

	sm.init()
}

func (sm *StableMap[K, V]) Set(key K, value V) {
	sm.Lock()
	defer sm.Unlock()

	if e, ok := sm.refMap[key]; ok && e.deleted == false {
		e.value = value
		return
	}

	e := newEntry(key, value)

	sm.len++
	sm.entries = append(sm.entries, e)
	sm.refMap[key] = e
}

func (sm *StableMap[K, V]) Get(key K) (V, bool) {
	sm.RLock()
	defer sm.RUnlock()

	var found bool
	var res V

	if e, ok := sm.refMap[key]; ok {
		res = e.value
		found = true
	}

	return res, found
}

func (sm *StableMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, e := range sm.entries {
			if e.deleted {
				continue
			}
			if !yield(e.key, e.value) {
				break
			}
		}
	}
}

func (sm *StableMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, e := range sm.entries {
			if e.deleted {
				continue
			}
			if !yield(e.key) {
				break
			}
		}
	}
}

func (sm *StableMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, e := range sm.entries {
			if e.deleted {
				continue
			}
			if !yield(e.value) {
				break
			}
		}
	}
}

func SequenceToSlice[V any](seq iter.Seq[V]) []V {
	return slices.Collect(seq)
}
