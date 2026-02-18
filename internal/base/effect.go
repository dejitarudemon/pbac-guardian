package base

import "slices"

var (
	// Список доступных эффектов от политик
	AVALIABLE_EFFECTS = []Effect{
		Effect_ALLOW, // разрешить
		Effect_DENY,  // запретить
	}
)

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

// Проверяет валидность значения эффекта
func (e Effect) IsValid() bool {
	return slices.Contains(AVALIABLE_EFFECTS, e)
}
