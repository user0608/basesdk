---
name: repositories
description: Use when creating or modifying repositories.
---

# Repositories

## Base rule

Repositories live in:

```txt
internal/<feature>/repositories/
```

Repositories contain persistence logic only.

Repositories use concrete structs. Do not create repository interfaces.

## Naming

Use singular type names:

```txt
UserRepository
ProductRepository
OrderRepository
SystemUserRepository
```

File names:

```txt
user_repository.go
user_create.go
user_find.go
user_update.go
user_delete.go
```

For small repositories, keep everything in:

```txt
user_repository.go
```

Split only when readability improves.

## Constructor

Repositories receive `connection.StorageManager`.

```go
type UserRepository struct {
	manager connection.StorageManager
}

func NewUserRepository(manager connection.StorageManager) *UserRepository {
	return &UserRepository{
		manager: manager,
	}
}
```

Connection details are defined by the `connection` skill.

## DB access

Always get the DB connection from the storage manager:

```go
tx := r.manager.Conn(ctx)
```

Do not create DB connections inside repositories.

Do not manage connection lifecycle inside repositories.

## Query style

Prefer raw SQL execution over GORM query builder when practical.

Use GORM mainly as the managed DB instance/connection provider.

Preferred:

```go
rs := tx.Raw(`
	SELECT id, name, created_at
	FROM users
	WHERE id = ?
`, id).Scan(&user)
```

Allowed when simple and readable:

```go
rs := tx.Where("username = ?", username).First(&user)
```

Avoid complex GORM builder chains when raw SQL is clearer.

## Errors

Repository errors are handled with `errs`.

Wrap database errors with:

```go
errs.Pgf(rs.Error)
```

When an update/delete affects zero rows and that means the record does not exist, return a not found error using `errs`.

Example:

```go
if rs.RowsAffected == 0 {
	return errs.NotFoundDirect("record not found")
}
```

Do not compare error strings.

Do not return HTTP errors.

Error mapping belongs to the `errs` skill.

## Tenant scope

Repositories for tenant-scoped data must receive the tenant from the usecase.

Tenant routes:

```txt
handler -> usecase receives tenant -> repository receives tenant
```

Repositories must not infer tenant from context by themselves.

Good:

```go
func (r *UserRepository) FindUser(ctx context.Context, tenantID string, id string) (*models.User, error)
```

Good:

```go
rs := tx.Raw(`
	SELECT id, tenant_id, name
	FROM users
	WHERE tenant_id = ? AND id = ?
`, tenantID, id).Scan(&user)
```

Bad:

```go
func (r *UserRepository) FindUser(ctx context.Context, id string) (*models.User, error)
```

when the record is tenant-scoped.

System repositories do not require tenant unless the system operation explicitly targets tenant data.

## Responsibilities

Repositories can:

```txt
- execute SQL
- call the managed DB instance
- map DB rows into domain models
- persist domain models
- receive filters from usecases
- apply tenant filters passed by usecases
- return errs errors
```

Repositories must not:

```txt
- contain business logic
- call usecases
- call handlers
- import Echo
- import httpapi
- import answer
- import binds
- decide HTTP status codes
- infer tenant scope on their own
```

## Models

Repositories use models from:

```txt
internal/<feature>/domain/models/
```

Do not use handler DTOs as persistence models.

Do not return raw database maps unless explicitly required.

## Method shape

Use context as the first parameter.

System/global data:

```go
func (r *SystemUserRepository) FindSystemUser(ctx context.Context, username string) (*models.SystemUser, error)
```

Tenant-scoped data:

```go
func (r *UserRepository) FindUser(ctx context.Context, tenantID string, id string) (*models.User, error)
```

Create/update:

```go
func (r *UserRepository) CreateUser(ctx context.Context, tenantID string, user *models.User) error
```

List:

```go
func (r *UserRepository) FindUsers(ctx context.Context, tenantID string) ([]models.User, error)
```

Exists:

```go
func (r *UserRepository) ExistsUser(ctx context.Context, tenantID string, username string) (bool, error)
```

## Method naming

Use explicit method names.

Examples:

```txt
FindUser
FindUsers
CreateUser
UpdateUser
DeleteUser
ExistsUser
ChangeUserPassword
EnableUser
DisableUser
```

For system resources:

```txt
FindSystemUser
CreateSystemUser
UpdateSystemUser
DeleteSystemUser
```

## File splitting rules

Base file:

```txt
user_repository.go
```

Contains:

```txt
- repository struct
- constructor
- small methods if repository is simple
```

Split when needed:

```txt
user_find.go
user_create.go
user_update.go
user_delete.go
user_status.go
```

Do not split a small CRUD into many files.

## Dependencies

Allowed:

```txt
repositories -> connection
repositories -> errs
repositories -> domain/models
repositories -> context
```

Allowed when needed:

```txt
repositories -> domain/dtos
```

Not allowed:

```txt
repositories -> handlers
repositories -> usecases
repositories -> httpapi
repositories -> echo
repositories -> answer
repositories -> binds
```

## Do not

- Do not create repository interfaces.
- Do not create DB connections.
- Do not put business logic in repositories.
- Do not call usecases.
- Do not infer tenant from context.
- Do not return HTTP-specific errors.
- Do not use GORM builder chains for complex queries when raw SQL is clearer.
- Do not add tests for now.