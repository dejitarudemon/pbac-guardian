/*
Package tests contains tests for time.Time support in policies.

Tests check the work of time.Time comparison through policy evaluation, verifying that
time.Time values can be used in policy conditions (Eq, Neq, Lt, Gt). The tests verify:
  - Using time.Time in conditions (Eq, Neq, Lt, Gt)
  - Comparing time.Time with literal time values
  - Comparing time.Time fields from source and target structures
  - Error handling for incompatible types
*/
package tests

import (
	"context"
	"testing"
	"time"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

// UserWithTime represents a user structure with time field for testing
type UserWithTime struct {
	Name      string    `pbac-guardian:"name"`
	Role      string    `pbac-guardian:"role"`
	CreatedAt time.Time `pbac-guardian:"created_at"`
	UpdatedAt time.Time `pbac-guardian:"updated_at"`
}

// DocumentWithTime represents a document structure with time field for testing
type DocumentWithTime struct {
	Owner     string    `pbac-guardian:"owner"`
	Type      string    `pbac-guardian:"type"`
	ExpiresAt time.Time `pbac-guardian:"expires_at"`
	Published time.Time `pbac-guardian:"published"`
}

/*
TestTime_Eq tests equality condition with time.Time values.

The test checks that time.Time values can be used in Eq conditions,
comparing them with literal time values and structure fields.
*/
func TestTime_Eq(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	past := now.Add(-24 * time.Hour)

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
			name: "time equals literal - match",
			// Test checks that time field equals literal time value
			policy: base.RawPolicy{
				Name:   "time-eq-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Eq: now, // source:created_at == now
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time equals literal - no match",
			// Test checks that time field doesn't equal different literal time value
			policy: base.RawPolicy{
				Name:   "time-eq-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Eq: past, // source:created_at == past (false)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "time equals target field - match",
			// Test checks that source time field equals target time field
			policy: base.RawPolicy{
				Name:   "time-eq-target-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Eq: "target:published", // source:created_at == target:published
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public", Published: now},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time equals target field - no match",
			// Test checks that source time field doesn't equal target time field
			policy: base.RawPolicy{
				Name:   "time-eq-target-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Eq: "target:published", // source:created_at == target:published
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public", Published: past},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
TestTime_Neq tests inequality condition with time.Time values.

The test checks that time.Time values can be used in Neq conditions.
*/
func TestTime_Neq(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	past := now.Add(-24 * time.Hour)

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
			name: "time not equals literal - match",
			// Test checks that time field doesn't equal literal time value
			policy: base.RawPolicy{
				Name:   "time-neq-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Neq: past, // source:created_at != past (true)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time not equals literal - no match",
			// Test checks that time field equals literal time value (condition fails)
			policy: base.RawPolicy{
				Name:   "time-neq-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Neq: now, // source:created_at != now (false)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
TestTime_Lt tests less-than condition with time.Time values.

The test checks that time.Time values can be used in Lt conditions.
*/
func TestTime_Lt(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	past := now.Add(-24 * time.Hour)
	future := now.Add(24 * time.Hour)

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
			name: "time less than literal - match",
			// Test checks that time field is less than literal time value
			policy: base.RawPolicy{
				Name:   "time-lt-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Lt: future, // source:created_at < future (true)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time less than literal - no match",
			// Test checks that time field is not less than literal time value
			policy: base.RawPolicy{
				Name:   "time-lt-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Lt: past, // source:created_at < past (false, now > past)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "time less than literal - equal",
			// Test checks that time field is not less than equal literal time value
			policy: base.RawPolicy{
				Name:   "time-lt-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Lt: now, // source:created_at < now (false, equal)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "time less than target field - match",
			// Test checks that source time field is less than target time field
			policy: base.RawPolicy{
				Name:   "time-lt-target-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Lt: "target:expires_at", // source:created_at < target:expires_at
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public", ExpiresAt: future},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time less than target field - no match",
			// Test checks that source time field is not less than target time field
			policy: base.RawPolicy{
				Name:   "time-lt-target-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Lt: "target:expires_at", // source:created_at < target:expires_at
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public", ExpiresAt: past},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
TestTime_Gt tests greater-than condition with time.Time values.

The test checks that time.Time values can be used in Gt conditions.
*/
func TestTime_Gt(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	past := now.Add(-24 * time.Hour)
	future := now.Add(24 * time.Hour)

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
			name: "time greater than literal - match",
			// Test checks that time field is greater than literal time value
			policy: base.RawPolicy{
				Name:   "time-gt-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Gt: past, // source:created_at > past (true)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time greater than literal - no match",
			// Test checks that time field is not greater than literal time value
			policy: base.RawPolicy{
				Name:   "time-gt-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Gt: future, // source:created_at > future (false, now < future)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "time greater than literal - equal",
			// Test checks that time field is not greater than equal literal time value
			policy: base.RawPolicy{
				Name:   "time-gt-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Gt: now, // source:created_at > now (false, equal)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "time greater than target field - match",
			// Test checks that source time field is greater than target time field
			policy: base.RawPolicy{
				Name:   "time-gt-target-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:updated_at": {
						Gt: "target:published", // source:updated_at > target:published
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", UpdatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public", Published: past},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time greater than target field - no match",
			// Test checks that source time field is not greater than target time field
			policy: base.RawPolicy{
				Name:   "time-gt-target-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:updated_at": {
						Gt: "target:published", // source:updated_at > target:published
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", UpdatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public", Published: future},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
TestTime_MultipleConditions tests combining time.Time conditions with other conditions.

The test checks that time.Time values can be combined with other conditions
in a single policy (logical AND).
*/
func TestTime_MultipleConditions(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	past := now.Add(-24 * time.Hour)
	future := now.Add(24 * time.Hour)

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
			name: "time and role - both match",
			// Test checks that both time condition and role condition are met
			policy: base.RawPolicy{
				Name:   "time-multiple-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Gt: past, // source:created_at > past
					},
					"source:role": {
						Eq: "user", // source:role == "user"
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time and role - one fails",
			// Test checks that policy fails when one condition doesn't match
			policy: base.RawPolicy{
				Name:   "time-multiple-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Gt: past, // source:created_at > past (true)
					},
					"source:role": {
						Eq: "admin", // source:role == "admin" (false, role is "user")
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "multiple time conditions - both match",
			// Test checks that multiple time conditions can be used in one policy
			policy: base.RawPolicy{
				Name:   "time-multiple-time-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Gt: past, // source:created_at > past
					},
					"source:updated_at": {
						Lt: future, // source:updated_at < future
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now, UpdatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time range check - within range",
			// Test checks that time is within a range using Lt and Gt
			policy: base.RawPolicy{
				Name:   "time-range-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Gt: past,   // source:created_at > past
						Lt: future, // source:created_at < future
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time range check - outside range",
			// Test checks that time outside range fails
			policy: base.RawPolicy{
				Name:   "time-range-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						Gt: now,    // source:created_at > now (false, equal)
						Lt: future, // source:created_at < future (true)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
TestTime_In tests In condition with time.Time values.

The test checks that time.Time values can be used in In conditions.
*/
func TestTime_In(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	past := now.Add(-24 * time.Hour)
	future := now.Add(24 * time.Hour)

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
			name: "time in list - found",
			// Test checks that time value is found in In list
			policy: base.RawPolicy{
				Name:   "time-in-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						In: []any{past, now, future}, // source:created_at in [past, now, future]
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time in list - not found",
			// Test checks that time value is not found in In list
			policy: base.RawPolicy{
				Name:   "time-in-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:created_at": {
						In: []any{past, future}, // source:created_at not in [past, future]
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  DocumentWithTime{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
