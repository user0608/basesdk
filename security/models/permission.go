package models

type Permission struct {
	Code        string `gorm:"primaryKey"`
	Description *string
	RolesCount  int64 `gorm:"->"`
	GroupsCount int64 `gorm:"->"`
	UsersCount  int64 `gorm:"->"`
}

func (*Permission) TableName() string {
	return "permission"
}
