package response

import (
	"github.com/jacky-htg/inventory/models"
)

// SupplierResponse : format json response for supplier
type SupplierResponse struct {
	ID      uint64          `json:"id"`
	Code    string          `json:"code"`
	Name    string          `json:"name"`
	Address string          `json:"address"`
	Company CompanyResponse `json:"company"`
}

// Transform from Supplier model to Supplier response
func (u *SupplierResponse) Transform(s *models.Supplier) {
	u.ID = s.ID
	u.Code = s.Code
	u.Name = s.Name
	u.Address = s.Address.String
	u.Company.Transform(&s.Company)
}
