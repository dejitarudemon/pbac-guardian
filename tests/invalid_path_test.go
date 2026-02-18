/*
Package tests contains tests for handling invalid field paths.

Tests check that the library correctly handles errors for invalid paths
to fields in policy conditions.
*/
package tests

import (
	"context"
	"testing"

	noctisguard "github.com/dejitarudemon/noctis-guard"
	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
TestEvaluateInvalidFieldPath tests error handling for invalid field path.

The test creates a policy with non-existent field in condition and checks
that Evaluate returns an error when trying to get value of non-existent field.
*/
func TestEvaluateInvalidFieldPath(t *testing.T) {
	// Create policy with invalid field path
	// Field "nonexistent" does not exist in User structure
	policies := []base.Policy{
		{
			Name:   "invalid-path-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:nonexistent": {Eq: "value"}, // Non-existent field
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user"}
	target := Document{Owner: "user", Type: "public"}

	// Expect error, as field "nonexistent" does not exist in User structure
	_, err = engine.Evaluate(ctx, source, target, "user:read")
	if err == nil {
		t.Errorf("expected error for invalid field path, got nil")
	}
}

/*
TestEvaluateInvalidNestedPath tests error handling for invalid nested path.

The test creates a policy with invalid nested path and checks
that Evaluate returns an error when trying to get value of non-existent nested field.
*/
func TestEvaluateInvalidNestedPath(t *testing.T) {
	// Create policy with invalid nested path
	// Path "source:user:nonexistent" is invalid, as field "nonexistent" does not exist
	policies := []base.Policy{
		{
			Name:   "invalid-nested-path-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:user:nonexistent": {Eq: "value"}, // Non-existent nested field
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user"}
	target := Document{Owner: "user", Type: "public"}

	// Expect error, as nested field "nonexistent" does not exist
	_, err = engine.Evaluate(ctx, source, target, "user:read")
	if err == nil {
		t.Errorf("expected error for invalid nested field path, got nil")
	}
}

/*
TestEvaluateInvalidTargetPath tests error handling for invalid field path in target.

The test creates a policy with invalid field path in target structure and checks
that Evaluate returns an error when trying to get value of non-existent field.
*/
func TestEvaluateInvalidTargetPath(t *testing.T) {
	// Create policy with invalid field path in target
	// Field "nonexistent" does not exist in Document structure
	policies := []base.Policy{
		{
			Name:   "invalid-target-path-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"target:nonexistent": {Eq: "value"}, // Non-existent field in target
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user"}
	target := Document{Owner: "user", Type: "public"}

	// Expect error, as field "nonexistent" does not exist in Document structure
	_, err = engine.Evaluate(ctx, source, target, "user:read")
	if err == nil {
		t.Errorf("expected error for invalid target field path, got nil")
	}
}

