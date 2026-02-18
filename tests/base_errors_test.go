/*
Пакет tests содержит тесты для публичных конструкторов ошибок пакета base.

Тесты проверяют создание ошибок, их сообщения
и корректность работы с errors.As() и errors.Is().
*/
package tests

import (
	"errors"
	"testing"

	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
TestErrInvalidType тестирует публичный конструктор NewErrInvalidType.

Тест проверяет:
  - Создание ошибки с ожидаемым и полученным типами
  - Корректность типа ошибки через errors.As()
  - Наличие сообщения об ошибке
*/
func TestErrInvalidType(t *testing.T) {
	// Создаем ошибку с ожидаемым типом "struct" и полученным "string"
	err := base.NewErrInvalidType("struct", "string")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Проверяем тип ошибки
	var invalidTypeErr base.ErrInvalidType
	if !errors.As(err, &invalidTypeErr) {
		t.Errorf("expected ErrInvalidType, got %T", err)
	}

	// Проверяем, что сообщение об ошибке не пустое
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

/*
TestErrUncomparable тестирует публичный конструктор NewErrUncomparable.

Тест проверяет:
  - Создание ошибки с левым и правым значениями
  - Корректность типа ошибки через errors.As()
  - Наличие сообщения об ошибке
*/
func TestErrUncomparable(t *testing.T) {
	// Создаем ошибку для несовместимых типов (int и string)
	err := base.NewErrUncomparable(5, "string")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Проверяем тип ошибки
	var uncomparableErr base.ErrUncomparable
	if !errors.As(err, &uncomparableErr) {
		t.Errorf("expected ErrUncomparable, got %T", err)
	}

	// Проверяем, что сообщение об ошибке не пустое
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

/*
TestErrInvalidPath тестирует публичный конструктор NewErrInvalidPath.

Тест проверяет:
  - Создание ошибки с путем и деталями
  - Корректность типа ошибки через errors.As()
  - Наличие сообщения об ошибке
*/
func TestErrInvalidPath(t *testing.T) {
	// Создаем ошибку с невалидным путем и описанием проблемы
	err := base.NewErrInvalidPath("source:invalid", "field not found")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Проверяем тип ошибки
	var invalidPathErr base.ErrInvalidPath
	if !errors.As(err, &invalidPathErr) {
		t.Errorf("expected ErrInvalidPath, got %T", err)
	}

	// Проверяем, что сообщение об ошибке не пустое
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

/*
TestErrInexpectedBehavior тестирует публичный конструктор NewErrInexpectedBehavior.

Тест проверяет:
  - Создание ошибки с источником и деталями
  - Корректность типа ошибки через errors.As()
  - Наличие сообщения об ошибке
*/
func TestErrInexpectedBehavior(t *testing.T) {
	// Создаем ошибку с указанием источника и деталей проблемы
	err := base.NewErrInexpectedBehavior("TestFunction", "unexpected condition")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Проверяем тип ошибки
	var unexpectedErr base.ErrInexpectedBehavior
	if !errors.As(err, &unexpectedErr) {
		t.Errorf("expected ErrInexpectedBehavior, got %T", err)
	}

	// Проверяем, что сообщение об ошибке не пустое
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

/*
TestErrCancelled тестирует глобальную публичную переменную ErrCancelled.

Тест проверяет:
  - Что ErrCancelled не является nil
  - Корректность сообщения об ошибке
*/
func TestErrCancelled(t *testing.T) {
	// Проверяем, что ошибка определена
	if base.ErrCancelled == nil {
		t.Fatal("ErrCancelled should not be nil")
	}

	// Проверяем сообщение об ошибке
	if base.ErrCancelled.Error() != "cancelled by context" {
		t.Errorf("expected error message 'cancelled by context', got %q", base.ErrCancelled.Error())
	}
}

/*
TestErrNotComparableStruct тестирует глобальную публичную переменную ErrNotComparableStruct.

Тест проверяет:
  - Что ErrNotComparableStruct не является nil
  - Корректность сообщения об ошибке
*/
func TestErrNotComparableStruct(t *testing.T) {
	// Проверяем, что ошибка определена
	if base.ErrNotComparableStruct == nil {
		t.Fatal("ErrNotComparableStruct should not be nil")
	}

	// Проверяем сообщение об ошибке
	expectedMsg := "left argument is a struct, but it doesn't implement Comapre() method"
	if base.ErrNotComparableStruct.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, base.ErrNotComparableStruct.Error())
	}
}

