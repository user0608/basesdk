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

func SystemTenantsListHandler(usecase *usecases.TenantsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants",
		Handler: func(c echo.Context) error {
			response, err := usecase.List(c.Request().Context())
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemTenantCreateHandler(usecase *usecases.TenantsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPost,
		Path:   "/api/v1/system/tenants",
		Handler: func(c echo.Context) error {
			var payload dtos.CreateTenantInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := kcheck.Valid(payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Create(c.Request().Context(), payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Created(c)
		},
	}
}

func SystemTenantFindHandler(usecase *usecases.TenantsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/tenants/:tenantCodigo",
		Handler: func(c echo.Context) error {
			response, err := usecase.Find(c.Request().Context(), c.Param("tenantCodigo"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemTenantUpdateHandler(usecase *usecases.TenantsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPut,
		Path:   "/api/v1/system/tenants/:tenantCodigo",
		Handler: func(c echo.Context) error {
			var payload dtos.UpdateTenantInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := kcheck.Valid(payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Update(c.Request().Context(), c.Param("tenantCodigo"), payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantsEnableHandler(usecase *usecases.TenantsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPatch,
		Path:   "/api/v1/system/tenants/enable",
		Handler: func(c echo.Context) error {
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Enable(c.Request().Context(), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemTenantsDisableHandler(usecase *usecases.TenantsUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPatch,
		Path:   "/api/v1/system/tenants/disable",
		Handler: func(c echo.Context) error {
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Disable(c.Request().Context(), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}
