/*
Package noctisguard provides an engine for checking structures against policies.

The library allows you to define access policies and check structures against
these policies using a flexible system of conditions and effects.

Example usage:

	package main

	import (
		"context"
		"fmt"
		"github.com/dejitarudemon/noctis-guard"
		"github.com/dejitarudemon/noctis-guard/internal/base"
	)

	type User struct {
		Name string `noctis-guard:"name"`
		Role string `noctis-guard:"role"`
	}

	type Document struct {
		Owner string `noctis-guard:"owner"`
		Type  string `noctis-guard:"type"`
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

		// Create engine
		engine, err := noctisguard.NewNoctisFromPolices(policies)
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
package noctisguard

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
Noctis is the main engine of the library for checking structures against access policies.

Noctis stores policies organized by actions and provides the Evaluate method
to check structures against these policies.

To create a Noctis instance, use:
  - NewNoctisFromPolices - create from a list of policies passed programmatically
  - NewNoctisFromFile - create from a JSON file with policies
*/
type Noctis struct {
	// хранит политики, разделенные по действиям (action)
	polices map[string][]base.Policy
}

/*
NewNoctisFromPolices creates a new Noctis instance from a list of policies passed programmatically.

The function performs policy validation and checks for duplicate names. Policies
are grouped by actions for subsequent fast checking.

Parameters:
  - polices - list of policies to initialize the engine

Returns:
  - *Noctis - created engine instance, ready to use
  - error - creation error if policies contain duplicate names or are invalid

Possible errors:
  - ErrExport - policy export error. May contain:
  - ErrDuplicateName - if policies with the same names are found
  - validation errors from base (ErrInvalidPath) - if policies are invalid

Example usage:

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

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		// handle error
	}
*/
func NewNoctisFromPolices(polices []base.Policy) (*Noctis, error) {
	mapped, err := export(polices)
	if err != nil {
		return nil, NewErrExport(err)
	}

	return &Noctis{polices: mapped}, nil
}

/*
NewNoctisFromFile creates a new Noctis instance from a JSON file containing an array of policies.

The function reads the file, parses JSON and creates the engine similar to NewNoctisFromPolices.
The file must contain a valid JSON array of Policy objects.

Parameters:
  - path - path to the file with policies in JSON format

Returns:
  - *Noctis - created engine instance, ready to use
  - error - creation error if the file is unavailable or contains invalid data

Possible errors:
  - ErrExport - policy export error. May contain:
  - file open/read errors (os.PathError, etc.)
  - JSON parsing errors (json.SyntaxError, etc.)
  - ErrDuplicateName - if policies with the same names are found
  - validation errors from base (ErrInvalidPath) - if policies are invalid

Example usage:

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

	engine, err := noctisguard.NewNoctisFromFile("policies.json")
	if err != nil {
		// handle error
	}
*/
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

/*
Evaluate checks whether source and target structures match policies for the specified action.

The function finds all policies related to the specified action and checks their
conditions. The result is determined by the effect logic: policies with DENY effect
have priority and deny the action if their conditions are not met. Policies with
ALLOW effect allow the action if at least one of them passes the check.

The function supports cancellation through context.Context, allowing to interrupt
long-running condition checking operations.

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
 1. If there are no policies for the specified action, returns (false, nil)
 2. For each policy, conditions are checked:
    - If context is cancelled, returns (false, ErrCancelled)
    - If policy has DENY effect and conditions are not met, returns (false, nil)
    - If policy has ALLOW effect and conditions are met, sets allowed = true flag
 3. Returns result: (allowed, nil) or (false, error) on error

Example usage:

	import (
		"context"
		"time"
	)

	type User struct {
		Name string `noctis-guard:"name"`
		Role string `noctis-guard:"role"`
	}

	type Document struct {
		Owner string `noctis-guard:"owner"`
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
func (n *Noctis) Evaluate(ctx context.Context, source, target any, action string) (bool, error) {
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

	for _, policy := range polices {
		ok, err := policy.Evaluate(ctx, source, target, action)
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
