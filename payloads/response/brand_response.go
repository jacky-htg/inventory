package response

import "github.com/jacky-htg/inventory/models"

// BrandResponse json
type BrandResponse struct {
	ID      uint64          `json:"id"`
	Company CompanyResponse `json:"company"`
	Code    string          `json:"code"`
	Name    string          `json:"name"`
}

// Transform Brand models to Brand response
func (u *BrandResponse) Transform(c *models.Brand) {
	u.ID = c.ID
	u.Code = c.Code
	u.Name = c.Name
	u.Company.Transform(&c.Company)
}
