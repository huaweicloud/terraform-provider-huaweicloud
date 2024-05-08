package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateQuotasOrderRequestInfo 创建订单订购配额请求体对象
type CreateQuotasOrderRequestInfo struct {

	// 资源规格   - hss.version.basic ：基础版。   - hss.version.advanced ：专业版。   - hss.version.enterprise ：企业版。   - hss.version.premium ：旗舰版。   - hss.version.wtp ：网页防篡改版。   - hss.version.container.enterprise：容器版。
	ResourceSpecCode string `json:"resource_spec_code"`

	// 订购周期类型   - 2 : 月   - 3 : 年
	PeriodType int32 `json:"period_type"`

	// 订购周期数
	PeriodNum int32 `json:"period_num"`

	// 是否支持自动续订，true表示自动续订，false表示不自动续订，默认值为false
	IsAutoRenew *bool `json:"is_auto_renew,omitempty"`

	// 是否支持自动支付，true表示支持，false表示不支持，默认值为false
	IsAutoPay *bool `json:"is_auto_pay,omitempty"`

	// 订购数量
	SubscriptionNum int32 `json:"subscription_num"`
}

func (o CreateQuotasOrderRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateQuotasOrderRequestInfo struct{}"
	}

	return strings.Join([]string{"CreateQuotasOrderRequestInfo", string(data)}, " ")
}
