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

func tenantGroupsRoute(path string, method string, system bool, usecase *usecases.TenantGroupsUsecase, action string) httpapi.Route {
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
			var payload dtos.CreateGroupInput
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
			var payload dtos.UpdateGroupInput
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
		case "users":
			response, err := usecase.FindUsers(ctx, tenantCodigo, c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		case "replace-users":
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplaceUsers(ctx, tenantCodigo, c.Param("code"), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		case "roles":
			response, err := usecase.FindRoles(ctx, tenantCodigo, c.Param("code"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		case "replace-roles":
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.ReplaceRoles(ctx, tenantCodigo, c.Param("code"), values, actor(c)); err != nil {
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

func TenantGroupsListHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/groups", http.MethodGet, false, usecase, "list")
}
func TenantGroupCreateHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/groups", http.MethodPost, false, usecase, "create")
}
func TenantGroupFindHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/groups/:code", http.MethodGet, false, usecase, "find")
}
func TenantGroupUpdateHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/groups/:code", http.MethodPut, false, usecase, "update")
}
func TenantGroupsEnableHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/groups/enable", http.MethodPatch, false, usecase, "enable")
}
func TenantGroupsDisableHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/groups/disable", http.MethodPatch, false, usecase, "disable")
}
func TenantGroupsDeleteHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/groups", http.MethodDelete, false, usecase, "delete")
}
func TenantGroupUsersHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/groups/:code/users", http.MethodGet, false, usecase, "users")
}
func TenantGroupReplaceUsersHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/groups/:code/users", http.MethodPut, false, usecase, "replace-users")
}
func TenantGroupRolesHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/groups/:code/roles", http.MethodGet, false, usecase, "roles")
}
func TenantGroupReplaceRolesHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/groups/:code/roles", http.MethodPut, false, usecase, "replace-roles")
}

func SystemTenantGroupsListHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/system/tenants/:tenantCodigo/groups", http.MethodGet, true, usecase, "list")
}
func SystemTenantGroupCreateHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/system/tenants/:tenantCodigo/groups", http.MethodPost, true, usecase, "create")
}
func SystemTenantGroupFindHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/system/tenants/:tenantCodigo/groups/:code", http.MethodGet, true, usecase, "find")
}
func SystemTenantGroupUpdateHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/system/tenants/:tenantCodigo/groups/:code", http.MethodPut, true, usecase, "update")
}
func SystemTenantGroupsEnableHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/system/tenants/:tenantCodigo/groups/enable", http.MethodPatch, true, usecase, "enable")
}
func SystemTenantGroupsDisableHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/system/tenants/:tenantCodigo/groups/disable", http.MethodPatch, true, usecase, "disable")
}
func SystemTenantGroupsDeleteHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/system/tenants/:tenantCodigo/groups", http.MethodDelete, true, usecase, "delete")
}
func SystemTenantGroupUsersHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/system/tenants/:tenantCodigo/groups/:code/users", http.MethodGet, true, usecase, "users")
}
func SystemTenantGroupReplaceUsersHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/system/tenants/:tenantCodigo/groups/:code/users", http.MethodPut, true, usecase, "replace-users")
}
func SystemTenantGroupRolesHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/system/tenants/:tenantCodigo/groups/:code/roles", http.MethodGet, true, usecase, "roles")
}
func SystemTenantGroupReplaceRolesHandler(usecase *usecases.TenantGroupsUsecase) httpapi.Route {
	return tenantGroupsRoute("/api/v1/system/tenants/:tenantCodigo/groups/:code/roles", http.MethodPut, true, usecase, "replace-roles")
}
