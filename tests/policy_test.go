/*
Package tests contains tests for public methods of Policy structure in base package.

Tests check only exported methods:
  - Evaluate - policy evaluation for given source, target and action
    (requires ConditionsMap to be set with condition functions)
  - IsValid - policy validation

The tests verify that Policy methods work correctly with various data types,
handle multiple conditions, and properly validate policy structure.
*/
package tests

import (
	"context"
	"testing"

	"github.com/dejitarudemon/pbac-guardian/internal/base"
	"github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

/*
Test structures for checking policy functionality.

These structures are used for testing Policy methods
with various data types and nested structures.
*/

// PolicyTestUser represents a user for testing policies
type PolicyTestUser struct {
	Name string `pbac-guardian:"name"` // User name
	Role string `pbac-guardian:"role"` // User role
	Age  int    `pbac-guardian:"age"`  // User age
}

// PolicyTestDocument represents a document for testing policies
type PolicyTestDocument struct {
	Owner string   `pbac-guardian:"owner"` // Document owner
	Type  string   `pbac-guardian:"type"`  // Document type
	Tags  []string `pbac-guardian:"tags"`  // Document tags
}

/*
TestPolicyEvaluate tests the public Evaluate method for policy evaluation.

The test checks:
  - Policy action matching passed action
  - Policy condition checking
  - Handling of multiple conditions (logical AND)
  - Field comparison from different structures (source and target)
  - Return of correct result in various scenarios
*/
func TestPolicyEvaluate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		raw     base.RawPolicy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "matching action and condition",
			// Test checks successful policy evaluation when action matches and condition is met
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
				},
			},
			source:  PolicyTestUser{Name: "admin", Role: "admin"},
			target:  PolicyTestDocument{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "non-matching action",
			// Test checks that policy is not applied when action does not match
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
				},
			},
			source:  PolicyTestUser{Name: "admin", Role: "admin"},
			target:  PolicyTestDocument{Owner: "user", Type: "public"},
			action:  "user:write", // Different action
			want:    false,
			wantErr: false,
		},
		{
			name: "non-matching condition",
			// Test checks that policy does not pass when condition is not met
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
				},
			},
			source:  PolicyTestUser{Name: "user", Role: "user"}, // Role is not "admin"
			target:  PolicyTestDocument{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "multiple conditions - all match",
			// Test checks that all conditions must be met (logical AND)
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
					"source:age":  {Lt: 100}, // Age less than 100
				},
			},
			source:  PolicyTestUser{Name: "admin", Role: "admin", Age: 25}, // Both conditions met
			target:  PolicyTestDocument{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "multiple conditions - one fails",
			// Test checks that if at least one condition is not met, policy does not pass
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
					"source:age":  {Lt: 10}, // Age must be less than 10
				},
			},
			source:  PolicyTestUser{Name: "admin", Role: "admin", Age: 25}, // Age 25, condition not met
			target:  PolicyTestDocument{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "compare fields from different structures",
			// Test checks field comparison from different structures (source and target)
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:name": {Eq: "target:owner"}, // Compare name from source with owner from target
				},
			},
			source:  PolicyTestUser{Name: "alice", Role: "user"},
			target:  PolicyTestDocument{Owner: "alice", Type: "public"}, // Names match
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create policy using NewPolicy with default conditions map and nil cache
			policy, err := base.NewPolicy(tt.raw, &implemented.DefaultConditionsMap, nil, nil)
			if err != nil {
				t.Fatalf("failed to create policy: %v", err)
			}
			// Use empty sessionID for direct Policy.Evaluate tests
			got, err := policy.Evaluate(ctx, tt.source, tt.target, tt.action, "")
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got != tt.want {
					t.Errorf("Evaluate() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

/*
TestPolicyIsValid tests the public IsValid method for policy validation.

The test checks:
  - Action format validation (minimum 2 parts separated by ":")
  - Check for absence of empty parts in action
  - Validation of all paths in conditions
  - Correctness of handling valid and invalid policies
*/
func TestPolicyIsValid(t *testing.T) {
	tests := []struct {
		name    string
		raw     base.RawPolicy
		wantErr bool
	}{
		{
			name: "valid policy",
			// Test checks valid policy with all correct parameters
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid action - too short",
			// Test checks action validation with insufficient number of parts
			raw: base.RawPolicy{
				Name:   "test",
				Action: "read", // Only one part, need minimum 2
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid action - empty part",
			// Test checks action validation with empty part
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user::read", // Empty part between separators
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid path in conditions",
			// Test checks path validation in policy conditions
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"invalid": {Eq: "admin"}, // Invalid path without separator ":"
				},
			},
			wantErr: true,
		},
		{
			name: "valid policy with path in condition value",
			// Test checks that paths in condition values (Eq, Neq, etc.) are validated
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "target:owner"}, // Valid path in condition value
				},
			},
			wantErr: false,
		},
		{
			name: "invalid path in condition value",
			// Test checks that invalid paths (containing ":") in condition values are rejected
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "invalid:path"}, // Invalid path (contains ":" but invalid entity)
				},
			},
			wantErr: true,
		},
		{
			name: "valid policy with Any condition using item:",
			// Test checks that item: can be used within Any conditions
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						Any: map[string]base.Condition{
							"item:name": {Eq: "admin"}, // item: is allowed in Any
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid policy with All condition using item:",
			// Test checks that item: can be used within All conditions
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						All: map[string]base.Condition{
							"item:name": {Eq: "admin"}, // item: is allowed in All
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - item: used outside Any/All",
			// Test checks that item: cannot be used outside Any/All conditions
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"item:name": {Eq: "admin"}, // item: is not allowed at top level
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - item: in condition value outside Any/All",
			// Test checks that item: cannot be used in condition values outside Any/All
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "item:name"}, // item: is not allowed in condition value at top level
				},
			},
			wantErr: true,
		},
		{
			name: "valid policy with nested Any and paths",
			// Test checks that all paths in nested Any are validated
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						Any: map[string]base.Condition{
							"item:name": {Eq: "target:owner"}, // Valid path in nested condition
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid path in nested Any condition value",
			// Test checks that invalid paths (containing ":") in nested Any condition values are rejected
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:groups": {
						Any: map[string]base.Condition{
							"item:name": {Eq: "invalid:path"}, // Invalid path (contains ":" but invalid entity)
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid policy with In condition containing paths",
			// Test checks that paths in In condition slice are validated
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						In: []any{"admin", "target:owner"}, // Valid path in In slice
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid path in In condition",
			// Test checks that invalid paths (containing ":") in In condition slice are rejected
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						In: []any{"admin", "invalid:path"}, // Invalid path (contains ":" but invalid entity)
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid policy with multiple condition types and paths",
			// Test checks that all condition types with paths are validated
			raw: base.RawPolicy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Eq:  "target:minAge",
						Neq: "target:maxAge",
						Lt:  "target:limit",
						Gt:  "target:threshold",
						Le:  "target:max",
						Ge:  "target:min",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// NewPolicy calls IsValid internally, so we test validation through it
			_, err := base.NewPolicy(tt.raw, &implemented.DefaultConditionsMap, nil, nil)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
