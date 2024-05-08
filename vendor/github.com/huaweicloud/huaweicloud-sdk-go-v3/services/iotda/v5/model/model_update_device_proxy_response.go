package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDeviceProxyResponse Response Object
type UpdateDeviceProxyResponse struct {

	// **参数说明**：代理ID。用来唯一标识一个代理规则
	ProxyId *string `json:"proxy_id,omitempty"`

	// **参数说明**：设备代理名称
	ProxyName *string `json:"proxy_name,omitempty"`

	// **参数说明**：代理设备组，组内所有设备共享网关权限，即组内任意一个网关下的子设备可以通过组里任意一个网关上线然后进行数据上报。
	ProxyDevices *[]string `json:"proxy_devices,omitempty"`

	EffectiveTimeRange *EffectiveTimeRangeResponseDto `json:"effective_time_range,omitempty"`

	// **参数说明**：资源空间ID。
	AppId          *string `json:"app_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateDeviceProxyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDeviceProxyResponse struct{}"
	}

	return strings.Join([]string{"UpdateDeviceProxyResponse", string(data)}, " ")
}
