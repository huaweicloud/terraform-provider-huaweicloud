package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetChargeModesBody 设置计费模式请求体
type SetChargeModesBody struct {

	// 计费模式，支持flux（流量），v2及以上客户支持bw（带宽）
	ChargeMode string `json:"charge_mode"`

	// 产品模式，仅支持base（基础加速）
	ProductType string `json:"product_type"`

	// 服务区域，仅支持mainland_china（国内）
	ServiceArea string `json:"service_area"`
}

func (o SetChargeModesBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetChargeModesBody struct{}"
	}

	return strings.Join([]string{"SetChargeModesBody", string(data)}, " ")
}
