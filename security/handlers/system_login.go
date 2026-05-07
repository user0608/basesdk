package handlers

import (
	"basesdk/binds"
	"basesdk/httpapi"
	"basesdk/security/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/user0608/goones/answer"
	"github.com/user0608/goones/errs"
	"github.com/user0608/goones/kcheck"
)

func SystemUserHandler(usecase *usecases.SecurityUsecase) httpapi.Route {
	return &httpapi.PublicHandler{
		Method: http.MethodPost,
		Path:   "/api/v1/system/auth/login",
		// Path: "/api/v1/tenants/:tenantId/auth/login"
		Handler: func(c echo.Context) error {
			var ctx = c.Request().Context()
			var payload struct {
				Username string `json:"username" chk:"nonil"`
				Password string `json:"password" chk:"nonil"`
			}
			if err := binds.JSON(c, &payload); err != nil {
				return answer.Err(c, err)
			}
			if err := kcheck.Valid(payload); err != nil {
				return answer.Err(c, errs.BadRequestDirect(err.Error()))
			}
			token, err := usecase.SystemUserLogin(ctx, payload.Username, payload.Password)
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, echo.Map{"token": token})
		},
	}
}
