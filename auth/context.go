package auth

import (
	"context"
	"strings"
	"time"

	"basesdk/errs"
)

type securityContextKey string

const (
	contextUsernameKey securityContextKey = "security.username"
	contextTenantKey   securityContextKey = "security.tenant"
	contextTimeZoneKey securityContextKey = "security.time_zone"

	undefinedSecurityContextValue = "__security_context_undefined__"
)

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

func IsUndefined(value string) bool {
	return value == undefinedSecurityContextValue
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
