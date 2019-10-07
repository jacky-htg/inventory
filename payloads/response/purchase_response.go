package response

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// PurchaseResponse : format json response for purchase
type PurchaseResponse struct {
	ID              uint64                   `json:"id"`
	Code            string                   `json:"code"`
	Date            time.Time                `json:"name"`
	Price           float64                  `json:"price"`
	Disc            float64                  `json:"disc"`
	AdditionalDisc  float64                  `json:"additional_disc"`
	Total           float64                  `json:"total"`
	Supplier        SupplierResponse         `json:"supplier"`
	Company         CompanyResponse          `json:"company"`
	Branch          BranchResponse           `json:"branch"`
	PurchaseDetails []PurchaseDetailResponse `json:"purchase_details"`
}

// Transform from Purchase model to Purchase response
func (u *PurchaseResponse) Transform(purchase *models.Purchase) {
	u.ID = purchase.ID
	u.Code = purchase.Code
	u.Date = purchase.Date
	u.Price = purchase.Price
	u.Disc = purchase.Disc
	u.AdditionalDisc = purchase.AdditionalDisc
	u.Total = purchase.Total
	u.Supplier.Transform(&purchase.Supplier)
	u.Company.Transform(&purchase.Company)
	u.Branch.Transform(&purchase.Branch)

	for _, d := range purchase.PurchaseDetails {
		var p PurchaseDetailResponse
		p.Transform(&d)
		u.PurchaseDetails = append(u.PurchaseDetails, p)
	}
}

// PurchaseListResponse : format json response for purchase list
type PurchaseListResponse struct {
	ID             uint64           `json:"id"`
	Code           string           `json:"code"`
	Date           time.Time        `json:"date"`
	Price          float64          `json:"price"`
	Disc           float64          `json:"disc"`
	AdditionalDisc float64          `json:"additional_disc"`
	Total          float64          `json:"total"`
	Supplier       SupplierResponse `json:"supplier"`
	Company        CompanyResponse  `json:"company"`
	Branch         BranchResponse   `json:"branch"`
}

// Transform from Purchase model to Purchase List response
func (u *PurchaseListResponse) Transform(purchase *models.Purchase) {
	u.ID = purchase.ID
	u.Code = purchase.Code
	u.Date = purchase.Date
	u.Price = purchase.Price
	u.Disc = purchase.Disc
	u.AdditionalDisc = purchase.AdditionalDisc
	u.Total = purchase.Total
	u.Supplier.Transform(&purchase.Supplier)
	u.Company.Transform(&purchase.Company)
	u.Branch.Transform(&purchase.Branch)
}

// PurchaseDetailResponse : format json response for purchase detail
type PurchaseDetailResponse struct {
	ID      uint64          `json:"id"`
	Price   float64         `json:"price"`
	Disc    float64         `json:"disc"`
	Qty     uint            `json:"qty"`
	Product ProductResponse `json:"product"`
}

// Transform from PurchaseDetail model to PurchaseDetail response
func (u *PurchaseDetailResponse) Transform(pd *models.PurchaseDetail) {
	u.ID = pd.ID
	u.Price = pd.Price
	u.Disc = pd.Disc
	u.Qty = pd.Qty
	u.Product.Transform(&pd.Product)
}
