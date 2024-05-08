package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddDeviceProxy 添加设备代理结构体。
type AddDeviceProxy struct {

	// **参数说明**：设备代理名称
	ProxyName string `json:"proxy_name"`

	// **参数说明**：代理设备列表，列表内所有设备共享网关权限，即列表内任意一个网关下的子设备可以通过组里任意一个网关上线然后进行数据上报。 **取值范围**：列表内填写设备id，列表内最少有2个设备id，最多有10个设备id，设备id取值范围：长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合，建议不少于4个字符。
	ProxyDevices []string `json:"proxy_devices"`

	EffectiveTimeRange *EffectiveTimeRange `json:"effective_time_range"`

	// **参数说明**：资源空间ID。携带该参数指定创建的设备归属到哪个资源空间下。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	AppId string `json:"app_id"`
}

func (o AddDeviceProxy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddDeviceProxy struct{}"
	}

	return strings.Join([]string{"AddDeviceProxy", string(data)}, " ")
}
