package handlers

import (
	"basesdk/httpapi"
	"basesdk/security/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/user0608/goones/answer"
)

func permissionsRoute(path string, system bool, find bool, usecase *usecases.PermissionsUsecase) httpapi.Route {
	handler := func(c echo.Context) error {
		if find {
			response, err := usecase.Find(c.Request().Context(), c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		}

		response, err := usecase.List(c.Request().Context())
		if err != nil {
			return answer.Err(c, err)
		}
		return answer.Ok(c, response)
	}
	if system {
		return &httpapi.SystemHandler{Method: http.MethodGet, Path: path, Handler: handler}
	}
	return &httpapi.TenantHandler{Method: http.MethodGet, Path: path, Handler: handler}
}

func TenantPermissionsListHandler(usecase *usecases.PermissionsUsecase) httpapi.Route {
	return permissionsRoute("/api/v1/permissions", false, false, usecase)
}
func TenantPermissionFindHandler(usecase *usecases.PermissionsUsecase) httpapi.Route {
	return permissionsRoute("/api/v1/permissions/:code", false, true, usecase)
}
func SystemPermissionsListHandler(usecase *usecases.PermissionsUsecase) httpapi.Route {
	return permissionsRoute("/api/v1/system/permissions", true, false, usecase)
}
func SystemPermissionFindHandler(usecase *usecases.PermissionsUsecase) httpapi.Route {
	return permissionsRoute("/api/v1/system/permissions/:code", true, true, usecase)
}
