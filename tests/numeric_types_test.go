/*
Пакет tests содержит тесты для различных числовых типов в условиях сравнения.

Тесты проверяют работу ltPrimitives и ltConditionFunc для всех поддерживаемых
числовых типов (int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64),
что улучшает покрытие функции ltPrimitives с 38.9% до более высокого уровня.
*/
package tests

import (
	"context"
	"testing"

	noctisguard "github.com/dejitarudemon/noctis-guard"
	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
Тестовые структуры для проверки различных числовых типов.
*/

// UserWithInt8 представляет пользователя с полем типа int8
type UserWithInt8 struct {
	Name string `noctis-guard:"name"`
	Age  int8   `noctis-guard:"age"`
}

// UserWithInt16 представляет пользователя с полем типа int16
type UserWithInt16 struct {
	Name string `noctis-guard:"name"`
	Age  int16  `noctis-guard:"age"`
}

// UserWithInt32 представляет пользователя с полем типа int32
type UserWithInt32 struct {
	Name string `noctis-guard:"name"`
	Age  int32  `noctis-guard:"age"`
}

// UserWithInt64 представляет пользователя с полем типа int64
type UserWithInt64 struct {
	Name string `noctis-guard:"name"`
	Age  int64  `noctis-guard:"age"`
}

// UserWithUint8 представляет пользователя с полем типа uint8
type UserWithUint8 struct {
	Name string `noctis-guard:"name"`
	Age  uint8  `noctis-guard:"age"`
}

// UserWithUint16 представляет пользователя с полем типа uint16
type UserWithUint16 struct {
	Name string `noctis-guard:"name"`
	Age  uint16 `noctis-guard:"age"`
}

// UserWithUint32 представляет пользователя с полем типа uint32
type UserWithUint32 struct {
	Name string `noctis-guard:"name"`
	Age  uint32 `noctis-guard:"age"`
}

// UserWithUint64 представляет пользователя с полем типа uint64
type UserWithUint64 struct {
	Name string `noctis-guard:"name"`
	Age  uint64 `noctis-guard:"age"`
}

// UserWithFloat32 представляет пользователя с полем типа float32
type UserWithFloat32 struct {
	Name string  `noctis-guard:"name"`
	Age  float32 `noctis-guard:"age"`
}

// UserWithFloat64 представляет пользователя с полем типа float64
type UserWithFloat64 struct {
	Name string  `noctis-guard:"name"`
	Age  float64 `noctis-guard:"age"`
}

/*
TestLtNumericTypes тестирует условие Lt для различных числовых типов.

Тест проверяет работу ltPrimitives для всех поддерживаемых числовых типов,
что значительно улучшает покрытие функции ltPrimitives.
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
			// Тест проверяет условие Lt для типа int8
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
			// Тест проверяет условие Lt для типа int16
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
			// Тест проверяет условие Lt для типа int32
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
			// Тест проверяет условие Lt для типа int64
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
			// Тест проверяет условие Lt для типа uint8
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
			// Тест проверяет условие Lt для типа uint16
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
			// Тест проверяет условие Lt для типа uint32
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
			// Тест проверяет условие Lt для типа uint64
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
			// Тест проверяет условие Lt для типа float32
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
			// Тест проверяет условие Lt для типа float64
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
			// Тест проверяет, что условие Lt возвращает false при равенстве для float64
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
			// Тест проверяет, что условие Lt возвращает false при большем значении для uint64
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
			engine, err := noctisguard.NewNoctisFromPolices([]base.Policy{tt.policy})
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

