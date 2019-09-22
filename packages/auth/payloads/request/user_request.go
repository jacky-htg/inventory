package request

import (
	"github.com/jacky-htg/inventory/packages/auth/models"
)

//NewUserRequest : format json request for new user
type NewUserRequest struct {
	Username   string        `json:"username" validate:"required"`
	Email      string        `json:"email" validate:"required"`
	Password   string        `json:"password" validate:"required"`
	RePassword string        `json:"re_password" validate:"required"`
	IsActive   bool          `json:"is_active"`
	Roles      []models.Role `json:"roles"`
}

//Transform NewUserRequest to User
func (u *NewUserRequest) Transform() *models.User {
	var user models.User
	user.Username = u.Username
	user.Email = u.Email
	user.Password = u.Password
	user.IsActive = u.IsActive
	user.Roles = u.Roles

	return &user
}

//UserRequest : format json request for user
type UserRequest struct {
	ID         uint64        `json:"id,omitempty" validate:"required"`
	Username   string        `json:"username,omitempty" validate:"required"`
	Email      string        `json:"email,omitempty" validate:"required"`
	Password   string        `json:"password,omitempty" validate:"required"`
	RePassword string        `json:"re_password,omitempty" validate:"required"`
	IsActive   bool          `json:"is_active,omitempty"`
	Roles      []models.Role `json:"roles,omitempty"`
}

//Transform NewUserRequest to User
func (u *UserRequest) Transform(user *models.User) *models.User {
	if u.ID == user.ID {
		if len(u.Username) > 0 {
			user.Username = u.Username
		}

		if len(u.Email) > 0 {
			user.Email = u.Email
		}

		if len(u.Password) > 0 {
			user.Password = u.Password
		}

		if len(u.Roles) > 0 {
			user.Roles = u.Roles
		}

		user.IsActive = u.IsActive
	}
	return user
}
