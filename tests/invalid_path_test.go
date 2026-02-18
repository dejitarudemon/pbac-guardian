/*
Пакет tests содержит тесты для обработки невалидных путей к полям.

Тесты проверяют, что библиотека корректно обрабатывает ошибки при невалидных путях
к полям в условиях политик.
*/
package tests

import (
	"context"
	"testing"

	noctisguard "github.com/dejitarudemon/noctis-guard"
	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
TestEvaluateInvalidFieldPath тестирует обработку ошибки при невалидном пути к полю.

Тест создает политику с несуществующим полем в условии и проверяет,
что Evaluate возвращает ошибку при попытке получить значение несуществующего поля.
*/
func TestEvaluateInvalidFieldPath(t *testing.T) {
	// Создаем политику с невалидным путем к полю
	// Поле "nonexistent" не существует в структуре User
	policies := []base.Policy{
		{
			Name:   "invalid-path-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:nonexistent": {Eq: "value"}, // Несуществующее поле
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user"}
	target := Document{Owner: "user", Type: "public"}

	// Ожидаем ошибку, так как поле "nonexistent" не существует в структуре User
	_, err = engine.Evaluate(ctx, source, target, "user:read")
	if err == nil {
		t.Errorf("expected error for invalid field path, got nil")
	}
}

/*
TestEvaluateInvalidNestedPath тестирует обработку ошибки при невалидном вложенном пути.

Тест создает политику с невалидным вложенным путем и проверяет,
что Evaluate возвращает ошибку при попытке получить значение несуществующего вложенного поля.
*/
func TestEvaluateInvalidNestedPath(t *testing.T) {
	// Создаем политику с невалидным вложенным путем
	// Путь "source:user:nonexistent" невалиден, так как поле "nonexistent" не существует
	policies := []base.Policy{
		{
			Name:   "invalid-nested-path-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:user:nonexistent": {Eq: "value"}, // Несуществующее вложенное поле
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user"}
	target := Document{Owner: "user", Type: "public"}

	// Ожидаем ошибку, так как вложенное поле "nonexistent" не существует
	_, err = engine.Evaluate(ctx, source, target, "user:read")
	if err == nil {
		t.Errorf("expected error for invalid nested field path, got nil")
	}
}

/*
TestEvaluateInvalidTargetPath тестирует обработку ошибки при невалидном пути к полю в target.

Тест создает политику с невалидным путем к полю в структуре target и проверяет,
что Evaluate возвращает ошибку при попытке получить значение несуществующего поля.
*/
func TestEvaluateInvalidTargetPath(t *testing.T) {
	// Создаем политику с невалидным путем к полю в target
	// Поле "nonexistent" не существует в структуре Document
	policies := []base.Policy{
		{
			Name:   "invalid-target-path-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"target:nonexistent": {Eq: "value"}, // Несуществующее поле в target
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	source := User{Name: "user", Role: "user"}
	target := Document{Owner: "user", Type: "public"}

	// Ожидаем ошибку, так как поле "nonexistent" не существует в структуре Document
	_, err = engine.Evaluate(ctx, source, target, "user:read")
	if err == nil {
		t.Errorf("expected error for invalid target field path, got nil")
	}
}

