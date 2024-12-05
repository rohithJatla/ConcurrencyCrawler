package crawler

import (
	"sync"
	"testing"
	"time"
)

func TestCacheConcurrency(t *testing.T) {
	fetcher := &fakeFetcher{
		data: map[string]*fakeResult{
			"http://example.com": {
				body: "Example Domain",
				urls: []string{},
			},
		},
	}
	cache := NewCache(fetcher)
	var wg sync.WaitGroup

	// Simulate concurrent access to the cache.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := cache.Get("http://example.com")
			if err != nil || res.Body != "Example Domain" {
				t.Errorf("Failed to fetch or incorrect data")
			}
		}()
	}
	wg.Wait()
}

func TestCacheFetchOnce(t *testing.T) {
	baseFetcher := &fakeFetcher{
		data: map[string]*fakeResult{
			"http://example.com": {
				body: "Example Domain",
				urls: []string{},
			},
		},
	}

	// Use countingFetcher to count fetches
	cf := &countingFetcher{
		fetcher: baseFetcher,
	}

	cache := NewCache(cf)

	var wg sync.WaitGroup
	// Simulate concurrent access to the same URL
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := cache.Get("http://example.com")
			if err != nil || res.Body != "Example Domain" {
				t.Errorf("Failed to fetch or incorrect data")
			}
		}()
	}
	wg.Wait()

	// Check that fetch was called only once
	cf.mu.Lock()
	fetchCount := cf.fetchCount
	cf.mu.Unlock()
	if fetchCount != 1 {
		t.Errorf("Expected fetch to be called once, but got %d", fetchCount)
	}
}

func TestCacheMultipleURLs(t *testing.T) {
	fetcher := &fakeFetcher{
		data: map[string]*fakeResult{
			"http://example.com": {
				body: "Example Domain",
				urls: []string{},
			},
			"http://example.org": {
				body: "Example Organization",
				urls: []string{},
			},
			"http://example.net": {
				body: "Example Network",
				urls: []string{},
			},
		},
	}
	cache := NewCache(fetcher)

	var wg sync.WaitGroup
	urls := []string{
		"http://example.com",
		"http://example.org",
		"http://example.net",
	}
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			res, err := cache.Get(u)
			if err != nil {
				t.Errorf("Failed to fetch %s: %v", u, err)
			}
			expected := fetcher.data[u].body
			if res.Body != expected {
				t.Errorf("Expected body %q for %s, got %q", expected, u, res.Body)
			}
		}(url)
	}
	wg.Wait()
}

func TestCacheErrorHandling(t *testing.T) {
	fetcher := &fakeFetcher{
		data: map[string]*fakeResult{},
	}
	cache := NewCache(fetcher)

	// Attempt to fetch a URL that doesn't exist in the fetcher
	res, err := cache.Get("http://nonexistent.com")
	if err == nil {
		t.Errorf("Expected an error for nonexistent URL, got result: %v", res)
	}
}

func (cf *countingFetcher) Fetch(url string) (string, []string, error) {
	cf.mu.Lock()
	cf.fetchCount++
	cf.mu.Unlock()
	time.Sleep(100 * time.Millisecond) // Simulate network delay
	return cf.fetcher.Fetch(url)
}
