package binds

import (
	"basesdk/errs"

	"github.com/labstack/echo/v4"
)

func From(c echo.Context, payload any) error {
	return JSON(c, payload)
}

func JSON(c echo.Context, payload any) error {
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		return errs.BadRequestError(err, "json document invalido")
	}
	return nil
}

func Query(c echo.Context, payload any) error {
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, payload); err != nil {
		return errs.BadRequestError(err, "json document invalido")
	}
	return nil
}
