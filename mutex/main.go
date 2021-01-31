package main

import (
	"fmt"
	"sync"
)

// AtomicInt struct
type AtomicInt struct {
	value int
	lock  sync.Mutex
}

var (
	counter       = 0
	lock          sync.Mutex
	atomicCounter = AtomicInt{}
)

func (a *AtomicInt) increase() {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.value++
}

func (a *AtomicInt) decrease() {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.value--
}

func (a *AtomicInt) getValue() int {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.value
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go updateCounter(&wg)
	}
	wg.Wait()
	fmt.Printf("final counter value: %d\n", counter)
	fmt.Printf("final atomic counter value: %d\n", atomicCounter.value)
}

func updateCounter(wg *sync.WaitGroup) {
	lock.Lock()
	defer lock.Unlock()

	counter++
	atomicCounter.value++
	wg.Done()
}
