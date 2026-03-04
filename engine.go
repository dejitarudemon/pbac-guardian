/*
Package guardian provides an engine for checking structures against policies.

The library allows you to define access policies and check structures against
these policies using a flexible system of conditions and effects.

The library supports L1 caching mechanism to avoid re-searching struct fields
by reflect package. Each evaluation session uses a unique sessionID that identifies
the cache scope for a single policy application. The cache can significantly improve
performance when the same fields are accessed multiple times within an evaluation
session (typically 3+ accesses provide benefits).

Example usage:

	package main

	import (
		"context"
		"fmt"
		"github.com/dejitarudemon/pbac-guardian"
		"github.com/dejitarudemon/pbac-guardian/internal/base"
		"github.com/dejitarudemon/pbac-guardian/internal/implemented"
	)

	type User struct {
		Name string `pbac-guardian:"name"`
		Role string `pbac-guardian:"role"`
	}

	type Document struct {
		Owner string `pbac-guardian:"owner"`
		Type  string `pbac-guardian:"type"`
	}

	func main() {
		// Create policies
		policies := []base.RawPolicy{
			{
				Name:   "admin-read",
				Action: "user:read:document",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Eq: "admin",
					},
				},
			},
			{
				Name:   "owner-read",
				Action: "user:read:document",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:name": {
						Eq: "target:owner",
					},
				},
			},
		}

		// Create cache instance (can be nil to disable caching)
		casher := implemented.NewDefaultCasher()

		// Create engine
		engine, err := guardian.NewGuardianFromPolices(casher, policies, nil)
		if err != nil {
			panic(err)
		}

		// Create context
		ctx := context.Background()

		// Check access
		user := User{Name: "alice", Role: "user"}
		doc := Document{Owner: "alice", Type: "private"}

		allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
		if err != nil {
			if err == base.ErrCancelled {
				fmt.Println("Operation cancelled")
			} else {
				panic(err)
			}
		}

		fmt.Printf("Access allowed: %v\n", allowed) // true (owner-read policy passed)
	}
*/
package guardian

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dejitarudemon/pbac-guardian/internal/base"
	"github.com/dejitarudemon/pbac-guardian/internal/cashing"
	"github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

/*
Guardian is the main engine of the library for checking structures against access policies.

Guardian stores policies organized by actions and effects, providing efficient
lookup and evaluation of access control rules. The engine uses L1 caching to
optimize field access during policy evaluation, significantly improving performance
when the same fields are accessed multiple times within an evaluation session.

To create a Guardian instance, use:
  - NewGuardianFromPolices - create from a list of policies passed programmatically
  - NewGuardianFromFile - create from a JSON file with policies

The engine supports both ALLOW and DENY effects, with DENY policies taking
precedence over ALLOW policies. Policies are grouped by action for fast lookup
during evaluation.

Fields:
  - polices - stores policy pointers organized by actions and effects
    (action -> effect -> []*Policy)
*/
type Guardian struct {
	polices map[string]map[base.Effect][]*base.Policy
}

/*
NewGuardianFromPolices creates a new Guardian instance from a list of policies passed programmatically.

The function performs comprehensive policy validation and checks for duplicate names.
Policies are grouped by actions and effects for subsequent fast lookup during evaluation.
Each policy is validated to ensure it contains valid paths, conditions, and effect values.

The engine uses the provided Casher instance for L1 caching to optimize field access
during policy evaluation. If nil is passed, caching will be disabled and field values
will be retrieved directly via reflection on each evaluation. For optimal performance
in production scenarios with multiple policies accessing the same fields, it is
recommended to use DefaultCasher. The cache becomes beneficial when the same field
is accessed 3+ times within a single evaluation session.

The config parameter allows customizing condition functions and cache behavior.
If config.ConditionsMap is nil, default condition functions (Eq, Neq, In, Lt, Gt, Le, Ge)
will be used. The cache tree tracks field access counts and can automatically disable
caching for fields that are accessed less than the threshold number of times, optimizing
memory usage.

Parameters:
  - cash - instance of cashing.Casher for L1 caching (can be nil to disable caching)
  - polices - list of raw policies (RawPolicy) to initialize the engine
  - config - configuration containing condition functions map and cache disable threshold

Returns:
  - *Guardian - created engine instance, ready to use
  - error - creation error if policies contain duplicate names or are invalid

Possible errors:
  - ErrExport - policy export error wrapping the original error. May contain:
  - ErrDuplicateName - if policies with the same names are found
  - validation errors from base (ErrInvalidPath, ErrInvalidType) - if policies are invalid
  - file I/O errors - if reading from file fails (for NewGuardianFromFile)
  - JSON parsing errors - if JSON is malformed (for NewGuardianFromFile)

Example usage:

	import (
		"github.com/dejitarudemon/pbac-guardian"
		"github.com/dejitarudemon/pbac-guardian/internal/implemented"
		"github.com/dejitarudemon/pbac-guardian/internal/base"
	)

	casher := implemented.NewDefaultCasher()
	policies := []base.RawPolicy{
		{
			Name:   "allow-admin",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	config := base.Config{
		ConditionsMap:        nil, // use defaults
		CashDisableThreShold: 3,   // disable cache for fields accessed less than 3 times
	}

	engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
	if err != nil {
		// handle error
	}
*/
func NewGuardianFromPolices(cash cashing.Casher, polices []base.RawPolicy, config base.Config) (*Guardian, error) {
	if config.ConditionsMap == nil {
		config.ConditionsMap = &implemented.DefaultConditionsMap
	}
	cashTree := cashing.NewCashTree(config.CashDisableThreShold)

	mapped, err := export(polices, cash, config, &cashTree)
	if err != nil {
		return nil, NewErrExport(err)
	}

	return &Guardian{polices: mapped}, nil
}

/*
NewGuardianFromFile creates a new Guardian instance from a JSON file containing an array of policies.

The function reads the file, parses JSON and creates the engine similar to NewGuardianFromPolices.
The file must contain a valid JSON array of RawPolicy objects. The file is opened in read-only
mode and read completely into memory before parsing.

The engine uses the provided Casher instance for L1 caching. If nil is passed,
caching will be disabled and field values will be retrieved directly via reflection
on each evaluation. For optimal performance in production scenarios with multiple
policies accessing the same fields, it is recommended to use DefaultCasher.

The config parameter allows customizing condition functions and cache behavior.
If config.ConditionsMap is nil, default condition functions (Eq, Neq, In, Lt, Gt, Le, Ge)
will be used. The cache tree tracks field access counts and can automatically disable
caching for fields that are accessed less than the threshold number of times.

Parameters:
  - cash - instance of cashing.Casher for L1 caching (can be nil to disable caching)
  - path - path to the file with policies in JSON format
  - config - configuration containing condition functions map and cache disable threshold

Returns:
  - *Guardian - created engine instance, ready to use
  - error - creation error if the file is unavailable or contains invalid data

Possible errors:
  - ErrExport - policy export error wrapping the original error. May contain:
  - file open/read errors (os.PathError, etc.) - if the file cannot be opened or read
  - JSON parsing errors (json.SyntaxError, etc.) - if the JSON is malformed
  - ErrDuplicateName - if policies with the same names are found
  - validation errors from base (ErrInvalidPath, ErrInvalidType) - if policies are invalid

Example usage:

	import (
		"github.com/dejitarudemon/pbac-guardian"
		"github.com/dejitarudemon/pbac-guardian/internal/implemented"
		"github.com/dejitarudemon/pbac-guardian/internal/base"
	)

	// File policies.json contains an array of RawPolicy structures:
	// [
	//   {
	//     "name": "allow-admin",
	//     "action": "user:read",
	//     "effect": "allow",
	//     "conditions": {
	//       "source:role": {"eq": "admin"}
	//     }
	//   }
	// ]

	casher := implemented.NewDefaultCasher()
	config := base.Config{
		ConditionsMap:        nil, // use defaults
		CashDisableThreShold: 3,
	}
	engine, err := guardian.NewGuardianFromFile(casher, "policies.json", config)
	if err != nil {
		// handle error
	}
*/
func NewGuardianFromFile(cash cashing.Casher, path string, config base.Config) (*Guardian, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return nil, NewErrExport(err)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, NewErrExport(err)
	}

	var polices []base.RawPolicy

	if err := json.Unmarshal(content, &polices); err != nil {
		return nil, NewErrExport(err)
	}

	return NewGuardianFromPolices(cash, polices, config)
}

/*
Evaluate checks whether source and target structures match policies for the specified action.

The function finds all policies related to the specified action and checks their
conditions. The result is determined by the effect logic: policies with DENY effect
have priority and deny the action if their conditions are met. Policies with
ALLOW effect allow the action if at least one of them passes the check.

The evaluation process follows this order:
 1. First, all DENY policies are evaluated. If any DENY policy's conditions are met,
    the action is immediately denied (returns false, nil).
 2. Then, all ALLOW policies are evaluated. If at least one ALLOW policy's conditions
    are met, the action is allowed (returns true, nil).
 3. If no policies match or no ALLOW policies pass, the action is denied (returns false, nil).

The function supports cancellation through context.Context, allowing to interrupt
long-running condition checking operations. If the context is cancelled during
evaluation, the function returns (false, ErrCancelled).

Each evaluation session generates a unique sessionID that identifies the cache scope
for a single policy application. This sessionID is used to cache field values retrieved
via reflection, avoiding repeated field searches within the same evaluation session.
The cache is automatically cleared after the evaluation completes, ensuring thread-safety
and preventing cache pollution between different evaluations.

The function supports various condition types including:
  - Eq, Neq - equality and inequality checks
  - In - membership checks (value in list)
  - Lt, Gt, Le, Ge - comparison operations for numeric and time values
  - Support for environment variables (env:VARIABLE_NAME)
  - Support for time values (time:now, time:now:1|day, etc.)

Parameters:
  - ctx - context for operation cancellation and timeout control (must not be nil)
  - source - first structure to check (usually the action source, e.g., User)
  - target - second structure to check (usually the action target, e.g., Document)
  - action - action in format "entity:action:extra..." for which policies are checked

Returns:
  - bool - check result:
  - true - action is allowed (at least one ALLOW policy passed the check)
  - false - action is denied (no policies for the action, DENY policy passed, or no ALLOW policies passed)
  - error - execution error if a problem occurred during condition evaluation

Possible errors:
  - ErrEvaluate - policy evaluation error wrapping the original error. May contain errors from base:
  - ErrCancelled - operation was cancelled through context.Context
  - ErrInvalidPath - path parsing error or field not found in structure
  - ErrInvalidType - invalid structure or field type for the condition
  - ErrUncomparable - cannot compare values in condition (incompatible types)
  - ErrInexpectedBehavior - internal error (should not occur in normal operation)

Example usage:

	import (
		"context"
		"time"
		"github.com/dejitarudemon/pbac-guardian"
		"github.com/dejitarudemon/pbac-guardian/internal/base"
	)

	type User struct {
		Name string `pbac-guardian:"name"`
		Role string `pbac-guardian:"role"`
	}

	type Document struct {
		Owner string `pbac-guardian:"owner"`
	}

	user := User{Name: "alice", Role: "admin"}
	doc := Document{Owner: "alice"}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check access to read document
	allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
	if err != nil {
		if err == base.ErrCancelled {
			// operation was cancelled
		} else {
			// other error - unwrap to get details
			// unwrappedErr := errors.Unwrap(err)
		}
	}

	if allowed {
		// access granted
	} else {
		// access denied
	}
*/
func (n *Guardian) Evaluate(ctx context.Context, source, target any, action string) (bool, error) {
	if n == nil {
		return false, NewErrEvaluate(fmt.Errorf("guardian engine is nil"))
	}
	if ctx == nil {
		return false, NewErrEvaluate(fmt.Errorf("context is nil"))
	}

	polices, ok := n.polices[action]
	if !ok {
		return false, nil
	}

	sessionID := generateNewSesstionID()

	for _, policy := range polices[base.Effect_DENY] {
		denied, err := policy.Evaluate(ctx, source, target, action, sessionID)
		if err != nil {
			return false, NewErrEvaluate(err)
		}
		if denied {
			return false, nil
		}
	}

	for _, policy := range polices[base.Effect_ALLOW] {
		allowed, err := policy.Evaluate(ctx, source, target, action, sessionID)
		if err != nil {
			return false, NewErrEvaluate(err)
		}

		if allowed {
			return true, nil
		}
	}

	return false, nil
}
