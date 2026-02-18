/*
Package base предоставляет базовые типы и функции для работы с политиками,
условиями, эффектами и сущностями в системе проверки доступа.

Пакет содержит определения политик, условий сравнения, эффектов (allow/deny),
сущностей (source/target) и интерфейсов для кастомного сравнения.
*/
package base

import (
	"fmt"
	"reflect"
	"strings"
)

// испольуемый ключ тегизации
const TAG_KEY = "noctis-guard"

const (
	// разделитель пути до экспортируемого поля
	PATH_SEP = ":"

	// разделитель внутри тега
	TAG_SEP = ","

	// минимальный размер пути
	MIN_PATH_LEN = 2

	// минимальное кол-во элементов в действии (entity и action)
	MIN_ACTION_PARTS = 2
)

var (
	// карта функций для условий
	CONDITION_TO_FUNC = map[string]conditionFunc{
		"Contains": containsConditionFunc,
		"Eq":       eqConditionFunc,
		"Neq":      neqConditionFunc,
		"Lt":       ltConditionFunc,
	}
)

/*
Структура Policy представляет одну политику доступа с набором условий.

Политика определяет правила проверки структур source и target для конкретного
действия. Условия проверяются через поля структуры, помеченные тегом "noctis-guard".

Параметры:
  - Name - уникальное название политики (используется для идентификации)
  - Action - действие в формате "entity:action:extra1:extra2..." (например, "user:read:profile")
  - Effect - эффект от политики: Effect_ALLOW (разрешить) или Effect_DENY (запретить)
  - Conditions - словарь условий. Ключ - путь до поля в формате "source:field" или "target:field",
    значение - условие для проверки (Contains, Eq, Neq, Lt)

Пример использования:

	type User struct {
		Name string `noctis-guard:"name"`
		Role string `noctis-guard:"role"`
		Age  int    `noctis-guard:"age"`
	}

	type Document struct {
		Owner string   `noctis-guard:"owner"`
		Tags  []string `noctis-guard:"tags"`
	}

	// Политика: разрешить чтение документа админам
	policy1 := Policy{
		Name:   "admin-read",
		Action: "user:read:document",
		Effect: Effect_ALLOW,
		Conditions: map[string]Condition{
			"source:role": {
				Eq: "admin",
			},
		},
	}

	// Политика: разрешить чтение документа владельцу
	policy2 := Policy{
		Name:   "owner-read",
		Action: "user:read:document",
		Effect: Effect_ALLOW,
		Conditions: map[string]Condition{
			"source:name": {
				Eq: "target:owner", // сравнение полей двух структур
			},
		},
	}

	// Политика: запретить чтение документов с тегом "private" для пользователей младше 18
	policy3 := Policy{
		Name:   "age-restriction",
		Action: "user:read:document",
		Effect: Effect_DENY,
		Conditions: map[string]Condition{
			"source:age": {
				Lt: 18,
			},
			"target:tags": {
				Contains: []any{"private"},
			},
		},
	}
*/
type Policy struct {
	Name       string               `json:"name"`
	Action     string               `json:"action"`
	Effect     Effect               `json:"effect"`
	Conditions map[string]Condition `json:"conditions"`
}

/*
Функция isPath проверяет, является ли path путем до поля структуры.

Путь до поля должен содержать разделитель PATH_SEP (":"), что указывает на то,
что это не литеральное значение, а ссылка на поле в структуре source или target.

Входные параметры:
  - path - строка для проверки

Выходные параметры:
  - bool - true, если path содержит разделитель ":" и является путем, иначе false
*/
func (p *Policy) isPath(path string) bool {
	return strings.Contains(path, PATH_SEP)
}

/*
Функция parsePath парсит путь из Conditions в сущность и путь до поля.

Путь должен иметь формат "entity:field1:field2...", где entity - это "source"
или "target", а field1, field2... - иерархический путь до поля в структуре.

Примеры валидных путей:
  - "source:name" - поле name в структуре source
  - "target:user:email" - поле email во вложенной структуре user в target

Входные параметры:
  - path - путь для парсинга в формате "entity:field1:field2..."

Выходные параметры:
  - *Entity - указатель на сущность (Entity_SOURCE или Entity_TARGET)
  - []string - путь до поля в виде массива строк ["field1", "field2", ...]
  - error - ошибка парсинга пути

Возможные ошибки:
  - ErrInvalidPath - возникает если:
  - path не содержит разделитель ":" (не является путем)
  - путь содержит менее 2 частей (минимум entity и одно поле)
  - первая часть пути не является валидной сущностью (не "source" и не "target")
*/
func (p *Policy) parsePath(path string) (*Entity, []string, error) {
	if !p.isPath(path) {
		return nil, nil, NewErrInvalidPath(path, "it is not a path")
	}

	separeted := strings.Split(path, PATH_SEP)
	if len(separeted) < MIN_PATH_LEN {
		return nil, nil, NewErrInvalidPath(path, fmt.Sprintf("expected at least %v parts, but got %v", MIN_ACTION_PARTS, len(separeted)))
	}

	entity := Entity(separeted[0])
	fields := separeted[1:]

	if !entity.IsValid() {
		return nil, nil, NewErrInvalidPath(path, fmt.Sprintf("path allocates to unknown entity: %v", entity))
	}

	return &entity, fields, nil
}

/*
Функция getValue находит поле в структуре по пути и возвращает его значение.

Функция рекурсивно проходит по пути, используя теги "noctis-guard" для поиска
полей. Поле должно быть экспортируемым (с заглавной буквы) и иметь тег с
соответствующим именем.

Входные параметры:
  - entity - сущность (структура или указатель на структуру) для поиска поля
  - path - путь до поля в виде массива строк ["field1", "field2", ...]

Выходные параметры:
  - any - найденное значение поля
  - error - ошибка поиска поля

Возможные ошибки:
  - ErrInvalidType - сущность не является структурой или указателем на структуру
  - ErrInvalidPath - возникает если:
  - поле не найдено (нет тега с соответствующим именем)
  - поле является неэкспортируемым (недоступно через CanInterface())
*/
func (p *Policy) getValue(entity any, path []string) (any, error) {
	v := reflect.ValueOf(entity)

	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, NewErrInvalidType(fmt.Sprintf("%v or %v", reflect.Pointer.String(), reflect.Struct.String()), v.Kind().String())
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, NewErrInvalidType(fmt.Sprintf("%v or %v", reflect.Pointer.String(), reflect.Struct.String()), v.Kind().String())
	}

	t := v.Type()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get(TAG_KEY)

		if tag != "" {
			flags := strings.Split(tag, TAG_KEY)

			if len(flags) < 1 {
				continue
			}

			if flags[0] != path[0] {
				continue
			}

			if !field.CanInterface() {
				return nil, NewErrInvalidPath(path[0], "got unexpored field")
			}

			if len(path) > 1 {
				return p.getValue(field.Interface(), path[1:])
			}

			return field.Interface(), nil
		}
	}

	return nil, NewErrInvalidPath(path[0], "param doesn't exist")
}

/*
Функция get получает значение из пути или возвращает литеральное значение.

Если path является путем (содержит ":"), функция парсит его и извлекает значение
из соответствующей структуры (source или target). Если path не является путем,
возвращается само значение path как литеральное значение.

Входные параметры:
  - source - первая структура для поиска полей (используется для путей "source:...")
  - target - вторая структура для поиска полей (используется для путей "target:...")
  - path - путь до поля (например, "source:name") или литеральное значение (например, "admin")
  - mustBePath - флаг, указывающий, должен ли path обязательно быть путем (true) или может быть литералом (false)

Выходные параметры:
  - any - найденное значение поля из структуры или литеральное значение path
  - error - ошибка получения значения

Возможные ошибки:
  - ErrInvalidPath - возникает если:
  - mustBePath=true, но path не содержит ":" (не является путем)
  - ошибка парсинга пути (см. parsePath)
  - ошибка поиска поля (см. getValue)
  - ErrInvalidType - сущность не является структурой или указателем на структуру (см. getValue)
*/
func (p *Policy) get(source, target any, path string, mustBePath bool) (any, error) {
	if !p.isPath(path) {
		if mustBePath {
			return nil, NewErrInvalidPath(path, "must be a path, but it's literal value")
		}
		return path, nil
	}

	entity, parsedPath, err := p.parsePath(path)
	if err != nil {
		return false, err
	}

	var value any

	switch *entity {
	case Entity_SOURCE:
		value, err = p.getValue(source, parsedPath)
	case Entity_TARGET:
		value, err = p.getValue(target, parsedPath)
	default:
		err = NewErrInvalidPath(path, fmt.Sprintf("unxpected entity: %v", entity))
	}

	return value, err
}

/*
Функция Evaluate применяет политику к структурам source и target для указанного действия.

Функция проверяет, соответствует ли действие политики переданному action. Если да,
проверяются все условия политики. Политика считается выполненной, если все условия
выполнены (логическое И).

Входные параметры:
  - source - первая проверяемая структура (обычно источник действия)
  - target - вторая проверяемая структура (обычно цель действия)
  - action - действие в формате "entity:action:extra..." для проверки

Выходные параметры:
  - bool - результат применения политики:
  - true - политика соответствует action и все условия выполнены
  - false - политика не соответствует action или хотя бы одно условие не выполнено
  - error - ошибка выполнения, если возникла проблема при проверке условий

Возможные ошибки:
  - ErrInvalidPath - ошибка парсинга пути или поиска поля в структуре
  - ErrInvalidType - ошибка типа при получении значения поля (структура не того типа)
  - ErrUncomparable - невозможно сравнить значения в условии (несовместимые типы)
  - ErrInexpectedBehavior - внутренняя ошибка: функция условия не найдена в CONDITION_TO_FUNC
*/
func (p *Policy) Evaluate(source, target any, action string) (bool, error) {
	if p.Action != action {
		return false, nil
	}

	match := true
	t := reflect.TypeFor[Condition]()

	for field, condition := range p.Conditions {
		left, err := p.get(source, target, field, true)
		if err != nil {
			return false, err
		}

		c := reflect.ValueOf(condition)

		for i := range c.NumField() {
			if !c.Field(i).IsZero() {
				if f, ok := CONDITION_TO_FUNC[t.Field(i).Name]; ok {

					right := c.Field(i).Interface()

					if r, ok := right.(string); ok {
						right, err = p.get(source, target, r, false)
						if err != nil {
							return false, err
						}
					}

					m, err := f(left, right)
					if err != nil {
						return false, err
					}

					match = match && m
				} else {
					return false, NewErrInexpectedBehavior("Policy.Evaluate()", fmt.Sprintf("condition func for %v doesn't exist", t.Field(i).Name))
				}
			}
		}
	}

	return match, nil
}

/*
Функция IsValid проверяет валидность политики.

Функция проверяет:
  1. Формат действия (action) - должно быть минимум 2 части, разделенные ":"
  2. Отсутствие пустых частей в действии
  3. Валидность всех путей в условиях (через parsePath)

Выходные параметры:
  - error - ошибка валидности, если политика невалидна, nil если политика валидна

Возможные ошибки:
  - ErrInvalidPath - возникает если:
    * действие содержит менее 2 частей (минимум entity и action)
    * действие содержит пустые части
    * пути в условиях невалидны (см. parsePath)
*/

func (p *Policy) IsValid() error {
	actions := strings.Split(p.Action, PATH_SEP)
	if len(actions) < MIN_ACTION_PARTS {
		return NewErrInvalidPath(p.Action, "not enough parts of action. use: entity:action:extra1:extra2 etc")
	}

	for i, action := range actions {
		if action == "" {
			return NewErrInvalidPath(p.Action, fmt.Sprintf("empty part: %v", i))
		}
	}

	for field := range p.Conditions {
		if _, _, err := p.parsePath(field); err != nil {
			return err
		}
	}

	return nil
}
