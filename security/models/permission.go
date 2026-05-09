package models

type Permission struct {
	Code        string `gorm:"primaryKey"`
	Description *string
}

func (*Permission) TableName() string {
	return "permission"
}
