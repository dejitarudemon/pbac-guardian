/*
Package cashing provides interfaces and implementations for L1 caching in the policy evaluation engine.

This file contains CashTree implementation for tracking field access counts and disabling
cache for rarely accessed fields to optimize memory usage.
*/
package cashing

// DefaultCashDisableThreShold is the default threshold for disabling cache.
// Fields accessed less than this number of times will not be cached.
var DefaultCashDisableThreShold = 3

type cashTreeNode struct {
	IsDisabled bool
	Count      int
}

type cashTreeNodes map[string]*cashTreeNode

/*
CashTree tracks field access counts per action and disables caching for fields
that are accessed less frequently than the configured threshold.

The tree structure allows efficient lookup of whether a field should be cached
for a specific action. This optimization reduces memory usage by avoiding
cache storage for fields that are rarely accessed.
*/
type CashTree struct {
	nodes                map[string]cashTreeNodes
	cashDisableThreShold int
}

/*
NewCashTree creates a new CashTree instance with the specified threshold.

If cashDisableThreShold is less than 1, it will be set to 1 automatically.

Parameters:
  - cashDisableThreShold - threshold for disabling cache (fields accessed less than this will not be cached)

Returns:
  - CashTree - new cash tree instance
*/
func NewCashTree(cashDisableThreShold int) CashTree {
	if cashDisableThreShold < 1 {
		cashDisableThreShold = 1
	}

	return CashTree{cashDisableThreShold: cashDisableThreShold, nodes: make(map[string]cashTreeNodes)}
}

/*
Add records a field access for the specified action and path.

The function increments the access count for the field and updates the disabled
status based on whether the count has reached the threshold. Fields are initially
disabled until they reach the threshold.

Parameters:
  - action - action identifier (e.g., "user:read:document")
  - path - field path (e.g., "source:role" or "target:owner")
*/
func (tree *CashTree) Add(action, path string) {
	if _, ok := tree.nodes[action]; !ok {
		tree.nodes[action] = make(cashTreeNodes, 1)
	}

	if node, ok := tree.nodes[action][path]; ok {
		node.Count += 1
		node.IsDisabled = node.Count >= tree.cashDisableThreShold
	} else {
		tree.nodes[action][path] = &cashTreeNode{
			IsDisabled: true,
			Count:      1,
		}
	}
}

/*
IsDisabled checks if caching is disabled for the specified action and path.

Returns true if the field has not been accessed enough times to reach the threshold,
or if the field has not been recorded in the tree yet. Returns false if the field
has been accessed at least threshold times and caching should be enabled.

Parameters:
  - action - action identifier (e.g., "user:read:document")
  - path - field path (e.g., "source:role" or "target:owner")

Returns:
  - bool - true if caching is disabled for this field, false if caching should be enabled
*/
func (tree *CashTree) IsDisabled(action, path string) bool {
	if nodes, ok := tree.nodes[action]; ok {
		if node, ok := nodes[path]; ok {
			return node.IsDisabled
		}
	}

	return true
}
