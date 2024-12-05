package crawler

import (
	"fmt"
	"sync"
	"testing"
)

func TestCrawl(t *testing.T) {
	fetcher := &fakeFetcher{
		data: map[string]*fakeResult{
			"https://golang.org/": {
				body: "The Go Programming Language",
				urls: []string{
					"https://golang.org/pkg/",
					"https://golang.org/cmd/",
				},
			},
			"https://golang.org/pkg/": {
				body: "Packages",
				urls: []string{
					"https://golang.org/",
					"https://golang.org/cmd/",
				},
			},
			"https://golang.org/cmd/": {
				body: "Commands",
				urls: []string{
					"https://golang.org/",
					"https://golang.org/pkg/",
				},
			},
		},
	}
	cache := NewCache(fetcher)
	var wg sync.WaitGroup
	wg.Add(1)
	Crawl("https://golang.org/", 4, fetcher, cache, &wg)
	wg.Wait()
}

func (f *fakeFetcher) Fetch(url string) (string, []string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if res, ok := f.data[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}
