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

func tenantUsersListRoute(path string, system bool, usecase *usecases.TenantUsersUsecase) httpapi.Route {
	handler := func(c echo.Context) error {
		tenantCodigo := tenant(c)
		if system {
			var err error
			tenantCodigo, err = tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
		}
		response, err := usecase.List(c.Request().Context(), tenantCodigo)
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

func tenantUserCreateRoute(path string, system bool, usecase *usecases.TenantUsersUsecase) httpapi.Route {
	handler := func(c echo.Context) error {
		tenantCodigo := tenant(c)
		if system {
			var err error
			tenantCodigo, err = tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
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
	}
	if system {
		return &httpapi.SystemHandler{Method: http.MethodPost, Path: path, Handler: handler}
	}
	return &httpapi.TenantHandler{Method: http.MethodPost, Path: path, Handler: handler}
}

func tenantUserFindRoute(path string, system bool, usecase *usecases.TenantUsersUsecase) httpapi.Route {
	handler := func(c echo.Context) error {
		tenantCodigo := tenant(c)
		if system {
			var err error
			tenantCodigo, err = tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
		}
		response, err := usecase.Find(c.Request().Context(), tenantCodigo, c.Param("username"))
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

func tenantUserUpdateRoute(path string, system bool, usecase *usecases.TenantUsersUsecase) httpapi.Route {
	handler := func(c echo.Context) error {
		tenantCodigo := tenant(c)
		if system {
			var err error
			tenantCodigo, err = tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
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
	}
	if system {
		return &httpapi.SystemHandler{Method: http.MethodPut, Path: path, Handler: handler}
	}
	return &httpapi.TenantHandler{Method: http.MethodPut, Path: path, Handler: handler}
}

func tenantUserPasswordRoute(path string, system bool, usecase *usecases.TenantUsersUsecase) httpapi.Route {
	handler := func(c echo.Context) error {
		tenantCodigo := tenant(c)
		if system {
			var err error
			tenantCodigo, err = tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
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
	}
	if system {
		return &httpapi.SystemHandler{Method: http.MethodPatch, Path: path, Handler: handler}
	}
	return &httpapi.TenantHandler{Method: http.MethodPatch, Path: path, Handler: handler}
}

func tenantUsersBulkRoute(path string, method string, system bool, usecase *usecases.TenantUsersUsecase, action string) httpapi.Route {
	handler := func(c echo.Context) error {
		tenantCodigo := tenant(c)
		if system {
			var err error
			tenantCodigo, err = tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
		}
		values, err := requestStrings(c)
		if err != nil {
			return answer.Err(c, err)
		}
		switch action {
		case "enable":
			err = usecase.Enable(c.Request().Context(), tenantCodigo, values, actor(c))
		case "disable":
			err = usecase.Disable(c.Request().Context(), tenantCodigo, values, actor(c))
		case "delete":
			err = usecase.Delete(c.Request().Context(), tenantCodigo, values)
		}
		if err != nil {
			return answer.Err(c, err)
		}
		return answer.Success(c)
	}
	if system {
		return &httpapi.SystemHandler{Method: method, Path: path, Handler: handler}
	}
	return &httpapi.TenantHandler{Method: method, Path: path, Handler: handler}
}

func tenantUserPermissionsRoute(path string, system bool, usecase *usecases.TenantUsersUsecase) httpapi.Route {
	handler := func(c echo.Context) error {
		tenantCodigo := tenant(c)
		if system {
			var err error
			tenantCodigo, err = tenantParam(c)
			if err != nil {
				return answer.Err(c, err)
			}
		}
		response, err := usecase.FindUserPermissions(c.Request().Context(), tenantCodigo, c.Param("username"))
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

func TenantUsersListHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUsersListRoute("/api/v1/users", false, usecase)
}
func TenantUserCreateHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUserCreateRoute("/api/v1/users", false, usecase)
}
func TenantUserFindHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUserFindRoute("/api/v1/users/:username", false, usecase)
}
func TenantUserUpdateHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUserUpdateRoute("/api/v1/users/:username", false, usecase)
}
func TenantUserPasswordHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUserPasswordRoute("/api/v1/users/:username/password", false, usecase)
}
func TenantUsersEnableHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUsersBulkRoute("/api/v1/users/enable", http.MethodPatch, false, usecase, "enable")
}
func TenantUsersDisableHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUsersBulkRoute("/api/v1/users/disable", http.MethodPatch, false, usecase, "disable")
}
func TenantUsersDeleteHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUsersBulkRoute("/api/v1/users", http.MethodDelete, false, usecase, "delete")
}
func TenantUserPermissionsHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUserPermissionsRoute("/api/v1/users/:username/permissions", false, usecase)
}

func SystemTenantUsersListHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUsersListRoute("/api/v1/system/tenants/:tenantCodigo/users", true, usecase)
}
func SystemTenantUserCreateHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUserCreateRoute("/api/v1/system/tenants/:tenantCodigo/users", true, usecase)
}
func SystemTenantUserFindHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUserFindRoute("/api/v1/system/tenants/:tenantCodigo/users/:username", true, usecase)
}
func SystemTenantUserUpdateHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUserUpdateRoute("/api/v1/system/tenants/:tenantCodigo/users/:username", true, usecase)
}
func SystemTenantUserPasswordHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUserPasswordRoute("/api/v1/system/tenants/:tenantCodigo/users/:username/password", true, usecase)
}
func SystemTenantUsersEnableHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUsersBulkRoute("/api/v1/system/tenants/:tenantCodigo/users/enable", http.MethodPatch, true, usecase, "enable")
}
func SystemTenantUsersDisableHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUsersBulkRoute("/api/v1/system/tenants/:tenantCodigo/users/disable", http.MethodPatch, true, usecase, "disable")
}
func SystemTenantUsersDeleteHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUsersBulkRoute("/api/v1/system/tenants/:tenantCodigo/users", http.MethodDelete, true, usecase, "delete")
}
func SystemTenantUserPermissionsHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return tenantUserPermissionsRoute("/api/v1/system/tenants/:tenantCodigo/users/:username/permissions", true, usecase)
}
