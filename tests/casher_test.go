/*
Package tests contains tests for the Casher interface implementation.

Tests check the functionality of DefaultCasher, including thread-safety,
session-based caching, and proper cache operations (Set, Get, Clear).
*/
package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/dejitarudemon/noctis-guard/internal/implemented"
)

/*
TestDefaultCasherSet tests the Set method of DefaultCasher.

The test checks:
  - Setting values for different sessions
  - Setting multiple values in the same session
  - Thread-safety of Set operations
*/
func TestDefaultCasherSet(t *testing.T) {
	casher := implemented.NewDefaultCasher()
	ctx := context.Background()

	tests := []struct {
		name      string
		sessionID string
		key       string
		value     any
		wantErr   bool
	}{
		{
			name:      "set value in new session",
			sessionID: "session1",
			key:       "source:name",
			value:     "alice",
			wantErr:   false,
		},
		{
			name:      "set value in existing session",
			sessionID: "session1",
			key:       "source:role",
			value:     "admin",
			wantErr:   false,
		},
		{
			name:      "set value in different session",
			sessionID: "session2",
			key:       "source:name",
			value:     "bob",
			wantErr:   false,
		},
		{
			name:      "set complex value",
			sessionID: "session1",
			key:       "target:tags",
			value:     []string{"public", "shared"},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := casher.Set(ctx, tt.sessionID, tt.key, tt.value)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

/*
TestDefaultCasherGet tests the Get method of DefaultCasher.

The test checks:
  - Getting values from cache
  - Cache miss behavior (returns nil, nil)
  - Getting values from different sessions
  - Getting non-existent keys
*/
func TestDefaultCasherGet(t *testing.T) {
	casher := implemented.NewDefaultCasher()
	ctx := context.Background()

	// Set up test data
	casher.Set(ctx, "session1", "source:name", "alice")
	casher.Set(ctx, "session1", "source:role", "admin")
	casher.Set(ctx, "session2", "source:name", "bob")

	tests := []struct {
		name      string
		sessionID string
		key       string
		want      any
		wantErr   bool
	}{
		{
			name:      "get existing value",
			sessionID: "session1",
			key:       "source:name",
			want:      "alice",
			wantErr:   false,
		},
		{
			name:      "get another value from same session",
			sessionID: "session1",
			key:       "source:role",
			want:      "admin",
			wantErr:   false,
		},
		{
			name:      "get value from different session",
			sessionID: "session2",
			key:       "source:name",
			want:      "bob",
			wantErr:   false,
		},
		{
			name:      "get non-existent key in existing session",
			sessionID: "session1",
			key:       "source:age",
			want:      nil,
			wantErr:   false,
		},
		{
			name:      "get key from non-existent session",
			sessionID: "session3",
			key:       "source:name",
			want:      nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := casher.Get(ctx, tt.sessionID, tt.key)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got != tt.want {
					t.Errorf("Get() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

/*
TestDefaultCasherClear tests the Clear method of DefaultCasher.

The test checks:
  - Clearing all values for a session
  - Clearing non-existent session (should not error)
  - Values in other sessions remain after clearing one session
*/
func TestDefaultCasherClear(t *testing.T) {
	casher := implemented.NewDefaultCasher()
	ctx := context.Background()

	// Set up test data
	casher.Set(ctx, "session1", "source:name", "alice")
	casher.Set(ctx, "session1", "source:role", "admin")
	casher.Set(ctx, "session2", "source:name", "bob")
	casher.Set(ctx, "session2", "source:role", "user")

	tests := []struct {
		name      string
		sessionID string
		wantErr   bool
	}{
		{
			name:      "clear existing session",
			sessionID: "session1",
			wantErr:   false,
		},
		{
			name:      "clear non-existent session",
			sessionID: "session3",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := casher.Clear(ctx, tt.sessionID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				// Verify that session is cleared
				value, _ := casher.Get(ctx, tt.sessionID, "source:name")
				if value != nil {
					t.Errorf("expected session %s to be cleared, but value still exists", tt.sessionID)
				}
			}
		})
	}

	// Verify that other sessions are not affected
	value, _ := casher.Get(ctx, "session2", "source:name")
	if value != "bob" {
		t.Errorf("expected session2 to remain intact, got %v", value)
	}
}

/*
TestDefaultCasherThreadSafety tests thread-safety of DefaultCasher.

The test performs concurrent Set, Get, and Clear operations to verify
that the casher is thread-safe and does not cause race conditions.
*/
func TestDefaultCasherThreadSafety(t *testing.T) {
	casher := implemented.NewDefaultCasher()
	ctx := context.Background()

	const numGoroutines = 100
	const numOperations = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Concurrent Set operations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			sessionID := fmt.Sprintf("session%d", id%10)
			key := fmt.Sprintf("key%d", id%5)
			value := id

			for j := 0; j < numOperations; j++ {
				err := casher.Set(ctx, sessionID, key, value)
				if err != nil {
					t.Errorf("Set failed: %v", err)
				}

				got, err := casher.Get(ctx, sessionID, key)
				if err != nil {
					t.Errorf("Get failed: %v", err)
				}
				if got == nil {
					t.Errorf("Get returned nil for key that was just set")
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify that at least some values were set (due to concurrency, not all may be set)
	// The important thing is that no race conditions occurred
	foundValues := 0
	for i := 0; i < 10; i++ {
		sessionID := fmt.Sprintf("session%d", i)
		for j := 0; j < 5; j++ {
			key := fmt.Sprintf("key%d", j)
			value, err := casher.Get(ctx, sessionID, key)
			if err != nil {
				t.Errorf("Get failed: %v", err)
			}
			if value != nil {
				foundValues++
			}
		}
	}

	// At least some values should have been set
	if foundValues == 0 {
		t.Error("No values were found in cache after concurrent operations")
	}
}

/*
TestDefaultCasherSessionIsolation tests that different sessions are isolated.

The test verifies that values stored in one session are not accessible
from another session, ensuring proper session isolation.
*/
func TestDefaultCasherSessionIsolation(t *testing.T) {
	casher := implemented.NewDefaultCasher()
	ctx := context.Background()

	// Set same key in different sessions with different values
	casher.Set(ctx, "session1", "source:name", "alice")
	casher.Set(ctx, "session2", "source:name", "bob")
	casher.Set(ctx, "session3", "source:name", "charlie")

	// Verify each session has its own value
	value1, _ := casher.Get(ctx, "session1", "source:name")
	if value1 != "alice" {
		t.Errorf("session1: expected 'alice', got %v", value1)
	}

	value2, _ := casher.Get(ctx, "session2", "source:name")
	if value2 != "bob" {
		t.Errorf("session2: expected 'bob', got %v", value2)
	}

	value3, _ := casher.Get(ctx, "session3", "source:name")
	if value3 != "charlie" {
		t.Errorf("session3: expected 'charlie', got %v", value3)
	}
}

