package request

import "github.com/jacky-htg/inventory/models"

// NewShelveRequest is json request for new Shelve
type NewShelveRequest struct {
	Code     string `json:"code" validate:"required"`
	Capacity uint   `json:"capacity" validate:"required"`
}

// Transform NewSehelveRequest to Shelve Model
func (r *NewShelveRequest) Transform() models.Shelve {
	var s models.Shelve
	s.Code = r.Code
	s.Capacity = r.Capacity
	return s
}

// ShelveRequest is json request for update
type ShelveRequest struct {
	ID       uint64 `json:"id" validate:"required"`
	Code     string `json:"code"`
	Capacity uint   `json:"capacity"`
}

// Transform ShelveRequest to Shelve Model
func (r *ShelveRequest) Transform(s *models.Shelve) *models.Shelve {
	if s.ID == r.ID {
		if len(r.Code) > 0 {
			s.Code = r.Code
		}

		if r.Capacity != 0 {
			s.Capacity = r.Capacity
		}
	}

	return s
}
