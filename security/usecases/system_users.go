package usecases

import (
	"basesdk/errs"
	"basesdk/security/dtos"
	"basesdk/security/models"
	"basesdk/security/repositories"
	"context"
	"slices"
	"strings"
	"time"
)

type SystemUsersUsecase struct {
	repository *repositories.SystemUserRepository
}

func NewSystemUsersUsecase(repository *repositories.SystemUserRepository) *SystemUsersUsecase {
	return &SystemUsersUsecase{repository: repository}
}

func (u *SystemUsersUsecase) List(ctx context.Context) ([]dtos.SystemUserResponse, error) {
	users, err := u.repository.FindSystemUsers(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dtos.SystemUserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, toSystemUserResponse(user))
	}

	return response, nil
}

func (u *SystemUsersUsecase) Find(ctx context.Context, username string) (*dtos.SystemUserResponse, error) {
	user, err := u.repository.FindSystemUser(ctx, strings.TrimSpace(username))
	if err != nil {
		return nil, err
	}

	response := toSystemUserResponse(*user)
	return &response, nil
}

func (u *SystemUsersUsecase) Create(ctx context.Context, input dtos.CreateSystemUserInput, createdBy string) error {
	user := &models.SystemAccount{
		Username:  strings.TrimSpace(input.Username),
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}
	if user.Username == "" {
		return errs.BadRequestDirect("usuario requerido")
	}
	if err := user.ChangePassword(input.Password); err != nil {
		return err
	}

	return u.repository.CreateSystemUser(ctx, user)
}

func (u *SystemUsersUsecase) Update(ctx context.Context, username string, input dtos.UpdateSystemUserInput, updatedBy string) error {
	username = strings.TrimSpace(username)
	if input.Disabled {
		if username == strings.TrimSpace(updatedBy) {
			return errs.BadRequestDirect("no puedes deshabilitar tu propio usuario system")
		}
		if err := u.ensureActiveSystemUserRemains(ctx, []string{username}); err != nil {
			return err
		}
	}

	now := time.Now()
	return u.repository.UpdateSystemUser(ctx, &models.SystemAccount{
		Username:  username,
		Disabled:  input.Disabled,
		UpdatedBy: &updatedBy,
		UpdatedAt: &now,
	})
}

func (u *SystemUsersUsecase) Enable(ctx context.Context, usernames []string, updatedBy string) error {
	if len(usernames) == 0 {
		return errs.BadRequestDirect("usuarios requeridos")
	}
	return u.repository.SetSystemUsersDisabled(ctx, usernames, false, updatedBy)
}

func (u *SystemUsersUsecase) Disable(ctx context.Context, usernames []string, updatedBy string) error {
	if len(usernames) == 0 {
		return errs.BadRequestDirect("usuarios requeridos")
	}
	if slices.Contains(usernames, strings.TrimSpace(updatedBy)) {
		return errs.BadRequestDirect("no puedes deshabilitar tu propio usuario system")
	}
	if err := u.ensureActiveSystemUserRemains(ctx, usernames); err != nil {
		return err
	}
	return u.repository.SetSystemUsersDisabled(ctx, usernames, true, updatedBy)
}

func (u *SystemUsersUsecase) Delete(ctx context.Context, usernames []string, deletedBy string) error {
	if len(usernames) == 0 {
		return errs.BadRequestDirect("usuarios requeridos")
	}
	if slices.Contains(usernames, strings.TrimSpace(deletedBy)) {
		return errs.BadRequestDirect("no puedes eliminar tu propio usuario system")
	}
	if err := u.ensureActiveSystemUserRemains(ctx, usernames); err != nil {
		return err
	}
	return u.repository.DeleteSystemUsers(ctx, usernames)
}

func (u *SystemUsersUsecase) ensureActiveSystemUserRemains(ctx context.Context, excludedUsernames []string) error {
	count, err := u.repository.CountActiveSystemUsersExcluding(ctx, excludedUsernames)
	if err != nil {
		return err
	}
	if count == 0 {
		return errs.BadRequestDirect("debe quedar al menos un usuario system activo")
	}
	return nil
}
