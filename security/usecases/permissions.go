package usecases

import (
	"basesdk/security/dtos"
	"basesdk/security/repositories"
	"context"
	"strings"
)

type PermissionsUsecase struct {
	repository *repositories.PermissionRepository
}

func NewPermissionsUsecase(repository *repositories.PermissionRepository) *PermissionsUsecase {
	return &PermissionsUsecase{repository: repository}
}

func (u *PermissionsUsecase) List(ctx context.Context) ([]dtos.PermissionResponse, error) {
	permissions, err := u.repository.FindPermissions(ctx)
	if err != nil {
		return nil, err
	}
	response := make([]dtos.PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		response = append(response, toPermissionResponse(permission))
	}
	return response, nil
}

func (u *PermissionsUsecase) Find(ctx context.Context, code string) (*dtos.PermissionResponse, error) {
	permission, err := u.repository.FindPermission(ctx, strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	response := toPermissionResponse(*permission)
	return &response, nil
}
