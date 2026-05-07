package repositories

import (
	"basesdk/security/models"
	"context"
)

type SystemUserRepository interface {
	FindSystemUser(ctx context.Context, username string) (*models.SystemUser, error)
	FindSystemUsers(ctx context.Context) ([]models.SystemUser, error)

	CreateSystemUser(ctx context.Context, user *models.SystemUser) error
	UpdateSystemUser(ctx context.Context, user *models.SystemUser) error
	DeleteSystemUser(ctx context.Context, username string) error

	ExistsSystemUser(ctx context.Context, username string) (bool, error)

	ChangeSystemUserPassword(ctx context.Context, username string, passwordHash string) error

	EnableSystemUser(ctx context.Context, username string) error
	DisableSystemUser(ctx context.Context, username string) error
}
