package kmutex

import "sync"

type storage struct {
	m sync.Map
}

type Locker struct {
	key interface{}
	l   *sync.Mutex
	s   *sync.Map
}

func NewKmutex() storage {
	return storage{sync.Map{}}
}

func (l Locker) Unlock() {
	l.s.Delete(l.key)
	l.l.Unlock()
}

func (s storage) Acquire(key interface{}) Locker {
	m := sync.Mutex{}
	m_, _ := s.m.LoadOrStore(key, &m)
	mm := m_.(*sync.Mutex)
	mm.Lock()
	if mm != &m {
		mm.Unlock()
		return s.Acquire(key)
	}
	return Locker{l: &m, s: &s.m, key: key}
}
