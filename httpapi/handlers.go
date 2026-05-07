package httpapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TenantHandler struct {
	Method            string
	Path              string
	BeforeMiddlewares []echo.MiddlewareFunc
	Middlewares       []echo.MiddlewareFunc
	Handler           echo.HandlerFunc
}

var _ Route = (*TenantHandler)(nil)
var _ BeforeSecurityMiddlewareProvider = (*TenantHandler)(nil)
var _ AfterSecurityMiddlewareProvider = (*TenantHandler)(nil)

// GetMethod implements [Route].
func (h *TenantHandler) GetMethod() string {
	if h.Method == "" {
		return http.MethodGet
	}
	return h.Method
}

// GetPath implements [Route].
func (h *TenantHandler) GetPath() string {
	return h.Path
}

// HandleRequest implements [Route].
func (h *TenantHandler) HandleRequest(c echo.Context) error {
	if h.Handler == nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"route handler is not configured",
		)
	}
	return h.Handler(c)
}

// BeforeSecurityMiddlewares implements [BeforeSecurityMiddlewareProvider].
func (h *TenantHandler) BeforeSecurityMiddlewares() []echo.MiddlewareFunc {
	return h.BeforeMiddlewares
}

// AfterSecurityMiddlewares implements [AfterSecurityMiddlewareProvider].
func (h *TenantHandler) AfterSecurityMiddlewares() []echo.MiddlewareFunc {
	return h.Middlewares
}

// System
type SystemHandler struct {
	SystemRoute
	Method            string
	Path              string
	BeforeMiddlewares []echo.MiddlewareFunc
	Middlewares       []echo.MiddlewareFunc
	Handler           echo.HandlerFunc
}

var _ Route = (*SystemHandler)(nil)
var _ BeforeSecurityMiddlewareProvider = (*SystemHandler)(nil)
var _ AfterSecurityMiddlewareProvider = (*SystemHandler)(nil)

// GetMethod implements [Route].
func (h *SystemHandler) GetMethod() string {
	if h.Method == "" {
		return http.MethodGet
	}
	return h.Method
}

// GetPath implements [Route].
func (h *SystemHandler) GetPath() string {
	return h.Path
}

// HandleRequest implements [Route].
func (h *SystemHandler) HandleRequest(c echo.Context) error {
	if h.Handler == nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"route handler is not configured",
		)
	}
	return h.Handler(c)
}

// BeforeSecurityMiddlewares implements [BeforeSecurityMiddlewareProvider].
func (h *SystemHandler) BeforeSecurityMiddlewares() []echo.MiddlewareFunc {
	return h.BeforeMiddlewares
}

// AfterSecurityMiddlewares implements [AfterSecurityMiddlewareProvider].
func (h *SystemHandler) AfterSecurityMiddlewares() []echo.MiddlewareFunc {
	return h.Middlewares
}

// Public
type PublicHandler struct {
	PublicRoute
	Method            string
	Path              string
	BeforeMiddlewares []echo.MiddlewareFunc
	Middlewares       []echo.MiddlewareFunc
	Handler           echo.HandlerFunc
}

var _ Route = (*PublicHandler)(nil)
var _ BeforeSecurityMiddlewareProvider = (*PublicHandler)(nil)
var _ AfterSecurityMiddlewareProvider = (*PublicHandler)(nil)

// GetMethod implements [Route].
func (h *PublicHandler) GetMethod() string {
	if h.Method == "" {
		return http.MethodGet
	}
	return h.Method
}

// GetPath implements [Route].
func (h *PublicHandler) GetPath() string {
	return h.Path
}

// HandleRequest implements [Route].
func (h *PublicHandler) HandleRequest(c echo.Context) error {
	if h.Handler == nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"route handler is not configured",
		)
	}
	return h.Handler(c)
}

// BeforeSecurityMiddlewares implements [BeforeSecurityMiddlewareProvider].
func (h *PublicHandler) BeforeSecurityMiddlewares() []echo.MiddlewareFunc {
	return h.BeforeMiddlewares
}

// AfterSecurityMiddlewares implements [AfterSecurityMiddlewareProvider].
func (h *PublicHandler) AfterSecurityMiddlewares() []echo.MiddlewareFunc {
	return h.Middlewares
}
