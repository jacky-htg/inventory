package request

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// NewSalesOrderReturnRequest : format json request for new salesOrder return
type NewSalesOrderReturnRequest struct {
	Date                    string                             `json:"date" validate:"required"`
	AdditionalDisc          float64                            `json:"additional_disc"`
	SalesOrderReturnDetails []NewSalesOrderReturnDetailRequest `json:"sales_order_return_details" validate:"required"`
	SalesOrderID            uint64                             `json:"sales_order" validate:"required"`
}

// Transform NewSalesOrderReturnRequest to SalesOrder
func (u *NewSalesOrderReturnRequest) Transform() *models.SalesOrderReturn {
	var p models.SalesOrderReturn
	p.Date, _ = time.Parse("2006-01-02", u.Date)
	p.SalesOrder.ID = u.SalesOrderID
	p.AdditionalDisc = u.AdditionalDisc

	for _, pd := range u.SalesOrderReturnDetails {
		if pd.Qty < 1 {
			pd.Qty = 1
		}

		p.SalesOrderReturnDetails = append(p.SalesOrderReturnDetails, pd.Transform())
	}

	return &p
}

// NewSalesOrderReturnDetailRequest : format json request for salesOrder return detail
type NewSalesOrderReturnDetailRequest struct {
	Price     float64 `json:"price"`
	Disc      float64 `json:"disc"`
	Qty       uint    `json:"qty" validate:"required"`
	ProductID uint64  `json:"product"`
}

// Transform NewSalesOrderReturnDetailRequest to SalesOrderReturnDetail
func (u *NewSalesOrderReturnDetailRequest) Transform() models.SalesOrderReturnDetail {
	var pd models.SalesOrderReturnDetail
	pd.Price = u.Price
	pd.Disc = u.Disc
	pd.Qty = u.Qty
	pd.Product.ID = u.ProductID

	return pd
}

// SalesOrderReturnRequest : format json request for salesOrder return
type SalesOrderReturnRequest struct {
	ID                      uint64                          `json:"id" validate:"required"`
	Date                    string                          `json:"date"`
	AdditionalDisc          float64                         `json:"additional_disc"`
	SalesOrderReturnDetails []SalesOrderReturnDetailRequest `json:"sales_order_return_details"`
	SalesOrderID            uint64                          `json:"sales_order"`
}

// Transform SalesOrderReturnRequest to SalesOrderReturn
func (u *SalesOrderReturnRequest) Transform(p *models.SalesOrderReturn) *models.SalesOrderReturn {
	if u.ID == p.ID {
		p.Date, _ = time.Parse("2006-01-02", u.Date)
		p.SalesOrder.ID = u.SalesOrderID
		p.AdditionalDisc = u.AdditionalDisc

		var details []models.SalesOrderReturnDetail
		for _, pd := range u.SalesOrderReturnDetails {
			if pd.Qty < 1 {
				pd.Qty = 1
			}

			details = append(details, pd.Transform())
		}

		p.SalesOrderReturnDetails = details

	}
	return p
}

// SalesOrderReturnDetailRequest : format json request for salesOrder return detail
type SalesOrderReturnDetailRequest struct {
	ID        uint64  `json:"id"`
	Price     float64 `json:"price"`
	Disc      float64 `json:"disc"`
	Qty       uint    `json:"qty"`
	ProductID uint64  `json:"product"`
}

// Transform SalesOrderReturnDetailRequest to SalesOrderReturnDetail
func (u *SalesOrderReturnDetailRequest) Transform() models.SalesOrderReturnDetail {
	var pd models.SalesOrderReturnDetail
	pd.ID = u.ID
	pd.Price = u.Price
	pd.Disc = u.Disc
	pd.Qty = u.Qty
	pd.Product.ID = u.ProductID

	return pd
}
