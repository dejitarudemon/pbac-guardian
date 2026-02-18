package internal

import (
	"fmt"
	"reflect"
	"strings"
)

// испольуемый ключ тегизации
const TAG = "pce"

const (
	// разделитель пути до экспортируемого поля
	FIELD_SEP = ":"

	// минимальный размер пути
	MIN_PATH_LEN = 2

	// минимальное кол-во элементов в действии (entity и action)
	MIN_ACTION_PARTS = 2
)

var (
	// карта функций для условий
	CONDITION_TO_FUNC = map[string]ConditionFunc{
		"Contains": ContainsConditionFunc,
		"Eq":       EqConditionFunc,
		"Neq":      NeqConditionFunc,
		"Lt":       LtConditionFunc,
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
Функция isStructTag проверяет, что field не является указателем на поле структруы (экспортируемым значением)
*/
func (p *Policy) isStructTag(field string) bool {
	return strings.Contains(field, FIELD_SEP)
}

/*
Функция parseField пытается распарсить ключ из Conditions в указатель на целевую сущность и путь до поля
*/
func (p *Policy) parseField(raw string) (*Entity, []string, error) {
	if !p.isStructTag(raw) {
		return nil, nil, fmt.Errorf("expected struct tag, got literal value: %v", raw)
	}

	separeted := strings.Split(raw, FIELD_SEP)
	if len(separeted) < MIN_ACTION_PARTS {
		return nil, nil, fmt.Errorf("expected at least %v parts, but got %v", MIN_ACTION_PARTS, len(separeted))
	}

	entity := Entity(separeted[0])
	fields := separeted[1:]

	if !entity.IsValid() {
		return nil, nil, fmt.Errorf("got unexpected entity %v", entity)
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
			return nil, fmt.Errorf("nil pointer")
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, but got %v", v.Type().Name())
	}

	t := v.Type()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get(TAG)

		if tag != "" {
			flags := strings.Split(tag, ",")

			if len(flags) < 1 {
				continue
			}

			if flags[0] != path[0] {
				continue
			}

			if !field.CanInterface() {
				return nil, fmt.Errorf("unexportable field")
			}

			if len(path) > 1 {
				return p.getValue(field.Interface(), path[1:])
			}

			return field.Interface(), nil
		}
	}

	return nil, fmt.Errorf("not found")
}

/*
Функция get проверяет, является ли key - указателем на поле структуры.
Возвращает значение поле структуры или значение key (если это не указатель)
*/
func (p *Policy) get(source, target any, key string) (any, error) {
	if !p.isStructTag(key) {
		return key, nil
	}

	entity, path, err := p.parseField(key)
	if err != nil {
		return false, err
	}

	var value any

	switch *entity {
	case Entity_SOURCE:
		value, err = p.getValue(source, path)
	case Entity_TARGET:
		value, err = p.getValue(target, path)
	default:
		err = fmt.Errorf("unexpected entity")
	}

	return value, err
}

func (p *Policy) Evaluate(source, target any, action string) (bool, error) {
	if p.Action != action {
		return false, nil
	}

	match := true

	for field, condition := range p.Conditions {
		left, err := p.get(source, target, field)
		if err != nil {
			return false, err
		}

		c := reflect.ValueOf(condition)
		t := reflect.TypeOf(condition)

		for i := range c.NumField() {
			if !c.Field(i).IsZero() {
				if f, ok := CONDITION_TO_FUNC[t.Field(i).Name]; ok {

					right := c.Field(i).Interface()

					if r, ok := right.(string); ok {
						right, err = p.get(source, target, r)
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
					return false, fmt.Errorf("failed to find func")
				}
			}
		}
	}

	return match, nil
}
