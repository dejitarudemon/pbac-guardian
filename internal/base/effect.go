/*
Package base предоставляет базовые типы и функции для работы с политиками,
условиями, эффектами и сущностями в системе проверки доступа.

Пакет содержит определения политик, условий сравнения, эффектов (allow/deny),
сущностей (source/target) и интерфейсов для кастомного сравнения.
*/
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
Тип Effect представляет эффект от применения политики.

Effect определяет, что происходит, если все условия политики выполнены:
  - Effect_ALLOW разрешает действие
  - Effect_DENY запрещает действие

Политики с эффектом DENY имеют приоритет: если хотя бы одна политика DENY
не прошла проверку, действие запрещается, даже если есть политики ALLOW.

Валидные значения:
  - Effect_ALLOW ("allow") - разрешить действие при выполнении условий
  - Effect_DENY ("deny") - запретить действие при невыполнении условий
*/
type Effect string

const (
	Effect_ALLOW Effect = "allow"
	Effect_DENY  Effect = "deny"
)

/*
Метод IsValid проверяет, является ли значение Effect валидным.

Валидными считаются только Effect_ALLOW и Effect_DENY.

Выходные параметры:
  - bool - true, если значение валидно, false иначе
*/
func (e Effect) IsValid() bool {
	return slices.Contains(AVALIABLE_EFFECTS, e)
}
