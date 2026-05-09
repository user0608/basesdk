---
name: domain
description: Use when creating or modifying domain models or DTOs.
---

# Domain

## Base rule

Domain lives in:

```txt
internal/<feature>/domain/
├── models/
└── dtos/
```

Packages use plural names:

```go
package models
package dtos
```

## Models

Models live in:

```txt
internal/<feature>/domain/models/
```

Models represent business/database entities.

Example:

```go
type User struct {
	ID       string
	Email    string
	Name     string
	Disabled bool
}
```

Models can:

```txt
- represent entities
- contain domain fields
- contain simple domain helper methods
- be used by repositories and usecases
```

Models must not:

```txt
- import Echo
- import httpapi
- import answer
- import binds
- depend on handlers
- depend on usecases
- depend on repositories
```

## DTOs

DTOs live in:

```txt
internal/<feature>/domain/dtos/
```

DTOs represent input/output data between handlers and usecases.

Examples:

```go
type CreateUserInput struct {
	Email string `json:"email" chk:"nonil"`
	Name  string `json:"name" chk:"nonil"`
}
```

```go
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
```

DTOs can:

```txt
- contain request/input structs
- contain response/output structs
- contain validation tags for kcheck
- be used by handlers and usecases
```

DTOs must not:

```txt
- contain business logic
- execute queries
- call usecases
- call repositories
```

## Naming

Models use singular names:

```txt
User
Product
Order
SystemUser
```

DTOs use explicit purpose names:

```txt
CreateUserInput
UpdateUserInput
LoginInput
UserResponse
LoginResponse
ListUsersOutput
```

Avoid vague names:

```txt
UserDTO
Data
Payload
Request
Response
```

unless the file/context makes it completely clear.

## File naming

For small features:

```txt
domain/models/user.go
domain/dtos/user.go
```

For larger features:

```txt
domain/models/user.go
domain/models/profile.go

domain/dtos/user_create.go
domain/dtos/user_update.go
domain/dtos/user_login.go
domain/dtos/user_response.go
```

## Rules

```txt
- Do not use DTOs as database models.
- Do not return models directly from handlers if the API response shape is different.
- Use models between repositories and usecases.
- Use DTOs between handlers and usecases when the input/output has shape.
- Keep domain free from HTTP concerns.
- Error definitions are handled by the errs skill, not here.
- Basic validation tags can live in DTOs and are validated with kcheck.
```

## Dependencies

Allowed:

```txt
handlers     -> domain/dtos
usecases     -> domain/models, domain/dtos
repositories -> domain/models
```

Not allowed:

```txt
domain/models -> handlers, usecases, repositories
domain/dtos   -> handlers, usecases, repositories
```

## Do not

- Do not create generic `types.go` files unless there is no better name.
- Do not create `common.go` as a dumping ground.
- Do not mix models and DTOs in the same package.
- Do not add tests for now.