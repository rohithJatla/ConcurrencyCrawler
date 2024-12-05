package crawler

import (
	"fmt"
	"sync"
)

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum depth.
func Crawl(url string, depth int, fetcher Fetcher, cache *Cache, wg *sync.WaitGroup) {
	defer wg.Done()
	if depth <= 0 {
		return
	}
	// Fetch URL using the cache.
	res, err := cache.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, res.Body)
	// Crawl the URLs found on this page.
	var childWG sync.WaitGroup
	for _, u := range res.URLs {
		childWG.Add(1)
		go Crawl(u, depth-1, fetcher, cache, &childWG)
	}
	childWG.Wait()
}
