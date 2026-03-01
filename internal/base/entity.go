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

Entity is used in field paths to indicate where to extract the value from:
  - source or target structures (Entity_SOURCE, Entity_TARGET)
  - environment variables (Entity_ENV)

For example:
  - In path "source:name", Entity_SOURCE indicates that the name field should be taken from the source structure
  - In path "target:owner", Entity_TARGET indicates that the owner field should be taken from the target structure
  - In path "env:DATABASE_URL", Entity_ENV indicates that the value should be retrieved from environment variable DATABASE_URL

Valid values:
  - Entity_SOURCE ("source") - field belongs to source structure
  - Entity_TARGET ("target") - field belongs to target structure
  - Entity_ENV ("env") - value is retrieved from environment variables
*/
type Entity string

const (
	// Entity_TARGET represents the target structure in policy conditions.
	// Used in paths like "target:field" to indicate that the field should be
	// extracted from the target structure passed to Evaluate method.
	Entity_TARGET Entity = "target"

	// Entity_SOURCE represents the source structure in policy conditions.
	// Used in paths like "source:field" to indicate that the field should be
	// extracted from the source structure passed to Evaluate method.
	Entity_SOURCE Entity = "source"

	// Entity_ENV represents environment variables in policy conditions.
	// Used in paths like "env:VARIABLE_NAME" to indicate that the value should be
	// retrieved from environment variables using os.LookupEnv.
	Entity_ENV Entity = "env"
)

var (
	// AVALIABLE_ENTITIES is a list of all valid entity values.
	// Used internally for validation of entity values in paths.
	// Contains Entity_SOURCE, Entity_TARGET, and Entity_ENV.
	AVALIABLE_ENTITIES = []Entity{
		Entity_SOURCE,
		Entity_TARGET,
		Entity_ENV,
	}
)

/*
IsValid checks if the Entity value is valid.

Valid entity values are Entity_SOURCE, Entity_TARGET, and Entity_ENV.
The function checks if the entity is present in AVALIABLE_ENTITIES list.

Returns:
  - bool - true if value is valid (one of Entity_SOURCE, Entity_TARGET, Entity_ENV), false otherwise
*/
func (e Entity) IsValid() bool {
	return slices.Contains(AVALIABLE_ENTITIES, e)
}
