/*
Package tests contains benchmark tests for measuring library performance.

Benchmarks check performance of main operations:
  - Engine creation from policies (with and without cache)
  - Policy evaluation for various scenarios
  - Working with various condition types (Eq, Neq, Lt, In)
  - Handling nested structures
  - Cache performance comparison (with cache vs without cache)

The benchmarks help identify performance bottlenecks and measure
the impact of caching on overall evaluation speed.
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
BenchmarkNewGuardianFromPolices measures performance of engine creation from policies.

The benchmark checks time of engine creation with various number of policies.
*/
func BenchmarkNewGuardianFromPolices(b *testing.B) {
	// Create set of policies for testing
	policies := []base.RawPolicy{
		{
			Name:   "admin-read",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
		{
			Name:   "owner-read",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:name": {Eq: "target:owner"},
			},
		},
		{
			Name:   "deny-private",
			Action: "user:read:document",
			Effect: base.Effect_DENY,
			Conditions: map[string]base.Condition{
				"target:type": {Eq: "private"},
				"source:role": {Neq: "admin"},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
		_, err := guardian.NewGuardianFromPolices(nil, policies, config)
		if err != nil {
			b.Fatalf("failed to create engine: %v", err)
		}
	}
}

/*
BenchmarkEvaluateSimple measures performance of simple policy evaluation.

The benchmark checks time of policy evaluation with simple conditions (Eq).
*/
func BenchmarkEvaluateSimple(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "admin-read",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateMultipleConditions measures performance of policy evaluation
with multiple conditions.

The benchmark checks time of policy evaluation with combination of various conditions.
*/
func BenchmarkEvaluateMultipleConditions(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "complex-policy",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {
					Eq:       "admin",
					In: []any{"admin", "moderator"},
				},
				"source:age": {
					Lt: 100,
				},
				"target:tags": {
					In: []any{"public", "shared"},
				},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin", Age: 25}
	target := Document{Owner: "user", Type: "public", Tags: []string{"public"}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateIn measures performance of In condition.

The benchmark checks time of searching value in list through In condition.
*/
func BenchmarkEvaluateIn(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "contains-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {
					In: []any{"admin", "moderator", "user", "guest", "visitor"},
				},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateLt measures performance of Lt (less than) condition.

The benchmark checks time of value comparison through Lt condition.
*/
func BenchmarkEvaluateLt(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "lt-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:age": {
					Lt: 65,
				},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user", Age: 30}
	target := Document{Owner: "user", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateNestedStructures measures performance of working
with nested structures.

The benchmark checks time of getting values from nested structures.
*/
func BenchmarkEvaluateNestedStructures(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "nested-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:user:role": {
					Eq: "admin",
				},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := NestedUser{User: User{Name: "admin", Role: "admin"}}
	target := NestedDocument{Doc: Document{Owner: "user", Type: "public", Tags: []string{}}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateMultiplePolicies measures performance of evaluating
multiple policies for one action.

The benchmark checks time of evaluation when multiple policies are defined for action.
*/
func BenchmarkEvaluateMultiplePolicies(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "policy-1",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
		{
			Name:   "policy-2",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:name": {Eq: "target:owner"},
			},
		},
		{
			Name:   "policy-3",
			Action: "user:read:document",
			Effect: base.Effect_DENY,
			Conditions: map[string]base.Condition{
				"target:type": {Eq: "private"},
				"source:role": {Neq: "admin"},
			},
		},
		{
			Name:   "policy-4",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:age": {Lt: 18},
			},
		},
		{
			Name:   "policy-5",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {In: []any{"moderator", "editor"}},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin", Age: 25}
	target := Document{Owner: "admin", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateFieldComparison measures performance of comparing
fields from different structures.

The benchmark checks time of comparing fields from source and target structures.
*/
func BenchmarkEvaluateFieldComparison(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "field-comparison",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:name": {
					Eq: "target:owner",
				},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "alice", Role: "user"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateDenyPolicy measures performance of policies with DENY effect.

The benchmark checks time of evaluating policies that deny access.
*/
func BenchmarkEvaluateDenyPolicy(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "deny-policy",
			Action: "user:read:document",
			Effect: base.Effect_DENY,
			Conditions: map[string]base.Condition{
				"target:type": {Eq: "private"},
				"source:role": {Neq: "admin"},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user"}
	target := Document{Owner: "other", Type: "private"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateLargeSlice measures performance of In condition
with large list of values.

The benchmark checks time of searching in large list through In condition.
*/
func BenchmarkEvaluateLargeSlice(b *testing.B) {
	// Create large list of roles
	largeRoleList := make([]any, 1000)
	for i := range largeRoleList {
		largeRoleList[i] = "role" + string(rune(i%26+'a'))
	}
	largeRoleList[500] = "admin" // Search value in middle of list

	policies := []base.RawPolicy{
		{
			Name:   "large-slice-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {
					In: largeRoleList,
				},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateNoMatch measures performance of case
when policies do not match action.

The benchmark checks time of evaluation when there are no suitable policies for action.
*/
func BenchmarkEvaluateNoMatch(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "other-action",
			Action: "user:write:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user"}
	target := Document{Owner: "user", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateWithCache measures performance of policy evaluation with L1 cache enabled.

The benchmark compares evaluation performance with cache vs without cache to measure
the performance improvement from caching field values.
*/
func BenchmarkEvaluateWithCache(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "admin-read",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
		{
			Name:   "owner-read",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:name": {Eq: "target:owner"},
			},
		},
		{
			Name:   "deny-private",
			Action: "user:read:document",
			Effect: base.Effect_DENY,
			Conditions: map[string]base.Condition{
				"target:type": {Eq: "private"},
				"source:role": {Neq: "admin"},
			},
		},
	}

	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateWithoutCache measures performance of policy evaluation without L1 cache.

This benchmark is used for comparison with BenchmarkEvaluateWithCache to measure
the performance improvement from caching.
*/
func BenchmarkEvaluateWithoutCache(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "admin-read",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
		{
			Name:   "owner-read",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:name": {Eq: "target:owner"},
			},
		},
		{
			Name:   "deny-private",
			Action: "user:read:document",
			Effect: base.Effect_DENY,
			Conditions: map[string]base.Condition{
				"target:type": {Eq: "private"},
				"source:role": {Neq: "admin"},
			},
		},
	}

	// Use nil casher to disable caching
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateMultiplePoliciesWithCache measures performance of evaluating
multiple policies with cache enabled.
*/
func BenchmarkEvaluateMultiplePoliciesWithCache(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "policy-1",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
		{
			Name:   "policy-2",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:name": {Eq: "target:owner"},
			},
		},
		{
			Name:   "policy-3",
			Action: "user:read:document",
			Effect: base.Effect_DENY,
			Conditions: map[string]base.Condition{
				"target:type": {Eq: "private"},
				"source:role": {Neq: "admin"},
			},
		},
		{
			Name:   "policy-4",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:age": {Lt: 18},
			},
		},
		{
			Name:   "policy-5",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {In: []any{"moderator", "editor"}},
			},
		},
	}

	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin", Age: 25}
	target := Document{Owner: "admin", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateNestedStructuresWithCache measures performance of working
with nested structures with cache enabled.
*/
func BenchmarkEvaluateNestedStructuresWithCache(b *testing.B) {
	policies := []base.RawPolicy{
		{
			Name:   "nested-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:user:role": {
					Eq: "admin",
				},
			},
		},
	}

	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := NestedUser{User: User{Name: "admin", Role: "admin"}}
	target := NestedDocument{Doc: Document{Owner: "user", Type: "public", Tags: []string{}}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}
