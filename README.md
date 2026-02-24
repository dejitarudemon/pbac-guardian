# pbac-guardian

Current version: 0.4.1

[English](#english) | [Русский](#русский)

---

<a name="english"></a>
# English

## Overview

`pbac-guardian` is a lightweight, policy-based access control (PBAC) library for Go. It allows you to define access policies declaratively and check structures against these policies using a flexible system of conditions and effects.

## Features

- **Simple API**: Minimalistic and easy to use
- **Declarative Policies**: Define policies in JSON or programmatically
- **Flexible Conditions**: Support for equality, inequality, less-than, and contains operations
- **Context Support**: Built-in support for cancellation and timeouts via `context.Context`
- **Optional L1 Caching**: Optimized cache mechanism for improved performance with repeated field access
- **Thread-Safe**: Fully concurrent-safe implementation
- **No Dependencies**: Uses only the standard Go library
- **Type Safety**: Typed constants and validation

## Installation

```bash
go get github.com/dejitarudemon/pbac-guardian
```

## Quick Start

### 1. Define Your Structures

Tag your struct fields with `pbac-guardian` tags:

```go
type User struct {
    Name string `pbac-guardian:"name"`
    Role string `pbac-guardian:"role"`
    Age  int    `pbac-guardian:"age"`
}

type Document struct {
    Owner string   `pbac-guardian:"owner"`
    Type  string   `pbac-guardian:"type"`
    Tags  []string `pbac-guardian:"tags"`
}
```

### 2. Create Policies

```go
import (
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/base"
)

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
```

### 3. Create Engine and Evaluate

```go
import (
    "context"
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

// Create cache instance (optional - pass nil to disable caching)
casher := implemented.NewDefaultCasher()

// Create engine
engine, err := guardian.NewGuardianFromPolices(casher, policies)
if err != nil {
    panic(err)
}

// Create context
ctx := context.Background()

// Check access
user := User{Name: "alice", Role: "user"}
doc := Document{Owner: "alice", Type: "private"}

allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
if err != nil {
    // handle error
}

if allowed {
    // access granted
} else {
    // access denied
}
```

**Note**: The `casher` parameter is optional. Pass `nil` to disable caching if you don't need it. Caching is beneficial when the same fields are accessed multiple times within an evaluation session.

## Tutorials

### Tutorial 1: Basic Policy Evaluation

This tutorial shows how to create a simple policy and evaluate access.

```go
package main

import (
    "context"
    "fmt"
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/base"
)

type User struct {
    Name string `pbac-guardian:"name"`
    Role string `pbac-guardian:"role"`
}

type Document struct {
    Owner string `pbac-guardian:"owner"`
}

func main() {
    // Define policies
    policies := []base.Policy{
        {
            Name:   "admin-access",
            Action: "user:read:document",
            Effect: base.Effect_ALLOW,
            Conditions: map[string]base.Condition{
                "source:role": {
                    Eq: "admin",
                },
            },
        },
    }

    // Create engine (nil casher disables caching)
    engine, err := guardian.NewGuardianFromPolices(nil, policies)
    if err != nil {
        panic(err)
    }

    // Test cases
    ctx := context.Background()

    // Admin user - should be allowed
    admin := User{Name: "admin", Role: "admin"}
    doc := Document{Owner: "alice"}
    allowed, _ := engine.Evaluate(ctx, admin, doc, "user:read:document")
    fmt.Printf("Admin access: %v\n", allowed) // true

    // Regular user - should be denied
    user := User{Name: "bob", Role: "user"}
    allowed, _ = engine.Evaluate(ctx, user, doc, "user:read:document")
    fmt.Printf("User access: %v\n", allowed) // false
}
```

### Tutorial 2: Loading Policies from JSON

This tutorial shows how to load policies from a JSON file.

**policies.json:**
```json
[
  {
    "name": "admin-read",
    "action": "user:read:document",
    "effect": "allow",
    "conditions": {
      "source:role": {"eq": "admin"}
    }
  },
  {
    "name": "owner-read",
    "action": "user:read:document",
    "effect": "allow",
    "conditions": {
      "source:name": {"eq": "target:owner"}
    }
  }
]
```

**main.go:**
```go
package main

import (
    "context"
    "fmt"
    "github.com/dejitarudemon/pbac-guardian"
)

func main() {
    // Load policies from file (nil casher disables caching)
    engine, err := guardian.NewGuardianFromFile(nil, "policies.json")
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    user := User{Name: "alice", Role: "user"}
    doc := Document{Owner: "alice", Type: "private"}

    allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Access allowed: %v\n", allowed) // true (owner-read policy passed)
}
```

### Tutorial 3: Using Multiple Conditions

This tutorial shows how to combine multiple conditions in a single policy.

```go
policies := []base.Policy{
    {
        Name:   "age-restricted-read",
        Action: "user:read:document",
        Effect: base.Effect_ALLOW,
        Conditions: map[string]base.Condition{
            "source:age": {
                Lt: 18, // age < 18
            },
            "target:tags": {
                Contains: []any{"public"}, // "public" in tags
            },
        },
    },
}
```

### Tutorial 4: DENY Policies

DENY policies have priority. If a DENY policy's conditions are met, access is denied even if ALLOW policies pass.

```go
policies := []base.Policy{
    {
        Name:   "block-minors",
        Action: "user:read:document",
        Effect: base.Effect_DENY,
        Conditions: map[string]base.Condition{
            "source:age": {
                Lt: 18, // age < 18
            },
            "target:tags": {
                Contains: []any{"adult-only"},
            },
        },
    },
    {
        Name:   "allow-all",
        Action: "user:read:document",
        Effect: base.Effect_ALLOW,
        Conditions: map[string]base.Condition{
            "source:role": {
                Eq: "user",
            },
        },
    },
}
```

### Tutorial 5: Context Cancellation

This tutorial shows how to use context for cancellation and timeouts.

```go
import (
    "context"
    "time"
)

func main() {
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Evaluate with timeout
    allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
    if err != nil {
        if err == base.ErrCancelled {
            fmt.Println("Operation cancelled")
        } else {
            panic(err)
        }
    }
}
```

### Tutorial 6: Using L1 Cache for Performance

This tutorial shows how to enable the optional L1 cache for improved performance when fields are accessed multiple times.

```go
import (
    "context"
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/base"
    "github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

func main() {
    policies := []base.Policy{
        {
            Name:   "admin-read",
            Action: "user:read:document",
            Effect: base.Effect_ALLOW,
            Conditions: map[string]base.Condition{
                "source:role": {Eq: "admin"},
            },
        },
    }

    // Enable caching for better performance with repeated field access
    casher := implemented.NewDefaultCasher()
    engine, err := guardian.NewGuardianFromPolices(casher, policies)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    user := User{Name: "admin", Role: "admin"}
    doc := Document{Owner: "alice", Type: "public"}

    // Cache will store field values during evaluation
    // Subsequent accesses to the same fields within the session will be faster
    allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Access allowed: %v\n", allowed)
}
```

**When to use cache:**
- Multiple policies accessing the same fields
- Fields accessed multiple times within an evaluation session
- Applications prioritizing reduced allocations

**When to disable cache (pass nil):**
- Simple single-policy evaluations
- Fields accessed only once per session
- Maximum performance for single-access scenarios

**Performance Benchmarks:**

The cache provides significant performance improvements in production scenarios:

| Scenario | Policies | Field Accesses | Time Savings | Memory Savings | Allocation Reduction |
|----------|----------|----------------|--------------|----------------|---------------------|
| **20P-3A** | 20 | 3 | **14%** faster | +2.7% | **47%** fewer |
| **30P-5A** | 30 | 5 | **32%** faster | **25%** less | **59%** fewer |
| **30P-10A** | 30 | 10 | **41%** faster | **36%** less | **63%** fewer |
| **Mixed** | 30 | 5+3+7 | **24%** faster | **23%** less | **53%** fewer |

*Benchmarks conducted with 10+ actions and 20-30 policies per evaluation. Cache becomes beneficial at 3+ field accesses within a single action evaluation.*

### Tutorial 7: Custom Types with Comparable Interface

For custom types that need special comparison logic, implement the `Comparable` interface.

```go
type User struct {
    Name string
    Age  int
}

func (u User) Compare(other any) (int, bool) {
    o, ok := other.(User)
    if !ok {
        return 0, false
    }
    if u.Age < o.Age {
        return -1, true
    }
    if u.Age > o.Age {
        return 1, true
    }
    return 0, true
}
```

## API Reference

### Guardian Engine

#### `NewGuardianFromPolices(casher base.Casher, policies []base.Policy) (*Guardian, error)`

Creates a new Guardian instance from a list of policies. The `casher` parameter is optional - pass `nil` to disable caching, or use `implemented.NewDefaultCasher()` to enable optimized L1 caching.

#### `NewGuardianFromFile(casher base.Casher, path string) (*Guardian, error)`

Creates a new Guardian instance from a JSON file. The `casher` parameter is optional - pass `nil` to disable caching.

#### `Evaluate(ctx context.Context, source, target any, action string) (bool, error)`

Evaluates access for the given action. Returns `true` if access is allowed, `false` otherwise. Each evaluation uses a unique session ID for cache isolation.

### Policy Structure

```go
type Policy struct {
    Name       string               `json:"name"`
    Action     string               `json:"action"`
    Effect     Effect               `json:"effect"`
    Conditions map[string]Condition `json:"conditions"`
}
```

### Condition Types

- **Eq**: Equality check (`left == right`)
- **Neq**: Inequality check (`left != right`)
- **Lt**: Less than check (`left < right`)
- **Contains**: Checks if `left` is in `right` (slice)

### Effects

- **Effect_ALLOW**: Allow action if conditions are met
- **Effect_DENY**: Deny action if conditions are met (has priority)

## Error Handling

The library provides typed errors:

- `ErrDuplicateName`: Policy name already exists
- `ErrExport`: Error exporting policies
- `ErrEvaluate`: Error evaluating policies
- `ErrCancelled`: Operation cancelled via context
- `ErrInvalidPath`: Invalid field path
- `ErrInvalidType`: Invalid type
- `ErrUncomparable`: Cannot compare values

---

<a name="русский"></a>
# Русский

## Обзор

`pbac-guardian` — это легковесная библиотека для контроля доступа на основе политик (PBAC) для Go. Она позволяет декларативно определять политики доступа и проверять структуры на соответствие этим политикам с использованием гибкой системы условий и эффектов.

## Возможности

- **Простой API**: Минималистичный и простой в использовании
- **Декларативные политики**: Определение политик в JSON или программно
- **Гибкие условия**: Поддержка операций равенства, неравенства, меньше и содержит
- **Поддержка контекста**: Встроенная поддержка отмены и таймаутов через `context.Context`
- **Опциональное L1-кеширование**: Оптимизированный механизм кеша для улучшения производительности при повторном доступе к полям
- **Потокобезопасность**: Полностью потокобезопасная реализация
- **Без зависимостей**: Использует только стандартную библиотеку Go
- **Безопасность типов**: Типизированные константы и валидация

## Установка

```bash
go get github.com/dejitarudemon/pbac-guardian
```

## Быстрый старт

### 1. Определите свои структуры

Помечайте поля структур тегами `pbac-guardian`:

```go
type User struct {
    Name string `pbac-guardian:"name"`
    Role string `pbac-guardian:"role"`
    Age  int    `pbac-guardian:"age"`
}

type Document struct {
    Owner string   `pbac-guardian:"owner"`
    Type  string   `pbac-guardian:"type"`
    Tags  []string `pbac-guardian:"tags"`
}
```

### 2. Создайте политики

```go
import (
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/base"
)

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
```

### 3. Создайте движок и проверьте доступ

```go
import (
    "context"
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

// Создание экземпляра кеша (опционально - передайте nil для отключения кеширования)
casher := implemented.NewDefaultCasher()

// Создание движка
engine, err := guardian.NewGuardianFromPolices(casher, policies)
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
    // обработка ошибки
}

if allowed {
    // доступ разрешен
} else {
    // доступ запрещен
}
```

**Примечание**: Параметр `casher` опционален. Передайте `nil` для отключения кеширования, если оно не нужно. Кеширование полезно, когда одни и те же поля обращаются несколько раз в рамках сессии оценки.

## Туториалы

### Туториал 1: Базовая оценка политик

Этот туториал показывает, как создать простую политику и проверить доступ.

```go
package main

import (
    "context"
    "fmt"
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/base"
)

type User struct {
    Name string `pbac-guardian:"name"`
    Role string `pbac-guardian:"role"`
}

type Document struct {
    Owner string `pbac-guardian:"owner"`
}

func main() {
    // Определение политик
    policies := []base.Policy{
        {
            Name:   "admin-access",
            Action: "user:read:document",
            Effect: base.Effect_ALLOW,
            Conditions: map[string]base.Condition{
                "source:role": {
                    Eq: "admin",
                },
            },
        },
    }

    // Создание движка (nil casher отключает кеширование)
    engine, err := guardian.NewGuardianFromPolices(nil, policies)
    if err != nil {
        panic(err)
    }

    // Тестовые случаи
    ctx := context.Background()

    // Пользователь-админ - должен быть разрешен
    admin := User{Name: "admin", Role: "admin"}
    doc := Document{Owner: "alice"}
    allowed, _ := engine.Evaluate(ctx, admin, doc, "user:read:document")
    fmt.Printf("Доступ админа: %v\n", allowed) // true

    // Обычный пользователь - должен быть запрещен
    user := User{Name: "bob", Role: "user"}
    allowed, _ = engine.Evaluate(ctx, user, doc, "user:read:document")
    fmt.Printf("Доступ пользователя: %v\n", allowed) // false
}
```

### Туториал 2: Загрузка политик из JSON

Этот туториал показывает, как загрузить политики из JSON файла.

**policies.json:**
```json
[
  {
    "name": "admin-read",
    "action": "user:read:document",
    "effect": "allow",
    "conditions": {
      "source:role": {"eq": "admin"}
    }
  },
  {
    "name": "owner-read",
    "action": "user:read:document",
    "effect": "allow",
    "conditions": {
      "source:name": {"eq": "target:owner"}
    }
  }
]
```

**main.go:**
```go
package main

import (
    "context"
    "fmt"
    "github.com/dejitarudemon/pbac-guardian"
)

func main() {
    // Загрузка политик из файла (nil casher отключает кеширование)
    engine, err := guardian.NewGuardianFromFile(nil, "policies.json")
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    user := User{Name: "alice", Role: "user"}
    doc := Document{Owner: "alice", Type: "private"}

    allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Доступ разрешен: %v\n", allowed) // true (политика owner-read прошла)
}
```

### Туториал 3: Использование множественных условий

Этот туториал показывает, как комбинировать несколько условий в одной политике.

```go
policies := []base.Policy{
    {
        Name:   "age-restricted-read",
        Action: "user:read:document",
        Effect: base.Effect_ALLOW,
        Conditions: map[string]base.Condition{
            "source:age": {
                Lt: 18, // возраст < 18
            },
            "target:tags": {
                Contains: []any{"public"}, // "public" в тегах
            },
        },
    },
}
```

### Туториал 4: Политики DENY

Политики DENY имеют приоритет. Если условия политики DENY выполнены, доступ запрещается, даже если политики ALLOW прошли проверку.

```go
policies := []base.Policy{
    {
        Name:   "block-minors",
        Action: "user:read:document",
        Effect: base.Effect_DENY,
        Conditions: map[string]base.Condition{
            "source:age": {
                Lt: 18, // возраст < 18
            },
            "target:tags": {
                Contains: []any{"adult-only"},
            },
        },
    },
    {
        Name:   "allow-all",
        Action: "user:read:document",
        Effect: base.Effect_ALLOW,
        Conditions: map[string]base.Condition{
            "source:role": {
                Eq: "user",
            },
        },
    },
}
```

### Туториал 5: Отмена через контекст

Этот туториал показывает, как использовать контекст для отмены и таймаутов.

```go
import (
    "context"
    "time"
)

func main() {
    // Создание контекста с таймаутом
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Оценка с таймаутом
    allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
    if err != nil {
        if err == base.ErrCancelled {
            fmt.Println("Операция отменена")
        } else {
            panic(err)
        }
    }
}
```

### Туториал 6: Использование L1-кеша для производительности

Этот туториал показывает, как включить опциональный L1-кеш для улучшения производительности при повторном доступе к полям.

```go
import (
    "context"
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/base"
    "github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

func main() {
    policies := []base.Policy{
        {
            Name:   "admin-read",
            Action: "user:read:document",
            Effect: base.Effect_ALLOW,
            Conditions: map[string]base.Condition{
                "source:role": {Eq: "admin"},
            },
        },
    }

    // Включить кеширование для лучшей производительности при повторном доступе к полям
    casher := implemented.NewDefaultCasher()
    engine, err := guardian.NewGuardianFromPolices(casher, policies)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    user := User{Name: "admin", Role: "admin"}
    doc := Document{Owner: "alice", Type: "public"}

    // Кеш будет хранить значения полей во время оценки
    // Последующие обращения к тем же полям в рамках сессии будут быстрее
    allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Доступ разрешен: %v\n", allowed)
}
```

**Когда использовать кеш:**
- Несколько политик обращаются к одним и тем же полям
- Поля обращаются несколько раз в рамках сессии оценки
- Приложения, приоритизирующие уменьшенные аллокации

**Когда отключать кеш (передать nil):**
- Простые оценки одной политики
- Поля обращаются только один раз за сессию
- Максимальная производительность для сценариев с однократным доступом

**Результаты бенчмарков:**

Кеш обеспечивает значительное улучшение производительности в production-сценариях:

| Сценарий | Политики | Обращения к полю | Экономия времени | Экономия памяти | Уменьшение аллокаций |
|----------|----------|------------------|------------------|------------------|----------------------|
| **20P-3A** | 20 | 3 | **14%** быстрее | +2.7% | **47%** меньше |
| **30P-5A** | 30 | 5 | **32%** быстрее | **25%** меньше | **59%** меньше |
| **30P-10A** | 30 | 10 | **41%** быстрее | **36%** меньше | **63%** меньше |
| **Mixed** | 30 | 5+3+7 | **24%** быстрее | **23%** меньше | **53%** меньше |

*Бенчмарки проведены с 10+ действиями и 20-30 политиками на оценку. Кеш становится выгодным при 3+ обращениях к полю в рамках одной оценки действия.*

### Туториал 7: Кастомные типы с интерфейсом Comparable

Для кастомных типов, требующих специальной логики сравнения, реализуйте интерфейс `Comparable`.

```go
type User struct {
    Name string
    Age  int
}

func (u User) Compare(other any) (int, bool) {
    o, ok := other.(User)
    if !ok {
        return 0, false
    }
    if u.Age < o.Age {
        return -1, true
    }
    if u.Age > o.Age {
        return 1, true
    }
    return 0, true
}
```

## Справочник API

### Движок Guardian

#### `NewGuardianFromPolices(casher base.Casher, policies []base.Policy) (*Guardian, error)`

Создает новый экземпляр Guardian из списка политик. Параметр `casher` опционален - передайте `nil` для отключения кеширования или используйте `implemented.NewDefaultCasher()` для включения оптимизированного L1-кеширования.

#### `NewGuardianFromFile(casher base.Casher, path string) (*Guardian, error)`

Создает новый экземпляр Guardian из JSON файла. Параметр `casher` опционален - передайте `nil` для отключения кеширования.

#### `Evaluate(ctx context.Context, source, target any, action string) (bool, error)`

Оценивает доступ для указанного действия. Возвращает `true`, если доступ разрешен, `false` в противном случае. Каждая оценка использует уникальный идентификатор сессии для изоляции кеша.

### Структура Policy

```go
type Policy struct {
    Name       string               `json:"name"`
    Action     string               `json:"action"`
    Effect     Effect               `json:"effect"`
    Conditions map[string]Condition `json:"conditions"`
}
```

### Типы условий

- **Eq**: Проверка равенства (`left == right`)
- **Neq**: Проверка неравенства (`left != right`)
- **Lt**: Проверка меньше (`left < right`)
- **Contains**: Проверяет, находится ли `left` в `right` (slice)

### Эффекты

- **Effect_ALLOW**: Разрешить действие, если условия выполнены
- **Effect_DENY**: Запретить действие, если условия выполнены (имеет приоритет)

## Обработка ошибок

Библиотека предоставляет типизированные ошибки:

- `ErrDuplicateName`: Имя политики уже существует
- `ErrExport`: Ошибка экспорта политик
- `ErrEvaluate`: Ошибка оценки политик
- `ErrCancelled`: Операция отменена через контекст
- `ErrInvalidPath`: Невалидный путь до поля
- `ErrInvalidType`: Невалидный тип
- `ErrUncomparable`: Невозможно сравнить значения

