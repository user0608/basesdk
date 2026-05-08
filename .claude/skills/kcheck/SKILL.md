# Skill: `kcheck`

Lightweight tag-based validator for Go structs.

## What it solves

- Declarative DTO validation using struct tags, defaulting to `chk`.
- Nested struct validation.
- Options to validate only selected fields or skip specific fields.
- Aggregated field errors with a stable format.

## Quick usage

```go
type CreateUser struct {
	Name  string `chk:"required min=2 max=50"`
	Email string `chk:"required email"`
	Age   int    `chk:"gte=18 lte=120"`
}

if err := kcheck.Valid(CreateUser{...}); err != nil {
	// err.Error(): "Name: ...; Email: ..."
}
```

## API

- `kcheck.Valid(input, skips...)` validates the struct using the default validator, skipping the given fields.
- `kcheck.ValidSelect(input, selected...)` validates only the given fields.
- `kcheck.Struct(input)` validates the whole struct, equivalent to the default validator.

You can also create an isolated validator:

```go
v := kcheck.New()
v.Register("startsx", func(f kcheck.Field) error { ... })
err := v.Struct(dto)
```

## Tags and rules

The tag, defaulting to `chk`, accepts multiple rules separated by spaces:

- Without parameter: `required`, `email`, `uuid`, `url`, `ip`, `lower`, `upper`, etc.
- With parameter: `min=2`, `max=50`, `len=6`, `oneof=a,b,c`, `prefix=USR-`, `gte=18`.

If a rule is not registered, the error will be: `"validator [X] not registered"`.

## Nested structs

- It dives into `struct` fields, including non-nil pointers to structs.
- Exception: `time.Time` is not treated as a nested struct to traverse.

## Skip vs Select

- Skip: `kcheck.Valid(user, "Email")` skips `Email`.
- Select: `kcheck.ValidSelect(user, "Address.City")` validates only that path.
- You can pass both field names, such as `"Email"`, and paths, such as `"Address.City"`.

## Errors

The aggregated error is `kcheck.Errors`, containing `[]FieldError`.

Its string format is:

```text
"<Path>: <Message>; <Path>: <Message>"
```

`Errors.Err()` returns `nil` when there are no errors.

## Included validators

Registered by default through `RegisterDefaults()`:

- Required: `required`, `nonil`.
- Length/size: `len`, `min`, `max`.
- Numeric comparators: `gt`, `gte`, `lt`, `lte`.
- Strings: `alpha`, `alphanum`, `num`, `decimal`, `lower`, `upper`.
- Format: `email`, `uuid` v4, `url`, `ip`, `ipv4`, `ipv6`.
- Advanced strings: `oneof`, `prefix`, `suffix`, `contains`.
- Dates: `date` using `DateOnly`, `time` using `TimeOnly`, `datetime` using `DateTime`, `utc` using RFC3339 UTC or `time.Time` in UTC.

## Notes

- Invalid input returns `kcheck.ErrInvalidInput`, for example `nil`, nil pointer, or non-struct input.
- `required` supports `string`, `[]string`, and nil pointers through `IsNil`.