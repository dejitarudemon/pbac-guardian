/*
Package tests contains tests for map support in policies.

Tests verify that the library correctly handles map[string]any types
as source and target entities, allowing access to map values through
paths like "source:key" or "target:nested:key".

The tests check:
  - Access to map values through paths
  - Nested maps support
  - All condition types with maps (Eq, Neq, Lt, Gt, Le, Ge, Contains)
  - Error handling for non-existent keys
  - Mixed scenarios (map in source, struct in target, etc.)
*/
package tests

import (
	"context"
	"testing"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

/*
TestMapSourceBasic tests basic map usage as source entity.

The test verifies that map[string]any can be used as source
and values can be accessed through paths like "source:key".
*/
func TestMapSourceBasic(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-source-eq",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
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
		name    string
		source  map[string]any
		target  Document
		want    bool
		wantErr bool
	}{
		{
			name: "map with matching role",
			source: map[string]any{
				"role": "admin",
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "map with non-matching role",
			source: map[string]any{
				"role": "user",
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "map with missing key",
			source: map[string]any{
				"name": "alice",
			},
			target:  Document{Owner: "user", Type: "public"},
			want:    false,
			wantErr: true, // Missing key should return error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

/*
TestMapTargetBasic tests basic map usage as target entity.

The test verifies that map[string]any can be used as target
and values can be accessed through paths like "target:key".
*/
func TestMapTargetBasic(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-target-eq",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"target:type": {Eq: "public"},
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
		source User
		target map[string]any
		want   bool
	}{
		{
			name:   "map with matching type",
			source: User{Name: "alice", Role: "user"},
			target: map[string]any{
				"type": "public",
			},
			want: true,
		},
		{
			name:   "map with non-matching type",
			source: User{Name: "alice", Role: "user"},
			target: map[string]any{
				"type": "private",
			},
			want: false,
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

/*
TestMapNested tests nested map access.

The test verifies that nested maps can be accessed through
paths like "source:parent:child".
*/
func TestMapNested(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-nested",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:user:role": {Eq: "admin"},
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
		name    string
		source  map[string]any
		target  Document
		want    bool
		wantErr bool
	}{
		{
			name: "nested map with matching role",
			source: map[string]any{
				"user": map[string]any{
					"role": "admin",
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "nested map with non-matching role",
			source: map[string]any{
				"user": map[string]any{
					"role": "user",
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "nested map with missing parent key",
			source: map[string]any{
				"name": "alice",
			},
			target:  Document{Owner: "user", Type: "public"},
			want:    false,
			wantErr: true, // Missing key should return error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

/*
TestMapConditionsEq tests Eq condition with maps.
*/
func TestMapConditionsEq(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-eq",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()

	source := map[string]any{
		"role": "admin",
	}
	target := Document{Owner: "user", Type: "public"}

	allowed, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true, got false")
	}
}

/*
TestMapConditionsNeq tests Neq condition with maps.
*/
func TestMapConditionsNeq(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-neq",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Neq: "guest"},
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
		source map[string]any
		want   bool
	}{
		{
			name: "role not equal to guest",
			source: map[string]any{
				"role": "admin",
			},
			want: true,
		},
		{
			name: "role equal to guest",
			source: map[string]any{
				"role": "guest",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := Document{Owner: "user", Type: "public"}
			allowed, err := engine.Evaluate(ctx, tt.source, target, "user:read")
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

/*
TestMapConditionsLt tests Lt condition with maps.
*/
func TestMapConditionsLt(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-lt",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:age": {Lt: 18},
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
		source map[string]any
		want   bool
	}{
		{
			name: "age less than 18",
			source: map[string]any{
				"age": 15,
			},
			want: true,
		},
		{
			name: "age equal to 18",
			source: map[string]any{
				"age": 18,
			},
			want: false,
		},
		{
			name: "age greater than 18",
			source: map[string]any{
				"age": 25,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := Document{Owner: "user", Type: "public"}
			allowed, err := engine.Evaluate(ctx, tt.source, target, "user:read")
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

/*
TestMapConditionsGt tests Gt condition with maps.
*/
func TestMapConditionsGt(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-gt",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:age": {Gt: 18},
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
		source map[string]any
		want   bool
	}{
		{
			name: "age greater than 18",
			source: map[string]any{
				"age": 25,
			},
			want: true,
		},
		{
			name: "age equal to 18",
			source: map[string]any{
				"age": 18,
			},
			want: false,
		},
		{
			name: "age less than 18",
			source: map[string]any{
				"age": 15,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := Document{Owner: "user", Type: "public"}
			allowed, err := engine.Evaluate(ctx, tt.source, target, "user:read")
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

/*
TestMapConditionsLe tests Le condition with maps.
*/
func TestMapConditionsLe(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-le",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:age": {Le: 18},
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
		source map[string]any
		want   bool
	}{
		{
			name: "age less than 18",
			source: map[string]any{
				"age": 15,
			},
			want: true,
		},
		{
			name: "age equal to 18",
			source: map[string]any{
				"age": 18,
			},
			want: true,
		},
		{
			name: "age greater than 18",
			source: map[string]any{
				"age": 25,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := Document{Owner: "user", Type: "public"}
			allowed, err := engine.Evaluate(ctx, tt.source, target, "user:read")
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

/*
TestMapConditionsGe tests Ge condition with maps.
*/
func TestMapConditionsGe(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-ge",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:age": {Ge: 18},
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
		source map[string]any
		want   bool
	}{
		{
			name: "age greater than 18",
			source: map[string]any{
				"age": 25,
			},
			want: true,
		},
		{
			name: "age equal to 18",
			source: map[string]any{
				"age": 18,
			},
			want: true,
		},
		{
			name: "age less than 18",
			source: map[string]any{
				"age": 15,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := Document{Owner: "user", Type: "public"}
			allowed, err := engine.Evaluate(ctx, tt.source, target, "user:read")
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

/*
TestMapConditionsContains tests Contains condition with maps.
*/
func TestMapConditionsContains(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-contains",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {
					In: []any{"admin", "moderator"},
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
		source map[string]any
		want   bool
	}{
		{
			name: "role in list",
			source: map[string]any{
				"role": "admin",
			},
			want: true,
		},
		{
			name: "role not in list",
			source: map[string]any{
				"role": "user",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := Document{Owner: "user", Type: "public"}
			allowed, err := engine.Evaluate(ctx, tt.source, target, "user:read")
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

/*
TestMapFieldComparison tests field comparison between map and struct.
*/
func TestMapFieldComparison(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-field-comparison",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:name": {Eq: "target:owner"},
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
		source map[string]any
		target Document
		want   bool
	}{
		{
			name: "map source matches struct target",
			source: map[string]any{
				"name": "alice",
			},
			target: Document{Owner: "alice", Type: "public"},
			want:   true,
		},
		{
			name: "map source does not match struct target",
			source: map[string]any{
				"name": "bob",
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

/*
TestMapInvalidKey tests error handling for non-existent keys in maps.
*/
func TestMapInvalidKey(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-invalid-key",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:missing": {Eq: "value"},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()

	source := map[string]any{
		"role": "admin",
	}
	target := Document{Owner: "user", Type: "public"}

	// Should return error when key doesn't exist
	allowed, err := engine.Evaluate(ctx, source, target, "user:read")
	if err == nil {
		t.Errorf("expected error for missing key, got nil")
	}
	if allowed {
		t.Errorf("expected allowed=false for missing key, got true")
	}
}

/*
TestMapMixedSourceTarget tests mixed scenarios with map and struct.
*/
func TestMapMixedSourceTarget(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-mixed",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
				"target:type": {Eq: "public"},
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
		source map[string]any
		target Document
		want   bool
	}{
		{
			name: "map source and struct target - both match",
			source: map[string]any{
				"role": "admin",
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "map source and struct target - source doesn't match",
			source: map[string]any{
				"role": "user",
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "map source and struct target - target doesn't match",
			source: map[string]any{
				"role": "admin",
			},
			target: Document{Owner: "user", Type: "private"},
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

/*
TestMapBothSourceAndTarget tests both source and target as maps.
*/
func TestMapBothSourceAndTarget(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-both",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
				"target:type": {Eq: "public"},
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
		source map[string]any
		target map[string]any
		want   bool
	}{
		{
			name: "both maps - both match",
			source: map[string]any{
				"role": "admin",
			},
			target: map[string]any{
				"type": "public",
			},
			want: true,
		},
		{
			name: "both maps - source doesn't match",
			source: map[string]any{
				"role": "user",
			},
			target: map[string]any{
				"type": "public",
			},
			want: false,
		},
		{
			name: "both maps - target doesn't match",
			source: map[string]any{
				"role": "admin",
			},
			target: map[string]any{
				"type": "private",
			},
			want: false,
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

/*
TestMapDeeplyNested tests deeply nested map structures.
*/
func TestMapDeeplyNested(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-deeply-nested",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:user:profile:role": {Eq: "admin"},
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
		name    string
		source  map[string]any
		want    bool
		wantErr bool
	}{
		{
			name: "deeply nested map - match",
			source: map[string]any{
				"user": map[string]any{
					"profile": map[string]any{
						"role": "admin",
					},
				},
			},
			want: true,
		},
		{
			name: "deeply nested map - no match",
			source: map[string]any{
				"user": map[string]any{
					"profile": map[string]any{
						"role": "user",
					},
				},
			},
			want: false,
		},
		{
			name: "deeply nested map - missing intermediate key",
			source: map[string]any{
				"user": map[string]any{
					"name": "alice",
				},
			},
			want:    false,
			wantErr: true, // Missing key should return error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := Document{Owner: "user", Type: "public"}
			allowed, err := engine.Evaluate(ctx, tt.source, target, "user:read")
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

// Group represents a group with a name
type Group struct {
	Name string `pbac-guardian:"name"`
}

// UserWithGroups represents a user with groups stored in a map
type UserWithGroups struct {
	Login  string           `pbac-guardian:"login"`
	Groups map[string]Group `pbac-guardian:"groups"`
}

/*
TestMapFieldInStruct tests accessing fields in structures stored in map fields.

The test verifies that when a structure has a map field containing structures,
we can access fields in those structures through paths like "source:groups:admins:name".

Example:
  - User { Login string, Groups map[string]Group }
  - Group { Name string }
  - Path: "source:groups:admins:name"
  - groups - field in User (map[string]Group)
  - admins - key in map
  - name - field in Group
*/
func TestMapFieldInStruct(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-field-in-struct",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:groups:admins:name": {Eq: "Administrators"},
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
		name    string
		source  UserWithGroups
		target  Document
		want    bool
		wantErr bool
	}{
		{
			name: "map with matching group name",
			source: UserWithGroups{
				Login: "alice",
				Groups: map[string]Group{
					"admins": {
						Name: "Administrators",
					},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "map with non-matching group name",
			source: UserWithGroups{
				Login: "alice",
				Groups: map[string]Group{
					"admins": {
						Name: "Users",
					},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   false,
		},
		{
			name: "map without required key",
			source: UserWithGroups{
				Login: "alice",
				Groups: map[string]Group{
					"users": {
						Name: "Users",
					},
				},
			},
			target:  Document{Owner: "user", Type: "public"},
			want:    false,
			wantErr: true, // Missing key should return error
		},
		{
			name: "empty map",
			source: UserWithGroups{
				Login:  "alice",
				Groups: map[string]Group{},
			},
			target:  Document{Owner: "user", Type: "public"},
			want:    false,
			wantErr: true, // Missing key should return error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

/*
TestMapFieldInStructMultipleKeys tests accessing different keys in map field.
*/
func TestMapFieldInStructMultipleKeys(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-field-multiple-keys",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:groups:moderators:name": {Eq: "Moderators"},
			},
		},
	}

	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()

	source := UserWithGroups{
		Login: "bob",
		Groups: map[string]Group{
			"admins": {
				Name: "Administrators",
			},
			"moderators": {
				Name: "Moderators",
			},
			"users": {
				Name: "Users",
			},
		},
	}
	target := Document{Owner: "user", Type: "public"}

	allowed, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true, got false")
	}
}

/*
TestMapFieldInStructComparison tests field comparison with map fields in structures.
*/
func TestMapFieldInStructComparison(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "map-field-comparison",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:groups:admins:name": {Eq: "target:owner"},
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
		source UserWithGroups
		target Document
		want   bool
	}{
		{
			name: "group name matches document owner",
			source: UserWithGroups{
				Login: "alice",
				Groups: map[string]Group{
					"admins": {
						Name: "alice",
					},
				},
			},
			target: Document{Owner: "alice", Type: "public"},
			want:   true,
		},
		{
			name: "group name does not match document owner",
			source: UserWithGroups{
				Login: "alice",
				Groups: map[string]Group{
					"admins": {
						Name: "bob",
					},
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

/*
TestMapFieldInStructAllConditions tests all condition types with map fields in structures.
*/
func TestMapFieldInStructAllConditions(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.RawPolicy
		source  UserWithGroups
		target  Document
		want    bool
		wantErr bool
	}{
		{
			name: "Eq condition",
			policy: base.RawPolicy{
				Name:   "map-eq",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups:admins:name": {Eq: "Administrators"},
				},
			},
			source: UserWithGroups{
				Login: "alice",
				Groups: map[string]Group{
					"admins": {Name: "Administrators"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "Neq condition",
			policy: base.RawPolicy{
				Name:   "map-neq",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups:admins:name": {Neq: "Users"},
				},
			},
			source: UserWithGroups{
				Login: "alice",
				Groups: map[string]Group{
					"admins": {Name: "Administrators"},
				},
			},
			target: Document{Owner: "user", Type: "public"},
			want:   true,
		},
		{
			name: "Contains condition",
			policy: base.RawPolicy{
				Name:   "map-contains",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups:admins:name": {
						In: []any{"Administrators", "Moderators"},
					},
				},
			},
			source: UserWithGroups{
				Login: "alice",
				Groups: map[string]Group{
					"admins": {Name: "Administrators"},
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
