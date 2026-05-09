package usecases

import (
	"basesdk/auth/jwt"
	"basesdk/security/repositories"
	"context"
	"time"

	"basesdk/errs"
)

type SecurityUsecase struct {
	tokenService         *jwt.TokenService
	systemUserRepository *repositories.SystemUserRepository
	appUserRepository    *repositories.AppUserRepository
}

func NewSecurityUsecase(
	tokenService *jwt.TokenService,
	systemUserRepository *repositories.SystemUserRepository,
	appUserRepository *repositories.AppUserRepository,
) *SecurityUsecase {
	return &SecurityUsecase{
		tokenService:         tokenService,
		systemUserRepository: systemUserRepository,
		appUserRepository:    appUserRepository,
	}
}

func (u *SecurityUsecase) SystemUserLogin(ctx context.Context, username, password string) (string, error) {
	user, err := u.systemUserRepository.FindSystemUser(ctx, username)
	if err != nil {
		return "", err
	}

	if user.Disabled {
		return "", errs.BadRequestDirect("usuario deshabilitado")
	}

	if err := user.ValidatePassword(password); err != nil {
		return "", err
	}

	return u.tokenService.GenerateSystemToken(ctx, username)
}

func (u *SecurityUsecase) TenantUserLogin(ctx context.Context, tenantCodigo, username, password string) (string, error) {
	tenant, err := u.appUserRepository.FindTenant(ctx, tenantCodigo)
	if err != nil {
		return "", err
	}

	if err := tenant.ValidateLoginAccess(time.Now()); err != nil {
		return "", err
	}

	user, err := u.appUserRepository.FindAppUser(ctx, tenantCodigo, username)
	if err != nil {
		return "", err
	}

	if err := user.ValidatePassword(password); err != nil {
		return "", err
	}

	return u.tokenService.GenerateTenantToken(ctx, tenant.Codigo, username, tenant.Timezone)
}
