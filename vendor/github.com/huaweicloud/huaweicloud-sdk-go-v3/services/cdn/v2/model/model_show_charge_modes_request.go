package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowChargeModesRequest Request Object
type ShowChargeModesRequest struct {

	// 加速类型，base（基础加速）
	ProductType string `json:"product_type"`

	// 查询计费模式状态，active（已生效），upcoming（待生效），不传默认为active(已生效)
	Status *string `json:"status,omitempty"`

	// 服务区域，mainland_china（国内），outside_mainland_china（海外），不传默认为mainland_china(国内)
	ServiceArea *string `json:"service_area,omitempty"`
}

func (o ShowChargeModesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowChargeModesRequest struct{}"
	}

	return strings.Join([]string{"ShowChargeModesRequest", string(data)}, " ")
}
