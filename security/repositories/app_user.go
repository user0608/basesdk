package repositories

import (
	"basesdk/connection"
	"basesdk/errs"
	"basesdk/security/models"
	"context"
	"time"
)

type AppUserRepository struct {
	manager connection.StorageManager
}

func NewAppUserRepository(manager connection.StorageManager) *AppUserRepository {
	return &AppUserRepository{
		manager: manager,
	}
}

func (r *AppUserRepository) FindTenant(ctx context.Context, tenantCodigo string) (*models.Tenant, error) {
	tx := r.manager.Conn(ctx)
	var tenant models.Tenant

	rs := tx.Where("codigo = ?", tenantCodigo).First(&tenant)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return &tenant, nil
}

func (r *AppUserRepository) FindAppUser(ctx context.Context, tenantCodigo string, username string) (*models.AppUser, error) {
	var user models.AppUser

	rs := r.manager.Conn(ctx).Raw(`
		select
			u.*,
			(
				select count(*)
				from user_role ur
				where ur.tenant_codigo = u.tenant_codigo and ur.username = u.username
			) as roles_count,
			(
				select count(*)
				from user_group ug
				where ug.tenant_codigo = u.tenant_codigo and ug.username = u.username
			) as groups_count,
			(
				select count(distinct user_permission.code)
				from (
					select p.code
					from permission p
					join role_permission rp on rp.permission_code = p.code
					join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
					join user_role ur on ur.tenant_codigo = rp.tenant_codigo and ur.role_code = rp.role_code
					where ur.tenant_codigo = u.tenant_codigo and ur.username = u.username and r.disabled = false
					union
					select p.code
					from permission p
					join role_permission rp on rp.permission_code = p.code
					join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
					join group_role gr on gr.tenant_codigo = rp.tenant_codigo and gr.role_code = rp.role_code
					join app_group g on g.tenant_codigo = gr.tenant_codigo and g.code = gr.group_code
					join user_group ug on ug.tenant_codigo = gr.tenant_codigo and ug.group_code = gr.group_code
					where ug.tenant_codigo = u.tenant_codigo and ug.username = u.username and r.disabled = false and g.disabled = false
				) user_permission
			) as permissions_count
		from app_user u
		where u.tenant_codigo = ? and u.username = ?
	`, tenantCodigo, username).Scan(&user)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return nil, errs.NotFoundDirect("usuario no encontrado")
	}

	return &user, nil
}

func (r *AppUserRepository) FindAppUsers(ctx context.Context, tenantCodigo string) ([]models.AppUser, error) {
	var users []models.AppUser

	rs := r.manager.Conn(ctx).Raw(`
		select
			u.*,
			(
				select count(*)
				from user_role ur
				where ur.tenant_codigo = u.tenant_codigo and ur.username = u.username
			) as roles_count,
			(
				select count(*)
				from user_group ug
				where ug.tenant_codigo = u.tenant_codigo and ug.username = u.username
			) as groups_count,
			(
				select count(distinct user_permission.code)
				from (
					select p.code
					from permission p
					join role_permission rp on rp.permission_code = p.code
					join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
					join user_role ur on ur.tenant_codigo = rp.tenant_codigo and ur.role_code = rp.role_code
					where ur.tenant_codigo = u.tenant_codigo and ur.username = u.username and r.disabled = false
					union
					select p.code
					from permission p
					join role_permission rp on rp.permission_code = p.code
					join role r on r.tenant_codigo = rp.tenant_codigo and r.code = rp.role_code
					join group_role gr on gr.tenant_codigo = rp.tenant_codigo and gr.role_code = rp.role_code
					join app_group g on g.tenant_codigo = gr.tenant_codigo and g.code = gr.group_code
					join user_group ug on ug.tenant_codigo = gr.tenant_codigo and ug.group_code = gr.group_code
					where ug.tenant_codigo = u.tenant_codigo and ug.username = u.username and r.disabled = false and g.disabled = false
				) user_permission
			) as permissions_count
		from app_user u
		where u.tenant_codigo = ?
		order by u.username
	`, tenantCodigo).Scan(&users)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return users, nil
}

func (r *AppUserRepository) CreateAppUser(ctx context.Context, user *models.AppUser) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Create(user)
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}

func (r *AppUserRepository) UpdateAppUser(ctx context.Context, user *models.AppUser) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Model(&models.AppUser{}).
		Where("tenant_codigo = ? and username = ?", user.TenantCodigo, user.Username).
		Updates(map[string]any{
			"full_name":            user.FullName,
			"must_change_password": user.MustChangePassword,
			"disabled":             user.Disabled,
			"updated_by":           user.UpdatedBy,
			"updated_at":           user.UpdatedAt,
		})
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return errs.NotFoundDirect("usuario no encontrado")
	}

	return nil
}

func (r *AppUserRepository) ChangeAppUserPassword(ctx context.Context, tenantCodigo string, username string, passwordHash string, updatedBy string) error {
	tx := r.manager.Conn(ctx)
	now := time.Now()

	rs := tx.Model(&models.AppUser{}).
		Where("tenant_codigo = ? and username = ?", tenantCodigo, username).
		Updates(map[string]any{
			"password_hash": passwordHash,
			"updated_by":    updatedBy,
			"updated_at":    now,
		})
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return errs.NotFoundDirect("usuario no encontrado")
	}

	return nil
}

func (r *AppUserRepository) SetAppUsersDisabled(ctx context.Context, tenantCodigo string, usernames []string, disabled bool, updatedBy string) error {
	tx := r.manager.Conn(ctx)
	now := time.Now()

	rs := tx.Model(&models.AppUser{}).
		Where("tenant_codigo = ? and username in ?", tenantCodigo, usernames).
		Updates(map[string]any{
			"disabled":   disabled,
			"updated_by": updatedBy,
			"updated_at": now,
		})
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}

func (r *AppUserRepository) DeleteAppUsers(ctx context.Context, tenantCodigo string, usernames []string) error {
	return r.manager.WithTx(ctx, func(ctx context.Context) error {
		tx := r.manager.Conn(ctx)

		if err := errs.Pgf(tx.Delete(&models.UserRole{}, "tenant_codigo = ? and username in ?", tenantCodigo, usernames).Error); err != nil {
			return err
		}
		if err := errs.Pgf(tx.Delete(&models.UserGroup{}, "tenant_codigo = ? and username in ?", tenantCodigo, usernames).Error); err != nil {
			return err
		}

		return errs.Pgf(tx.Delete(&models.AppUser{}, "tenant_codigo = ? and username in ?", tenantCodigo, usernames).Error)
	})
}

func (r *AppUserRepository) FindAppUserRoles(ctx context.Context, tenantCodigo string, username string) ([]models.Role, error) {
	var roles []models.Role
	rs := r.manager.Conn(ctx).Raw(`
		select r.*
		from role r
		join user_role ur on ur.tenant_codigo = r.tenant_codigo and ur.role_code = r.code
		where ur.tenant_codigo = ? and ur.username = ?
		order by r.code
	`, tenantCodigo, username).Scan(&roles)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return roles, nil
}

func (r *AppUserRepository) ReplaceAppUserRoles(ctx context.Context, tenantCodigo string, username string, roleCodes []string, createdBy string) error {
	return r.manager.WithTx(ctx, func(ctx context.Context) error {
		tx := r.manager.Conn(ctx)
		if err := errs.Pgf(tx.Delete(&models.UserRole{}, "tenant_codigo = ? and username = ?", tenantCodigo, username).Error); err != nil {
			return err
		}

		now := time.Now()
		for _, roleCode := range roleCodes {
			entry := &models.UserRole{TenantCodigo: tenantCodigo, Username: username, RoleCode: roleCode, CreatedBy: createdBy, CreatedAt: now}
			if err := errs.Pgf(tx.Create(entry).Error); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *AppUserRepository) FindAppUserGroups(ctx context.Context, tenantCodigo string, username string) ([]models.AppGroup, error) {
	var groups []models.AppGroup
	rs := r.manager.Conn(ctx).Raw(`
		select g.*
		from app_group g
		join user_group ug on ug.tenant_codigo = g.tenant_codigo and ug.group_code = g.code
		where ug.tenant_codigo = ? and ug.username = ?
		order by g.code
	`, tenantCodigo, username).Scan(&groups)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return groups, nil
}

func (r *AppUserRepository) ReplaceAppUserGroups(ctx context.Context, tenantCodigo string, username string, groupCodes []string, createdBy string) error {
	return r.manager.WithTx(ctx, func(ctx context.Context) error {
		tx := r.manager.Conn(ctx)
		if err := errs.Pgf(tx.Delete(&models.UserGroup{}, "tenant_codigo = ? and username = ?", tenantCodigo, username).Error); err != nil {
			return err
		}

		now := time.Now()
		for _, groupCode := range groupCodes {
			entry := &models.UserGroup{TenantCodigo: tenantCodigo, Username: username, GroupCode: groupCode, CreatedBy: createdBy, CreatedAt: now}
			if err := errs.Pgf(tx.Create(entry).Error); err != nil {
				return err
			}
		}

		return nil
	})
}
