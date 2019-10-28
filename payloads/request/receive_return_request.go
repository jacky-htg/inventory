package request

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// NewReceiveReturnRequest : format json request for new Receive return
type NewReceiveReturnRequest struct {
	Date                 string                          `json:"date" validate:"required"`
	Remark               string                          `json:"remark"`
	ReceiveReturnDetails []NewReceiveReturnDetailRequest `json:"receive_return_details" validate:"required"`
	ReceiveID            uint64                          `json:"receive" validate:"required"`
}

// Transform NewReceiveReturnRequest to Receive return
func (u *NewReceiveReturnRequest) Transform() *models.ReceiveReturn {
	var p models.ReceiveReturn
	p.Date, _ = time.Parse("2006-01-02", u.Date)
	p.Receive.ID = u.ReceiveID
	p.Remark = u.Remark

	for _, pd := range u.ReceiveReturnDetails {
		p.ReceiveReturnDetails = append(p.ReceiveReturnDetails, pd.Transform())
	}

	return &p
}

// NewReceiveReturnDetailRequest : format json request for Receive return detail
type NewReceiveReturnDetailRequest struct {
	ProductID uint64 `json:"product" validate:"required"`
	Code      string `json:"code" validate:"required"`
}

// Transform NewReceiveDetailRequest to ReceiveDetail
func (u *NewReceiveReturnDetailRequest) Transform() models.ReceiveReturnDetail {
	var pd models.ReceiveReturnDetail
	pd.Qty = 1
	pd.Product.ID = u.ProductID
	pd.Code = u.Code

	return pd
}

// ReceiveReturnRequest : format json request for Receive return
type ReceiveReturnRequest struct {
	ID                   uint64                       `json:"id" validate:"required"`
	Date                 string                       `json:"date"`
	Remark               string                       `json:"remark"`
	ReceiveReturnDetails []ReceiveReturnDetailRequest `json:"receive_return_details"`
	ReceiveID            uint64                       `json:"receive"`
}

// Transform ReceiveReturnRequest to ReceiveReturn
func (u *ReceiveReturnRequest) Transform(p *models.ReceiveReturn) *models.ReceiveReturn {
	if u.ID == p.ID {
		p.Date, _ = time.Parse("2006-01-02", u.Date)
		p.Receive.ID = u.ReceiveID
		p.Remark = u.Remark

		var details []models.ReceiveReturnDetail
		for _, pd := range u.ReceiveReturnDetails {
			details = append(details, pd.Transform())
		}

		p.ReceiveReturnDetails = details

	}
	return p
}

// ReceiveReturnDetailRequest : format json request for Receive return detail
type ReceiveReturnDetailRequest struct {
	ID        uint64 `json:"id"`
	ProductID uint64 `json:"product"`
	Code      string `json:"code"`
}

// Transform ReceiveReturnDetailRequest to ReceiveReturnDetail
func (u *ReceiveReturnDetailRequest) Transform() models.ReceiveReturnDetail {
	var pd models.ReceiveReturnDetail
	pd.ID = u.ID
	pd.Qty = 1
	pd.Product.ID = u.ProductID
	pd.Code = u.Code

	return pd
}
