## v 0.4.0

### Русский

**Добавлено:**

- In-memory кеширование для уже распаршенных полей в рамках сессий проверки правил
- Встроенный кеш-движок `DefaultCasher` с оптимизированной реализацией на основе `sync.RWMutex`
- Интерфейс `base.Casher` для кастомных реализаций кеша
- Автоматическая генерация уникальных `sessionID` для изоляции кеша между сессиями
- Параметр `casher` в `NewNoctisFromPolices` и `NewNoctisFromFile` (опциональный, можно передать `nil`)

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
- `casher` parameter in `NewNoctisFromPolices` and `NewNoctisFromFile` (optional, can pass `nil`)

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