package response

import (
	"time"

	"github.com/jacky-htg/inventory/models"
)

// DeliveryResponse : format json response for Delivery
type DeliveryResponse struct {
	ID              uint64                   `json:"id"`
	Code            string                   `json:"code"`
	Date            time.Time                `json:"name"`
	Remark          string                   `json:"remark"`
	SalesOrder      SalesOrderResponse       `json:"sales_order"`
	Company         CompanyResponse          `json:"company"`
	Branch          BranchResponse           `json:"branch"`
	DeliveryDetails []DeliveryDetailResponse `json:"delivery_details"`
}

// Transform from Delivery model to Delivery response
func (u *DeliveryResponse) Transform(delivery *models.Delivery) {
	u.ID = delivery.ID
	u.Code = delivery.Code
	u.Date = delivery.Date
	u.Remark = delivery.Remark
	u.SalesOrder.Transform(&delivery.SalesOrder)
	u.Company.Transform(&delivery.Company)
	u.Branch.Transform(&delivery.Branch)

	for _, d := range delivery.DeliveryDetails {
		var p DeliveryDetailResponse
		p.Transform(&d)
		u.DeliveryDetails = append(u.DeliveryDetails, p)
	}
}

// DeliveryListResponse : format json response for Delivery list
type DeliveryListResponse struct {
	ID         uint64             `json:"id"`
	Code       string             `json:"code"`
	Date       time.Time          `json:"date"`
	Remark     string             `json:"remark"`
	SalesOrder SalesOrderResponse `json:"sales_order"`
	Company    CompanyResponse    `json:"company"`
	Branch     BranchResponse     `json:"branch"`
}

// Transform from Delivery model to Delivery List response
func (u *DeliveryListResponse) Transform(delivery *models.Delivery) {
	u.ID = delivery.ID
	u.Code = delivery.Code
	u.Date = delivery.Date
	u.Remark = delivery.Remark
	u.SalesOrder.Transform(&delivery.SalesOrder)
	u.Company.Transform(&delivery.Company)
	u.Branch.Transform(&delivery.Branch)
}

// DeliveryDetailResponse : format json response for Delivery detail
type DeliveryDetailResponse struct {
	ID      uint64          `json:"id"`
	Qty     uint            `json:"qty"`
	Product ProductResponse `json:"product"`
	Code    string          `json:"code"`
	Shelve  ShelveResponse  `json:"shelve"`
}

// Transform from DeliveryDetail model to DeliveryDetail response
func (u *DeliveryDetailResponse) Transform(pd *models.DeliveryDetail) {
	u.ID = pd.ID
	u.Qty = pd.Qty
	u.Product.Transform(&pd.Product)
	u.Code = pd.Code
	u.Shelve.Transform(&pd.Shelve)
}
