package dtos

import "time"

type CreateTenantInput struct {
	Codigo         string     `json:"codigo" chk:"nonil"`
	Name           string     `json:"name" chk:"nonil"`
	Timezone       string     `json:"timezone" chk:"nonil"`
	MaxActiveUsers int        `json:"maxActiveUsers"`
	ExpiresAt      *time.Time `json:"expiresAt"`
}

type UpdateTenantInput struct {
	Name           string     `json:"name" chk:"nonil"`
	Timezone       string     `json:"timezone" chk:"nonil"`
	MaxActiveUsers int        `json:"maxActiveUsers"`
	Disabled       bool       `json:"disabled"`
	ExpiresAt      *time.Time `json:"expiresAt"`
}

type TenantResponse struct {
	Codigo         string     `json:"codigo"`
	Name           string     `json:"name"`
	Timezone       string     `json:"timezone"`
	MaxActiveUsers int        `json:"maxActiveUsers"`
	Disabled       bool       `json:"disabled"`
	ExpiresAt      *time.Time `json:"expiresAt"`
	UsersCount     int64      `json:"usersCount"`
	RolesCount     int64      `json:"rolesCount"`
	GroupsCount    int64      `json:"groupsCount"`
	CreatedBy      string     `json:"createdBy"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedBy      *string    `json:"updatedBy"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}
