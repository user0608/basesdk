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

func TenantMeHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{Method: http.MethodGet, Path: "/api/v1/me", Handler: func(c echo.Context) error {
		response, err := usecase.Find(c.Request().Context(), tenant(c), actor(c))
		if err != nil {
			return answer.Err(c, err)
		}
		return answer.Ok(c, response)
	}}
}

func TenantMePasswordHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{Method: http.MethodPatch, Path: "/api/v1/me/password", Handler: func(c echo.Context) error {
		var payload dtos.ChangePasswordInput
		if err := binds.JSON(c, &payload); err != nil {
			return answer.Err(c, err)
		}
		if err := kcheck.Valid(payload); err != nil {
			return answer.Err(c, err)
		}
		if err := usecase.ChangePassword(c.Request().Context(), tenant(c), actor(c), payload, actor(c)); err != nil {
			return answer.Err(c, err)
		}
		return answer.Success(c)
	}}
}

func TenantMePermissionsHandler(usecase *usecases.TenantUsersUsecase) httpapi.Route {
	return &httpapi.TenantHandler{Method: http.MethodGet, Path: "/api/v1/me/permissions", Handler: func(c echo.Context) error {
		response, err := usecase.FindUserPermissions(c.Request().Context(), tenant(c), actor(c))
		if err != nil {
			return answer.Err(c, err)
		}
		return answer.Ok(c, response)
	}}
}
