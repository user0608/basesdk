# Skill: `testdb`

Package for creating PostgreSQL-backed test database connections.

Use this skill when writing integration tests that need a real database, migrations, repositories, or use cases backed by `connection.StorageManager`.

## Purpose

`testdb` is the standard way to create database connections in tests.

It starts a PostgreSQL container with testcontainers, creates a `connection.StorageManager`, runs migrations, and cleans all public tables after each test.

This package is directly related to the `connection` skill: tests should receive and use the same `connection.StorageManager` abstraction used by production code.

## Main API

```go
func NewPostgresStorage(t *testing.T, additionalFileSystems ...fs.FS) connection.StorageManager
```

Use it at the start of tests that need a database:

```go
func TestProductRepositoryCreate(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	repo := repositories.NewProductRepository(storage)

	ctx := context.Background()

	err := repo.Create(ctx, &models.Product{...})
	require.NoError(t, err)
}
```

## What It Does

`NewPostgresStorage(t)` performs these steps:

- Starts one shared PostgreSQL container using `postgres:18-alpine`.
- Creates a `connection.StorageManager` using `connection.NewConnection`.
- Runs SDK migrations from `basesdk.MigrationsFS`.
- Runs any additional migration file systems passed by the test.
- Registers `t.Cleanup` to drop all public tables after the test.
- Returns the shared `connection.StorageManager`.

The container is started only once per test process using `sync.Once`.

## Database Configuration

The test database uses:

- Database: `testdb`
- Username: `testuser`
- Password: `testpass`
- Log level: `error`
- Image: `postgres:18-alpine`

The internal `testdb.Config` implements `configs.DatabaseConfig`, so the connection is created through the same production path:

```go
connection.NewConnection(&testdb.Config{...})
```

Do not bypass this with raw `gorm.Open` in tests.

## Repository Tests

Repositories should be tested with `connection.StorageManager` exactly as production code uses it.

```go
func TestSystemUserRepositoryFind(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	repo := repositories.NewSystemUserRepository(storage)

	user, err := repo.FindSystemUser(context.Background(), "kevin")

	require.NoError(t, err)
	require.Equal(t, "kevin", user.Username)
}
```

Use `storage.Conn(ctx)` only when the test needs direct setup or assertions that are not exposed through a repository.

## Transaction Tests

Use `storage.WithTx(ctx, fn)` when testing transactional behavior.

```go
func TestCreateRollback(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	repo := repositories.NewProductRepository(storage)

	err := storage.WithTx(context.Background(), func(ctx context.Context) error {
		if err := repo.Create(ctx, &models.Product{...}); err != nil {
			return err
		}

		return errors.New("force rollback")
	})

	require.Error(t, err)
}
```

Pass the transaction callback context to repositories. Do not replace it with `context.Background()` inside the transaction.

## Additional Migrations

Pass extra embedded migration file systems when the test belongs to an application or module with migrations beyond the SDK base migrations.

```go
//go:embed all:migrations
var MigrationsFS embed.FS

func TestInvoiceRepository(t *testing.T) {
	storage := testdb.NewPostgresStorage(t, MigrationsFS)
	// SDK migrations and app migrations are applied.
}
```

Additional migrations are combined with `basesdk.MigrationsFS` through `setup/migrations.MigrationRunner`.

Migration file names must be unique across all file systems.

## Cleanup Behavior

`NewPostgresStorage(t)` registers cleanup with `t.Cleanup`.

After each test, it drops every table in the public schema with `CASCADE`.

This keeps the shared container fast while giving each test a clean schema after it finishes.

Because cleanup happens after the test, each call to `NewPostgresStorage(t)` reruns migrations before returning storage.

## Container Lifecycle

The PostgreSQL container is shared for the test process.

Use `TerminatePostgres(ctx)` only when a suite or external test harness needs to explicitly stop the container.

```go
func TestMain(m *testing.M) {
	code := m.Run()
	_ = testdb.TerminatePostgres(context.Background())
	os.Exit(code)
}
```

Most tests do not need to call `TerminatePostgres` directly.

## Requirements

Tests using `testdb` require a working container runtime supported by testcontainers, such as Docker or Podman.

If the container runtime is unavailable, `NewPostgresStorage(t)` fails the test through `require.NoError`.

## Recommendations

- Use `testdb.NewPostgresStorage(t)` for all database integration tests.
- Inject the returned `connection.StorageManager` into repositories and use cases.
- Use the same `connection` patterns as production code: `Conn(ctx)` and `WithTx(ctx, fn)`.
- Pass additional migration file systems for application-specific tables.
- Use repositories for setup when possible; use direct SQL only for focused test setup or assertions.
- Keep test data explicit inside each test.
- Prefer `require.NoError` for database setup failures.

## Common Mistakes

- Do not create test database connections with raw `gorm.Open`.
- Do not skip migrations when testing repositories that depend on real tables.
- Do not pass duplicate migration file names through additional file systems.
- Do not use `context.Background()` inside a `WithTx` callback when the operation must participate in the transaction.
- Do not assume test data persists across tests; cleanup drops all public tables after each test.
- Do not call `TerminatePostgres` from individual tests unless the test explicitly owns the full process lifecycle.
