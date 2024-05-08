package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDeviceProxy 添加设备代理结构体。
type UpdateDeviceProxy struct {

	// **参数说明**：设备代理名称
	ProxyName *string `json:"proxy_name,omitempty"`

	// **参数说明**：代理设备列表，列表内所有设备共享网关权限，即列表内任意一个网关下的子设备可以通过组里任意一个网关上线然后进行数据上报。 **取值范围**：列表内填写设备id，列表内最少有2个设备id，最多有10个设备id，设备id取值范围：长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合，建议不少于4个字符。
	ProxyDevices *[]string `json:"proxy_devices,omitempty"`

	EffectiveTimeRange *EffectiveTimeRange `json:"effective_time_range,omitempty"`
}

func (o UpdateDeviceProxy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDeviceProxy struct{}"
	}

	return strings.Join([]string{"UpdateDeviceProxy", string(data)}, " ")
}
