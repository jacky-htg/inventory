package request

import (
	"strconv"

	"github.com/jacky-htg/inventory/models"
)

// NewProductRequest : format json request for new product
type NewProductRequest struct {
	Code              string  `json:"code" validate:"required"`
	Name              string  `json:"name" validate:"required"`
	SalePrice         float64 `json:"price" validate:"required"`
	MinimumStock      string  `json:"minimum_stock" validate:"required"`
	BrandID           string  `json:"brand" validate:"required"`
	ProductCategoryID string  `json:"product_category" validate:"required"`
}

// Transform NewProductRequest to Product
func (u *NewProductRequest) Transform() *models.Product {
	var product models.Product
	product.Code = u.Code
	product.Name = u.Name
	product.SalePrice = u.SalePrice

	minStock, _ := strconv.Atoi(u.MinimumStock)
	product.MinimumStock = uint(minStock)

	brandID, _ := strconv.Atoi(u.BrandID)
	product.Brand.ID = uint64(brandID)

	productCategoryID, _ := strconv.Atoi(u.ProductCategoryID)
	product.ProductCategory.ID = uint64(productCategoryID)

	return &product
}

// ProductRequest : format json request for product
type ProductRequest struct {
	ID                uint64  `json:"id,omitempty" validate:"required"`
	Code              string  `json:"code,omitempty"`
	Name              string  `json:"name,omitempty"`
	SalePrice         float64 `json:"price,omitempty"`
	MinimumStock      string  `json:"minimum_stock,omitempty"`
	BrandID           string  `json:"brand"`
	ProductCategoryID string  `json:"product_category"`
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

		if len(u.MinimumStock) > 0 {
			minStock, _ := strconv.Atoi(u.MinimumStock)
			product.MinimumStock = uint(minStock)
		}

		if len(u.BrandID) > 0 {
			brandID, _ := strconv.Atoi(u.BrandID)
			product.Brand.ID = uint64(brandID)
		}

		if len(u.ProductCategoryID) > 0 {
			productCategoryID, _ := strconv.Atoi(u.ProductCategoryID)
			product.ProductCategory.ID = uint64(productCategoryID)
		}
	}
	return product
}
