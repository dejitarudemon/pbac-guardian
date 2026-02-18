package noctisguard

import "fmt"

/*
ErrDuplicateName represents an error that occurs when trying to create a policy
with a name that is already used by another policy.

Each policy must have a unique name. When trying to add a policy with an
already existing name, this error is returned.
*/
type ErrDuplicateName struct {
	name string
}

/*
NewErrDuplicateName creates a new ErrDuplicateName error.

Parameters:
  - name - policy name that is already used by another policy

Returns:
  - error - created ErrDuplicateName error
*/
func NewErrDuplicateName(name string) error {
	return ErrDuplicateName{name: name}
}

func (e ErrDuplicateName) Error() string {
	return fmt.Sprintf("%v is already used by another policy", e.name)
}

/*
ErrExport represents an error that occurs when exporting policies to the Noctis engine.

The error wraps the original error (source), allowing the use of errors.Unwrap()
to get problem details. Used when creating the engine from policies or a file.
*/
type ErrExport struct {
	source error
}

/*
NewErrExport creates a new ErrExport error wrapping the original error.

Parameters:
  - source - original error that led to the export error (may be nil)

Returns:
  - error - created ErrExport error that can be unwrapped via errors.Unwrap()
*/
func NewErrExport(source error) error {
	return ErrExport{source: source}
}

func (e ErrExport) Unwrap() error {
	return e.source
}

func (e ErrExport) Error() string {
	return fmt.Sprintf("failed to export polices: %v", e.source)
}

/*
ErrEvaluate represents an error that occurs when evaluating policies during
the execution of the Evaluate method.

The error wraps the original error (source), allowing the use of errors.Unwrap()
to get problem details. Occurs when checking policy conditions or accessing structure fields.
*/
type ErrEvaluate struct {
	source error
}

/*
NewErrEvaluate creates a new ErrEvaluate error wrapping the original error.

Parameters:
  - source - original error that led to the evaluation error (may be nil)

Returns:
  - error - created ErrEvaluate error that can be unwrapped via errors.Unwrap()
*/
func NewErrEvaluate(source error) error {
	return ErrEvaluate{source: source}
}

func (e ErrEvaluate) Unwrap() error {
	return e.source
}

func (e ErrEvaluate) Error() string {
	return fmt.Sprintf("failed to evaluate: %v", e.source)
}
