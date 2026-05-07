package security

import (
	"basesdk/security/jwt"
	"context"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/user0608/goones/answer"
	"github.com/user0608/goones/errs"
)

type securityContextKey string

const (
	contextUsernameKey securityContextKey = "security.username"
	contextTenantKey   securityContextKey = "security.tenant"
	contextTimeZoneKey securityContextKey = "security.time_zone"

	undefinedSecurityContextValue = "__security_context_undefined__"
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
		ctx = context.WithValue(ctx, contextUsernameKey, claims.Audience)
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
		ctx = context.WithValue(ctx, contextUsernameKey, claims.Audience)
		ctx = context.WithValue(ctx, contextTimeZoneKey, claims.TimeZone)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func Username(ctx context.Context) string {
	username, ok := ctx.Value(contextUsernameKey).(string)
	if !ok || username == "" {
		return undefinedSecurityContextValue
	}

	return username
}

func Tenant(ctx context.Context) string {
	tenant, ok := ctx.Value(contextTenantKey).(string)
	if !ok || tenant == "" {
		return undefinedSecurityContextValue
	}

	return tenant
}

func Tz(ctx context.Context) (*time.Location, error) {
	timeZone, ok := ctx.Value(contextTimeZoneKey).(string)
	if !ok || strings.TrimSpace(timeZone) == "" {
		return nil, errs.BadRequestDirect("falta la zona horaria")
	}

	location, err := time.LoadLocation(timeZone)
	if err != nil {
		return nil, errs.BadRequestDirect("zona horaria inválida")
	}

	return location, nil
}
