package implemented

import "github.com/dejitarudemon/pbac-guardian/internal/base"

// DEFAULT_CASH_DISABLE_THRESHOLD is the default threshold for disabling cache
// for rarely accessed fields. Fields accessed less than this number of times
// will not be cached to optimize memory usage.
const DEFAULT_CASH_DISABLE_THRESHOLD = 3

/*
NewDefaultConfig creates a default configuration for Guardian engine.

The function returns a Config with:
  - DefaultConditionsMap - default implementations of all condition functions (In, Eq, Neq, Lt, Gt, Le, Ge)
  - DEFAULT_CASH_DISABLE_THRESHOLD (3) - cache will be disabled for fields accessed less than 3 times

This is a convenience function for creating a ready-to-use configuration with sensible defaults.
You can use it when you don't need to customize condition functions or cache behavior.

Returns:
  - base.Config - default configuration ready for use in NewGuardianFromPolices or NewGuardianFromFile

Example usage:

	config := implemented.NewDefaultConfig()
	engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
	if err != nil {
		// handle error
	}
*/
func NewDefaultConfig() base.Config {
	return base.Config{
		ConditionsMap:        &DefaultConditionsMap,
		CashDisableThreShold: DEFAULT_CASH_DISABLE_THRESHOLD,
	}
}
