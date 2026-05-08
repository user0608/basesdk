# Skill: `types-time`

Focused guide for `DateOnly`, `JustTime`, and `DateTimeOnly` from the `types` package.

Use this skill when handling date-only, time-only, date-time-only, SQL scan/value, JSON, query params, or day-range logic.

## Type summary

- `DateOnly` embeds `time.Time` but keeps only `YYYY-MM-DD`.
- `JustTime` is a `time.Duration` representing time since midnight.
- `DateTimeOnly` embeds `time.Time` but keeps only `YYYY-MM-DD HH:MM:SS`.

## Constructors and parsing

```go
date, err := types.NewDateOnlyFromString("2026-05-08")
hour, err := types.NewJustTimeFromString("09:30:00")

date = types.NewDateOnly(time.Now())
hour = types.NewJustTime(time.Now())
dt := types.NewDateTimeOnly(time.Now())
```

- `DateOnly` accepts `time.DateOnly`: `2006-01-02`.
- `JustTime` accepts `HH`, `HH:MM`, `HH:MM:SS`, and optional fractional seconds in the seconds component.
- `DateTimeOnly` accepts `time.DateTime`: `2006-01-02 15:04:05`.
- Empty string, quoted empty string, and `null` produce zero values.

## JSON formats

- `DateOnly`: `"2026-05-08"` or `null` when zero.
- `JustTime`: `"09:30:00"`, `"09:30:00.123456789"`, or `null` when zero.
- `DateTimeOnly`: `"2026-05-08 09:30:00"` or `null` when zero.

## Location rules

- `DateOnly.ToTimeInLocation(loc)` returns midnight for that date in `loc`; nil `loc` means UTC.
- `DateOnly.StartOfDayInLocation(loc)` is the same as `ToTimeInLocation(loc)`.
- `DateOnly.EndOfDayInLocation(loc)` returns `23:59:59.999999999` in `loc`.
- `DateTimeOnly.ToTimeInLocation(loc)` rebuilds the same date and clock fields in `loc`; nil `loc` means UTC.
- `JustTime.ToTimeInLocation(loc)` applies the time to today in `loc`; nil `loc` means UTC.
- `JustTime.ToTimeOnDateInLocation(date, loc)` applies the time to the provided date in `loc`.

## Date ranges

Use `DateOnly.ToUTCDayRange(loc)` for inclusive UTC day bounds from a local business date.

```go
startUTC, endUTC := date.ToUTCDayRange(location)
```

Use `BuildUTCDayRange(loc, start, end)` when the day has a custom time window.

```go
start, _ := types.NewJustTimeFromString("22:00:00")
end, _ := types.NewJustTimeFromString("06:00:00")
fromUTC, toUTC := date.BuildUTCDayRange(location, start, end)
```

If `end` is not after `start`, the range crosses midnight and `end` is moved to the next day.

## SQL scan/value

- `DateOnly.Scan` accepts `nil`, `time.Time`, `*time.Time`, `DateOnly`, `*DateOnly`, `string`, and `[]byte`.
- `JustTime.Scan` accepts `nil`, `time.Time`, `*time.Time`, `JustTime`, `*JustTime`, `string`, and `[]byte`.
- `DateTimeOnly.Scan` accepts `nil`, `time.Time`, `*time.Time`, `DateTimeOnly`, `*DateTimeOnly`, `string`, and `[]byte`.
- SQL `Value()` returns formatted strings, not raw `time.Time`.

## Common mistakes

- Do not use RFC3339 for `DateOnly` or `DateTimeOnly` JSON.
- Do not store timezone assumptions inside `DateOnly` or `DateTimeOnly`; choose the location at conversion time.
- Do not use `JustTime.Sub` for cross-midnight elapsed durations without handling negative values yourself.
- Do not compare `JustTime` JSON strings if fractional seconds may appear; compare typed values.
- Do not assume zero `JustTime` means midnight in JSON output; it marshals as `null`.
