// Package answer provides standardized HTTP JSON responses for handlers.
//
// It centralizes success and error response formatting, so the application can
// keep a consistent API contract. It also maps domain errors into HTTP status
// codes and logs unexpected internal failures.
package answer

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"basesdk/errs"
)

type Target interface {
	JSON(code int, i any) error
	NoContent(code int) error
}

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Ok(c Target, payload any) error {
	return c.JSON(http.StatusOK, &Response{
		Data: payload,
	})
}

func Created(c Target, payload any) error {
	return c.JSON(http.StatusCreated, &Response{
		Data: payload,
	})
}

func Accepted(c Target, payload any) error {
	return c.JSON(http.StatusAccepted, &Response{
		Data: payload,
	})
}

func Message(c Target, message string) error {
	return c.JSON(http.StatusOK, &Response{
		Message: message,
	})
}

func CreatedMessage(c Target, message string) error {
	return c.JSON(http.StatusCreated, &Response{
		Message: message,
	})
}

func AcceptedMessage(c Target, message string) error {
	return c.JSON(http.StatusAccepted, &Response{
		Message: message,
	})
}

func Success(c Target) error {
	return c.JSON(http.StatusOK, &Response{
		Message: "Operación completada exitosamente",
	})
}

func CreatedSuccess(c Target) error {
	return c.JSON(http.StatusCreated, &Response{
		Message: "Recurso creado exitosamente",
	})
}

func AcceptedSuccess(c Target) error {
	return c.JSON(http.StatusAccepted, &Response{
		Message: "Operación aceptada exitosamente",
	})
}

func NoContent(c Target) error {
	return c.NoContent(http.StatusNoContent)
}

// UnwrapErr converts application errors into standardized HTTP responses.
//
// Supported behaviors:
//   - errs.Err: extracts HTTP status code and public message.
//   - errors prefixed with ":" are treated as client errors (400).
//   - unexpected errors are logged as internal failures.
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

func Err(c Target, err error) error {
	code, message := UnwrapErr(err)

	return c.JSON(code, &Response{
		Message: message,
	})
}

func Auto(c Target, err error) error {
	if err != nil {
		return Err(c, err)
	}

	return Success(c)
}
