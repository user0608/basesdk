package handlers

import (
	"basesdk/binds"
	"basesdk/httpapi"
	"basesdk/security/dtos"
	"basesdk/security/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/user0608/goones/answer"
	"github.com/user0608/goones/kcheck"
)

func tenantRolesRoute(path string, method string, system bool, usecase *usecases.TenantRolesUsecase, action string) httpapi.Route {
	handler := func(c echo.Context) error {
		tenantCodigo := tenant(c)
		if system {
			var err error
			tenantCodigo, err = tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
		}
		ctx := c.Request().Context()
		switch action {
		case "list":
			response, err := usecase.List(ctx, tenantCodigo)
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		case "find":
			response, err := usecase.Find(ctx, tenantCodigo, c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		case "create":
			var payload dtos.CreateRoleInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := kcheck.Valid(payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Create(ctx, tenantCodigo, payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Created(c)
		case "update":
			var payload dtos.UpdateRoleInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Update(ctx, tenantCodigo, c.Param("code"), payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		case "enable", "disable", "delete":
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if action == "enable" {
				err = usecase.Enable(ctx, tenantCodigo, values, actor(c))
			}
			if action == "disable" {
				err = usecase.Disable(ctx, tenantCodigo, values, actor(c))
			}
			if action == "delete" {
				err = usecase.Delete(ctx, tenantCodigo, values)
			}
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		case "permissions":
			response, err := usecase.FindPermissions(ctx, tenantCodigo, c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		case "replace-permissions":
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplacePermissions(ctx, tenantCodigo, c.Param("code"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		}
		return answer.Success(c)
	}
	if system {
		return &httpapi.SystemHandler{Method: method, Path: path, Handler: handler}
	}
	return &httpapi.TenantHandler{Method: method, Path: path, Handler: handler}
}

func TenantRolesListHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/roles", http.MethodGet, false, usecase, "list")
}
func TenantRoleCreateHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/roles", http.MethodPost, false, usecase, "create")
}
func TenantRoleFindHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/roles/:code", http.MethodGet, false, usecase, "find")
}
func TenantRoleUpdateHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/roles/:code", http.MethodPut, false, usecase, "update")
}
func TenantRolesEnableHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/roles/enable", http.MethodPatch, false, usecase, "enable")
}
func TenantRolesDisableHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/roles/disable", http.MethodPatch, false, usecase, "disable")
}
func TenantRolesDeleteHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/roles", http.MethodDelete, false, usecase, "delete")
}
func TenantRolePermissionsHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/roles/:code/permissions", http.MethodGet, false, usecase, "permissions")
}
func TenantRoleReplacePermissionsHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/roles/:code/permissions", http.MethodPut, false, usecase, "replace-permissions")
}

func SystemTenantRolesListHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/system/tenants/:tenantCodigo/roles", http.MethodGet, true, usecase, "list")
}
func SystemTenantRoleCreateHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/system/tenants/:tenantCodigo/roles", http.MethodPost, true, usecase, "create")
}
func SystemTenantRoleFindHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/system/tenants/:tenantCodigo/roles/:code", http.MethodGet, true, usecase, "find")
}
func SystemTenantRoleUpdateHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/system/tenants/:tenantCodigo/roles/:code", http.MethodPut, true, usecase, "update")
}
func SystemTenantRolesEnableHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/system/tenants/:tenantCodigo/roles/enable", http.MethodPatch, true, usecase, "enable")
}
func SystemTenantRolesDisableHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/system/tenants/:tenantCodigo/roles/disable", http.MethodPatch, true, usecase, "disable")
}
func SystemTenantRolesDeleteHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/system/tenants/:tenantCodigo/roles", http.MethodDelete, true, usecase, "delete")
}
func SystemTenantRolePermissionsHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/system/tenants/:tenantCodigo/roles/:code/permissions", http.MethodGet, true, usecase, "permissions")
}
func SystemTenantRoleReplacePermissionsHandler(usecase *usecases.TenantRolesUsecase) httpapi.Route {
	return tenantRolesRoute("/api/v1/system/tenants/:tenantCodigo/roles/:code/permissions", http.MethodPut, true, usecase, "replace-permissions")
}
