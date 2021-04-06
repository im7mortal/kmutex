package kmutex

import (
	"sync"
	"testing"
)

// Number of unique resources to access
const number = 100

func makeIds(count int) []int {
	ids := make([]int, count)
	for i := 0; i < count; i++ {
		ids[i] = i
	}
	return ids
}

func TestKmutex(t *testing.T) {
	km := New()
	ids := makeIds(number)
	resources := make([]int, number)
	wg := sync.WaitGroup{}

	// Start 10n goroutines accessing n resources 10 times each
	for i := 0; i < 10*number; i++ {
		wg.Add(1)
		go func(k int) {
			for j := 0; j < 10; j++ {
				km.Lock(ids[k])
				// read and write resource to check for race
				resources[k] = resources[k] + 1
				km.Unlock(ids[k])
			}
			wg.Done()
		}(i % len(ids))
	}
	wg.Wait()
}

func TestWithLock(t *testing.T) {
	l := sync.Mutex{}
	km := WithLock(&l)
	ids := makeIds(number)
	resources := make([]int, number)
	wg := sync.WaitGroup{}

	// Start 10n goroutines accessing n resources 10 times each
	for i := 0; i < 10*number; i++ {
		wg.Add(1)
		go func(k int) {
			for j := 0; j < 10; j++ {
				km.Lock(ids[k])
				// read and write resource to check for race
				resources[k] = resources[k] + 1
				km.Unlock(ids[k])
			}
			wg.Done()
		}(i % len(ids))
	}
	wg.Wait()
}

func TestLockerInterface(t *testing.T) {
	km := New()
	locker := km.Locker("TEST")
	cond := sync.NewCond(locker)

	if false {
		cond.Wait()
	}
}

func BenchmarkKmutex1000(b *testing.B) {
	km := New()
	ids := makeIds(number)
	resources := make([]int, number)
	wg := sync.WaitGroup{}

	// Start 1000 goroutines accessing 100 resources N times each
	b.ResetTimer()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(k int) {
			for j := 0; j < b.N; j++ {
				km.Lock(ids[k])
				// read and write resource to check for race
				resources[k] = resources[k] + 1
				km.Unlock(ids[k])
			}
			wg.Done()
		}(i % len(ids))
	}
	wg.Wait()
}
