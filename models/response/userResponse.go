package response

import (
	"go-rest-api/models"
	"time"
)

type UserLoginResponse struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type UserInfoResponse struct {
	UserId      uint                 `json:"user_id"`
	Username    string               `json:"username"`
	RoleId      uint                 `json:"role_id"`
	RoleName    string               `json:"role_name"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

func (userResponse UserInfoResponse) ToFormat(user *models.User) UserInfoResponse {
	userResponse.UserId = user.ID
	userResponse.Username = user.Username
	userResponse.RoleId = user.RoleID
	userResponse.RoleName = user.Role.Name
	userResponse.Permissions = PermissionResponse{}.ToFormat(user.Role.Permissions)
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	return userResponse

}
