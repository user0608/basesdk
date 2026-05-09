# Skill: `binds`

Package for binding Echo HTTP request input into Go payloads.

Use this skill when implementing handlers that need to read JSON bodies, query parameters, flexible array inputs, JSON embedded in multipart form fields, or uploaded form files.

## What it solves

- JSON body binding through Echo's `DefaultBinder`.
- Query parameter binding through Echo's `DefaultBinder`.
- Request fields that may be a single value or an array.
- Multipart forms that include a JSON field plus optional or required files.
- Consistent request errors using the `errs` package.

## Basic JSON Binding

Use `binds.JSON(c, &payload)` when the endpoint expects a JSON body.

```go
type CreateProductRequest struct {
	Name string `json:"name"`
	SKU  string `json:"sku"`
}

func CreateProduct(usecase *usecases.ProductUsecase) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CreateProductRequest
		if err := binds.JSON(c, &req); err != nil {
			return answer.Err(c, err)
		}

		product, err := usecase.Create(c.Request().Context(), req)
		if err != nil {
			return answer.Err(c, err)
		}

		return answer.Ok(c, product)
	}
}
```

`binds.From(c, &payload)` is currently an alias for `binds.JSON(c, &payload)`.

## Query Binding

Use `binds.Query(c, &payload)` when the endpoint reads values from query parameters.

```go
type ListProductsQuery struct {
	Search string `query:"search"`
	Page   int    `query:"page"`
}

var query ListProductsQuery
if err := binds.Query(c, &query); err != nil {
	return answer.Err(c, err)
}
```

## Flexible Array Requests

Use `binds.RequestUUIDs(c)` for endpoints whose body should contain one UUID or a list of UUIDs.

Accepted field names, in priority order:

- `code`
- `codes`
- `id`
- `ids`
- `uuid`
- `uuids`
- `values`

Accepted JSON shapes:

```json
{ "id": "550e8400-e29b-41d4-a716-446655440000" }
{ "ids": ["550e8400-e29b-41d4-a716-446655440000"] }
```

The result is deduplicated while preserving first occurrence order.

```go
ids, err := binds.RequestUUIDs(c)
if err != nil {
	return answer.Err(c, err)
}
```

Use `binds.RequestStrings(c)` for endpoints whose body should contain one string or a list of strings.

Accepted field names, in priority order:

- `code`
- `codes`
- `codigo`
- `codigos`
- `value`
- `values`
- `valor`
- `valores`
- `string`
- `strings`
- `text`
- `texts`
- `key`
- `keys`
- `id`
- `ids`

Accepted JSON shapes:

```json
{ "value": "admin" }
{ "values": ["admin", "manager"] }
```

The result is deduplicated and empty strings are removed.

```go
values, err := binds.RequestStrings(c)
if err != nil {
	return answer.Err(c, err)
}
```

## Multipart JSON Field

Use `binds.FormFieldJSON(c, fieldName, &payload)` when a multipart form contains a field whose value is a JSON document.

```go
type Metadata struct {
	Name string `json:"name"`
}

var metadata Metadata
if err := binds.FormFieldJSON(c, "metadata", &metadata); err != nil {
	return answer.Err(c, err)
}
```

The form field is required. Missing, empty, blank, or invalid JSON values return a bad request error.

## Multipart Files

Use `binds.FormFileBytesRequired(c, fieldName)` when the upload must include the file.

```go
content, err := binds.FormFileBytesRequired(c, "file")
if err != nil {
	return answer.Err(c, err)
}
```

Use `binds.FormFileBytesOptional(c, fieldName)` when the file may be omitted.

```go
content, err := binds.FormFileBytesOptional(c, "avatar")
if err != nil {
	return answer.Err(c, err)
}

if content != nil {
	// Process uploaded content.
}
```

`FormFileBytesOptional` returns `nil, nil` when the file is not present.

## Handler Pattern

Handlers should bind input, delegate work to the use case, and return errors through `answer.Err`.

```go
func UploadProductImage(usecase *usecases.ProductUsecase) echo.HandlerFunc {
	return func(c echo.Context) error {
		var metadata UploadImageRequest
		if err := binds.FormFieldJSON(c, "metadata", &metadata); err != nil {
			return answer.Err(c, err)
		}

		image, err := binds.FormFileBytesRequired(c, "image")
		if err != nil {
			return answer.Err(c, err)
		}

		if err := usecase.UploadImage(c.Request().Context(), metadata, image); err != nil {
			return answer.Err(c, err)
		}

		return answer.Success(c)
	}
}
```

## Recommendations

- Pass pointers to payload structs when binding.
- Use `JSON` for JSON body DTOs and `Query` for query DTOs.
- Use `RequestUUIDs` or `RequestStrings` only for simple bulk-action request bodies.
- Use `FormFieldJSON` for multipart metadata instead of manually reading `FormValue`.
- Use `FormFileBytesRequired` when missing files should be a client error.
- Use `FormFileBytesOptional` when missing files are valid and should produce `nil` content.
- Return binding errors with `answer.Err(c, err)`.
- Add validation after binding when DTO fields have business or format requirements.

## Common Mistakes

- Do not pass a non-pointer payload to bind functions.
- Do not expect `RequestStrings` to split comma-separated strings; use a JSON array for multiple values.
- Do not use `FormFileBytesOptional` when a missing file should fail the request.
- Do not manually expose raw binding errors to clients; let `errs` and `answer` standardize them.
