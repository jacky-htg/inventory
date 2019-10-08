package response

import "github.com/jacky-htg/inventory/models"

// ProductCategoryResponse json
type ProductCategoryResponse struct {
	ID       uint64           `json:"id"`
	Company  CompanyResponse  `json:"company"`
	Name     string           `json:"name"`
	Category CategoryResponse `json:"category"`
}

// Transform ProductCategory models to ProductCategory response
func (u *ProductCategoryResponse) Transform(c *models.ProductCategory) {
	u.ID = c.ID
	u.Name = c.Name
	u.Company.Transform(&c.Company)
	u.Category.Transform(&c.Category)
}

// CategoryResponse json
type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Transform Category models to Category response
func (u *CategoryResponse) Transform(c *models.Category) {
	u.ID = c.ID
	u.Name = c.Name
}
