package request

import (
	"github.com/jacky-htg/inventory/models"
)

//NewRegionRequest : format json request for new region
type NewRegionRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

//Transform NewRegionRequest to Region
func (u *NewRegionRequest) Transform() *models.Region {
	var region models.Region
	region.Code = u.Code
	region.Name = u.Name

	return &region
}

//RegionRequest : format json request for region
type RegionRequest struct {
	ID   uint32 `json:"id,omitempty" validate:"required"`
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

//Transform RegionRequest to Region
func (u *RegionRequest) Transform(region *models.Region) *models.Region {
	if u.ID == region.ID {
		if len(u.Code) > 0 {
			region.Code = u.Code
		}

		if len(u.Name) > 0 {
			region.Name = u.Name
		}
	}
	return region
}
