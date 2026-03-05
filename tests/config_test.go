/*
Package tests contains tests for ConditionsMap.

Tests check the functionality of ConditionsMap.Select method
for retrieving condition functions by name.
*/
package tests

import (
	"testing"

	"github.com/dejitarudemon/pbac-guardian/internal/base"
	"github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

/*
TestConditionFuncsConfigSelect tests the Select method of ConditionsMap.

The test checks that Select correctly returns condition functions by their names
and returns nil for unknown function names.

Note: This test uses ConditionsMap directly, not Config structure.
*/
func TestConditionFuncsConfigSelect(t *testing.T) {
	config := implemented.DefaultConditionsMap

	tests := []struct {
		name     string
		key      string
		wantNil  bool
		wantFunc string
	}{
		{
			name:     "select In",
			key:      "In",
			wantNil:  false,
			wantFunc: "In",
		},
		{
			name:     "select Eq",
			key:      "Eq",
			wantNil:  false,
			wantFunc: "Eq",
		},
		{
			name:     "select Neq",
			key:      "Neq",
			wantNil:  false,
			wantFunc: "Neq",
		},
		{
			name:     "select Lt",
			key:      "Lt",
			wantNil:  false,
			wantFunc: "Lt",
		},
		{
			name:    "select unknown function",
			key:     "Unknown",
			wantNil: true,
		},
		{
			name:    "select empty string",
			key:     "",
			wantNil: true,
		},
		{
			name:    "select lowercase key",
			key:     "eq",
			wantNil: true, // Keys are case-sensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.Select(tt.key)
			if tt.wantNil {
				if got != nil {
					t.Errorf("Select(%q) = %v, want nil", tt.key, got)
				}
			} else {
				if got == nil {
					t.Errorf("Select(%q) = nil, want function", tt.key)
				}
			}
		})
	}
}

/*
TestConditionFuncsConfigAllFunctionsPresent tests that all required functions are present in DefaultConditionsMap.

The test checks that DefaultConditionsMap contains all four condition functions:
In, Eq, Neq, and Lt.
*/
func TestConditionFuncsConfigAllFunctionsPresent(t *testing.T) {
	config := implemented.DefaultConditionsMap

	if config.In == nil {
		t.Error("DefaultConditionsMap.In is nil")
	}
	if config.Eq == nil {
		t.Error("DefaultConditionsMap.Eq is nil")
	}
	if config.Neq == nil {
		t.Error("DefaultConditionsMap.Neq is nil")
	}
	if config.Lt == nil {
		t.Error("DefaultConditionsMap.Lt is nil")
	}
}

/*
TestConditionFuncsConfigCustomConfig tests creating a custom ConditionsMap.

The test checks that custom configuration can be created with custom functions
and that Select method works correctly with custom configuration.
*/
func TestConditionFuncsConfigCustomConfig(t *testing.T) {
	// Create custom config with only Eq function
	customConfig := base.ConditionsMap{
		Eq: implemented.DefaultConditionsMap.Eq,
		// Other functions are nil
	}

	// Test that Select returns correct function for present key
	eqFunc := customConfig.Select("Eq")
	if eqFunc == nil {
		t.Error("Select(\"Eq\") returned nil, expected function")
	}

	// Test that Select returns nil for absent functions
	inFunc := customConfig.Select("In")
	if inFunc != nil {
		t.Error("Select(\"In\") returned function, expected nil")
	}

	neqFunc := customConfig.Select("Neq")
	if neqFunc != nil {
		t.Error("Select(\"Neq\") returned function, expected nil")
	}

	ltFunc := customConfig.Select("Lt")
	if ltFunc != nil {
		t.Error("Select(\"Lt\") returned function, expected nil")
	}
}
