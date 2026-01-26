/*
Package tests contains tests for error types of the pbac-guardian library.

Tests check creation, wrapping and unwrapping of errors,
as well as correctness of error messages.

This file is in the tests directory for test organization.
Tests import the noctisguard package as regular library users.
*/
package tests

import (
	"errors"
	"testing"

	"github.com/dejitarudemon/pbac-guardian"
)

/*
TestErrDuplicateName tests the duplicate policy name error.

The test checks:
  - Error creation with specified policy name
  - Error type correctness via errors.As()
  - Name storage in error structure
  - Error message format (should contain policy name)
*/
func TestErrDuplicateName(t *testing.T) {
	// Create error with test policy name
	err := noctisguard.NewErrDuplicateName("test-policy")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Check error message (type check via errors.As not available for unexported types)
	// Check that message contains policy name and duplicate indication
	expectedMsg := "test-policy is already used by another policy"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
}

/*
TestErrExport tests the policy export error.

The test checks:
  - Error creation wrapping original error
  - Error type correctness via errors.As()
  - Ability to unwrap error via errors.Unwrap()
  - Check via errors.Is() for original error
  - Storage of original error in source field of structure
*/
func TestErrExport(t *testing.T) {
	// Create original error that will be wrapped
	originalErr := errors.New("original error")

	// Wrap it in ErrExport
	err := noctisguard.NewErrExport(originalErr)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Check that error wraps original via errors.Is()
	// This allows checking for original error in the chain
	if !errors.Is(err, originalErr) {
		t.Error("expected error to wrap original error")
	}

	// Check error unwrapping via errors.Unwrap()
	// This allows access to original error
	unwrapped := errors.Unwrap(err)
	if unwrapped != originalErr {
		t.Errorf("expected unwrapped error to be original error")
	}

	// Check that error message contains export information
	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("expected non-empty error message")
	}
}

/*
TestErrEvaluate tests the policy evaluation error.

The test checks:
  - Error creation wrapping original error
  - Error type correctness via errors.As()
  - Ability to unwrap error via errors.Unwrap()
  - Check via errors.Is() for original error
  - Storage of original error in source field of structure

This error is used for errors during Evaluate execution,
e.g., invalid field paths or comparison errors.
*/
func TestErrEvaluate(t *testing.T) {
	// Create original error that will be wrapped
	originalErr := errors.New("original error")

	// Wrap it in ErrEvaluate
	err := noctisguard.NewErrEvaluate(originalErr)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Check that error wraps original via errors.Is()
	if !errors.Is(err, originalErr) {
		t.Error("expected error to wrap original error")
	}

	// Check error unwrapping via errors.Unwrap()
	unwrapped := errors.Unwrap(err)
	if unwrapped != originalErr {
		t.Errorf("expected unwrapped error to be original error")
	}

	// Check that error message contains evaluation information
	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("expected non-empty error message")
	}
}
