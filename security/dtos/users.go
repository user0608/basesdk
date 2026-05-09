package dtos

import "time"

type CreateSystemUserInput struct {
	Username string `json:"username" chk:"nonil"`
	Password string `json:"password" chk:"nonil"`
}

type UpdateSystemUserInput struct {
	Disabled bool `json:"disabled"`
}

type ChangePasswordInput struct {
	Password string `json:"password" chk:"nonil"`
}

type SystemUserResponse struct {
	Username  string     `json:"username"`
	Disabled  bool       `json:"disabled"`
	CreatedBy string     `json:"createdBy"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedBy *string    `json:"updatedBy"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type CreateTenantUserInput struct {
	Username           string  `json:"username" chk:"nonil"`
	Email              string  `json:"email" chk:"nonil"`
	FullName           *string `json:"fullName"`
	Password           string  `json:"password" chk:"nonil"`
	EmailVerified      bool    `json:"emailVerified"`
	MustChangePassword bool    `json:"mustChangePassword"`
}

type UpdateTenantUserInput struct {
	Email              string  `json:"email" chk:"nonil"`
	FullName           *string `json:"fullName"`
	EmailVerified      bool    `json:"emailVerified"`
	MustChangePassword bool    `json:"mustChangePassword"`
	Disabled           bool    `json:"disabled"`
}

type TenantUserResponse struct {
	TenantCodigo       string     `json:"tenantCodigo"`
	Username           string     `json:"username"`
	Email              string     `json:"email"`
	FullName           *string    `json:"fullName"`
	EmailVerified      bool       `json:"emailVerified"`
	MustChangePassword bool       `json:"mustChangePassword"`
	LastLoginAt        *time.Time `json:"lastLoginAt"`
	Disabled           bool       `json:"disabled"`
}

type UserPermissionsResponse struct {
	TenantCodigo string               `json:"tenantCodigo"`
	Username     string               `json:"username"`
	Permissions  []PermissionResponse `json:"permissions"`
}
