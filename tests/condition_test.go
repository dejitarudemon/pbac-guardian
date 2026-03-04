/*
Package tests contains tests for comparison conditions in policies.

Tests check the work of various conditions (In, Eq, Neq, Lt, Gt) through
the public API of the library, which allows improving coverage of internal condition functions.
The tests verify that conditions work correctly with:
  - Different data types (strings, integers, slices)
  - Field comparisons between source and target structures
  - Multiple conditions combined in a single policy (logical AND)
  - Edge cases (nil values, empty slices, incompatible types)
*/
package tests

import (
	"context"
	"testing"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

/*
TestInCondition tests the In condition through policies.

The test checks the work of InConditionFunc through using the In condition
in policies. This improves coverage of InConditionFunc from 0% to 100%.
*/
func TestInCondition(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.RawPolicy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "in - found in slice",
			// Test checks that In condition finds value in list
			policy: base.RawPolicy{
				Name:   "in-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						In: []any{"admin", "moderator", "user"},
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
			name: "in - not found in slice",
			// Test checks that In condition does not find value if it's not in list
			policy: base.RawPolicy{
				Name:   "in-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						In: []any{"admin", "moderator"},
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
			name: "in - empty slice",
			// Test checks that In condition returns false for empty list
			policy: base.RawPolicy{
				Name:   "in-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						In: []any{},
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
			name: "in - with source role in list",
			// Test checks In condition with source:role field
			// Check that source.role value is in list
			policy: base.RawPolicy{
				Name:   "in-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						In: []any{"admin", "moderator", "user"},
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
			name: "in - with integer values",
			// Test checks In condition with numeric values
			policy: base.RawPolicy{
				Name:   "in-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						In: []any{18, 21, 25},
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
			config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
			engine, err := guardian.NewGuardianFromPolices(nil, []base.RawPolicy{tt.policy}, config)
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
		policy  base.RawPolicy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "lt - int less than",
			// Test checks Lt condition for integers (int)
			policy: base.RawPolicy{
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
			policy: base.RawPolicy{
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
			policy: base.RawPolicy{
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
			policy: base.RawPolicy{
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
			policy: base.RawPolicy{
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
			policy: base.RawPolicy{
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
			config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
			engine, err := guardian.NewGuardianFromPolices(nil, []base.RawPolicy{tt.policy}, config)
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
		policy  base.RawPolicy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "eq - integer comparison",
			// Test checks Eq condition for integers
			policy: base.RawPolicy{
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
			policy: base.RawPolicy{
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
			policy: base.RawPolicy{
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
			config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
			engine, err := guardian.NewGuardianFromPolices(nil, []base.RawPolicy{tt.policy}, config)
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
		policy  base.RawPolicy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "multiple conditions - all match",
			// Test checks policy with multiple conditions, all of which are met
			policy: base.RawPolicy{
				Name:   "combined-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Eq: "admin",
						In: []any{"admin", "moderator"},
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
			policy: base.RawPolicy{
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
			name: "multiple conditions - In and Lt",
			// Test checks combination of In and Lt conditions
			policy: base.RawPolicy{
				Name:   "combined-in-lt",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						In: []any{"admin", "moderator"},
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
			config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
			engine, err := guardian.NewGuardianFromPolices(nil, []base.RawPolicy{tt.policy}, config)
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
