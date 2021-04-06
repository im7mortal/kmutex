/*
Package kmutex is a synchronization primitive that allows locking individual
resources by unique ID.

Key + Mutex = Kmutex

*/
package kmutex

import "sync"

// Can be locked by unique ID
type Kmutex struct {
	c *sync.Cond
	l sync.Locker
	s map[interface{}]struct{}
}

// Create new Kmutex
func New() *Kmutex {
	l := sync.Mutex{}
	return &Kmutex{
		c: sync.NewCond(&l),
		l: &l,
		s: make(map[interface{}]struct{}),
	}
}

// Create new Kmutex with user provided lock
func WithLock(l sync.Locker) *Kmutex {
	return &Kmutex{
		c: sync.NewCond(l),
		l: l,
		s: make(map[interface{}]struct{}),
	}
}

// Unlock Kmutex by unique ID
func (km *Kmutex) Unlock(key interface{}) {
	km.l.Lock()
	defer km.l.Unlock()
	delete(km.s, key)
	km.c.Signal()
}

// Lock Kmutex by unique ID
func (km *Kmutex) Lock(key interface{}) {
	km.l.Lock()
	defer km.l.Unlock()
	for km.locked(key) {
		km.c.Wait()
	}
	km.s[key] = struct{}{}
	return
}

func (km *Kmutex) locked(key interface{}) bool {
	_, ok := km.s[key]
	return ok
}

// satisfy sync.Locker interface
type locker struct {
	km  *Kmutex
	key interface{}
}

// Lock this locker. If already locked, Lock blocks until it is available.
func (l locker) Lock() {
	l.km.Lock(l.key)
}

// Unlock this locker. It is a run-time error if already locked.
func (l locker) Unlock() {
	l.km.Unlock(l.key)
}

// Return a object which implement sync.Locker interface
// A Locker represents an object that can be locked and unlocked.
func (km Kmutex) Locker(key interface{}) sync.Locker {
	return locker{
		key: key,
		km:  &km,
	}
}
