package crawler

import "sync"

// Fetcher defines the behavior of a URL fetcher.
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// countingFetcher wraps a Fetcher and counts how many times Fetch is called.
type countingFetcher struct {
	fetcher    Fetcher
	fetchCount int
	mu         sync.Mutex
}

// fakeFetcher is a Fetcher that returns predefined results.
type fakeFetcher struct {
	mu   sync.RWMutex
	data map[string]*fakeResult
}

type fakeResult struct {
	body string
	urls []string
}

// CacheResult stores the result of a Fetch operation.
type CacheResult struct {
	Body string
	URLs []string
	Err  error
}
