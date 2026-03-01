package guardian

import (
	"github.com/dejitarudemon/pbac-guardian/internal/base"
	"github.com/dejitarudemon/pbac-guardian/internal/cashing"
	"github.com/google/uuid"
)

/*
export converts a list of raw policies into a map of Policy instances organized by actions.

The function performs the following checks:
 1. Check for duplicate policy names (using map for O(1) lookup)
 2. Create Policy instances from RawPolicy using NewPolicy (which validates each policy)
 3. Group policies by actions for fast access
 4. Attach condition functions configuration, cache instance, and cache tree to each policy

Parameters:
  - rawPolices - list of raw policies to export and group
  - cash - L1 cache instance for storing field values (can be nil to disable caching)
  - config - configuration containing condition functions map and cache disable threshold
  - tree - cache tree for tracking field access counts and disabling cache for rarely accessed fields

Returns:
  - map[string][]*base.Policy - map of policy pointers where:
  - key - action in format "entity:action:extra..."
  - value - list of all policies for this action
  - error - export error if duplicate names or invalid policies are found

Possible errors:
  - ErrDuplicateName - policy name is already used by another policy in the list
  - errors from base.NewPolicy() - ErrInvalidType if conditions map is invalid, ErrInvalidPath when path in policy conditions is invalid
*/
func export(rawPolices []base.RawPolicy, cash cashing.Casher, config base.Config, tree *cashing.CashTree) (map[string][]*base.Policy, error) {
	mappedPolices := make(map[string][]*base.Policy)
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
			mappedPolices[rawPolicy.Action] = make([]*base.Policy, 0, 1)
		}

		mappedPolices[rawPolicy.Action] = append(mappedPolices[rawPolicy.Action], policy)
	}

	return mappedPolices, nil
}

/*
generateNewSesstionID generates a unique identifier for an evaluation session.

The function creates a new UUID string that serves as a session identifier.
This sessionID is used to scope the L1 cache for a single policy application,
allowing multiple concurrent evaluations to use the same cache instance without
interference. Each evaluation session gets its own cache scope, ensuring thread-safety
and preventing cache pollution between different evaluations.

Returns:
  - string - unique session identifier (UUID v4 string)
*/
func generateNewSesstionID() string {
	return uuid.NewString()
}
