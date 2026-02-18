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
Структура Policy обозначает конкретную политику

Параметры:
  - Name - название политики
  - Action - действие в формате action:entity:extra...
  - Effect - эффект от политики
  - Conditions - словарь условий. Каждый ключ - экспортруемое значение, к которому примменяются Condition-операции, на которое оно указывает
*/
type Policy struct {
	Name       string               `json:"name"`
	Action     string               `json:"action"`
	Effect     Effect               `json:"effect"`
	Conditions map[string]Condition `json:"conditions"`
}

/*
Функция isPath проверяет, что path является указателем на поле структруы (экспортируемым значением)
*/
func (p *Policy) isPath(path string) bool {
	return strings.Contains(path, PATH_SEP)
}

/*
Функция parsePath пытается распарсить ключ из Conditions в указатель на целевую сущность и путь до поля
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
Функция getValue пытается найти поле в сущности по пути и получить значение
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
Функция get проверяет, является ли path - указателем на поле структуры.
Возвращает значение поле структуры или значение path (если это не указатель)
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
Функция Evaluate для Policy предназначена для проверки политики по отношению
к target и source, которые должны являться структурами

Входные параметры:
  - source - первая проверяемая структура
  - target - вторая проверяемая структура
  - action - искомая операция

Выходные параметры:
  - bool - результат применения политики
  - error - ошибка выполнения
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
