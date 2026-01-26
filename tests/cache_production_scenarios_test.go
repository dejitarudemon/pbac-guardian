/*
Package tests contains benchmarks for production-like scenarios.

These benchmarks simulate real-world usage with:
- 10+ different actions
- 20-30 policies per action
- Fields accessed 3-10 times within policies
*/
package tests

import (
	"context"
	"testing"

	noctisguard "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
	"github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

// Helper to create policies with repeated field access
// Since map keys must be unique, we create multiple policies that all access the same field
// This simulates the scenario where the same field is accessed multiple times across policies
func createPoliciesWithFieldAccess(action string, field string, count int, accessCount int) []base.Policy {
	// Create 'count' policies, each accessing the field once
	// The total number of accesses to the field will be 'count' (one per policy)
	// To simulate 'accessCount' accesses, we need to create more policies
	totalPolicies := count * accessCount
	policies := make([]base.Policy, totalPolicies)

	for i := 0; i < totalPolicies; i++ {
		conditions := make(map[string]base.Condition)
		// Each policy accesses the field once, but we create multiple policies
		// to simulate multiple accesses to the same field
		switch i % 4 {
		case 0:
			conditions[field] = base.Condition{Eq: "admin"}
		case 1:
			conditions[field] = base.Condition{Neq: "guest"}
		case 2:
			conditions[field] = base.Condition{Eq: "moderator"}
		case 3:
			conditions[field] = base.Condition{Neq: "banned"}
		}
		policies[i] = base.Policy{
			Name:       "policy-" + string(rune('a'+(i%26))) + string(rune('0'+(i/26))),
			Action:     action,
			Effect:     base.Effect_ALLOW,
			Conditions: conditions,
		}
	}
	return policies
}

// Helper to create multiple actions with policies
// fieldAccessCount = how many policies per action (each policy accesses the field once)
// This simulates the field being accessed 'fieldAccessCount' times within one action evaluation
func createMultipleActions(actions []string, policiesPerAction int, fieldAccessCount int) []base.Policy {
	allPolicies := make([]base.Policy, 0)
	for _, action := range actions {
		// Create 'policiesPerAction' groups, each with 'fieldAccessCount' policies
		// Total: policiesPerAction × fieldAccessCount policies per action
		for group := 0; group < policiesPerAction; group++ {
			for access := 0; access < fieldAccessCount; access++ {
				conditions := make(map[string]base.Condition)
				switch access % 4 {
				case 0:
					conditions["source:role"] = base.Condition{Eq: "admin"}
				case 1:
					conditions["source:role"] = base.Condition{Neq: "guest"}
				case 2:
					conditions["source:role"] = base.Condition{Eq: "moderator"}
				case 3:
					conditions["source:role"] = base.Condition{Neq: "banned"}
				}
				allPolicies = append(allPolicies, base.Policy{
					Name:       "policy-" + action + "-" + string(rune('a'+group)) + string(rune('0'+access)),
					Action:     action,
					Effect:     base.Effect_ALLOW,
					Conditions: conditions,
				})
			}
		}
	}
	return allPolicies
}

/*
BenchmarkProductionScenario_20Policies_3Accesses measures performance
with 20 policies accessing the same field 3 times (at break-even point).
*/
func BenchmarkProductionScenario_20Policies_3Accesses(b *testing.B) {
	actions := []string{"user:read", "user:write", "user:delete", "user:update", "user:create",
		"doc:read", "doc:write", "doc:delete", "doc:update", "doc:create"}
	policies := createMultipleActions(actions, 2, 3) // 10 actions × 2 policies = 20 policies

	casher := implemented.NewDefaultCasher()
	engine, _ := noctisguard.NewNoctisFromPolices(casher, policies)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		action := actions[i%len(actions)]
		_, _ = engine.Evaluate(ctx, source, target, action)
	}
}

func BenchmarkProductionScenario_20Policies_3Accesses_NoCache(b *testing.B) {
	actions := []string{"user:read", "user:write", "user:delete", "user:update", "user:create",
		"doc:read", "doc:write", "doc:delete", "doc:update", "doc:create"}
	policies := createMultipleActions(actions, 2, 3)

	engine, _ := noctisguard.NewNoctisFromPolices(nil, policies)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		action := actions[i%len(actions)]
		_, _ = engine.Evaluate(ctx, source, target, action)
	}
}

/*
BenchmarkProductionScenario_30Policies_5Accesses measures performance
with 30 policies accessing the same field 5 times (well above break-even).
*/
func BenchmarkProductionScenario_30Policies_5Accesses(b *testing.B) {
	actions := []string{"user:read", "user:write", "user:delete", "user:update", "user:create",
		"doc:read", "doc:write", "doc:delete", "doc:update", "doc:create"}
	policies := createMultipleActions(actions, 3, 5) // 10 actions × 3 policies = 30 policies

	casher := implemented.NewDefaultCasher()
	engine, _ := noctisguard.NewNoctisFromPolices(casher, policies)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		action := actions[i%len(actions)]
		_, _ = engine.Evaluate(ctx, source, target, action)
	}
}

func BenchmarkProductionScenario_30Policies_5Accesses_NoCache(b *testing.B) {
	actions := []string{"user:read", "user:write", "user:delete", "user:update", "user:create",
		"doc:read", "doc:write", "doc:delete", "doc:update", "doc:create"}
	policies := createMultipleActions(actions, 3, 5)

	engine, _ := noctisguard.NewNoctisFromPolices(nil, policies)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		action := actions[i%len(actions)]
		_, _ = engine.Evaluate(ctx, source, target, action)
	}
}

/*
BenchmarkProductionScenario_30Policies_10Accesses measures performance
with 30 policies accessing the same field 10 times (maximum benefit scenario).
*/
func BenchmarkProductionScenario_30Policies_10Accesses(b *testing.B) {
	actions := []string{"user:read", "user:write", "user:delete", "user:update", "user:create",
		"doc:read", "doc:write", "doc:delete", "doc:update", "doc:create"}
	policies := createMultipleActions(actions, 3, 10) // 10 actions × 3 policies = 30 policies, 10 accesses each

	casher := implemented.NewDefaultCasher()
	engine, _ := noctisguard.NewNoctisFromPolices(casher, policies)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		action := actions[i%len(actions)]
		_, _ = engine.Evaluate(ctx, source, target, action)
	}
}

func BenchmarkProductionScenario_30Policies_10Accesses_NoCache(b *testing.B) {
	actions := []string{"user:read", "user:write", "user:delete", "user:update", "user:create",
		"doc:read", "doc:write", "doc:delete", "doc:update", "doc:create"}
	policies := createMultipleActions(actions, 3, 10)

	engine, _ := noctisguard.NewNoctisFromPolices(nil, policies)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		action := actions[i%len(actions)]
		_, _ = engine.Evaluate(ctx, source, target, action)
	}
}

/*
BenchmarkProductionScenario_MixedFields measures performance with
multiple fields accessed different number of times (realistic scenario).
*/
func BenchmarkProductionScenario_MixedFields(b *testing.B) {
	// Create policies accessing different fields with varying access counts
	// Each policy accesses a field once, multiple policies = multiple accesses
	allPolicies := make([]base.Policy, 0, 30)
	actions := []string{"action0:read", "action1:read", "action2:read", "action3:read", "action4:read",
		"action5:read", "action6:read", "action7:read", "action8:read", "action9:read"}

	policyIdx := 0
	for _, action := range actions {
		// source:role accessed 5 times (5 policies)
		for i := 0; i < 5; i++ {
			conditions := map[string]base.Condition{
				"source:role": {Eq: "admin"},
			}
			allPolicies = append(allPolicies, base.Policy{
				Name:       "p-role-" + string(rune('0'+policyIdx)),
				Action:     action,
				Effect:     base.Effect_ALLOW,
				Conditions: conditions,
			})
			policyIdx++
		}

		// source:name accessed 3 times (3 policies)
		for i := 0; i < 3; i++ {
			conditions := map[string]base.Condition{
				"source:name": {Eq: "target:owner"},
			}
			allPolicies = append(allPolicies, base.Policy{
				Name:       "p-name-" + string(rune('0'+policyIdx)),
				Action:     action,
				Effect:     base.Effect_ALLOW,
				Conditions: conditions,
			})
			policyIdx++
		}

		// target:type accessed 7 times (7 policies)
		for i := 0; i < 7; i++ {
			conditions := map[string]base.Condition{
				"target:type": {Eq: "public"},
			}
			allPolicies = append(allPolicies, base.Policy{
				Name:       "p-type-" + string(rune('0'+policyIdx)),
				Action:     action,
				Effect:     base.Effect_ALLOW,
				Conditions: conditions,
			})
			policyIdx++
		}

		if len(allPolicies) >= 30 {
			break
		}
	}

	casher := implemented.NewDefaultCasher()
	engine, _ := noctisguard.NewNoctisFromPolices(casher, allPolicies)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		action := actions[i%len(actions)]
		_, _ = engine.Evaluate(ctx, source, target, action)
	}
}

func BenchmarkProductionScenario_MixedFields_NoCache(b *testing.B) {
	allPolicies := make([]base.Policy, 0, 30)
	actions := []string{"action0:read", "action1:read", "action2:read", "action3:read", "action4:read",
		"action5:read", "action6:read", "action7:read", "action8:read", "action9:read"}

	policyIdx := 0
	for _, action := range actions {
		// source:role accessed 5 times
		for i := 0; i < 5; i++ {
			conditions := map[string]base.Condition{
				"source:role": {Eq: "admin"},
			}
			allPolicies = append(allPolicies, base.Policy{
				Name:       "p-role-" + string(rune('0'+policyIdx)),
				Action:     action,
				Effect:     base.Effect_ALLOW,
				Conditions: conditions,
			})
			policyIdx++
		}

		// source:name accessed 3 times
		for i := 0; i < 3; i++ {
			conditions := map[string]base.Condition{
				"source:name": {Eq: "target:owner"},
			}
			allPolicies = append(allPolicies, base.Policy{
				Name:       "p-name-" + string(rune('0'+policyIdx)),
				Action:     action,
				Effect:     base.Effect_ALLOW,
				Conditions: conditions,
			})
			policyIdx++
		}

		// target:type accessed 7 times
		for i := 0; i < 7; i++ {
			conditions := map[string]base.Condition{
				"target:type": {Eq: "public"},
			}
			allPolicies = append(allPolicies, base.Policy{
				Name:       "p-type-" + string(rune('0'+policyIdx)),
				Action:     action,
				Effect:     base.Effect_ALLOW,
				Conditions: conditions,
			})
			policyIdx++
		}

		if len(allPolicies) >= 30 {
			break
		}
	}

	engine, _ := noctisguard.NewNoctisFromPolices(nil, allPolicies)
	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		action := actions[i%len(actions)]
		_, _ = engine.Evaluate(ctx, source, target, action)
	}
}
