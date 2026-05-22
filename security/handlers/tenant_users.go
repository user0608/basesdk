package handlers

import (
	"basesdk/binds"
	"basesdk/httpapi"
	"basesdk/security/dtos"
	securitypermissions "basesdk/security/permissions"
	"basesdk/security/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/user0608/goones/answer"
	"github.com/user0608/goones/kcheck"
)

func TenantUsersListHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/users",
		RequiredPerms: []string{securitypermissions.SecurityUsersRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.List(c.Request().Context(), tenant(c))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantUserCreateHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPost,
		Path:          "/api/v1/users",
		RequiredPerms: []string{securitypermissions.SecurityUsersCreate},
		Handler: func(c echo.Context) error {
			var payload dtos.CreateTenantUserInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := kcheck.Valid(payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Create(c.Request().Context(), tenant(c), payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Created(c)
		},
	}
}

func TenantUserFindHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/users/:username",
		RequiredPerms: []string{securitypermissions.SecurityUsersRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.Find(c.Request().Context(), tenant(c), c.Param("username"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantUserUpdateHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPut,
		Path:          "/api/v1/users/:username",
		RequiredPerms: []string{securitypermissions.SecurityUsersUpdate},
		Handler: func(c echo.Context) error {
			var payload dtos.UpdateTenantUserInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := kcheck.Valid(payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Update(c.Request().Context(), tenant(c), c.Param("username"), payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func TenantUserPasswordHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPatch,
		Path:          "/api/v1/users/:username/password",
		RequiredPerms: []string{securitypermissions.SecurityUsersPasswordUpdate},
		Handler: func(c echo.Context) error {
			var payload dtos.ChangePasswordInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := kcheck.Valid(payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ChangePassword(c.Request().Context(), tenant(c), c.Param("username"), payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func TenantUsersEnableHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPatch,
		Path:          "/api/v1/users/enable",
		RequiredPerms: []string{securitypermissions.SecurityUsersEnable},
		Handler: func(c echo.Context) error {
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Enable(c.Request().Context(), tenant(c), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func TenantUsersDisableHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPatch,
		Path:          "/api/v1/users/disable",
		RequiredPerms: []string{securitypermissions.SecurityUsersDisable},
		Handler: func(c echo.Context) error {
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Disable(c.Request().Context(), tenant(c), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func TenantUsersDeleteHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodDelete,
		Path:          "/api/v1/users",
		RequiredPerms: []string{securitypermissions.SecurityUsersDelete},
		Handler: func(c echo.Context) error {
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Delete(c.Request().Context(), tenant(c), values); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func TenantUserPermissionsHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/users/:username/permissions",
		RequiredPerms: []string{securitypermissions.SecurityUsersPermissionsRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.FindUserPermissions(c.Request().Context(), tenant(c), c.Param("username"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantUserRolesHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/users/:username/roles",
		RequiredPerms: []string{securitypermissions.SecurityUsersRolesRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.FindRoles(c.Request().Context(), tenant(c), c.Param("username"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantUserReplaceRolesHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPut,
		Path:          "/api/v1/users/:username/roles",
		RequiredPerms: []string{securitypermissions.SecurityUsersRolesReplace},
		Handler: func(c echo.Context) error {
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplaceRoles(c.Request().Context(), tenant(c), c.Param("username"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func TenantUserGroupsHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/users/:username/groups",
		RequiredPerms: []string{securitypermissions.SecurityUsersGroupsRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.FindGroups(c.Request().Context(), tenant(c), c.Param("username"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantUserReplaceGroupsHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPut,
		Path:          "/api/v1/users/:username/groups",
		RequiredPerms: []string{securitypermissions.SecurityUsersGroupsReplace},
		Handler: func(c echo.Context) error {
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplaceGroups(c.Request().Context(), tenant(c), c.Param("username"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantUsersListHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			response, err := usecase.List(c.Request().Context(), tenantCodigo)
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemTenantUserCreateHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPost,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			var payload dtos.CreateTenantUserInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := kcheck.Valid(payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Create(c.Request().Context(), tenantCodigo, payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Created(c)
		},
	}
}

func SystemTenantUserFindHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users/:username",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			response, err := usecase.Find(c.Request().Context(), tenantCodigo, c.Param("username"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemTenantUserUpdateHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPut,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users/:username",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			var payload dtos.UpdateTenantUserInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := kcheck.Valid(payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Update(c.Request().Context(), tenantCodigo, c.Param("username"), payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantUserPasswordHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPatch,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users/:username/password",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			var payload dtos.ChangePasswordInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := kcheck.Valid(payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ChangePassword(c.Request().Context(), tenantCodigo, c.Param("username"), payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantUsersEnableHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPatch,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users/enable",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Enable(c.Request().Context(), tenantCodigo, values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantUsersDisableHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPatch,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users/disable",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Disable(c.Request().Context(), tenantCodigo, values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantUsersDeleteHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodDelete,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Delete(c.Request().Context(), tenantCodigo, values); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantUserPermissionsHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users/:username/permissions",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			response, err := usecase.FindUserPermissions(c.Request().Context(), tenantCodigo, c.Param("username"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemTenantUserRolesHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users/:username/roles",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			response, err := usecase.FindRoles(c.Request().Context(), tenantCodigo, c.Param("username"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemTenantUserReplaceRolesHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPut,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users/:username/roles",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplaceRoles(c.Request().Context(), tenantCodigo, c.Param("username"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantUserGroupsHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users/:username/groups",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			response, err := usecase.FindGroups(c.Request().Context(), tenantCodigo, c.Param("username"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemTenantUserReplaceGroupsHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPut,
		Path:   "/api/v1/system/tenants/:tenantCodigo/users/:username/groups",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplaceGroups(c.Request().Context(), tenantCodigo, c.Param("username"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}
