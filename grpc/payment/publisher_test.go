package main

import (
	"sync"
	"testing"
)

func TestPublish(t *testing.T) {
	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			PublishMessage(i)
		}(wg, i)
	}
	wg.Wait()
}
