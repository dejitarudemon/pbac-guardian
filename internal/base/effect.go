/*
Package base provides basic types and functions for working with policies,
conditions, effects, and entities in the access control system.

The package contains definitions of policies, comparison conditions, effects (allow/deny),
entities (source/target) and interfaces for custom comparison.
*/
package base

import "slices"

var (
	// AVAILABLE_EFFECTS is a list of all valid policy effect values.
	// Used internally for validation of effect values in policies.
	AVAILABLE_EFFECTS = []Effect{
		Effect_ALLOW, // allow
		Effect_DENY,  // deny
	}
)

/*
Effect represents the effect of applying a policy.

Effect determines what happens if all policy conditions are met:
  - Effect_ALLOW allows the action
  - Effect_DENY denies the action

Policies with DENY effect have priority: if at least one DENY policy
does not pass the check, the action is denied, even if there are ALLOW policies.

Valid values:
  - Effect_ALLOW ("allow") - allow action when conditions are met
  - Effect_DENY ("deny") - deny action when conditions are not met
*/
type Effect string

const (
	Effect_ALLOW Effect = "allow"
	Effect_DENY  Effect = "deny"
)

/*
IsValid checks if the Effect value is valid.

Only Effect_ALLOW and Effect_DENY are considered valid.

Returns:
  - bool - true if value is valid, false otherwise
*/
func (e Effect) IsValid() bool {
	return slices.Contains(AVAILABLE_EFFECTS, e)
}
