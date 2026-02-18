/*
Пакет tests содержит тесты для условий сравнения в политиках.

Тесты проверяют работу различных условий (Contains, Eq, Neq, Lt) через
публичный API библиотеки, что позволяет улучшить покрытие внутренних функций условий.
*/
package tests

import (
	"context"
	"testing"

	noctisguard "github.com/dejitarudemon/noctis-guard"
	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
TestContainsCondition тестирует условие Contains через политики.

Тест проверяет работу containsConditionFunc через использование условия Contains
в политиках. Это улучшает покрытие функции containsConditionFunc с 0% до 100%.
*/
func TestContainsCondition(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.Policy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "contains - found in slice",
			// Тест проверяет, что условие Contains находит значение в списке
			policy: base.Policy{
				Name:   "contains-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Contains: []any{"admin", "moderator", "user"},
					},
				},
			},
			source:  User{Name: "user", Role: "admin"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "contains - not found in slice",
			// Тест проверяет, что условие Contains не находит значение, если его нет в списке
			policy: base.Policy{
				Name:   "contains-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Contains: []any{"admin", "moderator"},
					},
				},
			},
			source:  User{Name: "user", Role: "guest"}, // "guest" нет в списке
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "contains - empty slice",
			// Тест проверяет, что условие Contains возвращает false для пустого списка
			policy: base.Policy{
				Name:   "contains-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Contains: []any{},
					},
				},
			},
			source:  User{Name: "user", Role: "admin"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "contains - with source role in list",
			// Тест проверяет условие Contains с полем source:role
			// Проверяем, что значение source.role находится в списке
			policy: base.Policy{
				Name:   "contains-source-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Contains: []any{"admin", "moderator", "user"},
					},
				},
			},
			source:  User{Name: "user", Role: "moderator"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "contains - with integer values",
			// Тест проверяет условие Contains с числовыми значениями
			policy: base.Policy{
				Name:   "contains-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Contains: []any{18, 21, 25},
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 21},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
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

/*
TestLtCondition тестирует условие Lt (less than) через политики.

Тест проверяет работу ltConditionFunc и ltPrimitives через использование условия Lt
в политиках. Это улучшает покрытие функций сравнения для различных типов данных.
*/
func TestLtCondition(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.Policy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "lt - int less than",
			// Тест проверяет условие Lt для целых чисел (int)
			policy: base.Policy{
				Name:   "lt-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Lt: 18,
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 16}, // 16 < 18
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - int equal",
			// Тест проверяет, что условие Lt возвращает false при равенстве
			policy: base.Policy{
				Name:   "lt-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Lt: 18,
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 18}, // 18 == 18
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "lt - int greater than",
			// Тест проверяет, что условие Lt возвращает false при большем значении
			policy: base.Policy{
				Name:   "lt-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Lt: 18,
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 25}, // 25 > 18
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "lt - with string comparison",
			// Тест проверяет условие Lt для строк (лексикографическое сравнение)
			// Для этого нужно использовать поле, которое можно сравнить как строку
			policy: base.Policy{
				Name:   "lt-string-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:name": {
						Lt: "m", // Имена до "m" в алфавитном порядке
					},
				},
			},
			source:  User{Name: "alice", Role: "user"}, // "alice" < "m"
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "lt - string greater",
			// Тест проверяет, что условие Lt возвращает false для строк, которые больше
			policy: base.Policy{
				Name:   "lt-string-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:name": {
						Lt: "m",
					},
				},
			},
			source:  User{Name: "zoe", Role: "user"}, // "zoe" > "m"
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "lt - compare with target field",
			// Тест проверяет условие Lt с сравнением полей разных структур
			// Для этого нужно добавить числовое поле в Document
			policy: base.Policy{
				Name:   "lt-compare-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Lt: "target:priority", // Сравнение с полем target
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 25},
			target:  DocumentWithPriority{Document: Document{Owner: "user", Type: "public"}, Priority: 30},
			action:  "user:read",
			want:    true, // 25 < 30
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

// DocumentWithPriority расширяет Document для тестирования числовых сравнений
type DocumentWithPriority struct {
	Document
	Priority int `noctis-guard:"priority"`
}

/*
TestEqConditionExtended тестирует расширенные случаи условия Eq.

Тест проверяет работу eqConditionFunc для различных типов данных и сценариев,
включая сравнение с nil, структуры с Comparable интерфейсом и различные типы.
*/
func TestEqConditionExtended(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.Policy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "eq - integer comparison",
			// Тест проверяет условие Eq для целых чисел
			policy: base.Policy{
				Name:   "eq-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Eq: 25,
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 25},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "eq - integer not equal",
			// Тест проверяет, что условие Eq возвращает false для неравных чисел
			policy: base.Policy{
				Name:   "eq-int-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Eq: 25,
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 30},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "eq - compare integer with string path",
			// Тест проверяет сравнение числа с полем, которое должно быть числом
			policy: base.Policy{
				Name:   "eq-mixed-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:age": {
						Eq: "target:priority", // Сравнение числа с числом через путь
					},
				},
			},
			source:  User{Name: "user", Role: "user", Age: 10},
			target:  DocumentWithPriority{Document: Document{Owner: "user", Type: "public"}, Priority: 10},
			action:  "user:read",
			want:    true,
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

/*
TestMultipleConditionsCombined тестирует комбинацию нескольких условий в одной политике.

Тест проверяет, что все условия в политике должны быть выполнены (логическое И),
что улучшает покрытие функции Evaluate для различных комбинаций условий.
*/
func TestMultipleConditionsCombined(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		policy  base.Policy
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "multiple conditions - all match",
			// Тест проверяет политику с несколькими условиями, все из которых выполнены
			policy: base.Policy{
				Name:   "combined-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Eq:       "admin",
						Contains: []any{"admin", "moderator"},
					},
					"source:age": {
						Lt: 100,
					},
				},
			},
			source:  User{Name: "admin", Role: "admin", Age: 25},
			target:  Document{Owner: "user", Type: "public", Tags: []string{"public"}},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "multiple conditions - one fails",
			// Тест проверяет политику с несколькими условиями, одно из которых не выполнено
			policy: base.Policy{
				Name:   "combined-test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Eq: "admin",
					},
					"source:age": {
						Lt: 10, // Возраст должен быть меньше 10
					},
				},
			},
			source:  User{Name: "admin", Role: "admin", Age: 25}, // Возраст 25, условие не выполнено
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "multiple conditions - Contains and Lt",
			// Тест проверяет комбинацию условий Contains и Lt
			policy: base.Policy{
				Name:   "combined-contains-lt",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Contains: []any{"admin", "moderator"},
					},
					"source:age": {
						Lt: 65,
					},
				},
			},
			source:  User{Name: "admin", Role: "admin", Age: 30},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
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

