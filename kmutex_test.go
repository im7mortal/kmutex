package kmutex_test

import (
	"sync"
	"testing"

	"github.com/im7mortal/kmutex"
)

const number  = 100

func TestKmutex(t *testing.T) {
	wg := sync.WaitGroup{}

	km := kmutex.NewKmutex()

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
