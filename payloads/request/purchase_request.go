package request

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// NewPurchaseRequest : format json request for new purchase
type NewPurchaseRequest struct {
	Date            string                     `json:"date" validate:"required"`
	AdditionalDisc  float64                    `json:"additional_disc"`
	PurchaseDetails []NewPurchaseDetailRequest `json:"purchase_details" validate:"required"`
	SupplierID      uint64                     `json:"supplier" validate:"required"`
}

// Transform NewPurchaseRequest to Purchase
func (u *NewPurchaseRequest) Transform() *models.Purchase {
	var p models.Purchase
	p.Date, _ = time.Parse("2006-01-02", u.Date)
	p.Supplier.ID = u.SupplierID
	p.AdditionalDisc = u.AdditionalDisc

	for _, pd := range u.PurchaseDetails {
		if pd.Qty < 1 {
			pd.Qty = 1
		}

		p.PurchaseDetails = append(p.PurchaseDetails, pd.Transform())
	}

	return &p
}

// NewPurchaseDetailRequest : format json request for purchase detail
type NewPurchaseDetailRequest struct {
	Price     float64 `json:"price" validate:"required"`
	Disc      float64 `json:"disc"`
	Qty       uint    `json:"qty" validate:"required"`
	ProductID uint64  `json:"product" validate:"required"`
}

// Transform NewPurchaseDetailRequest to PurchaseDetail
func (u *NewPurchaseDetailRequest) Transform() models.PurchaseDetail {
	var pd models.PurchaseDetail
	pd.Price = u.Price
	pd.Disc = u.Disc
	pd.Qty = u.Qty
	pd.Product.ID = u.ProductID

	return pd
}

// PurchaseRequest : format json request for purchase
type PurchaseRequest struct {
	ID              uint64                  `json:"id" validate:"required"`
	Date            string                  `json:"date"`
	AdditionalDisc  float64                 `json:"additional_disc"`
	PurchaseDetails []PurchaseDetailRequest `json:"purchase_details"`
	SupplierID      uint64                  `json:"supplier"`
}

// Transform PurchaseRequest to Purchase
func (u *PurchaseRequest) Transform(p *models.Purchase) *models.Purchase {
	if u.ID == p.ID {
		p.Date, _ = time.Parse("2006-01-02", u.Date)
		p.Supplier.ID = u.SupplierID
		p.AdditionalDisc = u.AdditionalDisc

		var details []models.PurchaseDetail
		for _, pd := range u.PurchaseDetails {
			if pd.Qty < 1 {
				pd.Qty = 1
			}

			details = append(details, pd.Transform())
		}

		p.PurchaseDetails = details

	}
	return p
}

// PurchaseDetailRequest : format json request for purchase detail
type PurchaseDetailRequest struct {
	ID        uint64  `json:"id"`
	Price     float64 `json:"price"`
	Disc      float64 `json:"disc"`
	Qty       uint    `json:"qty"`
	ProductID uint64  `json:"product"`
}

// Transform PurchaseDetailRequest to PurchaseDetail
func (u *PurchaseDetailRequest) Transform() models.PurchaseDetail {
	var pd models.PurchaseDetail
	pd.ID = u.ID
	pd.Price = u.Price
	pd.Disc = u.Disc
	pd.Qty = u.Qty
	pd.Product.ID = u.ProductID

	return pd
}
