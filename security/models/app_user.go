package models

import (
	"errors"
	"time"

	"basesdk/errs"

	"golang.org/x/crypto/bcrypt"
)

type Tenant struct {
	Codigo         string `gorm:"primaryKey"`
	Name           string
	Timezone       string
	MaxActiveUsers int
	Disabled       bool
	ExpiresAt      *time.Time
	UsersCount     int64 `gorm:"->"`
	RolesCount     int64 `gorm:"->"`
	GroupsCount    int64 `gorm:"->"`
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedBy      *string
	UpdatedAt      *time.Time
}

func (*Tenant) TableName() string {
	return "tenant"
}

func (t *Tenant) ValidateLoginAccess(now time.Time) error {
	if t.Disabled {
		return errs.BadRequestDirect("tenant deshabilitado")
	}

	if t.ExpiresAt != nil && now.After(*t.ExpiresAt) {
		return errs.BadRequestDirect("tenant expirado")
	}

	return nil
}

type AppUser struct {
	TenantCodigo       string `gorm:"primaryKey"`
	Username           string `gorm:"primaryKey"`
	FullName           *string
	PasswordHash       *string
	MustChangePassword bool
	LastLoginAt        *time.Time
	Disabled           bool
	RolesCount         int64 `gorm:"->"`
	GroupsCount        int64 `gorm:"->"`
	PermissionsCount   int64 `gorm:"->"`
	CreatedBy          string
	CreatedAt          time.Time
	UpdatedBy          *string
	UpdatedAt          *time.Time
}

func (*AppUser) TableName() string {
	return "app_user"
}

func (u *AppUser) ValidateLoginAccess() error {
	if u.Disabled {
		return errs.BadRequestDirect("usuario deshabilitado")
	}

	if u.PasswordHash == nil || *u.PasswordHash == "" {
		return errs.BadRequestDirect("usuario o contraseña invalidos")
	}

	return nil
}

func (u *AppUser) ValidatePassword(plain string) error {
	if err := u.ValidateLoginAccess(); err != nil {
		return err
	}

	err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(plain))
	if err != nil {
		return errs.BadRequestDirect("usuario o contraseña invalidos")
	}

	return nil
}

func (u *AppUser) ChangePassword(plain string) error {
	if plain == "" {
		return errs.BadRequestDirect("la contraseña no puede estar vacía")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return errs.BadRequestDirect("la contraseña es demasiado larga")
		}
		return errs.BadRequestDirect("no se pudo generar la contraseña")
	}

	passwordHash := string(hash)
	u.PasswordHash = &passwordHash

	return nil
}
