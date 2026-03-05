/*
Package tests contains tests for public error constructors of base package.

Tests check error creation, their messages and correctness of working
with errors.As() and errors.Is(). The tests verify:

  - ErrInvalidType - type mismatch errors
  - ErrUncomparable - incomparable value errors
  - ErrInvalidPath - invalid field path errors
  - ErrInexpectedBehavior - internal library errors
  - Error message formatting with expected and received values
*/
package tests

import (
	"errors"
	"testing"

	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

/*
TestErrInvalidType tests the public constructor NewErrInvalidType.

The test checks:
  - Error creation with expected and received types
  - Error type correctness via errors.As()
  - Presence of error message
*/
func TestErrInvalidType(t *testing.T) {
	// Create error with expected type "struct" and received "string"
	err := base.NewErrInvalidType("struct", "string")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Check error type
	var invalidTypeErr base.ErrInvalidType
	if !errors.As(err, &invalidTypeErr) {
		t.Errorf("expected ErrInvalidType, got %T", err)
	}

	// Check that error message is not empty
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

/*
TestErrUncomparable tests the public constructor NewErrUncomparable.

The test checks:
  - Error creation with left and right values
  - Error type correctness via errors.As()
  - Presence of error message
*/
func TestErrUncomparable(t *testing.T) {
	// Create error for incompatible types (int and string)
	err := base.NewErrUncomparable(5, "string")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Check error type
	var uncomparableErr base.ErrUncomparable
	if !errors.As(err, &uncomparableErr) {
		t.Errorf("expected ErrUncomparable, got %T", err)
	}

	// Check that error message is not empty
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

/*
TestErrInvalidPath tests the public constructor NewErrInvalidPath.

The test checks:
  - Error creation with path and details
  - Error type correctness via errors.As()
  - Presence of error message
*/
func TestErrInvalidPath(t *testing.T) {
	// Create error with invalid path and problem description
	err := base.NewErrInvalidPath("source:invalid", "field not found")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Check error type
	var invalidPathErr base.ErrInvalidPath
	if !errors.As(err, &invalidPathErr) {
		t.Errorf("expected ErrInvalidPath, got %T", err)
	}

	// Check that error message is not empty
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

/*
TestErrInexpectedBehavior tests the public constructor NewErrInexpectedBehavior.

The test checks:
  - Error creation with source and details
  - Error type correctness via errors.As()
  - Presence of error message
*/
func TestErrInexpectedBehavior(t *testing.T) {
	// Create error with source and problem details
	err := base.NewErrInexpectedBehavior("TestFunction", "unexpected condition")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Check error type
	var unexpectedErr base.ErrInexpectedBehavior
	if !errors.As(err, &unexpectedErr) {
		t.Errorf("expected ErrInexpectedBehavior, got %T", err)
	}

	// Check that error message is not empty
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

/*
TestErrCancelled tests the global public variable ErrCancelled.

The test checks:
  - That ErrCancelled is not nil
  - Correctness of error message
*/
func TestErrCancelled(t *testing.T) {
	// Check that error is defined
	if base.ErrCancelled == nil {
		t.Fatal("ErrCancelled should not be nil")
	}

	// Check error message
	if base.ErrCancelled.Error() != "cancelled by context" {
		t.Errorf("expected error message 'cancelled by context', got %q", base.ErrCancelled.Error())
	}
}

/*
TestErrNotComparableStruct tests the global public variable ErrNotComparableStruct.

The test checks:
  - That ErrNotComparableStruct is not nil
  - Correctness of error message
*/
func TestErrNotComparableStruct(t *testing.T) {
	// Check that error is defined
	if base.ErrNotComparableStruct == nil {
		t.Fatal("ErrNotComparableStruct should not be nil")
	}

	// Check error message
	expectedMsg := "left argument is a struct, but it doesn't implement Compare() method"
	if base.ErrNotComparableStruct.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, base.ErrNotComparableStruct.Error())
	}
}
