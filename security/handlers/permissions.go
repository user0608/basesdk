package handlers

import (
	"basesdk/httpapi"
	securitypermissions "basesdk/security/permissions"
	"basesdk/security/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/user0608/goones/answer"
)

func TenantPermissionsListHandler(usecase *usecases.PermissionsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/permissions",
		RequiredPerms: []string{securitypermissions.SecurityPermissionsRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.List(c.Request().Context())
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantPermissionFindHandler(usecase *usecases.PermissionsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/permissions/:code",
		RequiredPerms: []string{securitypermissions.SecurityPermissionsRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.Find(c.Request().Context(), c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemPermissionsListHandler(usecase *usecases.PermissionsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/permissions",
		Handler: func(c echo.Context) error {
			response, err := usecase.List(c.Request().Context())
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemPermissionFindHandler(usecase *usecases.PermissionsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/permissions/:code",
		Handler: func(c echo.Context) error {
			response, err := usecase.Find(c.Request().Context(), c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}
