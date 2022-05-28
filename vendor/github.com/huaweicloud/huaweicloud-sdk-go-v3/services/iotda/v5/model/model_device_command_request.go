package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeviceCommandRequest struct {

	// **参数说明**：设备命令所属的设备服务ID，在设备关联的产品模型中定义。 **取值范围**：长度不超过64的字符串。
	ServiceId *string `json:"service_id,omitempty"`

	// **参数说明**：设备命令名称，在设备关联的产品模型中定义。 **取值范围**：长度不超过128的字符串。
	CommandName *string `json:"command_name,omitempty"`

	// **参数说明**：设备执行的命令，Json格式，里面是一个个键值对，如果serviceId不为空，每个键都是profile中命令的参数名（paraName）;如果serviceId为空则由用户自定义命令格式。设备命令示例：{\"value\":\"1\"}，具体格式需要应用和设备约定。此参数仅支持Json格式，暂不支持字符串。
	Paras *interface{} `json:"paras"`
}

func (o DeviceCommandRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeviceCommandRequest struct{}"
	}

	return strings.Join([]string{"DeviceCommandRequest", string(data)}, " ")
}
