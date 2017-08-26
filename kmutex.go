package kmutex

import "sync"

type Kmutex struct {
	mut *sync.Mutex
	m   map[interface{}]*sync.Mutex
}

func NewKmutex() Kmutex {
	m := sync.Mutex{}
	return Kmutex{&m, map[interface{}]*sync.Mutex{}}
}

func (s Kmutex) Unlock(key interface{}) {
	s.mut.Lock()
	defer s.mut.Unlock()
	m := s.m[key]
	delete(s.m, key)
	m.Unlock()
}

func (s Kmutex) Lock(key interface{}) {
	m := sync.Mutex{}
	s.mut.Lock()
	m_, busy := s.m[key]
	if busy {
		s.mut.Unlock()
		m_.Lock()
		m_.Unlock()
		s.Lock(key)
		return
	}
	m.Lock()
	s.m[key] = &m
	s.mut.Unlock()
	return
}
