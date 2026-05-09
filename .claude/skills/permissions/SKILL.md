# Skill: `permissions`

Guide for defining and synchronizing application permissions through `setup/permissions`.

Use this skill when adding permission catalogs, embedding permission CSV files, registering permission sources, or changing permission synchronization behavior.

## Purpose

Permissions are declared in CSV files and synchronized into the `permission` database table when the service starts.

This is for static permission definitions such as API actions, menu capabilities, or authorization scopes.

## CSV Format

Each `.csv` row must have exactly two columns:

```csv
users.read,Read users
users.update,Update users
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

- `users.read`
- `users.create`
- `users.update`
- `users.delete`
- `roles.permissions.replace`
- `groups.users.replace`

Prefer action names that match behavior rather than HTTP methods.

## When Adding Permissions

1. Add the permission row to an embedded `.csv` file.
2. Keep codes unique across SDK and application permission sources.
3. Register the embedded FS with `setup.WithPermissions` if it is application-owned.
4. Ensure authorization checks use the same exact code string.

## Common Mistakes

- Do not add headers to permission CSV files unless the header is meant to be a real permission row.
- Do not use three columns, comments as extra columns, or malformed CSV rows.
- Do not leave the permission code empty.
- Do not duplicate codes across multiple CSV files or embedded sources.
- Do not expect deleted CSV rows to delete database permissions; synchronization only inserts or updates.
- Do not put runtime or tenant-specific grants here; this catalog defines available permissions, not user assignments.
- Do not ignore startup sync errors; they indicate the permission catalog is invalid or the database upsert failed.
