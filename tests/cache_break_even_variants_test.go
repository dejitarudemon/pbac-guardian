/*
Package tests contains benchmarks to test break-even point with different numbers of field accesses.

These benchmarks test scenarios with 1, 2, 3, 4, and 5 policies accessing the same field
to determine when cache becomes beneficial. The benchmarks compare performance with cache
enabled (DefaultCasher) versus without cache (nil casher) to measure the break-even point.

Results show that cache becomes beneficial when the same field is accessed 3+ times
within a single action evaluation.
*/
package tests

import (
	"context"
	"testing"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
	"github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

// Helper function to create policies that access the same field N times
func createPoliciesWithRepeatedField(field string, count int) []base.RawPolicy {
	policies := make([]base.RawPolicy, count)
	for i := 0; i < count; i++ {
		policies[i] = base.RawPolicy{
			Name:   "policy-" + string(rune('1'+i)),
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				field: {Eq: "admin"},
			},
		}
	}
	return policies
}

func BenchmarkEvaluate1PolicyWithCache(b *testing.B) {
	policies := createPoliciesWithRepeatedField("source:role", 1)
	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, _ := guardian.NewGuardianFromPolices(casher, policies, config)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}

func BenchmarkEvaluate1PolicyWithoutCache(b *testing.B) {
	policies := createPoliciesWithRepeatedField("source:role", 1)
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, _ := guardian.NewGuardianFromPolices(nil, policies, config)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}

func BenchmarkEvaluate2PoliciesWithCache(b *testing.B) {
	policies := createPoliciesWithRepeatedField("source:role", 2)
	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, _ := guardian.NewGuardianFromPolices(casher, policies, config)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}

func BenchmarkEvaluate2PoliciesWithoutCache(b *testing.B) {
	policies := createPoliciesWithRepeatedField("source:role", 2)
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, _ := guardian.NewGuardianFromPolices(nil, policies, config)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}

func BenchmarkEvaluate3PoliciesWithCache(b *testing.B) {
	policies := createPoliciesWithRepeatedField("source:role", 3)
	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, _ := guardian.NewGuardianFromPolices(casher, policies, config)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}

func BenchmarkEvaluate3PoliciesWithoutCache(b *testing.B) {
	policies := createPoliciesWithRepeatedField("source:role", 3)
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, _ := guardian.NewGuardianFromPolices(nil, policies, config)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}

func BenchmarkEvaluate4PoliciesWithCache(b *testing.B) {
	policies := createPoliciesWithRepeatedField("source:role", 4)
	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, _ := guardian.NewGuardianFromPolices(casher, policies, config)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}

func BenchmarkEvaluate4PoliciesWithoutCache(b *testing.B) {
	policies := createPoliciesWithRepeatedField("source:role", 4)
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, _ := guardian.NewGuardianFromPolices(nil, policies, config)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}

func BenchmarkEvaluate5PoliciesWithCache(b *testing.B) {
	policies := createPoliciesWithRepeatedField("source:role", 5)
	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, _ := guardian.NewGuardianFromPolices(casher, policies, config)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}

func BenchmarkEvaluate5PoliciesWithoutCache(b *testing.B) {
	policies := createPoliciesWithRepeatedField("source:role", 5)
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, _ := guardian.NewGuardianFromPolices(nil, policies, config)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}
