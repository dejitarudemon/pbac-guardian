/*
Пакет tests содержит тесты для публичных методов структуры Policy пакета base.

Тесты проверяют только экспортируемые методы:
  - Evaluate - оценка политики для заданных source, target и action
  - IsValid - валидация политики
*/
package tests

import (
	"context"
	"testing"

	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
Тестовые структуры для проверки работы с политиками.

Эти структуры используются для тестирования методов Policy
с различными типами данных и вложенными структурами.
*/

// PolicyTestUser представляет пользователя для тестирования политик
type PolicyTestUser struct {
	Name string `noctis-guard:"name"` // Имя пользователя
	Role string `noctis-guard:"role"` // Роль пользователя
	Age  int    `noctis-guard:"age"`  // Возраст пользователя
}

// PolicyTestDocument представляет документ для тестирования политик
type PolicyTestDocument struct {
	Owner string   `noctis-guard:"owner"` // Владелец документа
	Type  string   `noctis-guard:"type"`  // Тип документа
	Tags  []string `noctis-guard:"tags"`  // Теги документа
}

/*
TestPolicyEvaluate тестирует публичный метод Evaluate для оценки политики.

Тест проверяет:
  - Соответствие действия политики переданному action
  - Проверку условий политики
  - Обработку нескольких условий (логическое И)
  - Сравнение полей из разных структур (source и target)
  - Возврат корректного результата при различных сценариях
*/
func TestPolicyEvaluate(t *testing.T) {
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
			name: "matching action and condition",
			// Тест проверяет успешную оценку политики при совпадении действия и выполнении условия
			policy: base.Policy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
				},
			},
			source:  PolicyTestUser{Name: "admin", Role: "admin"},
			target:  PolicyTestDocument{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "non-matching action",
			// Тест проверяет, что политика не применяется при несовпадении действия
			policy: base.Policy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
				},
			},
			source:  PolicyTestUser{Name: "admin", Role: "admin"},
			target:  PolicyTestDocument{Owner: "user", Type: "public"},
			action:  "user:write", // Другое действие
			want:    false,
			wantErr: false,
		},
		{
			name: "non-matching condition",
			// Тест проверяет, что политика не проходит при невыполнении условия
			policy: base.Policy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
				},
			},
			source:  PolicyTestUser{Name: "user", Role: "user"}, // Роль не "admin"
			target:  PolicyTestDocument{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "multiple conditions - all match",
			// Тест проверяет, что все условия должны быть выполнены (логическое И)
			policy: base.Policy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
					"source:age":  {Lt: 100}, // Возраст меньше 100
				},
			},
			source:  PolicyTestUser{Name: "admin", Role: "admin", Age: 25}, // Оба условия выполнены
			target:  PolicyTestDocument{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
		{
			name: "multiple conditions - one fails",
			// Тест проверяет, что при невыполнении хотя бы одного условия политика не проходит
			policy: base.Policy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
					"source:age":  {Lt: 10}, // Возраст должен быть меньше 10
				},
			},
			source:  PolicyTestUser{Name: "admin", Role: "admin", Age: 25}, // Возраст 25, условие не выполнено
			target:  PolicyTestDocument{Owner: "user", Type: "public"},
			action:  "user:read",
			want:    false,
			wantErr: false,
		},
		{
			name: "compare fields from different structures",
			// Тест проверяет сравнение полей из разных структур (source и target)
			policy: base.Policy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:name": {Eq: "target:owner"}, // Сравнение name из source с owner из target
				},
			},
			source:  PolicyTestUser{Name: "alice", Role: "user"},
			target:  PolicyTestDocument{Owner: "alice", Type: "public"}, // Имена совпадают
			action:  "user:read",
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.policy.Evaluate(ctx, tt.source, tt.target, tt.action)
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
TestPolicyIsValid тестирует публичный метод IsValid для валидации политики.

Тест проверяет:
  - Валидацию формата действия (минимум 2 части через ":")
  - Проверку отсутствия пустых частей в действии
  - Валидацию всех путей в условиях
  - Корректность обработки валидных и невалидных политик
*/
func TestPolicyIsValid(t *testing.T) {
	tests := []struct {
		name    string
		policy  base.Policy
		wantErr bool
	}{
		{
			name: "valid policy",
			// Тест проверяет валидную политику со всеми корректными параметрами
			policy: base.Policy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid action - too short",
			// Тест проверяет валидацию действия с недостаточным количеством частей
			policy: base.Policy{
				Name:   "test",
				Action: "read", // Только одна часть, нужно минимум 2
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid action - empty part",
			// Тест проверяет валидацию действия с пустой частью
			policy: base.Policy{
				Name:   "test",
				Action: "user::read", // Пустая часть между разделителями
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {Eq: "admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid path in conditions",
			// Тест проверяет валидацию путей в условиях политики
			policy: base.Policy{
				Name:   "test",
				Action: "user:read",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"invalid": {Eq: "admin"}, // Невалидный путь без разделителя ":"
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policy.IsValid()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
