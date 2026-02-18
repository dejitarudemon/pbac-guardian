/*
Package base provides basic types and functions for working with policies,
conditions, effects, and entities in the access control system.

The package contains definitions of policies, comparison conditions, effects (allow/deny),
entities (source/target) and interfaces for custom comparison.
*/
package base

import (
	"context"
	"fmt"
	"reflect"
)

/*
This file contains various conditions for policies
*/

const SUPPORTED_LT_PRIMITIVES = "int|uint|float|string"

/*
Condition represents a set of rules for comparing values in policies.

A condition can contain one or more fields. If multiple fields are specified,
they are checked independently of each other. Values can be either literals
or paths to structure fields (e.g., "source:name" or "target:role").

Fields:
  - Contains - checks if left value is in right list (right must be a slice)
  - Eq - checks equality of values (left == right)
  - Neq - checks inequality of values (left != right)
  - Lt - checks if left value is less than right (left < right)

Example usage:

	type User struct {
		Name string   `noctis-guard:"name"`
		Role string   `noctis-guard:"role"`
		Age  int      `noctis-guard:"age"`
		Tags []string `noctis-guard:"tags"`
	}

	// Equality check with literal
	condition1 := Condition{
		Eq: "admin", // source:role == "admin"
	}

	// Equality check of fields from two structures
	condition2 := Condition{
		Eq: "target:owner", // source:name == target:owner
	}

	// Inequality check
	condition3 := Condition{
		Neq: "guest", // source:role != "guest"
	}

	// Check if value is less
	condition4 := Condition{
		Lt: 18, // source:age < 18
	}

	// Check if value is in list
	condition5 := Condition{
		Contains: []any{"admin", "moderator"}, // source:role in ["admin", "moderator"]
	}

	// Combining conditions (all must be met)
	condition6 := Condition{
		Eq:       "user",           // source:role == "user"
		Contains: []any{"read"},    // "read" in source:tags
		Lt:       100,             // source:age < 100
	}
*/
type Condition struct {
	Contains []any `json:"contains,omitempty"`
	Eq       any   `json:"eq,omitempty"`
	Neq      any   `json:"neq,omitempty"`
	Lt       any   `json:"lt,omitempty"`
}

/*
conditionFunc represents a function for checking a condition between two values.

Functions of this type are used to check conditions in policies. Argument order
matters for Contains and Lt operations (left and right are not interchangeable).

Functions must support cancellation through context.Context and return ErrCancelled
when context is cancelled.

Parameters:
  - ctx - context for operation cancellation and timeout control
  - left - left value to compare (what is being checked)
  - right - right value to compare (what is being compared against)

Returns:
  - bool - comparison result (true if condition is met, false otherwise)
  - err - comparison execution error (nil if comparison successful, ErrCancelled on context cancellation)
*/
type conditionFunc func(ctx context.Context, left, right any) (bool, error)

/*
containsConditionFunc checks if value left is in list right.

The function uses reflect.DeepEqual to compare elements, allowing it to
work with any data types. Supports cancellation through context.Context.

Parameters:
  - ctx - context for operation cancellation and timeout control
  - left - value to search for in the list
  - right - list (slice) or pointer to list in which to search

Returns:
  - bool - true if left is found in right, false otherwise
  - err - execution error if right is not a list or operation is cancelled

Possible errors:
  - ErrCancelled - operation was cancelled through context.Context
  - ErrInvalidType - right is not a slice or pointer to slice (may be nil)
*/
func containsConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	if right == nil {
		return false, NewErrInvalidType(reflect.Slice.String(), nil)
	}

	slice := reflect.ValueOf(right)

	if slice.Kind() == reflect.Pointer {
		if slice.IsNil() {
			return false, NewErrInvalidType(reflect.Slice.String(), nil)
		}

		slice = slice.Elem()
	}

	if slice.Kind() != reflect.Slice {
		return false, NewErrInvalidType(reflect.Slice.String(), slice.Kind().String())
	}

	for i := range slice.Len() {
		select {
		case <-ctx.Done():
			return false, ErrCancelled
		default:
			if reflect.DeepEqual(left, slice.Index(i).Interface()) {
				return true, nil
			}
		}
	}

	return false, nil
}

/*
eqConditionFunc checks equality of two values.

If one of the values implements the Comparable interface, the Compare() method
is used for comparison. If Compare() returns false (comparison impossible), or
neither value implements Comparable, reflect.DeepEqual is used.

Parameters:
  - ctx - context for operation cancellation and timeout control (not used, but required for compatibility)
  - left - left value to compare
  - right - right value to compare

Returns:
  - bool - true if values are equal, false otherwise
  - err - execution error (always nil, function does not return errors)
*/
func eqConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	// nil значения обрабатываются через reflect.DeepEqual
	// reflect.DeepEqual(nil, nil) == true, reflect.DeepEqual(nil, non-nil) == false

	if l, ok := left.(Comparable); ok {
		if result, acceptable := l.Compare(right); acceptable {
			return result == 0, nil
		}
	}
	if r, ok := right.(Comparable); ok {
		if result, acceptable := r.Compare(left); acceptable {
			return result == 0, nil
		}
	}
	return reflect.DeepEqual(left, right), nil
}

/*
neqConditionFunc checks inequality of two values.

The function is the inverse of eqConditionFunc: returns !eqConditionFunc(ctx, left, right).
Uses the same comparison logic through Comparable or DeepEqual.

Parameters:
  - ctx - context for operation cancellation and timeout control (passed to eqConditionFunc)
  - left - left value to compare
  - right - right value to compare

Returns:
  - bool - true if values are not equal, false if equal
  - err - execution error (always nil, function does not return errors)
*/
func neqConditionFunc(ctx context.Context, left, right any) (bool, error) {
	ok, err := eqConditionFunc(ctx, left, right)
	return !ok, err
}

/*
ltConditionFunc checks if left is less than right.

If left is a structure, it must implement the Comparable interface.
For primitive types (int, uint, float, string), comparison is done through reflect.
Types must match for correct comparison.

Parameters:
  - ctx - context for operation cancellation and timeout control (passed to ltPrimitives)
  - left - left value to compare
  - right - right value to compare

Returns:
  - bool - true if left < right, false otherwise
  - err - execution error if comparison is impossible

Possible errors:
  - ErrNotComparableStruct - left is a structure but does not implement Comparable interface
  - ErrUncomparable - cannot compare left and right (incompatible types or Compare returned false)
  - ErrInvalidType - left or right is neither a structure nor a supported primitive
    (int, uint, float, string and their variants)
*/
func ltConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	if left == nil {
		return false, NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, nil)
	}
	if right == nil {
		return false, NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, nil)
	}

	if reflect.TypeOf(left).Kind() != reflect.Struct {
		return ltPrimitives(left, right)
	}

	l, ok := left.(Comparable)
	if !ok {
		return false, ErrNotComparableStruct
	}

	result, ok := l.Compare(right)

	if !ok {
		return false, NewErrUncomparable(left, right)
	}

	if result < 0 {
		return true, nil
	}

	return false, nil
}

/*
ltPrimitives compares primitive types left and right based on reflection.

The function supports comparison of the following types:
  - int, int8, int16, int32, int64
  - uint, uint8, uint16, uint32, uint64
  - float32, float64
  - string

Types left and right must match. If a pointer is passed, it is dereferenced.

Parameters:
  - left - left value to compare (primitive type)
  - right - right value to compare (primitive type)

Returns:
  - bool - true if left < right, false otherwise
  - err - execution error if comparison is impossible

Possible errors:
  - ErrUncomparable - types left and right do not match or one of them is nil
  - ErrInvalidType - left or right is not a supported primitive
    (int, uint, float, string and their variants)
*/
func ltPrimitives(left, right any) (bool, error) {
	v1 := reflect.ValueOf(left)
	v2 := reflect.ValueOf(right)

	if v1.Kind() == reflect.Pointer {
		if v1.IsNil() {
			return false, NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, nil)
		}

		v1 = v1.Elem()
	}

	if v2.Kind() == reflect.Pointer {
		if v2.IsNil() {
			return false, NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, nil)
		}

		v2 = v2.Elem()
	}

	if v1.Type() != v2.Type() {
		return false, NewErrUncomparable(left, right)
	}

	switch v1.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int8, reflect.Int64:
		return v1.Int() < v2.Int(), nil
	case reflect.Float32, reflect.Float64:
		return v1.Float() < v2.Float(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v1.Uint() < v2.Uint(), nil
	case reflect.String:
		return v1.String() < v2.String(), nil
	}

	return false, NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, v1.Kind().String())
}
