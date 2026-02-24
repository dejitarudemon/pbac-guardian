package base

/*
ConditionFuncsConfig provides configuration for condition functions used in policy evaluation.

This structure allows customizing the behavior of condition checks (Contains, Eq, Neq, Lt)
by providing custom implementations. If nil is passed to NewGuardianFromPolices or
NewGuardianFromFile, the default implementations from implemented.DefaultConditionsFuncs
will be used.

Fields:
  - Contains - function for checking if a value is in a list
  - Eq - function for equality comparison
  - Neq - function for inequality comparison
  - Lt - function for less-than comparison
*/
type ConditionFuncsConfig struct {
	Contains ConditionFunc
	Eq       ConditionFunc
	Neq      ConditionFunc
	Lt       ConditionFunc
}

/*
Select returns a condition function by its name.

The function is used internally to retrieve the appropriate condition function
based on the condition type specified in the policy.

Parameters:
  - key - condition function name ("Contains", "Eq", "Neq", or "Lt")

Returns:
  - ConditionFunc - condition function if found, nil otherwise
*/
func (c ConditionFuncsConfig) Select(key string) ConditionFunc {
	switch key {
	case "Contains":
		return c.Contains
	case "Eq":
		return c.Eq
	case "Neq":
		return c.Neq
	case "Lt":
		return c.Lt
	}

	return nil
}
