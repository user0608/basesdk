package repositories

import (
	"basesdk/connection"
	"basesdk/errs"
	"basesdk/security/models"
	"context"
)

type PermissionRepository struct {
	manager connection.StorageManager
}

func NewPermissionRepository(manager connection.StorageManager) *PermissionRepository {
	return &PermissionRepository{manager: manager}
}

func (r *PermissionRepository) FindPermissions(ctx context.Context) ([]models.Permission, error) {
	var permissions []models.Permission
	rs := r.manager.Conn(ctx).Order("code").Find(&permissions)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return permissions, nil
}

func (r *PermissionRepository) FindPermission(ctx context.Context, code string) (*models.Permission, error) {
	var permission models.Permission
	rs := r.manager.Conn(ctx).Where("code = ?", code).First(&permission)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return &permission, nil
}

func (r *PermissionRepository) FindUserPermissions(ctx context.Context, tenantCodigo string, username string) ([]models.Permission, error) {
	var permissions []models.Permission
	rs := r.manager.Conn(ctx).Raw(`
		select distinct p.code, p.description
		from permission p
		join role_permission rp on rp.permission_code = p.code
		join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
		join user_role ur on ur.tenant_codigo = rp.tenant_codigo and ur.role_code = rp.role_code
		where ur.tenant_codigo = ? and ur.username = ? and r.disabled = false
		union
		select distinct p.code, p.description
		from permission p
		join role_permission rp on rp.permission_code = p.code
		join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
		join group_role gr on gr.tenant_codigo = rp.tenant_codigo and gr.role_code = rp.role_code
		join app_group g on g.tenant_codigo = gr.tenant_codigo and g.code = gr.group_code
		join user_group ug on ug.tenant_codigo = gr.tenant_codigo and ug.group_code = gr.group_code
		where ug.tenant_codigo = ? and ug.username = ? and r.disabled = false and g.disabled = false
		order by code
	`, tenantCodigo, username, tenantCodigo, username).Scan(&permissions)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return permissions, nil
}

func (r *PermissionRepository) UserHasPermission(ctx context.Context, tenantCodigo string, username string, permissionCode string) (bool, error) {
	var count int64
	rs := r.manager.Conn(ctx).Raw(`
		select count(*)
		from (
			select p.code
			from permission p
			join role_permission rp on rp.permission_code = p.code
			join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
			join user_role ur on ur.tenant_codigo = rp.tenant_codigo and ur.role_code = rp.role_code
			where ur.tenant_codigo = ? and ur.username = ? and p.code = ? and r.disabled = false
			union
			select p.code
			from permission p
			join role_permission rp on rp.permission_code = p.code
			join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
			join group_role gr on gr.tenant_codigo = rp.tenant_codigo and gr.role_code = rp.role_code
			join app_group g on g.tenant_codigo = gr.tenant_codigo and g.code = gr.group_code
			join user_group ug on ug.tenant_codigo = gr.tenant_codigo and ug.group_code = gr.group_code
			where ug.tenant_codigo = ? and ug.username = ? and p.code = ? and r.disabled = false and g.disabled = false
		) user_permission
	`, tenantCodigo, username, permissionCode, tenantCodigo, username, permissionCode).Scan(&count)
	if rs.Error != nil {
		return false, errs.Pgf(rs.Error)
	}

	return count > 0, nil
}
