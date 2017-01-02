package bigpipe

import (
	"sync"
)

// LookupFunc is sample template for cache lookup. It is de-dupe across multiple goroutines.
// Example is http get where url is the key:
//	res, err := http.Get(url)
//	if err != nil {
//	return nil, err
//	}
//	defer res.Body.Close()
//	return ioutil.ReadAll(res.Body)
type LookupFunc func(key string) (interface{}, error)

// CacheContainerGenerator is used by application to generate cache for its pagelets.
type CacheContainerGenerator func(f LookupFunc)

type result struct {
	value interface{}
	err error
}

type entry struct {
	res result
	awaitResChannel chan struct{}
}

// CacheContainer is a cache designed to work for multiple goroutines.
// Its main purpose is to support de-duping network request in pagelets.
// It is supposed to be contained in request scope. Each requests needs to have a new CacheContainer.
type CacheContainer struct {
	f     LookupFunc
	mu    sync.Mutex
	cache map[string]*entry
}

func newCache(cacheContainer *CacheContainer) CacheContainerGenerator {
	cache := make(map[string] *entry)
	return func(f LookupFunc) {
		*cacheContainer = CacheContainer{f: f, cache: cache}
	}
}

// GetValueForKey is a simple function for key lookup to be used by pagelets.
// Calls to this function are deduped (in multiple goroutines).
// Pagelets runs in their own goroutines and it ensures the above for same.
func (cacheContainer *CacheContainer) GetValueForKey(key string) (value interface{}, err error) {

	cacheContainer.mu.Lock()

	if cacheContainer.f == nil {
		panic("No cache implemented")
	}
	e := cacheContainer.cache[key]
	if e == nil {
		e = &entry{awaitResChannel: make(chan struct{})}
		cacheContainer.cache[key] = e
		cacheContainer.mu.Unlock()
		e.res.value, e.res.err = cacheContainer.f(key)
		close(e.awaitResChannel)
	} else {
		cacheContainer.mu.Unlock()
		<- e.awaitResChannel
	}
	return e.res.value, e.res.err
}

