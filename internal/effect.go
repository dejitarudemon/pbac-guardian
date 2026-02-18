package internal

import "slices"

var (
	// Список доступных эффектов от политик
	AVALIABLE_EFFECTS = []Effect{
		Effect_ALLOW, // разрешить
		Effect_DENY,  // запретить
	}

	// Карта приориетов эффектов от политик
	EFFECT_PRIORITIES = map[Effect]int8{
		Effect_ALLOW: 0,
		Effect_DENY:  1,
	}
)

// Приоритет невалидных эффектов
const EFFECT_UNKNOWN_PRIORITY int8 = -127

/*
Структура Effect предназначена для типизации
эффекта от выполнения политики.

Валидные значения:

	allow - Effect_ALLOW
	deny - Effect_DENY
*/
type Effect string

const (
	Effect_ALLOW Effect = "allow"
	Effect_DENY  Effect = "deny"
)

/*
Функция toPriority для Effect преобразует эффект в числовой приоритет

Выходные параметры:
  - int8 - выходной приоритет
*/
func (e Effect) toPriority() int8 {
	if priority, ok := EFFECT_PRIORITIES[e]; ok {
		return priority
	}
	return EFFECT_UNKNOWN_PRIORITY
}

// Проверяет, что вызывающий эффект менее приоритетен, чем other
func (e Effect) Lt(other Effect) bool {
	return e.toPriority() < other.toPriority()
}

// Проверяет, что вызывающий эффект равен other по приоритету
func (e Effect) Eq(other Effect) bool {
	return e.toPriority() == other.toPriority()
}

// Проверяет, что вызывающий эффект менее приоритетен, чем other, или равен ему
func (e Effect) Le(other Effect) bool {
	return e.toPriority() <= other.toPriority()
}

// Проверяет валидность значения эффекта
func (e Effect) IsValid() bool {
	return slices.Contains(AVALIABLE_EFFECTS, e)
}
