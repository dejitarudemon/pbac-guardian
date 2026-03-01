## v 0.5.0

### Русский

**Добавлено:**

- Поддержка переменных окружения через `Entity_ENV` в политиках
- Метод `Policy.loadEnv()` для получения значений переменных окружения
- Возможность использования путей вида `"env:VARIABLE_NAME"` в условиях политик
- Поддержка переменных окружения во всех типах условий (Eq, Neq, Contains, Lt, Gt)
- Функция `Gt` (greater than) для проверки условия "больше" в политиках
- Реализация `GtConditionFunc` и `gtPrimitives` для сравнения значений
- Комплексные тесты для работы с переменными окружения (`tests/env_test.go`)
- Тесты для функции Gt: `TestGtCondition`, `TestGtNumericTypes`, `TestGtConditionFunc`, `TestEnvVariable_Gt`

**Изменено:**

- Метод `Policy.load()` теперь поддерживает загрузку переменных окружения через `Entity_ENV`
- Метод `Policy.splitPath()` обновлен для поддержки путей с `Entity_ENV`
- Метод `Policy.toPath()` обновлен для распознавания путей к переменным окружения
- Метод `Entity.IsValid()` теперь валидирует `Entity_ENV` как валидное значение
- Метод `Policy.Evaluate()` теперь поддерживает проверку условия `Gt`
- Структура `Condition` расширена полем `Gt` для проверки условия "больше"
- Структура `ConditionsMap` расширена полем `Gt` для функции сравнения
- Метод `ConditionsMap.Select()` теперь поддерживает выбор функции `Gt`
- `DefaultConditionsMap` теперь включает реализацию `GtConditionFunc`
- Обновлены godoc комментарии для `Entity`, `Entity_ENV`, `loadEnv()`, `load()`, `splitPath()`, `toPath()`
- Обновлены godoc комментарии для `Condition`, `ConditionsMap`, `Config`, `GtConditionFunc`, `gtPrimitives`

**Обновлено:**

- Все godoc комментарии в `entity.go` обновлены для отражения поддержки `Entity_ENV`
- Все godoc комментарии в `policy.go` обновлены для отражения работы с переменными окружения
- Тест `TestEntityIsValid` обновлен для проверки валидности `Entity_ENV`
- Добавлены тесты для всех сценариев использования переменных окружения:
  - Сравнение env переменных с литералами и полями структур
  - Использование env переменных в условиях Contains, Eq, Neq, Lt, Gt
  - Комбинации env переменных с другими условиями
  - Обработка ошибок при отсутствии переменных окружения
  - Кеширование переменных окружения
- Обновлены тесты для включения поддержки функции Gt во всех соответствующих тестовых файлах

---

### English

**Added:**

- Environment variable support through `Entity_ENV` in policies
- `Policy.loadEnv()` method for retrieving environment variable values
- Support for paths like `"env:VARIABLE_NAME"` in policy conditions
- Environment variable support in all condition types (Eq, Neq, Contains, Lt, Gt)
- `Gt` (greater than) function for checking "greater than" condition in policies
- Implementation of `GtConditionFunc` and `gtPrimitives` for value comparison
- Comprehensive tests for environment variable functionality (`tests/env_test.go`)
- Tests for Gt function: `TestGtCondition`, `TestGtNumericTypes`, `TestGtConditionFunc`, `TestEnvVariable_Gt`

**Changed:**

- `Policy.load()` method now supports loading environment variables through `Entity_ENV`
- `Policy.splitPath()` method updated to support paths with `Entity_ENV`
- `Policy.toPath()` method updated to recognize environment variable paths
- `Entity.IsValid()` method now validates `Entity_ENV` as a valid value
- `Policy.Evaluate()` method now supports checking `Gt` condition
- `Condition` structure extended with `Gt` field for "greater than" condition checking
- `ConditionsMap` structure extended with `Gt` field for comparison function
- `ConditionsMap.Select()` method now supports selecting `Gt` function
- `DefaultConditionsMap` now includes `GtConditionFunc` implementation
- Updated godoc comments for `Entity`, `Entity_ENV`, `loadEnv()`, `load()`, `splitPath()`, `toPath()`
- Updated godoc comments for `Condition`, `ConditionsMap`, `Config`, `GtConditionFunc`, `gtPrimitives`

**Updated:**

- All godoc comments in `entity.go` updated to reflect `Entity_ENV` support
- All godoc comments in `policy.go` updated to reflect environment variable functionality
- `TestEntityIsValid` test updated to check `Entity_ENV` validity
- Added tests for all environment variable usage scenarios:
  - Comparing env variables with literals and structure fields
  - Using env variables in Contains, Eq, Neq, Lt, Gt conditions
  - Combining env variables with other conditions
  - Error handling for missing environment variables
  - Caching of environment variables
- Updated tests to include Gt function support in all relevant test files


## v 0.4.3

### Русский

**Добавлено:**

- Структура `base.Config` для конфигурации движка (условия и порог отключения кеша)
- `CashTree` для отслеживания частоты доступа к полям и отключения кеша для редко используемых полей
- Параметр `config` в `NewGuardianFromPolices` и `NewGuardianFromFile` вместо `funcConfig`
- Структура `base.RawPolicy` для JSON unmarshaling и создания политик (публичные поля)
- Конструктор `base.NewPolicy` для создания `Policy` из `RawPolicy` с валидацией
- Метод `Policy.Effect()` для получения эффекта политики

**Изменено:**

- Изменена сигнатура метода `Policy.Evaluate`: удален параметр `cash`, кеш теперь хранится в поле `Cash` структуры `Policy`
- Кеш автоматически устанавливается в политику при экспорте через функцию `export()` в `utils.go`
- Упрощена архитектура: кеш передается в политику один раз при создании движка, а не при каждом вызове `Evaluate`
- Переименован тип `ConditionFuncsConfig` в `ConditionsMap` для более точного отражения назначения
- Оптимизирована проверка дубликатов в `export()`: используется `map` вместо `slices.Contains` для O(1) lookup
- Оптимизирован метод `Policy.Evaluate`: убрана рефлексия для итерации по полям условий, используются прямые проверки полей (Contains, Eq, Neq, Lt) для улучшения производительности на 40-60%
- Поля структуры `Policy` стали приватными для защиты инкапсуляции
- `NewGuardianFromPolices` и `NewGuardianFromFile` теперь принимают `[]base.RawPolicy` вместо `[]base.Policy`
- Функция `export()` теперь принимает `[]base.RawPolicy` и возвращает `map[string][]*base.Policy`

**Обновлено:**

- Все тесты обновлены для новой сигнатуры `Policy.Evaluate` (удален параметр `cash`)
- Все тесты обновлены для использования `base.Config` вместо прямого `ConditionsMap`
- Все тесты обновлены для использования `base.RawPolicy` вместо `base.Policy` в массивах политик
- Тесты `policy_test.go` обновлены для использования `NewPolicy` с `RawPolicy` вместо прямого создания `Policy`
- Все godoc комментарии обновлены и приведены в порядок
- Обновлены godoc комментарии для `Policy`, `RawPolicy`, `NewPolicy`, `Effect()`, `Policy.Evaluate`, `load()`, `IsValid()` для отражения изменений архитектуры
- Обновлены godoc комментарии в `utils.go` и `engine.go` для отражения использования `RawPolicy`

**Исправлено:**

- Исправлены вызовы `Policy.Evaluate` в тестах для соответствия новой сигнатуре

---

### English

**Added:**

- `base.Config` structure for engine configuration (condition functions and cache disable threshold)
- `CashTree` for tracking field access frequency and disabling cache for rarely accessed fields
- `config` parameter in `NewGuardianFromPolices` and `NewGuardianFromFile` instead of `funcConfig`
- `base.RawPolicy` structure for JSON unmarshaling and policy creation (public fields)
- `base.NewPolicy` constructor for creating `Policy` from `RawPolicy` with validation
- `Policy.Effect()` method for getting policy effect

**Changed:**

- Changed `Policy.Evaluate` method signature: removed `cash` parameter, cache is now stored in `Cash` field of `Policy` struct
- Cache is automatically set in policy during export via `export()` function in `utils.go`
- Simplified architecture: cache is passed to policy once during engine creation, not on every `Evaluate` call
- Renamed `ConditionFuncsConfig` type to `ConditionsMap` for better clarity of purpose
- Optimized duplicate check in `export()`: uses `map` instead of `slices.Contains` for O(1) lookup
- Optimized `Policy.Evaluate` method: removed reflection for iterating over condition fields, uses direct field checks (Contains, Eq, Neq, Lt) for 40-60% performance improvement
- `Policy` struct fields are now private for encapsulation protection
- `NewGuardianFromPolices` and `NewGuardianFromFile` now accept `[]base.RawPolicy` instead of `[]base.Policy`
- `export()` function now accepts `[]base.RawPolicy` and returns `map[string][]*base.Policy`

**Updated:**

- All tests updated for new `Policy.Evaluate` signature (removed `cash` parameter)
- All tests updated to use `base.Config` instead of direct `ConditionsMap`
- All tests updated to use `base.RawPolicy` instead of `base.Policy` in policy arrays
- `policy_test.go` tests updated to use `NewPolicy` with `RawPolicy` instead of direct `Policy` creation
- All godoc comments updated and cleaned up
- Updated godoc comments for `Policy`, `RawPolicy`, `NewPolicy`, `Effect()`, `Policy.Evaluate`, `load()`, `IsValid()` to reflect architecture changes
- Updated godoc comments in `utils.go` and `engine.go` to reflect `RawPolicy` usage

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