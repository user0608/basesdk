# Skill: `types-arrays`

Focused guide for `StrArray` and `UUIDArray` from the `types` package.

Use this skill when accepting JSON fields that may be either a single value, an array, or `null`.

## `StrArray`

`StrArray` is `[]string` with flexible JSON decoding.

Accepted JSON inputs:

```json
"admin"
["admin", "manager"]
null
```

Behavior:

- A JSON string becomes a one-element array.
- A JSON string array is used as-is.
- JSON `null` becomes an empty slice.
- Invalid JSON shapes return `ErrInvalidStrArrayInput`.
- A nil `StrArray` marshals as JSON `null`; a non-nil empty slice marshals as `[]`.

Useful methods:

- `Trimmed()` returns a copy with whitespace trimmed from each value.
- `NonEmpty()` returns a copy without empty strings.
- `Unique()` returns a copy preserving the first occurrence of each value.

Recommended cleanup:

```go
tags = tags.Trimmed().NonEmpty().Unique()
```

## `UUIDArray`

`UUIDArray` is `[]uuid.UUID` with flexible JSON decoding from strings.

Accepted JSON inputs:

```json
"550e8400-e29b-41d4-a716-446655440000"
["550e8400-e29b-41d4-a716-446655440000"]
null
```

Behavior:

- A JSON UUID string becomes a one-element array.
- A JSON array of UUID strings is parsed into UUID values.
- JSON `null` becomes an empty slice.
- Invalid JSON shapes return `ErrInvalidUUIDArrayInput`.
- Invalid UUID strings return the UUID parser error directly.
- A nil `UUIDArray` marshals as JSON `null`; a non-nil empty slice marshals as `[]`.

Useful methods:

- `Unique()` returns a copy preserving the first occurrence of each UUID.

## Common mistakes

- Do not expect `StrArray` to split comma-separated strings; `"a,b"` becomes one element.
- Do not expect `StrArray.NonEmpty()` to trim spaces; call `Trimmed()` first.
- Do not assume JSON `null` remains nil after unmarshal; it becomes an empty slice.
- Do not pass UUID objects in JSON; `UUIDArray` expects UUID strings.
- Do not ignore the unmarshal error for `UUIDArray`; invalid UUID text is reported immediately.
