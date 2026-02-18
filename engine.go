package noctisguard

import (
	"encoding/json"
	"io"
	"os"

	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
В этом файле представлен движок для работы всей библиотеки
*/

type Noctis struct {
	polices map[string][]base.Policy
}

func NewNoctisFromPolices(polices []base.Policy) (*Noctis, error) {
	mapped, err := export(polices)
	if err != nil {
		return nil, NewErrExport(err)
	}

	return &Noctis{polices: mapped}, nil
}

func NewNoctisFromFile(path string) (*Noctis, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return nil, NewErrExport(err)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, NewErrExport(err)
	}

	var polices []base.Policy

	if err := json.Unmarshal(content, &polices); err != nil {
		return nil, NewErrExport(err)
	}

	return NewNoctisFromPolices(polices)
}

func (n *Noctis) Evaluate(source, target any, action string) (bool, error) {
	polices, ok := n.polices[action]
	if !ok {
		return false, nil
	}

	allowed := false

	for _, policy := range polices {
		ok, err := policy.Evaluate(source, target, action)
		if err != nil {
			return false, NewErrEvaluate(err)
		}

		if policy.Effect == base.Effect_DENY {
			if !ok {
				return false, err
			}
		} else {
			allowed = allowed || ok
		}
	}

	return allowed, nil
}
