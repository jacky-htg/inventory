package request

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// NewReceiveRequest : format json request for new Receive
type NewReceiveRequest struct {
	Date           string                    `json:"date" validate:"required"`
	Remark         string                    `json:"remark"`
	ReceiveDetails []NewReceiveDetailRequest `json:"receive_details" validate:"required"`
	PurchaseID     uint64                    `json:"purchase" validate:"required"`
}

// Transform NewReceiveRequest to Receive
func (u *NewReceiveRequest) Transform() *models.Receive {
	var p models.Receive
	p.Date, _ = time.Parse("2006-01-02", u.Date)
	p.Purchase.ID = u.PurchaseID
	p.Remark = u.Remark

	for _, pd := range u.ReceiveDetails {
		p.ReceiveDetails = append(p.ReceiveDetails, pd.Transform())
	}

	return &p
}

// NewReceiveDetailRequest : format json request for Receive detail
type NewReceiveDetailRequest struct {
	ProductID uint64 `json:"product" validate:"required"`
	ShelveID  uint64 `json:"shelve" validate:"required"`
}

// Transform NewReceiveDetailRequest to ReceiveDetail
func (u *NewReceiveDetailRequest) Transform() models.ReceiveDetail {
	var pd models.ReceiveDetail
	pd.Qty = 1
	pd.Product.ID = u.ProductID
	pd.Shelve.ID = u.ShelveID

	return pd
}

// ReceiveRequest : format json request for Receive
type ReceiveRequest struct {
	ID             uint64                 `json:"id" validate:"required"`
	Date           string                 `json:"date"`
	Remark         string                 `json:"remark"`
	ReceiveDetails []ReceiveDetailRequest `json:"receive_details"`
	PurchaseID     uint64                 `json:"purchase"`
}

// Transform ReceiveRequest to Receive
func (u *ReceiveRequest) Transform(p *models.Receive) *models.Receive {
	if u.ID == p.ID {
		p.Date, _ = time.Parse("2006-01-02", u.Date)
		p.Purchase.ID = u.PurchaseID
		p.Remark = u.Remark

		var details []models.ReceiveDetail
		for _, pd := range u.ReceiveDetails {
			details = append(details, pd.Transform())
		}

		p.ReceiveDetails = details

	}
	return p
}

// ReceiveDetailRequest : format json request for Receive detail
type ReceiveDetailRequest struct {
	ID        uint64 `json:"id"`
	ProductID uint64 `json:"product"`
	ShelveID  uint64 `json:"shelve"`
}

// Transform ReceiveDetailRequest to ReceiveDetail
func (u *ReceiveDetailRequest) Transform() models.ReceiveDetail {
	var pd models.ReceiveDetail
	pd.ID = u.ID
	pd.Qty = 1
	pd.Product.ID = u.ProductID
	pd.Shelve.ID = u.ShelveID

	return pd
}
