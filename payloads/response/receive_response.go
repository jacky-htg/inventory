package response

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// ReceiveResponse : format json response for Receive
type ReceiveResponse struct {
	ID             uint64                  `json:"id"`
	Code           string                  `json:"code"`
	Date           time.Time               `json:"name"`
	Remark         string                  `json:"remark"`
	Purchase       PurchaseResponse        `json:"purchase"`
	Company        CompanyResponse         `json:"company"`
	Branch         BranchResponse          `json:"branch"`
	ReceiveDetails []ReceiveDetailResponse `json:"receive_details"`
}

// Transform from Receive model to Receive response
func (u *ReceiveResponse) Transform(Receive *models.Receive) {
	u.ID = Receive.ID
	u.Code = Receive.Code
	u.Date = Receive.Date
	u.Remark = Receive.Remark
	u.Purchase.Transform(&Receive.Purchase)
	u.Company.Transform(&Receive.Company)
	u.Branch.Transform(&Receive.Branch)

	for _, d := range Receive.ReceiveDetails {
		var p ReceiveDetailResponse
		p.Transform(&d)
		u.ReceiveDetails = append(u.ReceiveDetails, p)
	}
}

// ReceiveListResponse : format json response for Receive list
type ReceiveListResponse struct {
	ID       uint64           `json:"id"`
	Code     string           `json:"code"`
	Date     time.Time        `json:"date"`
	Remark   string           `json:"remark"`
	Purchase PurchaseResponse `json:"purchase"`
	Company  CompanyResponse  `json:"company"`
	Branch   BranchResponse   `json:"branch"`
}

// Transform from Receive model to Receive List response
func (u *ReceiveListResponse) Transform(Receive *models.Receive) {
	u.ID = Receive.ID
	u.Code = Receive.Code
	u.Date = Receive.Date
	u.Remark = Receive.Remark
	u.Purchase.Transform(&Receive.Purchase)
	u.Company.Transform(&Receive.Company)
	u.Branch.Transform(&Receive.Branch)
}

// ReceiveDetailResponse : format json response for Receive detail
type ReceiveDetailResponse struct {
	ID      uint64          `json:"id"`
	Qty     uint            `json:"qty"`
	Product ProductResponse `json:"product"`
	Code    string          `json:"code"`
	Shelve  ShelveResponse  `json:"shelve"`
}

// Transform from ReceiveDetail model to ReceiveDetail response
func (u *ReceiveDetailResponse) Transform(pd *models.ReceiveDetail) {
	u.ID = pd.ID
	u.Qty = pd.Qty
	u.Product.Transform(&pd.Product)
	u.Code = pd.Code
	u.Shelve.Transform(&pd.Shelve)
}
