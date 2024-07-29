package response

import "go-rest-api/models"

type PermissionResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Method string `json:"method"`
	Route  string `json:"route"`
}

func (permResponse PermissionResponse) ToFormat(permissions []models.Permission) []PermissionResponse {
	var responses []PermissionResponse
	for _, perm := range permissions {
		responses = append(responses, PermissionResponse{
			ID:     perm.ID,
			Name:   perm.Name,
			Method: perm.Method,
			Route:  perm.Route,
		})
	}
	return responses
}
