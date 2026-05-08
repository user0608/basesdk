# Skill: `types`

Compact guide for the `types` package in this repository.

Use this skill when selecting or using custom DTO, JSON, SQL, or validation-friendly types from `types`.

## What it provides

- Date-only values: `DateOnly`.
- Time-of-day values: `JustTime`.
- Date-time values without timezone offset in the wire format: `DateTimeOnly`.
- Flexible JSON string arrays: `StrArray`.
- Flexible JSON UUID arrays: `UUIDArray`.

## Main rule

Use these types at API, DTO, query parameter, and persistence boundaries when the value shape matters.

Do not replace them with raw `string`, `[]string`, `time.Time`, or `[]uuid.UUID` unless the boundary really accepts the raw Go representation.

## Pick the right type

- Use `DateOnly` for calendar dates such as birthdays, schedules by day, or date filters. JSON format: `"2006-01-02"`.
- Use `JustTime` for a time of day without a date, such as opening hours. JSON format: `"15:04:05"` or `"15:04:05.999999999"`.
- Use `DateTimeOnly` for local date-time values without timezone in JSON. JSON format: `"2006-01-02 15:04:05"`.
- Use `StrArray` when JSON may be either `"one"`, `["one", "two"]`, or `null`.
- Use `UUIDArray` when JSON may be either a UUID string, an array of UUID strings, or `null`.

## Null and zero behavior

- `DateOnly`, `JustTime`, and `DateTimeOnly` marshal their zero value as JSON `null`.
- Empty string, whitespace, and `null` inputs reset date/time types to their zero value.
- `StrArray` and `UUIDArray` unmarshal JSON `null` as an empty slice.
- A nil `StrArray` or `UUIDArray` marshals as JSON `null`.

## SQL behavior

- `DateOnly`, `JustTime`, and `DateTimeOnly` implement `driver.Valuer` and `sql.Scanner`.
- `DateOnly.Value()` writes `YYYY-MM-DD`.
- `JustTime.Value()` writes `HH:MM:SS` or `HH:MM:SS.fffffffff`.
- `DateTimeOnly.Value()` writes `YYYY-MM-DD HH:MM:SS`.

## Common mistakes

- Do not parse `DateOnly` with RFC3339; use `time.DateOnly` format.
- Do not parse `DateTimeOnly` with RFC3339; use `time.DateTime` format.
- Do not treat `JustTime` as elapsed business duration; it wraps as a time-of-day within 24 hours when using `Add`.
- Do not expect `StrArray.NonEmpty()` to trim values. Call `Trimmed().NonEmpty()` when both are needed.
- Do not assume JSON `null` preserves nil slices on unmarshal; array types convert it to empty slices.

## Related focused skills

- Load `types-time` for detailed date/time rules.
- Load `types-arrays` for detailed `StrArray` and `UUIDArray` rules.
