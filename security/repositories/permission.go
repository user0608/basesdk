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
	rs := r.manager.Conn(ctx).Raw(`
		select
			p.*,
			(
				select count(distinct rp.role_code)
				from role_permission rp
				where rp.permission_code = p.code
			) as roles_count,
			(
				select count(distinct gr.group_code)
				from role_permission rp
				join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
				join group_role gr on gr.tenant_codigo = rp.tenant_codigo and gr.role_code = rp.role_code
				join app_group g on g.tenant_codigo = gr.tenant_codigo and g.code = gr.group_code
				where rp.permission_code = p.code and r.disabled = false and g.disabled = false
			) as groups_count,
			(
				select count(distinct user_permission.username)
				from (
					select ur.tenant_codigo, ur.username
					from role_permission rp
					join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
					join user_role ur on ur.tenant_codigo = rp.tenant_codigo and ur.role_code = rp.role_code
					where rp.permission_code = p.code and r.disabled = false
					union
					select ug.tenant_codigo, ug.username
					from role_permission rp
					join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
					join group_role gr on gr.tenant_codigo = rp.tenant_codigo and gr.role_code = rp.role_code
					join app_group g on g.tenant_codigo = gr.tenant_codigo and g.code = gr.group_code
					join user_group ug on ug.tenant_codigo = gr.tenant_codigo and ug.group_code = gr.group_code
					where rp.permission_code = p.code and r.disabled = false and g.disabled = false
				) user_permission
			) as users_count
		from permission p
		order by p.code
	`).Scan(&permissions)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return permissions, nil
}

func (r *PermissionRepository) FindPermission(ctx context.Context, code string) (*models.Permission, error) {
	var permission models.Permission
	rs := r.manager.Conn(ctx).Raw(`
		select
			p.*,
			(
				select count(distinct rp.role_code)
				from role_permission rp
				where rp.permission_code = p.code
			) as roles_count,
			(
				select count(distinct gr.group_code)
				from role_permission rp
				join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
				join group_role gr on gr.tenant_codigo = rp.tenant_codigo and gr.role_code = rp.role_code
				join app_group g on g.tenant_codigo = gr.tenant_codigo and g.code = gr.group_code
				where rp.permission_code = p.code and r.disabled = false and g.disabled = false
			) as groups_count,
			(
				select count(distinct user_permission.username)
				from (
					select ur.tenant_codigo, ur.username
					from role_permission rp
					join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
					join user_role ur on ur.tenant_codigo = rp.tenant_codigo and ur.role_code = rp.role_code
					where rp.permission_code = p.code and r.disabled = false
					union
					select ug.tenant_codigo, ug.username
					from role_permission rp
					join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
					join group_role gr on gr.tenant_codigo = rp.tenant_codigo and gr.role_code = rp.role_code
					join app_group g on g.tenant_codigo = gr.tenant_codigo and g.code = gr.group_code
					join user_group ug on ug.tenant_codigo = gr.tenant_codigo and ug.group_code = gr.group_code
					where rp.permission_code = p.code and r.disabled = false and g.disabled = false
				) user_permission
			) as users_count
		from permission p
		where p.code = ?
	`, code).Scan(&permission)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return nil, errs.NotFoundDirect("permiso no encontrado")
	}

	return &permission, nil
}

func (r *PermissionRepository) FindPermissionRoles(ctx context.Context, code string) ([]models.Role, error) {
	var roles []models.Role
	rs := r.manager.Conn(ctx).Raw(`
		select distinct r.*
		from role r
		join role_permission rp on rp.tenant_codigo = r.tenant_codigo and rp.role_code = r.code
		where rp.permission_code = ?
		order by r.code
	`, code).Scan(&roles)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return roles, nil
}

func (r *PermissionRepository) FindPermissionGroups(ctx context.Context, code string) ([]models.AppGroup, error) {
	var groups []models.AppGroup
	rs := r.manager.Conn(ctx).Raw(`
		select distinct g.*
		from app_group g
		join group_role gr on gr.tenant_codigo = g.tenant_codigo and gr.group_code = g.code
		join role r on r.tenant_codigo = gr.tenant_codigo and r.code = gr.role_code
		join role_permission rp on rp.tenant_codigo = gr.tenant_codigo and rp.role_code = gr.role_code
		where rp.permission_code = ? and r.disabled = false and g.disabled = false
		order by g.code
	`, code).Scan(&groups)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return groups, nil
}

func (r *PermissionRepository) FindPermissionUsers(ctx context.Context, code string) ([]models.AppUser, error) {
	var users []models.AppUser
	rs := r.manager.Conn(ctx).Raw(`
		select distinct u.*
		from app_user u
		join (
			select ur.tenant_codigo, ur.username
			from role_permission rp
			join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
			join user_role ur on ur.tenant_codigo = rp.tenant_codigo and ur.role_code = rp.role_code
			where rp.permission_code = ? and r.disabled = false
			union
			select ug.tenant_codigo, ug.username
			from role_permission rp
			join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
			join group_role gr on gr.tenant_codigo = rp.tenant_codigo and gr.role_code = rp.role_code
			join app_group g on g.tenant_codigo = gr.tenant_codigo and g.code = gr.group_code
			join user_group ug on ug.tenant_codigo = gr.tenant_codigo and ug.group_code = gr.group_code
			where rp.permission_code = ? and r.disabled = false and g.disabled = false
		) permitted_user on permitted_user.tenant_codigo = u.tenant_codigo and permitted_user.username = u.username
		order by u.username
	`, code, code).Scan(&users)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return users, nil
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
