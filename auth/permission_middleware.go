package auth

import (
	"basesdk/answer"
	"context"

	"basesdk/errs"

	"github.com/labstack/echo/v4"
)

type PermissionValidator interface {
	HasAllPermissions(ctx context.Context, tenantCodigo string, username string, permissions []string) (bool, error)
	HasAnyPermission(ctx context.Context, tenantCodigo string, username string, permissions []string) (bool, error)
}

type PermissionMiddleware struct {
	validator PermissionValidator
}

func NewPermissionMiddleware(validator PermissionValidator) *PermissionMiddleware {
	return &PermissionMiddleware{validator: validator}
}

func (m *PermissionMiddleware) RequireAll(permissions []string) echo.MiddlewareFunc {
	return m.require(permissions, m.validator.HasAllPermissions)
}

func (m *PermissionMiddleware) RequireAny(permissions []string) echo.MiddlewareFunc {
	return m.require(permissions, m.validator.HasAnyPermission)
}

func (m *PermissionMiddleware) require(
	permissions []string,
	validate func(context.Context, string, string, []string) (bool, error),
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if len(permissions) == 0 {
				return next(c)
			}

			ctx := c.Request().Context()
			username := Username(ctx)
			tenantCodigo := Tenant(ctx)
			if IsUndefined(username) || IsUndefined(tenantCodigo) {
				return answer.Err(c, errs.ForbiddenDirect("no se pudo resolver el contexto de seguridad"))
			}

			allowed, err := validate(ctx, tenantCodigo, username, permissions)
			if err != nil {
				return answer.Err(c, err)
			}
			if !allowed {
				return answer.Err(c, errs.ForbiddenDirect("no tiene permisos para realizar esta operación"))
			}

			return next(c)
		}
	}
}
