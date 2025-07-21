package pokecache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Cache holds our cached data, a mutex for concurrent access, and a channel to stop the reapLoop.
type Cache struct {
	cacheMap map[string]cacheEntry
	mu       sync.RWMutex  // mutex to protect the map across goroutines
	stopCh   chan struct{} // Channel to signal the reapLoop to stop
	interval time.Duration // Stores the interval for the reapLoop
}

// cacheEntry represents a single item in the cache.
type cacheEntry struct {
	createdAt time.Time // A time.Time that represents when the entry was created.
	val       []byte    // A []byte that represents the raw data we're caching.
}

// creates a new cache with a configurable interval (time.Duration)
func NewCache(interval time.Duration) (*Cache, error) {
	if interval <= 0 {
		return nil, errors.New("interval must be greater than zero")
	}

	cache := &Cache{ // Use a pointer literal to initialize the struct
		cacheMap: make(map[string]cacheEntry),
		stopCh:   make(chan struct{}), // Initialize the stop channel
		interval: interval,
	}

	// Start the reapLoop in a separate goroutine
	go cache.reapLoop()

	return cache, nil
}

// adds a new entry to the cache
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// create new entry and add to the cache
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.cacheMap[key] = entry
}

// gets an entry from the cache.
// It should take a key (a string) and return a []byte and a bool.
// The bool should be true if the entry was found and false if it wasn't.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cacheEntry, ok := c.cacheMap[key]; ok {
		return cacheEntry.val, true
	} else {
		return nil, false
	}
}

// Stop signals the reapLoop to stop and waits for it to finish.
func (c *Cache) Stop() {
	close(c.stopCh) // Close the stop channel to signal the reapLoop to exit
	// In a more complex scenario, you might want to use a sync.WaitGroup
	// here to wait for the reapLoop goroutine to actually finish before returning.
	// For this example, closing the channel is sufficient to trigger its exit.
	fmt.Println("Cache stop signal sent.")
}

// cache.reapLoop() method that is called when the cache is created (by the NewCache function).
// Each time an interval (the time.Duration passed to NewCache) passes it should remove any entries that are older than the interval.
// This makes sure that the cache doesn't grow too large over time. For example, if the interval is 5 seconds, and an entry was added 7 seconds ago, that entry should be removed.
// I used a time.Ticker to make this happen.
// Maps are not thread-safe in Go.
// You should use a sync.Mutex to lock access to the map when you're adding, getting entries or reaping entries.
// It's unlikely that you'll have issues because reaping only happens every ~5 seconds, but it's still possible, so you should make your cache package safe for concurrent use.
func (c *Cache) reapLoop() {

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	fmt.Printf("reapLoop started with interval: %v\n", c.interval)

	for {
		select {
		case <-ticker.C:
			fmt.Printf("reapLoop 'Tick!' on Cache %v\n", c)
			c.mu.Lock()
			currTime := time.Now()
			for key, cacheEntry := range c.cacheMap {
				if currTime.Sub(cacheEntry.createdAt) >= c.interval {
					fmt.Println("deleting cache entry...")
					delete(c.cacheMap, key)
				}
			}
			c.mu.Unlock()

		case <-c.stopCh:
			// Received stop signal, exit the loop
			fmt.Println("reapLoop received stop signal, exiting.")
			return
		}
	}
}
