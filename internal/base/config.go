package base

type ConditionFuncsConfig struct {
	Contains ConditionFunc
	Eq       ConditionFunc
	Neq      ConditionFunc
	Lt       ConditionFunc
}

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
