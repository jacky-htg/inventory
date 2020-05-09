package response

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// SalesOrderReturnResponse : format json response for salesOrder return
type SalesOrderReturnResponse struct {
	ID                      uint64                           `json:"id"`
	Code                    string                           `json:"code"`
	Date                    time.Time                        `json:"name"`
	Price                   float64                          `json:"price"`
	Disc                    float64                          `json:"disc"`
	AdditionalDisc          float64                          `json:"additional_disc"`
	Total                   float64                          `json:"total"`
	SalesOrder              SalesOrderResponse               `json:"salesOrder"`
	Company                 CompanyResponse                  `json:"company"`
	Branch                  BranchResponse                   `json:"branch"`
	SalesOrderReturnDetails []SalesOrderReturnDetailResponse `json:"sales_order_return_details"`
}

// Transform from SalesOrderReturn model to SalesOrder Return response
func (u *SalesOrderReturnResponse) Transform(salesOrderReturn *models.SalesOrderReturn) {
	u.ID = salesOrderReturn.ID
	u.Code = salesOrderReturn.Code
	u.Date = salesOrderReturn.Date
	u.Price = salesOrderReturn.Price
	u.Disc = salesOrderReturn.Disc
	u.AdditionalDisc = salesOrderReturn.AdditionalDisc
	u.Total = salesOrderReturn.Total
	u.SalesOrder.Transform(&salesOrderReturn.SalesOrder)
	u.Company.Transform(&salesOrderReturn.Company)
	u.Branch.Transform(&salesOrderReturn.Branch)

	for _, d := range salesOrderReturn.SalesOrderReturnDetails {
		var p SalesOrderReturnDetailResponse
		p.Transform(&d)
		u.SalesOrderReturnDetails = append(u.SalesOrderReturnDetails, p)
	}
}

// SalesOrderReturnListResponse : format json response for salesOrder return list
type SalesOrderReturnListResponse struct {
	ID             uint64             `json:"id"`
	Code           string             `json:"code"`
	Date           time.Time          `json:"date"`
	Price          float64            `json:"price"`
	Disc           float64            `json:"disc"`
	AdditionalDisc float64            `json:"additional_disc"`
	Total          float64            `json:"total"`
	SalesOrder     SalesOrderResponse `json:"sales_order"`
	Company        CompanyResponse    `json:"company"`
	Branch         BranchResponse     `json:"branch"`
}

// Transform from SalesOrderReturn model to SalesOrderReturn List response
func (u *SalesOrderReturnListResponse) Transform(salesOrderReturn *models.SalesOrderReturn) {
	u.ID = salesOrderReturn.ID
	u.Code = salesOrderReturn.Code
	u.Date = salesOrderReturn.Date
	u.Price = salesOrderReturn.Price
	u.Disc = salesOrderReturn.Disc
	u.AdditionalDisc = salesOrderReturn.AdditionalDisc
	u.Total = salesOrderReturn.Total
	u.SalesOrder.Transform(&salesOrderReturn.SalesOrder)
	u.Company.Transform(&salesOrderReturn.Company)
	u.Branch.Transform(&salesOrderReturn.Branch)
}

// SalesOrderReturnDetailResponse : format json response for salesOrder return detail
type SalesOrderReturnDetailResponse struct {
	ID      uint64          `json:"id"`
	Price   float64         `json:"price"`
	Disc    float64         `json:"disc"`
	Qty     uint            `json:"qty"`
	Product ProductResponse `json:"product"`
}

// Transform from SalesOrderReturnDetail model to SalesOrderReturnDetail response
func (u *SalesOrderReturnDetailResponse) Transform(pd *models.SalesOrderReturnDetail) {
	u.ID = pd.ID
	u.Price = pd.Price
	u.Disc = pd.Disc
	u.Qty = pd.Qty
	u.Product.Transform(&pd.Product)
}
