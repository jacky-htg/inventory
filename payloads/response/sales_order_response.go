package response

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// SalesOrderResponse : format json response for sales order
type SalesOrderResponse struct {
	ID                uint64                     `json:"id"`
	Code              string                     `json:"code"`
	Date              time.Time                  `json:"name"`
	Price             float64                    `json:"price"`
	Disc              float64                    `json:"disc"`
	AdditionalDisc    float64                    `json:"additional_disc"`
	Total             float64                    `json:"total"`
	Salesman          SalesmanResponse           `json:"salesman"`
	Customer          CustomerResponse           `json:"customer"`
	Company           CompanyResponse            `json:"company"`
	Branch            BranchResponse             `json:"branch"`
	SalesOrderDetails []SalesOrderDetailResponse `json:"sales_order_details"`
}

// Transform from SalesOrder model to SalesOrderResponse
func (u *SalesOrderResponse) Transform(salesOrder *models.SalesOrder) {
	u.ID = salesOrder.ID
	u.Code = salesOrder.Code
	u.Date = salesOrder.Date
	u.Price = salesOrder.Price
	u.Disc = salesOrder.Disc
	u.AdditionalDisc = salesOrder.AdditionalDisc
	u.Total = salesOrder.Total
	u.Salesman.Transform(&salesOrder.Salesman)
	u.Customer.Transform(&salesOrder.Customer)
	u.Company.Transform(&salesOrder.Company)
	u.Branch.Transform(&salesOrder.Branch)

	for _, d := range salesOrder.SalesOrderDetails {
		var s SalesOrderDetailResponse
		s.Transform(&d)
		u.SalesOrderDetails = append(u.SalesOrderDetails, s)
	}
}

// SalesOrderListResponse : format json response for sales order list
type SalesOrderListResponse struct {
	ID             uint64           `json:"id"`
	Code           string           `json:"code"`
	Date           time.Time        `json:"date"`
	Price          float64          `json:"price"`
	Disc           float64          `json:"disc"`
	AdditionalDisc float64          `json:"additional_disc"`
	Total          float64          `json:"total"`
	Salesman       SalesmanResponse `json:"salesman"`
	Customer       CustomerResponse `json:"customer"`
	Company        CompanyResponse  `json:"company"`
	Branch         BranchResponse   `json:"branch"`
}

// Transform from SalesOrder model to Sales Order List response
func (u *SalesOrderListResponse) Transform(salesOrder *models.SalesOrder) {
	u.ID = salesOrder.ID
	u.Code = salesOrder.Code
	u.Date = salesOrder.Date
	u.Price = salesOrder.Price
	u.Disc = salesOrder.Disc
	u.AdditionalDisc = salesOrder.AdditionalDisc
	u.Total = salesOrder.Total
	u.Salesman.Transform(&salesOrder.Salesman)
	u.Customer.Transform(&salesOrder.Customer)
	u.Company.Transform(&salesOrder.Company)
	u.Branch.Transform(&salesOrder.Branch)
}

// SalesOrderDetailResponse : format json response for sales order detail
type SalesOrderDetailResponse struct {
	ID      uint64          `json:"id"`
	Price   float64         `json:"price"`
	Disc    float64         `json:"disc"`
	Qty     uint            `json:"qty"`
	Product ProductResponse `json:"product"`
}

// Transform from SalesOrderDetail model to SalesOrderDetailResponse
func (u *SalesOrderDetailResponse) Transform(sod *models.SalesOrderDetail) {
	u.ID = sod.ID
	u.Price = sod.Price
	u.Disc = sod.Disc
	u.Qty = sod.Qty
	u.Product.Transform(&sod.Product)
}
