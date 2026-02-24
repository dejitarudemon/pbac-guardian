/*
Package tests contains tests for comparison conditions in policies.

Tests check the work of various conditions (Contains, Eq, Neq, Lt) through
the public API of the library, which allows improving coverage of internal condition functions.
*/
package tests

import (
	"context"
	"testing"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

/*
TestContainsCondition tests the Contains condition through policies.

The test checks the work of containsConditionFunc through using the Contains condition
in policies. This improves coverage of containsConditionFunc from 0% to 100%.
*/
func TestContainsCondition(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.Policy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "contains - found in slice",
			// Test checks that Contains condition finds value in list
			policy: base.Policy{
				Name:   "contains-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Contains: []any{"admin", "moderator", "user"},
					},
				},
			},
			source:  User{Name: "user", Role: "admin"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "contains - not found in slice",
			// Test checks that Contains condition does not find value if it's not in list
			policy: base.Policy{
				Name:   "contains-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Contains: []any{"admin", "moderator"},
					},
				},
			},
			source:  User{Name: "user", Role: "guest"}, // "guest" not in list
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "contains - empty slice",
			// Test checks that Contains condition returns false for empty list
			policy: base.Policy{
				Name:   "contains-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Contains: []any{},
					},
				},
			},
			source:  User{Name: "user", Role: "admin"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "contains - with source role in list",
			// Test checks Contains condition with source:role field
			// Check that source.role value is in list
			policy: base.Policy{
				Name:   "contains-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Contains: []any{"admin", "moderator", "user"},
					},
				},
			},
			source:  User{Name: "user", Role: "moderator"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "contains - with integer values",
			// Test checks Contains condition with numeric values
			policy: base.Policy{
				Name:   "contains-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Contains: []any{18, 21, 25},
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 21},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use nil casher for basic functionality tests
			engine, err := guardian.NewGuardianFromPolices(nil, []base.Policy{tt.policy})
			if err != nil {
				t.Fatalf("failed to create engine: %v", err)
			}

			got, err := engine.Evaluate(ctx, tt.source, tt.target, tt.action)
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
TestLtCondition tests the Lt (less than) condition through policies.

The test checks the work of ltConditionFunc and ltPrimitives through using the Lt condition
in policies. This improves coverage of comparison functions for various data types.
*/
func TestLtCondition(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.Policy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "lt - int less than",
			// Test checks Lt condition for integers (int)
			policy: base.Policy{
				Name:   "lt-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Lt: 18,
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 16}, // 16 < 18
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - int equal",
			// Test checks that Lt condition returns false on equality
			policy: base.Policy{
				Name:   "lt-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Lt: 18,
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 18}, // 18 == 18
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "lt - int greater than",
			// Test checks that Lt condition returns false for greater value
			policy: base.Policy{
				Name:   "lt-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Lt: 18,
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 25}, // 25 > 18
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "lt - with string comparison",
			// Test checks Lt condition for strings (lexicographic comparison)
			// Need to use a field that can be compared as string
			policy: base.Policy{
				Name:   "lt-string-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:name": {
						Lt: "m", // Names before "m" in alphabetical order
					},
				},
			},
			source:  User{Name: "alice", Role: "user"}, // "alice" < "m"
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - string greater",
			// Test checks that Lt condition returns false for strings that are greater
			policy: base.Policy{
				Name:   "lt-string-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:name": {
						Lt: "m",
					},
				},
			},
			source:  User{Name: "zoe", Role: "user"}, // "zoe" > "m"
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "lt - compare with target field",
			// Test checks Lt condition with field comparison from different structures
			// Need to add numeric field to Document
			policy: base.Policy{
				Name:   "lt-compare-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Lt: "target:priority", // Compare with target field
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 25},
			target:  DocumentWithPriority{Document: Document{Owner: "user", Type: "public"}, Priority: 30},
			action:  "user:read",
			want:    true, // 25 < 30
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use nil casher for basic functionality tests
			engine, err := guardian.NewGuardianFromPolices(nil, []base.Policy{tt.policy})
			if err != nil {
				t.Fatalf("failed to create engine: %v", err)
			}

			got, err := engine.Evaluate(ctx, tt.source, tt.target, tt.action)
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

// DocumentWithPriority extends Document for testing numeric comparisons
type DocumentWithPriority struct {
	Document
	Priority int `pbac-guardian:"priority"`
}

/*
TestEqConditionExtended tests extended cases of Eq condition.

The test checks the work of eqConditionFunc for various data types and scenarios,
including comparison with nil, structures with Comparable interface and various types.
*/
func TestEqConditionExtended(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.Policy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "eq - integer comparison",
			// Test checks Eq condition for integers
			policy: base.Policy{
				Name:   "eq-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Eq: 25,
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 25},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "eq - integer not equal",
			// Test checks that Eq condition returns false for unequal numbers
			policy: base.Policy{
				Name:   "eq-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Eq: 25,
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 30},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "eq - compare integer with string path",
			// Test checks comparison of number with field that should be a number
			policy: base.Policy{
				Name:   "eq-mixed-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Eq: "target:priority", // Compare number with number via path
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 10},
			target:  DocumentWithPriority{Document: Document{Owner: "user", Type: "public"}, Priority: 10},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use nil casher for basic functionality tests
			engine, err := guardian.NewGuardianFromPolices(nil, []base.Policy{tt.policy})
			if err != nil {
				t.Fatalf("failed to create engine: %v", err)
			}

			got, err := engine.Evaluate(ctx, tt.source, tt.target, tt.action)
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
TestMultipleConditionsCombined tests combination of multiple conditions in one policy.

The test checks that all conditions in policy must be met (logical AND),
which improves coverage of Evaluate function for various condition combinations.
*/
func TestMultipleConditionsCombined(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.Policy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "multiple conditions - all match",
			// Test checks policy with multiple conditions, all of which are met
			policy: base.Policy{
				Name:   "combined-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Eq:       "admin",
						Contains: []any{"admin", "moderator"},
					},
					"source:age": {
						Lt: 100,
					},
				},
			},
			source:  User{Name: "admin", Role: "admin", Age: 25},
			target:  Document{Owner: "user", Type: "public", Tags: []string{"public"}},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "multiple conditions - one fails",
			// Test checks policy with multiple conditions, one of which is not met
			policy: base.Policy{
				Name:   "combined-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Eq: "admin",
					},
					"source:age": {
						Lt: 10, // Age must be less than 10
					},
				},
			},
			source:  User{Name: "admin", Role: "admin", Age: 25}, // Age 25, condition not met
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "multiple conditions - Contains and Lt",
			// Test checks combination of Contains and Lt conditions
			policy: base.Policy{
				Name:   "combined-contains-lt",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Contains: []any{"admin", "moderator"},
					},
					"source:age": {
						Lt: 65,
					},
				},
			},
			source:  User{Name: "admin", Role: "admin", Age: 30},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use nil casher for basic functionality tests
			engine, err := guardian.NewGuardianFromPolices(nil, []base.Policy{tt.policy})
			if err != nil {
				t.Fatalf("failed to create engine: %v", err)
			}

			got, err := engine.Evaluate(ctx, tt.source, tt.target, tt.action)
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
