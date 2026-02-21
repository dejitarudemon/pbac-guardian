package base

import "context"

/*
This file contains interface for engine cashing
*/

/*
Interface Casher is intended for L1-cash to avoid re-searching sruct fields
by reflect package. It must be thread-safe
*/
type Casher interface {
	/*
		Set function is used to set some value under a key in L1-cash

		Parameters:
			- ctx - perfoming context
			- key - any string value providing to value
			- value - struct, primitve or custom type etc

		Returns:
			- error - set value error.
	*/
	Set(ctx context.Context, sessionID, key string, value any) error

	/*
		Get function is used to get some value from L1-cash by a key

		Parameters:
			- ctx - perfoming context
			- key - any string value providing to value

		Returns:
			- any - storaged value if it's exist. If error occured or it's not found, the function must return Nil
			- error - get value error. If value is not found, must be Nil.
	*/
	Get(ctx context.Context, sessionID, key string) (any, error)

	/*
		Clear function is used to clear value by its key

		Parameters:
			- ctx - perfoming context
			- key - any string value providing to value

		Returns:
			- error - clear error. If value is not found, must be Nil.
	*/
	Clear(ctx context.Context, sessionID string) error
}
