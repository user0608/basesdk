# Skill: `errs`

Package for building application errors with:

- An associated HTTP status code.
- A public message suitable for API responses.
- An optional internal error for diagnostics.

Use this skill whenever an error is going to be propagated to a handler and returned through `answer.Err`.

## Base type

`*errs.Err` implements `error` and exposes:

- `Message() string` public API-safe message.
- `Code() int` HTTP status code.
- `Wrapped() error` internal error, if any.

The `answer` skill recognizes `*errs.Err` and converts it into a standardized HTTP JSON response.

## Main rule

If an error is going to cross application layers and eventually reach an HTTP handler, convert it to an `errs` error.

Do not propagate raw database, validation, or internal errors directly to handlers.

Recommended flow:

```go
result, err := usecase.Execute(ctx, req)
if err != nil {
	return answer.Err(c, err)
}

return answer.Ok(c, result)
```

The handler should only call `answer.Err(c, err)`. The error should already be correctly classified before reaching the handler.

## Common constructors

Use these when wrapping an internal error:

- `BadRequestError(err, format, args...)` → `400`.
- `NotFoundError(err, format, args...)` → `404`.
- `UnauthorizedError(err, format, args...)` → `401`.
- `ForbiddenError(err, format, args...)` → `403`.
- `UnsupportedMediaTypeError(err, format, args...)` → `415`.
- `InternalError(err, format, args...)` → `500`.

Use these when there is no internal error to wrap:

- `BadRequestf(format, args...)`.
- `NotFoundf(format, args...)`.
- `InternalErrorf(format, args...)`.

Use these when the message is already complete:

- `BadRequestDirect(message)`.
- `NotFoundDirect(message)`.
- `InternalErrorDirect(message)`.
- `UnauthorizedDirect(message)`.
- `ForbiddenDirect(message)`.
- `UnsupportedMediaTypeDirect(message)`.

## Other utilities

- `NewWithMessage(err, message)` changes the public message while preserving the code if `err` is already `*Err`; otherwise, it assumes `400`.
- `WrapError(err, message, httpCode)` returns `nil` if `err == nil`.
- `ContainsMessage(err, substr)` checks `Message()` when `err` is `*Err`.
- `IsBadRequest(err)` checks whether an error is a bad request.
- `IsInternalError(err)` checks whether an error is an internal error.
- `IsErr(err)` checks whether an error is `*Err`.
- `ToSummary(err)` returns a stable summary: `"Error <code>: <message>"`.

## Gorm and Postgres errors

All errors returned by Gorm must pass through `errs` before being propagated.

Use `errs.Pgf(err)` for database errors coming from Gorm or pgx.

```go
func (r *SystemUserRepository) FindSystemUser(ctx context.Context, username string) (*models.SystemUser, error) {
	tx := r.manager.Conn(ctx)

	var user models.SystemUser
	rs := tx.Where("username = ?", username).First(&user)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return &user, nil
}
```

`Pgf(err)` converts database errors into `*errs.Err` with public API-safe messages.

It handles:

- `gorm.ErrRecordNotFound`.
- `*pgconn.PgError` from pgx.
- SQLSTATE-based Postgres errors.
- Common constraint errors such as foreign key, unique, invalid input, and related database failures.

Use `IsPgErrCode(err, code)` only when the application needs to branch based on a specific SQLSTATE.

## Service and use case errors

Use `errs` for validation, authorization, domain, and business-rule errors that can reach a handler.

```go
func (u *ProductUsecase) Find(ctx context.Context, id string) (*models.Product, error) {
	if id == "" {
		return nil, errs.BadRequestDirect("product id is required")
	}

	product, err := u.repo.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
```

When adding context to an internal failure, wrap the original error:

```go
func (u *ProductUsecase) Create(ctx context.Context, req CreateProductRequest) error {
	if req.Name == "" {
		return errs.BadRequestDirect("product name is required")
	}

	if err := u.repo.Create(ctx, req); err != nil {
		return errs.InternalError(err, "could not create product")
	}

	return nil
}
```

## Handler usage

Handlers should not decide database or domain error details. They should return the error through `answer.Err`.

```go
func ListProducts(usecase *usecases.ProductUsecase) echo.HandlerFunc {
	return func(c echo.Context) error {
		products, err := usecase.List(c.Request().Context())
		if err != nil {
			return answer.Err(c, err)
		}

		return answer.Ok(c, products)
	}
}
```

## Recommendations

- Use `errs.Pgf` for all Gorm or pgx errors before propagating them.
- Use `errs` for every error that may reach an HTTP handler.
- Do not return raw database errors from repositories.
- Do not return raw internal errors from use cases when they may reach a handler.
- Keep public messages clear, stable, and safe for API clients.
- Wrap internal errors only when the caller needs a public message but diagnostics should be preserved internally.
- Let handlers delegate error responses to the `answer` skill.