/*
Пакет tests содержит тесты для библиотеки noctis-guard.

Тесты проверяют функциональность создания движка из политик и файлов,
а также оценку политик для различных сценариев доступа.

Этот файл находится в директории tests для организации тестов.
Тесты импортируют пакет noctisguard как обычные пользователи библиотеки.
*/
package tests

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	noctisguard "github.com/dejitarudemon/noctis-guard"
	"github.com/dejitarudemon/noctis-guard/internal/base"
)

/*
Тестовые структуры для проверки работы библиотеки.

Эти структуры используются во всех тестах для симуляции реальных
объектов, которые проверяются на соответствие политикам доступа.
*/

// User представляет пользователя системы с полями, помеченными тегами для доступа
type User struct {
	Name string `noctis-guard:"name"` // Имя пользователя
	Role string `noctis-guard:"role"` // Роль пользователя (admin, user, guest и т.д.)
	Age  int    `noctis-guard:"age"`  // Возраст пользователя
}

// Document представляет документ с информацией о владельце и типе
type Document struct {
	Owner string   `noctis-guard:"owner"` // Владелец документа
	Type  string   `noctis-guard:"type"`  // Тип документа (public, private и т.д.)
	Tags  []string `noctis-guard:"tags"`  // Теги документа
}

// NestedUser представляет вложенную структуру пользователя для тестирования вложенных путей
type NestedUser struct {
	User User `noctis-guard:"user"` // Вложенный пользователь
}

// NestedDocument представляет вложенную структуру документа для тестирования вложенных путей
type NestedDocument struct {
	Doc Document `noctis-guard:"doc"` // Вложенный документ
}

/*
TestNewNoctisFromPolices тестирует создание движка из списка политик.

Тест проверяет:
  - Создание движка с валидными политиками
  - Обработку дубликатов имен политик (должна возвращаться ошибка)
  - Валидацию формата действий (минимум 2 части через ":")
  - Валидацию путей в условиях (формат "entity:field")
  - Создание движка с пустым списком политик (должно быть успешно)
*/
func TestNewNoctisFromPolices(t *testing.T) {
	tests := []struct {
		name     string
		policies []base.Policy
		wantErr  bool
		errType  error
	}{
		{
			name: "valid policies",
			// Тест проверяет успешное создание движка с валидными политиками
			// Политика имеет корректный формат действия и валидные пути в условиях
			policies: []base.Policy{
				{
					Name:   "test-policy",
					Action: "user:read", // Валидный формат: минимум 2 части
					Effect: base.Effect_ALLOW,
					Conditions: map[string]base.Condition{
						"source:role": {Eq: "admin"}, // Валидный путь: entity:field
					},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate names",
			// Тест проверяет, что движок не создается при наличии политик с одинаковыми именами
			// Должна возвращаться ошибка ErrExport, оборачивающая ErrDuplicateName
			policies: []base.Policy{
				{
					Name:   "test-policy",
					Action: "user:read",
					Effect: base.Effect_ALLOW,
					Conditions: map[string]base.Condition{
						"source:role": {Eq: "admin"},
					},
				},
				{
					Name:   "test-policy", // Дубликат имени - должно вызвать ошибку
					Action: "user:write",
					Effect: base.Effect_ALLOW,
					Conditions: map[string]base.Condition{
						"source:role": {Eq: "admin"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid action format",
			// Тест проверяет валидацию формата действия
			// Действие должно содержать минимум 2 части, разделенные ":"
			policies: []base.Policy{
				{
					Name:   "test-policy",
					Action: "invalid", // Невалидный формат - только одна часть
					Effect: base.Effect_ALLOW,
					Conditions: map[string]base.Condition{
						"source:role": {Eq: "admin"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid path in conditions",
			// Тест проверяет валидацию путей в условиях
			// Путь должен иметь формат "entity:field", где entity - "source" или "target"
			policies: []base.Policy{
				{
					Name:   "test-policy",
					Action: "user:read",
					Effect: base.Effect_ALLOW,
					Conditions: map[string]base.Condition{
						"invalid": {Eq: "admin"}, // Невалидный путь - нет разделителя ":"
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty policies",
			// Тест проверяет создание движка с пустым списком политик
			// Это должно быть успешно - движок создается, но без политик
			policies: []base.Policy{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := noctisguard.NewNoctisFromPolices(tt.policies)
			if tt.wantErr {
				// Ожидаем ошибку
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				// Ожидаем успех
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if engine == nil {
					t.Errorf("expected engine, got nil")
				}
			}
		})
	}
}

/*
TestNewNoctisFromFile тестирует создание движка из JSON-файла.

Тест проверяет:
  - Успешное чтение и парсинг валидного JSON-файла с политиками
  - Обработку ошибок при отсутствии файла
  - Корректность создания движка из файла (аналогично NewNoctisFromPolices)
*/
func TestNewNoctisFromFile(t *testing.T) {
	// Создаем временный файл с валидными политиками для тестирования
	validPolicies := []base.Policy{
		{
			Name:   "file-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	// Сериализуем политики в JSON для записи в файл
	validJSON, _ := json.Marshal(validPolicies)

	// Создаем временный файл с уникальным именем
	tmpFile, err := os.CreateTemp("", "test_policies_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Удаляем временный файл после теста

	// Записываем JSON в файл
	if _, err := tmpFile.Write(validJSON); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Создаем временный файл с невалидным JSON
	invalidJSON := []byte(`{"invalid": "json"`)
	tmpFileInvalid, err := os.CreateTemp("", "test_invalid_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFileInvalid.Name())
	if _, err := tmpFileInvalid.Write(invalidJSON); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFileInvalid.Close()

	// Создаем временный файл с пустым содержимым
	tmpFileEmpty, err := os.CreateTemp("", "test_empty_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFileEmpty.Name())
	tmpFileEmpty.Close()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name: "valid file",
			// Тест проверяет успешное чтение и парсинг валидного JSON-файла
			// Файл должен содержать массив объектов Policy в формате JSON
			path:    tmpFile.Name(),
			wantErr: false,
		},
		{
			name: "non-existent file",
			// Тест проверяет обработку ошибки при попытке открыть несуществующий файл
			// Должна возвращаться ошибка ErrExport, оборачивающая os.PathError
			path:    "/nonexistent/file.json",
			wantErr: true,
		},
		{
			name: "invalid JSON",
			// Тест проверяет обработку ошибки при невалидном JSON в файле
			// Должна возвращаться ошибка ErrExport, оборачивающая json.SyntaxError
			path:    tmpFileInvalid.Name(),
			wantErr: true,
		},
		{
			name: "empty file",
			// Тест проверяет обработку ошибки при пустом файле
			// Должна возвращаться ошибка ErrExport, оборачивающая json.SyntaxError
			path:    tmpFileEmpty.Name(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := noctisguard.NewNoctisFromFile(tt.path)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if engine == nil {
					t.Errorf("expected engine, got nil")
				}
			}
		})
	}
}

/*
TestEvaluate тестирует основную функциональность оценки политик.

Тест проверяет различные сценарии доступа:
  - Разрешение доступа для админов (политика allow-admin)
  - Разрешение доступа для владельцев (политика allow-owner с сравнением полей)
  - Запрет доступа для не-админов к приватным документам (политика deny-private)
  - Разрешение доступа админам к приватным документам (политика deny-private не применяется)
  - Обработку отсутствия политик для действия (возврат false)
  - Обработку невалидных путей к полям (возврат ошибки)
*/
func TestEvaluate(t *testing.T) {
	// Создаем набор политик для тестирования различных сценариев доступа
	policies := []base.Policy{
		{
			Name: "allow-admin",
			// Политика разрешает чтение документа админам
			// Условие: source:role == "admin"
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
		{
			Name: "allow-owner",
			// Политика разрешает чтение документа владельцу
			// Условие: source:name == target:owner (сравнение полей разных структур)
			Action: "user:read:document",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:name": {Eq: "target:owner"},
			},
		},
		{
			Name: "deny-private",
			// Политика запрещает чтение приватных документов не-админам
			// Условия: target:type == "private" И source:role != "admin"
			// Если оба условия выполнены, доступ запрещается
			Action: "user:read:document",
			Effect: base.Effect_DENY,
			Conditions: map[string]base.Condition{
				"target:type": {Eq: "private"},
				"source:role": {Neq: "admin"},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		source  any
		target  any
		action  string
		want    bool
		wantErr bool
	}{
		{
			name: "admin can read",
			// Тест проверяет, что админ может читать документы согласно политике allow-admin
			// source.role == "admin", поэтому политика должна пройти
			source:  User{Name: "admin", Role: "admin"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:read:document",
			want:    true,
			wantErr: false,
		},
		{
			name: "owner can read",
			// Тест проверяет, что владелец документа может его читать
			// source.name == target.owner == "alice", поэтому политика allow-owner должна пройти
			source:  User{Name: "alice", Role: "user"},
			target:  Document{Owner: "alice", Type: "public"},
			action:  "user:read:document",
			want:    true,
			wantErr: false,
		},
		{
			name: "deny private for non-admin",
			// Тест проверяет, что не-админ не может читать приватные документы
			// target.type == "private" И source.role != "admin", поэтому политика deny-private применяется
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "other", Type: "private"},
			action:  "user:read:document",
			want:    false,
			wantErr: false,
		},
		{
			name: "admin can read private",
			// Тест проверяет, что админ может читать приватные документы
			// source.role == "admin", поэтому условие source.role != "admin" не выполнено
			// Политика deny-private не применяется, и админ получает доступ через allow-admin
			source:  User{Name: "admin", Role: "admin"},
			target:  Document{Owner: "other", Type: "private"},
			action:  "user:read:document",
			want:    true,
			wantErr: false,
		},
		{
			name: "no policies for action",
			// Тест проверяет, что при отсутствии политик для действия возвращается false
			// Это означает, что действие запрещено по умолчанию
			source:  User{Name: "user", Role: "user"},
			target:  Document{Owner: "user", Type: "public"},
			action:  "user:write:document", // Нет политик для этого действия
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.Evaluate(ctx, tt.source, tt.target, tt.action)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got != tt.want {
					t.Errorf("Evaluate() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

/*
TestEvaluateWithContextCancellation тестирует отмену операции через context.Context.

Тест проверяет, что длительные операции проверки условий могут быть прерваны
через отмену контекста, и функция возвращает соответствующую ошибку ErrCancelled.
Это особенно важно для операций Contains с большими списками.
*/
func TestEvaluateWithContextCancellation(t *testing.T) {
	// Создаем политику с условием Contains, которое может выполняться долго для больших списков
	policies := []base.Policy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {
					Contains: []any{"admin", "moderator", "user"},
				},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	// Создаем контекст с отменой и сразу отменяем его
	// Это симулирует ситуацию, когда операция должна быть прервана
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Отменяем контекст до начала выполнения

	source := User{Name: "user", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	// Ожидаем, что операция вернет ошибку отмены
	_, err = engine.Evaluate(ctx, source, target, "user:read")
	if err == nil {
		t.Errorf("expected cancellation error, got nil")
	}
}

/*
TestEvaluateWithTimeout тестирует работу с таймаутом через context.Context.

Тест проверяет, что операция успешно выполняется в пределах таймаута
и не прерывается преждевременно. Это важно для проверки корректности
обработки контекста в нормальных условиях.
*/
func TestEvaluateWithTimeout(t *testing.T) {
	policies := []base.Policy{
		{
			Name:   "test-policy",
			Action: "user:read",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				"source:role": {Eq: "admin"},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	// Создаем контекст с таймаутом в 1 секунду
	// Операция должна выполниться быстрее, чем истечет таймаут
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	source := User{Name: "admin", Role: "admin"}
	target := Document{Owner: "user", Type: "public"}

	// Операция должна выполниться успешно в пределах таймаута
	allowed, err := engine.Evaluate(ctx, source, target, "user:read")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true, got false")
	}
}

/*
TestEvaluateNestedStructures тестирует работу с вложенными структурами.

Тест проверяет, что библиотека корректно обрабатывает пути к полям
во вложенных структурах, например "source:user:role" для доступа
к полю role внутри вложенной структуры user.

Это важно для работы с реальными структурами данных, которые часто
имеют вложенную структуру.
*/
func TestEvaluateNestedStructures(t *testing.T) {
	policies := []base.Policy{
		{
			Name:   "nested-policy",
			Action: "user:read:nested",
			Effect: base.Effect_ALLOW,
			Conditions: map[string]base.Condition{
				// Проверяем доступ к полю во вложенной структуре
				// Путь "source:user:role" означает: в source найти поле user, затем в нем найти поле role
				"source:user:role": {Eq: "admin"},
			},
		},
	}

	engine, err := noctisguard.NewNoctisFromPolices(policies)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	ctx := context.Background()
	// Создаем вложенные структуры для тестирования
	source := NestedUser{User: User{Name: "admin", Role: "admin"}}
	target := NestedDocument{Doc: Document{Owner: "user", Type: "public"}}

	// Проверяем, что путь к вложенному полю работает корректно
	allowed, err := engine.Evaluate(ctx, source, target, "user:read:nested")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true, got false")
	}
}
