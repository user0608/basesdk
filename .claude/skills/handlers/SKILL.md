---
name: handlers
description: Use when creating or modifying HTTP handlers.
---

# Handlers

## Base rule

Handlers live in:

```txt
internal/<feature>/handlers/
```

Handlers are Echo route constructors that return `httpapi.Route`.

Do not create handler structs by default.

## Pattern

Use functions like:

```go
func CreateUserHandler(usecase *usecases.UserUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method: http.MethodPost,
		Path:   "/api/v1/users",
		Handler: func(c echo.Context) error {
			ctx := c.Request().Context()

			// Bind request with binds.
			// Call usecase.
			// Return with answer.
			_ = ctx

			return nil
		},
	}
}
```

## Security type

Use:

```txt
httpapi.PublicHandler  -> public endpoints
httpapi.SystemHandler  -> system authenticated endpoints
httpapi.TenantHandler  -> tenant authenticated endpoints
```

Default preference:

```txt
httpapi.TenantHandler
```

unless the endpoint is explicitly public or system-level.

## File naming

For small CRUDs:

```txt
handlers/user_handler.go
```

For larger handlers:

```txt
handlers/user_handler.go
handlers/user_create.go
handlers/user_update.go
handlers/user_delete.go
handlers/user_login.go
```

## Base file

`user_handler.go` may contain:

```txt
- small CRUD handlers if the file remains readable
- common handler-related helpers if needed
```

If the handler grows, split by action.

## Responsibilities

Handlers can:

```txt
- get request context
- bind request input using binds
- call usecases
- return success through answer
- return errors through answer.Err
```

Handlers must not:

```txt
- call repositories directly
- contain business logic
- manually format standard responses
- manually implement error response logic
- create custom route registration patterns
- create handler structs unless there is a specific framework-level need
```

## DTOs

Request and response DTOs live in:

```txt
internal/<feature>/domain/dtos/
```

Do not define large DTOs inside handlers.

Small anonymous payloads are allowed only for very small endpoints.

## Routes

Routes are registered using `httpapi.AsRoute`.

Example:

```go
fx.Provide(
	httpapi.AsRoute(handlers.CreateUserHandler),
)
```

Do not register routes manually in Echo.

## Required flow

```txt
1. ctx := c.Request().Context()
2. bind input with binds
3. validate if needed
4. call usecase
5. return success through answer
6. return errors through answer.Err
```

## Related skills

Use these when needed:

```txt
feature-structure
httpapi
binds
answer
errs
```