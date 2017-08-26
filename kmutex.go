// Package kmutex is synchronization primitive. Mutex which can be locked by unique ID.
package kmutex

import "sync"

// Can be locked by unique ID
type Kmutex struct {
	m *sync.Map
}

// Create new Kmutex
func NewKmutex() Kmutex {
	m := sync.Map{}
	return Kmutex{&m}
}

// Unlock Kmutex by unique ID
func (s Kmutex) Unlock(key interface{})  {
	l, exist := s.m.Load(key)
	if !exist {
		panic("kmutex: unlock of unlocked mutex")
	}
	l_ := l.(*sync.Mutex)
	s.m.Delete(key)
	l_.Unlock()
}

// Lock Kmutex by unique ID
func (s Kmutex) Lock(key interface{}) {
	m := sync.Mutex{}
	m_, _ := s.m.LoadOrStore(key, &m)
	mm := m_.(*sync.Mutex)
	mm.Lock()
	if mm != &m {
		mm.Unlock()
		s.Lock(key)
		return
	}
	return
}
