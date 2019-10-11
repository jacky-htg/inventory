package response

import (
	"github.com/jacky-htg/inventory/models"
)

// ShelveResponse : format json response for shelve
type ShelveResponse struct {
	ID       uint64 `json:"id"`
	Code     string `json:"code"`
	Capacity uint   `json:"capacity"`
}

//Transform from Shelve model to Shelve response
func (u *ShelveResponse) Transform(shelve *models.Shelve) {
	u.ID = shelve.ID
	u.Code = shelve.Code
	u.Capacity = shelve.Capacity
}
