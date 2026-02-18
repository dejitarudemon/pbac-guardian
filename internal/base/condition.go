package base

import (
	"reflect"
)

/*
В этом файле представлены различные условия для политик
*/

const SUPPORTED_LT_PRIMITIVES = "int|uint|float|string"

/*
Структура Condition представляет собой набор правил, по которым можно сравнивать некоторые значения

- Contains - находится ли значение в списке
- Eq - равны ли значения
- Neq - равны ли значения с реверсом
- Lt - меньше ли значение
*/
type Condition struct {
	Contains []any `json:"contains,omitempty"`
	Eq       any   `json:"eq,omitempty"`
	Neq      any   `json:"neq,omitempty"`
	Lt       any   `json:"lt,omitempty"`
}

/*
Тип conditionFunc предназначен для обозначения функций, которые могут использоваться
для проверки условия. Порядок имеет значение только в операциях типо Contains и Lt

Входные параметры:
  - left - левое сравниваемое значение
  - right - правое сравниваемое значение

Выходные параметры:
  - bool - результат сравнения
  - err - ошибка выполнения сравнения
*/
type conditionFunc func(left, right any) (bool, error)

/*
Фукниця проверки условия для Contains операции

Входные парамметры:
  - left - левое значение (что ищем)
  - right - правое значение (где ищем)

Выходные параметры:
  - bool - находится ли left среди right
  - err - ошибка выполнения

Возможные ошибки:
  - ErrInvalidType - right значения не является slice или не nil-указателем на slice
*/
func containsConditionFunc(left, right any) (bool, error) {
	slice := reflect.ValueOf(right)

	if slice.Kind() == reflect.Pointer {
		if slice.IsNil() {
			return false, NewErrInvalidType(reflect.Slice.String(), nil)
		}

		slice = slice.Elem()
	}

	if slice.Kind() != reflect.Slice {
		return false, NewErrInvalidType(reflect.Slice.String(), slice.Kind().String())
	}

	for i := range slice.Len() {
		if reflect.DeepEqual(left, slice.Index(i).Interface()) {
			return true, nil
		}
	}

	return false, nil
}

/*
Фукниця проверки условия для Eq операции. Если left или right
являются Comparable, сравнение сначала будет проводится через
функцию Compare. В случае неудачи, будет произведено DeepEqual-сравнение

Входные парамметры:
  - left - левый аргумент
  - right - правый аргумент

Выходные параметры:
  - bool - являются ли left и right одинаковыми
  - err - ошибка выполнения. Всегда nil
*/
func eqConditionFunc(left, right any) (bool, error) {
	if l, ok := left.(Comparable); ok {
		if result, acceptable := l.Compare(right); acceptable {
			return result == 0, nil
		}
	}
	if r, ok := right.(Comparable); ok {
		if result, acceptable := r.Compare(left); acceptable {
			return result == 0, nil
		}
	}
	return reflect.DeepEqual(left, right), nil
}

/*
Фукниця проверки условия для Neq операции. Если left или right
являются Comparable, сравнение сначала будет проводится через
функцию Compare. В случае неудачи, будет произведено DeepEqual-сравнение

Входные парамметры:
  - left - левый аргумент
  - right - правый аргумент

Выходные параметры:
  - bool - являются ли left и right не одинаковыми
  - err - ошибка выполнения. Всегда nil
*/
func neqConditionFunc(left, right any) (bool, error) {
	ok, err := eqConditionFunc(left, right)
	return !ok, err
}

/*
Фукниця проверки условия для Lt операции.
Если left и right являтся структурами, то они должны реализовывать
интерфейс Comparable. Примитивы (int, float, string и т.д.) сравниваются
через reflect.

Входные парамметры:
  - left - левый аргумент
  - right - правый аргумент

Выходные параметры:
  - bool - является ли left меньше right
  - err - ошибка выполнения.

Возможные ошибки:
  - ErrNotComparableStruct - левый аргумент является структурой, но метод Compare() не реализован
  - ErrUncomparable - левый аргумент невозмоно сравнить с правым
  - ErrInvalidType - если левый или правый не является ни структурой, ни примитовом, обозначенным в SUPPORTED_LT_PRIMITIVES
*/
func ltConditionFunc(left, right any) (bool, error) {
	if reflect.TypeOf(left).Kind() != reflect.Struct {
		return ltPrimitives(left, right)
	}

	l, ok := left.(Comparable)
	if !ok {
		return false, ErrNotComparableStruct
	}

	result, ok := l.Compare(right)

	if !ok {
		return false, NewErrUncomparable(left, right)
	}

	if result < 0 {
		return true, nil
	}

	return false, nil
}

/*
Функция ltPrimitives предназначена для сравнения left и right аргументов
на основе рефлексии.

Входные парамметры:
  - left - левый аргумент
  - right - правый аргумент

Выходные параметры:
  - bool - является ли left меньше right
  - err - ошибка выполнения.

Возможные ошибки:
  - ErrUncomparable - левый аргумент невозмоно сравнить с правым
  - ErrInvalidType - если левый или правый не является примитовом, обозначенным в SUPPORTED_LT_PRIMITIVES
*/
func ltPrimitives(left, right any) (bool, error) {
	v1 := reflect.ValueOf(left)
	v2 := reflect.ValueOf(right)

	if v1.Kind() == reflect.Pointer {
		if v1.IsNil() {
			return false, NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, nil)
		}

		v1 = v1.Elem()
	}

	if v2.Kind() == reflect.Pointer {
		if v2.IsNil() {
			return false, NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, nil)
		}

		v2 = v2.Elem()
	}

	if v1.Type() != v2.Type() {
		return false, NewErrUncomparable(left, right)
	}

	switch v1.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int8, reflect.Int64:
		return v1.Int() < v2.Int(), nil
	case reflect.Float32, reflect.Float64:
		return v1.Float() < v2.Float(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v1.Uint() < v2.Uint(), nil
	case reflect.String:
		return v1.String() < v2.String(), nil
	}

	return false, NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, v1.Kind().String())
}
