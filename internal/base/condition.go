/*
Package base provides basic types and functions for working with policies,
conditions, effects, and entities in the access control system.

The package contains definitions of policies, comparison conditions, effects (allow/deny),
entities (source/target) and interfaces for custom comparison.
*/
package base

import (
	"context"
)

/*
Condition represents a set of rules for comparing values in policies.

A condition can contain one or more fields. If multiple fields are specified,
they are checked independently of each other. Values can be either literals
or paths to structure fields, environment variables, or time values
(e.g., "source:name", "target:role", "env:VAR_NAME", "time:now").

Fields:
  - Contains - checks if left value is in right list (right must be a slice)
  - Eq - checks equality of values (left == right)
  - Neq - checks inequality of values (left != right)
  - Lt - checks if left value is less than right (left < right)
  - Gt - checks if left value is greater than right (left > right)
  - Le - checks if left value is less than or equal to right (left <= right)
  - Ge - checks if left value is greater than or equal to right (left >= right)

Example usage:

	type User struct {
		Name string   `pbac-guardian:"name"`
		Role string   `pbac-guardian:"role"`
		Age  int      `pbac-guardian:"age"`
		Tags []string `pbac-guardian:"tags"`
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

	// Check if value is greater
	condition5 := Condition{
		Gt: 18, // source:age > 18
	}

	// Check if value is in list
	condition6 := Condition{
		Contains: []any{"admin", "moderator"}, // source:role in ["admin", "moderator"]
	}

	// Combining conditions (all must be met)
	condition7 := Condition{
		Eq:       "user",           // source:role == "user"
		Contains: []any{"read"},    // "read" in source:tags
		Lt:       100,             // source:age < 100
		Gt:       18,              // source:age > 18
	}
*/
type Condition struct {
	Contains []any `json:"contains,omitempty"`
	Eq       any   `json:"eq,omitempty"`
	Neq      any   `json:"neq,omitempty"`
	Lt       any   `json:"lt,omitempty"`
	Gt       any   `json:"gt,omitempty"`
	Le       any   `json:"le,omitempty"`
	Ge       any   `json:"ge,omitempty"`
}

/*
ConditionFunc represents a function for checking a condition between two values.

Functions of this type are used to check conditions in policies. Argument order
matters for Contains, Lt, Gt, Le, and Ge operations (left and right are not interchangeable).

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
type ConditionFunc func(ctx context.Context, left, right any) (bool, error)
