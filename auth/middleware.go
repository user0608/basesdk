package auth

import (
	"basesdk/auth/jwt"
	"context"
	"strings"

	"basesdk/errs"

	"github.com/labstack/echo/v4"
	"github.com/user0608/goones/answer"
)

type SecurityMiddleware struct {
	tokenService *jwt.TokenService
}

func NewSecurityMiddleware(tokenService *jwt.TokenService) *SecurityMiddleware {
	return &SecurityMiddleware{
		tokenService: tokenService,
	}
}

func (s *SecurityMiddleware) readJWTToken(c echo.Context) (string, error) {
	authorization := strings.TrimSpace(c.Request().Header.Get("Authorization"))
	if authorization == "" {
		return "", errs.UnauthorizedDirect("no se encontró el token de autorización en la petición")
	}

	parts := strings.Fields(authorization)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errs.UnauthorizedDirect("el token de autorización debe enviarse con el formato Bearer")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errs.UnauthorizedDirect("el token de autorización está vacío")
	}

	return token, nil
}

func (s *SecurityMiddleware) Tenant(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := s.readJWTToken(c)
		if err != nil {
			return answer.Err(c, err)
		}

		claims, err := s.tokenService.ValidateToken(c.Request().Context(), token)
		if err != nil {
			return answer.Err(c, err)
		}

		if claims.Type != jwt.TypeTenant {
			return answer.Err(c, errs.UnauthorizedDirect("el token no corresponde a una sesión de tenant"))
		}

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, contextUsernameKey, claims.Subject)
		ctx = context.WithValue(ctx, contextTenantKey, claims.Tenant)
		ctx = context.WithValue(ctx, contextTimeZoneKey, claims.TimeZone)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func (s *SecurityMiddleware) System(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := s.readJWTToken(c)
		if err != nil {
			return answer.Err(c, err)
		}

		claims, err := s.tokenService.ValidateToken(c.Request().Context(), token)
		if err != nil {
			return answer.Err(c, err)
		}

		if claims.Type != jwt.TypeSystem {
			return answer.Err(c, errs.UnauthorizedDirect("el token no corresponde a una sesión de sistema"))
		}

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, contextUsernameKey, claims.Subject)
		ctx = context.WithValue(ctx, contextTimeZoneKey, claims.TimeZone)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
