# Skill: `answer`

Package for standardizing HTTP JSON responses in Echo handlers.

## What it solves

- Consistent success responses using `data` or `message`.
- Consistent error responses using `message`.
- Mapping domain errors to HTTP status codes.
- Logging unexpected internal errors.

## Usage

`answer` works with any handler that receives `echo.Context`.

```go
func ListProducts(usecase *usecases.ProductUsecase) echo.HandlerFunc {
	return func(c echo.Context) error {
		products, err := usecase.List()
		if err != nil {
			return answer.Err(c, err)
		}

		return answer.Ok(c, products)
	}
}
```

## Main API

- `Ok(c, payload)` → `200` with `{ data: payload }`.
- `Created(c)` → `201` with `{ message: "Resource created successfully" }`.
- `Message(c, msg)` → `200` with `{ message: msg }`.
- `Success(c)` → `200` with `{ message: "Operation completed successfully" }`.
- `NoContent(c)` → `204` with no body.
- `Err(c, err)` → responds with an error using the standard format.

## Examples

Response with data:

```go
return answer.Ok(c, products)
```

Creation response:

```go
return answer.Created(c)
```

Custom message response:

```go
return answer.Message(c, "Product updated successfully")
```

Simple success response:

```go
return answer.Success(c)
```

Response with no body:

```go
return answer.NoContent(c)
```

Error response:

```go
if err != nil {
	return answer.Err(c, err)
}
```

## Errors

Use `answer.Err(c, err)` when a use case returns an error.

```go
product, err := usecase.Find(id)
if err != nil {
	return answer.Err(c, err)
}

return answer.Ok(c, product)
```

Rules applied by `Err`:

- If the error is a domain error, it responds with its HTTP status code and public message.
- If the error starts with `":"`, it responds with `400`.
- If the error is unexpected, it responds with `500` using a generic message and logs the error internally.

## Recommendations

- Use `Ok` when returning data.
- Use `Created` when confirming creation without returning a payload.
- Use `Message` when you need a custom message.
- Use `Success` when you only need to confirm an operation.
- Use `NoContent` when the endpoint should not return a response body.
- Use `Err` to respond to errors.
- Use the `errs` skill to create public errors with HTTP status codes.
- Avoid exposing sensitive information through arbitrary error messages.