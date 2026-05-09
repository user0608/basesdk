package dtos

type CreateGroupInput struct {
	Code        string  `json:"code" chk:"nonil"`
	Description *string `json:"description"`
}

type UpdateGroupInput struct {
	Description *string `json:"description"`
	Disabled    bool    `json:"disabled"`
}

type GroupResponse struct {
	TenantCodigo string  `json:"tenantCodigo"`
	Code         string  `json:"code"`
	Description  *string `json:"description"`
	Disabled     bool    `json:"disabled"`
}
