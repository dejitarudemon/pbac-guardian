/*
Package cashing provides interfaces and implementations for L1 caching in the policy evaluation engine.

The package contains the Casher interface for cache implementations and related types.
*/
package cashing

import "context"

/*
Casher is an interface for L1 cache to avoid re-searching struct fields by reflect package.

The cache is used to store field values retrieved from structures during policy evaluation.
Each evaluation session uses a unique sessionID that identifies the cache scope for a single
policy application. This allows reusing field values within the same evaluation without
repeated reflection-based searches.

The cache provides significant performance improvements in production scenarios:
  - 14-41% faster execution when fields are accessed 3+ times
  - 47-63% fewer allocations, reducing GC pressure
  - 25-36% less memory usage with 5+ field accesses

The implementation must be thread-safe to support concurrent policy evaluations.
For a ready-to-use implementation, see implemented.DefaultCasher.

Example usage:

	type MyCasher struct {
		storage map[string]map[string]any
		mu      sync.RWMutex
	}

	func (c *MyCasher) Set(ctx context.Context, sessionID, key string, value any) error {
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.storage[sessionID] == nil {
			c.storage[sessionID] = make(map[string]any)
		}
		c.storage[sessionID][key] = value
		return nil
	}

	func (c *MyCasher) Get(ctx context.Context, sessionID, key string) (any, error) {
		c.mu.RLock()
		defer c.mu.RUnlock()
		if session, ok := c.storage[sessionID]; ok {
			if value, ok := session[key]; ok {
				return value, nil
			}
		}
		return nil, nil
	}

	func (c *MyCasher) Clear(ctx context.Context, sessionID string) error {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.storage, sessionID)
		return nil
	}
*/
type Casher interface {
	/*
		Set stores a value in L1 cache under the specified key within the session scope.

		The function stores the value using sessionID and key as composite identifier.
		This allows multiple evaluation sessions to use the same cache instance without
		interference.

		Parameters:
			- ctx - performing context for operation cancellation and timeout control
			- sessionID - unique identifier for the current evaluation session (cache scope)
			- key - cache key (typically a field path like "source:name" or "target:owner")
			- value - value to store (struct field value, primitive or custom type)

		Returns:
			- error - error if value storage failed, nil on success
	*/
	Set(ctx context.Context, sessionID, key string, value any) error

	/*
		Get retrieves a value from L1 cache by key within the session scope.

		The function looks up the value using sessionID and key as composite identifier.
		If the value is not found in cache, it returns (nil, nil) to indicate cache miss.

		Parameters:
			- ctx - performing context for operation cancellation and timeout control
			- sessionID - unique identifier for the current evaluation session (cache scope)
			- key - cache key (typically a field path like "source:name" or "target:owner")

		Returns:
			- any - stored value if it exists, nil if not found or error occurred
			- error - error if value retrieval failed, nil if value found or not found (cache miss)
	*/
	Get(ctx context.Context, sessionID, key string) (any, error)

	/*
		Clear removes all cached values for the specified session.

		The function clears all cache entries associated with the sessionID, effectively
		resetting the cache for that evaluation session. This is typically called after
		an evaluation session completes.

		Parameters:
			- ctx - performing context for operation cancellation and timeout control
			- sessionID - unique identifier for the evaluation session to clear

		Returns:
			- error - error if cache clearing failed, nil on success
	*/
	Clear(ctx context.Context, sessionID string) error
}
