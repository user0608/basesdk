package models

import "time"

type AppGroup struct {
	TenantCodigo     string `gorm:"primaryKey"`
	Code             string `gorm:"primaryKey"`
	Description      *string
	Disabled         bool
	UsersCount       int64 `gorm:"->"`
	RolesCount       int64 `gorm:"->"`
	PermissionsCount int64 `gorm:"->"`
	CreatedBy        string
	CreatedAt        time.Time
	UpdatedBy        *string
	UpdatedAt        *time.Time
}

func (*AppGroup) TableName() string {
	return "app_group"
}

type UserGroup struct {
	TenantCodigo string `gorm:"primaryKey"`
	Username     string `gorm:"primaryKey"`
	GroupCode    string `gorm:"primaryKey"`
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedBy    *string
	UpdatedAt    *time.Time
}

func (*UserGroup) TableName() string {
	return "user_group"
}

type GroupRole struct {
	TenantCodigo string `gorm:"primaryKey"`
	GroupCode    string `gorm:"primaryKey"`
	RoleCode     string `gorm:"primaryKey"`
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedBy    *string
	UpdatedAt    *time.Time
}

func (*GroupRole) TableName() string {
	return "group_role"
}
