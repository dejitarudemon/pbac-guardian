/*
Пакет tests содержит тесты для публичного API пакета base.

Тесты проверяют валидацию значений Effect и корректность работы
с эффектами ALLOW и DENY через публичный метод IsValid.
*/
package tests

import (
	"testing"

	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
TestEffectIsValid тестирует публичный метод IsValid для проверки валидности значения Effect.

Тест проверяет:
  - Валидность Effect_ALLOW ("allow")
  - Валидность Effect_DENY ("deny")
  - Невалидность произвольных строк
  - Невалидность пустой строки
*/
func TestEffectIsValid(t *testing.T) {
	tests := []struct {
		name   string
		effect base.Effect
		want   bool
	}{
		{"valid allow", base.Effect_ALLOW, true},        // Валидный эффект allow
		{"valid deny", base.Effect_DENY, true},          // Валидный эффект deny
		{"invalid effect", base.Effect("invalid"), false}, // Невалидный эффект
		{"empty effect", base.Effect(""), false},        // Пустая строка
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.effect.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

