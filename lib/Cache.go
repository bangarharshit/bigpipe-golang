package bigpipe

import (
	"sync"
)

type LookupFunc func(key string) (interface{}, error)

type CacheContainerGenerator func(f LookupFunc)

type result struct {
	value interface{}
	err error
}

type entry struct {
	res result
	awaitResChannel chan struct{}
}

type CacheContainer struct {
	f     LookupFunc
	mu    sync.Mutex
	cache map[string]*entry
}

func newCache(cacheContainer *CacheContainer) CacheContainerGenerator {
	mu := sync.Mutex{}
	cache := make(map[string] *entry)
	return func(f LookupFunc) {
		*cacheContainer = CacheContainer{f, mu, cache}
	}
}

func (cacheContainer *CacheContainer) GetValueForKey(key string) (value interface{}, err error) {

	if cacheContainer.f == nil {
		panic("No cache implemented")
	}

	cacheContainer.mu.Lock()
	e := cacheContainer.cache[key]
	if e == nil {
		e = &entry{awaitResChannel: make(chan struct{})}
		cacheContainer.cache[key] = e
		cacheContainer.mu.Unlock()
		e.res.value, e.res.err = cacheContainer.f(key)
		close(cacheContainer.cache[key].awaitResChannel)
	} else {
		cacheContainer.mu.Unlock()
		<- cacheContainer.cache[key].awaitResChannel
	}
	return e.res.value, e.res.err
}

