package response

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// DeliveryReturnResponse : format json response for Delivery Return
type DeliveryReturnResponse struct {
	ID                    uint64                         `json:"id"`
	Code                  string                         `json:"code"`
	Date                  time.Time                      `json:"name"`
	Remark                string                         `json:"remark"`
	Delivery              DeliveryResponse               `json:"delivery"`
	Company               CompanyResponse                `json:"company"`
	Branch                BranchResponse                 `json:"branch"`
	DeliveryReturnDetails []DeliveryReturnDetailResponse `json:"delivery_return_details"`
}

// Transform from Delivery Return model to Delivery return response
func (u *DeliveryReturnResponse) Transform(deliveryReturn *models.DeliveryReturn) {
	u.ID = deliveryReturn.ID
	u.Code = deliveryReturn.Code
	u.Date = deliveryReturn.Date
	u.Remark = deliveryReturn.Remark
	u.Delivery.Transform(&deliveryReturn.Delivery)
	u.Company.Transform(&deliveryReturn.Company)
	u.Branch.Transform(&deliveryReturn.Branch)

	for _, d := range deliveryReturn.DeliveryReturnDetails {
		var p DeliveryReturnDetailResponse
		p.Transform(&d)
		u.DeliveryReturnDetails = append(u.DeliveryReturnDetails, p)
	}
}

// DeliveryReturnListResponse : format json response for Delivery Return list
type DeliveryReturnListResponse struct {
	ID       uint64           `json:"id"`
	Code     string           `json:"code"`
	Date     time.Time        `json:"date"`
	Remark   string           `json:"remark"`
	Delivery DeliveryResponse `json:"delivery"`
	Company  CompanyResponse  `json:"company"`
	Branch   BranchResponse   `json:"branch"`
}

// Transform from Delivery Return model to Delivery Return List response
func (u *DeliveryReturnListResponse) Transform(deliveryReturn *models.DeliveryReturn) {
	u.ID = deliveryReturn.ID
	u.Code = deliveryReturn.Code
	u.Date = deliveryReturn.Date
	u.Remark = deliveryReturn.Remark
	u.Delivery.Transform(&deliveryReturn.Delivery)
	u.Company.Transform(&deliveryReturn.Company)
	u.Branch.Transform(&deliveryReturn.Branch)
}

// DeliveryReturnDetailResponse : format json response for Delivery Return detail
type DeliveryReturnDetailResponse struct {
	ID      uint64          `json:"id"`
	Qty     uint            `json:"qty"`
	Product ProductResponse `json:"product"`
	Code    string          `json:"code"`
}

// Transform from DeliveryReturnDetail model to DeliveryReturnDetail response
func (u *DeliveryReturnDetailResponse) Transform(pd *models.DeliveryReturnDetail) {
	u.ID = pd.ID
	u.Qty = pd.Qty
	u.Product.Transform(&pd.Product)
	u.Code = pd.Code
}
