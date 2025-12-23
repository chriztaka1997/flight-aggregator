package cache

import (
	"crypto/sha256"
	"encoding/json"
	"flight-aggregator/internal/models"
	"fmt"
	"sync"
	"time"
)

// CacheEntry represents a cached item with expiration
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

// IsExpired checks if the cache entry has expired
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache represents an in-memory cache with TTL
type Cache struct {
	data  map[string]*CacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
	stats Stats
}

// Stats tracks cache performance metrics
type Stats struct {
	Hits          int64
	Misses        int64
	Evictions     int64
	Expirations   int64
	CurrentSize   int
	TotalRequests int64
}

// New creates a new cache with the specified TTL
func New(ttl time.Duration) *Cache {
	c := &Cache{
		data: make(map[string]*CacheEntry),
		ttl:  ttl,
	}

	// Start background cleanup goroutine
	go c.cleanupExpired()

	return c
}

// Set stores a value in the cache with TTL
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &CacheEntry{
		Data:      value,
		ExpiresAt: time.Now().Add(c.ttl),
	}

	c.stats.CurrentSize = len(c.data)
}

// Get retrieves a value from the cache
// Returns (value, true) if found and not expired, (nil, false) otherwise
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.stats.TotalRequests++

	entry, exists := c.data[key]
	if !exists {
		c.stats.Misses++
		return nil, false
	}

	if entry.IsExpired() {
		c.stats.Misses++
		c.stats.Expirations++
		// Don't delete here to avoid deadlock, will be cleaned up by background goroutine
		return nil, false
	}

	c.stats.Hits++
	return entry.Data, true
}

// cleanupExpired removes expired entries periodically
func (c *Cache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.removeExpiredEntries()
	}
}

// removeExpiredEntries scans and removes expired entries
func (c *Cache) removeExpiredEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	expiredCount := 0

	for key, entry := range c.data {
		if now.After(entry.ExpiresAt) {
			delete(c.data, key)
			expiredCount++
		}
	}

	if expiredCount > 0 {
		c.stats.Expirations += int64(expiredCount)
		c.stats.CurrentSize = len(c.data)
	}
}

// GenerateKey creates a cache key from an object by hashing its JSON representation
func GenerateKey(prefix string, obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		// Fallback to string representation if marshaling fails
		return fmt.Sprintf("%s:%v", prefix, obj)
	}

	hash := sha256.Sum256(data)
	return fmt.Sprintf("%s:%x", prefix, hash[:8]) // Use first 8 bytes of hash
}

// GenerateKey generates a cache key from a search request
func (c *Cache) GenerateKey(req models.SearchRequest) string {
	return GenerateKey("search", req)
}
