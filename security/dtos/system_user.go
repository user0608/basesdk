package dtos

import (
	"basesdk/security/models"
	"time"
)

type SystemUserDTO struct {
	Username  string     `json:"username"`
	Disabled  bool       `json:"disabled"`
	CreatedBy string     `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedBy *string    `json:"updated_by,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func NewSystemUserDTO(user *models.SystemUser) *SystemUserDTO {
	if user == nil {
		return nil
	}

	return &SystemUserDTO{
		Username:  user.Username,
		Disabled:  user.Disabled,
		CreatedBy: user.CreatedBy,
		CreatedAt: user.CreatedAt,
		UpdatedBy: user.UpdatedBy,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewSystemUserDTOs(users []models.SystemUser) []SystemUserDTO {
	result := make([]SystemUserDTO, 0, len(users))

	for i := range users {
		dto := NewSystemUserDTO(&users[i])
		if dto != nil {
			result = append(result, *dto)
		}
	}

	return result
}
