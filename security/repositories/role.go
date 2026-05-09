package repositories

import (
	"basesdk/connection"
	"basesdk/errs"
	"basesdk/security/models"
	"context"
	"time"
)

type RoleRepository struct {
	manager connection.StorageManager
}

func NewRoleRepository(manager connection.StorageManager) *RoleRepository {
	return &RoleRepository{manager: manager}
}

func (r *RoleRepository) FindRoles(ctx context.Context, tenantCodigo string) ([]models.Role, error) {
	tx := r.manager.Conn(ctx)
	var roles []models.Role

	rs := tx.Where("tenant_codigo = ?", tenantCodigo).Order("code").Find(&roles)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return roles, nil
}

func (r *RoleRepository) FindRole(ctx context.Context, tenantCodigo string, code string) (*models.Role, error) {
	tx := r.manager.Conn(ctx)
	var role models.Role

	rs := tx.Where("tenant_codigo = ? and code = ?", tenantCodigo, code).First(&role)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return &role, nil
}

func (r *RoleRepository) CreateRole(ctx context.Context, role *models.Role) error {
	return errs.Pgf(r.manager.Conn(ctx).Create(role).Error)
}

func (r *RoleRepository) UpdateRole(ctx context.Context, role *models.Role) error {
	rs := r.manager.Conn(ctx).Model(&models.Role{}).
		Where("tenant_codigo = ? and code = ?", role.TenantCodigo, role.Code).
		Updates(map[string]any{
			"description": role.Description,
			"disabled":    role.Disabled,
			"updated_by":  role.UpdatedBy,
			"updated_at":  role.UpdatedAt,
		})
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return errs.NotFoundDirect("rol no encontrado")
	}

	return nil
}

func (r *RoleRepository) SetRolesDisabled(ctx context.Context, tenantCodigo string, codes []string, disabled bool, updatedBy string) error {
	rs := r.manager.Conn(ctx).Model(&models.Role{}).
		Where("tenant_codigo = ? and code in ?", tenantCodigo, codes).
		Updates(map[string]any{
			"disabled":   disabled,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		})
	return errs.Pgf(rs.Error)
}

func (r *RoleRepository) DeleteRoles(ctx context.Context, tenantCodigo string, codes []string) error {
	return r.manager.WithTx(ctx, func(ctx context.Context) error {
		tx := r.manager.Conn(ctx)

		if err := errs.Pgf(tx.Delete(&models.RolePermission{}, "tenant_codigo = ? and role_code in ?", tenantCodigo, codes).Error); err != nil {
			return err
		}
		if err := errs.Pgf(tx.Delete(&models.UserRole{}, "tenant_codigo = ? and role_code in ?", tenantCodigo, codes).Error); err != nil {
			return err
		}
		if err := errs.Pgf(tx.Delete(&models.GroupRole{}, "tenant_codigo = ? and role_code in ?", tenantCodigo, codes).Error); err != nil {
			return err
		}

		return errs.Pgf(tx.Delete(&models.Role{}, "tenant_codigo = ? and code in ?", tenantCodigo, codes).Error)
	})
}

func (r *RoleRepository) FindRolePermissions(ctx context.Context, tenantCodigo string, roleCode string) ([]models.Permission, error) {
	var permissions []models.Permission
	rs := r.manager.Conn(ctx).Raw(`
		select p.*
		from permission p
		join role_permission rp on rp.permission_code = p.code
		where rp.tenant_codigo = ? and rp.role_code = ?
		order by p.code
	`, tenantCodigo, roleCode).Scan(&permissions)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return permissions, nil
}

func (r *RoleRepository) ReplaceRolePermissions(ctx context.Context, tenantCodigo string, roleCode string, permissionCodes []string, createdBy string) error {
	return r.manager.WithTx(ctx, func(ctx context.Context) error {
		tx := r.manager.Conn(ctx)

		if err := errs.Pgf(tx.Delete(&models.RolePermission{}, "tenant_codigo = ? and role_code = ?", tenantCodigo, roleCode).Error); err != nil {
			return err
		}

		now := time.Now()
		for _, permissionCode := range permissionCodes {
			entry := &models.RolePermission{
				TenantCodigo:   tenantCodigo,
				RoleCode:       roleCode,
				PermissionCode: permissionCode,
				CreatedBy:      createdBy,
				CreatedAt:      now,
			}
			if err := errs.Pgf(tx.Create(entry).Error); err != nil {
				return err
			}
		}

		return nil
	})
}
