package request

import (
	"github.com/jacky-htg/inventory/models"
)

// NewProductRequest : format json request for new product
type NewProductRequest struct {
	Code      string  `json:"code" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	SalePrice float64 `json:"price" validate:"required"`
}

// Transform NewProductRequest to Product
func (u *NewProductRequest) Transform() *models.Product {
	var product models.Product
	product.Code = u.Code
	product.Name = u.Name
	product.SalePrice = u.SalePrice

	return &product
}

// ProductRequest : format json request for product
type ProductRequest struct {
	ID        uint64  `json:"id,omitempty" validate:"required"`
	Code      string  `json:"code,omitempty"`
	Name      string  `json:"name,omitempty"`
	SalePrice float64 `json:"price,omitempty"`
}

// Transform ProductRequest to Product
func (u *ProductRequest) Transform(product *models.Product) *models.Product {
	if u.ID == product.ID {
		if len(u.Code) > 0 {
			product.Code = u.Code
		}

		if len(u.Name) > 0 {
			product.Name = u.Name
		}

		if u.SalePrice > 0 {
			product.SalePrice = u.SalePrice
		}
	}
	return product
}
