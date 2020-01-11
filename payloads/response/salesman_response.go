package response

import "github.com/jacky-htg/inventory/models"

// SalesmanResponse json
type SalesmanResponse struct {
	ID      uint64          `json:"id"`
	Company CompanyResponse `json:"company"`
	Name    string          `json:"name"`
	Email   string          `json:"email"`
	Address string          `json:"address"`
	Hp      string          `json:"hp"`
}

// Transform Salesman models to salesman response
func (u *SalesmanResponse) Transform(c *models.Salesman) {
	u.ID = c.ID
	u.Name = c.Name
	u.Email = c.Email
	u.Address = c.Address
	u.Hp = c.Hp
	u.Company.Transform(&c.Company)
}
