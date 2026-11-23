# Migrating from v2 to v3

v3 is a breaking release. Two themes dominate:

1. **Constructors clamp instead of wrapping.** Out-of-range values are clamped to the
   nearest valid value instead of being normalized by `time.Date` (which wrapped/carried
   them). Arithmetic (`Add`, `AddDate`, `AddMonths`, `AddYearMonth`) still wraps.
2. **Renames and removals** to make the API consistent (`Nano` suffix, `ClockNano`,
   `AtClockNano`) and drop things that don't belong (GORM methods, `*Now` helpers).

## 1. Overflow behavior: wrap → clamp

| Expression                                 | v2 result             | v3 result                                 |
|--------------------------------------------|-----------------------|-------------------------------------------|
| `NewDate(2020, 2, 31)`                     | `2020-03-02`          | `2020-02-29`                              |
| `NewDate(2020, 13, 1)`                     | `2021-01-01`          | `2020-12-01`                              |
| `NewTime(25, 70, 70, 0)`                   | `02:11:10`            | `NewTimeNano(25, 70, 70, 0)` → `23:59:59` |
| `NewDateTime(2025, 12, 31, 25, 70, 70, 0)` | `2026-01-01 02:11:10` | `2025-12-31 23:59:59`                     |

Clamp ranges: month `[1,12]`, day `[1, days in month]`, hour `[0,23]`,
minute/second `[0,59]`, nanosecond `[0,999999999]`. Year is passed through unchanged.

**Arithmetic still wraps/carries** - `Add`, `AddDate`, `AddDays`, `AddMonths`, `AddYearMonth` keep their previous
behavior.

## 2. Changes

| v2                         | v3                                                                      |
|----------------------------|-------------------------------------------------------------------------|
| `NewTime(h, m, s, ns)`     | `NewTimeNano(h, m, s, ns)`                                              |
| `NewClock(h, m, s)`        | `NewTime(h, m, s)`                                                      |
| `t.Time()`                 | `t.ClockNano()`                                                         |
| `d.AtTime(h, m, s, ns)`    | `d.AtClockNano(h, m, s, ns)`                                            |
| `dt.AtTime(h, m, s, ns)`   | `dt.AtClockNano(h, m, s, ns)`                                           |
| `DateNow()`                | `DateOf(time.Now())`                                                    |
| `DateTimeNow()`            | `DateTimeOf(time.Now())`                                                |
| `DateTimeFrom(date, time)` | `NewDateTime(...)` or `date.AtClockNano(time.ClockNano())`              |
| `ClockOf(t)`               | `TimeOf(t)` (includes nanoseconds)                                      |
| `GormDataType()`           | GORM struct tag: `gorm:"type:date"` / `"type:time"` / `"type:datetime"` |

## 5. Migration checklist

- [ ] `date.DateNow()` → `date.DateOf(time.Now())`
- [ ] `date.DateTimeNow()` → `date.DateTimeOf(time.Now())`
- [ ] `NewTime(h,m,s,ns)` → `NewTimeNano(h,m,s,ns)`
- [ ] `NewClock(h,m,s)` → `NewTime(h,m,s)`
- [ ] `.Time()` → `.ClockNano()`
- [ ] `.AtTime(h,m,s,ns)` → `.AtClockNano(h,m,s,ns)`
- [ ] Replace/remove `DateTimeFrom`, `ClockOf`
- [ ] Remove `GormDataType` usages; use `gorm:"type:..."` tags
- [ ] Audit constructors that relied on wrap/carry; use `Add`/`AddDate`/`AddMonths` for arithmetic
