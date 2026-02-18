/*
Package base предоставляет базовые типы и функции для работы с политиками,
условиями, эффектами и сущностями в системе проверки доступа.

Пакет содержит определения политик, условий сравнения, эффектов (allow/deny),
сущностей (source/target) и интерфейсов для кастомного сравнения.
*/
package base

import "slices"

/*
Тип Entity представляет сущность, к которой относится поле в политике.

Entity используется в путях до полей для указания, из какой структуры (source или target)
нужно извлекать значение. Например, в пути "source:name" Entity_SOURCE указывает,
что поле name нужно брать из структуры source.

Валидные значения:
  - Entity_SOURCE ("source") - поле относится к структуре source
  - Entity_TARGET ("target") - поле относится к структуре target
*/
type Entity string

const (
	Entity_TARGET Entity = "target"
	Entity_SOURCE Entity = "source"
)

var (
	// Список валидных целей
	AVALIABLE_ENTITIES = []Entity{
		Entity_SOURCE,
		Entity_TARGET,
	}
)

/*
Метод IsValid проверяет, является ли значение Entity валидным.

Валидными считаются только Entity_SOURCE и Entity_TARGET.

Выходные параметры:
  - bool - true, если значение валидно, false иначе
*/
func (e Entity) IsValid() bool {
	return slices.Contains(AVALIABLE_ENTITIES, e)
}
