package base

/*
ConditionsMap provides configuration for condition functions used in policy evaluation.

This structure allows customizing the behavior of condition checks (Contains, Eq, Neq, Lt, Gt)
by providing custom implementations. If nil is passed in Config.ConditionsMap to
NewGuardianFromPolices or NewGuardianFromFile, the default implementations from
implemented.DefaultConditionsMap will be used.

Fields:
  - Contains - function for checking if a value is in a list
  - Eq - function for equality comparison
  - Neq - function for inequality comparison
  - Lt - function for less-than comparison
  - Gt - function for greater-than comparison
*/
type ConditionsMap struct {
	Contains ConditionFunc
	Eq       ConditionFunc
	Neq      ConditionFunc
	Lt       ConditionFunc
	Gt       ConditionFunc
}

/*
Select returns a condition function by its name.

The function is used internally to retrieve the appropriate condition function
based on the condition type specified in the policy.

Parameters:
  - key - condition function name ("Contains", "Eq", "Neq", "Lt", or "Gt")

Returns:
  - ConditionFunc - condition function if found, nil otherwise
*/
func (c ConditionsMap) Select(key string) ConditionFunc {
	switch key {
	case "Contains":
		return c.Contains
	case "Eq":
		return c.Eq
	case "Neq":
		return c.Neq
	case "Lt":
		return c.Lt
	case "Gt":
		return c.Gt
	}

	return nil
}

/*
Config provides configuration for Guardian engine initialization.

The structure contains settings for condition functions and cache behavior.
It is used as a parameter in NewGuardianFromPolices and NewGuardianFromFile.

Fields:
  - CashDisableThreShold - threshold for disabling cache for rarely accessed fields.
    Fields accessed less than this number of times will not be cached.
    Must be at least 1. If less than 1, it will be set to 1 automatically.
  - ConditionsMap - configuration for condition functions (Contains, Eq, Neq, Lt, Gt).
    If nil, default implementations from implemented.DefaultConditionsMap will be used.
*/
type Config struct {
	CashDisableThreShold int
	ConditionsMap        *ConditionsMap
}
