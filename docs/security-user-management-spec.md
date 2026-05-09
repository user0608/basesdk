# Security User Management Spec

## Goal

Implement handler end-to-end user administration for the security domain.

The API must cover the entities already defined by migrations:

- `system_account`
- `tenant`
- `app_user`
- `role`
- `permission`
- `user_role`
- `app_group`
- `user_group`
- `group_role`
- `role_permission`

There are two endpoint groups:

- System endpoints: accessible by authenticated system users.
- Tenant endpoints: accessible by authenticated tenant users.

## Non-Goals

- No endpoints to create, update, disable, enable, or delete permissions.
- No standalone repository tests.
- No standalone usecase tests.
- No custom Echo route registration outside `httpapi.NewServer` and `httpapi.AsRoute`.

## Route Security

System endpoints use:

```go
httpapi.SystemHandler
```

Tenant endpoints use:

```go
httpapi.TenantHandler
```

Tenant endpoints must resolve the tenant from the authenticated request context:

```go
tenantCodigo := auth.Tenant(ctx)
```

System endpoints that operate on tenant data receive the tenant explicitly in the route path:

```txt
:tenantCodigo
```

## Package Layout

Keep this feature under the existing `security/` package.

```txt
security/
├── models/
│   ├── system_user.go
│   ├── app_user.go
│   ├── role.go
│   ├── permission.go
│   └── group.go
├── dtos/
│   ├── users.go
│   ├── roles.go
│   ├── groups.go
│   └── permissions.go
├── repositories/
│   ├── system_account.go
│   ├── app_user.go
│   ├── role.go
│   ├── permission.go
│   └── group.go
├── usecases/
│   ├── security.go
│   ├── system_users.go
│   ├── tenant_users.go
│   ├── tenant_roles.go
│   ├── tenant_groups.go
│   └── permissions.go
└── handlers/
    ├── system_login.go
    ├── tenant_login.go
    ├── system_users.go
    ├── system_tenant_users.go
    ├── system_tenant_roles.go
    ├── system_tenant_groups.go
    ├── system_permissions.go
    ├── tenant_users.go
    ├── tenant_roles.go
    ├── tenant_groups.go
    ├── tenant_permissions.go
    └── tenant_me.go
```

Small files may be combined when readability is better. Split by entity/action only when files grow.

## Binding Rules

Use `binds.JSON(c, &payload)` for create, update, and password payloads.

Every operation that can naturally be executed in bulk must use `binds.RequestStrings(c)`.

Use bulk string bodies for:

- delete users
- enable users
- disable users
- delete roles
- enable roles
- disable roles
- delete groups
- enable groups
- disable groups
- assign role permissions
- assign group users
- assign group roles

Handlers for these operations must not read comma-separated path params or custom array DTOs.

The handler should bind once and pass the resulting slice to the usecase:

```go
values, err := binds.RequestStrings(c)
if err != nil {
	return answer.Err(c, err)
}
```

Preferred request shapes:

```json
{ "codes": ["ADMIN", "SELLER"] }
```

```json
{ "values": ["kevin", "ana"] }
```

`binds.RequestStrings` currently accepts `code`, `codes`, `codigo`, `codigos`, `value`, `values`, `valor`, `valores`, `string`, `strings`, `text`, `texts`, `key`, `keys`, `id`, and `ids`.

If user-specific bodies should use `username` or `usernames`, extend `binds.RequestStrings` before implementing those handlers.

Bulk endpoints should be idempotent when practical.

Examples:

- Disabling an already disabled user should not fail only because it was already disabled.
- Enabling an already enabled role should not fail only because it was already enabled.
- Replacing role permissions with the same list should succeed.
- Deleting an entity that is referenced by constraints should return the standardized database/business error instead of partially succeeding silently.

## System Account Endpoints

Manage `system_account`.

```txt
GET    /api/v1/system/users
POST   /api/v1/system/users
GET    /api/v1/system/users/:username
PUT    /api/v1/system/users/:username
PATCH  /api/v1/system/users/:username/password
PATCH  /api/v1/system/users/enable
PATCH  /api/v1/system/users/disable
DELETE /api/v1/system/users
```

No permissions endpoint is defined for `system_account` because migrations do not define system roles or system permissions.

## System Tenant User Endpoints

System users can manage `app_user` for any tenant.

```txt
GET    /api/v1/system/tenants/:tenantCodigo/users
POST   /api/v1/system/tenants/:tenantCodigo/users
GET    /api/v1/system/tenants/:tenantCodigo/users/:username
PUT    /api/v1/system/tenants/:tenantCodigo/users/:username
PATCH  /api/v1/system/tenants/:tenantCodigo/users/:username/password
PATCH  /api/v1/system/tenants/:tenantCodigo/users/enable
PATCH  /api/v1/system/tenants/:tenantCodigo/users/disable
DELETE /api/v1/system/tenants/:tenantCodigo/users
GET    /api/v1/system/tenants/:tenantCodigo/users/:username/permissions
```

## System Tenant Role Endpoints

System users can manage tenant roles for any tenant.

```txt
GET    /api/v1/system/tenants/:tenantCodigo/roles
POST   /api/v1/system/tenants/:tenantCodigo/roles
GET    /api/v1/system/tenants/:tenantCodigo/roles/:code
PUT    /api/v1/system/tenants/:tenantCodigo/roles/:code
PATCH  /api/v1/system/tenants/:tenantCodigo/roles/enable
PATCH  /api/v1/system/tenants/:tenantCodigo/roles/disable
DELETE /api/v1/system/tenants/:tenantCodigo/roles
GET    /api/v1/system/tenants/:tenantCodigo/roles/:code/permissions
PUT    /api/v1/system/tenants/:tenantCodigo/roles/:code/permissions
```

`PUT /roles/:code/permissions` replaces the complete permission set assigned to the role.

## System Tenant Group Endpoints

System users can manage tenant groups for any tenant.

```txt
GET    /api/v1/system/tenants/:tenantCodigo/groups
POST   /api/v1/system/tenants/:tenantCodigo/groups
GET    /api/v1/system/tenants/:tenantCodigo/groups/:code
PUT    /api/v1/system/tenants/:tenantCodigo/groups/:code
PATCH  /api/v1/system/tenants/:tenantCodigo/groups/enable
PATCH  /api/v1/system/tenants/:tenantCodigo/groups/disable
DELETE /api/v1/system/tenants/:tenantCodigo/groups
GET    /api/v1/system/tenants/:tenantCodigo/groups/:code/users
PUT    /api/v1/system/tenants/:tenantCodigo/groups/:code/users
GET    /api/v1/system/tenants/:tenantCodigo/groups/:code/roles
PUT    /api/v1/system/tenants/:tenantCodigo/groups/:code/roles
```

`PUT /groups/:code/users` replaces the complete user set assigned to the group.

`PUT /groups/:code/roles` replaces the complete role set assigned to the group.

## System Permission Endpoints

Permissions are created by migrations and are read-only through the API.

```txt
GET    /api/v1/system/permissions
GET    /api/v1/system/permissions/:code
```

## Tenant User Endpoints

Tenant users can manage users in their authenticated tenant.

```txt
GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/:username
PUT    /api/v1/users/:username
PATCH  /api/v1/users/:username/password
PATCH  /api/v1/users/enable
PATCH  /api/v1/users/disable
DELETE /api/v1/users
GET    /api/v1/users/:username/permissions
```

The tenant is not accepted from the path or body. It comes from auth context.

## Tenant Role Endpoints

Tenant users can manage roles in their authenticated tenant.

```txt
GET    /api/v1/roles
POST   /api/v1/roles
GET    /api/v1/roles/:code
PUT    /api/v1/roles/:code
PATCH  /api/v1/roles/enable
PATCH  /api/v1/roles/disable
DELETE /api/v1/roles
GET    /api/v1/roles/:code/permissions
PUT    /api/v1/roles/:code/permissions
```

`PUT /roles/:code/permissions` replaces the complete permission set assigned to the role.

## Tenant Group Endpoints

Tenant users can manage groups in their authenticated tenant.

```txt
GET    /api/v1/groups
POST   /api/v1/groups
GET    /api/v1/groups/:code
PUT    /api/v1/groups/:code
PATCH  /api/v1/groups/enable
PATCH  /api/v1/groups/disable
DELETE /api/v1/groups
GET    /api/v1/groups/:code/users
PUT    /api/v1/groups/:code/users
GET    /api/v1/groups/:code/roles
PUT    /api/v1/groups/:code/roles
```

`PUT /groups/:code/users` replaces the complete user set assigned to the group.

`PUT /groups/:code/roles` replaces the complete role set assigned to the group.

## Tenant Permission Endpoints

Permissions are global and read-only.

```txt
GET    /api/v1/permissions
GET    /api/v1/permissions/:code
```

## Current User Endpoints

Tenant users can inspect and update their own account.

```txt
GET    /api/v1/me
PATCH  /api/v1/me/password
GET    /api/v1/me/permissions
```

`GET /api/v1/me/permissions` returns effective permissions for the authenticated user.

## Effective User Permissions

Effective permissions for a tenant user are the distinct union of:

- Direct role permissions through `user_role -> role_permission -> permission`.
- Group role permissions through `user_group -> group_role -> role_permission -> permission`.

Disabled roles and disabled groups must not contribute permissions.

Recommended repository method:

```go
func (r *PermissionRepository) FindUserPermissions(ctx context.Context, tenantCodigo string, username string) ([]models.Permission, error)
```

Recommended usecase method:

```go
func (u *TenantUsersUsecase) FindUserPermissions(ctx context.Context, tenantCodigo string, username string) ([]dtos.PermissionResponse, error)
```

Conceptual SQL:

```sql
select distinct p.code, p.description
from permission p
join role_permission rp on rp.permission_code = p.code
join role r on r.tenant_codigo = rp.tenant_codigo
    and r.code = rp.role_code
join user_role ur on ur.tenant_codigo = rp.tenant_codigo
    and ur.role_code = rp.role_code
where ur.tenant_codigo = ?
  and ur.username = ?
  and r.disabled = false

union

select distinct p.code, p.description
from permission p
join role_permission rp on rp.permission_code = p.code
join role r on r.tenant_codigo = rp.tenant_codigo
    and r.code = rp.role_code
join group_role gr on gr.tenant_codigo = rp.tenant_codigo
    and gr.role_code = rp.role_code
join app_group g on g.tenant_codigo = gr.tenant_codigo
    and g.code = gr.group_code
join user_group ug on ug.tenant_codigo = gr.tenant_codigo
    and ug.group_code = gr.group_code
where ug.tenant_codigo = ?
  and ug.username = ?
  and r.disabled = false
  and g.disabled = false
```

## DTO Shape

Initial response DTOs should be simple and not expose password hashes.

User response:

```json
{
  "tenantCodigo": "local",
  "username": "kevin",
  "email": "kevin@local",
  "fullName": "Kevin",
  "emailVerified": true,
  "mustChangePassword": false,
  "lastLoginAt": null,
  "disabled": false
}
```

Role response:

```json
{
  "tenantCodigo": "local",
  "code": "SUPER_ADMIN",
  "description": "Full access role.",
  "disabled": false
}
```

Group response:

```json
{
  "tenantCodigo": "local",
  "code": "OPERATIONS",
  "description": "Operations team",
  "disabled": false
}
```

Permission response:

```json
{
  "code": "users.read",
  "description": "Read users"
}
```

Effective user permissions response:

```json
{
  "tenantCodigo": "local",
  "username": "kevin",
  "permissions": [
    {
      "code": "users.read",
      "description": "Read users"
    }
  ]
}
```

## Implementation Phases

### Phase 1: Shared Domain And Repositories

- Add missing models for role, permission, group, and relation tables as needed.
- Add DTOs for users, roles, groups, and permissions.
- Extend repositories with list, find, create, update, enable, disable, password, and assignment operations.
- Add delete operations for users, roles, and groups using bulk string bodies.
- Add effective user permissions repository query.

### Phase 2: Tenant Endpoints

- Implement tenant users.
- Implement tenant roles.
- Implement tenant role permissions.
- Implement tenant groups.
- Implement tenant group users and group roles.
- Implement tenant permissions read-only endpoints.
- Implement current user endpoints.

### Phase 3: System Endpoints

- Implement system users.
- Implement system tenant users.
- Implement system tenant roles.
- Implement system tenant groups.
- Implement system permissions read-only endpoints.

### Phase 4: Handler E2E Tests

- Add handler e2e tests using `httpapi.NewServer`.
- Use `testdb.NewPostgresStorage(t)`.
- Use `auth.NewDefaultAuthenticatedSecurityMiddleware()` for tenant and system authenticated server tests.
- Do not add standalone repository or usecase tests.

## Open Decisions

- Decide whether to extend `binds.RequestStrings` to accept `username` and `usernames`.
- Decide whether tenant user administration should require permission checks before handlers are exposed broadly.
- Decide whether effective permissions responses need source tracing later.
