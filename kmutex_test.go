package kmutex

import (
	"sync"
	"testing"
	"time"
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

	lc := make(chan int)
	uc := make(chan int)
	// Start 10n goroutines accessing n resources 10 times each
	for i := 0; i < 10*number; i++ {
		wg.Add(1)
		go func(k int) {
			for j := 0; j < 10; j++ {
				lc <- k
				km.Lock(ids[k])
				// read and write resource to check for race
				resources[k] = resources[k] + 1
				km.Unlock(ids[k])
				uc <- k
			}
			wg.Done()
		}(i % len(ids))
	}

	to := time.After(time.Second)
	counts := make(map[int]int)
	var lCount, ulCount int
loop:
	for {
		select {
		case k := <-lc:
			counts[k] = counts[k] + 1
			lCount++
		case k := <-uc:
			counts[k] = counts[k] - 1
			ulCount++
		case <-to:
			t.Fatal("timed out waiting for results")
			break loop
		}
		expectCount := 100 * number
		if lCount == expectCount && ulCount == expectCount {
			// Have all results
			break
		}
	}
	for k, c := range counts {
		if c != 0 {
			t.Errorf("Key %d count != 0: %d\n", k, c)
		}
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

	// Verify correct hit count for each resource
	// expecting: (10n hits / n resources) * 10 == 100 hits/resource
	for i := range resources {
		if resources[i] != 100 {
			t.Errorf("resource-%d expected 100 hits, got %d", i, resources[i])
		}
	}
}

func TestLockerInterface(t *testing.T) {
	km := New()
	locker := km.Locker("TEST")
	cond := sync.NewCond(locker)

	if false {
		cond.Wait()
	}
}

func TestCondDeadlock(t *testing.T) {
	l := sync.Mutex{}
	km := WithLock(&l)
	ids := makeIds(10)

	timeout := time.NewTimer(time.Second)
	defer timeout.Stop()

	for checks := 0; checks < 5; checks++ {
		done := make(chan struct{})
		go func() {
			for i := 0; i < len(ids); i++ {
				km.Lock(ids[i])
			}
			close(done)
		}()

		select {
		case <-done:
		case <-timeout.C:
			t.Fatal("timeout while locking all locks")
		}

		var wg, wgReady sync.WaitGroup
		unlocked := make(chan int, len(ids))
		for i := 0; i < len(ids); i++ {
			wg.Add(1)
			wgReady.Add(1)
			go func(k int) {
				wgReady.Done()
				km.Lock(ids[k])
				unlocked <- k
				wg.Done()
			}(i)
		}
		wgReady.Wait()

		km.Unlock(ids[0])
		select {
		case u := <-unlocked:
			if u != 0 {
				t.Fatal("unlocked wrong key, expected 0 but got", u)
			}
		case <-timeout.C:
			t.Fatal("timed out waiting for ids[0] to unlock")
		}

		if !timeout.Stop() {
			<-timeout.C
		}

		for i := 1; i < len(ids); i++ {
			km.Unlock(ids[i])
			u := <-unlocked
			if u != ids[i] {
				t.Fatal("unlocked wrong key, expected", ids[i], "got", u)
			}
		}
		for i := 0; i < len(ids); i++ {
			km.Unlock(ids[i])
		}
		wg.Wait()
		timeout.Reset(time.Second)
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
