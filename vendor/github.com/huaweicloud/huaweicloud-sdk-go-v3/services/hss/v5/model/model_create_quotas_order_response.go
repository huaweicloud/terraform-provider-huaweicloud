package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateQuotasOrderResponse Response Object
type CreateQuotasOrderResponse struct {

	// 订单ID
	OrderId        *string `json:"order_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateQuotasOrderResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateQuotasOrderResponse struct{}"
	}

	return strings.Join([]string{"CreateQuotasOrderResponse", string(data)}, " ")
}
