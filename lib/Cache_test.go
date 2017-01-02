package bigpipe

import (
	"time"
	"sync"
	"testing"
)

func TestCacheContainer_GetValueForKey(t *testing.T) {
	var noOfLookups int = 0
	var mu sync.Mutex = sync.Mutex{}
	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(4)
	cacheContainer := &CacheContainer{func(key string) (interface{}, error) {
		defer wg.Done()
		mu.Lock()
		noOfLookups++
		mu.Unlock()
		time.Sleep(50 * time.Millisecond)
		return key, nil
	}, sync.Mutex{}, make(map[string]*entry)}
	keys := []string{"key1", "key2", "key3", "key4", "key4", "key1"}
	for _, key := range keys {
		go cacheContainer.GetValueForKey(key)
	}
	wg.Wait()
	if noOfLookups != 4 {
		t.Errorf("expected 4 lookups, found %d lookups", noOfLookups)
	}
}


