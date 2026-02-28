/*
Package tests contains tests for various numeric types in comparison conditions.

Tests check the work of Lt condition for all supported numeric types:
  - Signed integers: int, int8, int16, int32, int64
  - Unsigned integers: uint, uint8, uint16, uint32, uint64
  - Floating point: float32, float64
  - Strings (lexicographic comparison)

The tests verify that numeric comparisons work correctly across different
type sizes and signedness, ensuring type safety and correct comparison logic.
*/
package tests

import (
	"context"
	"testing"

	guardian "github.com/dejitarudemon/pbac-guardian"
	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

/*
Test structures for checking various numeric types.
*/

// UserWithInt8 represents a user with int8 field
type UserWithInt8 struct {
	Name string `pbac-guardian:"name"`
	Age  int8   `pbac-guardian:"age"`
}

// UserWithInt16 represents a user with int16 field
type UserWithInt16 struct {
	Name string `pbac-guardian:"name"`
	Age  int16  `pbac-guardian:"age"`
}

// UserWithInt32 represents a user with int32 field
type UserWithInt32 struct {
	Name string `pbac-guardian:"name"`
	Age  int32  `pbac-guardian:"age"`
}

// UserWithInt64 represents a user with int64 field
type UserWithInt64 struct {
	Name string `pbac-guardian:"name"`
	Age  int64  `pbac-guardian:"age"`
}

// UserWithUint8 represents a user with uint8 field
type UserWithUint8 struct {
	Name string `pbac-guardian:"name"`
	Age  uint8  `pbac-guardian:"age"`
}

// UserWithUint16 represents a user with uint16 field
type UserWithUint16 struct {
	Name string `pbac-guardian:"name"`
	Age  uint16 `pbac-guardian:"age"`
}

// UserWithUint32 represents a user with uint32 field
type UserWithUint32 struct {
	Name string `pbac-guardian:"name"`
	Age  uint32 `pbac-guardian:"age"`
}

// UserWithUint64 represents a user with uint64 field
type UserWithUint64 struct {
	Name string `pbac-guardian:"name"`
	Age  uint64 `pbac-guardian:"age"`
}

// UserWithFloat32 represents a user with float32 field
type UserWithFloat32 struct {
	Name string  `pbac-guardian:"name"`
	Age  float32 `pbac-guardian:"age"`
}

// UserWithFloat64 represents a user with float64 field
type UserWithFloat64 struct {
	Name string  `pbac-guardian:"name"`
	Age  float64 `pbac-guardian:"age"`
}

/*
TestLtNumericTypes tests the Lt condition for various numeric types.

The test checks the work of ltPrimitives for all supported numeric types,
which significantly improves coverage of ltPrimitives function.
*/
func TestLtNumericTypes(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.Policy
		source  any
		target  Document
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "lt - int8",
			// Test checks Lt condition for int8 type
			policy: base.Policy{
				Name:   "lt-int8-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {Lt: int8(20)},
				},
			},
			source:  UserWithInt8{Name: "user", Age: 18},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - int16",
			// Test checks Lt condition for int16 type
			policy: base.Policy{
				Name:   "lt-int16-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {Lt: int16(20)},
				},
			},
			source:  UserWithInt16{Name: "user", Age: 18},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - int32",
			// Test checks Lt condition for int32 type
			policy: base.Policy{
				Name:   "lt-int32-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {Lt: int32(20)},
				},
			},
			source:  UserWithInt32{Name: "user", Age: 18},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - int64",
			// Test checks Lt condition for int64 type
			policy: base.Policy{
				Name:   "lt-int64-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {Lt: int64(20)},
				},
			},
			source:  UserWithInt64{Name: "user", Age: 18},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - uint8",
			// Test checks Lt condition for uint8 type
			policy: base.Policy{
				Name:   "lt-uint8-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {Lt: uint8(20)},
				},
			},
			source:  UserWithUint8{Name: "user", Age: 18},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - uint16",
			// Test checks Lt condition for uint16 type
			policy: base.Policy{
				Name:   "lt-uint16-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {Lt: uint16(20)},
				},
			},
			source:  UserWithUint16{Name: "user", Age: 18},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - uint32",
			// Test checks Lt condition for uint32 type
			policy: base.Policy{
				Name:   "lt-uint32-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {Lt: uint32(20)},
				},
			},
			source:  UserWithUint32{Name: "user", Age: 18},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - uint64",
			// Test checks Lt condition for uint64 type
			policy: base.Policy{
				Name:   "lt-uint64-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {Lt: uint64(20)},
				},
			},
			source:  UserWithUint64{Name: "user", Age: 18},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - float32",
			// Test checks Lt condition for float32 type
			policy: base.Policy{
				Name:   "lt-float32-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {Lt: float32(20.5)},
				},
			},
			source:  UserWithFloat32{Name: "user", Age: 18.5},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - float64",
			// Test checks Lt condition for float64 type
			policy: base.Policy{
				Name:   "lt-float64-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {Lt: 20.5},
				},
			},
			source:  UserWithFloat64{Name: "user", Age: 18.5},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - float64 equal",
			// Test checks that Lt condition returns false on equality for float64
			policy: base.Policy{
				Name:   "lt-float64-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {Lt: 20.5},
				},
			},
			source:  UserWithFloat64{Name: "user", Age: 20.5},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "lt - uint64 greater",
			// Test checks that Lt condition returns false for greater value for uint64
			policy: base.Policy{
				Name:   "lt-uint64-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {Lt: uint64(20)},
				},
			},
			source:  UserWithUint64{Name: "user", Age: 25},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use nil casher for basic functionality tests
			config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
			engine, err := guardian.NewGuardianFromPolices(nil, []base.Policy{tt.policy}, config)
			if err != nil {
				t.Fatalf("failed to create engine: %v", err)
			}

			got, err := engine.Evaluate(ctx, tt.source, tt.target, tt.action)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got != tt.want {
					t.Errorf("Evaluate() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
