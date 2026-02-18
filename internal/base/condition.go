/*
Package base предоставляет базовые типы и функции для работы с политиками,
условиями, эффектами и сущностями в системе проверки доступа.

Пакет содержит определения политик, условий сравнения, эффектов (allow/deny),
сущностей (source/target) и интерфейсов для кастомного сравнения.
*/
package base

import (
	"context"
	"fmt"
	"reflect"
)

/*
В этом файле представлены различные условия для политик
*/

const SUPPORTED_LT_PRIMITIVES = "int|uint|float|string"

/*
Структура Condition представляет набор правил для сравнения значений в политиках.

Условие может содержать одно или несколько полей. Если указано несколько полей,
они проверяются независимо друг от друга. Значения могут быть как литералами,
так и путями до полей структур (например, "source:name" или "target:role").

Поля:
  - Contains - проверяет, находится ли левое значение в правом списке (right должен быть slice)
  - Eq - проверяет равенство значений (left == right)
  - Neq - проверяет неравенство значений (left != right)
  - Lt - проверяет, меньше ли левое значение правого (left < right)

Пример использования:

	type User struct {
		Name string   `noctis-guard:"name"`
		Role string   `noctis-guard:"role"`
		Age  int      `noctis-guard:"age"`
		Tags []string `noctis-guard:"tags"`
	}

	// Проверка равенства с литералом
	condition1 := Condition{
		Eq: "admin", // source:role == "admin"
	}

	// Проверка равенства полей двух структур
	condition2 := Condition{
		Eq: "target:owner", // source:name == target:owner
	}

	// Проверка неравенства
	condition3 := Condition{
		Neq: "guest", // source:role != "guest"
	}

	// Проверка, что значение меньше
	condition4 := Condition{
		Lt: 18, // source:age < 18
	}

	// Проверка, что значение находится в списке
	condition5 := Condition{
		Contains: []any{"admin", "moderator"}, // source:role в ["admin", "moderator"]
	}

	// Комбинирование условий (все должны быть выполнены)
	condition6 := Condition{
		Eq:       "user",           // source:role == "user"
		Contains: []any{"read"},    // "read" в source:tags
		Lt:       100,             // source:age < 100
	}
*/
type Condition struct {
	Contains []any `json:"contains,omitempty"`
	Eq       any   `json:"eq,omitempty"`
	Neq      any   `json:"neq,omitempty"`
	Lt       any   `json:"lt,omitempty"`
}

/*
Тип conditionFunc представляет функцию для проверки условия между двумя значениями.

Функции этого типа используются для проверки условий в политиках. Порядок аргументов
имеет значение для операций Contains и Lt (left и right не взаимозаменяемы).

Функции должны поддерживать отмену через context.Context и возвращать ErrCancelled
при отмене контекста.

Входные параметры:
  - ctx - контекст для отмены операции и контроля таймаутов
  - left - левое сравниваемое значение (то, что проверяется)
  - right - правое сравниваемое значение (то, с чем сравнивается)

Выходные параметры:
  - bool - результат сравнения (true если условие выполнено, false иначе)
  - err - ошибка выполнения сравнения (nil если сравнение успешно, ErrCancelled при отмене контекста)
*/
type conditionFunc func(ctx context.Context, left, right any) (bool, error)

/*
Функция containsConditionFunc проверяет, находится ли значение left в списке right.

Функция использует reflect.DeepEqual для сравнения элементов, что позволяет
работать с любыми типами данных. Поддерживает отмену через context.Context.

Входные параметры:
  - ctx - контекст для отмены операции и контроля таймаутов
  - left - значение, которое ищется в списке
  - right - список (slice) или указатель на список, в котором выполняется поиск

Выходные параметры:
  - bool - true, если left найден в right, false иначе
  - err - ошибка выполнения, если right не является списком или операция отменена

Возможные ошибки:
  - ErrCancelled - операция была отменена через context.Context
  - ErrInvalidType - right не является slice или указателем на slice (может быть nil)
*/
func containsConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	if right == nil {
		return false, NewErrInvalidType(reflect.Slice.String(), nil)
	}

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
		select {
		case <-ctx.Done():
			return false, ErrCancelled
		default:
			if reflect.DeepEqual(left, slice.Index(i).Interface()) {
				return true, nil
			}
		}
	}

	return false, nil
}

/*
Функция eqConditionFunc проверяет равенство двух значений.

Если одно из значений реализует интерфейс Comparable, используется метод Compare()
для сравнения. Если Compare() возвращает false (сравнение невозможно), или
ни одно значение не реализует Comparable, используется reflect.DeepEqual.

Входные параметры:
  - ctx - контекст для отмены операции и контроля таймаутов (не используется, но требуется для совместимости)
  - left - левое сравниваемое значение
  - right - правое сравниваемое значение

Выходные параметры:
  - bool - true, если значения равны, false иначе
  - err - ошибка выполнения (всегда nil, функция не возвращает ошибок)
*/
func eqConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	// nil значения обрабатываются через reflect.DeepEqual
	// reflect.DeepEqual(nil, nil) == true, reflect.DeepEqual(nil, non-nil) == false

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
Функция neqConditionFunc проверяет неравенство двух значений.

Функция является инверсией eqConditionFunc: возвращает !eqConditionFunc(ctx, left, right).
Использует ту же логику сравнения через Comparable или DeepEqual.

Входные параметры:
  - ctx - контекст для отмены операции и контроля таймаутов (передается в eqConditionFunc)
  - left - левое сравниваемое значение
  - right - правое сравниваемое значение

Выходные параметры:
  - bool - true, если значения не равны, false если равны
  - err - ошибка выполнения (всегда nil, функция не возвращает ошибок)
*/
func neqConditionFunc(ctx context.Context, left, right any) (bool, error) {
	ok, err := eqConditionFunc(ctx, left, right)
	return !ok, err
}

/*
Функция ltConditionFunc проверяет, меньше ли left значения right.

Если left является структурой, она должна реализовывать интерфейс Comparable.
Для примитивных типов (int, uint, float, string) используется сравнение через reflect.
Типы должны совпадать для корректного сравнения.

Входные параметры:
  - ctx - контекст для отмены операции и контроля таймаутов (передается в ltPrimitives)
  - left - левое сравниваемое значение
  - right - правое сравниваемое значение

Выходные параметры:
  - bool - true, если left < right, false иначе
  - err - ошибка выполнения, если сравнение невозможно

Возможные ошибки:
  - ErrNotComparableStruct - left является структурой, но не реализует интерфейс Comparable
  - ErrUncomparable - невозможно сравнить left и right (несовместимые типы или Compare вернул false)
  - ErrInvalidType - left или right не является ни структурой, ни поддерживаемым примитивом
    (int, uint, float, string и их варианты)
*/
func ltConditionFunc(ctx context.Context, left, right any) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	if left == nil {
		return false, NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, nil)
	}
	if right == nil {
		return false, NewErrInvalidType(SUPPORTED_LT_PRIMITIVES, nil)
	}

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
Функция ltPrimitives сравнивает примитивные типы left и right на основе рефлексии.

Функция поддерживает сравнение следующих типов:
  - int, int8, int16, int32, int64
  - uint, uint8, uint16, uint32, uint64
  - float32, float64
  - string

Типы left и right должны совпадать. Если передан указатель, он разыменовывается.

Входные параметры:
  - left - левое сравниваемое значение (примитивный тип)
  - right - правое сравниваемое значение (примитивный тип)

Выходные параметры:
  - bool - true, если left < right, false иначе
  - err - ошибка выполнения, если сравнение невозможно

Возможные ошибки:
  - ErrUncomparable - типы left и right не совпадают или один из них nil
  - ErrInvalidType - left или right не является поддерживаемым примитивом
    (int, uint, float, string и их варианты)
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
