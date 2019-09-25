package request

import (
	"database/sql"

	"github.com/jacky-htg/inventory/models"
)

//NewBranchRequest : format json request for new branch
type NewBranchRequest struct {
	Code    string `json:"code" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Type    string `json:"type" validate:"required"`
	Address string `json:"address,omitempty"`
}

//Transform NewBranchRequest to Branch
func (u *NewBranchRequest) Transform() *models.Branch {
	var branch models.Branch
	branch.Code = u.Code
	branch.Name = u.Name
	branch.Type = u.Type
	branch.Address = sql.NullString{Valid: true, String: u.Address}

	return &branch
}

//BranchRequest : format json request for branch
type BranchRequest struct {
	ID      uint32 `json:"id,omitempty" validate:"required"`
	Code    string `json:"code,omitempty" validate:"required"`
	Name    string `json:"name,omitempty" validate:"required"`
	Address string `json:"address,omitempty"`
}

//Transform BranchRequest to Branch
func (u *BranchRequest) Transform(branch *models.Branch) *models.Branch {
	if u.ID == branch.ID {
		if len(u.Code) > 0 {
			branch.Code = u.Code
		}

		if len(u.Name) > 0 {
			branch.Name = u.Name
		}

		if len(u.Address) > 0 {
			branch.Address = sql.NullString{Valid: true, String: u.Address}
		}
	}
	return branch
}
