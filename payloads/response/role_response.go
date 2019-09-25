package response

import (
	"github.com/jacky-htg/inventory/models"
)

//RoleResponse : format json response for role
type RoleResponse struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

//Transform from Role model to Role response
func (u *RoleResponse) Transform(role *models.Role) {
	u.ID = role.ID
	u.Name = role.Name
}
