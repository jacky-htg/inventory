package request

import "github.com/jacky-htg/inventory/models"

// NewCustomerRequest is json request for new customer and validation
type NewCustomerRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required"`
	Address string `json:"address" validate:"required"`
	Hp      string `json:"hp" validate:"required"`
}

// Transform NewCustomerRequest to Customer model
func (u *NewCustomerRequest) Transform() models.Customer {
	var c models.Customer
	c.Name = u.Name
	c.Email = u.Email
	c.Address = u.Address
	c.Hp = u.Hp

	return c
}

// CustomerRequest is json request for update customer and validation
type CustomerRequest struct {
	ID      uint64 `json:"id" validate:"required"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Hp      string `json:"hp"`
}

// Transform CustomerRequest to Customer model
func (u *CustomerRequest) Transform(c *models.Customer) *models.Customer {
	if c.ID == u.ID {
		if len(u.Name) > 0 {
			c.Name = u.Name
		}
		if len(u.Email) > 0 {
			c.Email = u.Email
		}
		if len(u.Address) > 0 {
			c.Address = u.Address
		}
		if len(u.Hp) > 0 {
			c.Hp = u.Hp
		}
	}
	return c
}
