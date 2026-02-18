/*
Пакет tests содержит тесты для публичного API пакета base.

Тесты проверяют валидацию значений Entity и корректность работы
с сущностями source и target через публичный метод IsValid.
*/
package tests

import (
	"testing"

	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
TestEntityIsValid тестирует публичный метод IsValid для проверки валидности значения Entity.

Тест проверяет:
  - Валидность Entity_SOURCE ("source")
  - Валидность Entity_TARGET ("target")
  - Невалидность произвольных строк
  - Невалидность пустой строки
*/
func TestEntityIsValid(t *testing.T) {
	tests := []struct {
		name    string
		entity  base.Entity
		want    bool
	}{
		{"valid source", base.Entity_SOURCE, true},        // Валидная сущность source
		{"valid target", base.Entity_TARGET, true},        // Валидная сущность target
		{"invalid entity", base.Entity("invalid"), false}, // Невалидная сущность
		{"empty entity", base.Entity(""), false},          // Пустая строка
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.entity.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

