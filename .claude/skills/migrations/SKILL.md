# Skill: `migrations`

Guide for creating and managing SQL database migrations in this project.

Use this skill when adding tables, changing schemas, inserting seed data, or working with the migration runner under `setup/migrations`.

## Migration System

Migrations are SQL files executed with Goose.

The SDK embeds its migrations through:

```go
//go:embed all:migrations
var MigrationsFS embed.FS
```

Application-level migrations can also be provided with `setup.WithMigrations(...)`.

The runner combines all provided `fs.FS` sources into a virtual `migrations/` directory before passing them to Goose.

## Current Base Files

The SDK keeps core migrations grouped by scope:

- `migrations/20260503162721_system.sql`: global system tables and seed data.
- `migrations/20260507234730_tenant.sql`: tenant tables, tenant properties, app users, roles, permissions, groups, and tenant seed data.

Keep this grouping unless a change clearly belongs in a new later migration.

## File Naming

Use timestamp-prefixed file names so Goose can order migrations lexicographically.

Recommended format:

```text
YYYYMMDDHHMMSS_short_description.sql
```

Examples:

```text
20260508103000_add_invoice_tables.sql
20260508104500_add_external_sync_properties.sql
```

Migration file names must be unique across all migration sources. The runner returns an error for duplicate file names after merging sources.

## Goose Format

Every executable migration must contain `Up` and `Down` blocks.

```sql
-- +goose Up
-- +goose StatementBegin

create table example
(
    id varchar(100) not null primary key,
    name varchar(255) not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists example;

-- +goose StatementEnd
```

Files without a valid `-- +goose Up` to `-- +goose Down` block are ignored by the `script` action.

## Creating Tables

Use explicit constraints and stable names.

```sql
create table invoice
(
    tenant_codigo varchar(100) not null,
    code varchar(100) not null,
    total numeric(12, 2) not null,
    created_by varchar(100) not null,
    created_at timestamp without time zone not null,

    constraint pk_invoice
    primary key (tenant_codigo, code),

    constraint fk_invoice_tenant
    foreign key (tenant_codigo)
    references tenant (codigo)
);
```

Prefer named constraints for primary keys, foreign keys, and unique constraints on non-trivial tables.

## Tenant Tables

Tenant-owned tables should include `tenant_codigo` and reference `tenant(codigo)`.

For tenant-scoped entities, prefer composite primary keys that start with `tenant_codigo`.

```sql
constraint pk_some_table
primary key (tenant_codigo, code),

constraint fk_some_table_tenant
foreign key (tenant_codigo)
references tenant (codigo)
```

This keeps tenant data isolated and prevents accidental cross-tenant collisions.

## Property Tables

Global runtime parameters belong in `system_properties`.

Tenant-specific runtime parameters belong in `tenant_system_properties`.

`tenant_system_properties` uses a composite primary key:

```sql
primary key (tenant_codigo, key)
```

and must reference the tenant table:

```sql
constraint fk_tenant_system_properties_tenant
foreign key (tenant_codigo) references tenant (codigo)
```

Use properties for configurable values such as API keys, external system URLs, file paths, feature toggles, limits, and tenant-specific behavior.

## Down Blocks

`Down` must reverse `Up` safely.

Drop dependent tables before parent tables.

Example order for tenant tables:

```sql
drop table if exists group_role;
drop table if exists user_group;
drop table if exists app_group;
drop table if exists role_permission;
drop table if exists user_role;
drop table if exists permission;
drop table if exists role;
drop table if exists app_user;
drop table if exists tenant_system_properties;
drop table if exists tenant;
```

Do not drop a parent table before tables that reference it.

## Migration Commands

The service parses migration commands before starting the web server.

Supported actions:

- `migrate up`: applies migrations.
- `migrate down`: rolls back one migration using Goose.
- `migrate status`: prints Goose migration status.
- `migrate script`: prints a combined SQL script from all `Up` blocks.

The command parser tolerates global flags before, between, or after the command and action.

Examples:

```bash
go run ./example migrate up
go run ./example --config app.yml migrate status
go run ./example migrate script
```

## SQL Script Action

`migrate script` is project-specific behavior implemented by `MigrationRunner.SQLScript()`.

It:

- Combines all migration sources.
- Sorts migration files by name.
- Extracts only the content between `-- +goose Up` and `-- +goose Down`.
- Removes blank lines and SQL comment lines.
- Prefixes each section with `-- <filename>`.

Use this action when a deployment needs a plain SQL script instead of applying migrations directly.

## Multiple Migration Sources

The runner supports multiple migration sources through Fx groups.

SDK migrations are provided by default:

```go
migrations.ProvideFSSources(basesdk.MigrationsFS)
```

Extra application migrations are provided with service options:

```go
service := setup.NewService(
	setup.WithMigrations(app.MigrationsFS),
)
```

The runner reads only first-level `.sql` files from each root directory in each source. Nested SQL files are ignored.

For an embedded app migration FS, keep SQL files under a top-level directory such as `migrations/`.

```go
//go:embed all:migrations
var MigrationsFS embed.FS
```

## Editing Existing Migrations

Only edit existing migrations when they have not been released or applied in any shared environment.

Once a migration has been applied outside local development, create a new migration instead of modifying history.

For this SDK during early development, core migrations may be consolidated when explicitly requested.

## Recommendations

- Keep migration names unique across SDK and application sources.
- Use timestamp prefixes to control execution order.
- Keep `Up` and `Down` blocks symmetric.
- Define foreign keys explicitly.
- For tenant tables, include `tenant_codigo` and reference `tenant(codigo)`.
- Drop tables in reverse dependency order in `Down`.
- Put seed data close to the table it initializes.
- Use `migrate status` to inspect applied state.
- Use `migrate script` when producing deployable SQL for manual execution.

## Common Mistakes

- Do not create duplicate migration file names across different embedded sources.
- Do not put migration SQL in nested directories if it should be executed.
- Do not omit Goose markers.
- Do not create tenant tables without `tenant_codigo`.
- Do not forget foreign keys to `tenant(codigo)` for tenant-owned tables.
- Do not drop parent tables before child tables in `Down`.
- Do not edit already-applied migrations unless the environment is explicitly disposable.
