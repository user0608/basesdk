package dtos

import (
	"basesdk/security/models"
	"time"
)

type SystemAccountDTO struct {
	Username  string     `json:"username"`
	Disabled  bool       `json:"disabled"`
	CreatedBy string     `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedBy *string    `json:"updated_by,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func NewSystemUserDTO(user *models.SystemAccount) *SystemAccountDTO {
	if user == nil {
		return nil
	}

	return &SystemAccountDTO{
		Username:  user.Username,
		Disabled:  user.Disabled,
		CreatedBy: user.CreatedBy,
		CreatedAt: user.CreatedAt,
		UpdatedBy: user.UpdatedBy,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewSystemUserDTOs(users []models.SystemAccount) []SystemAccountDTO {
	result := make([]SystemAccountDTO, 0, len(users))

	for i := range users {
		dto := NewSystemUserDTO(&users[i])
		if dto != nil {
			result = append(result, *dto)
		}
	}

	return result
}
