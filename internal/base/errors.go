/*
Package base предоставляет базовые типы и функции для работы с политиками,
условиями, эффектами и сущностями в системе проверки доступа.

Пакет содержит определения политик, условий сравнения, эффектов (allow/deny),
сущностей (source/target) и интерфейсов для кастомного сравнения.
*/
package base

import (
	"errors"
	"fmt"
)

/*
В этом файле представлены кастомные ошибки для базы движка
*/

var (
	// ErrNotComparableStruct представляет ошибку, возникающую когда левый аргумент
	// является структурой, но не реализует метод Compare()
	ErrNotComparableStruct = errors.New("left argument is a struct, but it doesn't implement Comapre() method")

	ErrCancelled = errors.New("cancelled by context")
)

/*
Структура ErrInvalidType представляет ошибку, возникающую при несоответствии
ожидаемого и полученного типов.

Ошибка используется, когда функция ожидает определенный тип (например, структуру
или slice), но получает другой тип.
*/
type ErrInvalidType struct {
	expected any
	got      any
}

/*
Функция NewErrInvalidType создает новую ошибку ErrInvalidType.

Входные параметры:
  - expected - ожидаемый тип (может быть строкой или типом)
  - got - полученный тип (может быть строкой, типом или nil)

Выходные параметры:
  - error - созданная ошибка типа ErrInvalidType
*/
func NewErrInvalidType(expected, got any) error {
	return ErrInvalidType{expected: expected, got: got}
}

func (e ErrInvalidType) Error() string {
	return fmt.Sprintf("expected %v, but got %v", e.expected, e.got)
}

/*
Структура ErrUncomparable представляет ошибку, возникающую при попытке
сравнить два значения, которые невозможно сравнить между собой.

Ошибка возникает, когда типы значений несовместимы для сравнения (например,
разные типы примитивов в операции Lt) или метод Compare() вернул false.
*/
type ErrUncomparable struct {
	left  any
	right any
}

/*
Функция NewErrUncomparable создает новую ошибку ErrUncomparable.

Входные параметры:
  - left - левое сравниваемое значение
  - right - правое сравниваемое значение

Выходные параметры:
  - error - созданная ошибка типа ErrUncomparable
*/
func NewErrUncomparable(left, right any) error {
	return ErrUncomparable{left: left, right: right}
}

func (e ErrUncomparable) Error() string {
	return fmt.Sprintf("can't compare %v with %v", e.left, e.right)
}

/*
Структура ErrInvalidPath представляет ошибку, возникающую при работе
с невалидным путем до поля структуры.

Ошибка используется при парсинге путей в формате "entity:field1:field2..."
или при поиске полей в структурах. Содержит путь и детали проблемы.
*/
type ErrInvalidPath struct {
	path    string
	details string
}

/*
Функция NewErrInvalidPath создает новую ошибку ErrInvalidPath.

Входные параметры:
  - path - невалидный путь, который вызвал ошибку
  - details - детали ошибки (описание проблемы)

Выходные параметры:
  - error - созданная ошибка типа ErrInvalidPath
*/
func NewErrInvalidPath(path, details string) error {
	return ErrInvalidPath{path: path, details: details}
}

func (e ErrInvalidPath) Error() string {
	return fmt.Sprintf("invalid path %v: %v", e.path, e.details)
}

/*
Структура ErrInexpectedBehavior представляет внутреннюю ошибку библиотеки,
возникающую при неожиданном поведении в коде.

Ошибка указывает на проблему в логике библиотеки (например, отсутствие функции
условия в CONDITION_TO_FUNC) и обычно не должна возникать при правильном использовании.
*/
type ErrInexpectedBehavior struct {
	source  string
	details string
}

/*
Функция NewErrInexpectedBehavior создает новую ошибку ErrInexpectedBehavior.

Входные параметры:
  - source - источник ошибки (название функции или места возникновения, например "Policy.Evaluate()")
  - details - детали неожиданного поведения (описание проблемы)

Выходные параметры:
  - error - созданная ошибка типа ErrInexpectedBehavior
*/
func NewErrInexpectedBehavior(source, details string) error {
	return ErrInexpectedBehavior{source: source, details: details}
}

func (e ErrInexpectedBehavior) Error() string {
	return fmt.Sprintf("unexpected behavior in %v : %v", e.source, e.details)
}
