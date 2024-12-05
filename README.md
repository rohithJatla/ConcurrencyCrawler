# Concurrency Crawler

## Introduction
This project demonstrates concurrency in Go using goroutines, mutexes, and various synchronization primitives covered in Chapter 9. It implements a concurrent web crawler that avoids visiting the same URL multiple times, utilizes a concurrent non-blocking cache, and demonstrates lazy initialization.
 It simulates a web crawler that:

- **Fetches web pages concurrently** using goroutines.
- Ensures **thread-safe access** to shared resources with mutexes.
- Implements a **concurrent non-blocking cache** to prevent redundant fetches.
- Demonstrates **lazy initialization** using synchronization primitives.
- Handles **race conditions** and ensures memory synchronization.

## Project Components


1. crawler.go
   Purpose: Contains the Crawl function, which recursively crawls web pages starting from a given URL up to a specified depth.
   
   ### Key Concepts:
   - Goroutines: Each call to Crawl may spawn new goroutines to fetch URLs concurrently.

   - Wait Groups: Uses sync.WaitGroup to wait for all goroutines to finish before exiting.

2. cache.go
   Purpose: Implements a concurrent non-blocking cache that stores fetched URLs and their content to avoid redundant network calls.
   
   ### Key Concepts:
   - Mutexes (sync.Mutex): Protects access to the shared cache map.

   - Lazy Initialization: Uses the ready channel in the entry struct to signal when a fetch operation is complete.

3. types.go
   Purpose: Defines common types and interfaces used across the project.

   ### Components:

   - Fetcher Interface: Represents the behavior required to fetch URLs.

   - CacheResult Struct: Stores the result of a fetch operation.

4. cache_test.go
   Purpose: Contains tests for the cache implementation.

   ### Tests Include:

   - Concurrency Tests: Ensures the cache handles concurrent access correctly.

   - Fetch Count Test: Verifies that each URL is fetched only once, even with concurrent requests.

   - Error Handling Test: Checks the cache's behavior when the Fetcher returns an error.

5. crawler_test.go
   Purpose: Contains tests for the Crawl function.

   ### Tests Include:
   - Functionality Test: Verifies that Crawl correctly fetches pages and discovers links.
---
   - Run `make -B test` to see results.