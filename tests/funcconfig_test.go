/*
Package tests contains tests for custom condition functions configuration.

Tests check the functionality of providing custom condition functions
through the funcConfig parameter in NewGuardianFromPolices and NewGuardianFromFile.
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
TestCustomConditionFuncsConfig tests engine creation with custom condition functions.

The test checks that custom condition functions can be provided through funcConfig
parameter, and they are used instead of default functions during policy evaluation.
*/
func TestCustomConditionFuncsConfig(t *testing.T) {
	// Create custom condition function that always returns true for Eq
	customEqFunc := func(ctx context.Context, left, right any) (bool, error) {
		// Custom function: always return true for equality check
		return true, nil
	}

	// Create custom funcConfig with custom Eq function
	customConfig := &base.ConditionsMap{
		Contains: implemented.DefaultConditionsMap.Contains,
		Eq:       customEqFunc, // Custom function
		Neq:      implemented.DefaultConditionsMap.Neq,
		Lt:       implemented.DefaultConditionsMap.Lt,
	}

	policies := []base.Policy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"}, // This should always pass with custom function
			},
		},
	}

	// Create engine with custom funcConfig
	engine, err := guardian.NewGuardianFromPolices(nil, policies, customConfig)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user"} // Role is "user", not "admin"
	target := Document{Owner: "user", Type: "public"}

	// With custom function, Eq should always return true, so policy should pass
	allowed, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true with custom Eq function, got false")
	}
}

/*
TestNilFuncConfigUsesDefaults tests that nil funcConfig uses default functions.

The test checks that when nil is passed as funcConfig, the engine uses
default condition functions from implemented.DefaultConditionsMap.
*/
func TestNilFuncConfigUsesDefaults(t *testing.T) {
	policies := []base.Policy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	// Create engine with nil funcConfig (should use defaults)
	engine, err := guardian.NewGuardianFromPolices(nil, policies, nil)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()

	// Test with admin role - should pass
	source1 := User{Name: "admin", Role: "admin"}
	target1 := Document{Owner: "user", Type: "public"}
	allowed1, err := engine.Evaluate(ctx, source1, target1, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed1 {
		t.Errorf("expected allowed=true for admin, got false")
	}

	// Test with user role - should fail (default function checks actual equality)
	source2 := User{Name: "user", Role: "user"}
	target2 := Document{Owner: "user", Type: "public"}
	allowed2, err := engine.Evaluate(ctx, source2, target2, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if allowed2 {
		t.Errorf("expected allowed=false for user (not admin), got true")
	}
}

/*
TestCustomContainsFunc tests custom Contains condition function.

The test checks that custom Contains function can be provided and used
instead of the default implementation.
*/
func TestCustomContainsFunc(t *testing.T) {
	// Custom Contains function that always returns true
	customContainsFunc := func(ctx context.Context, left, right any) (bool, error) {
		return true, nil
	}

	customConfig := &base.ConditionsMap{
		Contains: customContainsFunc,
		Eq:       implemented.DefaultConditionsMap.Eq,
		Neq:      implemented.DefaultConditionsMap.Neq,
		Lt:       implemented.DefaultConditionsMap.Lt,
	}

	policies := []base.Policy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {
					Contains: []any{"admin", "moderator"}, // Should always pass with custom function
				},
			},
		},
	}

	engine, err := guardian.NewGuardianFromPolices(nil, policies, customConfig)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "guest"} // Role is "guest", not in list
	target := Document{Owner: "user", Type: "public"}

	// With custom function, Contains should always return true
	allowed, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true with custom Contains function, got false")
	}
}

/*
TestCustomLtFunc tests custom Lt condition function.

The test checks that custom Lt function can be provided and used
instead of the default implementation.
*/
func TestCustomLtFunc(t *testing.T) {
	// Custom Lt function that always returns true
	customLtFunc := func(ctx context.Context, left, right any) (bool, error) {
		return true, nil
	}

	customConfig := &base.ConditionsMap{
		Contains: implemented.DefaultConditionsMap.Contains,
		Eq:       implemented.DefaultConditionsMap.Eq,
		Neq:      implemented.DefaultConditionsMap.Neq,
		Lt:       customLtFunc,
	}

	policies := []base.Policy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:age": {Lt: 18}, // Should always pass with custom function
			},
		},
	}

	engine, err := guardian.NewGuardianFromPolices(nil, policies, customConfig)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user", Age: 25} // Age is 25, not < 18
	target := Document{Owner: "user", Type: "public"}

	// With custom function, Lt should always return true
	allowed, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true with custom Lt function, got false")
	}
}

/*
TestNewGuardianFromFileWithFuncConfig tests engine creation from file with custom funcConfig.

The test checks that NewGuardianFromFile accepts funcConfig parameter
and uses it when creating the engine.
*/
func TestNewGuardianFromFileWithFuncConfig(t *testing.T) {
	// Create temporary JSON file with policies
	content := `[
		{
			"name": "test-policy",
			"action": "user:read",
			"effect": "allow",
			"conditions": {
				"source:role": {"eq": "admin"}
			}
		}
	]`

	tmpFile, err := os.CreateTemp("", "test-policies-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	// Create custom funcConfig
	customEqFunc := func(ctx context.Context, left, right any) (bool, error) {
		return true, nil
	}

	customConfig := &base.ConditionsMap{
		Contains: implemented.DefaultConditionsMap.Contains,
		Eq:       customEqFunc,
		Neq:      implemented.DefaultConditionsMap.Neq,
		Lt:       implemented.DefaultConditionsMap.Lt,
	}

	// Create engine from file with custom funcConfig
	engine, err := guardian.NewGuardianFromFile(nil, tmpFile.Name(), customConfig)
	if err != nil {
		t.Fatalf("failed to create engine from file: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user"} // Role is "user", not "admin"
	target := Document{Owner: "user", Type: "public"}

	// With custom function, should always pass
	allowed, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true with custom funcConfig, got false")
	}
}
