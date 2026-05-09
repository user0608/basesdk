package usecases

import (
	"basesdk/security/repositories"
	"context"
	"strings"
)

type AuthorizationUsecase struct {
	permissionRepository *repositories.PermissionRepository
}

func NewAuthorizationUsecase(permissionRepository *repositories.PermissionRepository) *AuthorizationUsecase {
	return &AuthorizationUsecase{permissionRepository: permissionRepository}
}

func (u *AuthorizationUsecase) HasAllPermissions(ctx context.Context, tenantCodigo string, username string, permissions []string) (bool, error) {
	for _, permission := range permissions {
		permission = strings.TrimSpace(permission)
		if permission == "" {
			continue
		}

		allowed, err := u.permissionRepository.UserHasPermission(ctx, tenantCodigo, username, permission)
		if err != nil || !allowed {
			return allowed, err
		}
	}

	return true, nil
}

func (u *AuthorizationUsecase) HasAnyPermission(ctx context.Context, tenantCodigo string, username string, permissions []string) (bool, error) {
	for _, permission := range permissions {
		permission = strings.TrimSpace(permission)
		if permission == "" {
			continue
		}

		allowed, err := u.permissionRepository.UserHasPermission(ctx, tenantCodigo, username, permission)
		if err != nil {
			return false, err
		}
		if allowed {
			return true, nil
		}
	}

	return false, nil
}
