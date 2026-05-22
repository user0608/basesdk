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

type TenantGroupsUsecase struct {
	repository *repositories.GroupRepository
}

func NewTenantGroupsUsecase(repository *repositories.GroupRepository) *TenantGroupsUsecase {
	return &TenantGroupsUsecase{repository: repository}
}

func (u *TenantGroupsUsecase) List(ctx context.Context, tenantCodigo string) ([]dtos.GroupResponse, error) {
	groups, err := u.repository.FindGroups(ctx, strings.TrimSpace(tenantCodigo))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.GroupResponse, 0, len(groups))
	for _, group := range groups {
		response = append(response, toGroupResponse(group))
	}
	return response, nil
}

func (u *TenantGroupsUsecase) Find(ctx context.Context, tenantCodigo string, code string) (*dtos.GroupResponse, error) {
	group, err := u.repository.FindGroup(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	response := toGroupResponse(*group)
	return &response, nil
}

func (u *TenantGroupsUsecase) Create(ctx context.Context, tenantCodigo string, input dtos.CreateGroupInput, createdBy string) error {
	group := &models.AppGroup{TenantCodigo: strings.TrimSpace(tenantCodigo), Code: strings.TrimSpace(input.Code), Description: input.Description, CreatedBy: createdBy, CreatedAt: time.Now()}
	if group.TenantCodigo == "" || group.Code == "" {
		return errs.BadRequestDirect("tenant y grupo son requeridos")
	}
	return u.repository.CreateGroup(ctx, group)
}

func (u *TenantGroupsUsecase) Update(ctx context.Context, tenantCodigo string, code string, input dtos.UpdateGroupInput, updatedBy string) error {
	now := time.Now()
	return u.repository.UpdateGroup(ctx, &models.AppGroup{TenantCodigo: strings.TrimSpace(tenantCodigo), Code: strings.TrimSpace(code), Description: input.Description, Disabled: input.Disabled, UpdatedBy: &updatedBy, UpdatedAt: &now})
}

func (u *TenantGroupsUsecase) Enable(ctx context.Context, tenantCodigo string, codes []string, updatedBy string) error {
	if len(codes) == 0 {
		return errs.BadRequestDirect("grupos requeridos")
	}
	return u.repository.SetGroupsDisabled(ctx, strings.TrimSpace(tenantCodigo), codes, false, updatedBy)
}

func (u *TenantGroupsUsecase) Disable(ctx context.Context, tenantCodigo string, codes []string, updatedBy string) error {
	if len(codes) == 0 {
		return errs.BadRequestDirect("grupos requeridos")
	}
	return u.repository.SetGroupsDisabled(ctx, strings.TrimSpace(tenantCodigo), codes, true, updatedBy)
}

func (u *TenantGroupsUsecase) Delete(ctx context.Context, tenantCodigo string, codes []string) error {
	if len(codes) == 0 {
		return errs.BadRequestDirect("grupos requeridos")
	}
	return u.repository.DeleteGroups(ctx, strings.TrimSpace(tenantCodigo), codes)
}

func (u *TenantGroupsUsecase) FindUsers(ctx context.Context, tenantCodigo string, code string) ([]dtos.TenantUserResponse, error) {
	users, err := u.repository.FindGroupUsers(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.TenantUserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, toTenantUserResponse(user))
	}
	return response, nil
}

func (u *TenantGroupsUsecase) ReplaceUsers(ctx context.Context, tenantCodigo string, code string, usernames []string, createdBy string) error {
	return u.repository.ReplaceGroupUsers(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code), usernames, createdBy)
}

func (u *TenantGroupsUsecase) FindRoles(ctx context.Context, tenantCodigo string, code string) ([]dtos.RoleResponse, error) {
	roles, err := u.repository.FindGroupRoles(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.RoleResponse, 0, len(roles))
	for _, role := range roles {
		response = append(response, toRoleResponse(role))
	}
	return response, nil
}

func (u *TenantGroupsUsecase) ReplaceRoles(ctx context.Context, tenantCodigo string, code string, roleCodes []string, createdBy string) error {
	return u.repository.ReplaceGroupRoles(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code), roleCodes, createdBy)
}

func (u *TenantGroupsUsecase) FindPermissions(ctx context.Context, tenantCodigo string, code string) ([]dtos.PermissionResponse, error) {
	permissions, err := u.repository.FindGroupPermissions(ctx, strings.TrimSpace(tenantCodigo), strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	response := make([]dtos.PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		response = append(response, toPermissionResponse(permission))
	}
	return response, nil
}
