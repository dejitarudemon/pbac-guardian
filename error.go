package noctisguard

import "fmt"

type ErrDublicateName struct {
	name string
}

func NewErrDuplicateName(name string) error {
	return ErrDublicateName{name: name}
}

func (e ErrDublicateName) Error() string {
	return fmt.Sprintf("%v is already used by another policy", e.name)
}

type ErrExport struct {
	source error
}

func NewErrExport(source error) error {
	return ErrExport{source: source}
}

func (e ErrExport) Unwrap() error {
	return e.source
}

func (e ErrExport) Error() string {
	return fmt.Sprintf("failed to export polices: %v", e.source)
}

type ErrEvaluate struct {
	source error
}

func NewErrEvaluate(source error) error {
	return ErrEvaluate{source: source}
}

func (e ErrEvaluate) Unwrap() error {
	return e.source
}

func (e ErrEvaluate) Error() string {
	return fmt.Sprintf("failed to evaluate: %v", e.source)
}
