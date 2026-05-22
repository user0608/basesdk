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

func SystemUsersListHandler(usecase *usecases.SystemUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/users",
		Handler: func(c echo.Context) error {
			response, err := usecase.List(c.Request().Context())
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemUserCreateHandler(usecase *usecases.SystemUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPost,
		Path:   "/api/v1/system/users",
		Handler: func(c echo.Context) error {
			var payload dtos.CreateSystemUserInput
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

func SystemUserFindHandler(usecase *usecases.SystemUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodGet,
		Path:   "/api/v1/system/users/:username",
		Handler: func(c echo.Context) error {
			response, err := usecase.Find(c.Request().Context(), c.Param("username"))
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, response)
		},
	}
}

func SystemUserUpdateHandler(usecase *usecases.SystemUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPut,
		Path:   "/api/v1/system/users/:username",
		Handler: func(c echo.Context) error {
			var payload dtos.UpdateSystemUserInput
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Update(c.Request().Context(), c.Param("username"), payload, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}

func SystemUsersEnableHandler(usecase *usecases.SystemUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPatch,
		Path:   "/api/v1/system/users/enable",
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

func SystemUsersDisableHandler(usecase *usecases.SystemUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodPatch,
		Path:   "/api/v1/system/users/disable",
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

func SystemUsersDeleteHandler(usecase *usecases.SystemUsersUsecase) httpapi.Route {
	return &httpapi.SystemHandler{
		Method: http.MethodDelete,
		Path:   "/api/v1/system/users",
		Handler: func(c echo.Context) error {
			values, err := requestStrings(c)
			if err != nil {
				return answer.Err(c, err)
			}
			if err := usecase.Delete(c.Request().Context(), values, actor(c)); err != nil {
				return answer.Err(c, err)
			}
			return answer.Success(c)
		},
	}
}
