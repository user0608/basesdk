package postgres

import (
	"basesdk/connection"
	"basesdk/security/models"
	"basesdk/security/repositories"
	"context"

	"basesdk/errs"
)

type SystemUserRepository struct {
	manager connection.StorageManager
}

var _ repositories.SystemUserRepository = (*SystemUserRepository)(nil)

func NewSystemUserRepository(manager connection.StorageManager) *SystemUserRepository {
	return &SystemUserRepository{
		manager: manager,
	}
}

func (r *SystemUserRepository) FindSystemUser(ctx context.Context, username string) (*models.SystemUser, error) {
	tx := r.manager.Conn(ctx)
	var user models.SystemUser

	rs := tx.Where("username = ?", username).First(&user)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return &user, nil
}

func (r *SystemUserRepository) FindSystemUsers(ctx context.Context) ([]models.SystemUser, error) {
	tx := r.manager.Conn(ctx)
	var users []models.SystemUser

	rs := tx.Order("username ASC").Find(&users)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return users, nil
}

func (r *SystemUserRepository) CreateSystemUser(ctx context.Context, user *models.SystemUser) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Create(user)
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}

func (r *SystemUserRepository) UpdateSystemUser(ctx context.Context, user *models.SystemUser) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Model(&models.SystemUser{}).
		Where("username = ?", user.Username).
		Select("*").
		Updates(user)

	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	if rs.RowsAffected == 0 {
		return errs.NotFoundDirect("no se pudo actualizar el usuario porque no existe")
	}

	return nil
}

func (r *SystemUserRepository) DeleteSystemUser(ctx context.Context, username string) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Delete(&models.SystemUser{}, "username = ?", username)
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}

func (r *SystemUserRepository) ExistsSystemUser(ctx context.Context, username string) (bool, error) {
	tx := r.manager.Conn(ctx)
	var count int64

	rs := tx.Model(&models.SystemUser{}).
		Where("username = ?", username).
		Count(&count)

	if rs.Error != nil {
		return false, errs.Pgf(rs.Error)
	}

	return count > 0, nil
}

func (r *SystemUserRepository) ChangeSystemUserPassword(ctx context.Context, username string, passwordHash string) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Model(&models.SystemUser{}).
		Where("username = ?", username).
		Update("password_hash", passwordHash)

	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	if rs.RowsAffected == 0 {
		return errs.NotFoundDirect("no se pudo cambiar la contraseña porque el usuario no existe")
	}

	return nil
}

func (r *SystemUserRepository) EnableSystemUser(ctx context.Context, username string) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Model(&models.SystemUser{}).
		Where("username = ?", username).
		Update("disabled", false)

	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	if rs.RowsAffected == 0 {
		return errs.NotFoundDirect("no se pudo habilitar el usuario porque no existe")
	}

	return nil
}

func (r *SystemUserRepository) DisableSystemUser(ctx context.Context, username string) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Model(&models.SystemUser{}).
		Where("username = ?", username).
		Update("disabled", true)

	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	if rs.RowsAffected == 0 {
		return errs.NotFoundDirect("no se pudo deshabilitar el usuario porque no existe")
	}

	return nil
}
