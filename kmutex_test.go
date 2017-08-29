package kmutex_test

import (
	"sync"
	"testing"

	"github.com/im7mortal/kmutex"
)

const number  = 100

func TestKmutex(t *testing.T) {
	wg := sync.WaitGroup{}

	km := kmutex.New()

	ids := []int{}

	for i := 0; i < number; i++ {
		ids = append(ids, i)
	}

	ii := 0
	for i := 0; i < number * number; i++ {
		wg.Add(1)
		go func(iii int) {
			km.Lock(ids[iii])
			km.Unlock(ids[iii])
			wg.Done()
		}(ii)
		ii++
		if ii == number {
			ii = 0
		}
	}
	wg.Wait()
}

func TestWithLock(t *testing.T) {
	wg := sync.WaitGroup{}
	l := sync.Mutex{}
	km := kmutex.WithLock(&l)

	ids := []int{}

	for i := 0; i < number; i++ {
		ids = append(ids, i)
	}

	ii := 0
	for i := 0; i < number * number; i++ {
		wg.Add(1)
		go func(iii int) {
			km.Lock(ids[iii])
			km.Unlock(ids[iii])
			wg.Done()
		}(ii)
		ii++
		if ii == number {
			ii = 0
		}
	}
	wg.Wait()
}


func TestLockerInterface(t *testing.T) {
	km := kmutex.New()

	locker := km.Locker("TEST")

	cond := sync.NewCond(locker)

	if false {
		cond.Wait()
	}
}
