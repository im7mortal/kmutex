package kmutex_test

import (
	"sync"
	"time"
	"testing"

	"github.com/im7mortal/kmutex"
)


func TestKmutex(t *testing.T) {
	wg := sync.WaitGroup{}

	km := kmutex.NewKmutex()

	ids := []string{
		"red",
		"blue",
		"yellow",
	}


	ii := 0
	for i := 0; i < 90; i++ {
		wg.Add(1)
		go func(iii int) {
			km.Lock(ids[iii])
			time.Sleep(time.Second)
			println(ii, ids[iii])
			km.Unlock(ids[iii])
			wg.Done()
		}(ii)
		ii++
		if ii == 3 {
			ii = 0
		}
	}
	wg.Wait()

}
