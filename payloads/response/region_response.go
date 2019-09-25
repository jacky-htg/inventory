package response

import (
	"github.com/jacky-htg/inventory/models"
)

//RegionResponse : format json response for region
type RegionResponse struct {
	ID      uint32          `json:"id"`
	Code    string          `json:"code"`
	Name    string          `json:"name"`
	Company CompanyResponse `json:"company"`
}

//Transform from Region model to Region response
func (u *RegionResponse) Transform(region *models.Region) {
	u.ID = region.ID
	u.Code = region.Code
	u.Name = region.Name
	u.Company.Transform(&region.Company)
}
