package response

type RoleRes struct {
	RoleCode    string          `json:"role_code"`
	Name        string          `json:"name"`
	Description string          `json:"description" `
	Status      string          `json:"status"`
	Permissions []PermissionRes `json:"permissions"`
}
