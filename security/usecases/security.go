package usecases

import (
	"basesdk/security/dtos"
	"basesdk/security/jwt"
	"basesdk/security/models"
	"basesdk/security/repositories"
	"context"

	"basesdk/errs"
)

type SecurityUsecase struct {
	tokenService         *jwt.TokenService
	systemUserRepository repositories.SystemUserRepository
}

func NewSecurityUsecase(
	tokenService *jwt.TokenService,
	systemUserRepository repositories.SystemUserRepository,
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

func (u *SecurityUsecase) FindSystemUser(ctx context.Context, username string) (*dtos.SystemUserDTO, error) {
	user, err := u.systemUserRepository.FindSystemUser(ctx, username)
	if err != nil {
		return nil, err
	}

	return dtos.NewSystemUserDTO(user), nil
}

func (u *SecurityUsecase) FindSystemUsers(ctx context.Context) ([]dtos.SystemUserDTO, error) {
	users, err := u.systemUserRepository.FindSystemUsers(ctx)
	if err != nil {
		return nil, err
	}

	return dtos.NewSystemUserDTOs(users), nil
}

func (u *SecurityUsecase) CreateSystemUser(ctx context.Context, username, password, createdBy string) (*dtos.SystemUserDTO, error) {
	if username == "" {
		return nil, errs.BadRequestDirect("el usuario es requerido")
	}

	exists, err := u.systemUserRepository.ExistsSystemUser(ctx, username)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errs.BadRequestDirect("el usuario ya existe")
	}

	user := &models.SystemUser{
		Username:  username,
		Disabled:  false,
		CreatedBy: createdBy,
	}

	if err := user.ChangePassword(password); err != nil {
		return nil, err
	}

	if err := u.systemUserRepository.CreateSystemUser(ctx, user); err != nil {
		return nil, err
	}

	return dtos.NewSystemUserDTO(user), nil
}

func (u *SecurityUsecase) ChangeSystemUserPassword(ctx context.Context, username, password string) error {
	user, err := u.systemUserRepository.FindSystemUser(ctx, username)
	if err != nil {
		return err
	}

	if err := user.ChangePassword(password); err != nil {
		return err
	}

	return u.systemUserRepository.ChangeSystemUserPassword(ctx, username, user.PasswordHash)
}

func (u *SecurityUsecase) EnableSystemUser(ctx context.Context, username string) error {
	return u.systemUserRepository.EnableSystemUser(ctx, username)
}

func (u *SecurityUsecase) DisableSystemUser(ctx context.Context, username string) error {
	return u.systemUserRepository.DisableSystemUser(ctx, username)
}

func (u *SecurityUsecase) DeleteSystemUser(ctx context.Context, username string) error {
	return u.systemUserRepository.DeleteSystemUser(ctx, username)
}
