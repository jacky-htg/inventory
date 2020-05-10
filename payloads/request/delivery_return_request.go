package request

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// NewDeliveryReturnRequest : format json request for new Delivery return
type NewDeliveryReturnRequest struct {
	Date                  string                           `json:"date" validate:"required"`
	Remark                string                           `json:"remark"`
	DeliveryReturnDetails []NewDeliveryReturnDetailRequest `json:"delivery_return_details" validate:"required"`
	DeliveryID            uint64                           `json:"delivery" validate:"required"`
}

// Transform NewDeliveryReturnRequest to Delivery return
func (u *NewDeliveryReturnRequest) Transform() *models.DeliveryReturn {
	var p models.DeliveryReturn
	p.Date, _ = time.Parse("2006-01-02", u.Date)
	p.Delivery.ID = u.DeliveryID
	p.Remark = u.Remark

	for _, pd := range u.DeliveryReturnDetails {
		p.DeliveryReturnDetails = append(p.DeliveryReturnDetails, pd.Transform())
	}

	return &p
}

// NewDeliveryReturnDetailRequest : format json request for Delivery return detail
type NewDeliveryReturnDetailRequest struct {
	ProductID uint64 `json:"product" validate:"required"`
	Code      string `json:"code" validate:"required"`
}

// Transform NewDeliveryDetailRequest to DeliveryDetail
func (u *NewDeliveryReturnDetailRequest) Transform() models.DeliveryReturnDetail {
	var pd models.DeliveryReturnDetail
	pd.Qty = 1
	pd.Product.ID = u.ProductID
	pd.Code = u.Code

	return pd
}

// DeliveryReturnRequest : format json request for Delivery return
type DeliveryReturnRequest struct {
	ID                    uint64                        `json:"id" validate:"required"`
	Date                  string                        `json:"date"`
	Remark                string                        `json:"remark"`
	DeliveryReturnDetails []DeliveryReturnDetailRequest `json:"delivery_return_details"`
	DeliveryID            uint64                        `json:"delivery"`
}

// Transform DeliveryReturnRequest to DeliveryReturn
func (u *DeliveryReturnRequest) Transform(p *models.DeliveryReturn) *models.DeliveryReturn {
	if u.ID == p.ID {
		p.Date, _ = time.Parse("2006-01-02", u.Date)
		p.Delivery.ID = u.DeliveryID
		p.Remark = u.Remark

		var details []models.DeliveryReturnDetail
		for _, pd := range u.DeliveryReturnDetails {
			details = append(details, pd.Transform())
		}

		p.DeliveryReturnDetails = details

	}
	return p
}

// DeliveryReturnDetailRequest : format json request for Delivery return detail
type DeliveryReturnDetailRequest struct {
	ID        uint64 `json:"id"`
	ProductID uint64 `json:"product"`
	Code      string `json:"code"`
}

// Transform DeliveryReturnDetailRequest to DeliveryReturnDetail
func (u *DeliveryReturnDetailRequest) Transform() models.DeliveryReturnDetail {
	var pd models.DeliveryReturnDetail
	pd.ID = u.ID
	pd.Qty = 1
	pd.Product.ID = u.ProductID
	pd.Code = u.Code

	return pd
}
