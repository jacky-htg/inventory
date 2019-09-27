package response

import (
	"github.com/jacky-htg/inventory/models"
)

// ProductResponse : format json response for product
type ProductResponse struct {
	ID        uint64          `json:"id"`
	Code      string          `json:"code"`
	Name      string          `json:"name"`
	SalePrice float64         `json:"price"`
	Company   CompanyResponse `json:"company"`
}

// Transform from Product model to Product response
func (u *ProductResponse) Transform(product *models.Product) {
	u.ID = product.ID
	u.Code = product.Code
	u.Name = product.Name
	u.SalePrice = product.SalePrice
	u.Company.Transform(&product.Company)
}
