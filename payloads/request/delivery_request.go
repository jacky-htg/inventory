package request

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// NewDeliveryRequest : format json request for new Delivery
type NewDeliveryRequest struct {
	Date            string                     `json:"date" validate:"required"`
	Remark          string                     `json:"remark"`
	DeliveryDetails []NewDeliveryDetailRequest `json:"delivery_details" validate:"required"`
	SalesOrderID    uint64                     `json:"sales_order" validate:"required"`
}

// Transform NewDeliveryRequest to Delivery
func (u *NewDeliveryRequest) Transform() *models.Delivery {
	var p models.Delivery
	p.Date, _ = time.Parse("2006-01-02", u.Date)
	p.SalesOrder.ID = u.SalesOrderID
	p.Remark = u.Remark

	for _, pd := range u.DeliveryDetails {
		p.DeliveryDetails = append(p.DeliveryDetails, pd.Transform())
	}

	return &p
}

// NewDeliveryDetailRequest : format json request for Delivery detail
type NewDeliveryDetailRequest struct {
	ProductID uint64 `json:"product" validate:"required"`
	ShelveID  uint64 `json:"shelve" validate:"required"`
}

// Transform NewDeliveryDetailRequest to DeliveryDetail
func (u *NewDeliveryDetailRequest) Transform() models.DeliveryDetail {
	var pd models.DeliveryDetail
	pd.Qty = 1
	pd.Product.ID = u.ProductID
	pd.Shelve.ID = u.ShelveID

	return pd
}

// DeliveryRequest : format json request for Delivery
type DeliveryRequest struct {
	ID              uint64                  `json:"id" validate:"required"`
	Date            string                  `json:"date"`
	Remark          string                  `json:"remark"`
	DeliveryDetails []DeliveryDetailRequest `json:"delivery_details"`
	SalesOrderID    uint64                  `json:"sales_order"`
}

// Transform DeliveryRequest to Delivery
func (u *DeliveryRequest) Transform(p *models.Delivery) *models.Delivery {
	if u.ID == p.ID {
		p.Date, _ = time.Parse("2006-01-02", u.Date)
		p.SalesOrder.ID = u.SalesOrderID
		p.Remark = u.Remark

		var details []models.DeliveryDetail
		for _, pd := range u.DeliveryDetails {
			details = append(details, pd.Transform())
		}

		p.DeliveryDetails = details

	}
	return p
}

// DeliveryDetailRequest : format json request for Delivery detail
type DeliveryDetailRequest struct {
	ID        uint64 `json:"id"`
	ProductID uint64 `json:"product"`
	ShelveID  uint64 `json:"shelve"`
}

// Transform DeliveryDetailRequest to DeliveryDetail
func (u *DeliveryDetailRequest) Transform() models.DeliveryDetail {
	var pd models.DeliveryDetail
	pd.ID = u.ID
	pd.Qty = 1
	pd.Product.ID = u.ProductID
	pd.Shelve.ID = u.ShelveID

	return pd
}
