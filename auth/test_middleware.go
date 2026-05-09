package auth

import (
	"context"

	"github.com/labstack/echo/v4"
)

const (
	testUsername = "kevin"
	testTenant   = "tenant_default"
	testTimeZone = "America/Lima"
)

type testSecurityMiddleware struct{}

var _ SecurityMiddleware = (*testSecurityMiddleware)(nil)

func NewTestSecurityMiddleware() SecurityMiddleware {
	return &testSecurityMiddleware{}
}

func (*testSecurityMiddleware) Tenant(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, contextUsernameKey, testUsername)
		ctx = context.WithValue(ctx, contextTenantKey, testTenant)
		ctx = context.WithValue(ctx, contextTimeZoneKey, testTimeZone)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func (*testSecurityMiddleware) System(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, contextUsernameKey, testUsername)
		ctx = context.WithValue(ctx, contextTimeZoneKey, testTimeZone)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
