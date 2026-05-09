package dtos

type CreateRoleInput struct {
	Code        string  `json:"code" chk:"nonil"`
	Description *string `json:"description"`
}

type UpdateRoleInput struct {
	Description *string `json:"description"`
	Disabled    bool    `json:"disabled"`
}

type RoleResponse struct {
	TenantCodigo string  `json:"tenantCodigo"`
	Code         string  `json:"code"`
	Description  *string `json:"description"`
	Disabled     bool    `json:"disabled"`
}
