package httpapi

import "github.com/labstack/echo/v4"

type BeforeSecurityMiddlewareProvider interface {
	BeforeSecurityMiddlewares() []echo.MiddlewareFunc
}

type AfterSecurityMiddlewareProvider interface {
	AfterSecurityMiddlewares() []echo.MiddlewareFunc
}

type PermissionsProvider interface {
	Permissions() []string
}

type AnyPermissionsProvider interface {
	AnyPermissions() []string
}
