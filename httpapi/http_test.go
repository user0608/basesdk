package httpapi_test

import (
	"basesdk/auth"
	"basesdk/httpapi"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type permissionValidatorStub struct {
	permissions map[string]bool
}

func (v *permissionValidatorStub) HasAllPermissions(_ context.Context, _ string, _ string, permissions []string) (bool, error) {
	for _, permission := range permissions {
		if !v.permissions[permission] {
			return false, nil
		}
	}

	return true, nil
}

func (v *permissionValidatorStub) HasAnyPermission(_ context.Context, _ string, _ string, permissions []string) (bool, error) {
	for _, permission := range permissions {
		if v.permissions[permission] {
			return true, nil
		}
	}

	return false, nil
}

func TestNewServerAppliesRequiredPermissions(t *testing.T) {
	server := httpapi.NewServer(
		[]httpapi.Route{
			&httpapi.TenantHandler{
				Method:        http.MethodGet,
				Path:          "/users",
				RequiredPerms: []string{"users.read", "users.create"},
				Handler: func(c echo.Context) error {
					return c.NoContent(http.StatusNoContent)
				},
			},
		},
		auth.NewTestSecurityMiddleware(),
		auth.NewPermissionMiddleware(&permissionValidatorStub{permissions: map[string]bool{
			"users.read": true,
		}}),
	)

	request := httptest.NewRequest(http.MethodGet, "/users", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	require.Equal(t, http.StatusForbidden, response.Code)
}

func TestNewServerAppliesAnyRequiredPermissions(t *testing.T) {
	server := httpapi.NewServer(
		[]httpapi.Route{
			&httpapi.TenantHandler{
				Method:           http.MethodGet,
				Path:             "/users",
				AnyRequiredPerms: []string{"users.create", "users.read"},
				Handler: func(c echo.Context) error {
					return c.NoContent(http.StatusNoContent)
				},
			},
		},
		auth.NewTestSecurityMiddleware(),
		auth.NewPermissionMiddleware(&permissionValidatorStub{permissions: map[string]bool{
			"users.read": true,
		}}),
	)

	request := httptest.NewRequest(http.MethodGet, "/users", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	require.Equal(t, http.StatusNoContent, response.Code)
}
