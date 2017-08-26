package kmutex

import "sync"

type Kmutex struct {
	m sync.Map
}

func NewKmutex() Kmutex {
	return Kmutex{sync.Map{}}
}

func (s Kmutex) Unlock(key interface{})  {
	l, _ := s.m.Load(key)
	l_ := l.(*sync.Mutex)
	s.m.Delete(key)
	l_.Unlock()
}

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
