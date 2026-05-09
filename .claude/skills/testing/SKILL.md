---
name: testing
description: Use when creating or modifying handler end-to-end tests using httpapi, httptest, real usecases, and the test database.
---

# Testing

## Purpose

Use this skill for handler end-to-end tests.

The standard test boundary is:

```txt
HTTP request -> httpapi.NewServer -> handler -> usecase -> repository -> test database
```

Do not create standalone repository tests or standalone usecase tests by default.

Repositories and usecases are real dependencies used to exercise handlers end-to-end.

## Base Rules

- Test handlers through `httpapi.NewServer`.
- Use `httptest` to execute HTTP requests.
- Use real usecases and repositories.
- Use `testdb.NewPostgresStorage(t)` for database-backed e2e tests.
- Use `github.com/stretchr/testify/require` for setup failures and assertions.
- Use `t.Helper()` in every test helper.
- Use `t.Cleanup` for resources owned by helpers.
- Keep setup explicit and close to the test.
- Do not mock repositories or usecases unless there is a concrete reason.
- Do not call handler functions directly when `httpapi.NewServer` can be used.
- Do not create raw database connections in tests.

## Test Database

Use `testdb.NewPostgresStorage(t)` to get a real migrated PostgreSQL database for the e2e server.

```go
storage := testdb.NewPostgresStorage(t)
```

`NewPostgresStorage(t)` performs these steps:

- Starts one shared PostgreSQL container using `postgres:18-alpine`.
- Creates a `connection.StorageManager` through `connection.NewConnection`.
- Runs SDK migrations from `basesdk.MigrationsFS`.
- Runs any additional migration file systems passed by the test.
- Registers `t.Cleanup` to drop all public tables after the test.
- Returns the shared `connection.StorageManager`.

The container is shared for the test process, but each test gets a clean schema because migrations are rerun before returning and cleanup drops public tables after the test.

## Test Database Configuration

The test database uses:

- Database: `testdb`.
- Username: `testuser`.
- Password: `testpass`.
- Log level: `error`.
- Image: `postgres:18-alpine`.

Do not bypass `testdb` with raw `gorm.Open`.

## Default Seed Data

Base migrations seed default security data.

Use these values in handler e2e tests when a real migrated database is used:

- System username: `kevin`.
- System password: `maira002`.
- Tenant code: `local`.
- Tenant timezone: `America/Lima`.
- Tenant username: `kevin`.
- Tenant password: `maira002`.

Prefer these defaults over recreating equivalent fixtures unless the test specifically needs different data.

## Additional Migrations

Pass extra embedded migration file systems when the application has migrations beyond the SDK base migrations.

```go
//go:embed all:migrations
var MigrationsFS embed.FS

func newInvoiceTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	storage := testdb.NewPostgresStorage(t, MigrationsFS)
	_ = storage

	// Build repositories, usecases, routes, and server.
}
```

Additional migrations are combined with `basesdk.MigrationsFS` through `setup/migrations.MigrationRunner`.

Migration file names must be unique across all file systems.

## E2E Server Pattern

Build a small server helper per handler group.

```go
func newLoginTestServer(t *testing.T) (*httptest.Server, *jwt.TokenService) {
	t.Helper()

	storage := testdb.NewPostgresStorage(t)
	usecase, tokenService := newSecurityUsecase(t, storage)
	securityMiddleware := auth.NewTestSecurityMiddleware()

	server := httpapi.NewServer(
		[]httpapi.Route{
			handlers.SystemUserHandler(usecase),
			handlers.TenantUserHandler(usecase),
		},
		securityMiddleware,
	)

	testServer := httptest.NewServer(server)
	t.Cleanup(testServer.Close)

	return testServer, tokenService
}
```

The helper can construct repositories and usecases, but the test still targets the handler through HTTP.

Keep dependency helpers small and explicit.

```go
func newSecurityUsecase(t *testing.T, storage connection.StorageManager) (*usecases.SecurityUsecase, *jwt.TokenService) {
	t.Helper()

	dir := t.TempDir()
	keyStore, err := jwt.NewKeyStore(&jwtConfigStub{
		privateKeyPath: filepath.Join(dir, "private.pem"),
		publicKeyPath:  filepath.Join(dir, "public.pem"),
	})
	require.NoError(t, err)

	systemProps := properties.NewSystemProperties(storage)
	tokenService := jwt.NewTokenService(keyStore, systemProps)
	usecase := usecases.NewSecurityUsecase(
		tokenService,
		repositories.NewSystemUserRepository(storage),
		repositories.NewAppUserRepository(storage),
	)

	return usecase, tokenService
}
```

## Requests

Use `httptest.NewRequest` and `httptest.NewRecorder` when the test does not need real network I/O.

```go
func executeLoginRequest(t *testing.T, server *httptest.Server, path string, body string) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	server.Config.Handler.ServeHTTP(response, request)

	return response
}
```

Use the real `server.URL` only when the test specifically needs client/network behavior.

## Servers By Security Type

`httpapi.NewServer` receives `auth.SecurityMiddleware`, which is an interface.

Use focused servers for each endpoint type:

- Public endpoints: include only public routes and pass `auth.NewTestSecurityMiddleware()` as the required server dependency.
- System endpoints: include system routes and pass `auth.NewTestSecurityMiddleware()` to simulate a valid system session.
- Tenant endpoints: include tenant routes and pass `auth.NewTestSecurityMiddleware()` to simulate `kevin/tenant_default/America/Lima`.

Example:

```go
server := httpapi.NewServer(
	[]httpapi.Route{
		handlers.SomeTenantHandler(usecase),
	},
	auth.NewTestSecurityMiddleware(),
)
```

## Real JWT Middleware

Use `auth.NewSecurityMiddleware(tokenService)` only when the e2e test is specifically validating JWT behavior, token validation, token type rejection, or Authorization header handling.

For normal handler e2e tests, prefer `auth.NewTestSecurityMiddleware()` so the test focuses on endpoint behavior and not JWT mechanics.

## Public Endpoints

Public endpoints do not apply security middleware, but `httpapi.NewServer` still requires an `auth.SecurityMiddleware` value.

Use the default authenticated middleware for this dependency:

```go
server := httpapi.NewServer(
	[]httpapi.Route{
		handlers.LoginHandler(usecase),
	},
	auth.NewTestSecurityMiddleware(),
)
```

This keeps public endpoint tests simple while preserving the same server construction path used by production.

## Assertions

Use `require` when a failed assertion makes the rest of the test invalid.

Examples:

- `require.NoError(t, err)` after setup or execution.
- `require.Equal(t, expected, actual)` for HTTP status, response fields, and claims.
- `require.NotEmpty(t, value)` for generated tokens and IDs.
- `require.Error(t, err)` for expected failures.

Do not continue after failed setup.

## Cleanup

Use `t.Cleanup` for resources created by helpers.

Examples:

- Close `httptest.Server`.
- Remove temporary files if they are not under `t.TempDir()`.
- Stop explicitly owned external resources.

Most tests should not call `testdb.TerminatePostgres` directly because the PostgreSQL container is shared for the test process.

Use `TerminatePostgres(ctx)` only when a suite or external harness owns the full process lifecycle.

## Requirements

Tests using `testdb.NewPostgresStorage(t)` require a working container runtime supported by testcontainers, such as Docker or Podman.

If the container runtime is unavailable, `NewPostgresStorage(t)` fails the test through `require.NoError`.

## Common Mistakes

- Do not create standalone repository tests by default.
- Do not create standalone usecase tests by default.
- Do not call handler functions directly when `httpapi.NewServer` can be used.
- Do not manually register Echo routes in tests when `httpapi.NewServer` can be used.
- Do not mock repositories or usecases unless there is a concrete reason.
- Do not create test database connections with raw `gorm.Open`.
- Do not skip migrations when handlers depend on real tables.
- Do not pass duplicate migration file names through additional file systems.
- Do not assume test data persists across tests; cleanup drops all public tables after each test.
- Do not use real JWT middleware unless JWT behavior is the subject of the test.

## Related Skills

Use these with testing when needed:

- `handlers`
- `httpapi`
- `connection`
- `migrations`
- `answer`
- `errs`
