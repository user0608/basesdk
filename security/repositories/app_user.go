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
	tx := r.manager.Conn(ctx)
	var user models.AppUser

	rs := tx.Where("tenant_codigo = ? and username = ?", tenantCodigo, username).First(&user)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return &user, nil
}

func (r *AppUserRepository) FindAppUsers(ctx context.Context, tenantCodigo string) ([]models.AppUser, error) {
	tx := r.manager.Conn(ctx)
	var users []models.AppUser

	rs := tx.Where("tenant_codigo = ?", tenantCodigo).Order("username").Find(&users)
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
			"email":                user.Email,
			"full_name":            user.FullName,
			"email_verified":       user.EmailVerified,
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
