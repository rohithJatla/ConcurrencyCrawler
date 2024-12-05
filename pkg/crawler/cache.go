package crawler

import (
	"sync"
)

// Cache is a concurrent non-blocking cache.
type Cache struct {
	fetcher Fetcher
	mu      sync.Mutex
	cache   map[string]*entry
}

type entry struct {
	res   CacheResult
	ready chan struct{}
}

// NewCache creates a new Cache instance.
func NewCache(fetcher Fetcher) *Cache {
	return &Cache{
		fetcher: fetcher,
		cache:   make(map[string]*entry),
	}
}

// Get fetches the URL either from the cache or by using the fetcher.
func (c *Cache) Get(url string) (CacheResult, error) {
	c.mu.Lock()
	e := c.cache[url]
	if e != nil {
		// Entry exists; wait for the result to be ready.
		c.mu.Unlock()
		<-e.ready
	} else {
		// Lazy initialization of the entry.
		e = &entry{ready: make(chan struct{})}
		c.cache[url] = e
		c.mu.Unlock()

		// Fetch the data.
		body, urls, err := c.fetcher.Fetch(url)
		e.res = CacheResult{Body: body, URLs: urls, Err: err}
		close(e.ready)
	}
	return e.res, e.res.Err
}
