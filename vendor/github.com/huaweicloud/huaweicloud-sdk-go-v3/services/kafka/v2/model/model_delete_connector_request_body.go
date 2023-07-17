package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeleteConnectorRequestBody struct {

	// cbc生成实例变更的订单 按需实例不传入此参数。 包周期实例传入删除connector时生成的订单，由cbc调用时传入。
	OrderId *string `json:"order_id,omitempty"`

	// 包周期实例变更时的product id 按需实例不传入此参数。 包周期实例传入对变更实例规格的product，由cbc调用时传入。
	CsbInstanceProductId *string `json:"csb_instance_product_id,omitempty"`
}

func (o DeleteConnectorRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteConnectorRequestBody struct{}"
	}

	return strings.Join([]string{"DeleteConnectorRequestBody", string(data)}, " ")
}
