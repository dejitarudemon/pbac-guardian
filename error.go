package noctisguard

import "fmt"

/*
Структура ErrDuplicateName представляет ошибку, возникающую при попытке
создать политику с именем, которое уже используется другой политикой.

Каждая политика должна иметь уникальное имя. При попытке добавить политику
с уже существующим именем возвращается эта ошибка.
*/
type ErrDuplicateName struct {
	name string
}

/*
Функция NewErrDuplicateName создает новую ошибку ErrDuplicateName.

Входные параметры:
  - name - имя политики, которое уже используется другой политикой

Выходные параметры:
  - error - созданная ошибка типа ErrDuplicateName
*/
func NewErrDuplicateName(name string) error {
	return ErrDuplicateName{name: name}
}

func (e ErrDuplicateName) Error() string {
	return fmt.Sprintf("%v is already used by another policy", e.name)
}

/*
Структура ErrExport представляет ошибку, возникающую при экспорте политик
в движок Noctis.

Ошибка оборачивает исходную ошибку (source), что позволяет использовать
errors.Unwrap() для получения деталей проблемы. Используется при создании
движка из политик или файла.
*/
type ErrExport struct {
	source error
}

/*
Функция NewErrExport создает новую ошибку ErrExport, оборачивающую исходную ошибку.

Входные параметры:
  - source - исходная ошибка, которая привела к ошибке экспорта (может быть nil)

Выходные параметры:
  - error - созданная ошибка типа ErrExport, которую можно развернуть через errors.Unwrap()
*/
func NewErrExport(source error) error {
	return ErrExport{source: source}
}

func (e ErrExport) Unwrap() error {
	return e.source
}

func (e ErrExport) Error() string {
	return fmt.Sprintf("failed to export polices: %v", e.source)
}

/*
Структура ErrEvaluate представляет ошибку, возникающую при оценке политик
во время выполнения метода Evaluate.

Ошибка оборачивает исходную ошибку (source), что позволяет использовать
errors.Unwrap() для получения деталей проблемы. Возникает при проверке
условий политик или доступа к полям структур.
*/
type ErrEvaluate struct {
	source error
}

/*
Функция NewErrEvaluate создает новую ошибку ErrEvaluate, оборачивающую исходную ошибку.

Входные параметры:
  - source - исходная ошибка, которая привела к ошибке оценки (может быть nil)

Выходные параметры:
  - error - созданная ошибка типа ErrEvaluate, которую можно развернуть через errors.Unwrap()
*/
func NewErrEvaluate(source error) error {
	return ErrEvaluate{source: source}
}

func (e ErrEvaluate) Unwrap() error {
	return e.source
}

func (e ErrEvaluate) Error() string {
	return fmt.Sprintf("failed to evaluate: %v", e.source)
}
