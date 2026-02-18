package noctisguard

import (
	"slices"

	"github.com/dejitarudemon/noctis-guard/internal/base"
)

func export(polices []base.Policy) (map[string][]base.Policy, error) {
	mappedPolices := make(map[string][]base.Policy)
	usedNames := make([]string, 0, len(polices))

	for _, policy := range polices {
		if slices.Contains(usedNames, policy.Name) {
			return nil, NewErrDuplicateName(policy.Name)
		}

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
