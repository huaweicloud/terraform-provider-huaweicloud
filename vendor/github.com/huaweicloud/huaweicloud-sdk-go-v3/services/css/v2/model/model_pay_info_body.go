package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 包周期集群相关参数订购包周期集群时使用。
type PayInfoBody struct {

	// 订购周期类型。 - 2: 包月。 - 3: 包年。
	PayModel int32 `json:"payModel"`

	// 订购周期数。 - 若payModel为2，则有效值为1-9。 - 若payModel为3，则有效值为1-3。
	Period int32 `json:"period"`

	// 是否自动续订，为空时表示不自动续订。 - 1: 自动续订。 - 2：不自动续订（默认）。
	IsAutoRenew *int32 `json:"isAutoRenew,omitempty"`

	// 是否自动支付：下单订购后，是否自动从客户的华为云账户中支付，而不需要客户手动去进行支付。  - 1: 是（会自动选择折扣和优惠券进行优惠，然后自动从客户华为云账户中支付），自动支付失败后会生成订单成功(该订单应付金额是优惠后金额)、但订单状态为“待支付”，等待客户手动支付(手动支付时，客户还可以修改系统自动选择的折扣和优惠券)。  - 2: 否（需要客户手动去支付，客户可以选择折扣和优惠券）。默认值为“0”。
	IsAutoPay *int32 `json:"isAutoPay,omitempty"`
}

func (o PayInfoBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PayInfoBody struct{}"
	}

	return strings.Join([]string{"PayInfoBody", string(data)}, " ")
}
