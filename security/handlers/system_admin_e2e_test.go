package handlers_test

import (
	"basesdk/auth"
	"basesdk/httpapi"
	"basesdk/security/handlers"
	"basesdk/security/repositories"
	"basesdk/security/usecases"
	"basesdk/testdb"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemTenantManagementE2E(t *testing.T) {
	server := newSystemAdminTestServer(t)

	createTenant := executeJSONRequest(t, server, http.MethodPost, "/api/v1/system/tenants", `{
		"codigo": "tenant_e2e",
		"name": "Tenant E2E",
		"timezone": "America/Lima",
		"maxActiveUsers": 25
	}`)
	require.Equal(t, http.StatusCreated, createTenant.Code)

	findTenant := executeJSONRequest(t, server, http.MethodGet, "/api/v1/system/tenants/tenant_e2e", ``)
	require.Equal(t, http.StatusOK, findTenant.Code)
	requireJSONField(t, findTenant, "codigo", "tenant_e2e")

	updateTenant := executeJSONRequest(t, server, http.MethodPut, "/api/v1/system/tenants/tenant_e2e", `{
		"name": "Tenant E2E Updated",
		"timezone": "America/Lima",
		"maxActiveUsers": 50,
		"disabled": false
	}`)
	require.Equal(t, http.StatusOK, updateTenant.Code)

	disableTenant := executeJSONRequest(t, server, http.MethodPatch, "/api/v1/system/tenants/disable", `{"codigos": ["tenant_e2e"]}`)
	require.Equal(t, http.StatusOK, disableTenant.Code)

	findDisabledTenant := executeJSONRequest(t, server, http.MethodGet, "/api/v1/system/tenants/tenant_e2e", ``)
	require.Equal(t, http.StatusOK, findDisabledTenant.Code)
	requireJSONField(t, findDisabledTenant, "disabled", true)

	enableTenant := executeJSONRequest(t, server, http.MethodPatch, "/api/v1/system/tenants/enable", `{"codigos": ["tenant_e2e"]}`)
	require.Equal(t, http.StatusOK, enableTenant.Code)

	listTenants := executeJSONRequest(t, server, http.MethodGet, "/api/v1/system/tenants", ``)
	require.Equal(t, http.StatusOK, listTenants.Code)
}

func TestSystemUsersManagementE2E(t *testing.T) {
	server := newSystemAdminTestServer(t)

	createUser := executeJSONRequest(t, server, http.MethodPost, "/api/v1/system/users", `{
		"username": "system_e2e",
		"password": "secret123"
	}`)
	require.Equal(t, http.StatusCreated, createUser.Code)

	findUser := executeJSONRequest(t, server, http.MethodGet, "/api/v1/system/users/system_e2e", ``)
	require.Equal(t, http.StatusOK, findUser.Code)
	requireJSONField(t, findUser, "username", "system_e2e")

	disableUser := executeJSONRequest(t, server, http.MethodPatch, "/api/v1/system/users/disable", `{"usernames": ["system_e2e"]}`)
	require.Equal(t, http.StatusOK, disableUser.Code)

	findDisabledUser := executeJSONRequest(t, server, http.MethodGet, "/api/v1/system/users/system_e2e", ``)
	require.Equal(t, http.StatusOK, findDisabledUser.Code)
	requireJSONField(t, findDisabledUser, "disabled", true)

	enableUser := executeJSONRequest(t, server, http.MethodPatch, "/api/v1/system/users/enable", `{"usernames": ["system_e2e"]}`)
	require.Equal(t, http.StatusOK, enableUser.Code)

	deleteUser := executeJSONRequest(t, server, http.MethodDelete, "/api/v1/system/users", `{"usernames": ["system_e2e"]}`)
	require.Equal(t, http.StatusOK, deleteUser.Code)

	missingUser := executeJSONRequest(t, server, http.MethodGet, "/api/v1/system/users/system_e2e", ``)
	require.GreaterOrEqual(t, missingUser.Code, http.StatusBadRequest)
}

func TestSystemUsersRejectSelfDisableAndDeleteE2E(t *testing.T) {
	server := newSystemAdminTestServer(t)

	selfDisable := executeJSONRequest(t, server, http.MethodPatch, "/api/v1/system/users/disable", `{"usernames": ["kevin"]}`)
	require.GreaterOrEqual(t, selfDisable.Code, http.StatusBadRequest)

	selfDelete := executeJSONRequest(t, server, http.MethodDelete, "/api/v1/system/users", `{"usernames": ["kevin"]}`)
	require.GreaterOrEqual(t, selfDelete.Code, http.StatusBadRequest)

	findCurrentUser := executeJSONRequest(t, server, http.MethodGet, "/api/v1/system/users/kevin", ``)
	require.Equal(t, http.StatusOK, findCurrentUser.Code)
	requireJSONField(t, findCurrentUser, "disabled", false)
}

func newSystemAdminTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	storage := testdb.NewPostgresStorage(t)
	systemUsersUsecase := usecases.NewSystemUsersUsecase(repositories.NewSystemUserRepository(storage))
	tenantsUsecase := usecases.NewTenantsUsecase(repositories.NewTenantRepository(storage))

	server := httpapi.NewServer(
		[]httpapi.Route{
			handlers.SystemUsersListHandler(systemUsersUsecase),
			handlers.SystemUserCreateHandler(systemUsersUsecase),
			handlers.SystemUserFindHandler(systemUsersUsecase),
			handlers.SystemUserUpdateHandler(systemUsersUsecase),
			handlers.SystemUsersEnableHandler(systemUsersUsecase),
			handlers.SystemUsersDisableHandler(systemUsersUsecase),
			handlers.SystemUsersDeleteHandler(systemUsersUsecase),
			handlers.SystemTenantsListHandler(tenantsUsecase),
			handlers.SystemTenantCreateHandler(tenantsUsecase),
			handlers.SystemTenantFindHandler(tenantsUsecase),
			handlers.SystemTenantUpdateHandler(tenantsUsecase),
			handlers.SystemTenantsEnableHandler(tenantsUsecase),
			handlers.SystemTenantsDisableHandler(tenantsUsecase),
		},
		auth.NewTestSecurityMiddleware(),
		nil,
	)

	testServer := httptest.NewServer(server)
	t.Cleanup(testServer.Close)

	return testServer
}
