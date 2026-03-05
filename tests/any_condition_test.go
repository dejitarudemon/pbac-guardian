/*
Package tests contains tests for Any condition functionality.

Tests verify that the library correctly handles Any conditions for checking
if at least one element in a collection satisfies all nested conditions.

Logic:
  - OR logic for elements: at least one element in the collection must satisfy
  - AND logic for conditions: all conditions in Any must be satisfied for one element

Example:

  - Collection: [Group{name: "admin", role: "user"}, Group{name: "user", role: "admin"}]

  - Condition: Any { item:name: {Eq: "admin"}, item:role: {Eq: "admin"} }

  - Result: false (no element has both name="admin" AND role="admin")

  - Collection: [Group{name: "admin", role: "admin"}, Group{name: "user", role: "user"}]

  - Condition: Any { item:name: {Eq: "admin"}, item:role: {Eq: "admin"} }

  - Result: true (first element has both name="admin" AND role="admin")
*/
package tests

import (
	"context"
	"testing"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

// AnyGroup represents a group with name and role for Any condition tests
type AnyGroup struct {
	Name string `pbac-guardian:"name"`
	Role string `pbac-guardian:"role"`
}

// UserWithAnyGroups represents a user with groups for Any condition tests
type UserWithAnyGroups struct {
	Login  string     `pbac-guardian:"login"`
	Groups []AnyGroup `pbac-guardian:"groups"`
}

func TestAnyConditionBasic(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "any-group-admin",
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
			name: "group with admin name",
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
			name: "group without admin name",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "multiple groups, one with admin",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "user"},
					{Name: "admin", Role: "admin"},
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

func TestAnyConditionMultipleFields(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "any-group-admin-name-and-role",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:groups": {
					Any: map[string]base.Condition{
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
			name: "group with admin name but not admin role",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false, // Both name and role must be "admin"
		},
		{
			name: "group with admin role but not admin name",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "admin"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false, // Both name and role must be "admin"
		},
		{
			name: "group with both admin name and role",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "admin"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true, // Both conditions satisfied
		},
		{
			name: "group with neither admin name nor role",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "multiple groups, one with both admin name and role",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "user"},
					{Name: "admin", Role: "admin"},
					{Name: "guest", Role: "guest"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true, // At least one group satisfies both conditions
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

func TestAnyConditionWithOtherConditions(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "any-group-with-role-check",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:login": {Eq: "alice"},
				"source:groups": {
					Any: map[string]base.Condition{
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
			name: "correct login and admin group",
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
			name: "wrong login but admin group",
			source: UserWithAnyGroups{
				Login: "bob",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "correct login but no admin group",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
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

func TestAnyConditionWithComparison(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "any-group-name-eq-target",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:groups": {
					Any: map[string]base.Condition{
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
			name: "group name matches target owner",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "alice", Role: "user"},
				},
			},
			target: Document{Owner: "alice", Type: "public"},
			want:   true,
		},
		{
			name: "group name does not match target owner",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
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

func TestAnyConditionEmptyCollection(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "any-empty-collection",
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
	if allowed {
		t.Errorf("expected allowed=false for empty collection, got true")
	}
}

func TestAnyConditionAllConditionsTypes(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.RawPolicy
		source  UserWithAnyGroups
		target  Document
		want    bool
		wantErr bool
	}{
		{
			name: "Any with Eq condition",
			policy: base.RawPolicy{
				Name:   "any-eq",
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
					{Name: "admin", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "Any with Neq condition",
			policy: base.RawPolicy{
				Name:   "any-neq",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						Any: map[string]base.Condition{
							"item:name": {Neq: "user"},
						},
					},
				},
			},
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
			name: "Any with Lt condition",
			policy: base.RawPolicy{
				Name:   "any-lt",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						Any: map[string]base.Condition{
							"item:name": {Lt: "z"},
						},
					},
				},
			},
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
			name: "Any with Gt condition",
			policy: base.RawPolicy{
				Name:   "any-gt",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						Any: map[string]base.Condition{
							"item:name": {Gt: "a"},
						},
					},
				},
			},
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
			name: "Any with In condition",
			policy: base.RawPolicy{
				Name:   "any-in",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						Any: map[string]base.Condition{
							"item:name": {In: []any{"admin", "moderator"}},
						},
					},
				},
			},
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
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

func TestAnyConditionMultipleElements(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "any-multiple-elements",
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
			name: "first element matches",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "admin", Role: "user"},
					{Name: "user", Role: "user"},
					{Name: "guest", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "middle element matches",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "user"},
					{Name: "admin", Role: "user"},
					{Name: "guest", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "last element matches",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "user"},
					{Name: "guest", Role: "user"},
					{Name: "admin", Role: "user"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "no element matches",
			source: UserWithAnyGroups{
				Login: "alice",
				Groups: []AnyGroup{
					{Name: "user", Role: "user"},
					{Name: "guest", Role: "user"},
					{Name: "moderator", Role: "user"},
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
