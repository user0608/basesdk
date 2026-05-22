package models

import "time"

type Role struct {
	TenantCodigo     string `gorm:"primaryKey"`
	Code             string `gorm:"primaryKey"`
	Description      *string
	Disabled         bool
	UsersCount       int64 `gorm:"->"`
	GroupsCount      int64 `gorm:"->"`
	PermissionsCount int64 `gorm:"->"`
	CreatedBy        string
	CreatedAt        time.Time
	UpdatedBy        *string
	UpdatedAt        *time.Time
}

func (*Role) TableName() string {
	return "role"
}

type UserRole struct {
	TenantCodigo string `gorm:"primaryKey"`
	Username     string `gorm:"primaryKey"`
	RoleCode     string `gorm:"primaryKey"`
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedBy    *string
	UpdatedAt    *time.Time
}

func (*UserRole) TableName() string {
	return "user_role"
}

type RolePermission struct {
	TenantCodigo   string `gorm:"primaryKey"`
	RoleCode       string `gorm:"primaryKey"`
	PermissionCode string `gorm:"primaryKey"`
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedBy      *string
	UpdatedAt      *time.Time
}

func (*RolePermission) TableName() string {
	return "role_permission"
}
