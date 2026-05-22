package handlers_test

import (
	"basesdk/auth"
	"basesdk/connection"
	"basesdk/httpapi"
	"basesdk/security/handlers"
	"basesdk/security/repositories"
	"basesdk/security/usecases"
	"basesdk/testdb"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type responsePayload struct {
	Data json.RawMessage `json:"data"`
}

func TestTenantUserManagementE2E(t *testing.T) {
	server := newUserManagementTestServer(t)

	createUser := executeJSONRequest(t, server, http.MethodPost, "/api/v1/users", `{
		"username": "ana",
		"fullName": "Ana",
		"password": "secret123"
	}`)
	require.Equal(t, http.StatusCreated, createUser.Code)

	findUser := executeJSONRequest(t, server, http.MethodGet, "/api/v1/users/ana", ``)
	require.Equal(t, http.StatusOK, findUser.Code)
	requireJSONField(t, findUser, "username", "ana")

	disableUser := executeJSONRequest(t, server, http.MethodPatch, "/api/v1/users/disable", `{"usernames": ["ana"]}`)
	require.Equal(t, http.StatusOK, disableUser.Code)

	findDisabledUser := executeJSONRequest(t, server, http.MethodGet, "/api/v1/users/ana", ``)
	require.Equal(t, http.StatusOK, findDisabledUser.Code)
	requireJSONField(t, findDisabledUser, "disabled", true)

	enableUser := executeJSONRequest(t, server, http.MethodPatch, "/api/v1/users/enable", `{"usernames": ["ana"]}`)
	require.Equal(t, http.StatusOK, enableUser.Code)

	deleteUser := executeJSONRequest(t, server, http.MethodDelete, "/api/v1/users", `{"usernames": ["ana"]}`)
	require.Equal(t, http.StatusOK, deleteUser.Code)

	missingUser := executeJSONRequest(t, server, http.MethodGet, "/api/v1/users/ana", ``)
	require.GreaterOrEqual(t, missingUser.Code, http.StatusBadRequest)
}

func TestTenantRoleGroupAndPermissionsE2E(t *testing.T) {
	server := newUserManagementTestServer(t)

	createRole := executeJSONRequest(t, server, http.MethodPost, "/api/v1/roles", `{
		"code": "DEV",
		"description": "Developers"
	}`)
	require.Equal(t, http.StatusCreated, createRole.Code)

	createGroup := executeJSONRequest(t, server, http.MethodPost, "/api/v1/groups", `{
		"code": "ENGINEERING",
		"description": "Engineering team"
	}`)
	require.Equal(t, http.StatusCreated, createGroup.Code)

	replaceGroupUsers := executeJSONRequest(t, server, http.MethodPut, "/api/v1/groups/ENGINEERING/users", `{"usernames": ["kevin"]}`)
	require.Equal(t, http.StatusOK, replaceGroupUsers.Code)

	replaceGroupRoles := executeJSONRequest(t, server, http.MethodPut, "/api/v1/groups/ENGINEERING/roles", `{"codes": ["DEV"]}`)
	require.Equal(t, http.StatusOK, replaceGroupRoles.Code)

	replaceRolePermissions := executeJSONRequest(t, server, http.MethodPut, "/api/v1/roles/DEV/permissions", `{"codes": ["security.users.read"]}`)
	require.Equal(t, http.StatusOK, replaceRolePermissions.Code)

	replaceUserRoles := executeJSONRequest(t, server, http.MethodPut, "/api/v1/users/kevin/roles", `{"codes": ["DEV"]}`)
	require.Equal(t, http.StatusOK, replaceUserRoles.Code)

	replaceUserGroups := executeJSONRequest(t, server, http.MethodPut, "/api/v1/users/kevin/groups", `{"codes": ["ENGINEERING"]}`)
	require.Equal(t, http.StatusOK, replaceUserGroups.Code)

	groupUsers := executeJSONRequest(t, server, http.MethodGet, "/api/v1/groups/ENGINEERING/users", ``)
	require.Equal(t, http.StatusOK, groupUsers.Code)
	requireJSONArrayLen(t, groupUsers, 1)

	groupRoles := executeJSONRequest(t, server, http.MethodGet, "/api/v1/groups/ENGINEERING/roles", ``)
	require.Equal(t, http.StatusOK, groupRoles.Code)
	requireJSONArrayLen(t, groupRoles, 1)

	groupPermissions := executeJSONRequest(t, server, http.MethodGet, "/api/v1/groups/ENGINEERING/permissions", ``)
	require.Equal(t, http.StatusOK, groupPermissions.Code)
	requireJSONArrayLen(t, groupPermissions, 1)

	rolePermissions := executeJSONRequest(t, server, http.MethodGet, "/api/v1/roles/DEV/permissions", ``)
	require.Equal(t, http.StatusOK, rolePermissions.Code)
	requireJSONArrayLen(t, rolePermissions, 1)

	userRoles := executeJSONRequest(t, server, http.MethodGet, "/api/v1/users/kevin/roles", ``)
	require.Equal(t, http.StatusOK, userRoles.Code)
	requireJSONArrayLen(t, userRoles, 1)

	userGroups := executeJSONRequest(t, server, http.MethodGet, "/api/v1/users/kevin/groups", ``)
	require.Equal(t, http.StatusOK, userGroups.Code)
	requireJSONArrayLen(t, userGroups, 1)

	roleUsers := executeJSONRequest(t, server, http.MethodGet, "/api/v1/roles/DEV/users", ``)
	require.Equal(t, http.StatusOK, roleUsers.Code)
	requireJSONArrayLen(t, roleUsers, 1)

	roleGroups := executeJSONRequest(t, server, http.MethodGet, "/api/v1/roles/DEV/groups", ``)
	require.Equal(t, http.StatusOK, roleGroups.Code)
	requireJSONArrayLen(t, roleGroups, 1)

	permissions := executeJSONRequest(t, server, http.MethodGet, "/api/v1/permissions", ``)
	require.Equal(t, http.StatusOK, permissions.Code)

	userPermissions := executeJSONRequest(t, server, http.MethodGet, "/api/v1/me/permissions", ``)
	require.Equal(t, http.StatusOK, userPermissions.Code)

	permissionRoles := executeJSONRequest(t, server, http.MethodGet, "/api/v1/permissions/security.users.read/roles", ``)
	require.Equal(t, http.StatusOK, permissionRoles.Code)
	requireJSONArrayLen(t, permissionRoles, 1)

	permissionGroups := executeJSONRequest(t, server, http.MethodGet, "/api/v1/permissions/security.users.read/groups", ``)
	require.Equal(t, http.StatusOK, permissionGroups.Code)
	requireJSONArrayLen(t, permissionGroups, 1)

	permissionUsers := executeJSONRequest(t, server, http.MethodGet, "/api/v1/permissions/security.users.read/users", ``)
	require.Equal(t, http.StatusOK, permissionUsers.Code)
	requireJSONArrayLen(t, permissionUsers, 1)
}

func newUserManagementTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	storage := testdb.NewPostgresStorage(t)
	seedTestPermission(t, storage)

	systemUsersUsecase := usecases.NewSystemUsersUsecase(repositories.NewSystemUserRepository(storage))
	tenantUsersUsecase := usecases.NewTenantUsersUsecase(
		repositories.NewAppUserRepository(storage),
		repositories.NewPermissionRepository(storage),
	)
	tenantRolesUsecase := usecases.NewTenantRolesUsecase(repositories.NewRoleRepository(storage))
	tenantGroupsUsecase := usecases.NewTenantGroupsUsecase(repositories.NewGroupRepository(storage))
	permissionsUsecase := usecases.NewPermissionsUsecase(repositories.NewPermissionRepository(storage))

	server := httpapi.NewServer(
		[]httpapi.Route{
			handlers.SystemUsersListHandler(systemUsersUsecase),
			handlers.TenantUsersListHandler(tenantUsersUsecase),
			handlers.TenantUserCreateHandler(tenantUsersUsecase),
			handlers.TenantUserFindHandler(tenantUsersUsecase),
			handlers.TenantUserUpdateHandler(tenantUsersUsecase),
			handlers.TenantUserPasswordHandler(tenantUsersUsecase),
			handlers.TenantUsersEnableHandler(tenantUsersUsecase),
			handlers.TenantUsersDisableHandler(tenantUsersUsecase),
			handlers.TenantUsersDeleteHandler(tenantUsersUsecase),
			handlers.TenantUserPermissionsHandler(tenantUsersUsecase),
			handlers.TenantUserRolesHandler(tenantUsersUsecase),
			handlers.TenantUserReplaceRolesHandler(tenantUsersUsecase),
			handlers.TenantUserGroupsHandler(tenantUsersUsecase),
			handlers.TenantUserReplaceGroupsHandler(tenantUsersUsecase),
			handlers.TenantMeHandler(tenantUsersUsecase),
			handlers.TenantMePermissionsHandler(tenantUsersUsecase),
			handlers.TenantRolesListHandler(tenantRolesUsecase),
			handlers.TenantRoleCreateHandler(tenantRolesUsecase),
			handlers.TenantRoleFindHandler(tenantRolesUsecase),
			handlers.TenantRoleUpdateHandler(tenantRolesUsecase),
			handlers.TenantRolesEnableHandler(tenantRolesUsecase),
			handlers.TenantRolesDisableHandler(tenantRolesUsecase),
			handlers.TenantRolesDeleteHandler(tenantRolesUsecase),
			handlers.TenantRolePermissionsHandler(tenantRolesUsecase),
			handlers.TenantRoleReplacePermissionsHandler(tenantRolesUsecase),
			handlers.TenantRoleUsersHandler(tenantRolesUsecase),
			handlers.TenantRoleReplaceUsersHandler(tenantRolesUsecase),
			handlers.TenantRoleGroupsHandler(tenantRolesUsecase),
			handlers.TenantRoleReplaceGroupsHandler(tenantRolesUsecase),
			handlers.TenantGroupsListHandler(tenantGroupsUsecase),
			handlers.TenantGroupCreateHandler(tenantGroupsUsecase),
			handlers.TenantGroupFindHandler(tenantGroupsUsecase),
			handlers.TenantGroupUpdateHandler(tenantGroupsUsecase),
			handlers.TenantGroupsEnableHandler(tenantGroupsUsecase),
			handlers.TenantGroupsDisableHandler(tenantGroupsUsecase),
			handlers.TenantGroupsDeleteHandler(tenantGroupsUsecase),
			handlers.TenantGroupUsersHandler(tenantGroupsUsecase),
			handlers.TenantGroupReplaceUsersHandler(tenantGroupsUsecase),
			handlers.TenantGroupRolesHandler(tenantGroupsUsecase),
			handlers.TenantGroupReplaceRolesHandler(tenantGroupsUsecase),
			handlers.TenantGroupPermissionsHandler(tenantGroupsUsecase),
			handlers.TenantPermissionsListHandler(permissionsUsecase),
			handlers.TenantPermissionFindHandler(permissionsUsecase),
			handlers.TenantPermissionRolesHandler(permissionsUsecase),
			handlers.TenantPermissionGroupsHandler(permissionsUsecase),
			handlers.TenantPermissionUsersHandler(permissionsUsecase),
		},
		auth.NewTestSecurityMiddleware(),
		nil,
	)

	testServer := httptest.NewServer(server)
	t.Cleanup(testServer.Close)

	return testServer
}

func seedTestPermission(t *testing.T, storage connection.StorageManager) {
	t.Helper()

	err := storage.Conn(context.Background()).Exec(`
		insert into permission (code, description)
		values ('security.users.read', 'Consultar usuarios')
	`).Error
	require.NoError(t, err)
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

func requireJSONArrayLen(t *testing.T, response *httptest.ResponseRecorder, expected int) {
	t.Helper()

	var payload responsePayload
	require.NoError(t, json.NewDecoder(response.Body).Decode(&payload))

	var data []any
	require.NoError(t, json.Unmarshal(payload.Data, &data))
	require.Len(t, data, expected)
}
