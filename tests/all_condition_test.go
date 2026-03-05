/*
Package tests contains tests for All condition functionality.

Tests verify that the library correctly handles All conditions for checking
if all elements in a collection satisfy all nested conditions.

Logic:
  - AND logic for elements: all elements in the collection must satisfy
  - AND logic for conditions: all conditions in All must be satisfied for each element

Example:
  - Collection: [Group{name: "admin", role: "user"}, Group{name: "admin", role: "admin"}]
  - Condition: All { item:name: {Eq: "admin"}, item:role: {Eq: "admin"} }
  - Result: false (first element has name="admin" but role="user", not "admin")
  
  - Collection: [Group{name: "admin", role: "admin"}, Group{name: "admin", role: "admin"}]
  - Condition: All { item:name: {Eq: "admin"}, item:role: {Eq: "admin"} }
  - Result: true (all elements have both name="admin" AND role="admin")
*/
package tests

import (
	"context"
	"testing"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

func TestAllConditionBasic(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "all-group-admin",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:groups": {
					All: map[string]base.Condition{
						"item:name": {Eq: "admin"},
					},
				},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name   string
		source UserWithAnyGroups
		target Document
		want   bool
	}{
		{
			name: "all groups have admin name",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "admin", Role: "admin"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "not all groups have admin name",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "user", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "single group with admin name",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "empty groups",
			source: UserWithAnyGroups{
				Login:  "alice",
				Groups: []AnyGroup{},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true, // Empty collection satisfies All (vacuous truth)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, err := engine.Evaluate(ctx, tt.source, tt.target, "user:read")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if allowed != tt.want {
				t.Errorf("Evaluate() = %v, want %v", allowed, tt.want)
			}
		})
	}
}

func TestAllConditionMultipleFields(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "all-group-admin-name-and-role",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:groups": {
					All: map[string]base.Condition{
						"item:name": {Eq: "admin"},
						"item:role": {Eq: "admin"},
					},
				},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name   string
		source UserWithAnyGroups
		target Document
		want   bool
	}{
		{
			name: "all groups have both admin name and role",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "admin"},
					{Name: "admin", Role: "admin"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "one group missing admin name",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "admin"},
					{Name: "user", Role: "admin"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "one group missing admin role",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "admin"},
					{Name: "admin", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "all groups missing both conditions",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "user"},
					{Name: "guest", Role: "guest"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "single group with both conditions",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "admin"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, err := engine.Evaluate(ctx, tt.source, tt.target, "user:read")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if allowed != tt.want {
				t.Errorf("Evaluate() = %v, want %v", allowed, tt.want)
			}
		})
	}
}

func TestAllConditionWithOtherConditions(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "all-group-with-role-check",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:login": {Eq: "alice"},
				"source:groups": {
					All: map[string]base.Condition{
						"item:name": {Eq: "admin"},
					},
				},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name   string
		source UserWithAnyGroups
		target Document
		want   bool
	}{
		{
			name: "correct login and all groups are admin",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "admin", Role: "admin"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "wrong login but all groups are admin",
			source: UserWithAnyGroups{
				Login: "bob",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "admin", Role: "admin"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "correct login but not all groups are admin",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "user", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, err := engine.Evaluate(ctx, tt.source, tt.target, "user:read")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if allowed != tt.want {
				t.Errorf("Evaluate() = %v, want %v", allowed, tt.want)
			}
		})
	}
}

func TestAllConditionWithComparison(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "all-group-name-eq-target",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:groups": {
					All: map[string]base.Condition{
						"item:name": {Eq: "target:owner"},
					},
				},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name   string
		source UserWithAnyGroups
		target Document
		want   bool
	}{
		{
			name: "all group names match target owner",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "alice", Role: "user"},
					{Name: "alice", Role: "admin"},
				},
			},
			target: Document{Owner: "alice", Type: "public"},
			want:   true,
		},
		{
			name: "one group name does not match target owner",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "alice", Role: "user"},
					{Name: "bob", Role: "user"},
				},
			},
			target: Document{Owner: "alice", Type: "public"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, err := engine.Evaluate(ctx, tt.source, tt.target, "user:read")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if allowed != tt.want {
				t.Errorf("Evaluate() = %v, want %v", allowed, tt.want)
			}
		})
	}
}

func TestAllConditionEmptyCollection(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "all-empty-collection",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:groups": {
					All: map[string]base.Condition{
						"item:name": {Eq: "admin"},
					},
				},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := UserWithAnyGroups{
		Login:  "alice",
		Groups: []AnyGroup{},
	}
	target := Document{Owner: "user", Type: "public"}

	allowed, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true for empty collection (vacuous truth), got false")
	}
}

func TestAllConditionAllConditionsTypes(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		policy   base.RawPolicy
		source   UserWithAnyGroups
		target   Document
		want     bool
		wantErr  bool
	}{
		{
			name: "All with Eq condition",
			policy: base.RawPolicy{
				Name:   "all-eq",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						All: map[string]base.Condition{
							"item:name": {Eq: "admin"},
						},
					},
				},
			},
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "admin", Role: "admin"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "All with Neq condition",
			policy: base.RawPolicy{
				Name:   "all-neq",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						All: map[string]base.Condition{
							"item:name": {Neq: "user"},
						},
					},
				},
			},
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "admin", Role: "admin"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "All with In condition",
			policy: base.RawPolicy{
				Name:   "all-in",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						All: map[string]base.Condition{
							"item:name": {In: []any{"admin", "moderator"}},
						},
					},
				},
			},
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "moderator", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
			engine, err := guardian.NewGuardianFromPolices(nil, []base.RawPolicy{tt.policy}, config)
			if err != nil {
				t.Fatalf("failed to create engine: %v", err)
			}

			allowed, err := engine.Evaluate(ctx, tt.source, tt.target, "user:read")
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if allowed != tt.want {
				t.Errorf("Evaluate() = %v, want %v", allowed, tt.want)
			}
		})
	}
}

func TestAllConditionMultipleElements(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "all-multiple-elements",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:groups": {
					All: map[string]base.Condition{
						"item:name": {Eq: "admin"},
					},
				},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name   string
		source UserWithAnyGroups
		target Document
		want   bool
	}{
		{
			name: "all elements match",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "admin", Role: "admin"},
					{Name: "admin", Role: "guest"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "first element does not match",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "user"},
					{Name: "admin", Role: "admin"},
					{Name: "admin", Role: "guest"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "middle element does not match",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "user", Role: "user"},
					{Name: "admin", Role: "guest"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "last element does not match",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "admin", Role: "admin"},
					{Name: "user", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "no elements match",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "user"},
					{Name: "guest", Role: "guest"},
					{Name: "moderator", Role: "moderator"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, err := engine.Evaluate(ctx, tt.source, tt.target, "user:read")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if allowed != tt.want {
				t.Errorf("Evaluate() = %v, want %v", allowed, tt.want)
			}
		})
	}
}

func TestAllVsAnyCondition(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name   string
		policy base.RawPolicy
		source UserWithAnyGroups
		target Document
		want   bool
	}{
		{
			name: "Any - at least one element satisfies",
			policy: base.RawPolicy{
				Name:   "any-condition",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						Any: map[string]base.Condition{
							"item:name": {Eq: "admin"},
						},
					},
				},
			},
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "user"},
					{Name: "admin", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true, // At least one group has name="admin"
		},
		{
			name: "All - all elements must satisfy",
			policy: base.RawPolicy{
				Name:   "all-condition",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						All: map[string]base.Condition{
							"item:name": {Eq: "admin"},
						},
					},
				},
			},
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "user"},
					{Name: "admin", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false, // Not all groups have name="admin"
		},
		{
			name: "All - all elements satisfy",
			policy: base.RawPolicy{
				Name:   "all-condition-all-match",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						All: map[string]base.Condition{
							"item:name": {Eq: "admin"},
						},
					},
				},
			},
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "admin", Role: "admin"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true, // All groups have name="admin"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
			engine, err := guardian.NewGuardianFromPolices(nil, []base.RawPolicy{tt.policy}, config)
			if err != nil {
				t.Fatalf("failed to create engine: %v", err)
			}

			allowed, err := engine.Evaluate(ctx, tt.source, tt.target, "user:read")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if allowed != tt.want {
				t.Errorf("Evaluate() = %v, want %v", allowed, tt.want)
			}
		})
	}
}
