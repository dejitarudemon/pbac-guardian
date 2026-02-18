/*
Package base provides basic types and functions for working with policies,
conditions, effects, and entities in the access control system.

The package contains definitions of policies, comparison conditions, effects (allow/deny),
entities (source/target) and interfaces for custom comparison.
*/
package base

import "slices"

/*
Entity represents the entity to which a field in a policy belongs.

Entity is used in field paths to indicate which structure (source or target)
to extract the value from. For example, in path "source:name", Entity_SOURCE indicates
that the name field should be taken from the source structure.

Valid values:
  - Entity_SOURCE ("source") - field belongs to source structure
  - Entity_TARGET ("target") - field belongs to target structure
*/
type Entity string

const (
	Entity_TARGET Entity = "target"
	Entity_SOURCE Entity = "source"
)

var (
	// List of valid entities
	AVALIABLE_ENTITIES = []Entity{
		Entity_SOURCE,
		Entity_TARGET,
	}
)

/*
IsValid checks if the Entity value is valid.

Only Entity_SOURCE and Entity_TARGET are considered valid.

Returns:
  - bool - true if value is valid, false otherwise
*/
func (e Entity) IsValid() bool {
	return slices.Contains(AVALIABLE_ENTITIES, e)
}
