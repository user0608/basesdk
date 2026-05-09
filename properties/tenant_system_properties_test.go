package properties_test

import (
	"basesdk/connection"
	"basesdk/properties"
	"basesdk/testdb"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTenantSystemPropertiesCRUD(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	props := properties.NewTenantSystemProperties(storage)
	ctx := context.Background()

	description := "Tenant file storage path"
	property := &properties.TenantSystemProperty{
		TenantCodigo: "tenant_default",
		Key:          "file_storage_path",
		Value:        "/var/app/tenant-default/files",
		DataType:     "string",
		Description:  &description,
	}

	require.NoError(t, props.Create(ctx, property))

	found, err := props.Get(ctx, property.TenantCodigo, property.Key)
	require.NoError(t, err)
	require.Equal(t, property.TenantCodigo, found.TenantCodigo)
	require.Equal(t, property.Key, found.Key)
	require.Equal(t, property.Value, found.Value)
	require.NotNil(t, found.Description)
	require.Equal(t, description, *found.Description)

	exists, err := props.Exists(ctx, property.TenantCodigo, property.Key)
	require.NoError(t, err)
	require.True(t, exists)

	property.Value = "/mnt/tenant-default/files"
	require.NoError(t, props.Update(ctx, property))

	value, err := props.GetString(ctx, property.TenantCodigo, property.Key, "")
	require.NoError(t, err)
	require.Equal(t, property.Value, value)

	require.NoError(t, props.Delete(ctx, property.TenantCodigo, property.Key))

	exists, err = props.Exists(ctx, property.TenantCodigo, property.Key)
	require.NoError(t, err)
	require.False(t, exists)
}

func TestTenantSystemPropertiesAreScopedByTenant(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	props := properties.NewTenantSystemProperties(storage)
	ctx := context.Background()

	createTenant(t, storage, "tenant_two")

	require.NoError(t, props.Create(ctx, &properties.TenantSystemProperty{
		TenantCodigo: "tenant_default",
		Key:          "external_sync_enabled",
		Value:        "true",
		DataType:     "bool",
	}))
	require.NoError(t, props.Create(ctx, &properties.TenantSystemProperty{
		TenantCodigo: "tenant_two",
		Key:          "external_sync_enabled",
		Value:        "false",
		DataType:     "bool",
	}))

	defaultEnabled, err := props.GetBool(ctx, "tenant_default", "external_sync_enabled", false)
	require.NoError(t, err)
	require.True(t, defaultEnabled)

	tenantTwoEnabled, err := props.GetBool(ctx, "tenant_two", "external_sync_enabled", true)
	require.NoError(t, err)
	require.False(t, tenantTwoEnabled)

	allDefault, err := props.GetAll(ctx, "tenant_default")
	require.NoError(t, err)
	require.Len(t, allDefault, 1)
	require.Equal(t, "tenant_default", allDefault[0].TenantCodigo)
}

func TestTenantSystemPropertiesTypedReaders(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	props := properties.NewTenantSystemProperties(storage)
	ctx := context.Background()
	tenantCodigo := "tenant_default"

	entries := []*properties.TenantSystemProperty{
		{TenantCodigo: tenantCodigo, Key: "tenant_max_upload_mb", Value: "25", DataType: "int"},
		{TenantCodigo: tenantCodigo, Key: "tenant_tax_rate", Value: "18.5", DataType: "float"},
		{TenantCodigo: tenantCodigo, Key: "tenant_feature_enabled", Value: "yes", DataType: "bool"},
		{TenantCodigo: tenantCodigo, Key: "tenant_provider_options", Value: `{"region":"us","retries":3}`, DataType: "json"},
		{TenantCodigo: tenantCodigo, Key: "tenant_cutoff_time", Value: "2026-05-08T20:15:00Z", DataType: "string"},
		{TenantCodigo: tenantCodigo, Key: "tenant_cache_ttl", Value: "15m", DataType: "string"},
	}

	for _, entry := range entries {
		require.NoError(t, props.Create(ctx, entry))
	}

	intValue, err := props.GetInt(ctx, tenantCodigo, "tenant_max_upload_mb", 0)
	require.NoError(t, err)
	require.Equal(t, 25, intValue)

	floatValue, err := props.GetFloat(ctx, tenantCodigo, "tenant_tax_rate", 0)
	require.NoError(t, err)
	require.InDelta(t, 18.5, floatValue, 0.001)

	boolValue, err := props.GetBool(ctx, tenantCodigo, "tenant_feature_enabled", false)
	require.NoError(t, err)
	require.True(t, boolValue)

	var jsonValue struct {
		Region  string `json:"region"`
		Retries int    `json:"retries"`
	}
	require.NoError(t, props.GetJSON(ctx, tenantCodigo, "tenant_provider_options", &jsonValue))
	require.Equal(t, "us", jsonValue.Region)
	require.Equal(t, 3, jsonValue.Retries)

	timeValue, err := props.GetTime(ctx, tenantCodigo, "tenant_cutoff_time", time.Time{})
	require.NoError(t, err)
	require.Equal(t, time.Date(2026, 5, 8, 20, 15, 0, 0, time.UTC), timeValue)

	durationValue, err := props.GetDuration(ctx, tenantCodigo, "tenant_cache_ttl", 0)
	require.NoError(t, err)
	require.Equal(t, 15*time.Minute, durationValue)
}

func TestTenantSystemPropertiesDefaultsAndInvalidValues(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	props := properties.NewTenantSystemProperties(storage)
	ctx := context.Background()
	tenantCodigo := "tenant_default"

	value, err := props.GetString(ctx, tenantCodigo, "missing_path", "/tmp/default")
	require.NoError(t, err)
	require.Equal(t, "/tmp/default", value)

	require.NoError(t, props.Create(ctx, &properties.TenantSystemProperty{
		TenantCodigo: tenantCodigo,
		Key:          "invalid_bool",
		Value:        "maybe",
		DataType:     "bool",
	}))

	_, err = props.GetBool(ctx, tenantCodigo, "invalid_bool", false)
	require.Error(t, err)
}

func createTenant(t *testing.T, storage connection.StorageManager, tenantCodigo string) {
	t.Helper()

	err := storage.Conn(context.Background()).Exec(`
		insert into tenant
		(
			codigo,
			name,
			timezone,
			max_active_users,
			disabled,
			expires_at,
			created_by,
			created_at
		)
		values
		(
			?,
			?,
			'America/Lima',
			100,
			false,
			null,
			'test',
			now()
		)
	`, tenantCodigo, tenantCodigo).Error
	require.NoError(t, err)
}
