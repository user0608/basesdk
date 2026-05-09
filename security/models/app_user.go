package models

import (
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
	Email              string
	FullName           *string
	PasswordHash       *string
	EmailVerified      bool
	MustChangePassword bool
	LastLoginAt        *time.Time
	Disabled           bool
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
