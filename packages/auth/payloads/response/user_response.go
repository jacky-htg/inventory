package response

import (
	"github.com/jacky-htg/inventory/packages/auth/models"
)

//UserResponse : format json response for user
type UserResponse struct {
	ID       uint64        `json:"id"`
	Username string        `json:"username"`
	IsActive bool          `json:"is_active"`
	Roles    []models.Role `json:"roles"`
}

//Transform from User model to User response
func (u *UserResponse) Transform(user *models.User) {
	u.ID = user.ID
	u.Username = user.Username
	u.IsActive = user.IsActive
	u.Roles = user.Roles
}
