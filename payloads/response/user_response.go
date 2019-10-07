package response

import (
	"github.com/jacky-htg/inventory/models"
)

//UserResponse : format json response for user
type UserResponse struct {
	ID       uint64          `json:"id"`
	Username string          `json:"username"`
	IsActive bool            `json:"is_active"`
	Roles    []models.Role   `json:"roles"`
	Company  CompanyResponse `json:"company"`
	Region   *RegionResponse `json:"region,omitempty"`
	Branch   *BranchResponse `json:"branch,omitempty"`
}

//Transform from User model to User response
func (u *UserResponse) Transform(user *models.User) {
	u.ID = user.ID
	u.Username = user.Username
	u.IsActive = user.IsActive
	u.Roles = user.Roles
	u.Company.Transform(&user.Company)

	if user.Region.ID > 0 {
		var regionResponse RegionResponse
		regionResponse.Transform(&user.Region)
		u.Region = &regionResponse
	}

	if user.Branch.ID > 0 {
		var branchResponse BranchResponse
		branchResponse.Transform(&user.Branch)
		u.Branch = &branchResponse
	}
}
