package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BindPublicReq struct {
	Eip *BindPublicReqEip `json:"eip"`

	// 是否自动支付。下单订购后，是否自动从客户的华为云账户中支付，而不需要客户手动去进行支付。该参数适用于包周期集群。  - 1: 是（会自动选择折扣和优惠券进行优惠，然后自动从客户华为云账户中支付），自动支付失败后会生成订单成功(该订单应付金额是优惠后金额)、但订单状态为“待支付”，等待客户手动支付(手动支付时，客户还可以修改系统自动选择的折扣和优惠券)。  - 0: 否（需要客户手动去支付，客户可以选择折扣和优惠券）。默认值为“0”。
	IsAutoPay *int32 `json:"isAutoPay,omitempty"`
}

func (o BindPublicReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BindPublicReq struct{}"
	}

	return strings.Join([]string{"BindPublicReq", string(data)}, " ")
}
