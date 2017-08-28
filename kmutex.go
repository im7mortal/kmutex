// Package kmutex is synchronization primitive. Mutex which can be locked by unique ID.
// Key + Mutex = Kmutex
package kmutex

import "sync"

// Can be locked by unique ID
type Kmutex struct {
	c *sync.Cond
	l sync.Locker
	s map[interface{}]struct{}
}

// Create new Kmutex
func NewKmutex() *Kmutex {
	l := sync.Mutex{}
	return &Kmutex{c: sync.NewCond(&l), l: &l, s: make(map[interface{}]struct{})}
}

// Create new Kmutex with user provided lock
func NewKmutexWithLock(l sync.Locker) *Kmutex {
	return &Kmutex{c: sync.NewCond(l), l: l, s: make(map[interface{}]struct{})}
}

func (km *Kmutex) locked(key interface{}) (ok bool) { _, ok = km.s[key]; return }

// Unlock Kmutex by unique ID
func (km *Kmutex) Unlock(key interface{})  {
	km.l.Lock()
	defer km.l.Unlock()
	delete(km.s, key)
	km.c.Broadcast()
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
