package response

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// ReceiveReturnResponse : format json response for Receive Return
type ReceiveReturnResponse struct {
	ID                   uint64                        `json:"id"`
	Code                 string                        `json:"code"`
	Date                 time.Time                     `json:"name"`
	Remark               string                        `json:"remark"`
	Receive              ReceiveResponse               `json:"receive"`
	Company              CompanyResponse               `json:"company"`
	Branch               BranchResponse                `json:"branch"`
	ReceiveReturnDetails []ReceiveReturnDetailResponse `json:"receive_return_details"`
}

// Transform from Receive Return model to Receive return response
func (u *ReceiveReturnResponse) Transform(receiveReturn *models.ReceiveReturn) {
	u.ID = receiveReturn.ID
	u.Code = receiveReturn.Code
	u.Date = receiveReturn.Date
	u.Remark = receiveReturn.Remark
	u.Receive.Transform(&receiveReturn.Receive)
	u.Company.Transform(&receiveReturn.Company)
	u.Branch.Transform(&receiveReturn.Branch)

	for _, d := range receiveReturn.ReceiveReturnDetails {
		var p ReceiveReturnDetailResponse
		p.Transform(&d)
		u.ReceiveReturnDetails = append(u.ReceiveReturnDetails, p)
	}
}

// ReceiveReturnListResponse : format json response for Receive Return list
type ReceiveReturnListResponse struct {
	ID      uint64          `json:"id"`
	Code    string          `json:"code"`
	Date    time.Time       `json:"date"`
	Remark  string          `json:"remark"`
	Receive ReceiveResponse `json:"receive"`
	Company CompanyResponse `json:"company"`
	Branch  BranchResponse  `json:"branch"`
}

// Transform from Receive Return model to Receive Return List response
func (u *ReceiveReturnListResponse) Transform(receiveReturn *models.ReceiveReturn) {
	u.ID = receiveReturn.ID
	u.Code = receiveReturn.Code
	u.Date = receiveReturn.Date
	u.Remark = receiveReturn.Remark
	u.Receive.Transform(&receiveReturn.Receive)
	u.Company.Transform(&receiveReturn.Company)
	u.Branch.Transform(&receiveReturn.Branch)
}

// ReceiveReturnDetailResponse : format json response for Receive Return detail
type ReceiveReturnDetailResponse struct {
	ID      uint64          `json:"id"`
	Qty     uint            `json:"qty"`
	Product ProductResponse `json:"product"`
	Code    string          `json:"code"`
}

// Transform from ReceiveReturnDetail model to ReceiveReturnDetail response
func (u *ReceiveReturnDetailResponse) Transform(pd *models.ReceiveReturnDetail) {
	u.ID = pd.ID
	u.Qty = pd.Qty
	u.Product.Transform(&pd.Product)
	u.Code = pd.Code
}
