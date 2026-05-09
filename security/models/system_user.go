package models

import (
	"errors"
	"time"

	"basesdk/errs"

	"golang.org/x/crypto/bcrypt"
)

type SystemAccount struct {
	Username     string `gorm:"primaryKey"`
	PasswordHash string
	Disabled     bool
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedBy    *string
	UpdatedAt    *time.Time
}

func (*SystemAccount) TableName() string {
	return "system_account"
}

func (u *SystemAccount) ValidatePassword(plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plain))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errs.BadRequestDirect("usuario o contraseña invalidos")
		}
		return errs.BadRequestDirect("no se pudo decodificar la session")
	}
	return nil
}

func (u *SystemAccount) ChangePassword(plain string) error {
	if plain == "" {
		return errs.BadRequestDirect("la contraseña no puede estar vacía")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return errs.BadRequestDirect("no se pudo generar la contraseña")
	}

	u.PasswordHash = string(hash)

	return nil
}
