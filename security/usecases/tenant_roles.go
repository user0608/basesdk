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

type TenantRolesUsecase struct {
	repository *repositories.RoleRepository
}

func NewTenantRolesUsecase(repository *repositories.RoleRepository) *TenantRolesUsecase {
	return &TenantRolesUsecase{repository: repository}
}

func (u *TenantRolesUsecase) List(ctx context.Context, tenantCodigo string) ([]dtos.RoleResponse, error) {
	roles, err := u.repository.FindRoles(ctx, strings.TrimSpace(tenantCodigo))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.RoleResponse, 0, len(roles))
	for _, role := range roles {
		response = append(response, toRoleResponse(role))
	}
	return response, nil
}

func (u *TenantRolesUsecase) Find(ctx context.Context, tenantCodigo string, code string) (*dtos.RoleResponse, error) {
	role, err := u.repository.FindRole(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	response := toRoleResponse(*role)
	return &response, nil
}

func (u *TenantRolesUsecase) Create(ctx context.Context, tenantCodigo string, input dtos.CreateRoleInput, createdBy string) error {
	role := &models.Role{TenantCodigo: strings.TrimSpace(tenantCodigo), Code: strings.TrimSpace(input.Code), Description: input.Description, CreatedBy: createdBy, CreatedAt: time.Now()}
	if role.TenantCodigo == "" || role.Code == "" {
		return errs.BadRequestDirect("tenant y rol son requeridos")
	}
	return u.repository.CreateRole(ctx, role)
}

func (u *TenantRolesUsecase) Update(ctx context.Context, tenantCodigo string, code string, input dtos.UpdateRoleInput, updatedBy string) error {
	now := time.Now()
	return u.repository.UpdateRole(ctx, &models.Role{TenantCodigo: strings.TrimSpace(tenantCodigo), Code: strings.TrimSpace(code), Description: input.Description, Disabled: input.Disabled, UpdatedBy: &updatedBy, UpdatedAt: &now})
}

func (u *TenantRolesUsecase) Enable(ctx context.Context, tenantCodigo string, codes []string, updatedBy string) error {
	if len(codes) == 0 {
		return errs.BadRequestDirect("roles requeridos")
	}
	return u.repository.SetRolesDisabled(ctx, strings.TrimSpace(tenantCodigo), codes, false, updatedBy)
}

func (u *TenantRolesUsecase) Disable(ctx context.Context, tenantCodigo string, codes []string, updatedBy string) error {
	if len(codes) == 0 {
		return errs.BadRequestDirect("roles requeridos")
	}
	return u.repository.SetRolesDisabled(ctx, strings.TrimSpace(tenantCodigo), codes, true, updatedBy)
}

func (u *TenantRolesUsecase) Delete(ctx context.Context, tenantCodigo string, codes []string) error {
	if len(codes) == 0 {
		return errs.BadRequestDirect("roles requeridos")
	}
	return u.repository.DeleteRoles(ctx, strings.TrimSpace(tenantCodigo), codes)
}

func (u *TenantRolesUsecase) FindPermissions(ctx context.Context, tenantCodigo string, code string) ([]dtos.PermissionResponse, error) {
	permissions, err := u.repository.FindRolePermissions(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		response = append(response, toPermissionResponse(permission))
	}
	return response, nil
}

func (u *TenantRolesUsecase) ReplacePermissions(ctx context.Context, tenantCodigo string, code string, permissionCodes []string, createdBy string) error {
	return u.repository.ReplaceRolePermissions(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code), permissionCodes, createdBy)
}

func (u *TenantRolesUsecase) FindUsers(ctx context.Context, tenantCodigo string, code string) ([]dtos.TenantUserResponse, error) {
	users, err := u.repository.FindRoleUsers(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.TenantUserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, toTenantUserResponse(user))
	}
	return response, nil
}

func (u *TenantRolesUsecase) ReplaceUsers(ctx context.Context, tenantCodigo string, code string, usernames []string, createdBy string) error {
	return u.repository.ReplaceRoleUsers(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code), usernames, createdBy)
}

func (u *TenantRolesUsecase) FindGroups(ctx context.Context, tenantCodigo string, code string) ([]dtos.GroupResponse, error) {
	groups, err := u.repository.FindRoleGroups(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.GroupResponse, 0, len(groups))
	for _, group := range groups {
		response = append(response, toGroupResponse(group))
	}
	return response, nil
}

func (u *TenantRolesUsecase) ReplaceGroups(ctx context.Context, tenantCodigo string, code string, groupCodes []string, createdBy string) error {
	return u.repository.ReplaceRoleGroups(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code), groupCodes, createdBy)
}
