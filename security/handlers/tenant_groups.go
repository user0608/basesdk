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

func TenantGroupsListHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/groups",
		RequiredPerms: []string{securitypermissions.SecurityGroupsRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.List(c.Request().Context(), tenant(c))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantGroupCreateHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPost,
		Path:          "/api/v1/groups",
		RequiredPerms: []string{securitypermissions.SecurityGroupsCreate},
		Handler: func(c echo.Context) error {
			var payload dtos.CreateGroupInput
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

func TenantGroupFindHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/groups/:code",
		RequiredPerms: []string{securitypermissions.SecurityGroupsRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.Find(c.Request().Context(), tenant(c), c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantGroupUpdateHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPut,
		Path:          "/api/v1/groups/:code",
		RequiredPerms: []string{securitypermissions.SecurityGroupsUpdate},
		Handler: func(c echo.Context) error {
			var payload dtos.UpdateGroupInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Update(c.Request().Context(), tenant(c), c.Param("code"), payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func TenantGroupsEnableHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPatch,
		Path:          "/api/v1/groups/enable",
		RequiredPerms: []string{securitypermissions.SecurityGroupsEnable},
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

func TenantGroupsDisableHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPatch,
		Path:          "/api/v1/groups/disable",
		RequiredPerms: []string{securitypermissions.SecurityGroupsDisable},
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

func TenantGroupsDeleteHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodDelete,
		Path:          "/api/v1/groups",
		RequiredPerms: []string{securitypermissions.SecurityGroupsDelete},
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

func TenantGroupUsersHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/groups/:code/users",
		RequiredPerms: []string{securitypermissions.SecurityGroupsUsersRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.FindUsers(c.Request().Context(), tenant(c), c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantGroupReplaceUsersHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPut,
		Path:          "/api/v1/groups/:code/users",
		RequiredPerms: []string{securitypermissions.SecurityGroupsUsersReplace},
		Handler: func(c echo.Context) error {
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplaceUsers(c.Request().Context(), tenant(c), c.Param("code"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func TenantGroupRolesHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/groups/:code/roles",
		RequiredPerms: []string{securitypermissions.SecurityGroupsRolesRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.FindRoles(c.Request().Context(), tenant(c), c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantGroupReplaceRolesHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPut,
		Path:          "/api/v1/groups/:code/roles",
		RequiredPerms: []string{securitypermissions.SecurityGroupsRolesReplace},
		Handler: func(c echo.Context) error {
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplaceRoles(c.Request().Context(), tenant(c), c.Param("code"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func TenantGroupPermissionsHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/groups/:code/permissions",
		RequiredPerms: []string{securitypermissions.SecurityGroupsPermissionsRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.FindPermissions(c.Request().Context(), tenant(c), c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemTenantGroupsListHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/groups",
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

func SystemTenantGroupCreateHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPost,
		Path:   "/api/v1/system/tenants/:tenantCodigo/groups",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			var payload dtos.CreateGroupInput
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

func SystemTenantGroupFindHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/groups/:code",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			response, err := usecase.Find(c.Request().Context(), tenantCodigo, c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemTenantGroupUpdateHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPut,
		Path:   "/api/v1/system/tenants/:tenantCodigo/groups/:code",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			var payload dtos.UpdateGroupInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Update(c.Request().Context(), tenantCodigo, c.Param("code"), payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantGroupsEnableHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPatch,
		Path:   "/api/v1/system/tenants/:tenantCodigo/groups/enable",
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

func SystemTenantGroupsDisableHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPatch,
		Path:   "/api/v1/system/tenants/:tenantCodigo/groups/disable",
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

func SystemTenantGroupsDeleteHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodDelete,
		Path:   "/api/v1/system/tenants/:tenantCodigo/groups",
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

func SystemTenantGroupUsersHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/groups/:code/users",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			response, err := usecase.FindUsers(c.Request().Context(), tenantCodigo, c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemTenantGroupReplaceUsersHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPut,
		Path:   "/api/v1/system/tenants/:tenantCodigo/groups/:code/users",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplaceUsers(c.Request().Context(), tenantCodigo, c.Param("code"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantGroupRolesHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/groups/:code/roles",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			response, err := usecase.FindRoles(c.Request().Context(), tenantCodigo, c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemTenantGroupReplaceRolesHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPut,
		Path:   "/api/v1/system/tenants/:tenantCodigo/groups/:code/roles",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplaceRoles(c.Request().Context(), tenantCodigo, c.Param("code"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantGroupPermissionsHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/groups/:code/permissions",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			response, err := usecase.FindPermissions(c.Request().Context(), tenantCodigo, c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}
