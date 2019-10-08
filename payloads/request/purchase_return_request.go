package request

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// NewPurchaseReturnRequest : format json request for new purchase return
type NewPurchaseReturnRequest struct {
	Date                  string                           `json:"date" validate:"required"`
	AdditionalDisc        float64                          `json:"additional_disc"`
	PurchaseReturnDetails []NewPurchaseReturnDetailRequest `json:"purchase_return_details" validate:"required"`
	PurchaseID            uint64                           `json:"purchase" validate:"required"`
}

// Transform NewPurchaseReturnRequest to Purchase
func (u *NewPurchaseReturnRequest) Transform() *models.PurchaseReturn {
	var p models.PurchaseReturn
	p.Date, _ = time.Parse("2006-01-02", u.Date)
	p.Purchase.ID = u.PurchaseID
	p.AdditionalDisc = u.AdditionalDisc

	for _, pd := range u.PurchaseReturnDetails {
		if pd.Qty < 1 {
			pd.Qty = 1
		}

		p.PurchaseReturnDetails = append(p.PurchaseReturnDetails, pd.Transform())
	}

	return &p
}

// NewPurchaseReturnDetailRequest : format json request for purchase return detail
type NewPurchaseReturnDetailRequest struct {
	Price     float64 `json:"price"`
	Disc      float64 `json:"disc"`
	Qty       uint    `json:"qty" validate:"required"`
	ProductID uint64  `json:"product"`
}

// Transform NewPurchaseReturnDetailRequest to PurchaseReturnDetail
func (u *NewPurchaseReturnDetailRequest) Transform() models.PurchaseReturnDetail {
	var pd models.PurchaseReturnDetail
	pd.Price = u.Price
	pd.Disc = u.Disc
	pd.Qty = u.Qty
	pd.Product.ID = u.ProductID

	return pd
}

// PurchaseReturnRequest : format json request for purchase return
type PurchaseReturnRequest struct {
	ID                    uint64                        `json:"id" validate:"required"`
	Date                  string                        `json:"date"`
	AdditionalDisc        float64                       `json:"additional_disc"`
	PurchaseReturnDetails []PurchaseReturnDetailRequest `json:"purchase_return_details"`
	PurchaseID            uint64                        `json:"purchase"`
}

// Transform PurchaseReturnRequest to PurchaseReturn
func (u *PurchaseReturnRequest) Transform(p *models.PurchaseReturn) *models.PurchaseReturn {
	if u.ID == p.ID {
		p.Date, _ = time.Parse("2006-01-02", u.Date)
		p.Purchase.ID = u.PurchaseID
		p.AdditionalDisc = u.AdditionalDisc

		var details []models.PurchaseReturnDetail
		for _, pd := range u.PurchaseReturnDetails {
			if pd.Qty < 1 {
				pd.Qty = 1
			}

			details = append(details, pd.Transform())
		}

		p.PurchaseReturnDetails = details

	}
	return p
}

// PurchaseReturnDetailRequest : format json request for purchase return detail
type PurchaseReturnDetailRequest struct {
	ID        uint64  `json:"id"`
	Price     float64 `json:"price"`
	Disc      float64 `json:"disc"`
	Qty       uint    `json:"qty"`
	ProductID uint64  `json:"product"`
}

// Transform PurchaseReturnDetailRequest to PurchaseReturnDetail
func (u *PurchaseReturnDetailRequest) Transform() models.PurchaseReturnDetail {
	var pd models.PurchaseReturnDetail
	pd.ID = u.ID
	pd.Price = u.Price
	pd.Disc = u.Disc
	pd.Qty = u.Qty
	pd.Product.ID = u.ProductID

	return pd
}
