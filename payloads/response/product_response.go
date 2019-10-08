package response

import (
	"github.com/jacky-htg/inventory/models"
)

// ProductResponse : format json response for product
type ProductResponse struct {
	ID              uint64                  `json:"id"`
	Code            string                  `json:"code"`
	Name            string                  `json:"name"`
	SalePrice       float64                 `json:"price"`
	MinimumStock    uint                    `json:"minimum_stock"`
	Company         CompanyResponse         `json:"company"`
	Brand           BrandResponse           `json:"brand"`
	ProductCategory ProductCategoryResponse `json:"product_category"`
}

// Transform from Product model to Product response
func (u *ProductResponse) Transform(product *models.Product) {
	u.ID = product.ID
	u.Code = product.Code
	u.Name = product.Name
	u.SalePrice = product.SalePrice
	u.MinimumStock = product.MinimumStock
	u.Company.Transform(&product.Company)
	u.Brand.Transform(&product.Brand)
	u.ProductCategory.Transform(&product.ProductCategory)
}
