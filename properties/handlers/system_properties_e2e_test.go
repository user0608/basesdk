package handlers_test

import (
	"basesdk/auth"
	"basesdk/httpapi"
	"basesdk/properties"
	"basesdk/properties/handlers"
	"basesdk/testdb"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type responsePayload struct {
	Data json.RawMessage `json:"data"`
}

func TestSystemPropertiesHandlersE2E(t *testing.T) {
	server := newPropertiesTestServer(t)

	createProperty := executeJSONRequest(t, server, http.MethodPost, "/api/v1/system/properties", `{
		"key": "e2e_global_property",
		"value": "true",
		"dataType": "bool",
		"description": "E2E global property"
	}`)
	require.Equal(t, http.StatusCreated, createProperty.Code)

	findProperty := executeJSONRequest(t, server, http.MethodGet, "/api/v1/system/properties/e2e_global_property", ``)
	require.Equal(t, http.StatusOK, findProperty.Code)
	requireJSONField(t, findProperty, "key", "e2e_global_property")

	updateProperty := executeJSONRequest(t, server, http.MethodPut, "/api/v1/system/properties/e2e_global_property", `{
		"key": "ignored_key",
		"value": "42",
		"dataType": "int",
		"description": "Updated E2E global property"
	}`)
	require.Equal(t, http.StatusOK, updateProperty.Code)

	findUpdatedProperty := executeJSONRequest(t, server, http.MethodGet, "/api/v1/system/properties/e2e_global_property", ``)
	require.Equal(t, http.StatusOK, findUpdatedProperty.Code)
	requireJSONField(t, findUpdatedProperty, "value", "42")

	invalidProperty := executeJSONRequest(t, server, http.MethodPost, "/api/v1/system/properties", `{
		"key": "e2e_invalid_property",
		"value": "not-an-int",
		"dataType": "int"
	}`)
	require.GreaterOrEqual(t, invalidProperty.Code, http.StatusBadRequest)

	deleteProperty := executeJSONRequest(t, server, http.MethodDelete, "/api/v1/system/properties/e2e_global_property", ``)
	require.Equal(t, http.StatusOK, deleteProperty.Code)
}

func TestSystemTenantPropertiesHandlersE2E(t *testing.T) {
	server := newPropertiesTestServer(t)

	createProperty := executeJSONRequest(t, server, http.MethodPost, "/api/v1/system/tenants/tenant_default/properties", `{
		"key": "e2e_tenant_property",
		"value": "{\"enabled\":true}",
		"dataType": "json",
		"description": "E2E tenant property"
	}`)
	require.Equal(t, http.StatusCreated, createProperty.Code)

	findProperty := executeJSONRequest(t, server, http.MethodGet, "/api/v1/system/tenants/tenant_default/properties/e2e_tenant_property", ``)
	require.Equal(t, http.StatusOK, findProperty.Code)
	requireJSONField(t, findProperty, "tenantCodigo", "tenant_default")

	updateProperty := executeJSONRequest(t, server, http.MethodPut, "/api/v1/system/tenants/tenant_default/properties/e2e_tenant_property", `{
		"key": "ignored_key",
		"value": "disabled",
		"dataType": "bool",
		"description": "Updated E2E tenant property"
	}`)
	require.Equal(t, http.StatusOK, updateProperty.Code)

	findUpdatedProperty := executeJSONRequest(t, server, http.MethodGet, "/api/v1/system/tenants/tenant_default/properties/e2e_tenant_property", ``)
	require.Equal(t, http.StatusOK, findUpdatedProperty.Code)
	requireJSONField(t, findUpdatedProperty, "value", "disabled")

	listProperties := executeJSONRequest(t, server, http.MethodGet, "/api/v1/system/tenants/tenant_default/properties", ``)
	require.Equal(t, http.StatusOK, listProperties.Code)

	deleteProperty := executeJSONRequest(t, server, http.MethodDelete, "/api/v1/system/tenants/tenant_default/properties/e2e_tenant_property", ``)
	require.Equal(t, http.StatusOK, deleteProperty.Code)
}

func newPropertiesTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	storage := testdb.NewPostgresStorage(t)
	systemProps := properties.NewSystemProperties(storage)
	tenantProps := properties.NewTenantSystemProperties(storage)

	server := httpapi.NewServer(
		[]httpapi.Route{
			handlers.SystemPropertiesListHandler(systemProps),
			handlers.SystemPropertyCreateHandler(systemProps),
			handlers.SystemPropertyFindHandler(systemProps),
			handlers.SystemPropertyUpdateHandler(systemProps),
			handlers.SystemPropertyDeleteHandler(systemProps),
			handlers.SystemTenantPropertiesListHandler(tenantProps),
			handlers.SystemTenantPropertyCreateHandler(tenantProps),
			handlers.SystemTenantPropertyFindHandler(tenantProps),
			handlers.SystemTenantPropertyUpdateHandler(tenantProps),
			handlers.SystemTenantPropertyDeleteHandler(tenantProps),
		},
		auth.NewTestSecurityMiddleware(),
		nil,
	)

	testServer := httptest.NewServer(server)
	t.Cleanup(testServer.Close)

	return testServer
}

func executeJSONRequest(t *testing.T, server *httptest.Server, method string, path string, body string) *httptest.ResponseRecorder {
	t.Helper()

	var requestBody *bytes.Buffer
	if body == "" {
		requestBody = bytes.NewBuffer(nil)
	} else {
		requestBody = bytes.NewBufferString(body)
	}

	request := httptest.NewRequest(method, path, requestBody)
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}

	response := httptest.NewRecorder()
	server.Config.Handler.ServeHTTP(response, request)

	return response
}

func requireJSONField[T comparable](t *testing.T, response *httptest.ResponseRecorder, field string, expected T) {
	t.Helper()

	var payload responsePayload
	require.NoError(t, json.NewDecoder(response.Body).Decode(&payload))

	var data map[string]any
	require.NoError(t, json.Unmarshal(payload.Data, &data))
	require.Equal(t, expected, data[field])
}
