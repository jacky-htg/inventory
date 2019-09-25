package response

import (
	"github.com/jacky-htg/inventory/models"
)

//BranchResponse : format json response for branch
type BranchResponse struct {
	ID      uint32          `json:"id"`
	Code    string          `json:"code"`
	Name    string          `json:"name"`
	Address string          `json:"address"`
	Type    string          `json:"type"`
	Company CompanyResponse `json:"company"`
}

//Transform from Branch model to Branch response
func (u *BranchResponse) Transform(branch *models.Branch) {
	u.ID = branch.ID
	u.Code = branch.Code
	u.Name = branch.Name
	u.Address = branch.Address.String
	u.Type = branch.Type
	u.Company.Transform(&branch.Company)
}
