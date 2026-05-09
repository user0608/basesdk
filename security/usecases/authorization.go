package usecases

import (
	"basesdk/security/models"
	"basesdk/security/repositories"
	"context"
	"strings"
)

const adminPermission = "admin"

type AuthorizationUsecase struct {
	permissionRepository *repositories.PermissionRepository
}

func NewAuthorizationUsecase(permissionRepository *repositories.PermissionRepository) *AuthorizationUsecase {
	return &AuthorizationUsecase{permissionRepository: permissionRepository}
}

func (u *AuthorizationUsecase) HasAllPermissions(ctx context.Context, tenantCodigo string, username string, permissions []string) (bool, error) {
	permissionSet, err := u.userPermissionSet(ctx, tenantCodigo, username)
	if err != nil {
		return false, err
	}
	if _, ok := permissionSet[adminPermission]; ok {
		return true, nil
	}

	for _, permission := range permissions {
		permission = strings.TrimSpace(permission)
		if permission == "" {
			continue
		}

		if _, ok := permissionSet[permission]; !ok {
			return false, nil
		}
	}

	return true, nil
}

func (u *AuthorizationUsecase) HasAnyPermission(ctx context.Context, tenantCodigo string, username string, permissions []string) (bool, error) {
	permissionSet, err := u.userPermissionSet(ctx, tenantCodigo, username)
	if err != nil {
		return false, err
	}
	if _, ok := permissionSet[adminPermission]; ok {
		return true, nil
	}

	for _, permission := range permissions {
		permission = strings.TrimSpace(permission)
		if permission == "" {
			continue
		}

		if _, ok := permissionSet[permission]; ok {
			return true, nil
		}
	}

	return false, nil
}

func (u *AuthorizationUsecase) userPermissionSet(ctx context.Context, tenantCodigo string, username string) (map[string]struct{}, error) {
	permissions, err := u.permissionRepository.FindUserPermissions(ctx, tenantCodigo, username)
	if err != nil {
		return nil, err
	}

	return permissionSet(permissions), nil
}

func permissionSet(permissions []models.Permission) map[string]struct{} {
	out := make(map[string]struct{}, len(permissions))
	for _, permission := range permissions {
		code := strings.TrimSpace(permission.Code)
		if code == "" {
			continue
		}

		out[code] = struct{}{}
	}

	return out
}
