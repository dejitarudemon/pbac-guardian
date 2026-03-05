/*
Package base provides basic types and functions for working with policies,
conditions, effects, and entities in the access control system.

The package contains definitions of policies, comparison conditions, effects (allow/deny),
entities (source/target) and interfaces for custom comparison.
*/
package base

import (
	"errors"
	"fmt"
)

/*
This file contains custom errors for the engine base
*/

var (
	// ErrNotComparableStruct represents an error that occurs when the left argument
	// is a structure but does not implement the Compare() method from base.Comparable interface.
	// This error is returned by LtConditionFunc when trying to compare structures
	// that don't implement custom comparison logic.
	ErrNotComparableStruct = errors.New("left argument is a struct, but it doesn't implement Compare() method")

	// ErrCancelled represents an error that occurs when an operation is cancelled through context.Context.
	// Used to interrupt long-running condition checking operations when context is cancelled.
	// This error is returned by condition functions and policy evaluation methods when
	// context.Context signals cancellation (via ctx.Done()).
	ErrCancelled = errors.New("cancelled by context")

	// ErrNilContext represents an error that occurs when a nil context is passed to a function
	// that requires a valid context.Context. This error is returned by policy evaluation methods
	// and condition functions when ctx parameter is nil instead of a valid context instance.
	// All methods that accept context.Context should validate it is not nil before use.
	ErrNilContext = errors.New("context is nil")
)

/*
ErrInvalidType represents an error that occurs when expected and received types do not match.

The error is used when a function expects a certain type (e.g., a structure
or slice) but receives a different type.
*/
type ErrInvalidType struct {
	expected any
	got      any
}

/*
NewErrInvalidType creates a new ErrInvalidType error.

Parameters:
  - expected - expected type (can be a string or type)
  - got - received type (can be a string, type, or nil)

Returns:
  - error - created ErrInvalidType error
*/
func NewErrInvalidType(expected, got any) error {
	return ErrInvalidType{expected: expected, got: got}
}

func (e ErrInvalidType) Error() string {
	return fmt.Sprintf("expected %v, but got %v", e.expected, e.got)
}

/*
ErrUncomparable represents an error that occurs when trying to
compare two values that cannot be compared with each other.

The error occurs when value types are incompatible for comparison (e.g.,
different primitive types in Lt operation) or the Compare() method returned false.
*/
type ErrUncomparable struct {
	left  any
	right any
}

/*
NewErrUncomparable creates a new ErrUncomparable error.

Parameters:
  - left - left value to compare
  - right - right value to compare

Returns:
  - error - created ErrUncomparable error
*/
func NewErrUncomparable(left, right any) error {
	return ErrUncomparable{left: left, right: right}
}

func (e ErrUncomparable) Error() string {
	return fmt.Sprintf("can't compare %v with %v", e.left, e.right)
}

/*
ErrInvalidPath represents an error that occurs when working
with an invalid path to a structure field.

The error is used when parsing paths in format "entity:field1:field2..."
or when searching for fields in structures. Contains the path and problem details.
*/
type ErrInvalidPath struct {
	path    any
	details string
}

/*
NewErrInvalidPath creates a new ErrInvalidPath error.

Parameters:
  - path - invalid path that caused the error
  - details - error details (problem description)

Returns:
  - error - created ErrInvalidPath error
*/
func NewErrInvalidPath(path any, details string) error {
	return ErrInvalidPath{path: path, details: details}
}

func (e ErrInvalidPath) Error() string {
	return fmt.Sprintf("invalid path %v: %v", e.path, e.details)
}

/*
ErrInexpectedBehavior represents an internal library error
that occurs with unexpected behavior in the code.

The error indicates a problem in the library logic (e.g., missing condition
function in CONDITION_TO_FUNC) and usually should not occur with proper usage.
*/
type ErrInexpectedBehavior struct {
	source  string
	details string
}

/*
NewErrInexpectedBehavior creates a new ErrInexpectedBehavior error.

Parameters:
  - source - error source (function name or location, e.g., "Policy.Evaluate()")
  - details - unexpected behavior details (problem description)

Returns:
  - error - created ErrInexpectedBehavior error
*/
func NewErrInexpectedBehavior(source, details string) error {
	return ErrInexpectedBehavior{source: source, details: details}
}

func (e ErrInexpectedBehavior) Error() string {
	return fmt.Sprintf("unexpected behavior in %v : %v", e.source, e.details)
}
