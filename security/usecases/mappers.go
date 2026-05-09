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
		Email:              user.Email,
		FullName:           user.FullName,
		EmailVerified:      user.EmailVerified,
		MustChangePassword: user.MustChangePassword,
		LastLoginAt:        user.LastLoginAt,
		Disabled:           user.Disabled,
	}
}

func toRoleResponse(role models.Role) dtos.RoleResponse {
	return dtos.RoleResponse{
		TenantCodigo: role.TenantCodigo,
		Code:         role.Code,
		Description:  role.Description,
		Disabled:     role.Disabled,
	}
}

func toGroupResponse(group models.AppGroup) dtos.GroupResponse {
	return dtos.GroupResponse{
		TenantCodigo: group.TenantCodigo,
		Code:         group.Code,
		Description:  group.Description,
		Disabled:     group.Disabled,
	}
}

func toPermissionResponse(permission models.Permission) dtos.PermissionResponse {
	return dtos.PermissionResponse{
		Code:        permission.Code,
		Description: permission.Description,
	}
}
