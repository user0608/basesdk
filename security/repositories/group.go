package repositories

import (
	"basesdk/connection"
	"basesdk/errs"
	"basesdk/security/models"
	"context"
	"time"
)

type GroupRepository struct {
	manager connection.StorageManager
}

func NewGroupRepository(manager connection.StorageManager) *GroupRepository {
	return &GroupRepository{manager: manager}
}

func (r *GroupRepository) FindGroups(ctx context.Context, tenantCodigo string) ([]models.AppGroup, error) {
	var groups []models.AppGroup
	rs := r.manager.Conn(ctx).Raw(`
		select
			g.*,
			(
				select count(*)
				from user_group ug
				where ug.tenant_codigo = g.tenant_codigo and ug.group_code = g.code
			) as users_count,
			(
				select count(*)
				from group_role gr
				where gr.tenant_codigo = g.tenant_codigo and gr.group_code = g.code
			) as roles_count,
			(
				select count(distinct rp.permission_code)
				from group_role gr
				join role r on r.tenant_codigo = gr.tenant_codigo and r.code = gr.role_code
				join role_permission rp on rp.tenant_codigo = gr.tenant_codigo and rp.role_code = gr.role_code
				where gr.tenant_codigo = g.tenant_codigo and gr.group_code = g.code and r.disabled = false and g.disabled = false
			) as permissions_count
		from app_group g
		where g.tenant_codigo = ?
		order by g.code
	`, tenantCodigo).Scan(&groups)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return groups, nil
}

func (r *GroupRepository) FindGroup(ctx context.Context, tenantCodigo string, code string) (*models.AppGroup, error) {
	var group models.AppGroup
	rs := r.manager.Conn(ctx).Raw(`
		select
			g.*,
			(
				select count(*)
				from user_group ug
				where ug.tenant_codigo = g.tenant_codigo and ug.group_code = g.code
			) as users_count,
			(
				select count(*)
				from group_role gr
				where gr.tenant_codigo = g.tenant_codigo and gr.group_code = g.code
			) as roles_count,
			(
				select count(distinct rp.permission_code)
				from group_role gr
				join role r on r.tenant_codigo = gr.tenant_codigo and r.code = gr.role_code
				join role_permission rp on rp.tenant_codigo = gr.tenant_codigo and rp.role_code = gr.role_code
				where gr.tenant_codigo = g.tenant_codigo and gr.group_code = g.code and r.disabled = false and g.disabled = false
			) as permissions_count
		from app_group g
		where g.tenant_codigo = ? and g.code = ?
	`, tenantCodigo, code).Scan(&group)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return nil, errs.NotFoundDirect("grupo no encontrado")
	}

	return &group, nil
}

func (r *GroupRepository) CreateGroup(ctx context.Context, group *models.AppGroup) error {
	return errs.Pgf(r.manager.Conn(ctx).Create(group).Error)
}

func (r *GroupRepository) UpdateGroup(ctx context.Context, group *models.AppGroup) error {
	rs := r.manager.Conn(ctx).Model(&models.AppGroup{}).
		Where("tenant_codigo = ? and code = ?", group.TenantCodigo, group.Code).
		Updates(map[string]any{
			"description": group.Description,
			"disabled":    group.Disabled,
			"updated_by":  group.UpdatedBy,
			"updated_at":  group.UpdatedAt,
		})
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return errs.NotFoundDirect("grupo no encontrado")
	}

	return nil
}

func (r *GroupRepository) SetGroupsDisabled(ctx context.Context, tenantCodigo string, codes []string, disabled bool, updatedBy string) error {
	rs := r.manager.Conn(ctx).Model(&models.AppGroup{}).
		Where("tenant_codigo = ? and code in ?", tenantCodigo, codes).
		Updates(map[string]any{
			"disabled":   disabled,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		})
	return errs.Pgf(rs.Error)
}

func (r *GroupRepository) DeleteGroups(ctx context.Context, tenantCodigo string, codes []string) error {
	return r.manager.WithTx(ctx, func(ctx context.Context) error {
		tx := r.manager.Conn(ctx)

		if err := errs.Pgf(tx.Delete(&models.UserGroup{}, "tenant_codigo = ? and group_code in ?", tenantCodigo, codes).Error); err != nil {
			return err
		}
		if err := errs.Pgf(tx.Delete(&models.GroupRole{}, "tenant_codigo = ? and group_code in ?", tenantCodigo, codes).Error); err != nil {
			return err
		}

		return errs.Pgf(tx.Delete(&models.AppGroup{}, "tenant_codigo = ? and code in ?", tenantCodigo, codes).Error)
	})
}

func (r *GroupRepository) FindGroupUsers(ctx context.Context, tenantCodigo string, groupCode string) ([]models.AppUser, error) {
	var users []models.AppUser
	rs := r.manager.Conn(ctx).Raw(`
		select u.*
		from app_user u
		join user_group ug on ug.tenant_codigo = u.tenant_codigo and ug.username = u.username
		where ug.tenant_codigo = ? and ug.group_code = ?
		order by u.username
	`, tenantCodigo, groupCode).Scan(&users)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return users, nil
}

func (r *GroupRepository) ReplaceGroupUsers(ctx context.Context, tenantCodigo string, groupCode string, usernames []string, createdBy string) error {
	return r.manager.WithTx(ctx, func(ctx context.Context) error {
		tx := r.manager.Conn(ctx)
		if err := errs.Pgf(tx.Delete(&models.UserGroup{}, "tenant_codigo = ? and group_code = ?", tenantCodigo, groupCode).Error); err != nil {
			return err
		}

		now := time.Now()
		for _, username := range usernames {
			entry := &models.UserGroup{TenantCodigo: tenantCodigo, GroupCode: groupCode, Username: username, CreatedBy: createdBy, CreatedAt: now}
			if err := errs.Pgf(tx.Create(entry).Error); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GroupRepository) FindGroupRoles(ctx context.Context, tenantCodigo string, groupCode string) ([]models.Role, error) {
	var roles []models.Role
	rs := r.manager.Conn(ctx).Raw(`
		select r.*
		from role r
		join group_role gr on gr.tenant_codigo = r.tenant_codigo and gr.role_code = r.code
		where gr.tenant_codigo = ? and gr.group_code = ?
		order by r.code
	`, tenantCodigo, groupCode).Scan(&roles)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return roles, nil
}

func (r *GroupRepository) ReplaceGroupRoles(ctx context.Context, tenantCodigo string, groupCode string, roleCodes []string, createdBy string) error {
	return r.manager.WithTx(ctx, func(ctx context.Context) error {
		tx := r.manager.Conn(ctx)
		if err := errs.Pgf(tx.Delete(&models.GroupRole{}, "tenant_codigo = ? and group_code = ?", tenantCodigo, groupCode).Error); err != nil {
			return err
		}

		now := time.Now()
		for _, roleCode := range roleCodes {
			entry := &models.GroupRole{TenantCodigo: tenantCodigo, GroupCode: groupCode, RoleCode: roleCode, CreatedBy: createdBy, CreatedAt: now}
			if err := errs.Pgf(tx.Create(entry).Error); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GroupRepository) FindGroupPermissions(ctx context.Context, tenantCodigo string, groupCode string) ([]models.Permission, error) {
	var permissions []models.Permission
	rs := r.manager.Conn(ctx).Raw(`
		select distinct p.code, p.description
		from permission p
		join role_permission rp on rp.permission_code = p.code
		join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
		join group_role gr on gr.tenant_codigo = rp.tenant_codigo and gr.role_code = rp.role_code
		join app_group g on g.tenant_codigo = gr.tenant_codigo and g.code = gr.group_code
		where gr.tenant_codigo = ? and gr.group_code = ? and r.disabled = false and g.disabled = false
		order by code
	`, tenantCodigo, groupCode).Scan(&permissions)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return permissions, nil
}
