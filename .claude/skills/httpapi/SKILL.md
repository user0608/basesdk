# Skill: `httpapi`

Package for declaring Echo HTTP routes and registering them through Fx.

Use this skill whenever creating handlers or adding HTTP endpoints to the project.

## Main Rule

Handlers in this project should be functions that return `httpapi.Route`.

Use the same pattern as `security/handlers/system_login.go`:

```go
func SystemUserHandler(usecase *usecases.SecurityUsecase) httpapi.Route {
	return &httpapi.PublicHandler{
		Method: http.MethodPost,
		Path:   "/api/v1/system/auth/login",
		Handler: func(c echo.Context) error {
			// Bind input, validate, call use case, return answer.
		},
	}
}
```

Do not create custom handler structs unless there is a specific framework-level need. Prefer function constructors that receive dependencies and return one of the built-in route handlers.

## Route Types

Choose the handler type based on the security policy for the route.

- `httpapi.PublicHandler`: no security middleware is applied.
- `httpapi.SystemHandler`: applies system security middleware.
- `httpapi.TenantHandler`: applies tenant security middleware.

If no marker is present, the server defaults to tenant security. In normal code, use one of the built-in handler types so the security behavior is explicit.

## Route Fields

The built-in handlers support these fields:

- `Method string`: HTTP method. Defaults to `GET` when empty.
- `Path string`: Echo route path.
- `BeforeMiddlewares []echo.MiddlewareFunc`: middleware applied before security.
- `Middlewares []echo.MiddlewareFunc`: middleware applied after security.
- `RequiredPerms []string`: all listed permissions are required.
- `AnyRequiredPerms []string`: at least one listed permission is required.
- `Handler echo.HandlerFunc`: the request handler function.

`Handler` is required. If it is nil, the route returns a `500` with `route handler is not configured`.

## Supported Methods

`httpapi.NewServer` registers only these methods:

- `http.MethodGet`
- `http.MethodPost`
- `http.MethodPut`
- `http.MethodDelete`
- `http.MethodPatch`

Unsupported methods are skipped and logged as warnings.

Prefer using `net/http` constants instead of string literals:

```go
Method: http.MethodPost,
```

## Middleware Order

For every route, middleware is applied in this order:

- `BeforeMiddlewares`
- route security middleware
- permission middleware
- `Middlewares`

Security middleware is selected from the route type:

- `PublicHandler`: no security middleware.
- `SystemHandler`: `securityMiddleware.System`.
- `TenantHandler`: `securityMiddleware.Tenant`.

Use `BeforeMiddlewares` only when something must run before authentication or tenant resolution. Use `Middlewares` for route-specific behavior that should run after security.

## Route Permissions

Protected routes can require permissions directly on the route struct.

Use `RequiredPerms` when the authenticated user must have every listed permission:

```go
return &httpapi.TenantHandler{
	Method:        http.MethodPost,
	Path:          "/api/v1/users",
	RequiredPerms: []string{"security.users.create"},
	Handler: func(c echo.Context) error {
		// Handler logic.
		return answer.Created(c)
	},
}
```

Use `AnyRequiredPerms` when the authenticated user can have any one of the listed permissions:

```go
return &httpapi.TenantHandler{
	Method:           http.MethodGet,
	Path:             "/api/v1/reports",
	AnyRequiredPerms: []string{"reports.read", "reports.admin"},
	Handler: func(c echo.Context) error {
		// Handler logic.
		return answer.Ok(c, nil)
	},
}
```

Rules:

- Permission checks are applied only to non-public routes.
- `SystemHandler` and `TenantHandler` can use `RequiredPerms` and `AnyRequiredPerms`.
- `PublicHandler` exposes the same fields but permission checks are intentionally skipped for public routes.
- If both `RequiredPerms` and `AnyRequiredPerms` are set, both checks run: all `RequiredPerms` must pass and at least one `AnyRequiredPerms` must pass.
- Permission middleware runs after security middleware, so `auth.Username(ctx)` and `auth.Tenant(ctx)` must already be available.
- Permission codes must exist in the permission catalog synchronized by the `permissions` skill.

Use stable dot-separated permission codes with a module prefix, for example `security.users.read`, `security.users.create`, or `security.roles.permissions.replace`.

For SDK routes, use generated constants from `basesdk/security/permissions` instead of string literals:

```go
import securitypermissions "basesdk/security/permissions"

RequiredPerms: []string{securitypermissions.SecurityUsersRead},
```

After changing the SDK permission CSV catalog, run `./scripts/generate-permissions.sh`.

## Registration With Fx

Routes must be registered with `httpapi.AsRoute` so they are included in the `http-api-routes` group.

```go
fx.Provide(
	httpapi.AsRoute(handlers.SystemUserHandler),
)
```

`httpapi.Module` receives all grouped routes and creates the Echo server.

## Handler Pattern

Handlers should be small orchestration functions:

- Get `ctx := c.Request().Context()`.
- Bind request input with `binds`.
- Validate DTOs with `kcheck` when needed.
- Call the use case.
- Return success through `answer`.
- Return errors through `answer.Err`.

Example:

```go
func CreateProductHandler(usecase *usecases.ProductUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPost,
		Path:          "/api/v1/products",
		RequiredPerms: []string{"products.create"},
		Handler: func(c echo.Context) error {
			ctx := c.Request().Context()

			var payload struct {
				Name string `json:"name" chk:"nonil"`
			}

			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}

			if err := kcheck.Valid(payload); err != nil {
				return answer.Err(c, errs.BadRequestDirect(err.Error()))
			}

			product, err := usecase.Create(ctx, payload.Name)
			if err != nil {
				return answer.Err(c, err)
			}

			return answer.Ok(c, product)
		},
	}
}
```

## Public Route Example

Use `PublicHandler` for login, health checks, or endpoints that must not require authentication.

```go
func LoginHandler(usecase *usecases.SecurityUsecase) httpapi.Route {
	return &httpapi.PublicHandler{
		Method: http.MethodPost,
		Path:   "/api/v1/auth/login",
		Handler: func(c echo.Context) error {
			// Public endpoint logic.
			return answer.Success(c)
		},
	}
}
```

## System Route Example

Use `SystemHandler` for endpoints that require system-level authentication.

```go
func SystemStatsHandler(usecase *usecases.SystemUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/stats",
		Handler: func(c echo.Context) error {
			stats, err := usecase.Stats(c.Request().Context())
			if err != nil {
				return answer.Err(c, err)
			}

			return answer.Ok(c, stats)
		},
	}
}
```

## Tenant Route Example

Use `TenantHandler` for tenant-scoped API endpoints.

```go
func ListProductsHandler(usecase *usecases.ProductUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/products",
		RequiredPerms: []string{"products.read"},
		Handler: func(c echo.Context) error {
			products, err := usecase.List(c.Request().Context())
			if err != nil {
				return answer.Err(c, err)
			}

			return answer.Ok(c, products)
		},
	}
}
```

## Recommendations

- Keep handlers as functions that return `httpapi.Route`.
- Return `&httpapi.PublicHandler{}`, `&httpapi.SystemHandler{}`, or `&httpapi.TenantHandler{}` directly.
- Use inline `Handler: func(c echo.Context) error { ... }` as the default project style.
- Use `http.Method*` constants for route methods.
- Always set `Path` explicitly.
- Register every route constructor with `httpapi.AsRoute` in the Fx module that owns it.
- Keep business logic in use cases, not in handlers.
- Use `binds`, `kcheck`, `errs`, and `answer` for the standard request flow.
- Add `RequiredPerms` or `AnyRequiredPerms` to protected routes that need authorization beyond authentication.
- Keep route permission codes aligned with the CSV catalog managed by the `permissions` skill.
- Use generated permission constants for SDK-owned handlers.

## Common Mistakes

- Do not forget `httpapi.AsRoute`; otherwise Fx will not provide the route to the server.
- Do not leave `Handler` nil.
- Do not use unsupported methods such as `OPTIONS` unless `NewServer` is extended first.
- Do not use `PublicHandler` for tenant or system-protected endpoints.
- Do not set permissions on `PublicHandler` and expect them to run.
- Do not use `RequiredPerms` when any one permission is enough; use `AnyRequiredPerms`.
- Do not define route permissions without adding the codes to the synchronized permission catalog.
- Do not hardcode SDK route permission strings when generated constants are available.
- Do not put route-specific middleware in global Echo setup unless it truly applies globally.
