---
name: usecases
description: Use when creating or modifying use cases.
---

# Usecases

## Base rule

Usecases live in:

```txt
internal/<feature>/usecases/
```

Usecases contain business logic and orchestration.

They can use:

```txt
- concrete repositories
- other usecases
- domain/models
- domain/dtos
- errs
- kcheck
```

Do not create repository interfaces.

## Naming

Use singular type names:

```txt
UserUsecase
ProductUsecase
OrderUsecase
```

For specific business flows, create a dedicated usecase:

```txt
UserLoginUsecase
UserPasswordResetUsecase
OrderCheckoutUsecase
```

## File naming

For small features:

```txt
usecases/user_usecase.go
```

For larger features:

```txt
usecases/user_usecase.go
usecases/user_login_usecase.go
usecases/user_create.go
usecases/user_update.go
```

## Base file

The base file contains:

```txt
- main usecase struct
- constructor
- shared dependencies
- small methods if the feature is simple
```

Example:

```go
type UserUsecase struct {
	userRepository *repositories.UserRepository
}

func NewUserUsecase(userRepository *repositories.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}
```

## Tenant and security context

Usecases must receive the scope required by the handler security type.

For tenant routes:

```txt
httpapi.TenantHandler -> usecase method must receive tenant context or tenant ID
```

For system routes:

```txt
httpapi.SystemHandler -> usecase method must receive system/admin context when needed
```

For public routes:

```txt
httpapi.PublicHandler -> usecase method does not receive tenant unless the endpoint explicitly resolves one
```

Tenant-scoped operations must not infer tenant inside repositories.

Handlers extract or resolve tenant data using the project standard mechanism, then pass it to the usecase.

Good:

```go
func (u *UserUsecase) List(ctx context.Context, tenantID string) ([]models.User, error)
```

Good:

```go
func (u *UserUsecase) Create(ctx context.Context, tenantID string, input dtos.CreateUserInput) (*models.User, error)
```

Bad:

```go
func (r *UserRepository) List(ctx context.Context) ([]models.User, error)
```

when the data is tenant-scoped.

Repositories may receive tenant ID from usecases as a query parameter.

## Usecase dependencies

A usecase may depend on another usecase when it needs to reuse an existing business flow.

Example:

```go
type UserLoginUsecase struct {
	userUsecase *UserUsecase
}

func NewUserLoginUsecase(userUsecase *UserUsecase) *UserLoginUsecase {
	return &UserLoginUsecase{
		userUsecase: userUsecase,
	}
}
```

Use another usecase when:

```txt
- the dependency represents a real business operation
- duplicating the logic would create inconsistency
- the dependency is part of the same feature or a clearly related feature
```

Do not create circular dependencies between usecases.

Do not use another usecase just to access its repository.

## Specific usecases

Create a specific usecase when the flow has its own business meaning or dependencies.

Good examples:

```txt
UserLoginUsecase
UserPasswordResetUsecase
UserActivationUsecase
OrderCheckoutUsecase
```

Do not create a specific usecase only because a method exists.

For simple CRUD, prefer:

```txt
UserUsecase.Create
UserUsecase.Update
UserUsecase.Delete
UserUsecase.Find
```

## Responsibilities

Usecases can:

```txt
- apply business rules
- validate business constraints
- call repositories
- call other usecases
- coordinate multiple repositories
- coordinate multiple usecases
- enforce tenant scope
- call external services when needed
- return domain models or DTOs
- return errs errors
```

Usecases must not:

```txt
- import Echo
- import httpapi
- import answer
- import binds
- parse HTTP requests
- build HTTP responses
- know about route paths or status codes
```

## Validation

Basic input validation is handled with `kcheck`.

Use `kcheck` for simple structural validations:

```txt
- required fields
- string length
- basic value checks
- simple DTO validation
```

Business validation belongs in the usecase:

```txt
- user already exists
- invalid credentials
- product unavailable
- order cannot be cancelled
- tenant cannot access resource
```

## Errors

Error creation and mapping belongs to the `errs` skill.

Usecases should return project errors from `errs` when business rules fail.

Do not return raw HTTP errors.

Do not compare error strings.

## DTOs and models

DTOs live in:

```txt
internal/<feature>/domain/dtos/
```

Models live in:

```txt
internal/<feature>/domain/models/
```

Usecases may receive DTOs or explicit params.

Prefer explicit params for small operations:

```go
func (u *UserUsecase) FindByID(ctx context.Context, tenantID string, id string) (*models.User, error)
```

Use DTOs when the input has multiple fields or represents a command:

```go
func (u *UserUsecase) Create(ctx context.Context, tenantID string, input dtos.CreateUserInput) (*models.User, error)
```

Do not use handler request DTOs as persistence models.

## Dependencies

Allowed:

```txt
usecases -> repositories
usecases -> usecases
usecases -> domain/models
usecases -> domain/dtos
usecases -> errs
usecases -> kcheck
```

Not allowed:

```txt
usecases -> handlers
usecases -> httpapi
usecases -> echo
usecases -> answer
usecases -> binds
```

## Method shape

Use context as the first parameter.

Tenant-scoped method:

```go
func (u *UserUsecase) Create(ctx context.Context, tenantID string, input dtos.CreateUserInput) (*models.User, error)
```

Public method:

```go
func (u *UserLoginUsecase) Login(ctx context.Context, input dtos.LoginInput) (*dtos.LoginOutput, error)
```

Return values should be:

```txt
- domain model
- DTO
- primitive only when it is the natural result
- error
```

## File splitting rules

Keep one file when readable.

Split when:

```txt
- the file is too large
- a business flow has several helper methods
- a specific usecase has different dependencies
- readability improves
```

Do not split just to follow a pattern.

## Do not

- Do not create repository interfaces.
- Do not import HTTP packages.
- Do not return HTTP-specific responses.
- Do not put database queries in usecases.
- Do not let repositories infer tenant scope on their own.
- Do not use another usecase just to reach its repository.
- Do not create circular usecase dependencies.
- Do not create many files for a simple CRUD.
- Do not add tests for now.
- Do not define error mapping here.