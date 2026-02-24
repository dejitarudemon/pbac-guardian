/*
Package tests contains tests for ConditionFuncsConfig.

Tests check the functionality of ConditionFuncsConfig.Select method
for retrieving condition functions by name.
*/
package tests

import (
	"testing"

	"github.com/dejitarudemon/pbac-guardian/internal/base"
	"github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

/*
TestConditionFuncsConfigSelect tests the Select method of ConditionFuncsConfig.

The test checks that Select correctly returns condition functions by their names
and returns nil for unknown function names.
*/
func TestConditionFuncsConfigSelect(t *testing.T) {
	config := implemented.DefaultConditionsFuncs

	tests := []struct {
		name     string
		key      string
		wantNil  bool
		wantFunc string
	}{
		{
			name:     "select Contains",
			key:      "Contains",
			wantNil:  false,
			wantFunc: "Contains",
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
TestConditionFuncsConfigAllFunctionsPresent tests that all required functions are present in DefaultConditionsFuncs.

The test checks that DefaultConditionsFuncs contains all four condition functions:
Contains, Eq, Neq, and Lt.
*/
func TestConditionFuncsConfigAllFunctionsPresent(t *testing.T) {
	config := implemented.DefaultConditionsFuncs

	if config.Contains == nil {
		t.Error("DefaultConditionsFuncs.Contains is nil")
	}
	if config.Eq == nil {
		t.Error("DefaultConditionsFuncs.Eq is nil")
	}
	if config.Neq == nil {
		t.Error("DefaultConditionsFuncs.Neq is nil")
	}
	if config.Lt == nil {
		t.Error("DefaultConditionsFuncs.Lt is nil")
	}
}

/*
TestConditionFuncsConfigCustomConfig tests creating a custom ConditionFuncsConfig.

The test checks that custom configuration can be created with custom functions
and that Select method works correctly with custom configuration.
*/
func TestConditionFuncsConfigCustomConfig(t *testing.T) {
	// Create custom config with only Eq function
	customConfig := base.ConditionFuncsConfig{
		Eq: implemented.DefaultConditionsFuncs.Eq,
		// Other functions are nil
	}

	// Test that Select returns correct function for present key
	eqFunc := customConfig.Select("Eq")
	if eqFunc == nil {
		t.Error("Select(\"Eq\") returned nil, expected function")
	}

	// Test that Select returns nil for absent functions
	containsFunc := customConfig.Select("Contains")
	if containsFunc != nil {
		t.Error("Select(\"Contains\") returned function, expected nil")
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

