/*
Пакет tests содержит бенчмарк тесты для измерения производительности библиотеки.

Бенчмарки проверяют производительность основных операций:
  - Создание движка из политик
  - Оценка политик для различных сценариев
  - Работа с различными типами условий
  - Обработка вложенных структур
*/
package tests

import (
	"context"
	"testing"

	noctisguard "github.com/dejitarudemon/noctis-guard"
	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
BenchmarkNewNoctisFromPolices измеряет производительность создания движка из политик.

Бенчмарк проверяет время создания движка с различным количеством политик.
*/
func BenchmarkNewNoctisFromPolices(b *testing.B) {
	// Создаем набор политик для тестирования
	policies := []base.Policy{
		{
			Name:   "admin-read",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
		{
			Name:   "owner-read",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:name": {Eq: "target:owner"},
			},
		},
		{
			Name:   "deny-private",
			Action: "user:read:document",
			Effect: base.Effect_DENY,
			Conditions: map[string]base.Condition{
				"target:type": {Eq: "private"},
				"source:role": {Neq: "admin"},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := noctisguard.NewNoctisFromPolices(policies)
		if err != nil {
			b.Fatalf("failed to create engine: %v", err)
		}
	}
}

/*
BenchmarkEvaluateSimple измеряет производительность простой оценки политик.

Бенчмарк проверяет время оценки политик с простыми условиями (Eq).
*/
func BenchmarkEvaluateSimple(b *testing.B) {
	policies := []base.Policy{
		{
			Name:   "admin-read",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateMultipleConditions измеряет производительность оценки политик
с несколькими условиями.

Бенчмарк проверяет время оценки политик с комбинацией различных условий.
*/
func BenchmarkEvaluateMultipleConditions(b *testing.B) {
	policies := []base.Policy{
		{
			Name:   "complex-policy",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {
					Eq:       "admin",
					Contains: []any{"admin", "moderator"},
				},
				"source:age": {
					Lt: 100,
				},
				"target:tags": {
					Contains: []any{"public", "shared"},
				},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin", Age: 25}
	target := Document{Owner: "user", Type: "public", Tags: []string{"public"}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateContains измеряет производительность условия Contains.

Бенчмарк проверяет время поиска значения в списке через условие Contains.
*/
func BenchmarkEvaluateContains(b *testing.B) {
	policies := []base.Policy{
		{
			Name:   "contains-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {
					Contains: []any{"admin", "moderator", "user", "guest", "visitor"},
				},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateLt измеряет производительность условия Lt (less than).

Бенчмарк проверяет время сравнения значений через условие Lt.
*/
func BenchmarkEvaluateLt(b *testing.B) {
	policies := []base.Policy{
		{
			Name:   "lt-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:age": {
					Lt: 65,
				},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user", Age: 30}
	target := Document{Owner: "user", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateNestedStructures измеряет производительность работы
с вложенными структурами.

Бенчмарк проверяет время получения значений из вложенных структур.
*/
func BenchmarkEvaluateNestedStructures(b *testing.B) {
	policies := []base.Policy{
		{
			Name:   "nested-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:user:role": {
					Eq: "admin",
				},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := NestedUser{User: User{Name: "admin", Role: "admin"}}
	target := NestedDocument{Doc: Document{Owner: "user", Type: "public", Tags: []string{}}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateMultiplePolicies измеряет производительность оценки
нескольких политик для одного действия.

Бенчмарк проверяет время оценки, когда для действия определено несколько политик.
*/
func BenchmarkEvaluateMultiplePolicies(b *testing.B) {
	policies := []base.Policy{
		{
			Name:   "policy-1",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
		{
			Name:   "policy-2",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:name": {Eq: "target:owner"},
			},
		},
		{
			Name:   "policy-3",
			Action: "user:read:document",
			Effect: base.Effect_DENY,
			Conditions: map[string]base.Condition{
				"target:type": {Eq: "private"},
				"source:role": {Neq: "admin"},
			},
		},
		{
			Name:   "policy-4",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:age": {Lt: 18},
			},
		},
		{
			Name:   "policy-5",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Contains: []any{"moderator", "editor"}},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "admin", Role: "admin", Age: 25}
	target := Document{Owner: "admin", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateFieldComparison измеряет производительность сравнения
полей разных структур.

Бенчмарк проверяет время сравнения полей source и target структур.
*/
func BenchmarkEvaluateFieldComparison(b *testing.B) {
	policies := []base.Policy{
		{
			Name:   "field-comparison",
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:name": {
					Eq: "target:owner",
				},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "alice", Role: "user"}
	target := Document{Owner: "alice", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateDenyPolicy измеряет производительность политик с эффектом DENY.

Бенчмарк проверяет время оценки политик, которые запрещают доступ.
*/
func BenchmarkEvaluateDenyPolicy(b *testing.B) {
	policies := []base.Policy{
		{
			Name:   "deny-policy",
			Action: "user:read:document",
			Effect: base.Effect_DENY,
			Conditions: map[string]base.Condition{
				"target:type": {Eq: "private"},
				"source:role": {Neq: "admin"},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user"}
	target := Document{Owner: "other", Type: "private"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateLargeSlice измеряет производительность условия Contains
с большим списком значений.

Бенчмарк проверяет время поиска в большом списке через условие Contains.
*/
func BenchmarkEvaluateLargeSlice(b *testing.B) {
	// Создаем большой список ролей
	largeRoleList := make([]any, 1000)
	for i := range largeRoleList {
		largeRoleList[i] = "role" + string(rune(i%26+'a'))
	}
	largeRoleList[500] = "admin" // Искомое значение в середине списка

	policies := []base.Policy{
		{
			Name:   "large-slice-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {
					Contains: largeRoleList,
				},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}

/*
BenchmarkEvaluateNoMatch измеряет производительность случая,
когда политики не соответствуют действию.

Бенчмарк проверяет время оценки, когда для действия нет подходящих политик.
*/
func BenchmarkEvaluateNoMatch(b *testing.B) {
	policies := []base.Policy{
		{
			Name:   "other-action",
			Action: "user:write:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		b.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user"}
	target := Document{Owner: "user", Type: "public"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(ctx, source, target, "user:read:document")
		if err != nil {
			b.Fatalf("failed to evaluate: %v", err)
		}
	}
}
