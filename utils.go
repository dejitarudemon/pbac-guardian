package guardian

import (
	"github.com/dejitarudemon/pbac-guardian/internal/base"
	"github.com/dejitarudemon/pbac-guardian/internal/cashing"
	"github.com/google/uuid"
)

/*
export converts a list of raw policies into a map of Policy instances organized by actions and effects.

The function performs comprehensive validation and organization of policies:
 1. Checks for duplicate policy names using a map for O(1) lookup efficiency
 2. Creates Policy instances from RawPolicy using base.NewPolicy, which validates each policy
    including path validation, condition validation, and effect validation
 3. Groups policies by actions and effects for fast lookup during evaluation
 4. Attaches condition functions configuration, cache instance, and cache tree to each policy
    for use during evaluation

The resulting map structure allows efficient lookup: first by action, then by effect (ALLOW or DENY),
enabling the Evaluate method to quickly find relevant policies and process them in the correct order
(DENY policies first, then ALLOW policies).

Parameters:
  - rawPolices - list of raw policies to export and group
  - cash - L1 cache instance for storing field values (can be nil to disable caching)
  - config - configuration containing condition functions map and cache disable threshold
  - tree - cache tree for tracking field access counts and disabling cache for rarely accessed fields

Returns:
  - map[string]map[base.Effect][]*base.Policy - nested map of policy pointers where:
  - outer key - action in format "entity:action:extra..."
  - inner key - effect (base.Effect_ALLOW or base.Effect_DENY)
  - value - list of all policies for this action and effect combination
  - error - export error if duplicate names or invalid policies are found

Possible errors:
  - ErrDuplicateName - policy name is already used by another policy in the list
  - errors from base.NewPolicy() - may include:
  - ErrInvalidType - if conditions map is invalid or condition types are incompatible
  - ErrInvalidPath - when path in policy conditions is invalid or field not found
  - other validation errors from base package

The function ensures that all policies are valid before being added to the map, providing
early error detection during engine initialization rather than during evaluation.
*/
func export(rawPolices []base.RawPolicy, cash cashing.Casher, config base.Config, tree *cashing.CashTree) (map[string]map[base.Effect][]*base.Policy, error) {
	mappedPolices := make(map[string]map[base.Effect][]*base.Policy)
	usedNames := make(map[string]struct{}, len(rawPolices))

	for _, rawPolicy := range rawPolices {
		if _, ok := usedNames[rawPolicy.Name]; ok {
			return nil, NewErrDuplicateName(rawPolicy.Name)
		}

		usedNames[rawPolicy.Name] = struct{}{}

		policy, err := base.NewPolicy(rawPolicy, config.ConditionsMap, cash, tree)
		if err != nil {
			return nil, err
		}

		if _, ok := mappedPolices[rawPolicy.Action]; !ok {
			mappedPolices[rawPolicy.Action] = map[base.Effect][]*base.Policy{
				base.Effect_ALLOW: make([]*base.Policy, 0),
				base.Effect_DENY:  make([]*base.Policy, 0),
			}
		}

		mappedPolices[rawPolicy.Action][rawPolicy.Effect] = append(mappedPolices[rawPolicy.Action][rawPolicy.Effect], policy)
	}

	return mappedPolices, nil
}

/*
generateNewSesstionID generates a unique identifier for an evaluation session.

The function creates a new UUID v4 string that serves as a session identifier.
This sessionID is used to scope the L1 cache for a single policy application,
allowing multiple concurrent evaluations to use the same cache instance without
interference. Each evaluation session gets its own cache scope, ensuring thread-safety
and preventing cache pollution between different evaluations.

The sessionID is generated once per Evaluate call and is used throughout the entire
evaluation process for that specific call. All field accesses within the same evaluation
session share the same sessionID, allowing the cache to efficiently store and retrieve
field values without conflicts between concurrent evaluations.

The UUID v4 format ensures uniqueness across all evaluation sessions, even in high-concurrency
scenarios. The sessionID is automatically cleared from the cache after the evaluation completes,
ensuring no memory leaks.

Returns:
  - string - unique session identifier (UUID v4 string in format "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")

Example:

	sessionID := generateNewSesstionID()
	// sessionID might be "550e8400-e29b-41d4-a716-446655440000"
*/
func generateNewSesstionID() string {
	return uuid.NewString()
}
