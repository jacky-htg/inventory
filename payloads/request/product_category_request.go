package request

import (
	"strconv"

	"github.com/jacky-htg/inventory/models"
)

// NewProductCategoryRequest is json request for new ProductCategory and validation
type NewProductCategoryRequest struct {
	Name       string `json:"name" validate:"required"`
	CategoryID uint   `json:"category" validate:"required"`
}

// Transform NewProductCategoryRequest to ProductCategory model
func (u *NewProductCategoryRequest) Transform() models.ProductCategory {
	var c models.ProductCategory
	c.Name = u.Name
	c.Category.ID = u.CategoryID

	return c
}

// ProductCategoryRequest is json request for update ProductCategory and validation
type ProductCategoryRequest struct {
	ID         uint64 `json:"id" validate:"required"`
	Name       string `json:"name"`
	CategoryID string `json:"category"`
}

// Transform ProductCategoryRequest to ProductCategory model
func (u *ProductCategoryRequest) Transform(c *models.ProductCategory) *models.ProductCategory {
	if c.ID == u.ID {
		if len(u.Name) > 0 {
			c.Name = u.Name
		}
		if len(u.CategoryID) > 0 {
			categoryID, _ := strconv.Atoi(u.CategoryID)
			c.Category.ID = uint(categoryID)
		}
	}
	return c
}
