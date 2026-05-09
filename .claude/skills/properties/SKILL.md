# Skill: `properties`

Package for runtime system parameters and tenant-specific configuration.

Use this skill when the application needs configurable values such as API keys, external service settings, file paths, feature behavior, thresholds, limits, or tenant-specific rules.

## Purpose

`properties` is the standard configuration mechanism for values that must be stored in the database and changed without recompiling the application.

Use it for parameters such as:

- API keys or credentials used to communicate with external systems.
- URLs, hosts, routes, or endpoints for integrations.
- File paths or storage locations where the system should save or read data.
- Operational limits, timeouts, flags, or thresholds.
- Feature toggles or behavior switches.
- Tenant-specific behavior and business rules.

When system behavior must change by tenant, configure it through tenant properties instead of hardcoding conditionals.

## Types Of Properties

There are two property scopes.

- `SystemProperties`: global system parameters stored in `system_properties`.
- `TenantSystemProperties`: tenant-specific parameters stored in `tenant_system_properties`.

Use global properties when the value applies to the whole installation.

Use tenant properties when the value can differ per tenant.

## System Properties

Inject `*properties.SystemProperties` when reading or writing global parameters.

```go
type TokenService struct {
	systemProps *properties.SystemProperties
}

func NewTokenService(systemProps *properties.SystemProperties) *TokenService {
	return &TokenService{systemProps: systemProps}
}
```

Read a typed value with a default:

```go
issuer, err := s.systemProps.GetString(ctx, "jwt_issuer", "basesdk")
if err != nil {
	return err
}
```

## Tenant Properties

Inject `*properties.TenantSystemProperties` when behavior depends on `tenant_codigo`.

```go
type InvoiceUsecase struct {
	tenantProps *properties.TenantSystemProperties
}

func NewInvoiceUsecase(tenantProps *properties.TenantSystemProperties) *InvoiceUsecase {
	return &InvoiceUsecase{tenantProps: tenantProps}
}
```

Read tenant-specific values by passing the tenant code:

```go
storagePath, err := u.tenantProps.GetString(ctx, tenantCodigo, "invoice_storage_path", "/var/app/invoices")
if err != nil {
	return err
}
```

Use tenant properties for behavior that changes per customer:

```go
enabled, err := u.tenantProps.GetBool(ctx, tenantCodigo, "external_sync_enabled", false)
if err != nil {
	return err
}

if enabled {
	// Run tenant-specific integration.
}
```

## Property Data Model

Both property scopes use the same base shape.

- `key`: property name.
- `value`: stored text value.
- `data_type`: expected logical type.
- `description`: optional explanation.

Tenant properties also include:

- `tenant_codigo`: tenant owner of the property.

Allowed `data_type` values are defined by the database migration:

- `string`
- `int`
- `float`
- `bool`
- `json`

## Typed Readers

Use typed readers instead of parsing values manually.

- `GetString(ctx, key, defaultValue)`
- `GetInt(ctx, key, defaultValue)`
- `GetFloat(ctx, key, defaultValue)`
- `GetBool(ctx, key, defaultValue)`
- `GetJSON(ctx, key, dest)`
- `GetTime(ctx, key, defaultValue)`
- `GetDuration(ctx, key, defaultValue)`

Tenant readers include `tenantCodigo` before `key`:

- `GetString(ctx, tenantCodigo, key, defaultValue)`
- `GetInt(ctx, tenantCodigo, key, defaultValue)`
- `GetFloat(ctx, tenantCodigo, key, defaultValue)`
- `GetBool(ctx, tenantCodigo, key, defaultValue)`
- `GetJSON(ctx, tenantCodigo, key, dest)`
- `GetTime(ctx, tenantCodigo, key, defaultValue)`
- `GetDuration(ctx, tenantCodigo, key, defaultValue)`

## Defaults

Typed readers with `defaultValue` return the default only when the property is missing.

They return an error when the property exists but cannot be converted to the requested type.

```go
maxRetries, err := props.GetInt(ctx, "external_api_max_retries", 3)
if err != nil {
	return err
}
```

`GetJSON` does not accept a default value. It requires the property to exist and unmarshals `value` into `dest`.

## Recommended Keys

Use stable, descriptive, snake_case keys.

Examples:

- `jwt_issuer`
- `time_zone`
- `invoice_storage_path`
- `external_api_base_url`
- `external_api_key`
- `external_sync_enabled`
- `max_upload_size_mb`
- `tenant_billing_mode`

## When To Use Properties

Use `properties` when a value should be configurable at runtime or may vary across installations or tenants.

Good use cases:

- Integration credentials and endpoints.
- Tenant-specific business behavior.
- Paths where files should be stored.
- Limits and operational knobs.
- Feature flags.
- JSON payloads for advanced provider configuration.

Avoid using properties for:

- Constants that are truly part of source code behavior.
- Values that must be available before the database connection exists.
- Request-specific or user-specific state.
- Secrets that require a dedicated secret manager if the deployment provides one.

## Recommendations

- Prefer `TenantSystemProperties` whenever behavior can differ by tenant.
- Prefer typed readers over manual parsing.
- Provide safe defaults for optional settings.
- Let invalid configured values return errors instead of silently ignoring them.
- Keep property keys stable because external tools or admin screens may depend on them.
- Document every operational property with a clear `description`.
- Use `json` properties for structured provider options, not for arbitrary application state.

## Common Mistakes

- Do not hardcode tenant-specific behavior when it can be configured through tenant properties.
- Do not parse `value` manually when a typed reader already exists.
- Do not use global `SystemProperties` for values that should vary per tenant.
- Do not ignore conversion errors from typed readers.
- Do not store per-request data in properties.
