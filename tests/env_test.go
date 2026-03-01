/*
Package tests contains tests for environment variable support in policies.

Tests check the work of Entity_ENV through policy evaluation, verifying that
environment variables can be used in policy conditions. The tests verify:
  - Using environment variables in conditions (Eq, Neq, Contains, Lt)
  - Comparing environment variables with structure fields
  - Comparing environment variables with literal values
  - Error handling when environment variables don't exist
  - Caching behavior for environment variables
*/
package tests

import (
	"context"
	"os"
	"testing"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
	"github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

/*
TestEnvVariable_Eq tests equality condition with environment variables.

The test checks that environment variables can be used in Eq conditions,
comparing them with literal values and structure fields.
*/
func TestEnvVariable_Eq(t *testing.T) {
	ctx := context.Background()

	// Set up test environment variable
	testEnvVar := "TEST_ENV_VAR"
	testEnvValue := "test-value"
	os.Setenv(testEnvVar, testEnvValue)
	defer os.Unsetenv(testEnvVar)

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
			name: "env equals literal - match",
			// Test checks that env variable equals literal value
			policy: base.RawPolicy{
				Name:   "env-eq-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + testEnvVar: {
						Eq: testEnvValue, // env:TEST_ENV_VAR == "test-value"
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "env equals literal - no match",
			// Test checks that env variable doesn't equal different literal value
			policy: base.RawPolicy{
				Name:   "env-eq-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + testEnvVar: {
						Eq: "different-value", // env:TEST_ENV_VAR == "different-value" (false)
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "env equals source field - match",
			// Test checks that env variable equals source field value
			policy: base.RawPolicy{
				Name:   "env-eq-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Eq: "env:" + testEnvVar, // source:role == env:TEST_ENV_VAR
					},
				},
			},
			source:  User{Name: "user", Role: testEnvValue}, // Role matches env value
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "env equals source field - no match",
			// Test checks that env variable doesn't equal source field value
			policy: base.RawPolicy{
				Name:   "env-eq-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Eq: "env:" + testEnvVar, // source:role == env:TEST_ENV_VAR
					},
				},
			},
			source:  User{Name: "user", Role: "different"}, // Role doesn't match env value
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "env equals target field - match",
			// Test checks that env variable equals target field value
			policy: base.RawPolicy{
				Name:   "env-eq-target-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"target:type": {
						Eq: "env:" + testEnvVar, // target:type == env:TEST_ENV_VAR
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: testEnvValue}, // Type matches env value
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
TestEnvVariable_Neq tests inequality condition with environment variables.

The test checks that environment variables can be used in Neq conditions.
*/
func TestEnvVariable_Neq(t *testing.T) {
	ctx := context.Background()

	// Set up test environment variable
	testEnvVar := "TEST_ENV_VAR_NEQ"
	testEnvValue := "allowed-value"
	os.Setenv(testEnvVar, testEnvValue)
	defer os.Unsetenv(testEnvVar)

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
			name: "env not equals literal - match",
			// Test checks that env variable doesn't equal literal value
			policy: base.RawPolicy{
				Name:   "env-neq-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + testEnvVar: {
						Neq: "blocked-value", // env:TEST_ENV_VAR_NEQ != "blocked-value" (true)
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "env not equals literal - no match",
			// Test checks that env variable equals literal value (condition fails)
			policy: base.RawPolicy{
				Name:   "env-neq-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + testEnvVar: {
						Neq: testEnvValue, // env:TEST_ENV_VAR_NEQ != "allowed-value" (false)
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "source field not equals env - match",
			// Test checks that source field doesn't equal env variable
			policy: base.RawPolicy{
				Name:   "env-neq-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Neq: "env:" + testEnvVar, // source:role != env:TEST_ENV_VAR_NEQ
					},
				},
			},
			source:  User{Name: "user", Role: "different"}, // Role doesn't match env value
			target:  Document{Owner: "user", Type: "public"},
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
TestEnvVariable_Contains tests Contains condition with environment variables.

The test checks that environment variables can be used in Contains conditions.
*/
func TestEnvVariable_Contains(t *testing.T) {
	ctx := context.Background()

	// Set up test environment variable
	testEnvVar := "TEST_ENV_VAR_CONTAINS"
	testEnvValue := "admin"
	os.Setenv(testEnvVar, testEnvValue)
	defer os.Unsetenv(testEnvVar)

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
			name: "env in list - found",
			// Test checks that env variable value is found in Contains list
			policy: base.RawPolicy{
				Name:   "env-contains-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + testEnvVar: {
						Contains: []any{"admin", "moderator", "user"}, // env:TEST_ENV_VAR_CONTAINS in list
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "env in list - not found",
			// Test checks that env variable value is not found in Contains list
			policy: base.RawPolicy{
				Name:   "env-contains-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + testEnvVar: {
						Contains: []any{"moderator", "user"}, // env:TEST_ENV_VAR_CONTAINS not in list
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "source field in list with env value",
			// Test checks that source field value is in list that includes env variable value
			// Note: env variable is resolved to its value before comparison
			policy: base.RawPolicy{
				Name:   "env-contains-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Contains: []any{testEnvValue, "moderator", "user"}, // source:role in ["admin", "moderator", "user"]
					},
				},
			},
			source:  User{Name: "user", Role: testEnvValue}, // Role matches env value
			target:  Document{Owner: "user", Type: "public"},
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
TestEnvVariable_Lt tests less-than condition with environment variables.

The test checks that environment variables can be used in Lt conditions with string values.
Note: Environment variables are always strings, so comparisons are done as string comparisons.
*/
func TestEnvVariable_Lt(t *testing.T) {
	ctx := context.Background()

	// Set up test environment variable with string value
	testEnvVar := "TEST_ENV_VAR_LT"
	testEnvValue := "m"
	os.Setenv(testEnvVar, testEnvValue)
	defer os.Unsetenv(testEnvVar)

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
			name: "env less than literal - match",
			// Test checks that env variable (as string) is less than literal value
			// String comparison "m" < "z" is true
			policy: base.RawPolicy{
				Name:   "env-lt-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + testEnvVar: {
						Lt: "z", // env:TEST_ENV_VAR_LT < "z" (string comparison "m" < "z")
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "env less than literal - no match",
			// Test checks that env variable is not less than literal value
			policy: base.RawPolicy{
				Name:   "env-lt-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + testEnvVar: {
						Lt: "a", // env:TEST_ENV_VAR_LT < "a" (string comparison "m" < "a" is false)
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "source field less than env - match",
			// Test checks that source field (string) is less than env variable (string)
			policy: base.RawPolicy{
				Name:   "env-lt-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:name": {
						Lt: "env:" + testEnvVar, // source:name < env:TEST_ENV_VAR_LT (string comparison)
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
			name: "source field less than env - no match",
			// Test checks that source field is not less than env variable
			policy: base.RawPolicy{
				Name:   "env-lt-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:name": {
						Lt: "env:" + testEnvVar, // source:name < env:TEST_ENV_VAR_LT
					},
				},
			},
			source:  User{Name: "zoe", Role: "user"}, // "zoe" >= "m"
			target:  Document{Owner: "user", Type: "public"},
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
TestEnvVariable_Gt tests greater-than condition with environment variables.

The test checks that environment variables can be used in Gt conditions with string values.
Note: Environment variables are always strings, so comparisons are done as string comparisons.
*/
func TestEnvVariable_Gt(t *testing.T) {
	ctx := context.Background()

	// Set up test environment variable with string value
	testEnvVar := "TEST_ENV_VAR_GT"
	testEnvValue := "m"
	os.Setenv(testEnvVar, testEnvValue)
	defer os.Unsetenv(testEnvVar)

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
			name: "env greater than literal - match",
			// Test checks that env variable (as string) is greater than literal value
			// String comparison "m" > "a" is true
			policy: base.RawPolicy{
				Name:   "gt-env-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + testEnvVar: {
						Gt: "a", // env:TEST_ENV_VAR_GT > "a" (string comparison "m" > "a")
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "env greater than literal - no match",
			// Test checks that env variable is not greater than literal value
			policy: base.RawPolicy{
				Name:   "gt-env-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + testEnvVar: {
						Gt: "z", // env:TEST_ENV_VAR_GT > "z" (string comparison "m" > "z" is false)
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "source field greater than env - match",
			// Test checks that source field (string) is greater than env variable (string)
			policy: base.RawPolicy{
				Name:   "gt-env-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:name": {
						Gt: "env:" + testEnvVar, // source:name > env:TEST_ENV_VAR_GT (string comparison)
					},
				},
			},
			source:  User{Name: "zoe", Role: "user"}, // "zoe" > "m"
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "source field greater than env - no match",
			// Test checks that source field is not greater than env variable
			policy: base.RawPolicy{
				Name:   "gt-env-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:name": {
						Gt: "env:" + testEnvVar, // source:name > env:TEST_ENV_VAR_GT
					},
				},
			},
			source:  User{Name: "alice", Role: "user"}, // "alice" <= "m"
			target:  Document{Owner: "user", Type: "public"},
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
TestEnvVariable_NotExists tests error handling when environment variable doesn't exist.

The test checks that appropriate error is returned when environment variable
specified in policy condition doesn't exist.
*/
func TestEnvVariable_NotExists(t *testing.T) {
	ctx := context.Background()

	// Ensure test environment variable doesn't exist
	testEnvVar := "NONEXISTENT_ENV_VAR"
	os.Unsetenv(testEnvVar)

	tests := []struct {
		name    string
		policy  base.RawPolicy
		source  any
		target  any
		action  string
		wantErr bool
	}{
		{
			name: "env variable doesn't exist - left side",
			// Test checks error when env variable on left side doesn't exist
			policy: base.RawPolicy{
				Name:   "env-not-exists-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + testEnvVar: {
						Eq: "some-value", // env:NONEXISTENT_ENV_VAR == "some-value"
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			wantErr: true, // Should return error because env variable doesn't exist
		},
		{
			name: "env variable doesn't exist - right side",
			// Test checks error when env variable on right side doesn't exist
			policy: base.RawPolicy{
				Name:   "env-not-exists-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Eq: "env:" + testEnvVar, // source:role == env:NONEXISTENT_ENV_VAR
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			wantErr: true, // Should return error because env variable doesn't exist
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
TestEnvVariable_MultipleConditions tests combining environment variables with other conditions.

The test checks that environment variables can be combined with other conditions
in a single policy (logical AND).
*/
func TestEnvVariable_MultipleConditions(t *testing.T) {
	ctx := context.Background()

	// Set up test environment variables
	envVar1 := "TEST_ENV_VAR_1"
	envValue1 := "admin"
	envVar2 := "TEST_ENV_VAR_2"
	envValue2 := "production"
	os.Setenv(envVar1, envValue1)
	os.Setenv(envVar2, envValue2)
	defer os.Unsetenv(envVar1)
	defer os.Unsetenv(envVar2)

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
			name: "env and source field - both match",
			// Test checks that both env variable and source field conditions are met
			policy: base.RawPolicy{
				Name:   "env-multiple-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + envVar1: {
						Eq: envValue1, // env:TEST_ENV_VAR_1 == "admin"
					},
					"source:role": {
						Eq: "admin", // source:role == "admin"
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
			name: "env and source field - one fails",
			// Test checks that policy fails when one condition doesn't match
			policy: base.RawPolicy{
				Name:   "env-multiple-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + envVar1: {
						Eq: envValue1, // env:TEST_ENV_VAR_1 == "admin" (true)
					},
					"source:role": {
						Eq: "user", // source:role == "user" (false, role is "admin")
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
			name: "multiple env variables - both match",
			// Test checks that multiple env variables can be used in one policy
			policy: base.RawPolicy{
				Name:   "env-multiple-env-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + envVar1: {
						Eq: envValue1, // env:TEST_ENV_VAR_1 == "admin"
					},
					"env:" + envVar2: {
						Eq: envValue2, // env:TEST_ENV_VAR_2 == "production"
					},
				},
			},
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "env, source, and target - all match",
			// Test checks combining env variable with source and target fields
			policy: base.RawPolicy{
				Name:   "env-complex-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"env:" + envVar1: {
						Eq: envValue1, // env:TEST_ENV_VAR_1 == "admin"
					},
					"source:role": {
						Eq: "admin", // source:role == "admin"
					},
					"target:type": {
						Eq: "public", // target:type == "public"
					},
				},
			},
			source:  User{Name: "user", Role: "admin"},
			target:  Document{Owner: "user", Type: "public"},
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
TestEnvVariable_WithCache tests caching behavior for environment variables.

The test checks that environment variables are properly cached when cache is enabled.
*/
func TestEnvVariable_WithCache(t *testing.T) {
	ctx := context.Background()

	// Set up test environment variable
	testEnvVar := "TEST_ENV_VAR_CACHE"
	testEnvValue := "cached-value"
	os.Setenv(testEnvVar, testEnvValue)
	defer os.Unsetenv(testEnvVar)

	policy := base.RawPolicy{
		Name:   "env-cache-test",
		Action: "user:read",
		Effect: base.Effect_ALLOW,
		Conditions: map[string]base.Condition{
			"env:" + testEnvVar: {
				Eq: testEnvValue, // env:TEST_ENV_VAR_CACHE == "cached-value"
			},
			"source:role": {
				Eq: "env:" + testEnvVar, // source:role == env:TEST_ENV_VAR_CACHE
			},
		},
	}

	// Test with cache enabled
	casher := implemented.NewDefaultCasher()
	config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
	engine, err := guardian.NewGuardianFromPolices(casher, []base.RawPolicy{policy}, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	source := User{Name: "user", Role: testEnvValue}
	target := Document{Owner: "user", Type: "public"}

	// First evaluation - should populate cache
	got1, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Fatalf("unexpected error on first evaluation: %v", err)
	}
	if !got1 {
		t.Errorf("first evaluation: Evaluate() = %v, want true", got1)
	}

	// Second evaluation - should use cache
	got2, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Fatalf("unexpected error on second evaluation: %v", err)
	}
	if !got2 {
		t.Errorf("second evaluation: Evaluate() = %v, want true", got2)
	}
}

