# Skill: `permissions`

Guide for defining and synchronizing application permissions through `setup/permissions`.

Use this skill when adding permission catalogs, embedding permission CSV files, registering permission sources, or changing permission synchronization behavior.

## Purpose

Permissions are declared in CSV files and synchronized into the `permission` database table when the service starts.

This is for static permission definitions such as API actions, menu capabilities, or authorization scopes.

The base SDK catalog lives in:

```text
permissions/security.csv
```

Go constants are generated from that CSV into:

```text
security/permissions/permissions.go
```

## CSV Format

Each `.csv` row must have exactly two columns:

```csv
security.users.read,Consultar usuarios
security.users.update,Actualizar usuarios
```

Column meaning:

- Column 1: permission `code`, required.
- Column 2: permission `description`, optional.

Parsing rules:

- Leading spaces are trimmed by the CSV reader.
- `code` and `description` are trimmed with `strings.TrimSpace`.
- Empty files return no permissions.
- Empty descriptions are allowed.
- Empty codes are rejected.
- Rows with anything other than two columns are rejected.

## File Discovery

`PermissionSynchronizer.Read()` walks every provided `fs.FS` source recursively.

Rules:

- Only files with `.csv` extension are read.
- Directories and non-CSV files are ignored.
- Nested CSV files are supported.
- Definitions are returned sorted by `Code`.
- Duplicate permission codes across all sources are rejected.

Duplicate error messages include both conflicting descriptions, so keep descriptions useful while debugging.

## Registering Permission Sources

The SDK registers its embedded base permission catalog by default through `basesdk.PermissionsFS`.

Applications should embed permission files and pass the embedded `fs.FS` through `setup.WithPermissions`.

```go
package app

import "embed"

//go:embed all:permissions
var PermissionsFS embed.FS
```

```go
service := setup.NewService(
	setup.WithPermissions(app.PermissionsFS),
)
```

The setup service provides permission sources to Fx with:

```go
setuppermissions.ProvideFSSources(basesdk.PermissionsFS)
setuppermissions.ProvideFSSources(s.permissions...)
```

Internally this uses the Fx group tag:

```go
group:"permission-fs-sources"
```

## Startup Synchronization

`Service.Run` invokes permission synchronization before starting the web server.

The synchronizer:

- Reads all configured CSV sources.
- Fails startup if any source has invalid rows, empty codes, read errors, close errors, or duplicate codes.
- Upserts each permission into the `permission` table.

Upsert behavior:

```sql
insert into permission (code, description)
values (?, ?)
on conflict (code) do update
set description = excluded.description
```

Existing permissions with the same code have their description updated.

## Generated Constants

After editing permission CSV files, regenerate Go constants:

```bash
./scripts/generate-permissions.sh
```

The script runs `cmd/permissiongen` and writes `security/permissions/permissions.go`.

Handlers should use generated constants instead of string literals:

```go
import securitypermissions "basesdk/security/permissions"

return &httpapi.TenantHandler{
	Method:        http.MethodPost,
	Path:          "/api/v1/users",
	RequiredPerms: []string{securitypermissions.SecurityUsersCreate},
	Handler:       handler,
}
```

This keeps route authorization tied to the CSV catalog. If a permission code is removed or renamed and constants are regenerated, stale handler references fail at compile time.

## Data Model

The in-memory definition is:

```go
type PermissionDefinition struct {
	Code        string
	Description string
}
```

The synchronizer expects the database to already have a `permission` table with a unique or primary key constraint on `code`.

## Recommended Codes

Use stable, descriptive, dot-separated codes.

Examples:

- `security.users.read`
- `security.users.create`
- `security.users.update`
- `security.users.delete`
- `security.roles.permissions.replace`
- `security.groups.users.replace`

Prefer action names that match behavior rather than HTTP methods.

For SDK-owned permissions, keep the module prefix such as `security.` so codes are clear when multiple modules define their own catalogs.

## When Adding Permissions

1. Add the permission row to an embedded `.csv` file.
2. Keep codes unique across SDK and application permission sources.
3. Run `./scripts/generate-permissions.sh` when editing SDK-owned CSV files.
4. Register the embedded FS with `setup.WithPermissions` if it is application-owned.
5. Use generated constants in handlers instead of raw strings.
6. Ensure authorization checks use the same exact code string.

## Common Mistakes

- Do not add headers to permission CSV files unless the header is meant to be a real permission row.
- Do not use three columns, comments as extra columns, or malformed CSV rows.
- Do not leave the permission code empty.
- Do not duplicate codes across multiple CSV files or embedded sources.
- Do not expect deleted CSV rows to delete database permissions; synchronization only inserts or updates.
- Do not put runtime or tenant-specific grants here; this catalog defines available permissions, not user assignments.
- Do not hardcode route permission strings when a generated constant exists.
- Do not forget to regenerate constants after changing SDK permission CSV files.
- Do not ignore startup sync errors; they indicate the permission catalog is invalid or the database upsert failed.
