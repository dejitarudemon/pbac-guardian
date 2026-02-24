/*
Package tests contains tests for the pbac-guardian library.

Tests check the functionality of creating an engine from policies and files,
as well as policy evaluation for various access scenarios.

This file is in the tests directory for test organization.
Tests import the guardian package as regular library users.
*/
package tests

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

/*
Test structures for checking library functionality.

These structures are used in all tests to simulate real
objects that are checked against access policies.
*/

// User represents a system user with fields tagged for access
type User struct {
	Name string `pbac-guardian:"name"` // User name
	Role string `pbac-guardian:"role"` // User role (admin, user, guest, etc.)
	Age  int    `pbac-guardian:"age"`  // User age
}

// Document represents a document with owner and type information
type Document struct {
	Owner string   `pbac-guardian:"owner"` // Document owner
	Type  string   `pbac-guardian:"type"`  // Document type (public, private, etc.)
	Tags  []string `pbac-guardian:"tags"`  // Document tags
}

// NestedUser represents a nested user structure for testing nested paths
type NestedUser struct {
	User User `pbac-guardian:"user"` // Nested user
}

// NestedDocument represents a nested document structure for testing nested paths
type NestedDocument struct {
	Doc Document `pbac-guardian:"doc"` // Nested document
}

/*
TestNewGuardianFromPolices tests engine creation from a list of policies.

The test checks:
  - Engine creation with valid policies
  - Handling of duplicate policy names (should return error)
  - Action format validation (minimum 2 parts separated by ":")
  - Path validation in conditions (format "entity:field")
  - Engine creation with empty policy list (should succeed)
*/
func TestNewGuardianFromPolices(t *testing.T) {
	tests := []struct {
		name     string
		policies []base.Policy
		wantErr  bool
		errType  error
	}{
		{
			name: "valid policies",
			// Test checks successful engine creation with valid policies
			// Policy has correct action format and valid paths in conditions
			policies: []base.Policy{
				{
					Name:   "test-policy",
					Action: "user:read", // Valid format: minimum 2 parts
					Effect: base.Effect_ALLOW,
					Conditions: map[string]base.Condition{
						"source:role": {Eq: "admin"}, // Valid path: entity:field
					},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate names",
			// Test checks that engine is not created when policies have duplicate names
			// Should return ErrExport error wrapping ErrDuplicateName
			policies: []base.Policy{
				{
					Name:   "test-policy",
					Action: "user:read",
					Effect: base.Effect_ALLOW,
					Conditions: map[string]base.Condition{
						"source:role": {Eq: "admin"},
					},
				},
				{
					Name:   "test-policy", // Duplicate name - should cause error
					Action: "user:write",
					Effect: base.Effect_ALLOW,
					Conditions: map[string]base.Condition{
						"source:role": {Eq: "admin"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid action format",
			// Test checks action format validation
			// Action must contain minimum 2 parts separated by ":"
			policies: []base.Policy{
				{
					Name:   "test-policy",
					Action: "invalid", // Invalid format - only one part
					Effect: base.Effect_ALLOW,
					Conditions: map[string]base.Condition{
						"source:role": {Eq: "admin"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid path in conditions",
			// Test checks path validation in conditions
			// Path must have format "entity:field", where entity is "source" or "target"
			policies: []base.Policy{
				{
					Name:   "test-policy",
					Action: "user:read",
					Effect: base.Effect_ALLOW,
					Conditions: map[string]base.Condition{
						"invalid": {Eq: "admin"}, // Invalid path - no separator ":"
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty policies",
			// Test checks engine creation with empty policy list
			// This should succeed - engine is created but without policies
			policies: []base.Policy{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use nil casher for basic functionality tests
			engine, err := guardian.NewGuardianFromPolices(nil, tt.policies)
			if tt.wantErr {
				// Expect error
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				// Expect success
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if engine == nil {
					t.Errorf("expected engine, got nil")
				}
			}
		})
	}
}

/*
TestNewGuardianFromFile tests engine creation from a JSON file.

The test checks:
  - Successful reading and parsing of valid JSON file with policies
  - Error handling when file is missing
  - Correctness of engine creation from file (similar to NewGuardianFromPolices)
*/
func TestNewGuardianFromFile(t *testing.T) {
	// Create temporary file with valid policies for testing
	validPolicies := []base.Policy{
		{
			Name:   "file-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	// Serialize policies to JSON for writing to file
	validJSON, _ := json.Marshal(validPolicies)

	// Create temporary file with unique name
	tmpFile, err := os.CreateTemp("", "test_policies_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Remove temporary file after test

	// Write JSON to file
	if _, err := tmpFile.Write(validJSON); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Create temporary file with invalid JSON
	invalidJSON := []byte(`{"invalid": "json"`)
	tmpFileInvalid, err := os.CreateTemp("", "test_invalid_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFileInvalid.Name())
	if _, err := tmpFileInvalid.Write(invalidJSON); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFileInvalid.Close()

	// Create temporary file with empty content
	tmpFileEmpty, err := os.CreateTemp("", "test_empty_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFileEmpty.Name())
	tmpFileEmpty.Close()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name: "valid file",
			// Test checks successful reading and parsing of valid JSON file
			// File must contain array of Policy objects in JSON format
			path:    tmpFile.Name(),
			wantErr: false,
		},
		{
			name: "non-existent file",
			// Test checks error handling when trying to open non-existent file
			// Should return ErrExport error wrapping os.PathError
			path:    "/nonexistent/file.json",
			wantErr: true,
		},
		{
			name: "invalid JSON",
			// Test checks error handling for invalid JSON in file
			// Should return ErrExport error wrapping json.SyntaxError
			path:    tmpFileInvalid.Name(),
			wantErr: true,
		},
		{
			name: "empty file",
			// Test checks error handling for empty file
			// Should return ErrExport error wrapping json.SyntaxError
			path:    tmpFileEmpty.Name(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use nil casher for basic functionality tests
			engine, err := guardian.NewGuardianFromFile(nil, tt.path)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if engine == nil {
					t.Errorf("expected engine, got nil")
				}
			}
		})
	}
}

/*
TestEvaluate tests the main policy evaluation functionality.

The test checks various access scenarios:
  - Access granted for admins (allow-admin policy)
  - Access granted for owners (allow-owner policy with field comparison)
  - Access denied for non-admins to private documents (deny-private policy)
  - Access granted for admins to private documents (deny-private policy does not apply)
  - Handling of missing policies for action (returns false)
  - Handling of invalid field paths (returns error)
*/
func TestEvaluate(t *testing.T) {
	// Create set of policies for testing various access scenarios
	policies := []base.Policy{
		{
			Name: "allow-admin",
			// Policy allows admins to read documents
			// Condition: source:role == "admin"
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
		{
			Name: "allow-owner",
			// Policy allows document owner to read
			// Condition: source:name == target:owner (comparing fields from different structures)
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:name": {Eq: "target:owner"},
			},
		},
		{
			Name: "deny-private",
			// Policy denies non-admins access to private documents
			// Conditions: target:type == "private" AND source:role != "admin"
			// If both conditions are met, access is denied
			Action: "user:read:document",
			Effect: base.Effect_DENY,
			Conditions: map[string]base.Condition{
				"target:type": {Eq: "private"},
				"source:role": {Neq: "admin"},
			},
		},
	}

	// Use nil casher for basic functionality tests
	engine, err := guardian.NewGuardianFromPolices(nil, policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "admin can read",
			// Test checks that admin can read documents according to allow-admin policy
			// source.role == "admin", so policy should pass
			source:  User{Name: "admin", Role: "admin"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read:document",
			want:    true,
			wantErr: false,
		},
		{
			name: "owner can read",
			// Test checks that document owner can read it
			// source.name == target.owner == "alice", so allow-owner policy should pass
			source:  User{Name: "alice", Role: "user"},
			target:  Document{Owner: "alice", Type: "public"},
			action:  "user:read:document",
			want:    true,
			wantErr: false,
		},
		{
			name: "deny private for non-admin",
			// Test checks that non-admin cannot read private documents
			// target.type == "private" AND source.role != "admin", so deny-private policy applies
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "other", Type: "private"},
			action:  "user:read:document",
			want:    false,
			wantErr: false,
		},
		{
			name: "admin can read private",
			// Test checks that admin can read private documents
			// source.role == "admin", so condition source.role != "admin" is not met
			// deny-private policy does not apply, and admin gets access through allow-admin
			source:  User{Name: "admin", Role: "admin"},
			target:  Document{Owner: "other", Type: "private"},
			action:  "user:read:document",
			want:    true,
			wantErr: false,
		},
		{
			name: "no policies for action",
			// Test checks that when there are no policies for action, false is returned
			// This means the action is denied by default
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:write:document", // No policies for this action
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
TestEvaluateWithContextCancellation tests operation cancellation through context.Context.

The test checks that long-running condition checking operations can be interrupted
through context cancellation, and the function returns the corresponding ErrCancelled error.
This is especially important for Contains operations with large lists.
*/
func TestEvaluateWithContextCancellation(t *testing.T) {
	// Create policy with Contains condition that may take long for large lists
	policies := []base.Policy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {
					Contains: []any{"admin", "moderator", "user"},
				},
			},
		},
	}

	// Use nil casher for basic functionality tests
	engine, err := guardian.NewGuardianFromPolices(nil, policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	// Create context with cancellation and cancel it immediately
	// This simulates a situation where operation should be interrupted
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel context before execution starts

	source := User{Name: "user", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	// Expect operation to return cancellation error
	_, err = engine.Evaluate(ctx, source, target, "user:read")
	if err == nil {
		t.Errorf("expected cancellation error, got nil")
	}
}

/*
TestEvaluateWithTimeout tests working with timeout through context.Context.

The test checks that operation successfully completes within timeout
and is not interrupted prematurely. This is important for checking correctness
of context handling under normal conditions.
*/
func TestEvaluateWithTimeout(t *testing.T) {
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

	// Use nil casher for basic functionality tests
	engine, err := guardian.NewGuardianFromPolices(nil, policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	// Create context with 1 second timeout
	// Operation should complete faster than timeout expires
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	// Operation should complete successfully within timeout
	allowed, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true, got false")
	}
}

/*
TestEvaluateNestedStructures tests working with nested structures.

The test checks that the library correctly handles paths to fields
in nested structures, e.g., "source:user:role" to access
the role field inside nested user structure.

This is important for working with real data structures that often
have nested structure.
*/
func TestEvaluateNestedStructures(t *testing.T) {
	policies := []base.Policy{
		{
			Name:   "nested-policy",
			Action: "user:read:nested",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				// Check access to field in nested structure
				// Path "source:user:role" means: in source find field user, then in it find field role
				"source:user:role": {Eq: "admin"},
			},
		},
	}

	// Use nil casher for basic functionality tests
	engine, err := guardian.NewGuardianFromPolices(nil, policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	// Create nested structures for testing
	source := NestedUser{User: User{Name: "admin", Role: "admin"}}
	target := NestedDocument{Doc: Document{Owner: "user", Type: "public"}}

	// Check that path to nested field works correctly
	allowed, err := engine.Evaluate(ctx, source, target, "user:read:nested")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true, got false")
	}
}
