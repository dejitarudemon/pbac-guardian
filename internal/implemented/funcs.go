/*
Package implemented provides default implementations of condition functions
and cache mechanisms for the policy evaluation engine.

The package contains:
  - Default condition functions (Contains, Eq, Neq, Lt) for policy evaluation
  - DefaultCasher implementation for L1 caching
*/
package implemented

import (
	"context"
	"fmt"
	"reflect"

	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

// SUPPORTED_LT_PRIMITIVES is a constant describing supported primitive types for Lt comparison.
const SUPPORTED_LT_PRIMITIVES = "int|uint|float|string"

/*
DefaultConditionsMap provides default implementations of condition functions.

This configuration is used by default when nil is passed as funcConfig parameter
to NewGuardianFromPolices or NewGuardianFromFile. It contains standard implementations
of all condition types: Contains, Eq, Neq, and Lt.

The default functions support:
  - Custom types implementing base.Comparable interface
  - Primitive types (int, uint, float, string and their variants)
  - Context cancellation for long-running operations
*/
var (
	DefaultConditionsMap = base.ConditionsMap{
		Contains: ContainsConditionFunc,
		Eq:       EqConditionFunc,
		Neq:      NeqConditionFunc,
		Lt:       LtConditionFunc,
	}
)

/*
ContainsConditionFunc checks if value left is in list right.

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
func ContainsConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	if right == nil {
		return false, base.NewErrInvalidType(reflect.Slice.String(), nil)
	}

	slice := reflect.ValueOf(right)

	if slice.Kind() == reflect.Pointer {
		if slice.IsNil() {
			return false, base.NewErrInvalidType(reflect.Slice.String(), nil)
		}

		slice = slice.Elem()
	}

	if slice.Kind() != reflect.Slice {
		return false, base.NewErrInvalidType(reflect.Slice.String(), slice.Kind().String())
	}

	for i := range slice.Len() {
		select {
		case <-ctx.Done():
			return false, base.ErrCancelled
		default:
			if reflect.DeepEqual(left, slice.Index(i).Interface()) {
				return true, nil
			}
		}
	}

	return false, nil
}

/*
EqConditionFunc checks equality of two values.

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
func EqConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	// nil значения обрабатываются через reflect.DeepEqual
	// reflect.DeepEqual(nil, nil) == true, reflect.DeepEqual(nil, non-nil) == false

	if l, ok := left.(base.Comparable); ok {
		if result, acceptable := l.Compare(right); acceptable {
			return result == 0, nil
		}
	}
	if r, ok := right.(base.Comparable); ok {
		if result, acceptable := r.Compare(left); acceptable {
			return result == 0, nil
		}
	}
	return reflect.DeepEqual(left, right), nil
}

/*
NeqConditionFunc checks inequality of two values.

The function is the inverse of EqConditionFunc: returns !EqConditionFunc(ctx, left, right).
Uses the same comparison logic through Comparable or DeepEqual.

Parameters:
  - ctx - context for operation cancellation and timeout control (passed to EqConditionFunc)
  - left - left value to compare
  - right - right value to compare

Returns:
  - bool - true if values are not equal, false if equal
  - err - execution error (always nil, function does not return errors)
*/
func NeqConditionFunc(ctx context.Context, left, right any) (bool, error) {
	ok, err := EqConditionFunc(ctx, left, right)
	return !ok, err
}

/*
LtConditionFunc checks if left is less than right.

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
func LtConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	if left == nil {
		return false, base.NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, nil)
	}
	if right == nil {
		return false, base.NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, nil)
	}

	if reflect.TypeOf(left).Kind() != reflect.Struct {
		return ltPrimitives(left, right)
	}

	l, ok := left.(base.Comparable)
	if !ok {
		return false, base.ErrNotComparableStruct
	}

	result, ok := l.Compare(right)

	if !ok {
		return false, base.NewErrUncomparable(left, right)
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
			return false, base.NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, nil)
		}

		v1 = v1.Elem()
	}

	if v2.Kind() == reflect.Pointer {
		if v2.IsNil() {
			return false, base.NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, nil)
		}

		v2 = v2.Elem()
	}

	if v1.Type() != v2.Type() {
		return false, base.NewErrUncomparable(left, right)
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

	return false, base.NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, v1.Kind().String())
}
