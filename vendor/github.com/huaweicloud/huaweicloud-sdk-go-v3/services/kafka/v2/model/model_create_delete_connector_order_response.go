package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateDeleteConnectorOrderResponse Response Object
type CreateDeleteConnectorOrderResponse struct {

	// 返回cbc生成的订单id。
	OrderId        *string `json:"order_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateDeleteConnectorOrderResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDeleteConnectorOrderResponse struct{}"
	}

	return strings.Join([]string{"CreateDeleteConnectorOrderResponse", string(data)}, " ")
}
