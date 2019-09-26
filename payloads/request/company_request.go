package request

import (
	"database/sql"

	"github.com/jacky-htg/inventory/models"
)

//NewCompanyRequest : format json request for new company
type NewCompanyRequest struct {
	Code    string `json:"code" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address,omitempty"`
}

//Transform NewCompanyRequest to Company
func (u *NewCompanyRequest) Transform() *models.Company {
	var company models.Company
	company.Code = u.Code
	company.Name = u.Name
	if len(u.Address) > 0 {
		company.Address = sql.NullString{Valid: true, String: u.Address}
	}
	return &company
}

//CompanyRequest : format json request for company
type CompanyRequest struct {
	ID      uint32 `json:"id,omitempty" validate:"required"`
	Code    string `json:"code,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

//Transform CompanyRequest to Company
func (u *CompanyRequest) Transform(company *models.Company) *models.Company {
	if u.ID == company.ID {
		if len(u.Code) > 0 {
			company.Code = u.Code
		}

		if len(u.Name) > 0 {
			company.Name = u.Name
		}

		if len(u.Address) > 0 {
			company.Address = sql.NullString{Valid: true, String: u.Address}
		}
	}
	return company
}
