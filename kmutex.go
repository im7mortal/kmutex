/*
Package kmutex is a synchronization primitive that allows locking individual
resources by unique ID.

Key + Mutex = Kmutex

*/
package kmutex

import "sync"

// Can be locked by unique ID
type Kmutex struct {
	l sync.Locker
	s map[interface{}]*klock
}

// klock is a per-key lock that conatins a sync.Cond to signal another
// goroutine that the lock is available, a reference count of the number of
// goroutines waiting for and using the lock, and a boolean to check if the
// lock is already unlocked.
//
// It is necessary to use a separate condition variable for each key to ensure
// that only a goroutine that is waiting for that key is awakened.
type klock struct {
	cond   *sync.Cond
	locked bool
	ref    uint64
}

// Create new Kmutex
func New() *Kmutex {
	l := sync.Mutex{}
	return &Kmutex{
		l: &l,
		s: make(map[interface{}]*klock),
	}
}

// Create new Kmutex with user provided lock
func WithLock(l sync.Locker) *Kmutex {
	return &Kmutex{
		l: l,
		s: make(map[interface{}]*klock),
	}
}

// Unlock Kmutex by unique ID
func (km *Kmutex) Unlock(key interface{}) {
	km.l.Lock()
	defer km.l.Unlock()
	kl, ok := km.s[key]
	if !ok || !kl.locked {
		panic("unlock of unlocked kmutex")
	}
	kl.ref--
	if kl.ref == 0 {
		delete(km.s, key)
		return
	}
	kl.locked = false
	kl.cond.Signal()
}

// Lock Kmutex by unique ID. Returns if waited for an Unlock.
func (km *Kmutex) Lock(key interface{}) (waited bool) {
	km.l.Lock()
	defer km.l.Unlock()
	for {
		kl, ok := km.s[key]
		if !ok {
			km.s[key] = &klock{
				cond:   sync.NewCond(km.l),
				locked: true,
				ref:    1,
			}
			return false
		}

		kl.ref++
		kl.cond.Wait()
		// No need to check if locked, since signal only given on unlock and
		// only delivered to one goroutine.
		kl.locked = true
		return true
	}
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
