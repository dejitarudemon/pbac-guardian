package implemented

import (
	"context"
	"sync"
)

/*
DefaultCasher is an optimized L1 cache implementation using RWMutex with regular map.

This implementation uses a two-level map structure (sessionID -> key -> value) with
RWMutex for thread-safety. RWMutex is optimized for read-heavy workloads, which
is typical for cache operations. This approach:
- Provides efficient concurrent reads with RLock
- Uses regular maps which are faster than sync.Map for known access patterns
- Reduces allocations compared to sync.Map
- Simplifies the code without type assertions

The implementation is optimized for the typical cache access pattern where reads
greatly outnumber writes.
*/
type DefaultCasher struct {
	mu      sync.RWMutex
	storage map[string]map[string]any // sessionID -> key -> value
}

/*
NewDefaultCasher creates a new optimized DefaultCasher instance.

The casher uses RWMutex with regular maps for optimal performance in
read-heavy concurrent scenarios.
*/
func NewDefaultCasher() *DefaultCasher {
	return &DefaultCasher{
		storage: make(map[string]map[string]any),
	}
}

/*
Set stores a value in the cache.

The function uses write lock (Lock) to ensure thread-safety during writes.
It creates the session map if it doesn't exist.
*/
func (c *DefaultCasher) Set(ctx context.Context, sessionID, key string, value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.storage[sessionID] == nil {
		c.storage[sessionID] = make(map[string]any)
	}
	c.storage[sessionID][key] = value
	return nil
}

/*
Get retrieves a value from the cache.

The function uses read lock (RLock) for efficient concurrent reads.
Returns (nil, nil) if the value is not found (cache miss).
*/
func (c *DefaultCasher) Get(ctx context.Context, sessionID, key string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if session, ok := c.storage[sessionID]; ok {
		if value, ok := session[key]; ok {
			return value, nil
		}
	}
	return nil, nil
}

/*
Clear removes all cached values for the specified session.

The function uses write lock (Lock) and efficiently deletes the entire
session map in a single operation.
*/
func (c *DefaultCasher) Clear(ctx context.Context, sessionID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.storage, sessionID)
	return nil
}
