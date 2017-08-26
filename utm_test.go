package kmutex_test

import (
	"sync"
	"testing"

	"github.com/im7mortal/kmutex"
)


func TestKmutex(t *testing.T) {
	wg := sync.WaitGroup{}

	km := kmutex.NewKmutex()

	id := "sfdjs839jnfd"


	//m := sync.Mutex{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			w := km.Acquire(id)
			println(i)
			w.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()

}
