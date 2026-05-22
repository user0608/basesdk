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

func TenantRolesListHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/roles",
		RequiredPerms: []string{securitypermissions.SecurityRolesRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.List(c.Request().Context(), tenant(c))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantRoleCreateHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPost,
		Path:          "/api/v1/roles",
		RequiredPerms: []string{securitypermissions.SecurityRolesCreate},
		Handler: func(c echo.Context) error {
			var payload dtos.CreateRoleInput
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

func TenantRoleFindHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/roles/:code",
		RequiredPerms: []string{securitypermissions.SecurityRolesRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.Find(c.Request().Context(), tenant(c), c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantRoleUpdateHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPut,
		Path:          "/api/v1/roles/:code",
		RequiredPerms: []string{securitypermissions.SecurityRolesUpdate},
		Handler: func(c echo.Context) error {
			var payload dtos.UpdateRoleInput
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

func TenantRolesEnableHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPatch,
		Path:          "/api/v1/roles/enable",
		RequiredPerms: []string{securitypermissions.SecurityRolesEnable},
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

func TenantRolesDisableHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPatch,
		Path:          "/api/v1/roles/disable",
		RequiredPerms: []string{securitypermissions.SecurityRolesDisable},
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

func TenantRolesDeleteHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodDelete,
		Path:          "/api/v1/roles",
		RequiredPerms: []string{securitypermissions.SecurityRolesDelete},
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

func TenantRolePermissionsHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/roles/:code/permissions",
		RequiredPerms: []string{securitypermissions.SecurityRolesPermissionsRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.FindPermissions(c.Request().Context(), tenant(c), c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantRoleReplacePermissionsHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPut,
		Path:          "/api/v1/roles/:code/permissions",
		RequiredPerms: []string{securitypermissions.SecurityRolesPermissionsReplace},
		Handler: func(c echo.Context) error {
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplacePermissions(c.Request().Context(), tenant(c), c.Param("code"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func TenantRoleUsersHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/roles/:code/users",
		RequiredPerms: []string{securitypermissions.SecurityRolesUsersRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.FindUsers(c.Request().Context(), tenant(c), c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantRoleReplaceUsersHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPut,
		Path:          "/api/v1/roles/:code/users",
		RequiredPerms: []string{securitypermissions.SecurityRolesUsersReplace},
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

func TenantRoleGroupsHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodGet,
		Path:          "/api/v1/roles/:code/groups",
		RequiredPerms: []string{securitypermissions.SecurityRolesGroupsRead},
		Handler: func(c echo.Context) error {
			response, err := usecase.FindGroups(c.Request().Context(), tenant(c), c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func TenantRoleReplaceGroupsHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.TenantHandler{
		Method:        http.MethodPut,
		Path:          "/api/v1/roles/:code/groups",
		RequiredPerms: []string{securitypermissions.SecurityRolesGroupsReplace},
		Handler: func(c echo.Context) error {
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplaceGroups(c.Request().Context(), tenant(c), c.Param("code"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantRolesListHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles",
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

func SystemTenantRoleCreateHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPost,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			var payload dtos.CreateRoleInput
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

func SystemTenantRoleFindHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles/:code",
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

func SystemTenantRoleUpdateHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPut,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles/:code",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			var payload dtos.UpdateRoleInput
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

func SystemTenantRolesEnableHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPatch,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles/enable",
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

func SystemTenantRolesDisableHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPatch,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles/disable",
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

func SystemTenantRolesDeleteHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodDelete,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles",
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

func SystemTenantRolePermissionsHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles/:code/permissions",
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

func SystemTenantRoleReplacePermissionsHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPut,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles/:code/permissions",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplacePermissions(c.Request().Context(), tenantCodigo, c.Param("code"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantRoleUsersHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles/:code/users",
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

func SystemTenantRoleReplaceUsersHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPut,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles/:code/users",
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

func SystemTenantRoleGroupsHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles/:code/groups",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			response, err := usecase.FindGroups(c.Request().Context(), tenantCodigo, c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemTenantRoleReplaceGroupsHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPut,
		Path:   "/api/v1/system/tenants/:tenantCodigo/roles/:code/groups",
		Handler: func(c echo.Context) error {
			tenantCodigo, err := tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplaceGroups(c.Request().Context(), tenantCodigo, c.Param("code"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}
