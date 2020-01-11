package request

import "github.com/jacky-htg/inventory/models"

// NewSalesmanRequest is json request for new salesman and validation
type NewSalesmanRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required"`
	Address string `json:"address" validate:"required"`
	Hp      string `json:"hp" validate:"required"`
}

// Transform NewSalesmanRequest to Salesman model
func (u *NewSalesmanRequest) Transform() models.Salesman {
	var c models.Salesman
	c.Name = u.Name
	c.Email = u.Email
	c.Address = u.Address
	c.Hp = u.Hp

	return c
}

// SalesmanRequest is json request for update salesman and validation
type SalesmanRequest struct {
	ID      uint64 `json:"id" validate:"required"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Hp      string `json:"hp"`
}

// Transform SalesmanRequest to Salesman model
func (u *SalesmanRequest) Transform(c *models.Salesman) *models.Salesman {
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
