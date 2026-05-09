package repositories

import (
	"basesdk/connection"
	"basesdk/errs"
	"basesdk/security/models"
	"context"
)

type SystemUserRepository struct {
	manager connection.StorageManager
}

func NewSystemUserRepository(manager connection.StorageManager) *SystemUserRepository {
	return &SystemUserRepository{
		manager: manager,
	}
}

func (r *SystemUserRepository) FindSystemUser(ctx context.Context, username string) (*models.SystemAccount, error) {
	tx := r.manager.Conn(ctx)
	var user models.SystemAccount

	rs := tx.Where("username = ?", username).First(&user)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return &user, nil
}
