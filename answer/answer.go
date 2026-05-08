package answer

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"basesdk/errs"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Ok(c echo.Context, payload any) error {
	return c.JSON(http.StatusOK, &Response{
		Data: payload,
	})
}

func Created(c echo.Context) error {
	return c.JSON(http.StatusCreated, &Response{
		Message: "Recurso creado exitosamente",
	})
}

func Message(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, &Response{
		Message: message,
	})
}

func Success(c echo.Context) error {
	return c.JSON(http.StatusOK, &Response{
		Message: "Operación completada exitosamente",
	})
}

func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func UnwrapErr(err error) (code int, message string) {
	var werr *errs.Err

	code = http.StatusInternalServerError
	message = "Ocurrió un problema. Se produjo un error inesperado."

	if errors.As(err, &werr) {
		code = werr.Code()
		message = werr.Message()
	}

	if werr == nil && err != nil {
		errSMS := strings.TrimSpace(err.Error())

		if strings.HasPrefix(errSMS, ":") {
			code = http.StatusBadRequest
			message = strings.TrimLeft(errSMS, ":")
			return
		}

		slog.Error("internal error", "error", err)
		return
	}

	if werr != nil && werr.Wrapped() != nil {
		slog.Error("internal error", "error", werr.Wrapped())
	}

	return code, message
}

func Err(c echo.Context, err error) error {
	code, message := UnwrapErr(err)

	return c.JSON(code, &Response{
		Message: message,
	})
}

func Auto(c echo.Context, err error) error {
	if err != nil {
		return Err(c, err)
	}

	return Success(c)
}
