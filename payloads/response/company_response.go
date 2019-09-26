package response

import (
	"github.com/jacky-htg/inventory/models"
)

//CompanyResponse : format json response for company
type CompanyResponse struct {
	ID      uint32 `json:"id"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

//Transform from Company model to Company response
func (u *CompanyResponse) Transform(company *models.Company) {
	u.ID = company.ID
	u.Name = company.Name
	u.Code = company.Code
	u.Address = company.Address.String
}
