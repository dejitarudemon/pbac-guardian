/*
Package tests contains benchmarks to determine cache break-even point.

These benchmarks measure when cache becomes beneficial compared to direct
reflection access. The goal is to find how many times a field must be
accessed within a single action evaluation for cache to be profitable.
*/
package tests

import (
	"context"
	"reflect"
	"testing"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
	"github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

/*
BenchmarkReflectSingleAccess measures the cost of a single field access via reflection.

This represents the baseline cost without caching.
*/
func BenchmarkReflectSingleAccess(b *testing.B) {
	user := User{Name: "alice", Role: "admin", Age: 25}
	userType := reflect.TypeOf(user)
	userValue := reflect.ValueOf(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate field access via reflection (similar to what the library does)
		for j := 0; j < userType.NumField(); j++ {
			field := userType.Field(j)
			tag := field.Tag.Get("pbac-guardian")
			if tag == "name" {
				_ = userValue.Field(j).Interface()
				break
			}
		}
	}
}

/*
BenchmarkCacheSingleAccess measures the cost of a single field access via cache.

This includes cache lookup overhead.
*/
func BenchmarkCacheSingleAccess(b *testing.B) {
	casher := implemented.NewDefaultCasher()
	ctx := context.Background()
	sessionID := "test-session"
	key := "source:name"
	value := "alice"

	// Pre-populate cache
	casher.Set(ctx, sessionID, key, value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = casher.Get(ctx, sessionID, key)
	}
}

/*
BenchmarkReflectMultipleAccess measures the cost of multiple field accesses via reflection.

This simulates accessing the same field N times (as would happen with multiple policies).
*/
func BenchmarkReflectMultipleAccess(b *testing.B) {
	user := User{Name: "alice", Role: "admin", Age: 25}
	userType := reflect.TypeOf(user)
	userValue := reflect.ValueOf(user)
	accessCount := 5 // Number of times to access the same field

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for k := 0; k < accessCount; k++ {
			for j := 0; j < userType.NumField(); j++ {
				field := userType.Field(j)
				tag := field.Tag.Get("pbac-guardian")
				if tag == "name" {
					_ = userValue.Field(j).Interface()
					break
				}
			}
		}
	}
}

/*
BenchmarkCacheMultipleAccess measures the cost of multiple field accesses via cache.

This simulates accessing the same cached field N times.
*/
func BenchmarkCacheMultipleAccess(b *testing.B) {
	casher := implemented.NewDefaultCasher()
	ctx := context.Background()
	sessionID := "test-session"
	key := "source:name"
	value := "alice"
	accessCount := 5 // Number of times to access the same field

	// Pre-populate cache
	casher.Set(ctx, sessionID, key, value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for k := 0; k < accessCount; k++ {
			_, _ = casher.Get(ctx, sessionID, key)
		}
	}
}

/*
BenchmarkEvaluateWithRepeatedFields measures actual policy evaluation
with policies that access the same field multiple times.
*/
func BenchmarkEvaluateWithRepeatedFields(b *testing.B) {
	// Create policies that all access the same field (source:role)
	policies := []base.Policy{
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
		{
			Name:   "policy-4",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Neq: "user"},
			},
		},
		{
			Name:   "policy-5",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	casher := implemented.NewDefaultCasher()
	engine, err := guardian.NewGuardianFromPolices(casher, policies)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}

/*
BenchmarkEvaluateWithoutCacheRepeatedFields measures policy evaluation
without cache when the same field is accessed multiple times.
*/
func BenchmarkEvaluateWithoutCacheRepeatedFields(b *testing.B) {
	// Same policies as above
	policies := []base.Policy{
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
		{
			Name:   "policy-4",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Neq: "user"},
			},
		},
		{
			Name:   "policy-5",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	// No cache
	engine, err := guardian.NewGuardianFromPolices(nil, policies)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}
