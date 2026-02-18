package internal

import "slices"

/*
Структура Effect предназначена для типизации
цели экспортируемых полей.

Валидные значения:

	target - Entity_TARGET
	source - Entity_SOURCE
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

// Проверяет валидность значения цели
func (e Entity) IsValid() bool {
	return slices.Contains(AVALIABLE_ENTITIES, e)
}
