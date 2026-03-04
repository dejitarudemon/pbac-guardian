[![Go Reference](https://pkg.go.dev/badge/github.com/username/project.svg)](https://pkg.go.dev/github.com/dejitarudemon/pbac-guardian)
[![Go Report Card](https://goreportcard.com/badge/github.com/dejitarudemon/pbac-guardian)](https://goreportcard.com/report/github.com/dejitarudemon/pbac-guardian)
[![GitHub License](https://img.shields.io/github/license/dejitarudemon/pbac-guardian)](https://github.com/dejitarudemon/pbac-guardian/blob/main/LICENSE)

# pbac-guardian

Current version: 1.1.0

[English](#english)

[Russian](#russian)

---

<a name="english"></a>
# English

## Overview

`pbac-guardian` is a lightweight, policy-based access control (PBAC) library for Go. It allows you to define access policies declaratively and check structures against these policies using a flexible system of conditions and effects.

## Features

- **Simple API**: Minimalistic and easy to use
- **Declarative Policies**: Define policies in JSON or programmatically
- **Flexible Conditions**: Support for equality, inequality, less-than, greater-than, less-than-or-equal, greater-than-or-equal, and contains operations
- **Time Support**: Built-in support for time values and time modifiers (e.g., "time:now", "time:now:1|day")
- **Environment Variables**: Support for environment variables in policy conditions (e.g., "env:VAR_NAME")
- **Context Support**: Built-in support for cancellation and timeouts via `context.Context`
- **Optional L1 Caching**: Optimized cache mechanism for improved performance with repeated field access
- **Thread-Safe**: Fully concurrent-safe implementation
- **Minimal Dependencies**: Uses only one external dependency (`github.com/google/uuid` for session ID generation)
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

policies := []base.RawPolicy{
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

    // Create engine with default config
    config := base.Config{
        ConditionsMap:        nil, // use default condition functions
        CashDisableThreShold: 3,   // disable cache for fields accessed < 3 times
    }
    engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
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
    policies := []base.RawPolicy{
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

    // Create engine with default config (nil casher disables caching)
    config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
    engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
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
    // Load policies from file with default config (nil casher disables caching)
    config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
    engine, err := guardian.NewGuardianFromFile(nil, "policies.json", config)
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
policies := []base.RawPolicy{
    {
        Name:   "age-restricted-read",
        Action: "user:read:document",
        Effect: base.Effect_ALLOW,
        Conditions: map[string]base.Condition{
            "source:age": {
                Lt: 18, // age < 18
            },
            "target:tags": {
                In: []any{"public"}, // "public" in tags
            },
        },
    },
}
```

### Tutorial 4: DENY Policies

DENY policies have priority. If a DENY policy's conditions are met, access is denied even if ALLOW policies pass.

```go
policies := []base.RawPolicy{
    {
        Name:   "block-minors",
        Action: "user:read:document",
        Effect: base.Effect_DENY,
        Conditions: map[string]base.Condition{
            "source:age": {
                Lt: 18, // age < 18
            },
            "target:tags": {
                In: []any{"adult-only"},
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
    policies := []base.RawPolicy{
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
    config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
    engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
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

*Note: Cache becomes beneficial at 3+ field accesses within a single action evaluation.*

### Tutorial 7: Custom Condition Functions

This tutorial shows how to customize condition functions using the `config` parameter.

```go
import (
    "context"
    "strings"
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/base"
    "github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

func main() {
    // Create custom condition function
    customEqFunc := func(ctx context.Context, left, right any) (bool, error) {
        // Custom logic: case-insensitive string comparison
        if ctx.Err() != nil {
            return false, base.ErrCancelled
        }
        leftStr, ok1 := left.(string)
        rightStr, ok2 := right.(string)
        if !ok1 || !ok2 {
            return false, base.NewErrInvalidType("string", left)
        }
        return strings.EqualFold(leftStr, rightStr), nil
    }

    // Create custom ConditionsMap
    customConditionsMap := &base.ConditionsMap{
        In:  implemented.DefaultConditionsMap.In,
        Eq:       customEqFunc, // Use custom function
        Neq:      implemented.DefaultConditionsMap.Neq,
        Lt:       implemented.DefaultConditionsMap.Lt,
        Gt:       implemented.DefaultConditionsMap.Gt,
        Le:       implemented.DefaultConditionsMap.Le,
        Ge:       implemented.DefaultConditionsMap.Ge,
    }

    // Create config with custom ConditionsMap
    config := base.Config{
        ConditionsMap:        customConditionsMap,
        CashDisableThreShold: 3,
    }

    policies := []base.RawPolicy{
        {
            Name:   "case-insensitive-admin",
            Action: "user:read:document",
            Effect: base.Effect_ALLOW,
            Conditions: map[string]base.Condition{
                "source:role": {Eq: "ADMIN"}, // Will match "admin", "Admin", "ADMIN", etc.
            },
        },
    }

    // Create engine with custom condition functions
    engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    user := User{Name: "alice", Role: "admin"} // lowercase "admin"
    doc := Document{Owner: "alice", Type: "public"}

    // Custom function will match case-insensitively
    allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Access allowed: %v\n", allowed) // true (case-insensitive match)
}
```

**When to use custom functions:**
- Need special comparison logic (case-insensitive, fuzzy matching, etc.)
- Integration with external validation services
- Custom business rules for condition evaluation

**When to use defaults (pass nil):**
- Standard equality, inequality, less-than, greater-than, less-than-or-equal, greater-than-or-equal, and contains operations
- Most common use cases
- When performance is critical (default functions are optimized)

### Tutorial 8: Environment Variables and Time Values

This tutorial shows how to use environment variables and time values in policy conditions.

#### Environment Variables

You can use environment variables in conditions using the `"env:VARIABLE_NAME"` path format:

```go
import (
    "os"
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/base"
)

// Set environment variable
os.Setenv("MIN_AGE", "18")
defer os.Unsetenv("MIN_AGE")

policies := []base.RawPolicy{
    {
        Name:   "age-check",
        Action: "user:read:document",
        Effect: base.Effect_ALLOW,
        Conditions: map[string]base.Condition{
            "source:age": {
                Gt: "env:MIN_AGE", // source:age > env:MIN_AGE
            },
        },
    },
}
```

#### Time Values

You can use time values in conditions using the `"time:now"` path format:

```go
import (
    "time"
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/base"
)

policies := []base.RawPolicy{
    {
        Name:   "recent-document",
        Action: "user:read:document",
        Effect: base.Effect_ALLOW,
        Conditions: map[string]base.Condition{
            "source:created_at": {
                Gt: "time:now:-7|day", // source:created_at > time:now - 7 days
            },
        },
    },
}
```

Time modifiers:
- `"time:now"` - current time
- `"time:now:1|day"` - current time + 1 day
- `"time:now:2|hour"` - current time + 2 hours
- `"time:now:30|minute"` - current time + 30 minutes
- Supported units: `day`, `hour`, `minute`, `second`, `milisecond`

### Tutorial 9: Custom Types with Comparable Interface

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

#### `NewGuardianFromPolices(casher cashing.Casher, policies []base.RawPolicy, config base.Config) (*Guardian, error)`

Creates a new Guardian instance from a list of policies. 
- The `casher` parameter is optional - pass `nil` to disable caching, or use `implemented.NewDefaultCasher()` to enable optimized L1 caching.
- The `policies` parameter accepts an array of `[]base.RawPolicy` (structures with public fields for JSON unmarshaling)
- The `config` parameter contains condition functions map and cache disable threshold. If `config.ConditionsMap` is `nil`, default condition functions will be used.

#### `NewGuardianFromFile(casher cashing.Casher, path string, config base.Config) (*Guardian, error)`

Creates a new Guardian instance from a JSON file.
- The `casher` parameter is optional - pass `nil` to disable caching.
- The `config` parameter contains condition functions map and cache disable threshold. If `config.ConditionsMap` is `nil`, default condition functions will be used.

#### `Evaluate(ctx context.Context, source, target any, action string) (bool, error)`

Evaluates access for the given action. Returns `true` if access is allowed, `false` otherwise. Each evaluation uses a unique session ID for cache isolation.

### Policy Structures

#### RawPolicy

`RawPolicy` is used for JSON unmarshaling and creating policies programmatically. It contains public fields that can be serialized/deserialized.

```go
type RawPolicy struct {
    Name       string               `json:"name"`
    Action     string               `json:"action"`
    Effect     Effect               `json:"effect"`
    Conditions map[string]Condition `json:"conditions"`
}
```

#### Policy

`Policy` is an internal structure with private fields. Policy instances are created automatically by the Guardian engine from `RawPolicy` structures using `NewPolicy` constructor. You should not create `Policy` instances directly.

Public methods:
- `Evaluate(ctx context.Context, source, target any, action string, sessionID string) (bool, error)` - applies the policy
- `IsValid() error` - validates the policy
- `Effect() Effect` - returns the policy effect

### Condition Types

- **Eq**: Equality check (`left == right`)
- **Neq**: Inequality check (`left != right`)
- **Lt**: Less than check (`left < right`)
- **Gt**: Greater than check (`left > right`)
- **Le**: Less than or equal check (`left <= right`)
- **Ge**: Greater than or equal check (`left >= right`)
- **In**: Checks if `left` is in `right` (slice)

### Path Types

Paths in conditions can reference:
- **Structure fields**: `"source:field"`, `"target:field"` - fields from source/target structures
- **Environment variables**: `"env:VARIABLE_NAME"` - values from environment variables
- **Time values**: `"time:now"` - current time, `"time:now:1|day"` - current time with modifiers
  - Supported modifiers: `day`, `hour`, `minute`, `second`, `milisecond`
  - Format: `"time:now:value|unit"` (e.g., `"time:now:2|hour"`)

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

<a name="russian"></a>
# Russian

## Обзор

`pbac-guardian` — это легковесная библиотека для контроля доступа на основе политик (PBAC) для Go. Она позволяет декларативно определять политики доступа и проверять структуры на соответствие этим политикам с использованием гибкой системы условий и эффектов.

## Возможности

- **Простой API**: Минималистичный и простой в использовании
- **Декларативные политики**: Определение политик в JSON или программно
- **Гибкие условия**: Поддержка операций равенства, неравенства, меньше, больше и содержит
- **Поддержка времени**: Встроенная поддержка значений времени и модификаторов времени (например, "time:now", "time:now:1|day")
- **Переменные окружения**: Поддержка переменных окружения в условиях политик (например, "env:VAR_NAME")
- **Поддержка контекста**: Встроенная поддержка отмены и таймаутов через `context.Context`
- **Опциональное L1-кеширование**: Оптимизированный механизм кеша для улучшения производительности при повторном доступе к полям
- **Потокобезопасность**: Полностью потокобезопасная реализация
- **Минимальные зависимости**: Использует только одну внешнюю зависимость (`github.com/google/uuid` для генерации session ID)
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

policies := []base.RawPolicy{
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

    // Создание движка с конфигурацией по умолчанию
    config := base.Config{
        ConditionsMap:        nil, // использовать функции условий по умолчанию
        CashDisableThreShold: 3,   // отключить кеш для полей, к которым обращаются < 3 раз
    }
    engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
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
    policies := []base.RawPolicy{
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

    // Создание движка с конфигурацией по умолчанию (nil casher отключает кеширование)
    config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
    engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
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
    // Загрузка политик из файла с конфигурацией по умолчанию (nil casher отключает кеширование)
    config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
    engine, err := guardian.NewGuardianFromFile(nil, "policies.json", config)
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
policies := []base.RawPolicy{
    {
        Name:   "age-restricted-read",
        Action: "user:read:document",
        Effect: base.Effect_ALLOW,
        Conditions: map[string]base.Condition{
            "source:age": {
                Lt: 18, // возраст < 18
            },
            "target:tags": {
                In: []any{"public"}, // "public" в тегах
            },
        },
    },
}
```

### Туториал 4: Политики DENY

Политики DENY имеют приоритет. Если условия политики DENY выполнены, доступ запрещается, даже если политики ALLOW прошли проверку.

```go
policies := []base.RawPolicy{
    {
        Name:   "block-minors",
        Action: "user:read:document",
        Effect: base.Effect_DENY,
        Conditions: map[string]base.Condition{
            "source:age": {
                Lt: 18, // возраст < 18
            },
            "target:tags": {
                In: []any{"adult-only"},
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
    policies := []base.RawPolicy{
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
    config := base.Config{ConditionsMap: nil, CashDisableThreShold: 3}
    engine, err := guardian.NewGuardianFromPolices(casher, policies, config)
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

*Примечание: Кеш становится выгодным при 3+ обращениях к полю в рамках одной оценки действия.*

### Туториал 7: Кастомные функции условий

Этот туториал показывает, как кастомизировать функции условий с помощью параметра `config`.

```go
import (
    "context"
    "strings"
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/base"
    "github.com/dejitarudemon/pbac-guardian/internal/implemented"
)

func main() {
    // Создание кастомной функции условия
    customEqFunc := func(ctx context.Context, left, right any) (bool, error) {
        // Кастомная логика: сравнение строк без учета регистра
        if ctx.Err() != nil {
            return false, base.ErrCancelled
        }
        leftStr, ok1 := left.(string)
        rightStr, ok2 := right.(string)
        if !ok1 || !ok2 {
            return false, base.NewErrInvalidType("string", left)
        }
        return strings.EqualFold(leftStr, rightStr), nil
    }

    // Создание кастомной ConditionsMap
    customConditionsMap := &base.ConditionsMap{
        In:  implemented.DefaultConditionsMap.In,
        Eq:       customEqFunc, // Использовать кастомную функцию
        Neq:      implemented.DefaultConditionsMap.Neq,
        Lt:       implemented.DefaultConditionsMap.Lt,
        Gt:       implemented.DefaultConditionsMap.Gt,
        Le:       implemented.DefaultConditionsMap.Le,
        Ge:       implemented.DefaultConditionsMap.Ge,
    }

    // Создание config с кастомной ConditionsMap
    config := base.Config{
        ConditionsMap:        customConditionsMap,
        CashDisableThreShold: 3,
    }

    policies := []base.RawPolicy{
        {
            Name:   "case-insensitive-admin",
            Action: "user:read:document",
            Effect: base.Effect_ALLOW,
            Conditions: map[string]base.Condition{
                "source:role": {Eq: "ADMIN"}, // Будет совпадать с "admin", "Admin", "ADMIN" и т.д.
            },
        },
    }

    // Создание движка с кастомными функциями условий
    engine, err := guardian.NewGuardianFromPolices(nil, policies, config)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    user := User{Name: "alice", Role: "admin"} // строчные буквы "admin"
    doc := Document{Owner: "alice", Type: "public"}

    // Кастомная функция будет совпадать без учета регистра
    allowed, err := engine.Evaluate(ctx, user, doc, "user:read:document")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Доступ разрешен: %v\n", allowed) // true (совпадение без учета регистра)
}
```

**Когда использовать кастомные функции:**
- Нужна специальная логика сравнения (без учета регистра, нечеткое совпадение и т.д.)
- Интеграция с внешними сервисами валидации
- Кастомные бизнес-правила для оценки условий

**Когда использовать значения по умолчанию (передать nil):**
- Стандартные операции равенства, неравенства, меньше, больше и содержит
- Большинство распространенных случаев использования
- Когда производительность критична (функции по умолчанию оптимизированы)

### Туториал 8: Переменные окружения и значения времени

Этот туториал показывает, как использовать переменные окружения и значения времени в условиях политик.

#### Переменные окружения

Вы можете использовать переменные окружения в условиях, используя формат пути `"env:VARIABLE_NAME"`:

```go
import (
    "os"
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/base"
)

// Установить переменную окружения
os.Setenv("MIN_AGE", "18")
defer os.Unsetenv("MIN_AGE")

policies := []base.RawPolicy{
    {
        Name:   "age-check",
        Action: "user:read:document",
        Effect: base.Effect_ALLOW,
        Conditions: map[string]base.Condition{
            "source:age": {
                Gt: "env:MIN_AGE", // source:age > env:MIN_AGE
            },
        },
    },
}
```

#### Значения времени

Вы можете использовать значения времени в условиях, используя формат пути `"time:now"`:

```go
import (
    "time"
    "github.com/dejitarudemon/pbac-guardian"
    "github.com/dejitarudemon/pbac-guardian/internal/base"
)

policies := []base.RawPolicy{
    {
        Name:   "recent-document",
        Action: "user:read:document",
        Effect: base.Effect_ALLOW,
        Conditions: map[string]base.Condition{
            "source:created_at": {
                Gt: "time:now:-7|day", // source:created_at > time:now - 7 days
            },
        },
    },
}
```

Модификаторы времени:
- `"time:now"` - текущее время
- `"time:now:1|day"` - текущее время + 1 день
- `"time:now:2|hour"` - текущее время + 2 часа
- `"time:now:30|minute"` - текущее время + 30 минут
- Поддерживаемые единицы: `day`, `hour`, `minute`, `second`, `milisecond`

### Туториал 9: Кастомные типы с интерфейсом Comparable

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

#### `NewGuardianFromPolices(casher cashing.Casher, policies []base.RawPolicy, config base.Config) (*Guardian, error)`

Создает новый экземпляр Guardian из списка политик.
- Параметр `casher` опционален - передайте `nil` для отключения кеширования или используйте `implemented.NewDefaultCasher()` для включения оптимизированного L1-кеширования.
- Параметр `policies` принимает массив `[]base.RawPolicy` (структуры с публичными полями для JSON unmarshaling)
- Параметр `config` содержит карту функций условий и порог отключения кеша. Если `config.ConditionsMap` равен `nil`, будут использованы функции условий по умолчанию.

#### `NewGuardianFromFile(casher cashing.Casher, path string, config base.Config) (*Guardian, error)`

Создает новый экземпляр Guardian из JSON файла.
- Параметр `casher` опционален - передайте `nil` для отключения кеширования.
- Параметр `config` содержит карту функций условий и порог отключения кеша. Если `config.ConditionsMap` равен `nil`, будут использованы функции условий по умолчанию.

#### `Evaluate(ctx context.Context, source, target any, action string) (bool, error)`

Оценивает доступ для указанного действия. Возвращает `true`, если доступ разрешен, `false` в противном случае. Каждая оценка использует уникальный идентификатор сессии для изоляции кеша.

### Структуры Policy

#### RawPolicy

`RawPolicy` используется для JSON unmarshaling и программного создания политик. Содержит публичные поля, которые могут быть сериализованы/десериализованы.

```go
type RawPolicy struct {
    Name       string               `json:"name"`
    Action     string               `json:"action"`
    Effect     Effect               `json:"effect"`
    Conditions map[string]Condition `json:"conditions"`
}
```

#### Policy

`Policy` - это внутренняя структура с приватными полями. Экземпляры `Policy` создаются автоматически движком Guardian из структур `RawPolicy` с помощью конструктора `NewPolicy`. Не следует создавать экземпляры `Policy` напрямую.

Публичные методы:
- `Evaluate(ctx context.Context, source, target any, action string, sessionID string) (bool, error)` - применяет политику
- `IsValid() error` - валидирует политику
- `Effect() Effect` - возвращает эффект политики

### Типы условий

- **Eq**: Проверка равенства (`left == right`)
- **Neq**: Проверка неравенства (`left != right`)
- **Lt**: Проверка меньше (`left < right`)
- **Gt**: Проверка больше (`left > right`)
- **Le**: Проверка меньше или равно (`left <= right`)
- **Ge**: Проверка больше или равно (`left >= right`)
- **In**: Проверяет, находится ли `left` в `right` (slice)

### Типы путей

Пути в условиях могут ссылаться на:
- **Поля структур**: `"source:field"`, `"target:field"` - поля из структур source/target
- **Переменные окружения**: `"env:VARIABLE_NAME"` - значения из переменных окружения
- **Значения времени**: `"time:now"` - текущее время, `"time:now:1|day"` - текущее время с модификаторами
  - Поддерживаемые модификаторы: `day`, `hour`, `minute`, `second`, `milisecond`
  - Формат: `"time:now:value|unit"` (например, `"time:now:2|hour"`)

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

