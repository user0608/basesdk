package handlers

import (
	"basesdk/binds"
	"basesdk/errs"
	"basesdk/httpapi"
	"basesdk/security/usecases"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/user0608/goones/answer"
	"github.com/user0608/goones/kcheck"
)

func TenantUserHandler(usecase *usecases.SecurityUsecase) httpapi.Route {
	return &httpapi.PublicHandler{
		Method: http.MethodPost,
		Path:   "/api/v1/tenants/:tenantCodigo/auth/login",
		Handler: func(c echo.Context) error {
			ctx := c.Request().Context()
			tenantCodigo := strings.TrimSpace(c.Param("tenantCodigo"))
			if tenantCodigo == "" {
				return answer.Err(c, errs.BadRequestDirect("tenant requerido"))
			}

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

			token, err := usecase.TenantUserLogin(ctx, tenantCodigo, payload.Username, payload.Password)
			if err != nil {
				return answer.Err(c, err)
			}

			return answer.Ok(c, echo.Map{"token": token})
		},
	}
}
