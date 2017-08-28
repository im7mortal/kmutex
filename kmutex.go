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

// satisfy sync.Locker interface
type Locker struct {
	km *Kmutex
	key interface{}
}

// Lock locks m. If the lock is already in use, the calling goroutine blocks until the mutex is available.
func (l Locker) Lock() {
	l.km.Lock(l.key)
}

// Unlock unlocks m. It is a run-time error if m is not locked on entry to Unlock.
func (l Locker) Unlock()  {
	l.km.Unlock(l.key)
}

// Return a object which implement sync.Locker interface
// A Locker represents an object that can be locked and unlocked.
func (km Kmutex) Locker(key interface{}) sync.Locker {
	return Locker{
		key: key,
		km: &km,
	}
}
