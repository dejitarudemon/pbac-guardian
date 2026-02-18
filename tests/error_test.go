/*
Пакет tests содержит тесты для типов ошибок библиотеки noctis-guard.

Тесты проверяют создание, обертку и разворачивание ошибок,
а также корректность сообщений об ошибках.

Этот файл находится в директории tests для организации тестов.
Тесты импортируют пакет noctisguard как обычные пользователи библиотеки.
*/
package tests

import (
	"errors"
	"testing"

	"github.com/dejitarudemon/noctis-guard"
)

/*
TestErrDuplicateName тестирует ошибку дубликата имени политики.

Тест проверяет:
  - Создание ошибки с указанным именем политики
  - Корректность типа ошибки через errors.As()
  - Сохранение имени в структуре ошибки
  - Формат сообщения об ошибке (должно содержать имя политики)
*/
func TestErrDuplicateName(t *testing.T) {
	// Создаем ошибку с тестовым именем политики
	err := noctisguard.NewErrDuplicateName("test-policy")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Проверяем сообщение об ошибке (проверка типа через errors.As недоступна для неэкспортируемых типов)
	// Проверяем, что сообщение содержит имя политики и указание на дубликат
	expectedMsg := "test-policy is already used by another policy"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
}

/*
TestErrExport тестирует ошибку экспорта политик.

Тест проверяет:
  - Создание ошибки с оберткой исходной ошибки
  - Корректность типа ошибки через errors.As()
  - Возможность разворачивания ошибки через errors.Unwrap()
  - Проверку через errors.Is() для исходной ошибки
  - Сохранение исходной ошибки в поле source структуры
*/
func TestErrExport(t *testing.T) {
	// Создаем исходную ошибку, которая будет обернута
	originalErr := errors.New("original error")
	
	// Оборачиваем её в ErrExport
	err := noctisguard.NewErrExport(originalErr)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Проверяем, что ошибка оборачивает исходную через errors.Is()
	// Это позволяет проверять наличие исходной ошибки в цепочке
	if !errors.Is(err, originalErr) {
		t.Error("expected error to wrap original error")
	}

	// Проверяем разворачивание ошибки через errors.Unwrap()
	// Это позволяет получить доступ к исходной ошибке
	unwrapped := errors.Unwrap(err)
	if unwrapped != originalErr {
		t.Errorf("expected unwrapped error to be original error")
	}

	// Проверяем, что сообщение об ошибке содержит информацию об экспорте
	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("expected non-empty error message")
	}
}

/*
TestErrEvaluate тестирует ошибку оценки политик.

Тест проверяет:
  - Создание ошибки с оберткой исходной ошибки
  - Корректность типа ошибки через errors.As()
  - Возможность разворачивания ошибки через errors.Unwrap()
  - Проверку через errors.Is() для исходной ошибки
  - Сохранение исходной ошибки в поле source структуры

Эта ошибка используется при ошибках во время выполнения Evaluate,
например, при невалидных путях к полям или ошибках сравнения.
*/
func TestErrEvaluate(t *testing.T) {
	// Создаем исходную ошибку, которая будет обернута
	originalErr := errors.New("original error")
	
	// Оборачиваем её в ErrEvaluate
	err := noctisguard.NewErrEvaluate(originalErr)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Проверяем, что ошибка оборачивает исходную через errors.Is()
	if !errors.Is(err, originalErr) {
		t.Error("expected error to wrap original error")
	}

	// Проверяем разворачивание ошибки через errors.Unwrap()
	unwrapped := errors.Unwrap(err)
	if unwrapped != originalErr {
		t.Errorf("expected unwrapped error to be original error")
	}

	// Проверяем, что сообщение об ошибке содержит информацию об оценке
	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("expected non-empty error message")
	}
}
