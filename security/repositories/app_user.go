package repositories

import (
	"basesdk/connection"
	"basesdk/errs"
	"basesdk/security/models"
	"context"
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
