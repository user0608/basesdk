package repositories

import (
	"basesdk/connection"
	"basesdk/errs"
	"basesdk/security/models"
	"context"
	"time"
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

func (r *SystemUserRepository) FindSystemUsers(ctx context.Context) ([]models.SystemAccount, error) {
	tx := r.manager.Conn(ctx)
	var users []models.SystemAccount

	rs := tx.Order("username").Find(&users)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return users, nil
}

func (r *SystemUserRepository) CreateSystemUser(ctx context.Context, user *models.SystemAccount) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Create(user)
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}

func (r *SystemUserRepository) UpdateSystemUser(ctx context.Context, user *models.SystemAccount) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Model(&models.SystemAccount{}).
		Where("username = ?", user.Username).
		Updates(map[string]any{
			"disabled":   user.Disabled,
			"updated_by": user.UpdatedBy,
			"updated_at": user.UpdatedAt,
		})
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return errs.NotFoundDirect("usuario no encontrado")
	}

	return nil
}

func (r *SystemUserRepository) SetSystemUsersDisabled(ctx context.Context, usernames []string, disabled bool, updatedBy string) error {
	tx := r.manager.Conn(ctx)
	now := time.Now()

	rs := tx.Model(&models.SystemAccount{}).
		Where("username in ?", usernames).
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

func (r *SystemUserRepository) DeleteSystemUsers(ctx context.Context, usernames []string) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Delete(&models.SystemAccount{}, "username in ?", usernames)
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}

func (r *SystemUserRepository) CountActiveSystemUsersExcluding(ctx context.Context, usernames []string) (int64, error) {
	tx := r.manager.Conn(ctx)
	var count int64

	query := tx.Model(&models.SystemAccount{}).Where("disabled = false")
	if len(usernames) > 0 {
		query = query.Where("username not in ?", usernames)
	}

	rs := query.Count(&count)
	if rs.Error != nil {
		return 0, errs.Pgf(rs.Error)
	}

	return count, nil
}
