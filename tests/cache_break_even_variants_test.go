/*
Package tests contains benchmarks to test break-even point with different numbers of field accesses.

These benchmarks test scenarios with 1, 2, 3, 4, and 5 policies accessing the same field
to determine when cache becomes beneficial.
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
func createPoliciesWithRepeatedField(field string, count int) []base.Policy {
	policies := make([]base.Policy, count)
	for i := 0; i < count; i++ {
		policies[i] = base.Policy{
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
	engine, _ := guardian.NewGuardianFromPolices(casher, policies)
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
	engine, _ := guardian.NewGuardianFromPolices(nil, policies)
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
	engine, _ := guardian.NewGuardianFromPolices(casher, policies)
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
	engine, _ := guardian.NewGuardianFromPolices(nil, policies)
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
	engine, _ := guardian.NewGuardianFromPolices(casher, policies)
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
	engine, _ := guardian.NewGuardianFromPolices(nil, policies)
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
	engine, _ := guardian.NewGuardianFromPolices(casher, policies)
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
	engine, _ := guardian.NewGuardianFromPolices(nil, policies)
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
	engine, _ := guardian.NewGuardianFromPolices(casher, policies)
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
	engine, _ := guardian.NewGuardianFromPolices(nil, policies)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, source, target, "user:read")
	}
}
