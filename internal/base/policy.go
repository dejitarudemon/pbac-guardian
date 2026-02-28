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

Fields:
  - Name - unique policy name (used for identification)
  - Action - action in format "entity:action:extra1:extra2..." (e.g., "user:read:profile")
  - Effect - policy effect: Effect_ALLOW (allow) or Effect_DENY (deny)
  - Conditions - map of conditions. Key - path to field in format "source:field" or "target:field",
    value - condition to check (Contains, Eq, Neq, Lt)

Example usage:

	type User struct {
		Name string `pbac-guardian:"name"`
		Role string `pbac-guardian:"role"`
		Age  int    `pbac-guardian:"age"`
	}

	type Document struct {
		Owner string   `pbac-guardian:"owner"`
		Tags  []string `pbac-guardian:"tags"`
	}

	// Policy: allow admins to read documents
	policy1 := Policy{
		Name:   "admin-read",
		Action: "user:read:document",
		Effect: Effect_ALLOW,
		Conditions: map[string]Condition{
			"source:role": {
				Eq: "admin",
			},
		},
	}

	// Policy: allow document owner to read
	policy2 := Policy{
		Name:   "owner-read",
		Action: "user:read:document",
		Effect: Effect_ALLOW,
		Conditions: map[string]Condition{
			"source:name": {
				Eq: "target:owner", // compare fields of two structures
			},
		},
	}

	// Policy: deny reading documents with "private" tag for users under 18
	policy3 := Policy{
		Name:   "age-restriction",
		Action: "user:read:document",
		Effect: Effect_DENY,
		Conditions: map[string]Condition{
			"source:age": {
				Lt: 18,
			},
			"target:tags": {
				Contains: []any{"private"},
			},
		},
	}
*/
type Policy struct {
	Name          string               `json:"name"`
	Action        string               `json:"action"`
	Effect        Effect               `json:"effect"`
	Conditions    map[string]Condition `json:"conditions"`
	ConditionsMap *ConditionsMap       `json:"-"`
	Cash          Casher               `json:"-"`
}

/*
isPath checks if path is a path to a structure field.

The path to a field must contain the PATH_SEP separator (":"), indicating that
this is not a literal value, but a reference to a field in the source or target structure.

Parameters:
  - path - string to check

Returns:
  - bool - true if path contains separator ":" and is a path, false otherwise
*/
func (p *Policy) isPath(path string) bool {
	return strings.Contains(path, PATH_SEP)
}

/*
parsePath parses a path from Conditions into an entity and a path to a field.

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
func (p *Policy) parsePath(path string) (*Entity, []string, error) {
	if !p.isPath(path) {
		return nil, nil, NewErrInvalidPath(path, "it is not a path")
	}

	separeted := strings.Split(path, PATH_SEP)
	if len(separeted) < MIN_PATH_LEN {
		return nil, nil, NewErrInvalidPath(path, fmt.Sprintf("expected at least %v parts, but got %v", MIN_ACTION_PARTS, len(separeted)))
	}

	entity := Entity(separeted[0])
	fields := separeted[1:]

	if !entity.IsValid() {
		return nil, nil, NewErrInvalidPath(path, fmt.Sprintf("path allocates to unknown entity: %v", entity))
	}

	return &entity, fields, nil
}

/*
getValue finds a field in a structure by path and returns its value.

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
func (p *Policy) getValue(entity any, path []string) (any, error) {
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
				return p.getValue(field.Interface(), path[1:])
			}

			return field.Interface(), nil
		}
	}

	return nil, NewErrInvalidPath(path[0], "param doesn't exist")
}

/*
get gets a value from a path or returns a literal value.

If path is a path (contains ":"), the function parses it and extracts the value
from the corresponding structure (source or target). If path is not a path,
returns the path value itself as a literal value.

The function uses L1 cache to optimize field value retrieval. Before searching
for a field via reflection, it checks the cache using sessionID and path as key.
If the value is found in cache, it is returned immediately. After retrieving
a value via reflection, it is stored in cache for subsequent use within the
same evaluation session.

Parameters:
  - ctx - context for operation cancellation and timeout control
  - source - first structure to search for fields (used for paths "source:...")
  - target - second structure to search for fields (used for paths "target:...")
  - path - path to field (e.g., "source:name") or literal value (e.g., "admin")
  - mustBePath - flag indicating whether path must be a path (true) or can be a literal (false)
  - cash - L1 cache instance for storing field values (can be nil to disable caching)
  - sessionID - unique identifier for the current evaluation session (used as cache scope)

Returns:
  - any - found field value from structure or literal path value
  - error - value retrieval error

Possible errors:
  - ErrInvalidPath - occurs if:
  - mustBePath=true, but path does not contain ":" (is not a path)
  - path parsing error (see parsePath)
  - field search error (see getValue)
  - ErrInvalidType - entity is not a structure or pointer to structure (see getValue)
*/
func (p *Policy) get(ctx context.Context, source, target any, path string, mustBePath bool, sessionID string) (any, error) {
	if !p.isPath(path) {
		if mustBePath {
			return nil, NewErrInvalidPath(path, "must be a path, but it's literal value")
		}
		return path, nil
	}

	// Если кеш не отключен, ищем в нем по id сессии и ключу (пути до искомого поля)
	if p.Cash != nil {
		value, err := p.Cash.Get(ctx, sessionID, path)
		if err == nil && value != nil {
			return value, nil
		}
	}

	entity, parsedPath, err := p.parsePath(path)
	if err != nil {
		return false, err
	}

	var value any

	switch *entity {
	case Entity_SOURCE:
		value, err = p.getValue(source, parsedPath)
	case Entity_TARGET:
		value, err = p.getValue(target, parsedPath)
	default:
		err = NewErrInvalidPath(path, fmt.Sprintf("unxpected entity: %v", entity))
	}

	if err != nil {
		return nil, err
	}

	if p.Cash != nil {
		p.Cash.Set(ctx, sessionID, path, value)
	}

	return value, nil
}

/*
Evaluate applies the policy to source and target structures for the specified action.

The function checks if the policy action matches the passed action. If yes,
all policy conditions are checked. The policy is considered fulfilled if all
conditions are met (logical AND).

The function supports cancellation through context.Context, allowing to interrupt
condition checking when context is cancelled.

Field values are retrieved using the get() method, which utilizes L1 cache to
optimize performance. The cache is scoped by sessionID, which is unique for each
evaluation session. This allows reusing field values within the same evaluation
without repeated reflection-based searches.

Parameters:
  - ctx - context for operation cancellation and timeout control
  - source - first structure to check (usually the action source)
  - target - second structure to check (usually the action target)
  - action - action in format "entity:action:extra..." to check
  - cash - L1 cache instance for storing field values (can be nil to disable caching)
  - sessionID - unique identifier for the current evaluation session (used as cache scope)

Returns:
  - bool - policy application result:
  - true - policy matches action and all conditions are met
  - false - policy does not match action or at least one condition is not met
  - error - execution error if a problem occurred during condition checking

Possible errors:
  - ErrCancelled - operation was cancelled through context.Context
  - ErrInvalidPath - path parsing error or field search error in structure
  - ErrInvalidType - type error when getting field value (structure is not of that type)
  - ErrUncomparable - cannot compare values in condition (incompatible types)
  - ErrInexpectedBehavior - internal error: condition function not found in CONDITION_TO_FUNC
*/
func (p *Policy) Evaluate(ctx context.Context, source, target any, action string, sessionID string) (bool, error) {
	if p == nil {
		return false, NewErrInexpectedBehavior("Policy.Evaluate()", "policy is nil")
	}
	if ctx == nil {
		return false, ErrNilContext
	}

	if p.Action != action {
		return false, nil
	}

	match := true
	t := reflect.TypeFor[Condition]()

	for field, condition := range p.Conditions {
		left, err := p.get(ctx, source, target, field, true, sessionID)
		if err != nil {
			return false, err
		}

		c := reflect.ValueOf(condition)

		for i := range c.NumField() {
			select {
			case <-ctx.Done():
				return false, ErrCancelled
			default:
				if !c.Field(i).IsZero() {
					if f := p.ConditionsMap.Select(t.Field(i).Name); f != nil {

						right := c.Field(i).Interface()

						if r, ok := right.(string); ok {
							right, err = p.get(ctx, source, target, r, false, sessionID)
							if err != nil {
								return false, err
							}
						}

						m, err := f(ctx, left, right)
						if err != nil {
							return false, err
						}

						match = match && m

						if !match {
							return false, nil
						}
					} else {
						return false, NewErrInexpectedBehavior("Policy.Evaluate()", fmt.Sprintf("condition func for %v doesn't exist", t.Field(i).Name))
					}
				}
			}
		}
	}

	return match, nil
}

/*
IsValid checks the validity of the policy.

The function performs comprehensive validation:
 1. Action format - must be at least 2 parts separated by ":" (entity:action:extra...)
 2. Absence of empty parts in action (no empty strings between separators)
 3. Validity of all paths in conditions (via parsePath) - each condition key must be a valid path

This method should be called before using the policy in the engine to ensure
all required fields are properly formatted and paths are valid.

Returns:
  - error - validity error if policy is invalid, nil if policy is valid

Possible errors:
  - ErrInvalidPath - occurs if:
  - action contains less than 2 parts (minimum entity and action)
  - action contains empty parts
  - paths in conditions are invalid (see parsePath for details)
*/
func (p *Policy) IsValid() error {
	actions := strings.Split(p.Action, PATH_SEP)
	if len(actions) < MIN_ACTION_PARTS {
		return NewErrInvalidPath(p.Action, "not enough parts of action. use: entity:action:extra1:extra2 etc")
	}

	for i, action := range actions {
		if action == "" {
			return NewErrInvalidPath(p.Action, fmt.Sprintf("empty part: %v", i))
		}
	}

	for field := range p.Conditions {
		if _, _, err := p.parsePath(field); err != nil {
			return err
		}
	}

	return nil
}
