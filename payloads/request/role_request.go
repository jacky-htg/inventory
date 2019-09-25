package request

import (
	"github.com/jacky-htg/inventory/models"
)

//NewRoleRequest : format json request for new role
type NewRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

//Transform NewRoleRequest to Role
func (u *NewRoleRequest) Transform() *models.Role {
	var role models.Role
	role.Name = u.Name

	return &role
}

//RoleRequest : format json request for role
type RoleRequest struct {
	ID   uint32 `json:"id,omitempty"  validate:"required"`
	Name string `json:"name,omitempty"  validate:"required"`
}

//Transform RoleRequest to Role
func (u *RoleRequest) Transform(role *models.Role) *models.Role {
	if u.ID == role.ID {
		if len(u.Name) > 0 {
			role.Name = u.Name
		}
	}
	return role
}
