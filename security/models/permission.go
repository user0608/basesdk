package models

import "time"

type Permission struct {
	Code        string `gorm:"primaryKey"`
	Description *string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   *string
	UpdatedAt   *time.Time
}

func (*Permission) TableName() string {
	return "permission"
}
