package properties_test

import (
	"basesdk/properties"
	"basesdk/testdb"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSystemPropertiesCRUD(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	props := properties.NewSystemProperties(storage)
	ctx := context.Background()

	description := "External API base URL"
	property := &properties.Property{
		Key:         "external_api_base_url",
		Value:       "https://api.example.test",
		DataType:    "string",
		Description: &description,
	}

	require.NoError(t, props.Create(ctx, property))

	found, err := props.Get(ctx, property.Key)
	require.NoError(t, err)
	require.Equal(t, property.Key, found.Key)
	require.Equal(t, property.Value, found.Value)
	require.Equal(t, property.DataType, found.DataType)
	require.NotNil(t, found.Description)
	require.Equal(t, description, *found.Description)

	exists, err := props.Exists(ctx, property.Key)
	require.NoError(t, err)
	require.True(t, exists)

	property.Value = "https://api-updated.example.test"
	require.NoError(t, props.Update(ctx, property))

	value, err := props.GetString(ctx, property.Key, "")
	require.NoError(t, err)
	require.Equal(t, property.Value, value)

	require.NoError(t, props.Delete(ctx, property.Key))

	exists, err = props.Exists(ctx, property.Key)
	require.NoError(t, err)
	require.False(t, exists)
}

func TestSystemPropertiesTypedReaders(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	props := properties.NewSystemProperties(storage)
	ctx := context.Background()

	entries := []*properties.Property{
		{Key: "test_int", Value: " 42 ", DataType: "int"},
		{Key: "test_float", Value: " 12.50 ", DataType: "float"},
		{Key: "test_bool", Value: "enabled", DataType: "bool"},
		{Key: "test_json", Value: `{"name":"storage","enabled":true}`, DataType: "json"},
		{Key: "test_time", Value: "2026-05-08T10:30:00Z", DataType: "string"},
		{Key: "test_duration", Value: "45m", DataType: "string"},
	}

	for _, entry := range entries {
		require.NoError(t, props.Create(ctx, entry))
	}

	intValue, err := props.GetInt(ctx, "test_int", 0)
	require.NoError(t, err)
	require.Equal(t, 42, intValue)

	floatValue, err := props.GetFloat(ctx, "test_float", 0)
	require.NoError(t, err)
	require.InDelta(t, 12.50, floatValue, 0.001)

	boolValue, err := props.GetBool(ctx, "test_bool", false)
	require.NoError(t, err)
	require.True(t, boolValue)

	var jsonValue struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}
	require.NoError(t, props.GetJSON(ctx, "test_json", &jsonValue))
	require.Equal(t, "storage", jsonValue.Name)
	require.True(t, jsonValue.Enabled)

	timeValue, err := props.GetTime(ctx, "test_time", time.Time{})
	require.NoError(t, err)
	require.Equal(t, time.Date(2026, 5, 8, 10, 30, 0, 0, time.UTC), timeValue)

	durationValue, err := props.GetDuration(ctx, "test_duration", 0)
	require.NoError(t, err)
	require.Equal(t, 45*time.Minute, durationValue)
}

func TestSystemPropertiesDefaultsAndInvalidValues(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	props := properties.NewSystemProperties(storage)
	ctx := context.Background()

	stringValue, err := props.GetString(ctx, "missing_string", "default")
	require.NoError(t, err)
	require.Equal(t, "default", stringValue)

	intValue, err := props.GetInt(ctx, "missing_int", 7)
	require.NoError(t, err)
	require.Equal(t, 7, intValue)

	require.NoError(t, props.Create(ctx, &properties.Property{
		Key:      "invalid_int",
		Value:    "abc",
		DataType: "int",
	}))

	_, err = props.GetInt(ctx, "invalid_int", 0)
	require.Error(t, err)
}

func TestSystemPropertiesSeededValues(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	props := properties.NewSystemProperties(storage)
	ctx := context.Background()

	durationValue, err := props.GetDuration(ctx, "jwt_token_ttl", 0)
	require.NoError(t, err)
	require.Equal(t, 720*time.Hour, durationValue)
}
