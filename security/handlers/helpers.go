package handlers

import (
	"basesdk/auth"
	"basesdk/binds"
	"basesdk/errs"

	"github.com/labstack/echo/v4"
)

func actor(c echo.Context) string {
	return auth.Username(c.Request().Context())
}

func tenant(c echo.Context) string {
	return auth.Tenant(c.Request().Context())
}

func tenantParam(c echo.Context) (string, error) {
	tenantCodigo := c.Param("tenantCodigo")
	if tenantCodigo == "" {
		return "", errs.BadRequestDirect("tenant requerido")
	}
	return tenantCodigo, nil
}

func requestStrings(c echo.Context) ([]string, error) {
	values, err := binds.RequestStrings(c)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, errs.BadRequestDirect("valores requeridos")
	}
	return values, nil
}
