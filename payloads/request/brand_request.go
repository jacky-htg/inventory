package request

import "github.com/jacky-htg/inventory/models"

// NewBrandRequest is json request for new Brand and validation
type NewBrandRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

// Transform NewBrandRequest to Brand model
func (u *NewBrandRequest) Transform() models.Brand {
	var c models.Brand
	c.Code = u.Code
	c.Name = u.Name

	return c
}

// BrandRequest is json request for update Brand and validation
type BrandRequest struct {
	ID   uint64 `json:"id" validate:"required"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// Transform BrandRequest to Brand model
func (u *BrandRequest) Transform(c *models.Brand) *models.Brand {
	if c.ID == u.ID {
		if len(u.Name) > 0 {
			c.Name = u.Name
		}
		if len(u.Code) > 0 {
			c.Code = u.Code
		}
	}
	return c
}
