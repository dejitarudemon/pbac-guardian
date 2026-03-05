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
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dejitarudemon/pbac-guardian/internal/cashing"
)

// TAG_KEY is the tag key used for struct field tagging.
// Fields in structures should be tagged with `pbac-guardian:"field_name"` to be
// accessible in policy conditions. The tag value specifies the field name used
// in paths like "source:field_name" or "target:field_name".
const TAG_KEY = "pbac-guardian"

const (
	// PATH_SEP is the separator used in paths to exported fields.
	// Paths have format "entity:field1:field2..." where ":" separates parts.
	// Example: "source:name", "target:user:email"
	PATH_SEP = ":"

	// TAG_SEP is the separator used inside struct field tags.
	// Tags can have format "field_name,flag1,flag2" where "," separates parts.
	// Currently only the first part (field name) is used.
	TAG_SEP = ","

	// MIN_PATH_LEN is the minimum number of parts in a valid path.
	// A path must have at least 2 parts: entity and one field name.
	// Example: "source:name" (2 parts) is valid, "source" (1 part) is invalid.
	MIN_PATH_LEN = 2

	// MIN_ACTION_PARTS is the minimum number of parts in a valid action.
	// An action must have at least 2 parts: entity and action name.
	// Example: "user:read" (2 parts) is valid, "read" (1 part) is invalid.
	MIN_ACTION_PARTS = 2

	TIME_MODIFIER_SEP = "|"
	MODIFIERS_PARTS   = 2
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
    value - condition to check (In, Eq, Neq, Lt, Gt, Le, Ge)
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

				if r, ok := p.toPath(right); ok && !strings.HasPrefix(r, string(Entity_TIME)) {
					p.cashTree.Add(p.action, r)
				}
			}

		}
	}
}

/*
toPath checks if a value is a path to a structure field, environment variable, or time value.

A path must contain the PATH_SEP separator (":"), indicating that this is not
a literal value, but a reference to:
  - a field in the source or target structure (e.g., "source:name", "target:owner")
  - an environment variable (e.g., "env:VARIABLE_NAME")
  - a time value (e.g., "time:now", "time:now:1|day")

Parameters:
  - path - value to check (can be any type, but must be a string to be a valid path)

Returns:
  - string - the path string if path is a string, empty string otherwise
  - bool - true if path is a string and contains separator ":" (is a path), false otherwise
*/
func (p *Policy) toPath(path any) (string, bool) {
	if str, ok := path.(string); ok {
		return str, strings.Contains(str, PATH_SEP)
	}
	return "", false
}

/*
splitPath parses a path from Conditions into an entity and a path to a field.

The path must have the format "entity:field1:field2...", where entity is "source",
"target", "env", or "time", and field1, field2... is a hierarchical path to a field
in the structure, environment variable name, or time specification.

Valid path examples:
  - "source:name" - field name in source structure
  - "target:user:email" - field email in nested user structure in target
  - "env:VARIABLE_NAME" - environment variable VARIABLE_NAME
  - "time:now" - current time
  - "time:now:1|day" - current time plus 1 day (with modifier)

Parameters:
  - value - path to parse in format "entity:field1:field2..." (can be any type, but must be a string to be valid)

Returns:
  - *Entity - pointer to entity (Entity_SOURCE, Entity_TARGET, Entity_ENV, or Entity_TIME)
  - []string - path to field as array of strings ["field1", "field2", ...]
  - For Entity_ENV: contains variable name
  - For Entity_TIME: contains time specification (e.g., ["now"] or ["now", "1|day"])
  - For Entity_SOURCE/Entity_TARGET: contains hierarchical path to field
  - error - path parsing error

Possible errors:
  - ErrInvalidPath - occurs if:
  - value is not a string or does not contain separator ":" (is not a path)
  - path contains less than 2 parts (minimum entity and one field/variable name)
  - first part of path is not a valid entity (not "source", "target", "env", or "time")
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

func (p *Policy) loadFieldFromStruct(v reflect.Value, t reflect.Type, path []string) (any, error) {
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
loadFieldFromMap retrieves a value from a map by key and continues path traversal if needed.

The function supports any map type with string keys (e.g., map[string]any, map[string]Group).
It uses reflection to access map values by key, allowing traversal through map fields
in structures. If the path has more segments, it recursively calls loadField to continue
traversing the found value.

This function enables access to fields in structures stored in map fields, for example:
  - Structure with map field: User { Groups map[string]Group }
  - Path: "source:groups:admins:name"
  - "groups" - map field in User
  - "admins" - key in map
  - "name" - field in Group structure

Parameters:
  - value - reflect.Value of the map to search in
  - path - array of path segments, where path[0] is the map key to look for

Returns:
  - any - found value from map, or result of recursive traversal if path continues
  - error - error if map key doesn't exist or value is unexported

Possible errors:
  - ErrInvalidType - value is not a map
  - ErrInvalidPath - occurs if:
  - map key doesn't exist (path[0] not found in map)
  - map value is unexported (not accessible via CanInterface())
*/
func (p *Policy) loadFieldFromMap(value reflect.Value, path []string) (any, error) {
	if value.Kind() != reflect.Map {
		return nil, NewErrInvalidType("map", value.Kind().String())
	}

	keyValue := reflect.ValueOf(path[0])
	mapValue := value.MapIndex(keyValue)

	if !mapValue.IsValid() {
		return nil, NewErrInvalidPath(path[0], "param doesn't exist")
	}

	if !mapValue.CanInterface() {
		return nil, NewErrInvalidPath(path[0], "got unexported value")
	}

	item := mapValue.Interface()

	if len(path) > 1 {
		return p.loadField(item, path[1:])
	}

	return item, nil
}

/*
loadField finds a field in a structure or map by path and returns its value.

The function recursively traverses the path using TAG_KEY ("pbac-guardian") tags to find
fields in structures, or by key lookup in maps. The field must be exported (capitalized)
and have a tag with the corresponding name. The function supports nested structures and
maps by recursively calling itself for each path segment.

The function supports:
  - Structures: access fields by tag names (e.g., "source:role")
  - Maps: access values by keys (e.g., "source:groups:admins")
  - Nested combinations: structures with map fields containing structures
    (e.g., "source:groups:admins:name" where groups is map[string]Group)

Examples:
  - Structure field: "source:role" - accesses role field in source structure
  - Map value: "source:metadata:key" - accesses key in metadata map
  - Map field in structure: "source:groups:admins:name" - accesses name field in Group
    structure stored in groups map under "admins" key

Parameters:
  - entity - entity (structure, pointer to structure, or map) to search for field
  - path - path to field as array of strings ["field1", "field2", ...]

Returns:
  - any - found field value
  - error - field search error

Possible errors:
  - ErrInvalidType - entity is not a structure, pointer to structure, or map
  - ErrInvalidPath - occurs if:
  - field not found (no tag with corresponding name in structure)
  - map key doesn't exist
  - field is unexported (not accessible via CanInterface())
  - entity is nil pointer
*/
func (p *Policy) loadField(entity any, path []string) (any, error) {
	v := reflect.ValueOf(entity)

	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, NewErrInvalidType(fmt.Sprintf("%v or %v or %v", reflect.Pointer.String(), reflect.Struct.String(), reflect.Map.String()), v.Kind().String())
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		return p.loadFieldFromStruct(v, reflect.TypeOf(entity), path)
	case reflect.Map:
		return p.loadFieldFromMap(v, path)
	}

	return nil, NewErrInvalidPath(path[0], "param doesn't exist")
}

/*
loadEnv retrieves a value from environment variables.

The function uses os.LookupEnv to get the value of an environment variable
by its name. This is used when Entity_ENV is specified in a path like "env:VARIABLE_NAME".

Parameters:
  - path - name of the environment variable to retrieve

Returns:
  - any - value of the environment variable as string
  - error - error if environment variable does not exist

Possible errors:
  - ErrInvalidPath - environment variable with the specified name does not exist
*/
func (p *Policy) loadEnv(path string) (any, error) {
	if value, ok := os.LookupEnv(path); ok {
		return value, nil
	}

	return nil, NewErrInvalidPath(path, "env doesn't exists")
}

/*
loadTimeModications applies time modifiers to a target time value.

The function processes modifiers in format "value|unit" where:
  - value is an integer (currently not used, MODIFIERS_PARTS is used instead)
  - unit is one of: "day", "hour", "minute", "second", "milisecond"

Modifiers are applied sequentially. Each modifier adds a fixed amount of time
based on MODIFIERS_PARTS constant (currently 2).

Parameters:
  - path - array of modifier strings in format "value|unit" (e.g., ["1|day", "2|hour"])
  - target - base time value to apply modifiers to

Returns:
  - time.Time - target time with modifiers applied
  - error - error if modifier format is invalid or unit is not supported

Possible errors:
  - ErrInvalidPath - modifier format is invalid (not "value|unit" or wrong number of parts)
  - ErrInvalidType - modifier value is not a number or unit is not supported
*/
func (p *Policy) loadTimeModications(path []string, target time.Time) (time.Time, error) {
	if len(path) == 0 {
		return target, nil
	}

	modifiers := strings.Split(path[0], TIME_MODIFIER_SEP)
	if len(modifiers) != MODIFIERS_PARTS {
		return target, NewErrInvalidPath(path[0], "expected int|day/hour/second/month/year")
	}

	modifierValue, err := strconv.Atoi(modifiers[0])
	if err != nil {
		return target, NewErrInvalidType(reflect.Int.String(), reflect.TypeOf(modifierValue).String())
	}

	switch modifiers[1] {
	case "day":
		return target.Add(time.Hour * 24 * MODIFIERS_PARTS), nil
	case "hour":
		return target.Add(time.Hour * MODIFIERS_PARTS), nil
	case "minute":
		return target.Add(time.Minute * MODIFIERS_PARTS), nil
	case "second":
		return target.Add(time.Second * MODIFIERS_PARTS), nil
	case "millisecond":
		return target.Add(time.Millisecond * MODIFIERS_PARTS), nil
	}

	return target, NewErrInvalidType("day|hour|minute|second|millisecond", modifiers[1])
}

/*
loadTime loads a time value from a path specification.

The function supports the following time specifications:
  - "now" - current time (time.Now())
  - "now:modifier1:modifier2..." - current time with modifiers applied

Modifiers are in format "value|unit" and are processed by loadTimeModications.
Supported units: day, hour, minute, second, millisecond.

Examples:
  - "time:now" -> current time
  - "time:now:1|day" -> current time + 2 days (MODIFIERS_PARTS = 2)
  - "time:now:2|hour" -> current time + 2 hours

Parameters:
  - path - array of strings specifying time value (e.g., ["now"] or ["now", "1|day"])

Returns:
  - time.Time - time value based on specification
  - error - error if time specification is invalid

Possible errors:
  - ErrInvalidPath - time base is not "now" or modifier format is invalid
  - ErrInvalidType - modifier value is not a number or unit is not supported
*/
func (p *Policy) loadTime(path []string) (time.Time, error) {
	var t time.Time

	switch path[0] {
	case "now":
		t = time.Now()
	default:
		return t, NewErrInvalidPath(path[0], "expected now, but got invalid")
	}

	return p.loadTimeModications(path[1:], t)
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
	case Entity_ENV:
		v, err = p.loadEnv(parsedPath[0])
	case Entity_TIME:
		v, err = p.loadTime(parsedPath)
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

Conditions are checked using direct field access (In, Eq, Neq, Lt, Gt, Le, Ge) without
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

			if condition.In != nil {
				if m, err := p.conditionsMap.In(ctx, left, condition.In); err != nil || !m {
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
			if condition.Gt != nil {
				right, err := p.load(ctx, source, target, condition.Gt, false, sessionID)
				if err != nil {
					return false, err
				}

				if m, err := p.conditionsMap.Gt(ctx, left, right); err != nil || !m {
					return false, err
				}
			}
			if condition.Le != nil {
				right, err := p.load(ctx, source, target, condition.Le, false, sessionID)
				if err != nil {
					return false, err
				}

				if m, err := p.conditionsMap.Le(ctx, left, right); err != nil || !m {
					return false, err
				}
			}
			if condition.Ge != nil {
				right, err := p.load(ctx, source, target, condition.Ge, false, sessionID)
				if err != nil {
					return false, err
				}

				if m, err := p.conditionsMap.Ge(ctx, left, right); err != nil || !m {
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
