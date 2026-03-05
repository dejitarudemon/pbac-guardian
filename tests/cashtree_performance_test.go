/*
Package tests contains performance benchmarks for CashTree optimization.

This file contains benchmarks to analyze the performance impact of CashTree
on memory usage and execution time when caching is selectively disabled
for rarely accessed fields.
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
BenchmarkCashTree_WithOptimization benchmarks performance with CashTree optimization.

This benchmark simulates a scenario with mixed field access patterns:
- Some fields accessed frequently (5+ times) - should be cached
- Some fields accessed rarely (1-2 times) - should not be cached

The CashTree optimization should reduce memory usage by avoiding cache storage
for rarely accessed fields while maintaining performance for frequently accessed ones.
*/
func BenchmarkCashTree_WithOptimization(b *testing.B) {
	// Create policies with mixed access patterns
	policies := []base.RawPolicy{
		// Frequently accessed field: source:role (5 times)
		{Name: "p1", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"source:role": {Eq: "admin"}}},
		{Name: "p2", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"source:role": {Eq: "admin"}}},
		{Name: "p3", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"source:role": {Eq: "admin"}}},
		{Name: "p4", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"source:role": {Eq: "admin"}}},
		{Name: "p5", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"source:role": {Eq: "admin"}}},

		// Rarely accessed fields: source:name (1 time), target:type (2 times)
		{Name: "p6", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"source:name": {Eq: "target:owner"}}},
		{Name: "p7", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"target:type": {Eq: "public"}}},
		{Name: "p8", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"target:type": {Eq: "public"}}},
	}

	casher := implemented.NewDefaultCasher()
	config := base.Config{
		ConditionsMap:        nil,
		CashDisableThreShold: 3, // Disable cache for fields accessed < 3 times
	}

	engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "admin", Type: "public"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}

/*
BenchmarkCashTree_WithoutOptimization benchmarks performance without CashTree optimization.

This benchmark uses the same scenario but with threshold=1, meaning all fields are cached.
This allows comparison of memory usage and performance impact.
*/
func BenchmarkCashTree_WithoutOptimization(b *testing.B) {
	// Create policies with mixed access patterns
	policies := []base.RawPolicy{
		// Frequently accessed field: source:role (5 times)
		{Name: "p1", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"source:role": {Eq: "admin"}}},
		{Name: "p2", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"source:role": {Eq: "admin"}}},
		{Name: "p3", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"source:role": {Eq: "admin"}}},
		{Name: "p4", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"source:role": {Eq: "admin"}}},
		{Name: "p5", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"source:role": {Eq: "admin"}}},

		// Rarely accessed fields: source:name (1 time), target:type (2 times)
		{Name: "p6", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"source:name": {Eq: "target:owner"}}},
		{Name: "p7", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"target:type": {Eq: "public"}}},
		{Name: "p8", Action: "user:read", Effect: base.Effect_ALLOW, Conditions: map[string]base.Condition{"target:type": {Eq: "public"}}},
	}

	casher := implemented.NewDefaultCasher()
	config := base.Config{
		ConditionsMap:        nil,
		CashDisableThreShold: 1, // Cache all fields (no optimization)
	}

	engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "admin", Type: "public"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}

/*
BenchmarkCashTree_ProductionScenario benchmarks a production-like scenario with 30 policies.

This benchmark tests CashTree optimization in a realistic scenario with:
- Multiple actions
- Mixed field access patterns
- Some fields accessed frequently, others rarely
*/
func BenchmarkCashTree_ProductionScenario(b *testing.B) {
	allPolicies := make([]base.RawPolicy, 0, 30)
	actions := []string{"action0:read", "action1:read", "action2:read"}

	policyIdx := 0
	for _, action := range actions {
		// source:role accessed 5 times per action (frequently accessed)
		for i := 0; i < 5; i++ {
			allPolicies = append(allPolicies, base.RawPolicy{
				Name:       "p-role-" + string(rune('0'+policyIdx)),
				Action:     action,
				Effect:     base.Effect_ALLOW,
				Conditions: map[string]base.Condition{"source:role": {Eq: "admin"}},
			})
			policyIdx++
		}

		// source:name accessed 1 time per action (rarely accessed)
		allPolicies = append(allPolicies, base.RawPolicy{
			Name:       "p-name-" + string(rune('0'+policyIdx)),
			Action:     action,
			Effect:     base.Effect_ALLOW,
			Conditions: map[string]base.Condition{"source:name": {Eq: "target:owner"}},
		})
		policyIdx++

		// target:type accessed 2 times per action (rarely accessed)
		for i := 0; i < 2; i++ {
			allPolicies = append(allPolicies, base.RawPolicy{
				Name:       "p-type-" + string(rune('0'+policyIdx)),
				Action:     action,
				Effect:     base.Effect_ALLOW,
				Conditions: map[string]base.Condition{"target:type": {Eq: "public"}},
			})
			policyIdx++
		}

		if len(allPolicies) >= 30 {
			break
		}
	}

	casher := implemented.NewDefaultCasher()
	config := base.Config{
		ConditionsMap:        nil,
		CashDisableThreShold: 3, // Disable cache for fields accessed < 3 times
	}

	engine, err := guardian.NewGuardianFromPolices(casher, allPolicies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "admin", Type: "public"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		action := actions[i%len(actions)]
		_, _ = engine.Evaluate(ctx, source, target, action)
	}
}
