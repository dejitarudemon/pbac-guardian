/*
Package base provides basic types and functions for working with policies,
conditions, effects, and entities in the access control system.

The package contains definitions of policies, comparison conditions, effects (allow/deny),
entities (source/target) and interfaces for custom comparison.
*/
package base

/*
This file contains interfaces for comparison operations
*/

/*
Comparable interface is intended for custom types that require
special comparison logic in policy conditions.

Implementing the interface allows using structures in Lt, Eq, and Neq conditions
with custom comparison logic instead of standard DeepEqual.

Example usage:

	type User struct {
		Name string
		Age  int
	}

	func (u User) Compare(other any) (int, bool) {
		o, ok := other.(User)
		if !ok {
			return 0, false
		}
		if u.Age < o.Age {
			return -1, true
		}
		if u.Age > o.Age {
			return 1, true
		}
		return 0, true
	}
*/
type Comparable interface {
	/*
		Compare compares the object implementing the interface with another object.

		The method is used in Eq, Neq, and Lt conditions for custom structure comparison.
		If the method returns false, standard comparison via DeepEqual is used.

		Parameters:
			- other - object to compare (can be of any type)

		Returns:
			- int - comparison result:
				* < 0 - current object is less than other
				* = 0 - objects are equal
				* > 0 - current object is greater than other
			- bool - comparison correctness flag:
				* true - comparison completed successfully, result is correct
				* false - comparison impossible (incompatible types), DeepEqual will be used
	*/
	Compare(other any) (int, bool)
}
