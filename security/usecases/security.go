package usecases

import (
	"basesdk/security/jwt"
	"basesdk/security/repositories"
	"context"

	"basesdk/errs"
)

type SecurityUsecase struct {
	tokenService         *jwt.TokenService
	systemUserRepository *repositories.SystemUserRepository
}

func NewSecurityUsecase(
	tokenService *jwt.TokenService,
	systemUserRepository *repositories.SystemUserRepository,
) *SecurityUsecase {
	return &SecurityUsecase{
		tokenService:         tokenService,
		systemUserRepository: systemUserRepository,
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
