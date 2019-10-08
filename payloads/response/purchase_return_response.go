package response

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// PurchaseReturnResponse : format json response for purchase return
type PurchaseReturnResponse struct {
	ID                    uint64                         `json:"id"`
	Code                  string                         `json:"code"`
	Date                  time.Time                      `json:"name"`
	Price                 float64                        `json:"price"`
	Disc                  float64                        `json:"disc"`
	AdditionalDisc        float64                        `json:"additional_disc"`
	Total                 float64                        `json:"total"`
	Purchase              PurchaseResponse               `json:"purchase"`
	Company               CompanyResponse                `json:"company"`
	Branch                BranchResponse                 `json:"branch"`
	PurchaseReturnDetails []PurchaseReturnDetailResponse `json:"purchase_return_details"`
}

// Transform from PurchaseReturn model to Purchase Return response
func (u *PurchaseReturnResponse) Transform(purchaseReturn *models.PurchaseReturn) {
	u.ID = purchaseReturn.ID
	u.Code = purchaseReturn.Code
	u.Date = purchaseReturn.Date
	u.Price = purchaseReturn.Price
	u.Disc = purchaseReturn.Disc
	u.AdditionalDisc = purchaseReturn.AdditionalDisc
	u.Total = purchaseReturn.Total
	u.Purchase.Transform(&purchaseReturn.Purchase)
	u.Company.Transform(&purchaseReturn.Company)
	u.Branch.Transform(&purchaseReturn.Branch)

	for _, d := range purchaseReturn.PurchaseReturnDetails {
		var p PurchaseReturnDetailResponse
		p.Transform(&d)
		u.PurchaseReturnDetails = append(u.PurchaseReturnDetails, p)
	}
}

// PurchaseReturnListResponse : format json response for purchase return list
type PurchaseReturnListResponse struct {
	ID             uint64           `json:"id"`
	Code           string           `json:"code"`
	Date           time.Time        `json:"date"`
	Price          float64          `json:"price"`
	Disc           float64          `json:"disc"`
	AdditionalDisc float64          `json:"additional_disc"`
	Total          float64          `json:"total"`
	Purchase       PurchaseResponse `json:"purchase"`
	Company        CompanyResponse  `json:"company"`
	Branch         BranchResponse   `json:"branch"`
}

// Transform from PurchaseReturn model to PurchaseReturn List response
func (u *PurchaseReturnListResponse) Transform(purchaseReturn *models.PurchaseReturn) {
	u.ID = purchaseReturn.ID
	u.Code = purchaseReturn.Code
	u.Date = purchaseReturn.Date
	u.Price = purchaseReturn.Price
	u.Disc = purchaseReturn.Disc
	u.AdditionalDisc = purchaseReturn.AdditionalDisc
	u.Total = purchaseReturn.Total
	u.Purchase.Transform(&purchaseReturn.Purchase)
	u.Company.Transform(&purchaseReturn.Company)
	u.Branch.Transform(&purchaseReturn.Branch)
}

// PurchaseReturnDetailResponse : format json response for purchase return detail
type PurchaseReturnDetailResponse struct {
	ID      uint64          `json:"id"`
	Price   float64         `json:"price"`
	Disc    float64         `json:"disc"`
	Qty     uint            `json:"qty"`
	Product ProductResponse `json:"product"`
}

// Transform from PurchaseReturnDetail model to PurchaseReturnDetail response
func (u *PurchaseReturnDetailResponse) Transform(pd *models.PurchaseReturnDetail) {
	u.ID = pd.ID
	u.Price = pd.Price
	u.Disc = pd.Disc
	u.Qty = pd.Qty
	u.Product.Transform(&pd.Product)
}
