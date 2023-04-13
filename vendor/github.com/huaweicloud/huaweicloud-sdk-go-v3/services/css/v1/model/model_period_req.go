package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PeriodReq struct {

	// 订购周期类型。 - 2: 包月。 - 3: 包年。
	PeriodType int32 `json:"periodType"`

	// 订购周期数。 - 若选择包月（参数范围：1-9）。 - 若选择包年（参数范围：1-3）。
	PeriodNum int32 `json:"periodNum"`

	// 是否自动续订，为空时表示不自动续订 - 1: 自动续订。 - 0: 不自动续订（默认）。
	IsAutoRenew *int32 `json:"isAutoRenew,omitempty"`

	//  是否自动支付。下单订购后，是否自动从客户的华为云账户中支付，而不需要客户手动去进行支付。该参数适用于包周期集群。   - 1: 是（会自动选择折扣和优惠券进行优惠，然后自动从客户华为云账户中支付），自动支付失败后会生成订单成功(该订单应付金额是优惠后金额)、但订单状态为“待支付”，等待客户手动支付(手动支付时，客户还可以修改系统自动选择的折扣和优惠券)。   - 0: 否（需要客户手动去支付，客户可以选择折扣和优惠券）。默认值为“0”。
	IsAutoPay *int32 `json:"isAutoPay,omitempty"`

	// 云服务ConsoleURL。 订购订单支付完成后，客户可以通过此URL跳转到云服务Console页面查看信息。（仅手动支付时涉及）。
	ConsoleURL *string `json:"consoleURL,omitempty"`
}

func (o PeriodReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PeriodReq struct{}"
	}

	return strings.Join([]string{"PeriodReq", string(data)}, " ")
}
