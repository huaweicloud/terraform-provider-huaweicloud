package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 设备影子数据结构体。
type DeviceShadowData struct {

	// 设备的服务ID，在设备关联的产品模型中定义。
	ServiceId string `json:"service_id"`

	Desired *DeviceShadowProperties `json:"desired,omitempty"`

	Reported *DeviceShadowProperties `json:"reported,omitempty"`

	// 设备影子的版本，携带该参数时平台会校验值必须等于当前影子版本，初始从0开始。
	Version *int64 `json:"version,omitempty"`
}

func (o DeviceShadowData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeviceShadowData struct{}"
	}

	return strings.Join([]string{"DeviceShadowData", string(data)}, " ")
}
