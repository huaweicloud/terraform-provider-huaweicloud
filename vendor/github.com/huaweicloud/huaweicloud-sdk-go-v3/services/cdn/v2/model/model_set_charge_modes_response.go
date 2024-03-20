package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetChargeModesResponse Response Object
type SetChargeModesResponse struct {

	// 账号的计费模式
	ChargeMode *string `json:"charge_mode,omitempty"`

	// 加速类型
	ProductType *string `json:"product_type,omitempty"`

	// 该模式生效时间
	EffectiveTime *int64 `json:"effective_time,omitempty"`

	// 创建时间
	CreateTime *int64 `json:"create_time,omitempty"`

	// 该模式的区域
	ServiceArea *string `json:"service_area,omitempty"`

	// 状态,首次开通状态为active,之后修改为upcoming
	Status         *string `json:"status,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o SetChargeModesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetChargeModesResponse struct{}"
	}

	return strings.Join([]string{"SetChargeModesResponse", string(data)}, " ")
}
