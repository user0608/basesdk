---
name: feature-structure
description: Use when creating or modifying a feature structure under internal/.
---

# Feature structure

## Base structure

Every feature goes under:

```txt
internal/<feature>/
├── domain/
│   ├── models/
│   └── dtos/
├── repositories/
├── usecases/
└── handlers/
```

## Naming

- Feature folder uses plural domain name: `users`, `products`, `orders`.
- Packages use plural names: `models`, `dtos`, `repositories`, `usecases`, `handlers`.
- Type names use singular names:
  - `User`
  - `UserRepository`
  - `UserUsecase`
  - `UserHandler`

## Simple feature

For small CRUDs, keep one file per layer:

```txt
internal/user/
├── domain/
│   ├── models/
│   │   └── user.go
│   └── dtos/
│       └── user.go
├── repositories/
│   └── user_repository.go
├── usecases/
│   └── user_usecase.go
└── handlers/
    └── user_handler.go
```

## Growing feature

If a layer grows, split by action or business case:

```txt
internal/user/
├── repositories/
│   ├── user_repository.go
│   ├── user_create.go
│   └── user_find.go
├── usecases/
│   ├── user_usecase.go
│   ├── user_login_usecase.go
│   ├── user_create.go
│   └── user_update.go
└── handlers/
    ├── user_handler.go
    ├── user_create.go
    ├── user_update.go
    └── user_login.go
```

## Dependency rules

Allowed:

```txt
handlers     -> usecases, domain/dtos
usecases     -> repositories, domain/models, domain/dtos
repositories -> domain/models
```

Not allowed:

```txt
domain       -> handlers, usecases, repositories
repositories -> handlers, usecases
usecases     -> handlers
```

## Layer rules

### domain/models

- Contains domain entities.
- Does not depend on HTTP.
- Does not depend on repositories, usecases, or handlers.

### domain/dtos

- Contains request/response/input/output DTOs.
- DTOs are not domain models.
- Do not use DTOs as persistence models.

### repositories

- Contains database access.
- Uses concrete structs.
- Do not create repository interfaces.
- Does not contain business logic.
- Does not depend on handlers or usecases.

### usecases

- Contains business logic and orchestration.
- Calls repositories directly using concrete structs.
- Does not depend on HTTP handlers.
- Can be split into specific business cases when needed.

### handlers

- Contains HTTP handlers.
- Parses requests.
- Calls usecases.
- Maps responses.
- Does not call repositories directly.
- Does not contain business logic.

## File splitting rules

For handlers, usecases, and repositories:

- Keep a single file when the feature is small.
- Split files only when it improves readability.
- The base file contains the struct, constructor, and shared dependencies.
- Action files contain related methods only.

Examples:

```txt
user_handler.go
user_create.go
user_update.go
user_login.go
```

```txt
user_usecase.go
user_login_usecase.go
user_create.go
user_update.go
```

```txt
user_repository.go
user_create.go
user_find.go
user_update.go
```

## Do not

- Do not create many files for a simple CRUD.
- Do not create repository interfaces.
- Do not add tests for now.
- Do not define error handling here.
- Error handling belongs to the `errs` skill.