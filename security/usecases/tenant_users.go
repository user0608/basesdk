package usecases

import (
	"basesdk/errs"
	"basesdk/security/dtos"
	"basesdk/security/models"
	"basesdk/security/repositories"
	"context"
	"strings"
	"time"
)

type TenantUsersUsecase struct {
	userRepository       *repositories.AppUserRepository
	permissionRepository *repositories.PermissionRepository
}

func NewTenantUsersUsecase(userRepository *repositories.AppUserRepository, permissionRepository *repositories.PermissionRepository) *TenantUsersUsecase {
	return &TenantUsersUsecase{userRepository: userRepository, permissionRepository: permissionRepository}
}

func (u *TenantUsersUsecase) List(ctx context.Context, tenantCodigo string) ([]dtos.TenantUserResponse, error) {
	users, err := u.userRepository.FindAppUsers(ctx, strings.TrimSpace(tenantCodigo))
	if err != nil {
		return nil, err
	}

	response := make([]dtos.TenantUserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, toTenantUserResponse(user))
	}

	return response, nil
}

func (u *TenantUsersUsecase) Find(ctx context.Context, tenantCodigo string, username string) (*dtos.TenantUserResponse, error) {
	user, err := u.userRepository.FindAppUser(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(username))
	if err != nil {
		return nil, err
	}

	response := toTenantUserResponse(*user)
	return &response, nil
}

func (u *TenantUsersUsecase) Create(ctx context.Context, tenantCodigo string, input dtos.CreateTenantUserInput, createdBy string) error {
	user := &models.AppUser{
		TenantCodigo:       strings.TrimSpace(tenantCodigo),
		Username:           strings.TrimSpace(input.Username),
		FullName:           input.FullName,
		MustChangePassword: input.MustChangePassword,
		Disabled:           false,
		CreatedBy:          createdBy,
		CreatedAt:          time.Now(),
	}
	if user.TenantCodigo == "" || user.Username == "" {
		return errs.BadRequestDirect("tenant y usuario son requeridos")
	}
	if err := user.ChangePassword(input.Password); err != nil {
		return err
	}

	return u.userRepository.CreateAppUser(ctx, user)
}

func (u *TenantUsersUsecase) Update(ctx context.Context, tenantCodigo string, username string, input dtos.UpdateTenantUserInput, updatedBy string) error {
	now := time.Now()
	return u.userRepository.UpdateAppUser(ctx, &models.AppUser{
		TenantCodigo:       strings.TrimSpace(tenantCodigo),
		Username:           strings.TrimSpace(username),
		FullName:           input.FullName,
		MustChangePassword: input.MustChangePassword,
		Disabled:           input.Disabled,
		UpdatedBy:          &updatedBy,
		UpdatedAt:          &now,
	})
}

func (u *TenantUsersUsecase) ChangePassword(ctx context.Context, tenantCodigo string, username string, input dtos.ChangePasswordInput, updatedBy string) error {
	user := &models.AppUser{}
	if err := user.ChangePassword(input.Password); err != nil {
		return err
	}

	return u.userRepository.ChangeAppUserPassword(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(username), *user.PasswordHash, updatedBy)
}

func (u *TenantUsersUsecase) Enable(ctx context.Context, tenantCodigo string, usernames []string, updatedBy string) error {
	if len(usernames) == 0 {
		return errs.BadRequestDirect("usuarios requeridos")
	}
	return u.userRepository.SetAppUsersDisabled(ctx, strings.TrimSpace(tenantCodigo), usernames, false, updatedBy)
}

func (u *TenantUsersUsecase) Disable(ctx context.Context, tenantCodigo string, usernames []string, updatedBy string) error {
	if len(usernames) == 0 {
		return errs.BadRequestDirect("usuarios requeridos")
	}
	return u.userRepository.SetAppUsersDisabled(ctx, strings.TrimSpace(tenantCodigo), usernames, true, updatedBy)
}

func (u *TenantUsersUsecase) Delete(ctx context.Context, tenantCodigo string, usernames []string) error {
	if len(usernames) == 0 {
		return errs.BadRequestDirect("usuarios requeridos")
	}
	return u.userRepository.DeleteAppUsers(ctx, strings.TrimSpace(tenantCodigo), usernames)
}

func (u *TenantUsersUsecase) FindUserPermissions(ctx context.Context, tenantCodigo string, username string) (*dtos.UserPermissionsResponse, error) {
	permissions, err := u.permissionRepository.FindUserPermissions(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(username))
	if err != nil {
		return nil, err
	}

	response := &dtos.UserPermissionsResponse{TenantCodigo: tenantCodigo, Username: username, Permissions: make([]dtos.PermissionResponse, 0, len(permissions))}
	for _, permission := range permissions {
		response.Permissions = append(response.Permissions, toPermissionResponse(permission))
	}

	return response, nil
}

func (u *TenantUsersUsecase) FindRoles(ctx context.Context, tenantCodigo string, username string) ([]dtos.RoleResponse, error) {
	roles, err := u.userRepository.FindAppUserRoles(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(username))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.RoleResponse, 0, len(roles))
	for _, role := range roles {
		response = append(response, toRoleResponse(role))
	}
	return response, nil
}

func (u *TenantUsersUsecase) ReplaceRoles(ctx context.Context, tenantCodigo string, username string, roleCodes []string, createdBy string) error {
	return u.userRepository.ReplaceAppUserRoles(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(username), roleCodes, createdBy)
}

func (u *TenantUsersUsecase) FindGroups(ctx context.Context, tenantCodigo string, username string) ([]dtos.GroupResponse, error) {
	groups, err := u.userRepository.FindAppUserGroups(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(username))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.GroupResponse, 0, len(groups))
	for _, group := range groups {
		response = append(response, toGroupResponse(group))
	}
	return response, nil
}

func (u *TenantUsersUsecase) ReplaceGroups(ctx context.Context, tenantCodigo string, username string, groupCodes []string, createdBy string) error {
	return u.userRepository.ReplaceAppUserGroups(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(username), groupCodes, createdBy)
}
