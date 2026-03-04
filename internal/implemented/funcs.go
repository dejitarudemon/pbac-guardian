/*
Package implemented provides default implementations of condition functions
and cache mechanisms for the policy evaluation engine.

The package contains:
  - Default condition functions (Contains, Eq, Neq, Lt, Gt) for policy evaluation
  - DefaultCasher implementation for L1 caching

The default condition functions support:
  - Primitive types (int, uint, float, string and their variants)
  - time.Time values (for Eq, Neq, Lt, Gt conditions)
  - Custom types implementing base.Comparable interface
*/
package implemented

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

// SUPPORTED_LT_PRIMITIVES is a constant describing supported primitive types for Lt comparison.
const SUPPORTED_LT_GT_PRIMITIVES = "int|uint|float|string"

/*
DefaultConditionsMap provides default implementations of condition functions.

This configuration is used by default when nil is passed as funcConfig parameter
to NewGuardianFromPolices or NewGuardianFromFile. It contains standard implementations
of all condition types: Contains, Eq, Neq, Lt, and Gt.

The default functions support:
  - Primitive types (int, uint, float, string and their variants)
  - time.Time values (for Eq, Neq, Lt, Gt conditions)
  - Custom types implementing base.Comparable interface
  - Context cancellation for long-running operations
*/
var (
	DefaultConditionsMap = base.ConditionsMap{
		Contains: ContainsConditionFunc,
		Eq:       EqConditionFunc,
		Neq:      NeqConditionFunc,
		Lt:       LtConditionFunc,
		Gt:       GtConditionFunc,
	}
)

/*
Comparison type constants for primitive value comparison.

These constants are used internally by comparePrimitives and compare functions
to specify the type of comparison operation to perform.

Values:
  - LE - less than or equal (<=)
  - LT - less than (<)
  - GE - greater than or equal (>=)
  - GT - greater than (>)
*/
const (
	LE = iota
	LT
	GE
	GT
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
			ok, err := EqConditionFunc(ctx, left, slice.Index(i).Interface())
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
	}

	return false, nil
}

/*
EqConditionFunc checks equality of two values.

The function uses the following comparison logic in order:
 1. If one of the values implements the Comparable interface, the Compare() method is used
 2. If both values are time.Time, they are compared using time.Time methods
 3. Otherwise, reflect.DeepEqual is used

Parameters:
  - ctx - context for operation cancellation and timeout control (not used, but required for compatibility)
  - left - left value to compare
  - right - right value to compare

Returns:
  - bool - true if values are equal, false otherwise
  - err - execution error (always nil, function does not return errors)

Supported types:
  - Primitive types (int, uint, float, string and their variants)
  - time.Time values
  - Custom types implementing base.Comparable interface
  - Any types supported by reflect.DeepEqual
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
	if result, ok := timesCompare(left, right); ok {
		return result == 0, nil
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

The function uses the following comparison logic:
 1. For primitive types (int, uint, float, string), comparison is done through reflection
 2. For time.Time values, comparison is done using time.Time methods (Before/After)
 3. For structures, the Comparable interface is used if implemented
 4. Types must match for correct comparison

Parameters:
  - ctx - context for operation cancellation and timeout control (passed to ltPrimitives)
  - left - left value to compare
  - right - right value to compare

Returns:
  - bool - true if left < right, false otherwise
  - err - execution error if comparison is impossible

Possible errors:
  - ErrNotComparableStruct - left is a structure but does not implement Comparable interface and is not time.Time
  - ErrUncomparable - cannot compare left and right (incompatible types or Compare returned false)
  - ErrInvalidType - left or right is neither a structure, time.Time, nor a supported primitive
    (int, uint, float, string and their variants)

Supported types:
  - Primitive types (int, uint, float, string and their variants)
  - time.Time values
  - Custom types implementing base.Comparable interface
*/
func LtConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	if left == nil {
		return false, base.NewErrInvalidType(SUPPORTED_LT_GT_PRIMITIVES, nil)
	}
	if right == nil {
		return false, base.NewErrInvalidType(SUPPORTED_LT_GT_PRIMITIVES, nil)
	}

	if reflect.TypeOf(left).Kind() != reflect.Struct {
		return comparePrimitives(left, right, LT)
	}

	l, ok := left.(base.Comparable)
	if !ok {
		if asTimeResult, ok := timesCompare(left, right); ok {
			return asTimeResult < 0, nil
		}
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
GtConditionFunc checks if left is greater than right.

The function uses the following comparison logic:
 1. For primitive types (int, uint, float, string), comparison is done through reflection
 2. For time.Time values, comparison is done using time.Time methods (Before/After)
 3. For structures, the Comparable interface is used if implemented
 4. Types must match for correct comparison

Parameters:
  - ctx - context for operation cancellation and timeout control (passed to gtPrimitives)
  - left - left value to compare
  - right - right value to compare

Returns:
  - bool - true if left > right, false otherwise
  - err - execution error if comparison is impossible

Possible errors:
  - ErrNotComparableStruct - left is a structure but does not implement Comparable interface and is not time.Time
  - ErrUncomparable - cannot compare left and right (incompatible types or Compare returned false)
  - ErrInvalidType - left or right is neither a structure, time.Time, nor a supported primitive
    (int, uint, float, string and their variants)

Supported types:
  - Primitive types (int, uint, float, string and their variants)
  - time.Time values
  - Custom types implementing base.Comparable interface
*/
func GtConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	if left == nil {
		return false, base.NewErrInvalidType(SUPPORTED_LT_GT_PRIMITIVES, nil)
	}
	if right == nil {
		return false, base.NewErrInvalidType(SUPPORTED_LT_GT_PRIMITIVES, nil)
	}

	if reflect.TypeOf(left).Kind() != reflect.Struct {
		return comparePrimitives(left, right, GT)
	}

	l, ok := left.(base.Comparable)
	if !ok {
		if asTimeResult, ok := timesCompare(left, right); ok {
			return asTimeResult > 0, nil
		}
		return false, base.ErrNotComparableStruct
	}

	result, ok := l.Compare(right)

	if !ok {
		return false, base.NewErrUncomparable(left, right)
	}

	if result > 0 {
		return true, nil
	}

	return false, nil
}

/*
LeConditionFunc checks if left is less than or equal to right.

The function uses the following comparison logic:
 1. For primitive types (int, uint, float, string), comparison is done through reflection
 2. For time.Time values, comparison is done using time.Time methods (Before/After)
 3. For structures, the Comparable interface is used if implemented
 4. Types must match for correct comparison

Parameters:
  - ctx - context for operation cancellation and timeout control (passed to ltPrimitives)
  - left - left value to compare
  - right - right value to compare

Returns:
  - bool - true if left <= right, false otherwise
  - err - execution error if comparison is impossible

Possible errors:
  - ErrNotComparableStruct - left is a structure but does not implement Comparable interface and is not time.Time
  - ErrUncomparable - cannot compare left and right (incompatible types or Compare returned false)
  - ErrInvalidType - left or right is neither a structure, time.Time, nor a supported primitive
    (int, uint, float, string and their variants)

Supported types:
  - Primitive types (int, uint, float, string and their variants)
  - time.Time values
  - Custom types implementing base.Comparable interface
*/
func LeConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	if left == nil {
		return false, base.NewErrInvalidType(SUPPORTED_LT_GT_PRIMITIVES, nil)
	}
	if right == nil {
		return false, base.NewErrInvalidType(SUPPORTED_LT_GT_PRIMITIVES, nil)
	}

	if reflect.TypeOf(left).Kind() != reflect.Struct {
		return comparePrimitives(left, right, LE)
	}

	l, ok := left.(base.Comparable)
	if !ok {
		if asTimeResult, ok := timesCompare(left, right); ok {
			return asTimeResult <= 0, nil
		}
		return false, base.ErrNotComparableStruct
	}

	result, ok := l.Compare(right)

	if !ok {
		return false, base.NewErrUncomparable(left, right)
	}

	if result <= 0 {
		return true, nil
	}

	return false, nil
}

/*
GeConditionFunc checks if left is greater than or equal to right.

The function uses the following comparison logic:
 1. For primitive types (int, uint, float, string), comparison is done through reflection
 2. For time.Time values, comparison is done using time.Time methods (Before/After)
 3. For structures, the Comparable interface is used if implemented
 4. Types must match for correct comparison

Parameters:
  - ctx - context for operation cancellation and timeout control (passed to gtPrimitives)
  - left - left value to compare
  - right - right value to compare

Returns:
  - bool - true if left > right, false otherwise
  - err - execution error if comparison is impossible

Possible errors:
  - ErrNotComparableStruct - left is a structure but does not implement Comparable interface and is not time.Time
  - ErrUncomparable - cannot compare left and right (incompatible types or Compare returned false)
  - ErrInvalidType - left or right is neither a structure, time.Time, nor a supported primitive
    (int, uint, float, string and their variants)

Supported types:
  - Primitive types (int, uint, float, string and their variants)
  - time.Time values
  - Custom types implementing base.Comparable interface
*/
func GeConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	if left == nil {
		return false, base.NewErrInvalidType(SUPPORTED_LT_GT_PRIMITIVES, nil)
	}
	if right == nil {
		return false, base.NewErrInvalidType(SUPPORTED_LT_GT_PRIMITIVES, nil)
	}

	if reflect.TypeOf(left).Kind() != reflect.Struct {
		return comparePrimitives(left, right, GE)
	}

	l, ok := left.(base.Comparable)
	if !ok {
		if asTimeResult, ok := timesCompare(left, right); ok {
			return asTimeResult >= 0, nil
		}
		return false, base.ErrNotComparableStruct
	}

	result, ok := l.Compare(right)

	if !ok {
		return false, base.NewErrUncomparable(left, right)
	}

	if result >= 0 {
		return true, nil
	}

	return false, nil
}

/*
timesCompare compares two time.Time values.

The function is used internally by condition functions (Eq, Neq, Lt, Gt) to support
time.Time comparisons without requiring the Comparable interface implementation.

Parameters:
  - left - first time.Time value to compare
  - right - second time.Time value to compare

Returns:
  - int - comparison result:
  - < 0 - left is before right
  - = 0 - left equals right
  - > 0 - left is after right
  - bool - true if both values are time.Time and comparison was successful, false otherwise

The function returns (0, false) if either left or right is not a time.Time value,
allowing condition functions to fall back to other comparison methods.
*/
func timesCompare(left, right any) (int, bool) {
	l, ok := left.(time.Time)
	if !ok {
		return 0, false
	}

	r, ok := right.(time.Time)
	if !ok {
		return 0, false
	}

	if l.Before(r) {
		return -1, true
	}

	if l.After(r) {
		return 1, true
	}

	return 0, true
}

/*
comparePrimitives compares two primitive values using reflection.

The function supports comparison of numeric types (int, uint, float) and strings.
It handles pointer types by dereferencing them before comparison. The comparison
type is specified by the compareType parameter (LE, LT, GE, GT).

Parameters:
  - left - first value to compare (must be a primitive type or pointer to primitive)
  - right - second value to compare (must be a primitive type or pointer to primitive)
  - compareType - type of comparison to perform (LE, LT, GE, or GT)

Returns:
  - bool - comparison result based on compareType
  - error - execution error if types don't match or are not supported primitives

Possible errors:
  - ErrUncomparable - left and right types don't match
  - ErrInvalidType - left or right is not a supported primitive type or is nil pointer

Supported types:
  - int, int8, int16, int32, int64
  - uint, uint8, uint16, uint32, uint64
  - float32, float64
  - string
  - Pointers to any of the above types
*/
func comparePrimitives(left, right any, compareType uint) (bool, error) {
	v1 := reflect.ValueOf(left)
	v2 := reflect.ValueOf(right)

	if v1.Kind() == reflect.Pointer {
		if v1.IsNil() {
			return false, base.NewErrInvalidType(SUPPORTED_LT_GT_PRIMITIVES, nil)
		}

		v1 = v1.Elem()
	}

	if v2.Kind() == reflect.Pointer {
		if v2.IsNil() {
			return false, base.NewErrInvalidType(SUPPORTED_LT_GT_PRIMITIVES, nil)
		}

		v2 = v2.Elem()
	}

	if v1.Type() != v2.Type() {
		return false, base.NewErrUncomparable(left, right)
	}

	switch v1.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int8, reflect.Int64:
		return compare(v1.Int(), v2.Int(), compareType), nil
	case reflect.Float32, reflect.Float64:
		return compare(v1.Float(), v2.Float(), compareType), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return compare(v1.Uint(), v2.Uint(), compareType), nil
	case reflect.String:
		return compare(v1.String(), v2.String(), compareType), nil
	}

	return false, base.NewErrInvalidType(SUPPORTED_LT_GT_PRIMITIVES, v1.Kind().String())
}

/*
ordered is a type constraint for ordered types that support comparison operators.

The constraint includes all numeric types (signed and unsigned integers, floats)
and strings. The tilde (~) prefix allows the constraint to match both the base
type and any custom types defined with that base type (e.g., type UserID int).

This constraint is used by the compare function to ensure type safety while
allowing comparison operations (<, <=, >, >=) on generic type parameters.

Supported types:
  - int, int8, int16, int32, int64
  - uint, uint8, uint16, uint32, uint64
  - float32, float64
  - string
  - Any custom types with these base types (e.g., type MyInt int)
*/
type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string
}

/*
compare performs comparison of two ordered values based on the comparison type.

The function is a generic helper used internally by comparePrimitives to perform
the actual comparison operation. It supports all comparison types: less than or
equal (LE), less than (LT), greater than or equal (GE), and greater than (GT).

Parameters:
  - left - first value to compare
  - right - second value to compare
  - compareType - type of comparison to perform (LE, LT, GE, or GT)

Returns:
  - bool - comparison result:
  - For LE: true if left <= right
  - For LT: true if left < right
  - For GE: true if left >= right
  - For GT: true if left > right
  - false for any other compareType value

The function uses the ordered constraint to ensure that only types supporting
comparison operators can be used, providing compile-time type safety.
*/
func compare[T ordered](left, right T, compareType uint) bool {
	switch compareType {
	case LE:
		return left <= right
	case LT:
		return left < right
	case GE:
		return left >= right
	case GT:
		return left > right
	}
	return false
}
