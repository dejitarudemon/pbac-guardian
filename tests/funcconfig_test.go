/*
Package tests contains tests for custom condition functions configuration.

Tests check the functionality of providing custom condition functions
through the config parameter in NewGuardianFromPolices and NewGuardianFromFile.
The config parameter is of type base.Config and contains ConditionsMap.
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

The test checks that custom condition functions can be provided through config
parameter, and they are used instead of default functions during policy evaluation.
*/
func TestCustomConditionFuncsConfig(t *testing.T) {
	// Create custom condition function that always returns true for Eq
	customEqFunc := func(ctx context.Context, left, right any) (bool, error) {
		// Custom function: always return true for equality check
		return true, nil
	}

	// Create custom ConditionsMap with custom Eq function
	customConditionsMap := &base.ConditionsMap{
		In:  implemented.DefaultConditionsMap.In,
		Eq:  customEqFunc, // Custom function
		Neq: implemented.DefaultConditionsMap.Neq,
		Lt:  implemented.DefaultConditionsMap.Lt,
		Gt:       implemented.DefaultConditionsMap.Gt,
		Le:       implemented.DefaultConditionsMap.Le,
		Ge:       implemented.DefaultConditionsMap.Ge,
	}

	// Create config with custom ConditionsMap
	config := base.Config{
		ConditionsMap:        customConditionsMap,
		CashDisableThreShold: 3,
	}

	policies := []base.RawPolicy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"}, // This should always pass with custom function
			},
		},
	}

	// Create engine with custom config
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
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
TestNilFuncConfigUsesDefaults tests that nil ConditionsMap in config uses default functions.

The test checks that when nil is passed in config.ConditionsMap, the engine uses
default condition functions from implemented.DefaultConditionsMap.
*/
func TestNilFuncConfigUsesDefaults(t *testing.T) {
	policies := []base.RawPolicy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	// Create config with nil ConditionsMap (should use defaults)
	config := base.Config{
		ConditionsMap:        nil, // use defaults
		CashDisableThreShold: 3,
	}

	// Create engine with config containing nil ConditionsMap
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
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
TestCustomInFunc tests custom In condition function.

The test checks that custom In function can be provided and used
instead of the default implementation.
*/
func TestCustomInFunc(t *testing.T) {
	// Custom In function that always returns true
	customInFunc := func(ctx context.Context, left, right any) (bool, error) {
		return true, nil
	}

	customConditionsMap := &base.ConditionsMap{
		In:  customInFunc,
		Eq:  implemented.DefaultConditionsMap.Eq,
		Neq: implemented.DefaultConditionsMap.Neq,
		Lt:  implemented.DefaultConditionsMap.Lt,
		Gt:       implemented.DefaultConditionsMap.Gt,
		Le:       implemented.DefaultConditionsMap.Le,
		Ge:       implemented.DefaultConditionsMap.Ge,
	}

	config := base.Config{
		ConditionsMap:        customConditionsMap,
		CashDisableThreShold: 3,
	}

	policies := []base.RawPolicy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {
					In: []any{"admin", "moderator"}, // Should always pass with custom function
				},
			},
		},
	}

	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "guest"} // Role is "guest", not in list
	target := Document{Owner: "user", Type: "public"}

	// With custom function, In should always return true
	allowed, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true with custom In function, got false")
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

	customConditionsMap := &base.ConditionsMap{
		In:  implemented.DefaultConditionsMap.In,
		Eq:  implemented.DefaultConditionsMap.Eq,
		Neq: implemented.DefaultConditionsMap.Neq,
		Lt:  customLtFunc,
		Gt:       implemented.DefaultConditionsMap.Gt,
		Le:       implemented.DefaultConditionsMap.Le,
		Ge:       implemented.DefaultConditionsMap.Ge,
	}

	config := base.Config{
		ConditionsMap:        customConditionsMap,
		CashDisableThreShold: 3,
	}

	policies := []base.RawPolicy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:age": {Lt: 18}, // Should always pass with custom function
			},
		},
	}

	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
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
TestNewGuardianFromFileWithFuncConfig tests engine creation from file with custom config.

The test checks that NewGuardianFromFile accepts config parameter
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

	// Create custom condition function
	customEqFunc := func(ctx context.Context, left, right any) (bool, error) {
		return true, nil
	}

	// Create custom ConditionsMap
	customConditionsMap := &base.ConditionsMap{
		In:  implemented.DefaultConditionsMap.In,
		Eq:  customEqFunc,
		Neq: implemented.DefaultConditionsMap.Neq,
		Lt:  implemented.DefaultConditionsMap.Lt,
		Gt:       implemented.DefaultConditionsMap.Gt,
		Le:       implemented.DefaultConditionsMap.Le,
		Ge:       implemented.DefaultConditionsMap.Ge,
	}

	// Create config with custom ConditionsMap
	config := base.Config{
		ConditionsMap:        customConditionsMap,
		CashDisableThreShold: 3,
	}

	// Create engine from file with custom config
	engine, err := guardian.NewGuardianFromFile(nil, tmpFile.Name(), config)
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
