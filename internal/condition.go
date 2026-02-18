package internal

import (
	"fmt"
	"reflect"
)

/*
В этом файле представлены условия для политик
*/

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
Тип ConditionFunc предназначен для обозначения функций, которые могут использоваться
для проверки условия. Порядок имеет значение только в операциях типо Contains и Lt

Входные параметры:
  - left - левое сравниваемое значение
  - right - правое сравниваемое значение

Выходные параметры:
  - bool - результат сравнения
  - err - ошибка выполнения сравнения
*/
type ConditionFunc func(left, right any) (bool, error)

/*
Фукниця проверки условия для Contains операции

Входные парамметры:
  - left - левое значение (что ищем)type Comparable
  - right - правое значение (где ищем)

Выходные параметры: аналогичны ConditionFunc
*/
func ContainsConditionFunc(left, right any) (bool, error) {
	slice := reflect.ValueOf(right)

	fmt.Printf("right is %v %v\n", right, slice.Kind())

	if slice.Kind() != reflect.Slice {
		return false, fmt.Errorf("invalid slice")
	}

	for i := range slice.Len() {
		if reflect.DeepEqual(left, slice.Index(i).Interface()) {
			return true, nil
		}
	}

	return false, nil
}

/*
Фукниця проверки условия для Eq операции

Входные парамметры: аналогичны ConditionFunc

Выходные параметры: аналогичны ConditionFunc
*/
func EqConditionFunc(left, right any) (bool, error) {
	return reflect.DeepEqual(left, right), nil
}

/*
Фукниця проверки условия для Neq операции

Входные парамметры: аналогичны ConditionFunc

Выходные параметры: аналогичны ConditionFunc
*/
func NeqConditionFunc(left, right any) (bool, error) {
	ok, err := EqConditionFunc(left, right)
	return !ok, err
}

func LtConditionFunc(left, right any) (bool, error) {
	l, ok := left.(Comparable)
	if !ok {
		return false, fmt.Errorf("cant cast left to Comparable")
	}

	result, ok := l.Compare(right)

	if !ok {
		return false, fmt.Errorf("uncomparable values")
	}

	switch result {
	case 1, 0:
		return false, nil
	case -1:
		return true, nil
	}

	return false, fmt.Errorf("unexpected behavior")
}
