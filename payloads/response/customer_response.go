package response

import "github.com/jacky-htg/inventory/models"

// CustomerResponse json
type CustomerResponse struct {
	ID      uint64          `json:"id"`
	Company CompanyResponse `json:"company"`
	Name    string          `json:"name"`
	Email   string          `json:"email"`
	Address string          `json:"address"`
	Hp      string          `json:"hp"`
}

// Transform Customer models to customer response
func (u *CustomerResponse) Transform(c *models.Customer) {
	u.ID = c.ID
	u.Name = c.Name
	u.Email = c.Email
	u.Address = c.Address
	u.Hp = c.Hp
	u.Company.Transform(&c.Company)
}
