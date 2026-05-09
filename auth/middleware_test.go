package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestDefaultAuthenticatedSecurityMiddlewareTenantContext(t *testing.T) {
	middleware := NewDefaultAuthenticatedSecurityMiddleware()

	server := echo.New()
	server.GET("/", func(c echo.Context) error {
		ctx := c.Request().Context()

		require.Equal(t, "kevin", Username(ctx))
		require.Equal(t, "local", Tenant(ctx))

		location, err := Tz(ctx)
		require.NoError(t, err)
		require.Equal(t, "America/Lima", location.String())

		return c.NoContent(http.StatusNoContent)
	}, middleware.Tenant)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	require.Equal(t, http.StatusNoContent, response.Code)
}

func TestDefaultAuthenticatedSecurityMiddlewareSystemContext(t *testing.T) {
	middleware := NewDefaultAuthenticatedSecurityMiddleware()

	server := echo.New()
	server.GET("/", func(c echo.Context) error {
		ctx := c.Request().Context()

		require.Equal(t, "kevin", Username(ctx))
		require.Equal(t, undefinedSecurityContextValue, Tenant(ctx))

		location, err := Tz(ctx)
		require.NoError(t, err)
		require.Equal(t, "America/Lima", location.String())

		return c.NoContent(http.StatusNoContent)
	}, middleware.System)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	require.Equal(t, http.StatusNoContent, response.Code)
}
