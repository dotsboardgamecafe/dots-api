package response

type PermissionRes struct {
	PermissionId   int64  `json:"permission_id"`
	PermissionCode string `json:"permission_code"`
	Name           string `json:"name"`
	RoutePattern   string `json:"route_pattern" `
	RouteMethod    string `json:"route_method" `
	Description    string `json:"description" `
	Status         string `json:"status"`
}
