package guardian

import (
	"slices"

	"github.com/dejitarudemon/pbac-guardian/internal/base"
	"github.com/google/uuid"
)

/*
export converts a list of policies into a map organized by actions.

The function performs the following checks:
 1. Check for duplicate policy names
 2. Validate each policy through Policy.IsValid()
 3. Group policies by actions for fast access

Parameters:
  - polices - list of policies to export and group

Returns:
  - map[string][]base.Policy - map of policies where:
  - key - action in format "entity:action:extra..."
  - value - list of all policies for this action
  - error - export error if duplicate names or invalid policies are found

Possible errors:
  - ErrDuplicateName - policy name is already used by another policy in the list
  - errors from base.Policy.IsValid() - ErrInvalidPath when path in policy conditions is invalid
*/
func export(polices []base.Policy) (map[string][]base.Policy, error) {
	mappedPolices := make(map[string][]base.Policy)
	usedNames := make([]string, 0, len(polices))

	for _, policy := range polices {
		if slices.Contains(usedNames, policy.Name) {
			return nil, NewErrDuplicateName(policy.Name)
		}

		usedNames = append(usedNames, policy.Name)

		if err := policy.IsValid(); err != nil {
			return nil, err
		}

		if _, ok := mappedPolices[policy.Action]; !ok {
			mappedPolices[policy.Action] = make([]base.Policy, 0, 1)
		}

		mappedPolices[policy.Action] = append(mappedPolices[policy.Action], policy)
	}

	return mappedPolices, nil
}

/*
generateNewSesstionID generates a unique identifier for an evaluation session.

The function creates a new UUID string that serves as a session identifier.
This sessionID is used to scope the L1 cache for a single policy application,
allowing multiple concurrent evaluations to use the same cache instance without
interference.

Returns:
  - string - unique session identifier (UUID string)
*/
func generateNewSesstionID() string {
	return uuid.NewString()
}
