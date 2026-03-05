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
	"strings"

	"github.com/dejitarudemon/pbac-guardian/internal/cashing"
)

const (
	// MIN_ACTION_PARTS is the minimum number of parts in a valid action.
	// An action must have at least 2 parts: entity and action name.
	// Example: "user:read" (2 parts) is valid, "read" (1 part) is invalid.
	MIN_ACTION_PARTS = 2
)

/*
Policy represents a single access policy with a set of conditions.

A policy defines rules for checking source and target structures for a specific
action. Conditions are checked through structure fields tagged with "pbac-guardian".

Policy is an internal structure with private fields. To create a Policy instance,
use NewPolicy function with a RawPolicy. Policy instances are created automatically
by Guardian engine during initialization from RawPolicy structures.

The Policy structure contains:
  - action - action in format "entity:action:extra1:extra2..." (e.g., "user:read:profile")
  - effect - policy effect: Effect_ALLOW (allow) or Effect_DENY (deny)
  - conditions - map of conditions. Key - path to field in format "source:field" or "target:field",
    value - condition to check (In, Eq, Neq, Lt, Gt, Le, Ge, Any, All)
  - conditionsMap - configuration for condition functions (In, Eq, Neq, Lt, Gt, Le, Ge)
  - cash - L1 cache instance for storing field values (can be nil to disable caching)
  - cashTree - cache tree for tracking field access counts and disabling cache for rarely accessed fields

Public methods:
  - Evaluate - applies the policy to source and target structures
  - IsValid - checks the validity of the policy
  - Effect - returns the policy effect

Example usage:

	raw := RawPolicy{
		Name:   "admin-read",
		Action: "user:read:document",
		Effect: Effect_ALLOW,
		Conditions: map[string]Condition{
			"source:role": {
				Eq: "admin",
			},
		},
	}

	policy, err := NewPolicy(raw, &DefaultConditionsMap, nil, nil)
	if err != nil {
		// handle error
	}

	allowed, err := policy.Evaluate(ctx, source, target, "user:read:document", sessionID)
*/
type Policy struct {
	name          string
	action        string
	effect        Effect
	conditions    map[string]Condition
	conditionsMap *ConditionsMap
	cash          cashing.Casher
	cashTree      *cashing.CashTree
	loader        *Loader
}

/*
NewPolicy creates a new Policy instance from a RawPolicy.

The function validates the raw policy, initializes internal fields, and prepares
the policy for evaluation. It performs the following steps:
 1. Validates that conditions map is not nil and all condition functions are set
 2. Creates Policy with internal fields initialized
 3. Validates policy structure (action format, paths in conditions)
 4. Records field paths to cash tree for cache optimization

Parameters:
  - raw - raw policy structure with public fields (Name, Action, Effect, Conditions)
  - conditions - configuration for condition functions (In, Eq, Neq, Lt, Gt, Le, Ge). Must not be nil.
  - cash - L1 cache instance for storing field values (can be nil to disable caching)
  - cashTree - cache tree for tracking field access counts (can be nil, cache tracking will be skipped)

Returns:
  - *Policy - created and validated policy instance, ready for evaluation
  - error - creation error if:
  - conditions is nil or any condition function is nil
  - policy validation fails (invalid action format or paths in conditions)

Example usage:

	raw := RawPolicy{
		Name:   "admin-read",
		Action: "user:read:document",
		Effect: Effect_ALLOW,
		Conditions: map[string]Condition{
			"source:role": {Eq: "admin"},
		},
	}

	policy, err := NewPolicy(raw, &DefaultConditionsMap, casher, cashTree)
	if err != nil {
		// handle error
	}
*/
func NewPolicy(raw RawPolicy, conditions *ConditionsMap, cash cashing.Casher, cashTree *cashing.CashTree) (*Policy, error) {
	if conditions == nil {
		return nil, NewErrInvalidType("ConditionsMap", nil)
	}

	if conditions.In == nil {
		return nil, NewErrInvalidType("ConditionsMap.In", nil)
	}
	if conditions.Eq == nil {
		return nil, NewErrInvalidType("ConditionsMap.Eq", nil)
	}
	if conditions.Lt == nil {
		return nil, NewErrInvalidType("ConditionsMap.Lt", nil)
	}
	if conditions.Neq == nil {
		return nil, NewErrInvalidType("ConditionsMap.Neq", nil)
	}
	if conditions.Le == nil {
		return nil, NewErrInvalidType("ConditionsMap.Le", nil)
	}
	if conditions.Ge == nil {
		return nil, NewErrInvalidType("ConditionsMap.Ge", nil)
	}

	p := Policy{
		name:          raw.Action,
		action:        raw.Action,
		effect:        raw.Effect,
		conditions:    raw.Conditions,
		conditionsMap: conditions,
		cash:          cash,
		cashTree:      cashTree,
		loader:        &Loader{},
	}

	if err := p.IsValid(); err != nil {
		return nil, err
	}

	p.addInformationToCashTree()

	return &p, nil
}

/*
addInformationToCashTree records all field paths used in policy conditions to the cache tree.

The function scans all conditions in the policy and adds field paths (both left and right sides)
to the CashTree for tracking access counts. This allows the cache tree to determine which
fields should be cached based on access frequency.

If cashTree is nil, the function returns immediately without doing anything.

This method is called automatically during policy creation in NewPolicy function.
*/
func (p *Policy) addInformationToCashTree() {
	if p.cashTree == nil {
		return
	}
	for left, condition := range p.conditions {
		if !strings.HasPrefix(left, string(Entity_TIME)) {
			p.cashTree.Add(p.action, left)
		}

		c := reflect.ValueOf(condition)

		for i := range c.NumField() {
			field := c.Field(i)
			if !field.IsZero() {
				right := field.Interface()

				if r, ok := right.(string); ok && p.loader.IsPath(r) && !strings.HasPrefix(r, string(Entity_TIME)) && !strings.HasPrefix(r, string(Entity_ENV)) {
					p.cashTree.Add(p.action, r)
				}
			}

		}
	}
}

/*
load gets a value from a path or returns a literal value.

If path is a path (contains ":"), the function parses it and extracts the value
from the corresponding structure (source or target), environment variable, or time value.
If path is not a path, returns the path value itself as a literal value.

The function uses L1 cache to optimize field value retrieval. Before searching
for a field via reflection, it checks the cache using sessionID and path as key.
If the value is found in cache, it is returned immediately. After retrieving
a value via reflection, it is stored in cache for subsequent use within the
same evaluation session.

Note: Time paths (Entity_TIME) are not cached to ensure they always return current time.

Cache behavior is controlled by CashTree attached to the policy, which tracks
field access counts and disables caching for rarely accessed fields to optimize
memory usage.

Parameters:
  - ctx - context for operation cancellation and timeout control
  - source - first structure to search for fields (used for paths "source:...")
  - target - second structure to search for fields (used for paths "target:...")
  - value - path to field (e.g., "source:name", "target:owner", "env:VAR_NAME", "time:now") or literal value (e.g., "admin")
  - mustBePath - flag indicating whether value must be a path (true) or can be a literal (false)
  - sessionID - unique identifier for the current evaluation session (used as cache scope)

Returns:
  - any - found field value from structure, environment variable, time value, or literal path value
  - error - value retrieval error

Possible errors:
  - ErrInvalidPath - occurs if:
  - mustBePath=true, but value does not contain ":" (is not a path)
  - path parsing error (see splitPath)
  - field search error (see loadField)
  - environment variable does not exist (see loadEnv)
  - time path is invalid (see loadTime)
  - ErrInvalidType - entity is not a structure or pointer to structure (see loadField)
*/
func (p *Policy) load(ctx context.Context, source, target, item any, path any, mustBePath bool, sessionID string) (any, error) {
	pathString, ok := path.(string)
	if !ok {
		if mustBePath {
			return nil, NewErrInvalidPath(path, "path is not a valid path")
		}
		return path, nil
	}

	if !p.loader.IsPath(pathString) && mustBePath {
		return nil, NewErrInvalidPath(pathString, "path is not a valid path")
	}

	if p.cash != nil && p.cashTree != nil && !p.cashTree.IsDisabled(p.action, pathString) {
		value, err := p.cash.Get(ctx, sessionID, pathString)
		if err == nil && value != nil {
			return value, nil
		}
	}

	value, err := p.loader.Load(ctx, source, target, item, pathString)
	if err != nil {
		return false, err
	}

	if p.cash != nil && p.cashTree != nil && !p.cashTree.IsDisabled(p.action, pathString) {
		p.cash.Set(ctx, sessionID, pathString, value)
	}

	return value, nil
}

func (p *Policy) evaluate(ctx context.Context, source, target, item any, conditions map[string]Condition, sessionID string) (bool, error) {
	for field, condition := range conditions {
		select {
		case <-ctx.Done():
			return false, ErrCancelled
		default:
			left, err := p.load(ctx, source, target, item, field, true, sessionID)
			if err != nil {
				return false, err
			}

			if condition.In != nil {
				if m, err := p.conditionsMap.In(ctx, left, condition.In); err != nil || !m {
					return false, err
				}
			}
			if condition.Eq != nil {
				right, err := p.load(ctx, source, target, item, condition.Eq, false, sessionID)
				if err != nil {
					return false, err
				}

				if m, err := p.conditionsMap.Eq(ctx, left, right); err != nil || !m {
					return false, err
				}
			}
			if condition.Lt != nil {
				right, err := p.load(ctx, source, target, item, condition.Lt, false, sessionID)
				if err != nil {
					return false, err
				}

				if m, err := p.conditionsMap.Lt(ctx, left, right); err != nil || !m {
					return false, err
				}
			}
			if condition.Neq != nil {
				right, err := p.load(ctx, source, target, item, condition.Neq, false, sessionID)
				if err != nil {
					return false, err
				}

				if m, err := p.conditionsMap.Neq(ctx, left, right); err != nil || !m {
					return false, err
				}
			}
			if condition.Gt != nil {
				right, err := p.load(ctx, source, target, item, condition.Gt, false, sessionID)
				if err != nil {
					return false, err
				}

				if m, err := p.conditionsMap.Gt(ctx, left, right); err != nil || !m {
					return false, err
				}
			}
			if condition.Le != nil {
				right, err := p.load(ctx, source, target, item, condition.Le, false, sessionID)
				if err != nil {
					return false, err
				}

				if m, err := p.conditionsMap.Le(ctx, left, right); err != nil || !m {
					return false, err
				}
			}
			if condition.Ge != nil {
				right, err := p.load(ctx, source, target, item, condition.Ge, false, sessionID)
				if err != nil {
					return false, err
				}

				if m, err := p.conditionsMap.Ge(ctx, left, right); err != nil || !m {
					return false, err
				}
			}
			if condition.Any != nil {
				l := reflect.ValueOf(left)

				if l.Kind() != reflect.Slice && l.Kind() != reflect.Array {
					return false, NewErrInvalidType("slice or array", l.Kind().String())
				}

				found := false
				for i := range l.Len() {
					select {
					case <-ctx.Done():
						return false, ErrCancelled
					default:
						item := l.Index(i).Interface()

						if ok, err := p.evaluate(ctx, source, target, item, condition.Any, sessionID); err != nil {
							return false, err
						} else if ok {
							found = true
							break
						}
					}
				}

				// If no element satisfied all conditions, return false
				if !found {
					return false, nil
				}
			}
			if condition.All != nil {
				l := reflect.ValueOf(left)

				if l.Kind() != reflect.Slice && l.Kind() != reflect.Array {
					return false, NewErrInvalidType("slice or array", l.Kind().String())
				}

				for i := range l.Len() {
					select {
					case <-ctx.Done():
						return false, ErrCancelled
					default:
						item := l.Index(i).Interface()

						if ok, err := p.evaluate(ctx, source, target, item, condition.All, sessionID); err != nil || !ok {
							return false, err
						}
					}
				}
			}
		}
	}

	return true, nil
}

/*
Evaluate applies the policy to source and target structures for the specified action.

The function checks if the policy action matches the passed action. If yes,
all policy conditions are checked. The policy is considered fulfilled if all
conditions are met (logical AND).

The function supports cancellation through context.Context, allowing to interrupt
condition checking when context is cancelled.

Conditions are checked using direct field access (In, Eq, Neq, Lt, Gt, Le, Ge, Any, All) without
reflection, providing optimal performance. Field values are retrieved using the
load() method, which utilizes L1 cache to optimize performance. The cache is
scoped by sessionID, which is unique for each evaluation session. This allows
reusing field values within the same evaluation without repeated reflection-based
searches.

The Any and All conditions allow checking collections (slices/arrays). Any checks if at least
one element satisfies all nested conditions (OR logic for elements, AND logic for conditions).
All checks if all elements satisfy all nested conditions (AND logic for both elements and conditions).
Within Any/All conditions, you can use "item:" prefix to reference the current element being checked.

The cache behavior is controlled by the CashTree attached to the policy, which
tracks field access counts and can disable caching for fields that are accessed
less frequently than the configured threshold.

Parameters:
  - ctx - context for operation cancellation and timeout control
  - source - first structure to check (usually the action source)
  - target - second structure to check (usually the action target)
  - action - action in format "entity:action:extra..." to check
  - sessionID - unique identifier for the current evaluation session (used as cache scope)

Returns:
  - bool - policy application result:
  - true - policy matches action and all conditions are met
  - false - policy does not match action or at least one condition is not met
  - error - execution error if a problem occurred during condition checking

Possible errors:
  - ErrNilContext - context parameter is nil
  - ErrCancelled - operation was cancelled through context.Context
  - ErrInvalidPath - path parsing error or field search error in structure
  - ErrInvalidType - type error when getting field value (structure is not of that type)
  - ErrUncomparable - cannot compare values in condition (incompatible types)
  - ErrInexpectedBehavior - internal error: policy is nil
*/
func (p *Policy) Evaluate(ctx context.Context, source, target any, action string, sessionID string) (bool, error) {
	if p == nil {
		return false, NewErrInexpectedBehavior("Policy.Evaluate()", "policy is nil")
	}
	if ctx == nil {
		return false, ErrNilContext
	}

	if p.action != action {
		return false, nil
	}

	return p.evaluate(ctx, source, target, nil, p.conditions, sessionID)
}

/*
IsValid checks the validity of the policy.

The function performs comprehensive validation:
 1. Action format - must be at least 2 parts separated by ":" (entity:action:extra...)
 2. Absence of empty parts in action (no empty strings between separators)
 3. Validity of all paths in conditions - each condition key must be a valid path
 4. Validity of all paths in condition values (Eq, Neq, Lt, Gt, Le, Ge, In)
 5. Validity of all paths in nested Any/All conditions
 6. Item entity ("item:") can only be used within Any/All conditions

This method should be called before using the policy in the engine to ensure
all required fields are properly formatted and paths are valid.

Returns:
  - error - validity error if policy is invalid, nil if policy is valid

Possible errors:
  - ErrInvalidPath - occurs if:
  - action contains less than 2 parts (minimum entity and action)
  - action contains empty parts
  - paths in conditions are invalid
  - item: entity is used outside of Any/All conditions
*/
func (p *Policy) IsValid() error {
	actions := strings.Split(p.action, PATH_SEP)
	if len(actions) < MIN_ACTION_PARTS {
		return NewErrInvalidPath(p.action, "not enough parts of action. use: entity:action:extra1:extra2 etc")
	}

	for i, action := range actions {
		if action == "" {
			return NewErrInvalidPath(p.action, fmt.Sprintf("empty part: %v", i))
		}
	}

	// Validate all conditions recursively
	return p.validateConditions(p.conditions, false)
}

// validateConditions recursively validates all conditions and their paths.
// allowItemEntity indicates whether "item:" entity is allowed in paths (true for Any/All, false for top-level).
func (p *Policy) validateConditions(conditions map[string]Condition, allowItemEntity bool) error {
	for left, condition := range conditions {
		// Validate left path (condition key)
		if !p.loader.IsPath(left) {
			return NewErrInvalidPath(left, "path is not a valid path")
		}

		// Check if item: is used outside of Any/All
		if !allowItemEntity && strings.HasPrefix(left, string(Entity_ITEM)+PATH_SEP) {
			return NewErrInvalidPath(left, "item: entity can only be used within Any/All conditions")
		}

		// Validate condition values that can be paths
		if err := p.validateConditionValue(condition.Eq, allowItemEntity); err != nil {
			return err
		}
		if err := p.validateConditionValue(condition.Neq, allowItemEntity); err != nil {
			return err
		}
		if err := p.validateConditionValue(condition.Lt, allowItemEntity); err != nil {
			return err
		}
		if err := p.validateConditionValue(condition.Gt, allowItemEntity); err != nil {
			return err
		}
		if err := p.validateConditionValue(condition.Le, allowItemEntity); err != nil {
			return err
		}
		if err := p.validateConditionValue(condition.Ge, allowItemEntity); err != nil {
			return err
		}

		// Validate In slice values
		if condition.In != nil {
			for _, val := range condition.In {
				if err := p.validateConditionValue(val, allowItemEntity); err != nil {
					return err
				}
			}
		}

		// Validate nested Any conditions (item: is allowed here)
		if condition.Any != nil {
			if err := p.validateConditions(condition.Any, true); err != nil {
				return err
			}
		}

		// Validate nested All conditions (item: is allowed here)
		if condition.All != nil {
			if err := p.validateConditions(condition.All, true); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateConditionValue validates a single condition value if it's a path.
// If a string contains ":" but is not a valid path, it's considered an error.
func (p *Policy) validateConditionValue(value any, allowItemEntity bool) error {
	if value == nil {
		return nil
	}

	// Check if value is a string
	if pathStr, ok := value.(string); ok {
		// If string contains ":", it should be a valid path
		if strings.Contains(pathStr, PATH_SEP) {
			if !p.loader.IsPath(pathStr) {
				return NewErrInvalidPath(pathStr, "path is not a valid path")
			}
			// Check if item: is used outside of Any/All
			if !allowItemEntity && strings.HasPrefix(pathStr, string(Entity_ITEM)+PATH_SEP) {
				return NewErrInvalidPath(pathStr, "item: entity can only be used within Any/All conditions")
			}
		}
		// If string doesn't contain ":", it's a literal value, no validation needed
	}

	return nil
}

/*
Effect returns the policy effect (Effect_ALLOW or Effect_DENY).

The effect determines how the policy result is interpreted:
  - Effect_ALLOW - policy allows the action if conditions are met
  - Effect_DENY - policy denies the action if conditions are not met

Returns:
  - Effect - policy effect value
*/
func (p *Policy) Effect() Effect {
	return p.effect
}
