package dtos

type PermissionResponse struct {
	Code        string  `json:"code"`
	Description *string `json:"description"`
	RolesCount  int64   `json:"rolesCount"`
	GroupsCount int64   `json:"groupsCount"`
	UsersCount  int64   `json:"usersCount"`
}
