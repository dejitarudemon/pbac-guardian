/*
Package noctisguard предоставляет движок для проверки структур на соответствие политикам.

Библиотека позволяет определять политики доступа и проверять соответствие структур
этим политикам с использованием гибкой системы условий и эффектов.

Пример использования:

	package main

	import (
		"context"
		"fmt"
		"github.com/dejitarudemon/noctis-guard"
		"github.com/dejitarudemon/noctis-guard/internal/base"
	)

	type User struct {
		Name string `noctis-guard:"name"`
		Role string `noctis-guard:"role"`
	}

	type Document struct {
		Owner string `noctis-guard:"owner"`
		Type  string `noctis-guard:"type"`
	}

	func main() {
		// Создание политик
		policies := []base.Policy{
			{
				Name:   "admin-read",
				Action: "user:read:document",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:role": {
						Eq: "admin",
					},
				},
			},
			{
				Name:   "owner-read",
				Action: "user:read:document",
				Effect: base.Effect_ALLOW,
				Conditions: map[string]base.Condition{
					"source:name": {
						Eq: "target:owner",
					},
				},
			},
		}

		// Создание движка
		engine, err := noctisguard.NewNoctisFromPolices(policies)
		if err != nil {
			panic(err)
		}

		// Создание контекста
		ctx := context.Background()

		// Проверка доступа
		user := User{Name: "alice", Role: "user"}
		doc := Document{Owner: "alice", Type: "private"}

		allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
		if err != nil {
			if err == base.ErrCancelled {
				fmt.Println("Operation cancelled")
			} else {
				panic(err)
			}
		}

		fmt.Printf("Access allowed: %v\n", allowed) // true (политика owner-read прошла)
	}
*/
package noctisguard

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
Структура Noctis представляет собой основной движок библиотеки для проверки
структур на соответствие политикам доступа.

Noctis хранит политики, организованные по действиям (action), и предоставляет
метод Evaluate для проверки соответствия структур этим политикам.

Для создания экземпляра Noctis используйте:
  - NewNoctisFromPolices - создание из списка политик, переданных вручную из кода
  - NewNoctisFromFile - создание из JSON-файла с политиками
*/
type Noctis struct {
	// хранит политики, разделенные по действиям (action)
	polices map[string][]base.Policy
}

/*
Функция NewNoctisFromPolices создает новый экземпляр Noctis из списка политик,
переданных программно.

Функция выполняет валидацию политик и проверку на дубликаты имен. Политики
группируются по действиям (action) для последующей быстрой проверки.

Входные параметры:
  - polices - список политик для инициализации движка

Выходные параметры:
  - *Noctis - созданный экземпляр движка, готовый к использованию
  - error - ошибка создания, если политики содержат дубликаты имен или невалидны

Возможные ошибки:
  - ErrExport - ошибка экспорта политик. Может содержать:
  - ErrDuplicateName - если найдены политики с одинаковыми именами
  - ошибки валидации из base (ErrInvalidPath) - если политики невалидны

Пример использования:

	policies := []base.Policy{
		{
			Name:   "allow-admin",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		// обработка ошибки
	}
*/
func NewNoctisFromPolices(polices []base.Policy) (*Noctis, error) {
	mapped, err := export(polices)
	if err != nil {
		return nil, NewErrExport(err)
	}

	return &Noctis{polices: mapped}, nil
}

/*
Функция NewNoctisFromFile создает новый экземпляр Noctis из JSON-файла,
содержащего массив политик.

Функция читает файл, парсит JSON и создает движок аналогично NewNoctisFromPolices.
Файл должен содержать валидный JSON-массив объектов Policy.

Входные параметры:
  - path - путь до файла с политиками в формате JSON

Выходные параметры:
  - *Noctis - созданный экземпляр движка, готовый к использованию
  - error - ошибка создания, если файл недоступен или содержит невалидные данные

Возможные ошибки:
  - ErrExport - ошибка экспорта политик. Может содержать:
  - ошибки открытия/чтения файла (os.PathError и т.д.)
  - ошибки JSON-парсинга (json.SyntaxError и т.д.)
  - ErrDuplicateName - если найдены политики с одинаковыми именами
  - ошибки валидации из base (ErrInvalidPath) - если политики невалидны

Пример использования:

	// Файл policies.json:
	// [
	//   {
	//     "name": "allow-admin",
	//     "action": "user:read",
	//     "effect": "allow",
	//     "conditions": {
	//       "source:role": {"eq": "admin"}
	//     }
	//   }
	// ]

	engine, err := noctisguard.NewNoctisFromFile("policies.json")
	if err != nil {
		// обработка ошибки
	}
*/
func NewNoctisFromFile(path string) (*Noctis, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return nil, NewErrExport(err)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, NewErrExport(err)
	}

	var polices []base.Policy

	if err := json.Unmarshal(content, &polices); err != nil {
		return nil, NewErrExport(err)
	}

	return NewNoctisFromPolices(polices)
}

/*
Функция Evaluate проверяет соответствие структур source и target политикам
для указанного действия.

Функция ищет все политики, связанные с указанным действием, и проверяет их
условия. Результат определяется логикой эффектов: политики с эффектом DENY
имеют приоритет и запрещают действие, если их условия не выполнены. Политики
с эффектом ALLOW разрешают действие, если хотя бы одна из них проходит проверку.

Функция поддерживает отмену через context.Context, что позволяет прервать
длительные операции проверки условий.

Входные параметры:
  - ctx - контекст для отмены операции и контроля таймаутов
  - source - первая проверяемая структура (обычно источник действия)
  - target - вторая проверяемая структура (обычно цель действия)
  - action - действие в формате "entity:action:extra..." для которого проверяются политики

Выходные параметры:
  - bool - результат проверки:
  - true - действие разрешено (хотя бы одна политика ALLOW прошла проверку)
  - false - действие запрещено (нет политик для действия, или политика DENY не прошла проверку)
  - error - ошибка выполнения проверки, если возникла проблема при оценке условий

Возможные ошибки:
  - ErrEvaluate - ошибка оценки политик. Может содержать ошибки из base:
  - ErrCancelled - операция была отменена через context.Context
  - ErrInvalidPath - ошибка парсинга пути до поля или поле не найдено
  - ErrInvalidType - неверный тип структуры или поля
  - ErrUncomparable - невозможно сравнить значения в условии
  - ErrInexpectedBehavior - внутренняя ошибка (функция условия не найдена)

Логика работы:
 1. Если для указанного action нет политик, возвращается (false, nil)
 2. Для каждой политики проверяются условия:
    - Если контекст отменен, возвращается (false, ErrCancelled)
    - Если политика имеет эффект DENY и условия не выполнены, возвращается (false, nil)
    - Если политика имеет эффект ALLOW и условия выполнены, устанавливается флаг allowed = true
 3. Возвращается результат: (allowed, nil) или (false, error) при ошибке

Пример использования:

	import (
		"context"
		"time"
	)

	type User struct {
		Name string `noctis-guard:"name"`
		Role string `noctis-guard:"role"`
	}

	type Document struct {
		Owner string `noctis-guard:"owner"`
	}

	user := User{Name: "alice", Role: "admin"}
	doc := Document{Owner: "alice"}

	// Создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверка доступа на чтение документа
	allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
	if err != nil {
		if err == base.ErrCancelled {
			// операция была отменена
		} else {
			// другая ошибка
		}
	}

	if allowed {
		// доступ разрешен
	} else {
		// доступ запрещен
	}
*/
func (n *Noctis) Evaluate(ctx context.Context, source, target any, action string) (bool, error) {
	polices, ok := n.polices[action]
	if !ok {
		return false, nil
	}

	allowed := false

	for _, policy := range polices {
		ok, err := policy.Evaluate(ctx, source, target, action)
		if err != nil {
			return false, NewErrEvaluate(err)
		}

		if policy.Effect == base.Effect_DENY {
			if !ok {
				return false, err
			}
		} else {
			allowed = allowed || ok
		}
	}

	return allowed, nil
}
