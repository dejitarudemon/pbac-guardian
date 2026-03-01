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

// tag key used for struct field tagging
const TAG_KEY = "pbac-guardian"

const (
	// separator for path to exported field
	PATH_SEP = ":"

	// separator inside tag
	TAG_SEP = ","

	// minimum path size
	MIN_PATH_LEN = 2

	// minimum number of elements in action (entity and action)
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
    value - condition to check (Contains, Eq, Neq, Lt)
  - conditionsMap - configuration for condition functions (Contains, Eq, Neq, Lt)
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
  - conditions - configuration for condition functions (Contains, Eq, Neq, Lt). Must not be nil.
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

	if conditions.Contains == nil {
		return nil, NewErrInvalidType("ConditionsMap.Contains", nil)
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

	p := Policy{
		name:          raw.Action,
		action:        raw.Action,
		effect:        raw.Effect,
		conditions:    raw.Conditions,
		conditionsMap: conditions,
		cash:          cash,
		cashTree:      cashTree,
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
		p.cashTree.Add(p.action, left)

		c := reflect.ValueOf(condition)

		for i := range c.NumField() {
			field := c.Field(i)
			if !field.IsZero() {
				right := field.Interface()

				if r, ok := p.toPath(right); ok {
					p.cashTree.Add(p.action, r)
				}
			}

		}
	}
}

/*
toPath checks if path is a path to a structure field.

The path to a field must contain the PATH_SEP separator (":"), indicating that
this is not a literal value, but a reference to a field in the source or target structure.

Parameters:
  - path - string to check

Returns:
  - bool - true if path contains separator ":" and is a path, false otherwise
*/
func (p *Policy) toPath(path any) (string, bool) {
	if str, ok := path.(string); ok {
		return str, strings.Contains(str, PATH_SEP)
	}
	return "", false
}

/*
splitPath parses a path from Conditions into an entity and a path to a field.

The path must have the format "entity:field1:field2...", where entity is "source"
or "target", and field1, field2... is a hierarchical path to a field in the structure.

Valid path examples:
  - "source:name" - field name in source structure
  - "target:user:email" - field email in nested user structure in target

Parameters:
  - path - path to parse in format "entity:field1:field2..."

Returns:
  - *Entity - pointer to entity (Entity_SOURCE or Entity_TARGET)
  - []string - path to field as array of strings ["field1", "field2", ...]
  - error - path parsing error

Possible errors:
  - ErrInvalidPath - occurs if:
  - path does not contain separator ":" (is not a path)
  - path contains less than 2 parts (minimum entity and one field)
  - first part of path is not a valid entity (not "source" and not "target")
*/
func (p *Policy) splitPath(value any) (*Entity, []string, error) {
	path, ok := p.toPath(value)
	if !ok {
		return nil, nil, NewErrInvalidPath(path, "it is not a path")
	}

	separeted := strings.Split(path, PATH_SEP)
	if len(separeted) < MIN_PATH_LEN {
		return nil, nil, NewErrInvalidPath(path, fmt.Sprintf("expected at least %v parts, but got %v", MIN_ACTION_PARTS, len(separeted)))
	}

	entity := Entity(separeted[0])

	if !entity.IsValid() {
		return nil, nil, NewErrInvalidPath(path, fmt.Sprintf("path allocates to unknown entity: %v", entity))
	}

	return &entity, separeted[1:], nil
}

/*
loadField finds a field in a structure by path and returns its value.

The function recursively traverses the path using TAG_KEY ("pbac-guardian") tags to find
fields. The field must be exported (capitalized) and have a tag with the
corresponding name. The function supports nested structures by recursively
calling itself for each path segment.

Parameters:
  - entity - entity (structure or pointer to structure) to search for field
  - path - path to field as array of strings ["field1", "field2", ...]

Returns:
  - any - found field value
  - error - field search error

Possible errors:
  - ErrInvalidType - entity is not a structure or pointer to structure
  - ErrInvalidPath - occurs if:
  - field not found (no tag with corresponding name)
  - field is unexported (not accessible via CanInterface())
  - entity is nil pointer
*/
func (p *Policy) loadField(entity any, path []string) (any, error) {
	v := reflect.ValueOf(entity)

	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, NewErrInvalidType(fmt.Sprintf("%v or %v", reflect.Pointer.String(), reflect.Struct.String()), v.Kind().String())
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, NewErrInvalidType(fmt.Sprintf("%v or %v", reflect.Pointer.String(), reflect.Struct.String()), v.Kind().String())
	}

	t := v.Type()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get(TAG_KEY)

		if tag != "" {
			flags := strings.Split(tag, TAG_SEP)
			tagValue := strings.TrimSpace(flags[0])

			if tagValue != path[0] {
				continue
			}

			if !field.CanInterface() {
				return nil, NewErrInvalidPath(path[0], "got unexpored field")
			}

			if len(path) > 1 {
				return p.loadField(field.Interface(), path[1:])
			}

			return field.Interface(), nil
		}
	}

	return nil, NewErrInvalidPath(path[0], "param doesn't exist")
}

/*
load gets a value from a path or returns a literal value.

If path is a path (contains ":"), the function parses it and extracts the value
from the corresponding structure (source or target). If path is not a path,
returns the path value itself as a literal value.

The function uses L1 cache to optimize field value retrieval. Before searching
for a field via reflection, it checks the cache using sessionID and path as key.
If the value is found in cache, it is returned immediately. After retrieving
a value via reflection, it is stored in cache for subsequent use within the
same evaluation session.

Cache behavior is controlled by CashTree attached to the policy, which tracks
field access counts and disables caching for rarely accessed fields to optimize
memory usage.

Parameters:
  - ctx - context for operation cancellation and timeout control
  - source - first structure to search for fields (used for paths "source:...")
  - target - second structure to search for fields (used for paths "target:...")
  - path - path to field (e.g., "source:name") or literal value (e.g., "admin")
  - mustBePath - flag indicating whether path must be a path (true) or can be a literal (false)
  - sessionID - unique identifier for the current evaluation session (used as cache scope)

Returns:
  - any - found field value from structure or literal path value
  - error - value retrieval error

Possible errors:
  - ErrInvalidPath - occurs if:
  - mustBePath=true, but path does not contain ":" (is not a path)
  - path parsing error (see splitPath)
  - field search error (see loadField)
  - ErrInvalidType - entity is not a structure or pointer to structure (see loadField)
*/
func (p *Policy) load(ctx context.Context, source, target any, value any, mustBePath bool, sessionID string) (any, error) {
	path, ok := p.toPath(value)
	if !ok {
		if mustBePath {
			return nil, NewErrInvalidPath(path, "must be a path, but it's literal value")
		}
		return value, nil
	}

	if p.cash != nil && p.cashTree != nil && !p.cashTree.IsDisabled(p.action, path) {
		value, err := p.cash.Get(ctx, sessionID, path)
		if err == nil && value != nil {
			return value, nil
		}

	}

	entity, parsedPath, err := p.splitPath(path)
	if err != nil {
		return false, err
	}

	var v any

	switch *entity {
	case Entity_SOURCE:
		v, err = p.loadField(source, parsedPath)
	case Entity_TARGET:
		v, err = p.loadField(target, parsedPath)
	default:
		err = NewErrInvalidPath(path, fmt.Sprintf("unxpected entity: %v", entity))
	}

	if err != nil {
		return nil, err
	}

	if p.cash != nil && p.cashTree != nil && !p.cashTree.IsDisabled(p.action, path) {
		p.cash.Set(ctx, sessionID, path, v)
	}

	return v, nil
}

/*
Evaluate applies the policy to source and target structures for the specified action.

The function checks if the policy action matches the passed action. If yes,
all policy conditions are checked. The policy is considered fulfilled if all
conditions are met (logical AND).

The function supports cancellation through context.Context, allowing to interrupt
condition checking when context is cancelled.

Conditions are checked using direct field access (Contains, Eq, Neq, Lt) without
reflection, providing optimal performance. Field values are retrieved using the
load() method, which utilizes L1 cache to optimize performance. The cache is
scoped by sessionID, which is unique for each evaluation session. This allows
reusing field values within the same evaluation without repeated reflection-based
searches.

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

	for field, condition := range p.conditions {
		select {
		case <-ctx.Done():
			return false, ErrCancelled
		default:
			left, err := p.load(ctx, source, target, field, true, sessionID)
			if err != nil {
				return false, err
			}

			if condition.Contains != nil {
				if m, err := p.conditionsMap.Contains(ctx, left, condition.Contains); err != nil || !m {
					return false, err
				}
			}
			if condition.Eq != nil {
				right, err := p.load(ctx, source, target, condition.Eq, false, sessionID)
				if err != nil {
					return false, err
				}

				if m, err := p.conditionsMap.Eq(ctx, left, right); err != nil || !m {
					return false, err
				}
			}
			if condition.Lt != nil {
				right, err := p.load(ctx, source, target, condition.Lt, false, sessionID)
				if err != nil {
					return false, err
				}

				if m, err := p.conditionsMap.Lt(ctx, left, right); err != nil || !m {
					return false, err
				}
			}
			if condition.Neq != nil {
				right, err := p.load(ctx, source, target, condition.Neq, false, sessionID)
				if err != nil {
					return false, err
				}

				if m, err := p.conditionsMap.Neq(ctx, left, right); err != nil || !m {
					return false, err
				}
			}
		}
	}

	return true, nil
}

/*
IsValid checks the validity of the policy.

The function performs comprehensive validation:
 1. Action format - must be at least 2 parts separated by ":" (entity:action:extra...)
 2. Absence of empty parts in action (no empty strings between separators)
 3. Validity of all paths in conditions (via splitPath) - each condition key must be a valid path

This method should be called before using the policy in the engine to ensure
all required fields are properly formatted and paths are valid.

Returns:
  - error - validity error if policy is invalid, nil if policy is valid

Possible errors:
  - ErrInvalidPath - occurs if:
  - action contains less than 2 parts (minimum entity and action)
  - action contains empty parts
  - paths in conditions are invalid (see splitPath for details)
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

	for field := range p.conditions {
		if _, _, err := p.splitPath(field); err != nil {
			return err
		}
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
