package noctisguard

import (
	"slices"

	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
Функция export преобразует список политик в карту, организованную по действиям (action).

Функция выполняет следующие проверки:
 1. Проверка на дубликаты имен политик
 2. Валидация каждой политики через Policy.IsValid()
 3. Группировка политик по действиям для быстрого доступа

Входные параметры:
  - polices - список политик для экспорта и группировки

Выходные параметры:
  - map[string][]base.Policy - карта политик, где:
  - ключ - действие (action) в формате "entity:action:extra..."
  - значение - список всех политик для этого действия
  - error - ошибка экспорта, если найдены дубликаты имен или невалидные политики

Возможные ошибки:
  - ErrDuplicateName - имя политики уже используется другой политикой в списке
  - ошибки из base.Policy.IsValid() - ErrInvalidPath при невалидном пути в условиях политики
*/
func export(polices []base.Policy) (map[string][]base.Policy, error) {
	mappedPolices := make(map[string][]base.Policy)
	usedNames := make([]string, 0, len(polices))

	for _, policy := range polices {
		if slices.Contains(usedNames, policy.Name) {
			return nil, NewErrDuplicateName(policy.Name)
		}

		usedNames = append(usedNames, policy.Name)

		if err := policy.IsValid(); err != nil {
			return nil, err
		}

		if _, ok := mappedPolices[policy.Action]; !ok {
			mappedPolices[policy.Action] = make([]base.Policy, 0, 1)
		}

		mappedPolices[policy.Action] = append(mappedPolices[policy.Action], policy)
	}

	return mappedPolices, nil
}
