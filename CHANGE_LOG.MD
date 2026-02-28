## v 0.4.3

### Русский

**Изменено:**

- Изменена сигнатура метода `Policy.Evaluate`: удален параметр `cash`, кеш теперь хранится в поле `Cash` структуры `Policy`
- Кеш автоматически устанавливается в политику при экспорте через функцию `export()` в `utils.go`
- Упрощена архитектура: кеш передается в политику один раз при создании движка, а не при каждом вызове `Evaluate`
- Переименован тип `ConditionFuncsConfig` в `ConditionsMap` для более точного отражения назначения

**Обновлено:**

- Все тесты обновлены для новой сигнатуры `Policy.Evaluate` (удален параметр `cash`)
- Тесты `policy_test.go` обновлены для установки поля `Cash` в структуре `Policy` перед вызовом `Evaluate`

**Исправлено:**

- Исправлены вызовы `Policy.Evaluate` в тестах для соответствия новой сигнатуре

---

### English

**Changed:**

- Changed `Policy.Evaluate` method signature: removed `cash` parameter, cache is now stored in `Cash` field of `Policy` struct
- Cache is automatically set in policy during export via `export()` function in `utils.go`
- Simplified architecture: cache is passed to policy once during engine creation, not on every `Evaluate` call
- Renamed `ConditionFuncsConfig` type to `ConditionsMap` for better clarity of purpose

**Updated:**

- All tests updated for new `Policy.Evaluate` signature (removed `cash` parameter)
- `policy_test.go` tests updated to set `Cash` field in `Policy` struct before calling `Evaluate`

**Fixed:**

- Fixed `Policy.Evaluate` calls in tests to match new signature


## v 0.4.2

### Русский

**Добавлено:**

- Параметр `funcConfig` в `NewGuardianFromPolices` и `NewGuardianFromFile` для кастомизации функций условий
- Новые тесты для проверки кастомных функций условий (`tests/funcconfig_test.go`)
- Тесты для `ConditionFuncsConfig.Select` метода (`tests/config_test.go`)
- Тесты для изоляции сессий и работы кеша (`tests/session_test.go`)

**Обновлено:**

- Все godoc-комментарии в основном коде обновлены и актуализированы
- Все godoc-комментарии в тестах обновлены и актуализированы
- Все тесты обновлены для поддержки нового параметра `funcConfig`
- Комментарии в тестах обновлены для отражения использования `DefaultConditionsFuncs`

**Исправлено:**

- Все вызовы `NewGuardianFromPolices` и `NewGuardianFromFile` в тестах обновлены с добавлением параметра `funcConfig`
- Тесты `condition_functions_test.go` обновлены для использования `DefaultConditionsFuncs` напрямую
- Тесты `policy_test.go` обновлены для установки `ConditionsMap` перед вызовом `Evaluate`

**Документация:**

- Обновлены все godoc-комментарии для отражения текущего API
- Добавлена информация о параметре `funcConfig` во все соответствующие комментарии
- Обновлены описания функций условий и их использования

---

### English

**Added:**

- `funcConfig` parameter in `NewGuardianFromPolices` and `NewGuardianFromFile` for customizing condition functions
- New tests for custom condition functions (`tests/funcconfig_test.go`)
- Tests for `ConditionFuncsConfig.Select` method (`tests/config_test.go`)
- Tests for session isolation and cache behavior (`tests/session_test.go`)

**Updated:**

- All godoc comments in main code updated and modernized
- All godoc comments in tests updated and modernized
- All tests updated to support new `funcConfig` parameter
- Test comments updated to reflect use of `DefaultConditionsFuncs`

**Fixed:**

- All calls to `NewGuardianFromPolices` and `NewGuardianFromFile` in tests updated with `funcConfig` parameter
- `condition_functions_test.go` tests updated to use `DefaultConditionsFuncs` directly
- `policy_test.go` tests updated to set `ConditionsMap` before calling `Evaluate`

**Documentation:**

- All godoc comments updated to reflect current API
- Added information about `funcConfig` parameter to all relevant comments
- Updated descriptions of condition functions and their usage


## v 0.4.1

### Русский

**Изменено:**

- noctis-guard переименован в pbac-guardian

---

### English

**Changed:**

- noctis-guard renamed to pbac-guardian


## v 0.4.0

### Русский

**Добавлено:**

- In-memory кеширование для уже распаршенных полей в рамках сессий проверки правил
- Встроенный кеш-движок `DefaultCasher` с оптимизированной реализацией на основе `sync.RWMutex`
- Интерфейс `base.Casher` для кастомных реализаций кеша
- Автоматическая генерация уникальных `sessionID` для изоляции кеша между сессиями
- Параметр `casher` в `NewGuardianFromPolices` и `NewGuardianFromFile` (опциональный, можно передать `nil`)

**Оптимизации:**

- Оптимизированная реализация кеша с использованием `sync.RWMutex` вместо `sync.Map` для лучшей производительности
- Уменьшение аллокаций на 47-63% в production-сценариях
- Экономия времени на 14-41% в зависимости от количества обращений к полям
- Экономия памяти до 36% при 10+ обращениях к полю

**Производительность:**

Результаты бенчмарков для production-сценариев (10+ действий, 20-30 политик):
- **20 политик, 3 обращения**: 14% быстрее, 47% меньше аллокаций
- **30 политик, 5 обращений**: 32% быстрее, 25% меньше памяти, 59% меньше аллокаций
- **30 политик, 10 обращений**: 41% быстрее, 36% меньше памяти, 63% меньше аллокаций
- **Смешанные поля**: 24% быстрее, 23% меньше памяти, 53% меньше аллокаций

**Точка безубыточности:**

Кеш становится выгодным при 3+ обращениях к одному полю в рамках одной оценки действия (action). При меньшем количестве обращений рекомендуется передавать `nil` для параметра `casher`.

**Тестирование:**

- Добавлены unit-тесты для `DefaultCasher` (потокобезопасность, изоляция сессий)
- Добавлены бенчмарки для production-сценариев
- Проведен анализ производительности и точки безубыточности

---

### English

**Added:**

- In-memory caching for already parsed fields within policy evaluation sessions
- Built-in cache engine `DefaultCasher` with optimized implementation based on `sync.RWMutex`
- `base.Casher` interface for custom cache implementations
- Automatic generation of unique `sessionID` for cache isolation between sessions
- `casher` parameter in `NewGuardianFromPolices` and `NewGuardianFromFile` (optional, can pass `nil`)

**Optimizations:**

- Optimized cache implementation using `sync.RWMutex` instead of `sync.Map` for better performance
- Reduction of allocations by 47-63% in production scenarios
- Time savings of 14-41% depending on field access count
- Memory savings up to 36% with 10+ field accesses

**Performance:**

Benchmark results for production scenarios (10+ actions, 20-30 policies):
- **20 policies, 3 accesses**: 14% faster, 47% fewer allocations
- **30 policies, 5 accesses**: 32% faster, 25% less memory, 59% fewer allocations
- **30 policies, 10 accesses**: 41% faster, 36% less memory, 63% fewer allocations
- **Mixed fields**: 24% faster, 23% less memory, 53% fewer allocations

**Break-even Point:**

Cache becomes beneficial at 3+ accesses to the same field within a single action evaluation. For fewer accesses, it's recommended to pass `nil` for the `casher` parameter.

**Testing:**

- Added unit tests for `DefaultCasher` (thread safety, session isolation)
- Added benchmarks for production scenarios
- Conducted performance analysis and break-even point analysis