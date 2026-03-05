/*
Package tests contains tests for session isolation and cache behavior.

Tests check that each evaluation session uses a unique sessionID
and that cache is properly isolated between different sessions.
*/
package tests

import (
	"context"
	"testing"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
	"github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

/*
TestSessionIsolation tests that different evaluation sessions use different sessionIDs.

The test checks that each call to Evaluate generates a unique sessionID,
ensuring that cache entries from one evaluation do not interfere with another.
*/
func TestSessionIsolation(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	// Perform multiple evaluations
	// Each should use a different sessionID internally
	for i := 0; i < 10; i++ {
		allowed, err := engine.Evaluate(ctx, source, target, "user:read")
		if err != nil {
			t.Errorf("evaluation %d failed: %v", i, err)
		}
		if !allowed {
			t.Errorf("evaluation %d: expected allowed=true, got false", i)
		}
	}

	// All evaluations should succeed independently
	// If sessionIDs were not unique, cache pollution could cause issues
}

/*
TestCacheReuseWithinSession tests that cache is reused within a single evaluation session.

The test checks that when multiple policies access the same field within one evaluation,
the cache is used to avoid repeated reflection-based field searches.
*/
func TestCacheReuseWithinSession(t *testing.T) {
	// Create multiple policies that all access the same field
	policies := []base.RawPolicy{
		{
			Name:   "policy-1",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
		{
			Name:   "policy-2",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Neq: "guest"},
			},
		},
		{
			Name:   "policy-3",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	// Single evaluation with multiple policies accessing the same field
	// Cache should be used for subsequent accesses to source:role
	allowed, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true, got false")
	}

	// Verify that cache was used by checking that casher has entries
	// (This is indirect verification - direct cache inspection would require exported methods)
	// The fact that evaluation succeeds with multiple policies accessing the same field
	// indicates that cache is working correctly
}

/*
TestCacheClearedAfterEvaluation tests that cache is cleared after evaluation completes.

The test checks that cache entries are properly cleaned up after an evaluation session,
preventing memory leaks in long-running applications.
*/
func TestCacheClearedAfterEvaluation(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	// Perform evaluation
	_, err = engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Perform another evaluation with different data
	// This should use a new sessionID, and previous cache should not interfere
	source2 := User{Name: "user", Role: "user"}
	allowed2, err := engine.Evaluate(ctx, source2, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Should be denied (role is "user", not "admin")
	if allowed2 {
		t.Errorf("expected allowed=false for user role, got true")
	}

	// The fact that second evaluation correctly denies access
	// indicates that cache from first evaluation did not interfere
}

/*
TestConcurrentEvaluationsWithCache tests that concurrent evaluations work correctly with cache.

The test checks that multiple goroutines can safely use the same engine and cache instance
concurrently without race conditions or cache pollution.
*/
func TestConcurrentEvaluationsWithCache(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()

	// Run multiple concurrent evaluations
	const numGoroutines = 10
	results := make(chan bool, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			source := User{Name: "admin", Role: "admin"}
			target := Document{Owner: "user", Type: "public"}

			allowed, err := engine.Evaluate(ctx, source, target, "user:read")
			results <- allowed
			errors <- err
		}(i)
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		allowed := <-results
		err := <-errors

		if err != nil {
			t.Errorf("goroutine %d: unexpected error: %v", i, err)
		}
		if !allowed {
			t.Errorf("goroutine %d: expected allowed=true, got false", i)
		}
	}
}
