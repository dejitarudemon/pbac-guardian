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
		policies := []base.Policy{
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
	"github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

/*
Guardian is the main engine of the library for checking structures against access policies.

Guardian stores policies organized by actions and provides the Evaluate method
to check structures against these policies.

The engine uses L1 cache (Casher interface) to optimize field value retrieval
by caching results of reflect-based field searches. Each evaluation session
generates a unique sessionID that identifies the cache scope for a single
policy application. The cache significantly improves performance when fields
are accessed multiple times (3+ accesses provide benefits, with up to 41%
time savings and 63% fewer allocations in production scenarios).

To create a Guardian instance, use:
  - NewGuardianFromPolices - create from a list of policies passed programmatically
  - NewGuardianFromFile - create from a JSON file with policies

Fields:
  - polices - stores policies organized by actions (action -> []Policy)
  - cash - L1 cache instance for storing field values (can be nil to disable caching)
*/
type Guardian struct {
	// хранит политики, разделенные по действиям (action)
	polices map[string][]base.Policy
	cash    base.Casher
}

/*
NewGuardianFromPolices creates a new Guardian instance from a list of policies passed programmatically.

The function performs policy validation and checks for duplicate names. Policies
are grouped by actions for subsequent fast checking.

The engine uses the provided Casher instance for L1 caching. If nil is passed,
caching will be disabled and field values will be retrieved directly via reflection
on each evaluation. For optimal performance in production scenarios with multiple
policies accessing the same fields, it is recommended to use DefaultCasher.

Parameters:
  - cash - instance of base.Casher for L1 caching (can be nil to disable caching)
  - polices - list of policies to initialize the engine
  - funcConfig - configuration for condition functions (can be nil to use default functions)

Returns:
  - *Guardian - created engine instance, ready to use
  - error - creation error if policies contain duplicate names or are invalid

Possible errors:
  - ErrExport - policy export error. May contain:
  - ErrDuplicateName - if policies with the same names are found
  - validation errors from base (ErrInvalidPath) - if policies are invalid

Example usage:

	import "github.com/dejitarudemon/pbac-guardian/internal/implemented"

	casher := implemented.NewDefaultCasher()
	policies := []base.Policy{
		{
			Name:   "allow-admin",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	engine, err := guardian.NewGuardianFromPolices(casher, policies, nil)
	if err != nil {
		// handle error
	}
*/
func NewGuardianFromPolices(cash base.Casher, polices []base.Policy, funcConfig *base.ConditionFuncsConfig) (*Guardian, error) {
	if funcConfig == nil {
		funcConfig = &implemented.DefaultConditionsFuncs
	}
	mapped, err := export(polices, *funcConfig)
	if err != nil {
		return nil, NewErrExport(err)
	}

	return &Guardian{polices: mapped, cash: cash}, nil
}

/*
NewGuardianFromFile creates a new Guardian instance from a JSON file containing an array of policies.

The function reads the file, parses JSON and creates the engine similar to NewGuardianFromPolices.
The file must contain a valid JSON array of Policy objects.

The engine uses the provided Casher instance for L1 caching. If nil is passed,
caching will be disabled and field values will be retrieved directly via reflection
on each evaluation. For optimal performance in production scenarios with multiple
policies accessing the same fields, it is recommended to use DefaultCasher.

Parameters:
  - cash - instance of base.Casher for L1 caching (can be nil to disable caching)
  - path - path to the file with policies in JSON format
  - funcConfig - configuration for condition functions (can be nil to use default functions)

Returns:
  - *Guardian - created engine instance, ready to use
  - error - creation error if the file is unavailable or contains invalid data

Possible errors:
  - ErrExport - policy export error. May contain:
  - file open/read errors (os.PathError, etc.)
  - JSON parsing errors (json.SyntaxError, etc.)
  - ErrDuplicateName - if policies with the same names are found
  - validation errors from base (ErrInvalidPath) - if policies are invalid

Example usage:

	import "github.com/dejitarudemon/pbac-guardian/internal/implemented"

	// File policies.json:
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
	engine, err := guardian.NewGuardianFromFile(casher, "policies.json", nil)
	if err != nil {
		// handle error
	}
*/
func NewGuardianFromFile(cash base.Casher, path string, funcConfig *base.ConditionFuncsConfig) (*Guardian, error) {
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

	return NewGuardianFromPolices(cash, polices, funcConfig)
}

/*
Evaluate checks whether source and target structures match policies for the specified action.

The function finds all policies related to the specified action and checks their
conditions. The result is determined by the effect logic: policies with DENY effect
have priority and deny the action if their conditions are not met. Policies with
ALLOW effect allow the action if at least one of them passes the check.

The function supports cancellation through context.Context, allowing to interrupt
long-running condition checking operations.

Each evaluation session generates a unique sessionID that identifies the cache scope
for a single policy application. This sessionID is used to cache field values retrieved
via reflection, avoiding repeated field searches within the same evaluation session.
The cache is cleared after the evaluation completes.

Parameters:
  - ctx - context for operation cancellation and timeout control
  - source - first structure to check (usually the action source)
  - target - second structure to check (usually the action target)
  - action - action in format "entity:action:extra..." for which policies are checked

Returns:
  - bool - check result:
  - true - action is allowed (at least one ALLOW policy passed the check)
  - false - action is denied (no policies for the action, or DENY policy did not pass the check)
  - error - execution error if a problem occurred during condition evaluation

Possible errors:
  - ErrEvaluate - policy evaluation error. May contain errors from base:
  - ErrCancelled - operation was cancelled through context.Context
  - ErrInvalidPath - path parsing error or field not found
  - ErrInvalidType - invalid structure or field type
  - ErrUncomparable - cannot compare values in condition
  - ErrInexpectedBehavior - internal error (condition function not found)

Logic:
 1. Generate unique sessionID for this evaluation session
 2. If there are no policies for the specified action, returns (false, nil)
 3. For each policy, conditions are checked:
    - If context is cancelled, returns (false, ErrCancelled)
    - Field values are retrieved using get() which checks cache first (if cash is not nil)
    - If policy has DENY effect and conditions are not met, returns (false, nil)
    - If policy has ALLOW effect and conditions are met, sets allowed = true flag
 4. Returns result: (allowed, nil) or (false, error) on error

Example usage:

	import (
		"context"
		"time"
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
			// other error
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
		return false, NewErrEvaluate(fmt.Errorf("noctis engine is nil"))
	}
	if ctx == nil {
		return false, NewErrEvaluate(fmt.Errorf("context is nil"))
	}

	polices, ok := n.polices[action]
	if !ok {
		return false, nil
	}

	allowed := false
	sessionID := generateNewSesstionID()

	for _, policy := range polices {
		ok, err := policy.Evaluate(ctx, source, target, action, n.cash, sessionID)
		if err != nil {
			return false, NewErrEvaluate(err)
		}

		if policy.Effect == base.Effect_DENY {
			if ok {
				return false, err
			}
		} else {
			allowed = allowed || ok
		}
	}

	return allowed, nil
}
