package request

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// NewSalesOrderRequest : format json request for new sales order
type NewSalesOrderRequest struct {
	Date              string                       `json:"date" validate:"required"`
	AdditionalDisc    float64                      `json:"additional_disc"`
	SalesOrderDetails []NewSalesOrderDetailRequest `json:"sales_order_details" validate:"required"`
	SalesmanID        uint64                       `json:"salesman" validate:"required"`
	CustomerID        uint64                       `json:"customer" validate:"required"`
}

// Transform NewSalesOrderRequest to SalesOrder
func (u *NewSalesOrderRequest) Transform() *models.SalesOrder {
	var p models.SalesOrder
	p.Date, _ = time.Parse("2006-01-02", u.Date)
	p.Salesman.ID = u.SalesmanID
	p.Customer.ID = u.CustomerID
	p.AdditionalDisc = u.AdditionalDisc

	for _, pd := range u.SalesOrderDetails {
		if pd.Qty < 1 {
			pd.Qty = 1
		}

		p.SalesOrderDetails = append(p.SalesOrderDetails, pd.Transform())
	}

	return &p
}

// NewSalesOrderDetailRequest : format json request for sales order detail
type NewSalesOrderDetailRequest struct {
	Price     float64 `json:"price" validate:"required"`
	Disc      float64 `json:"disc"`
	Qty       uint    `json:"qty" validate:"required"`
	ProductID uint64  `json:"product" validate:"required"`
}

// Transform NewSalesOrderDetailRequest to SalesOrderDetail
func (u *NewSalesOrderDetailRequest) Transform() models.SalesOrderDetail {
	var pd models.SalesOrderDetail
	pd.Price = u.Price
	pd.Disc = u.Disc
	pd.Qty = u.Qty
	pd.Product.ID = u.ProductID

	return pd
}

// SalesOrderRequest : format json request for sales order
type SalesOrderRequest struct {
	ID                uint64                    `json:"id" validate:"required"`
	Date              string                    `json:"date"`
	AdditionalDisc    float64                   `json:"additional_disc"`
	SalesOrderDetails []SalesOrderDetailRequest `json:"sales_order_details"`
	SalesmanID        uint64                    `json:"salesman"`
	CustomerID        uint64                    `json:"customer"`
}

// Transform SalesOrderRequest to SalesOrder
func (u *SalesOrderRequest) Transform(p *models.SalesOrder) *models.SalesOrder {
	if u.ID == p.ID {
		p.Date, _ = time.Parse("2006-01-02", u.Date)
		p.Salesman.ID = u.SalesmanID
		p.AdditionalDisc = u.AdditionalDisc

		var details []models.SalesOrderDetail
		for _, pd := range u.SalesOrderDetails {
			if pd.Qty < 1 {
				pd.Qty = 1
			}

			details = append(details, pd.Transform())
		}

		p.SalesOrderDetails = details

	}
	return p
}

// SalesOrderDetailRequest : format json request for sales order detail
type SalesOrderDetailRequest struct {
	ID        uint64  `json:"id"`
	Price     float64 `json:"price"`
	Disc      float64 `json:"disc"`
	Qty       uint    `json:"qty"`
	ProductID uint64  `json:"product"`
}

// Transform SalesOrderDetailRequest to SalesOrderDetail
func (u *SalesOrderDetailRequest) Transform() models.SalesOrderDetail {
	var pd models.SalesOrderDetail
	pd.ID = u.ID
	pd.Price = u.Price
	pd.Disc = u.Disc
	pd.Qty = u.Qty
	pd.Product.ID = u.ProductID

	return pd
}
