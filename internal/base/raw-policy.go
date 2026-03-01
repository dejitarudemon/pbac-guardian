package base

/*
RawPolicy represents a raw policy structure used for JSON unmarshaling and policy creation.

This structure contains the public fields that can be serialized/deserialized from JSON.
To create a Policy instance, use NewPolicy function which validates the raw policy
and creates a Policy with internal fields properly initialized.

Fields:
  - Name - unique policy name (used for identification)
  - Action - action in format "entity:action:extra1:extra2..." (e.g., "user:read:profile")
  - Effect - policy effect: Effect_ALLOW (allow) or Effect_DENY (deny)
  - Conditions - map of conditions. Key - path to field in format "source:field", "target:field", "env:VAR_NAME", or "time:now",
    value - condition to check (Contains, Eq, Neq, Lt, Gt)

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
*/
type RawPolicy struct {
	Name       string               `json:"name"`
	Action     string               `json:"action"`
	Effect     Effect               `json:"effect"`
	Conditions map[string]Condition `json:"conditions"`
}
