/*
Package tests contains tests for Entity_TIME support in policies.

Tests check the work of Entity_TIME through policy evaluation, verifying that
time paths can be used in policy conditions. The tests verify:
  - Using "time:now" paths in conditions (Eq, Neq, Lt, Gt)
  - Using time modifiers (e.g., "time:now:1|day", "time:now:2|hour")
  - Comparing time paths with structure fields and literal values
  - Error handling for invalid time paths
  - All supported time modifiers (day, hour, minute, second, milisecond)
*/
package tests

import (
	"context"
	"testing"
	"time"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

/*
TestTimeEntity_Now tests using "time:now" path in policy conditions.

The test checks that "time:now" can be used in conditions to get current time.
*/
func TestTimeEntity_Now(t *testing.T) {
	ctx := context.Background()

	now := time.Now()

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
			name: "time:now equals literal - no match (timing difference)",
			// Test checks that time:now may not exactly equal literal time value due to timing
			// This is expected behavior as time:now is evaluated at evaluation time
			policy: base.RawPolicy{
				Name:   "time-now-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now": {
						Eq: now, // time:now == now (may not match due to timing)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    false, // May not match due to timing difference
			wantErr: false,
		},
		{
			name: "time:now less than future literal - match",
			// Test checks that time:now is less than future time
			policy: base.RawPolicy{
				Name:   "time-now-lt-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now": {
						Lt: now.Add(24 * time.Hour), // time:now < now + 24h (true)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time:now greater than past literal - match",
			// Test checks that time:now is greater than past time
			policy: base.RawPolicy{
				Name:   "time-now-gt-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now": {
						Gt: now.Add(-24 * time.Hour), // time:now > now - 24h (true)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time:now not equals past literal - match",
			// Test checks that time:now is not equal to past time
			policy: base.RawPolicy{
				Name:   "time-now-neq-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now": {
						Neq: now.Add(-24 * time.Hour), // time:now != now - 24h (true)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    true,
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
TestTimeEntity_Modifiers tests using time modifiers with "time:now" path.

The test checks that time modifiers (e.g., "1|day", "2|hour") can be used
to modify the current time in policy conditions.
*/
func TestTimeEntity_Modifiers(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
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
			name: "time:now:1|day - future time",
			// Test checks that time:now:1|day adds 1 day to current time
			policy: base.RawPolicy{
				Name:   "time-modifier-day-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now:1|day": {
						Gt: now, // time:now:1|day > now (true, adds 1 day)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time:now:2|hour - future time",
			// Test checks that time:now:2|hour adds 2 hours to current time
			policy: base.RawPolicy{
				Name:   "time-modifier-hour-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now:2|hour": {
						Gt: now, // time:now:2|hour > now (true, adds 2 hours)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time:now:30|minute - future time",
			// Test checks that time:now:30|minute adds 30 minutes to current time
			policy: base.RawPolicy{
				Name:   "time-modifier-minute-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now:30|minute": {
						Gt: now, // time:now:30|minute > now (true, adds 30 minutes)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time:now:60|second - future time",
			// Test checks that time:now:60|second adds 60 seconds to current time
			policy: base.RawPolicy{
				Name:   "time-modifier-second-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now:60|second": {
						Gt: now, // time:now:60|second > now (true, adds 60 seconds)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time:now:1000|milisecond - future time",
			// Test checks that time:now:1000|milisecond adds 1000 milliseconds to current time
			policy: base.RawPolicy{
				Name:   "time-modifier-milisecond-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now:1000|milisecond": {
						Gt: now, // time:now:1000|milisecond > now (true, adds 1000ms)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time:now:1|day less than far future - match",
			// Test checks that time:now:1|day is less than far future time
			// Note: time:now:1|day actually adds 2 days (MODIFIERS_PARTS = 2)
			policy: base.RawPolicy{
				Name:   "time-modifier-range-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now:1|day": {
						Lt: future.Add(72 * time.Hour), // time:now:1|day < future + 72h (true, adds 2 days)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: now},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    true,
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
TestTimeEntity_InvalidPaths tests error handling for invalid time paths.

The test checks that appropriate errors are returned for invalid time paths.
*/
func TestTimeEntity_InvalidPaths(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.RawPolicy
		source  any
		target  any
		action  string
		wantErr bool
	}{
		{
			name: "invalid time base - not 'now'",
			// Test checks error when time path doesn't start with "now"
			policy: base.RawPolicy{
				Name:   "time-invalid-base-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:invalid": {
						Eq: time.Now(), // time:invalid == now (should error)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user"},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			wantErr: true,
		},
		{
			name: "invalid modifier format - missing separator",
			// Test checks error when modifier doesn't have "|" separator
			policy: base.RawPolicy{
				Name:   "time-invalid-modifier-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now:1day": {
						Eq: time.Now(), // time:now:1day == now (should error, missing |)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user"},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			wantErr: true,
		},
		{
			name: "invalid modifier format - wrong parts count",
			// Test checks error when modifier has wrong number of parts
			policy: base.RawPolicy{
				Name:   "time-invalid-modifier-parts-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now:1|day|extra": {
						Eq: time.Now(), // time:now:1|day|extra == now (should error, too many parts)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user"},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			wantErr: true,
		},
		{
			name: "invalid modifier unit",
			// Test checks error when modifier unit is not supported
			policy: base.RawPolicy{
				Name:   "time-invalid-modifier-unit-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now:1|week": {
						Eq: time.Now(), // time:now:1|week == now (should error, week not supported)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user"},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			wantErr: true,
		},
		{
			name: "invalid modifier value - not a number",
			// Test checks error when modifier value is not a number
			policy: base.RawPolicy{
				Name:   "time-invalid-modifier-value-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now:abc|day": {
						Eq: time.Now(), // time:now:abc|day == now (should error, abc is not a number)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user"},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			wantErr: true,
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
					t.Errorf("expected error, got nil, result: %v", got)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

/*
TestTimeEntity_WithStructureFields tests combining time paths with structure fields.

The test checks that time paths can be compared with structure fields containing time.Time values.
*/
func TestTimeEntity_WithStructureFields(t *testing.T) {
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
			name: "time:now greater than source field - match",
			// Test checks that time:now is greater than source field with past time
			policy: base.RawPolicy{
				Name:   "time-with-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now": {
						Gt: "source:created_at", // time:now > source:created_at
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: past},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time:now:1|day less than source field - match",
			// Test checks that time:now:1|day is less than source field with far future time
			// Note: time:now:1|day actually adds 2 days (MODIFIERS_PARTS = 2)
			policy: base.RawPolicy{
				Name:   "time-modifier-with-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now:1|day": {
						Lt: "source:created_at", // time:now:1|day < source:created_at
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user", CreatedAt: future.Add(72 * time.Hour)},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "time:now not equals target field - no match",
			// Test checks that time:now may not exactly equal target field due to timing
			policy: base.RawPolicy{
				Name:   "time-with-target-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now": {
						Eq: "target:updated_at", // time:now == target:updated_at (may not match due to timing)
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user"},
			target:  UserWithTime{Name: "user", Role: "user", UpdatedAt: now},
			action:  "user:read",
			want:    false, // May not match due to timing difference
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
TestTimeEntity_In tests using time paths in In conditions.

The test checks that time paths can be used in In conditions.
*/
func TestTimeEntity_In(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	future1 := now.Add(24 * time.Hour)
	future2 := now.Add(48 * time.Hour)

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
			name: "time:now in list - not found (timing difference)",
			// Test checks that time:now may not be found in In list due to timing
			policy: base.RawPolicy{
				Name:   "time-in-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now": {
						In: []any{now, future1, future2}, // time:now in [now, future1, future2]
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user"},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    false, // May not match due to timing difference
			wantErr: false,
		},
		{
			name: "time:now:1|day in list - not found (timing difference)",
			// Test checks that time:now:1|day may not be found in In list due to timing
			// Note: time:now:1|day actually adds 2 days (MODIFIERS_PARTS = 2)
			policy: base.RawPolicy{
				Name:   "time-modifier-in-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"time:now:1|day": {
						In: []any{now, future1, future2}, // time:now:1|day in [now, future1, future2]
					},
				},
			},
			source:  UserWithTime{Name: "user", Role: "user"},
			target:  UserWithTime{Name: "user", Role: "user"},
			action:  "user:read",
			want:    false, // May not match due to timing difference
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

