package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// QueryDeviceProxySimplify 设备代理的基本信息。
type QueryDeviceProxySimplify struct {

	// **参数说明**：代理ID。用来唯一标识一个代理规则
	ProxyId *string `json:"proxy_id,omitempty"`

	// **参数说明**：设备代理名称
	ProxyName *string `json:"proxy_name,omitempty"`

	EffectiveTimeRange *EffectiveTimeRangeResponseDto `json:"effective_time_range,omitempty"`

	// **参数说明**：资源空间ID。
	AppId *string `json:"app_id,omitempty"`
}

func (o QueryDeviceProxySimplify) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QueryDeviceProxySimplify struct{}"
	}

	return strings.Join([]string{"QueryDeviceProxySimplify", string(data)}, " ")
}
