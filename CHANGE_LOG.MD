## v 1.1.1

### Русский

**Добавлено:**

- Поддержка доступа к полям структур, хранящихся в map полях структур
- Возможность использования путей вида `"source:groups:admins:name"` для доступа к полям в структурах внутри map
- Функция `loadFieldFromMap` теперь поддерживает любые типы map с ключами типа string (не только `map[string]any`)
- Тесты для проверки работы с map полями в структурах: `TestMapFieldInStruct`, `TestMapFieldInStructMultipleKeys`, `TestMapFieldInStructComparison`, `TestMapFieldInStructAllConditions`

**Изменено:**

- Функция `loadFieldFromMap` переписана для использования рефлексии вместо type assertion, что позволяет работать с любыми `map[string]T`
- Обновлены godoc комментарии для `loadFieldFromMap` с описанием поддержки map полей в структурах
- Обновлены godoc комментарии для `loadField` с упоминанием поддержки map полей и примеров использования

**Обновлено:**

- Все godoc комментарии обновлены для отражения поддержки map полей в структурах
- Документация в README.md обновлена с примерами использования map полей в структурах

---

### English

**Added:**

- Support for accessing fields in structures stored in map fields of structures
- Ability to use paths like `"source:groups:admins:name"` to access fields in structures within maps
- `loadFieldFromMap` function now supports any map types with string keys (not just `map[string]any`)
- Tests for map fields in structures: `TestMapFieldInStruct`, `TestMapFieldInStructMultipleKeys`, `TestMapFieldInStructComparison`, `TestMapFieldInStructAllConditions`

**Changed:**

- `loadFieldFromMap` function rewritten to use reflection instead of type assertion, allowing work with any `map[string]T`
- Updated godoc comments for `loadFieldFromMap` with description of map field support in structures
- Updated godoc comments for `loadField` with mention of map field support and usage examples

**Updated:**

- All godoc comments updated to reflect support for map fields in structures
- README.md documentation updated with examples of using map fields in structures


## v 1.1.0

### Русский

**Добавлено:**

- Функции `Le` (less than or equal) и `Ge` (greater than or equal) для проверки условий "меньше или равно" и "больше или равно" в политиках
- Реализация `LeConditionFunc` и `GeConditionFunc` для сравнения значений
- Поддержка `Le` и `Ge` в структуре `Condition` и `ConditionsMap`
- Добавлены `LeConditionFunc` и `GeConditionFunc` в `DefaultConditionsMap`
- Поддержка `Le` и `Ge` во всех типах условий (примитивы, time.Time, Comparable)
- Тесты для функций Le и Ge: `TestLeConditionFunc`, `TestGeConditionFunc` (`tests/condition_functions_test.go`)

**Изменено:**

- `DefaultConditionsMap` теперь включает реализации `LeConditionFunc` и `GeConditionFunc`
- Обновлены godoc комментарии для `DefaultConditionsMap` с упоминанием поддержки `Le` и `Ge`
- Обновлены godoc комментарии для `LeConditionFunc` и `GeConditionFunc`
- Обновлены godoc комментарии в `condition.go` для упоминания `Le` и `Ge`
- Обновлены godoc комментарии в `config.go` для упоминания `Le` и `Ge`

**Обновлено:**

- Все godoc комментарии обновлены для отражения поддержки `Le` и `Ge`
- Тесты обновлены для включения проверки функций `Le` и `Ge`

---

### English

**Added:**

- `Le` (less than or equal) and `Ge` (greater than or equal) functions for checking "less than or equal" and "greater than or equal" conditions in policies
- Implementation of `LeConditionFunc` and `GeConditionFunc` for value comparison
- Support for `Le` and `Ge` in `Condition` and `ConditionsMap` structures
- Added `LeConditionFunc` and `GeConditionFunc` to `DefaultConditionsMap`
- Support for `Le` and `Ge` in all condition types (primitives, time.Time, Comparable)
- Tests for Le and Ge functions: `TestLeConditionFunc`, `TestGeConditionFunc` (`tests/condition_functions_test.go`)

**Changed:**

- `DefaultConditionsMap` now includes `LeConditionFunc` and `GeConditionFunc` implementations
- Updated godoc comments for `DefaultConditionsMap` with mention of `Le` and `Ge` support
- Updated godoc comments for `LeConditionFunc` and `GeConditionFunc`
- Updated godoc comments in `condition.go` to mention `Le` and `Ge`
- Updated godoc comments in `config.go` to mention `Le` and `Ge`

**Updated:**

- All godoc comments updated to reflect `Le` and `Ge` support
- Tests updated to include `Le` and `Ge` function checks


## v 0.5.0

### Русский

**Добавлено:**

- Поддержка переменных окружения через `Entity_ENV` в политиках
- Метод `Policy.loadEnv()` для получения значений переменных окружения
- Возможность использования путей вида `"env:VARIABLE_NAME"` в условиях политик
- Поддержка переменных окружения во всех типах условий (Eq, Neq, Contains, Lt, Gt)
- Функция `Gt` (greater than) для проверки условия "больше" в политиках
- Реализация `GtConditionFunc` и `gtPrimitives` для сравнения значений
- Поддержка `time.Time` значений в условиях (Eq, Neq, Lt, Gt) без необходимости реализации `Comparable`
- Функция `timesCompare` для сравнения `time.Time` значений
- Поддержка времени через `Entity_TIME` в политиках
- Метод `Policy.loadTime()` для получения значений времени (например, "time:now")
- Возможность использования путей вида `"time:now"` и `"time:now:1|day"` в условиях политик
- Поддержка модификаторов времени (day, hour, minute, second, milisecond)
- Комплексные тесты для работы с переменными окружения (`tests/env_test.go`)
- Тесты для функции Gt: `TestGtCondition`, `TestGtNumericTypes`, `TestGtConditionFunc`, `TestEnvVariable_Gt`
- Комплексные тесты для работы с `time.Time`: `TestTime_Eq`, `TestTime_Neq`, `TestTime_Lt`, `TestTime_Gt`, `TestTime_MultipleConditions`, `TestTime_Contains` (`tests/time_test.go`)
- Комплексные тесты для работы с `Entity_TIME`: `TestTimeEntity_Now`, `TestTimeEntity_Modifiers`, `TestTimeEntity_InvalidPaths`, `TestTimeEntity_WithStructureFields`, `TestTimeEntity_Contains` (`tests/time_entity_test.go`)

**Изменено:**

- Метод `Policy.load()` теперь поддерживает загрузку переменных окружения через `Entity_ENV` и времени через `Entity_TIME`
- Метод `Policy.splitPath()` обновлен для поддержки путей с `Entity_ENV` и `Entity_TIME`
- Метод `Policy.toPath()` обновлен для распознавания путей к переменным окружения и времени
- Метод `Entity.IsValid()` теперь валидирует `Entity_ENV` и `Entity_TIME` как валидные значения
- Добавлены методы `Policy.loadTime()` и `Policy.loadTimeModications()` для работы с time paths
- Метод `Policy.Evaluate()` теперь поддерживает проверку условия `Gt`
- Структура `Condition` расширена полем `Gt` для проверки условия "больше"
- Структура `ConditionsMap` расширена полем `Gt` для функции сравнения
- Метод `ConditionsMap.Select()` теперь поддерживает выбор функции `Gt`
- `DefaultConditionsMap` теперь включает реализацию `GtConditionFunc`
- Обновлены godoc комментарии для `Entity`, `Entity_ENV`, `Entity_TIME`, `loadEnv()`, `loadTime()`, `load()`, `splitPath()`, `toPath()`
- Обновлены godoc комментарии для `Condition`, `ConditionsMap`, `Config`, `GtConditionFunc`, `gtPrimitives`
- Обновлены godoc комментарии в `funcs.go` для явного указания поддержки `time.Time`:
  - Пакетный комментарий обновлен с упоминанием поддержки `time.Time`
  - `DefaultConditionsMap` обновлен с указанием поддержки `time.Time`
  - `EqConditionFunc` обновлен с описанием логики сравнения для `time.Time`
  - `LtConditionFunc` обновлен с указанием поддержки `time.Time`
  - `GtConditionFunc` обновлен с указанием поддержки `time.Time`
  - Добавлен godoc комментарий для функции `timesCompare`
- Обновлены godoc комментарии в `entity.go` для отражения поддержки `Entity_TIME`
- Обновлены godoc комментарии в `policy.go` для отражения работы с time paths
- Обновлены godoc комментарии в `condition.go` и `raw-policy.go` для упоминания time paths

**Обновлено:**

- Все godoc комментарии в `entity.go` обновлены для отражения поддержки `Entity_ENV` и `Entity_TIME`
- Все godoc комментарии в `policy.go` обновлены для отражения работы с переменными окружения и time paths
- Тест `TestEntityIsValid` обновлен для проверки валидности `Entity_ENV` и `Entity_TIME`
- Добавлены тесты для всех сценариев использования переменных окружения:
  - Сравнение env переменных с литералами и полями структур
  - Использование env переменных в условиях Contains, Eq, Neq, Lt, Gt
  - Комбинации env переменных с другими условиями
  - Обработка ошибок при отсутствии переменных окружения
  - Кеширование переменных окружения
- Обновлены тесты для включения поддержки функции Gt во всех соответствующих тестовых файлах
- Добавлены тесты для всех сценариев использования `time.Time`:
  - Сравнение `time.Time` с литералами и полями структур
  - Использование `time.Time` в условиях Eq, Neq, Lt, Gt
  - Комбинации `time.Time` с другими условиями
  - Проверка диапазонов времени (Gt и Lt вместе)
  - Использование `time.Time` в условиях Contains

---

### English

**Added:**

- Environment variable support through `Entity_ENV` in policies
- `Policy.loadEnv()` method for retrieving environment variable values
- Support for paths like `"env:VARIABLE_NAME"` in policy conditions
- Environment variable support in all condition types (Eq, Neq, Contains, Lt, Gt)
- `Gt` (greater than) function for checking "greater than" condition in policies
- Implementation of `GtConditionFunc` and `gtPrimitives` for value comparison
- Support for `time.Time` values in conditions (Eq, Neq, Lt, Gt) without requiring `Comparable` implementation
- `timesCompare` function for comparing `time.Time` values
- Time support through `Entity_TIME` in policies
- `Policy.loadTime()` method for retrieving time values (e.g., "time:now")
- Support for paths like `"time:now"` and `"time:now:1|day"` in policy conditions
- Support for time modifiers (day, hour, minute, second, milisecond)
- Comprehensive tests for environment variable functionality (`tests/env_test.go`)
- Tests for Gt function: `TestGtCondition`, `TestGtNumericTypes`, `TestGtConditionFunc`, `TestEnvVariable_Gt`
- Comprehensive tests for `time.Time` functionality: `TestTime_Eq`, `TestTime_Neq`, `TestTime_Lt`, `TestTime_Gt`, `TestTime_MultipleConditions`, `TestTime_Contains` (`tests/time_test.go`)
- Comprehensive tests for `Entity_TIME` functionality: `TestTimeEntity_Now`, `TestTimeEntity_Modifiers`, `TestTimeEntity_InvalidPaths`, `TestTimeEntity_WithStructureFields`, `TestTimeEntity_Contains` (`tests/time_entity_test.go`)

**Changed:**

- `Policy.load()` method now supports loading environment variables through `Entity_ENV` and time values through `Entity_TIME`
- `Policy.splitPath()` method updated to support paths with `Entity_ENV` and `Entity_TIME`
- `Policy.toPath()` method updated to recognize environment variable and time paths
- `Entity.IsValid()` method now validates `Entity_ENV` and `Entity_TIME` as valid values
- Added `Policy.loadTime()` and `Policy.loadTimeModications()` methods for working with time paths
- `Policy.Evaluate()` method now supports checking `Gt` condition
- `Condition` structure extended with `Gt` field for "greater than" condition checking
- `ConditionsMap` structure extended with `Gt` field for comparison function
- `ConditionsMap.Select()` method now supports selecting `Gt` function
- `DefaultConditionsMap` now includes `GtConditionFunc` implementation
- Updated godoc comments for `Entity`, `Entity_ENV`, `Entity_TIME`, `loadEnv()`, `loadTime()`, `load()`, `splitPath()`, `toPath()`
- Updated godoc comments for `Condition`, `ConditionsMap`, `Config`, `GtConditionFunc`, `gtPrimitives`
- Updated godoc comments in `funcs.go` to explicitly document `time.Time` support:
  - Package comment updated with mention of `time.Time` support
  - `DefaultConditionsMap` updated with `time.Time` support indication
  - `EqConditionFunc` updated with comparison logic description for `time.Time`
  - `LtConditionFunc` updated with `time.Time` support indication
  - `GtConditionFunc` updated with `time.Time` support indication
  - Added godoc comment for `timesCompare` function
- Updated godoc comments in `entity.go` to reflect `Entity_TIME` support
- Updated godoc comments in `policy.go` to reflect time paths functionality
- Updated godoc comments in `condition.go` and `raw-policy.go` to mention time paths

**Updated:**

- All godoc comments in `entity.go` updated to reflect `Entity_ENV` and `Entity_TIME` support
- All godoc comments in `policy.go` updated to reflect environment variable and time paths functionality
- `TestEntityIsValid` test updated to check `Entity_ENV` and `Entity_TIME` validity
- Added tests for all environment variable usage scenarios:
  - Comparing env variables with literals and structure fields
  - Using env variables in Contains, Eq, Neq, Lt, Gt conditions
  - Combining env variables with other conditions
  - Error handling for missing environment variables
  - Caching of environment variables
- Updated tests to include Gt function support in all relevant test files
- Added tests for all `time.Time` usage scenarios:
  - Comparing `time.Time` with literals and structure fields
  - Using `time.Time` in Eq, Neq, Lt, Gt conditions
  - Combining `time.Time` with other conditions
  - Checking time ranges (Gt and Lt together)
  - Using `time.Time` in Contains conditions


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