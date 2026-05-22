package usecases

import (
	"basesdk/security/dtos"
	"basesdk/security/models"
)

func toSystemUserResponse(user models.SystemAccount) dtos.SystemUserResponse {
	return dtos.SystemUserResponse{
		Username:  user.Username,
		Disabled:  user.Disabled,
		CreatedBy: user.CreatedBy,
		CreatedAt: user.CreatedAt,
		UpdatedBy: user.UpdatedBy,
		UpdatedAt: user.UpdatedAt,
	}
}

func toTenantUserResponse(user models.AppUser) dtos.TenantUserResponse {
	return dtos.TenantUserResponse{
		TenantCodigo:       user.TenantCodigo,
		Username:           user.Username,
		FullName:           user.FullName,
		MustChangePassword: user.MustChangePassword,
		LastLoginAt:        user.LastLoginAt,
		Disabled:           user.Disabled,
		RolesCount:         user.RolesCount,
		GroupsCount:        user.GroupsCount,
		PermissionsCount:   user.PermissionsCount,
	}
}

func toTenantResponse(tenant models.Tenant) dtos.TenantResponse {
	return dtos.TenantResponse{
		Codigo:         tenant.Codigo,
		Name:           tenant.Name,
		Timezone:       tenant.Timezone,
		MaxActiveUsers: tenant.MaxActiveUsers,
		Disabled:       tenant.Disabled,
		ExpiresAt:      tenant.ExpiresAt,
		UsersCount:     tenant.UsersCount,
		RolesCount:     tenant.RolesCount,
		GroupsCount:    tenant.GroupsCount,
		CreatedBy:      tenant.CreatedBy,
		CreatedAt:      tenant.CreatedAt,
		UpdatedBy:      tenant.UpdatedBy,
		UpdatedAt:      tenant.UpdatedAt,
	}
}

func toRoleResponse(role models.Role) dtos.RoleResponse {
	return dtos.RoleResponse{
		TenantCodigo:     role.TenantCodigo,
		Code:             role.Code,
		Description:      role.Description,
		Disabled:         role.Disabled,
		UsersCount:       role.UsersCount,
		GroupsCount:      role.GroupsCount,
		PermissionsCount: role.PermissionsCount,
	}
}

func toGroupResponse(group models.AppGroup) dtos.GroupResponse {
	return dtos.GroupResponse{
		TenantCodigo:     group.TenantCodigo,
		Code:             group.Code,
		Description:      group.Description,
		Disabled:         group.Disabled,
		UsersCount:       group.UsersCount,
		RolesCount:       group.RolesCount,
		PermissionsCount: group.PermissionsCount,
	}
}

func toPermissionResponse(permission models.Permission) dtos.PermissionResponse {
	return dtos.PermissionResponse{
		Code:        permission.Code,
		Description: permission.Description,
		RolesCount:  permission.RolesCount,
		GroupsCount: permission.GroupsCount,
		UsersCount:  permission.UsersCount,
	}
}
