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

func (u *PermissionsUsecase) FindRoles(ctx context.Context, code string) ([]dtos.RoleResponse, error) {
	roles, err := u.repository.FindPermissionRoles(ctx, strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.RoleResponse, 0, len(roles))
	for _, role := range roles {
		response = append(response, toRoleResponse(role))
	}
	return response, nil
}

func (u *PermissionsUsecase) FindGroups(ctx context.Context, code string) ([]dtos.GroupResponse, error) {
	groups, err := u.repository.FindPermissionGroups(ctx, strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.GroupResponse, 0, len(groups))
	for _, group := range groups {
		response = append(response, toGroupResponse(group))
	}
	return response, nil
}

func (u *PermissionsUsecase) FindUsers(ctx context.Context, code string) ([]dtos.TenantUserResponse, error) {
	users, err := u.repository.FindPermissionUsers(ctx, strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.TenantUserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, toTenantUserResponse(user))
	}
	return response, nil
}
