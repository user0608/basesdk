package httpapi

import (
	"basesdk/auth"
	"basesdk/configs"
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

func buildMiddlewares(route Route, securityMiddleware auth.SecurityMiddleware) []echo.MiddlewareFunc {
	var before []echo.MiddlewareFunc
	var after []echo.MiddlewareFunc

	if r, ok := route.(BeforeSecurityMiddlewareProvider); ok {
		before = r.BeforeSecurityMiddlewares()
	}

	if r, ok := route.(AfterSecurityMiddlewareProvider); ok {
		after = r.AfterSecurityMiddlewares()
	}

	_, isPublic := route.(publicRouteMarker)
	_, isSystem := route.(systemRouteMarker)

	var securityMdw []echo.MiddlewareFunc

	if isPublic && isSystem {
		slog.Warn(
			"conflicting route markers: public + system, defaulting to system",
			"path", route.GetPath(),
			"method", route.GetMethod(),
		)
		securityMdw = []echo.MiddlewareFunc{securityMiddleware.System}
	} else if isSystem {
		securityMdw = []echo.MiddlewareFunc{securityMiddleware.System}
	} else if isPublic {
		securityMdw = nil
	} else {
		securityMdw = []echo.MiddlewareFunc{securityMiddleware.Tenant}
	}

	middlewares := make([]echo.MiddlewareFunc, 0, len(before)+len(securityMdw)+len(after))
	middlewares = append(middlewares, before...)
	middlewares = append(middlewares, securityMdw...)
	middlewares = append(middlewares, after...)

	return middlewares
}

func NewServer(routes []Route, securityMiddleware auth.SecurityMiddleware) *echo.Echo {
	server := echo.New()
	server.HideBanner = true
	server.Use(middleware.RequestLogger())
	server.Use(middleware.Recover())

	for _, route := range routes {
		routeMiddlewares := buildMiddlewares(route, securityMiddleware)

		switch route.GetMethod() {
		case http.MethodGet:
			server.GET(route.GetPath(), route.HandleRequest, routeMiddlewares...)
		case http.MethodPost:
			server.POST(route.GetPath(), route.HandleRequest, routeMiddlewares...)
		case http.MethodPut:
			server.PUT(route.GetPath(), route.HandleRequest, routeMiddlewares...)
		case http.MethodDelete:
			server.DELETE(route.GetPath(), route.HandleRequest, routeMiddlewares...)
		case http.MethodPatch:
			server.PATCH(route.GetPath(), route.HandleRequest, routeMiddlewares...)
		default:
			slog.Warn("unsupported route method", "method", route.GetMethod(), "path", route.GetPath())
		}
	}

	return server
}

var Module = fx.Module("http-server",
	fx.Provide(
		auth.NewSecurityMiddleware,
		fx.Annotate(
			NewServer,
			fx.ParamTags(
				RouteTag,
			),
		),
	),
)

func StartWebServer(lc fx.Lifecycle, e *echo.Echo, c configs.ApplicationConfigs) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				e.Use(middleware.CORS())
				slog.Info("Starting HTTP server", "address", c.ListenAddress())
				if err := e.Start(c.ListenAddress()); err != nil && err != http.ErrServerClosed {
					slog.Error("HTTP server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Shutting down HTTP server")
			if err := e.Shutdown(ctx); err != nil {
				slog.Error("Error shutting down HTTP server", "error", err)
				return err
			}
			slog.Info("HTTP server stopped successfully")
			return nil
		},
	})
}
